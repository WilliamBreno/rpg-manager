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

// conModifier retorna o modificador de CON padrão D&D
func conModifier(constitution int) int {
    return (constitution - 10) / 2
}

// calcInitialHP calcula o HP inicial (nível 1) com base na edição
func calcInitialHP(c *domain.Character) int {
    switch c.Edition {
    case "5e":
        // 5e: máximo do hit die + modificador de CON
        hp := c.Class.HitDie + conModifier(c.Constitution)
        if hp < 1 {
            hp = 1
        }
        return hp
    case "4e":
        // 4e: BaseHP da classe + valor de CON (não modificador!)
        hp := c.Class.BaseHP + c.Constitution
        if hp < 1 {
            hp = 1
        }
        return hp
    default:
        return c.Class.HitDie + conModifier(c.Constitution)
    }
}

// calcLevelUpHP calcula o HP ganho ao subir de nível
func calcLevelUpHP(c *domain.Character) int {
    switch c.Edition {
    case "5e":
        // 5e: média do hit die (HitDie/2 + 1) + modificador de CON
        gain := (c.Class.HitDie/2 + 1) + conModifier(c.Constitution)
        if gain < 1 {
            gain = 1
        }
        return gain
    case "4e":
        // 4e: HPPerLevel fixo por nível
        return c.Class.HPPerLevel
    default:
        gain := c.Class.HitDie/2 + 1 + conModifier(c.Constitution)
        if gain < 1 {
            gain = 1
        }
        return gain
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

    // Calcula HP automaticamente com base na edição e classe
    initialHP := calcInitialHP(character)
    character.HitPoints = initialHP
    character.MaxHP = initialHP

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

    // HP gain baseado na edição
    hpGain := calcLevelUpHP(&character)
    character.HitPoints += hpGain
    character.MaxHP += hpGain

    // Busca novas habilidades desbloqueadas pelo novo nível
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