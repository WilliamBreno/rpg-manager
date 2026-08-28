package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/repository"

	"gorm.io/gorm"
)

type CharacterService struct {
	Repo        *repository.CharacterRepository
	SkillRepo   *repository.SkillRepository
	ClassRepo   *repository.ClassRepository
	RaceRepo    *repository.RaceRepository
	TalentoRepo *repository.TalentoRepository
	DB          *gorm.DB
}

func NewCharacterService(repo *repository.CharacterRepository, skillRepo *repository.SkillRepository, classRepo *repository.ClassRepository, raceRepo *repository.RaceRepository, talentoRepo *repository.TalentoRepository, db *gorm.DB) *CharacterService {
	return &CharacterService{Repo: repo, SkillRepo: skillRepo, ClassRepo: classRepo, RaceRepo: raceRepo, TalentoRepo: talentoRepo, DB: db}
}

// ── Tabelas de XP ────────────────────────────────────────────────────────────

// xpTable4e: XP acumulado necessário para atingir cada nível (1–30)
var xpTable4e = map[int]int{
	1: 0, 2: 1000, 3: 2250, 4: 3750, 5: 5500,
	6: 7500, 7: 10000, 8: 13000, 9: 16500, 10: 20500,
	11: 26000, 12: 32000, 13: 39000, 14: 47000, 15: 57000,
	16: 69000, 17: 83000, 18: 99000, 19: 119000, 20: 143000,
	21: 175000, 22: 210000, 23: 255000, 24: 310000, 25: 375000,
	26: 450000, 27: 550000, 28: 675000, 29: 825000, 30: 1000000,
}

// xpTable5e: XP acumulado necessário para atingir cada nível (1–20)
var xpTable5e = map[int]int{
	1: 0, 2: 300, 3: 900, 4: 2700, 5: 6500,
	6: 14000, 7: 23000, 8: 34000, 9: 48000, 10: 64000,
	11: 85000, 12: 100000, 13: 120000, 14: 140000, 15: 165000,
	16: 195000, 17: 225000, 18: 265000, 19: 305000, 20: 355000,
}

// Níveis 5e que concedem ASI (Ability Score Improvement). A maioria das
// classes usa a progressão padrão, mas Guerreiro e Ladino têm ASIs bônus
// extras no PHB 2024 (cap. 3): Guerreiro ganha em 6 e 14 além do padrão,
// Ladino ganha em 10 além do padrão. Antes disso havia um único mapa global
// pra todas as classes 5e — Guerreiros e Ladinos ficavam sem essas melhorias
// extras, uma característica de classe real perdida silenciosamente.
var asiLevels5e = map[int]bool{4: true, 8: true, 12: true, 16: true, 19: true}
var asiLevels5eGuerreiro = map[int]bool{4: true, 6: true, 8: true, 12: true, 14: true, 16: true, 19: true}
var asiLevels5eLadino = map[int]bool{4: true, 8: true, 10: true, 12: true, 16: true, 19: true}

// ── Helpers ──────────────────────────────────────────────────────────────────

// mod retorna o modificador de atributo padrão D&D — piso da divisão, não
// truncamento em direção a zero. Para valores ímpares abaixo de 10 isso
// importa: mod(9) deve ser -1 (tabela oficial do PHB), não 0. Go trunca
// divisão de inteiros em direção a zero, então "(attr-10)/2" sozinho erra
// exatamente esses casos (9,7,5,3,1) por +1.
func mod(attr int) int {
	diff := attr - 10
	if diff < 0 {
		return (diff - 1) / 2
	}
	return diff / 2
}

// maxInt retorna o maior entre dois inteiros
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// capAt20 limita o valor entre 1 e 20
func capAt20(v int) int {
	if v > 20 {
		return 20
	}
	if v < 1 {
		return 1
	}
	return v
}

// maxLevel retorna o nível máximo para a edição
func maxLevel(edition string) int {
	if edition == "4e" {
		return 30
	}
	return 20
}

// isASILevel verifica se o nível concede melhoria de atributo.
// 5e: níveis 4, 8, 12, 16, 19 — exceto Guerreiro (+6, +14) e Ladino (+10),
// que têm ASIs bônus extras no PHB 2024 (cap. 3).
// 4e: todo nível par (2, 4, 6, ... 30).
func isASILevel(edition, className string, level int) bool {
	if edition == "5e" {
		switch className {
		case "Guerreiro":
			return asiLevels5eGuerreiro[level]
		case "Ladino":
			return asiLevels5eLadino[level]
		default:
			return asiLevels5e[level]
		}
	}
	return level > 1 && level%2 == 0
}

// proficiencyBonus5e retorna o bônus de proficiência para o nível 5e
func proficiencyBonus5e(level int) int {
	switch {
	case level <= 4:
		return 2
	case level <= 8:
		return 3
	case level <= 12:
		return 4
	case level <= 16:
		return 5
	default:
		return 6
	}
}

