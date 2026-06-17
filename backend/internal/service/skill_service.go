package service

import (
    "errors"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/repository"
)

type SkillService struct {
    Repo *repository.SkillRepository
}

func NewSkillService(repo *repository.SkillRepository) *SkillService {
    return &SkillService{Repo: repo}
}

func (s *SkillService) GetAll() ([]domain.Skill, error) {
    return s.Repo.FindAll()
}

func (s *SkillService) GetByClassAndRace(classID uint, raceID *uint) ([]domain.Skill, error) {
    return s.Repo.FindByClassAndRace(classID, raceID)
}

func (s *SkillService) GetByLevel(classID uint, level int) ([]domain.Skill, error) {
    if level < 1 {
        return nil, errors.New("nível deve ser maior que zero")
    }
    return s.Repo.FindByLevel(classID, level)
}

func (s *SkillService) Create(skill *domain.Skill) error {
    if skill.Name == "" {
        return errors.New("nome da habilidade é obrigatório")
    }
    if skill.PowerType == "" {
        return errors.New("tipo do poder é obrigatório")
    }
    if skill.Edition == "" {
        return errors.New("edição é obrigatória")
    }
    return s.Repo.Create(skill)
}

func (s *SkillService) Update(skill *domain.Skill) error {
    if skill.Name == "" {
        return errors.New("nome da habilidade é obrigatório")
    }
    return s.Repo.Update(skill)
}

func (s *SkillService) Delete(id uint) error {
    return s.Repo.Delete(id)
}