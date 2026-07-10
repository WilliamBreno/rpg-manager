package service

import (
	"errors"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/repository"
)

type CharacterService struct {
	Repo      *repository.CharacterRepository
	SkillRepo *repository.SkillRepository
}

func NewCharacterService(repo *repository.CharacterRepository, skillRepo *repository.SkillRepository) *CharacterService {
	return &CharacterService{Repo: repo, SkillRepo: skillRepo}
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

// mod retorna o modificador de atributo padrão D&D
func mod(attr int) int {
	return (attr - 10) / 2
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

// ASIChoice: melhorias de atributo escolhidas pelo jogador
type ASIChoice struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
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

		// Adiciona habilidades que desbloqueiam neste nível
		newSkills, err := s.SkillRepo.FindByLevel(c.ClassID, c.Level)
		if err == nil {
			for _, skill := range newSkills {
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

	return s.Repo.Create(character)
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

	// Valida total de pontos
	// 5e: +2 em 1 atributo OU +1 em 2 atributos
	// 4e: +1 em 2 atributos (mesma lógica — total = 2)
	total := choice.Strength + choice.Dexterity + choice.Constitution +
		choice.Intelligence + choice.Wisdom + choice.Charisma

	if total < 1 || total > 2 {
		return character, LevelUpResult{}, errors.New("escolha +2 em um atributo ou +1 em dois atributos")
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