// ── Espaços de magia 5e ──────────────────────────────────────────────────────
// Tabelas extraídas de "Características de X" no PHB 2024 (cap. 3), uma por
// classe conjuradora, cruzadas com a extração real do PDF (não de memória —
// ver CLAUDE.md para a verificação linha a linha que confirmou que Paladino
// e Guardião ganham espaços de magia a partir do NÍVEL 1 no PHB 2024, uma
// mudança real em relação ao PHB 2014, onde meio-conjuradores só ganhavam
// magia no nível 2). Índice do slice = nível de personagem - 1; cada
// elemento é [espaços de 1º círculo, 2º, 3º, ...] (índice 0 = 1º círculo).

// fullCasterSlots5e: Bardo, Clérigo, Druida, Feiticeiro, Mago.
var fullCasterSlots5e = [21][]int{
	1:  {2},
	2:  {3},
	3:  {4, 2},
	4:  {4, 3},
	5:  {4, 3, 2},
	6:  {4, 3, 3},
	7:  {4, 3, 3, 1},
	8:  {4, 3, 3, 2},
	9:  {4, 3, 3, 3, 1},
	10: {4, 3, 3, 3, 2},
	11: {4, 3, 3, 3, 2, 1},
	12: {4, 3, 3, 3, 2, 1},
	13: {4, 3, 3, 3, 2, 1, 1},
	14: {4, 3, 3, 3, 2, 1, 1},
	15: {4, 3, 3, 3, 2, 1, 1, 1},
	16: {4, 3, 3, 3, 2, 1, 1, 1},
	17: {4, 3, 3, 3, 2, 1, 1, 1, 1},
	18: {4, 3, 3, 3, 3, 1, 1, 1, 1},
	19: {4, 3, 3, 3, 3, 2, 1, 1, 1},
	20: {4, 3, 3, 3, 3, 2, 2, 1, 1},
}

// halfCasterSlots5e: Paladino, Guardião — verificado contra a tabela real
// (ambas as classes começam com 1 espaço de 1º círculo já no nível 1).
var halfCasterSlots5e = [21][]int{
	1:  {2},
	2:  {2},
	3:  {3},
	4:  {3},
	5:  {4, 2},
	6:  {4, 2},
	7:  {4, 3},
	8:  {4, 3},
	9:  {4, 3, 2},
	10: {4, 3, 2},
	11: {4, 3, 3},
	12: {4, 3, 3},
	13: {4, 3, 3, 1},
	14: {4, 3, 3, 1},
	15: {4, 3, 3, 2},
	16: {4, 3, 3, 2},
	17: {4, 3, 3, 3, 1},
	18: {4, 3, 3, 3, 1},
	19: {4, 3, 3, 3, 2},
	20: {4, 3, 3, 3, 2},
}

// pactMagicSlots5e / pactMagicCircle5e: Magia de Pacto do Bruxo — mecanismo
// próprio (recupera em Descanso Curto ou Longo, não só Longo, e todos os
// espaços são sempre do mesmo círculo, que sobe com o nível).
var pactMagicSlots5e = [21]int{
	1: 1, 2: 2, 3: 2, 4: 2, 5: 2, 6: 2, 7: 2, 8: 2, 9: 2, 10: 2,
	11: 3, 12: 3, 13: 3, 14: 3, 15: 3, 16: 3,
	17: 4, 18: 4, 19: 4, 20: 4,
}
var pactMagicCircle5e = [21]int{
	1: 1, 2: 1, 3: 2, 4: 2, 5: 3, 6: 3, 7: 4, 8: 4, 9: 5, 10: 5,
	11: 5, 12: 5, 13: 5, 14: 5, 15: 5, 16: 5, 17: 5, 18: 5, 19: 5, 20: 5,
}

// thirdCasterSlots5e: Cavaleiro Místico (Guerreiro) e Trapaceiro Arcano
// (Ladino) — tabela idêntica pras duas subclasses, extraída verbatim das
// tabelas "Conjuração de Cavaleiro Místico"/"Conjuração de Trapaceiro
// Arcano" do PHB 2024 (não existia nenhuma tabela de terço-conjurador antes
// disso — Guerreiro/Ladino sempre caíam no "default: nil" abaixo).
var thirdCasterSlots5e = [21][]int{
	3: {2}, 4: {3}, 5: {3},
	6: {3}, 7: {4, 2}, 8: {4, 2}, 9: {4, 2}, 10: {4, 3},
	11: {4, 3}, 12: {4, 3}, 13: {4, 3, 2}, 14: {4, 3, 2}, 15: {4, 3, 2},
	16: {4, 3, 3}, 17: {4, 3, 3}, 18: {4, 3, 3}, 19: {4, 3, 3, 1}, 20: {4, 3, 3, 1},
}

// thirdCasterSubclass indica, por nome de subclasse, a qual classe base ela
// pertence — usado só pra confirmar que o personagem realmente escolheu essa
// subclasse antes de tratá-lo como conjurador (um Guerreiro comum não conjura).
var thirdCasterSubclass = map[string]string{
	"Cavaleiro Místico": "Guerreiro",
	"Trapaceiro Arcano": "Ladino",
}

