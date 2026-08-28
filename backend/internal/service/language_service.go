package service

import "rpg-manager/internal/domain"

type languageRepo interface {
	GetAll(edition string) ([]domain.Language, error)
	GetByCharacter(characterID uint) ([]domain.Language, error)
	Add(characterID, languageID uint) error
	Remove(characterID, languageID uint) error
}

type LanguageService struct{ repo languageRepo }

func NewLanguageService(repo languageRepo) *LanguageService { return &LanguageService{repo: repo} }

func (s *LanguageService) GetAll(edition string) ([]domain.Language, error) {
	return s.repo.GetAll(edition)
}
func (s *LanguageService) GetByCharacter(id uint) ([]domain.Language, error) {
	return s.repo.GetByCharacter(id)
}
func (s *LanguageService) Add(charID, languageID uint) error {
	return s.repo.Add(charID, languageID)
}
func (s *LanguageService) Remove(charID, languageID uint) error {
	return s.repo.Remove(charID, languageID)
}
