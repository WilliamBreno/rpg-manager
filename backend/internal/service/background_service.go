package service

import (
    "errors"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/repository"
)

type BackgroundService struct {
    Repo *repository.BackgroundRepository
}

func NewBackgroundService(repo *repository.BackgroundRepository) *BackgroundService {
    return &BackgroundService{Repo: repo}
}

func (s *BackgroundService) GetByCharacterID(characterID uint) (domain.Background, error) {
    return s.Repo.FindByCharacterID(characterID)
}

func (s *BackgroundService) Save(background *domain.Background) error {
    if background.CharacterID == 0 {
        return errors.New("personagem é obrigatório")
    }

    // Verifica se já existe um background para esse personagem
    existing, err := s.Repo.FindByCharacterID(background.CharacterID)
    if err == nil {
        // Já existe — atualiza
        background.ID = existing.ID
        return s.Repo.Update(background)
    }

    // Não existe — cria
    return s.Repo.Create(background)
}

func (s *BackgroundService) Delete(characterID uint) error {
    return s.Repo.Delete(characterID)
}