// spellSlots5e retorna os espaços de magia por círculo (índice 0 = 1º
// círculo) de uma classe 5e num dado nível de personagem. subclassName só
// importa pra Guerreiro/Ladino, pra distinguir um terço-conjurador (Cavaleiro
// Místico/Trapaceiro Arcano) de um personagem comum dessas classes, que não
// conjura. nil para classes sem espaços de magia tradicionais (Bruxo usa
// Magia de Pacto à parte; Bárbaro/Monge não conjuram).
func spellSlots5e(className, subclassName string, level int) []int {
	if level < 1 || level > 20 {
		return nil
	}
	switch className {
	case "Bardo", "Clérigo", "Druida", "Feiticeiro", "Mago":
		return fullCasterSlots5e[level]
	case "Paladino", "Guardião":
		return halfCasterSlots5e[level]
	default:
		if thirdCasterSubclass[subclassName] == className {
			return thirdCasterSlots5e[level]
		}
		return nil
	}
}

// pactMagic5e retorna (usos, círculo) da Magia de Pacto do Bruxo num nível.
func pactMagic5e(level int) (slots, circle int) {
	if level < 1 || level > 20 {
		return 0, 0
	}
	return pactMagicSlots5e[level], pactMagicCircle5e[level]
}

// cantripsKnown5e5e retorna quantos truques uma classe conjuradora conhece
// num dado nível — todas seguem o mesmo formato "base no nível 1, +1 no
// nível 4, +1 no nível 10" (ver a característica "Conjuração"/"Truques" de
// cada classe no cap. 3), exceto Paladino e Guardião, que não têm truques.
// spellcastingAbility5e: atributo de conjuração por classe (extraído das
// características de classe já seedadas — "Atributo de Conjuração" na
// descrição de cada uma). Guerreiro/Ladino não entram aqui: eles só conjuram
// se tiverem escolhido Cavaleiro Místico/Trapaceiro Arcano especificamente
// (ver buildConjuracao em pdf_export_service.go, que trata esse caso à parte
// porque depende da subclasse, não só da classe).
var spellcastingAbility5e = map[string]string{
	"Mago":       "INT",
	"Clérigo":    "SAB",
	"Druida":     "SAB",
	"Guardião":   "SAB",
	"Bardo":      "CAR",
	"Bruxo":      "CAR",
	"Feiticeiro": "CAR",
	"Paladino":   "CAR",
}

// ExpertiseSlotsFor retorna quantas perícias (já proficientes) um personagem
// pode ter com Expertise (Especialização — dobra o Bônus de Proficiência)
// no nível atual. Extraído verbatim da tabela de características de cada
// classe: Ladino "Especialista" nível 1 (2) e nível 6 (+2, total 4); Bardo
// "Especialista" nível 2 (2) e nível 9 (+2, total 4); Guardião "Especialista"
// nível 9 (2, único). Não é um nível fixo genérico — cada classe tem o seu.
func ExpertiseSlotsFor(className string, level int) int {
	switch className {
	case "Ladino":
		if level >= 6 {
			return 4
		}
		if level >= 1 {
			return 2
		}
	case "Bardo":
		if level >= 9 {
			return 4
		}
		if level >= 2 {
			return 2
		}
	case "Guardião":
		if level >= 9 {
			return 2
		}
	}
	return 0
}

func cantripsKnown5e(className string, level int) int {
	base, hasCantrips := map[string]int{
		"Bardo": 2, "Bruxo": 2, "Clérigo": 3, "Druida": 2, "Feiticeiro": 4, "Mago": 3,
	}[className]
	if !hasCantrips {
		return 0
	}
	n := base
	if level >= 4 {
		n++
	}
	if level >= 10 {
		n++
	}
	return n
}

// XPParaProximoNivel retorna o XP necessário para o próximo nível
func XPParaProximoNivel(edition string, currentLevel int) int {
	nextLevel := currentLevel + 1
	table := xpTable5e
	if edition == "4e" {
		table = xpTable4e
	}
	threshold, ok := table[nextLevel]
	if !ok {
		return 0 // nível máximo
	}
	return threshold
}

// ── Cálculo de HP ─────────────────────────────────────────────────────────────

// calcHP4e calcula HP inicial, surge value e surges/dia no 4e
func calcHP4e(c *domain.Character) (hp, surgeVal, surgesPerDay int) {
	hp = c.Class.BaseHP + c.Constitution
	if hp < 1 {
		hp = 1
	}
	surgeVal = hp / 4
	surgesPerDay = c.Class.SurgesPerDay + mod(c.Constitution)
	if surgesPerDay < 1 {
		surgesPerDay = 1
	}
	return
}

// calcHP5e calcula HP inicial no nível 1 (5e)
func calcHP5e(c *domain.Character) int {
	hp := c.Class.HitDie + mod(c.Constitution)
	if hp < 1 {
		hp = 1
	}
	return hp
}

