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

// Níveis 5e que concedem ASI (Ability Score Improvement)
var asiLevels5e = map[int]bool{4: true, 8: true, 12: true, 16: true, 19: true}

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

// isASILevel verifica se o nível concede melhoria de atributo
// 5e: níveis 4, 8, 12, 16, 19
// 4e: todo nível par (2, 4, 6, ... 30)
func isASILevel(edition string, level int) bool {
	if edition == "5e" {
		return asiLevels5e[level]
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
		newSkills, err := s.SkillRepo.FindByLevel(c.ClassID, c.Level)
		if err == nil {
			for _, skill := range newSkills {
				if skill.RequiresChoice {
					continue
				}
				s.Repo.AddSkill(c, &skill)
			}
		}

		// Se é nível ASI, para aqui e aguarda escolha do jogador
		if isASILevel(c.Edition, c.Level) {
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
			s.Repo.AddSkill(&character, &skill)
		}
	}

	err = s.Repo.Update(&character)
	return character, err
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
