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

// calcHP4e calcula HP inicial, valor do pulso e pulsos por dia no 4e
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

// calcHP5e calcula HP inicial no 5e (nível 1)
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
func calcDefenses4e(c *domain.Character, armorBonus int, maxDexBonus int) {
	half := c.Level / 2

	// Bônus de DEX respeitando o limite da armadura
	dexMod := mod(c.Dexterity)
	if maxDexBonus >= 0 && dexMod > maxDexBonus {
		dexMod = maxDexBonus
	}

	// CA = 10 + ½ nível + mod DES (limitado pela armadura) + bônus armadura
	c.Defense_AC = 10 + half + dexMod + armorBonus

	// FORT = 10 + ½ nível + max(mod FOR, mod CON) + bônus classe
	c.Defense_Fort = 10 + half + maxInt(mod(c.Strength), mod(c.Constitution)) + c.Class.FortBonus

	// REFL = 10 + ½ nível + max(mod DES, mod INT) + bônus classe
	c.Defense_Refl = 10 + half + maxInt(mod(c.Dexterity), mod(c.Intelligence)) + c.Class.ReflBonus

	// VONT = 10 + ½ nível + max(mod SAB, mod CAR) + bônus classe
	c.Defense_Will = 10 + half + maxInt(mod(c.Wisdom), mod(c.Charisma)) + c.Class.WillBonus
}

// recalculate recalcula HP, pulsos e defesas do personagem
func (s *CharacterService) recalculate(c *domain.Character) {
	armorBonus := 0
	maxDexBonus := -1 // -1 = sem limite (sem armadura)
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
	}
}

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

	character.Level = 1
	s.recalculate(character)
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

func (s *CharacterService) LevelUp(id uint) (domain.Character, error) {
	character, err := s.Repo.FindByID(id)
	if err != nil {
		return character, errors.New("personagem não encontrado")
	}

	if character.Level >= 20 {
		return character, errors.New("personagem já está no nível máximo")
	}

	character.Level++

	hpGain := calcLevelUpHP(&character)
	character.HitPoints += hpGain
	character.MaxHP += hpGain

	// Recalcula surges e defesas com o novo nível
	if character.Edition == "4e" {
		character.SurgeValue = character.MaxHP / 4

		armorBonus := 0
		maxDexBonus := -1
		if character.Armor != nil {
			armorBonus = character.Armor.BaseAC
			maxDexBonus = character.Armor.MaxDexBonus
		}
		calcDefenses4e(&character, armorBonus, maxDexBonus)
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