// calcLevelUpHP calcula o HP ganho ao subir de nível
func calcLevelUpHP(c *domain.Character) int {
	switch c.Edition {
	case "4e":
		return c.Class.HPPerLevel
	default: // 5e
		gain := (c.Class.HitDie/2 + 1) + mod(c.Constitution)
		if gain < 1 {
			gain = 1
		}
		return gain
	}
}

// calcDefenses4e calcula CA, FORT, REFL e VONT conforme regras do 4e
// CA usa o maior entre o modificador de DES e INT (regra 4e pós-errata),
// o que naturalmente aplica o atributo certo para cada classe conforme
// o personagem tenha investido em DES ou INT — mesma lógica já usada
// abaixo para Reflexo.
func calcDefenses4e(c *domain.Character, armorBonus int, maxDexBonus int) {
	half := c.Level / 2

	abilityMod := maxInt(mod(c.Dexterity), mod(c.Intelligence))
	if maxDexBonus >= 0 && abilityMod > maxDexBonus {
		abilityMod = maxDexBonus
	}

	c.Defense_AC = 10 + half + abilityMod + armorBonus
	c.Defense_Fort = 10 + half + maxInt(mod(c.Strength), mod(c.Constitution)) + c.Class.FortBonus
	c.Defense_Refl = 10 + half + maxInt(mod(c.Dexterity), mod(c.Intelligence)) + c.Class.ReflBonus
	c.Defense_Will = 10 + half + maxInt(mod(c.Wisdom), mod(c.Charisma)) + c.Class.WillBonus
}

// recalculate recalcula HP, surges e defesas do personagem (usado na criação)
func (s *CharacterService) recalculate(c *domain.Character) {
	armorBonus := 0
	maxDexBonus := -1
	if c.Armor != nil {
		armorBonus = c.Armor.BaseAC
		maxDexBonus = c.Armor.MaxDexBonus
	}

	switch c.Edition {
	case "4e":
		hp, surgeVal, surgesPerDay := calcHP4e(c)
		c.HitPoints = hp
		c.MaxHP = hp
		c.SurgeValue = surgeVal
		c.SurgesPerDay = surgesPerDay
		calcDefenses4e(c, armorBonus, maxDexBonus)
	default: // 5e
		hp := calcHP5e(c)
		c.HitPoints = hp
		c.MaxHP = hp
		c.ProficiencyBonus = proficiencyBonus5e(c.Level)
	}
}

// recalcDefensas4e recalcula apenas as defesas 4e (usado no level up)
func (s *CharacterService) recalcDefensas4e(c *domain.Character) {
	armorBonus := 0
	maxDexBonus := -1
	if c.Armor != nil {
		armorBonus = c.Armor.BaseAC
		maxDexBonus = c.Armor.MaxDexBonus
	}
	calcDefenses4e(c, armorBonus, maxDexBonus)
}

// ── LevelUpResult: retorno dos endpoints de XP e ASI ─────────────────────────

type LevelUpResult struct {
	LeveledUp bool `json:"leveled_up"`
	NeedsASI  bool `json:"needs_asi"` // true = jogador deve escolher melhoria de atributo
	NewLevel  int  `json:"new_level"`
}

