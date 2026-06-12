package service

import (
    "errors"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/repository"
)

type ClassService struct {
    Repo *repository.ClassRepository
}

func NewClassService(repo *repository.ClassRepository) *ClassService {
    return &ClassService{Repo: repo}
}

func (s *ClassService) GetAll() ([]domain.Class, error) {
    return s.Repo.FindAll()
}

func (s *ClassService) GetByID(id uint) (domain.Class, error) {
    return s.Repo.FindByID(id)
}

func (s *ClassService) GetByEdition(edition string) ([]domain.Class, error) {
    return s.Repo.FindByEdition(edition)
}

func (s *ClassService) Create(class *domain.Class) error {
    if class.Name == "" {
        return errors.New("nome da classe é obrigatório")
    }
    if class.Edition == "" {
        return errors.New("edição é obrigatória")
    }
    if class.HitDie == 0 {
        return errors.New("hit die é obrigatório")
    }
    return s.Repo.Create(class)
}

func (s *ClassService) Update(class *domain.Class) error {
    if class.Name == "" {
        return errors.New("nome da classe é obrigatório")
    }
    return s.Repo.Update(class)
}

func (s *ClassService) Delete(id uint) error {
    return s.Repo.Delete(id)
}