package service

import (
    "errors"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/repository"
)

type RaceService struct {
    Repo *repository.RaceRepository
}

func NewRaceService(repo *repository.RaceRepository) *RaceService {
    return &RaceService{Repo: repo}
}

func (s *RaceService) GetAll() ([]domain.Race, error) {
    return s.Repo.FindAll()
}

func (s *RaceService) GetByID(id uint) (domain.Race, error) {
    return s.Repo.FindByID(id)
}

func (s *RaceService) GetByEdition(edition string) ([]domain.Race, error) {
    return s.Repo.FindByEdition(edition)
}

func (s *RaceService) Create(race *domain.Race) error {
    if race.Name == "" {
        return errors.New("nome da raça é obrigatório")
    }
    if race.Edition == "" {
        return errors.New("edição é obrigatória")
    }
    return s.Repo.Create(race)
}

func (s *RaceService) Update(race *domain.Race) error {
    if race.Name == "" {
        return errors.New("nome da raça é obrigatório")
    }
    return s.Repo.Update(race)
}

func (s *RaceService) Delete(id uint) error {
    return s.Repo.Delete(id)
}