// ASIChoice: melhoria escolhida pelo jogador num nível de ASI — ou pontos de
// atributo, ou um talento no lugar (RAW: "+2 em um atributo, +1 em dois
// atributos, ou um talento à sua escolha" — nunca as duas coisas no mesmo
// nível). Ver ApplyASI.
type ASIChoice struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`

	// TalentoID: se informado, o jogador escolheu pegar um talento em vez de
	// melhorar atributos neste nível de ASI — nesse caso todos os campos de
	// atributo acima devem ficar zerados (ver validação em ApplyASI).
	TalentoID *uint `json:"talento_id"`
}

// ── checkAndApplyLevelUps ─────────────────────────────────────────────────────
// Aplica todos os level ups automáticos para o XP atual.
// Para em níveis ASI para aguardar a escolha do jogador.
// Chamado tanto pelo AddXP quanto pelo ApplyASI (para continuar após a escolha).
func (s *CharacterService) checkAndApplyLevelUps(c *domain.Character) (leveledUp bool, needsASI bool) {
	ml := maxLevel(c.Edition)
	table := xpTable5e
	if c.Edition == "4e" {
		table = xpTable4e
	}

	for c.Level < ml {
		nextLevel := c.Level + 1
		threshold, ok := table[nextLevel]
		if !ok || c.ExperiencePoints < threshold {
			break // XP insuficiente para o próximo nível
		}

		// ── Level up ──
		c.Level = nextLevel
		leveledUp = true

		// HP ganho com o CON atual (antes de qualquer ASI)
		hpGain := calcLevelUpHP(c)
		c.HitPoints += hpGain
		c.MaxHP += hpGain

		// Atualiza bônus de proficiência (5e)
		if c.Edition == "5e" {
			c.ProficiencyBonus = proficiencyBonus5e(c.Level)
		}

		// Recalcula surge e defesas (4e)
		if c.Edition == "4e" {
			c.SurgeValue = c.MaxHP / 4
			s.recalcDefensas4e(c)
		}

		// Adiciona habilidades que desbloqueiam neste nível. Habilidades que
		// exigem escolha (RequiresChoice) não são concedidas automaticamente
		// aqui — hoje nenhuma existe acima do nível 1 (todas são resolvidas
		// na criação, ver CharacterCreate.tsx), mas se uma for adicionada no
		// futuro isso evita conceder as duas opções de uma vez em vez de
		// esperar uma escolha do jogador que ainda não existe nesta tela.
		//
		// Habilidades de subclasse em níveis altos (ex.: "Ramos da Árvore" no
		// nível 6 do Bárbaro) usam ChoiceGroup para guardar o Name exato da
		// habilidade de entrada da subclasse (ex.: "Trilha da Árvore do
		// Mundo") em vez de um slug de grupo de escolha — RequiresChoice fica
		// false porque a escolha já foi feita no nível 3, não há nada para
		// escolher de novo aqui. requiresSubclass() confere se o personagem já
		// tem essa entrada de subclasse antes de conceder; sem essa checagem,
		// TODO bárbaro ganharia as características de TODAS as trilhas ao
		// subir de nível, não só da que escolheu.
		newSkills, err := s.SkillRepo.FindByLevel(c.ClassID, c.Level)
		if err == nil {
			for _, skill := range newSkills {
				if skill.RequiresChoice {
					continue
				}
				if !hasChosenSubclass(c.Skills, skill.ChoiceGroup) {
					continue
				}
				s.Repo.AddSkill(c, &skill)
			}
		}

		// Se é nível ASI, para aqui e aguarda escolha do jogador
		if isASILevel(c.Edition, c.Class.Name, c.Level) {
			needsASI = true
			break
		}
	}

	return leveledUp, needsASI
}

// ── CRUD ──────────────────────────────────────────────────────────────────────

func (s *CharacterService) GetAll(userID uint) ([]domain.Character, error) {
	return s.Repo.FindAll(userID)
}

func (s *CharacterService) GetByID(id uint) (domain.Character, error) {
	return s.Repo.FindByID(id)
}

func (s *CharacterService) Create(character *domain.Character) error {
	if character.Name == "" {
		return errors.New("nome do personagem é obrigatório")
	}
	if character.ClassID == 0 {
		return errors.New("classe é obrigatória")
	}
	if character.RaceID == 0 {
		return errors.New("raça é obrigatória")
	}

	// character.Class só tem o zero-value aqui (o handler só recebe class_id
	// do JSON, não o objeto Class inteiro) — sem isso, recalculate() calcula o
	// PV inicial como se HitDie fosse 0, dando PV = max(1, mod(CON)) pra
	// qualquer classe. Precisa recarregar a classe antes de calcular.
	class, err := s.ClassRepo.FindByID(character.ClassID)
	if err != nil {
		return errors.New("classe inválida")
	}
	character.Class = class

	// Mesmo motivo do Class acima: o handler só recebe race_id do JSON, não a
	// Race inteira. Sem isso, Character.Speed nunca era preenchido (não há
	// campo de velocidade no formulário de criação) e ficava sempre no
	// zero-value 0 — a ficha PDF exportada mostrava "Deslocamento: 0" pra
	// todo personagem 5e, em vez do valor real da raça (`Race.Speed`, que já
	// varia por raça no seed).
	race, err := s.RaceRepo.FindByID(character.RaceID)
	if err != nil {
		return errors.New("raça inválida")
	}
	character.Race = race
	if character.Edition == "5e" {
		character.Speed = race.Speed
	}

	// Regra 2024: o bônus de atributo vem do Antecedente, nunca da raça (por
	// isso não há nenhum campo de bônus em domain.Race). Se o personagem tem
	// um antecedente com AbilityBonusOptions, o bônus é obrigatório aqui —
	// os atributos que chegaram no payload são tratados como BASE (antes do
	// bônus), e o valor final gravado já inclui a distribuição escolhida.
	var antecedent domain.Antecedent
	if character.Edition == "5e" && character.AntecedentID != nil {
		if err := s.DB.First(&antecedent, *character.AntecedentID).Error; err != nil {
			return errors.New("antecedente inválido")
		}
		if err := s.applyAntecedentAbilityBonus(character, &antecedent); err != nil {
			return err
		}
	}

	manualHP := character.HitPoints

	character.Level = 1
	character.ExperiencePoints = 0
	s.recalculate(character)

	// Respeita o HP definido manualmente no formulário de criação, em vez de
	// sobrescrever com o valor calculado pela fórmula de classe/CON.
	if manualHP > 0 {
		character.HitPoints = manualHP
		character.MaxHP = manualHP
		if character.Edition == "4e" {
			character.SurgeValue = manualHP / 4
		}
	}

	if err := s.Repo.Create(character); err != nil {
		return err
	}

	// Concede o talento de Origem do antecedente. RAW 2024 normal: "seu
	// antecedente concede um talento", fixo por antecedente, não é escolha
	// livre. Exceção (caixa "Antecedentes e Espécies de Livros Antigos",
	// PHB 2024 cap. 2, p.38): se o antecedente é de um livro antigo e não
	// fornece talento nenhum (IsLegacy, OriginFeatName vazio), o jogador
	// escolhe livremente qualquer Talento de categoria "Origem".
	if antecedent.OriginFeatName != "" {
		talento, err := s.TalentoRepo.FindByName(antecedent.OriginFeatName, "5e")
		if err == nil {
			s.TalentoRepo.Add(character.ID, talento.ID)
		}
	} else if antecedent.IsLegacy {
		if character.OriginFeatChoiceID == nil {
			return errors.New("escolha um talento de Origem (antecedente de livro antigo, sem talento fixo)")
		}
		talento, err := s.TalentoRepo.FindByID(*character.OriginFeatChoiceID)
		if err != nil {
			return errors.New("talento de origem inválido")
		}
		if talento.Category != "Origem" || talento.Edition != "5e" {
			return errors.New("escolha um talento da categoria Origem (5e)")
		}
		s.TalentoRepo.Add(character.ID, talento.ID)
	}

	// Concede o "Equipamento Inicial" da classe (Livro do Jogador 2024, cap.
	// 3, tabela "Traços Básicos de X"), se o jogador escolheu uma das opções
	// lettered (A/B/C) na criação. Opcional — sem escolha, o personagem não
	// recebe nada aqui, igual ao comportamento anterior a esta feature.
	if character.EquipmentOptionID != nil {
		if err := s.grantStartingEquipment(character, *character.EquipmentOptionID); err != nil {
			return err
		}
	}

	// Comum é concedido automaticamente a todo personagem 5e — RAW 2024 (Livro
	// do Jogador, cap. 2, "Escolha Idiomas"): "seu personagem sabe pelo menos
	// três idiomas: Comum e mais dois...". Os outros 2 são escolha livre do
	// jogador (ver CharacterCreate.tsx), concedidos via POST
	// /characters/:id/languages/:language_id como os demais catálogos
	// (Talento/Spell) — aqui só o automático.
	if character.Edition == "5e" {
		var comum domain.Language
		if err := s.DB.Where("name = ? AND edition = ?", "Comum", "5e").First(&comum).Error; err == nil {
			s.DB.Exec(
				"INSERT INTO character_languages (character_id, language_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				character.ID, comum.ID,
			)
		}
	}

	return nil
}

// grantStartingEquipment aplica os itens/armaduras/PO de uma
// ClassEquipmentOption ao personagem recém-criado. Usa a mesma lógica de
// upsert por (character_id, item_id/armor_id) do InventoryService, pra não
// duplicar linha se o mesmo item aparecer em mais de um componente.
func (s *CharacterService) grantStartingEquipment(character *domain.Character, optionID uint) error {
	var option domain.ClassEquipmentOption
	if err := s.DB.Preload("Components").First(&option, optionID).Error; err != nil {
		return errors.New("opção de equipamento inicial inválida")
	}
	if option.ClassID != character.ClassID {
		return errors.New("opção de equipamento inicial não pertence à classe do personagem")
	}

	for _, comp := range option.Components {
		qty := comp.Quantity
		if qty <= 0 {
			qty = 1
		}
		switch {
		case comp.ItemID != nil:
			var existing domain.CharacterItem
			err := s.DB.Where("character_id = ? AND item_id = ?", character.ID, *comp.ItemID).First(&existing).Error
			if err != nil {
				s.DB.Create(&domain.CharacterItem{CharacterID: character.ID, ItemID: *comp.ItemID, Quantity: qty})
			} else {
				existing.Quantity += qty
				s.DB.Save(&existing)
			}
		case comp.ArmorID != nil:
			var existing domain.CharacterArmorOwned
			err := s.DB.Where("character_id = ? AND armor_id = ?", character.ID, *comp.ArmorID).First(&existing).Error
			if err != nil {
				s.DB.Create(&domain.CharacterArmorOwned{CharacterID: character.ID, ArmorID: *comp.ArmorID, Quantity: qty})
			} else {
				existing.Quantity += qty
				s.DB.Save(&existing)
			}
		}
	}

	if option.GoldPieces > 0 {
		character.GoldPieces += option.GoldPieces
		s.DB.Model(character).Update("gold_pieces", character.GoldPieces)
	}
	return nil
}

// applyAntecedentAbilityBonus aplica a distribuição de bônus de atributo do
// Antecedente 2024 (regra de ouro deste projeto: bônus de atributo NUNCA vem
// da raça/espécie, só do antecedente — ver domain.Race, que não tem nenhum
// campo de bônus). Espera exatamente +2 numa habilidade e +1 em outra (das
// opções permitidas), OU +1 em todas elas.
//
// Antecedentes 2024 normais restringem a 3 atributos nomeados
// (AbilityBonusOptions). Antecedentes de livro antigo (IsLegacy, sem
// AbilityBonusOptions) liberam a escolha entre os 6 atributos — caixa
// "Antecedentes e Espécies de Livros Antigos", PHB 2024 cap. 2, p.38.
func (s *CharacterService) applyAntecedentAbilityBonus(character *domain.Character, antecedent *domain.Antecedent) error {
	var allowed []string
	switch {
	case antecedent.AbilityBonusOptions != "":
		if err := json.Unmarshal([]byte(antecedent.AbilityBonusOptions), &allowed); err != nil || len(allowed) != 3 {
			return nil // dado do antecedente mal formado no seed — não bloqueia a criação
		}
	case antecedent.IsLegacy:
		allowed = []string{"FOR", "DES", "CON", "INT", "SAB", "CAR"}
	default:
		return nil
	}
	allowedSet := map[string]bool{}
	for _, a := range allowed {
		allowedSet[a] = true
	}

	choice := character.AbilityBonusChoice
	if len(choice) == 0 {
		return fmt.Errorf("escolha a distribuição do bônus de atributo do antecedente (opções: %v)", allowed)
	}

	total := 0
	for attr, val := range choice {
		if !allowedSet[attr] {
			return fmt.Errorf("%s não é uma opção de bônus deste antecedente (opções: %v)", attr, allowed)
		}
		if val < 1 || val > 2 {
			return errors.New("cada atributo escolhido deve receber +1 ou +2")
		}
		total += val
	}
	// +2/+1 em duas habilidades (soma 3, 2 escolhas) OU +1 nas três (soma 3, 3 escolhas)
	if total != 3 || (len(choice) != 2 && len(choice) != 3) {
		return errors.New("distribua +2 em uma habilidade e +1 em outra, ou +1 nas três habilidades do antecedente")
	}

	for attr, val := range choice {
		switch attr {
		case "FOR":
			character.Strength = capAt20(character.Strength + val)
		case "DES":
			character.Dexterity = capAt20(character.Dexterity + val)
		case "CON":
			character.Constitution = capAt20(character.Constitution + val)
		case "INT":
			character.Intelligence = capAt20(character.Intelligence + val)
		case "SAB":
			character.Wisdom = capAt20(character.Wisdom + val)
		case "CAR":
			character.Charisma = capAt20(character.Charisma + val)
		}
	}
	return nil
}

func (s *CharacterService) Update(character *domain.Character) error {
	if character.Name == "" {
		return errors.New("nome do personagem é obrigatório")
	}
	return s.Repo.Update(character)
}

func (s *CharacterService) Delete(id uint) error {
	return s.Repo.Delete(id)
}

// ── AddXP ─────────────────────────────────────────────────────────────────────
// Adiciona XP ao personagem e aplica level ups automáticos.
// Retorna o personagem atualizado, info de level up e erro.
func (s *CharacterService) AddXP(id uint, xp int) (domain.Character, LevelUpResult, error) {
	character, err := s.Repo.FindByID(id)
	if err != nil {
		return character, LevelUpResult{}, errors.New("personagem não encontrado")
	}

	if xp <= 0 {
		return character, LevelUpResult{}, errors.New("XP deve ser maior que zero")
	}

	ml := maxLevel(character.Edition)
	if character.Level >= ml {
		// Já no nível máximo — apenas registra o XP mas não sobe mais
		return character, LevelUpResult{NewLevel: character.Level}, nil
	}

	character.ExperiencePoints += xp

	leveledUp, needsASI := s.checkAndApplyLevelUps(&character)

	if err := s.Repo.Update(&character); err != nil {
		return character, LevelUpResult{}, err
	}

	result := LevelUpResult{
		LeveledUp: leveledUp,
		NeedsASI:  needsASI,
		NewLevel:  character.Level,
	}

	return character, result, nil
}

// ── ApplyASI ──────────────────────────────────────────────────────────────────
// Aplica melhorias de atributo (ASI) escolhidas pelo jogador após level up.
// Após aplicar, verifica se há mais level ups pendentes pelo XP atual.
func (s *CharacterService) ApplyASI(id uint, choice ASIChoice) (domain.Character, LevelUpResult, error) {
	character, err := s.Repo.FindByID(id)
	if err != nil {
		return character, LevelUpResult{}, errors.New("personagem não encontrado")
	}

	// 5e: +2 em 1 atributo OU +1 em 2 atributos OU um talento (nunca as duas
	// coisas no mesmo nível de ASI — RAW e decisão de produto explícita).
	// 4e: +1 em 2 atributos (mesma lógica — total = 2); não há troca por
	// talento no 4e, já que 4e não usa Talento (RAW 4e não tem essa opção).
	total := choice.Strength + choice.Dexterity + choice.Constitution +
		choice.Intelligence + choice.Wisdom + choice.Charisma

	if choice.TalentoID != nil {
		if total != 0 {
			return character, LevelUpResult{}, errors.New("não é possível escolher atributo e talento na mesma melhoria — escolha um ou outro")
		}
		if character.Edition != "5e" {
			return character, LevelUpResult{}, errors.New("troca de ASI por talento só existe no 5e")
		}
		talento, err := s.TalentoRepo.FindByID(*choice.TalentoID)
		if err != nil {
			return character, LevelUpResult{}, errors.New("talento inválido")
		}
		if err := s.TalentoRepo.Add(character.ID, talento.ID); err != nil {
			return character, LevelUpResult{}, err
		}
	} else {
		if total < 1 || total > 2 {
			return character, LevelUpResult{}, errors.New("escolha +2 em um atributo ou +1 em dois atributos (ou um talento, via talento_id)")
		}

		// Guarda o mod CON antigo para ajuste retroativo de HP (5e)
		oldConMod := mod(character.Constitution)

		// Aplica melhorias com cap em 20
		character.Strength = capAt20(character.Strength + choice.Strength)
		character.Dexterity = capAt20(character.Dexterity + choice.Dexterity)
		character.Constitution = capAt20(character.Constitution + choice.Constitution)
		character.Intelligence = capAt20(character.Intelligence + choice.Intelligence)
		character.Wisdom = capAt20(character.Wisdom + choice.Wisdom)
		character.Charisma = capAt20(character.Charisma + choice.Charisma)

		// 5e: se o mod de CON aumentou, +1 HP por nível retroativamente
		// Ex: CON 12→14 (mod +1→+2) no nível 4 = +4 HP permanentes
		if character.Edition == "5e" && choice.Constitution > 0 {
			newConMod := mod(character.Constitution)
			if newConMod > oldConMod {
				bonusHP := (newConMod - oldConMod) * character.Level
				character.MaxHP += bonusHP
				character.HitPoints += bonusHP
			}
		}

		// 4e: recalcula surges com novo CON
		if character.Edition == "4e" {
			character.SurgeValue = character.MaxHP / 4
			character.SurgesPerDay = character.Class.SurgesPerDay + mod(character.Constitution)
			if character.SurgesPerDay < 1 {
				character.SurgesPerDay = 1
			}
		}
	}

	// Continua verificando level ups pendentes (caso o XP já cubra mais níveis)
	leveledUp, needsASI := s.checkAndApplyLevelUps(&character)

	if err := s.Repo.Update(&character); err != nil {
		return character, LevelUpResult{}, err
	}

	result := LevelUpResult{
		LeveledUp: leveledUp,
		NeedsASI:  needsASI,
		NewLevel:  character.Level,
	}

	return character, result, nil
}

// ── LevelUp manual (mantido para compatibilidade) ────────────────────────────
func (s *CharacterService) LevelUp(id uint) (domain.Character, error) {
	character, err := s.Repo.FindByID(id)
	if err != nil {
		return character, errors.New("personagem não encontrado")
	}

	ml := maxLevel(character.Edition)
	if character.Level >= ml {
		return character, errors.New("personagem já está no nível máximo")
	}

	character.Level++

	hpGain := calcLevelUpHP(&character)
	character.HitPoints += hpGain
	character.MaxHP += hpGain

	if character.Edition == "5e" {
		character.ProficiencyBonus = proficiencyBonus5e(character.Level)
	}

	if character.Edition == "4e" {
		character.SurgeValue = character.MaxHP / 4
		s.recalcDefensas4e(&character)
	}

	newSkills, err := s.SkillRepo.FindByLevel(character.ClassID, character.Level)
	if err == nil {
		for _, skill := range newSkills {
			if skill.RequiresChoice {
				continue
			}
			if !hasChosenSubclass(character.Skills, skill.ChoiceGroup) {
				continue
			}
			s.Repo.AddSkill(&character, &skill)
		}
	}

	err = s.Repo.Update(&character)
	return character, err
}

// hasChosenSubclass reporta se o personagem já possui a habilidade de entrada
// de subclasse referenciada. group vazio significa que a habilidade é uma
// característica base da classe (sem subclasse envolvida) e sempre libera.
// Quando não vazio, group é o Name exato da habilidade de entrada da
// subclasse (ex.: "Trilha da Árvore do Mundo"), não um slug de ChoiceGroup —
// convenção usada só pelas habilidades de subclasse de nível alto seedadas a
// partir de 2026-08-27 (ver seed.go).
func hasChosenSubclass(skills []domain.Skill, group string) bool {
	if group == "" {
		return true
	}
	for _, s := range skills {
		if s.Name == group {
			return true
		}
	}
	return false
}

// ── Skills ────────────────────────────────────────────────────────────────────

func (s *CharacterService) AddSkill(characterID, skillID uint) error {
	character, err := s.Repo.FindByID(characterID)
	if err != nil {
		return errors.New("personagem não encontrado")
	}

	skill, err := s.SkillRepo.FindByID(skillID)
	if err != nil {
		return errors.New("habilidade não encontrada")
	}

	return s.Repo.AddSkill(&character, &skill)
}

func (s *CharacterService) RemoveSkill(characterID, skillID uint) error {
	character, err := s.Repo.FindByID(characterID)
	if err != nil {
		return errors.New("personagem não encontrado")
	}

	skill, err := s.SkillRepo.FindByID(skillID)
	if err != nil {
		return errors.New("habilidade não encontrada")
	}

	return s.Repo.RemoveSkill(&character, &skill)
}
