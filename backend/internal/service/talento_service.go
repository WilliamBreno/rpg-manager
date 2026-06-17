package service

import "rpg-manager/internal/domain"

type talentoRepo interface {
	GetAll(edition string) ([]domain.Talento, error)
	GetByCharacter(characterID uint) ([]domain.Talento, error)
	Add(characterID, talentoID uint) error
	Remove(characterID, talentoID uint) error
}

type TalentoService struct{ repo talentoRepo }

func NewTalentoService(repo talentoRepo) *TalentoService { return &TalentoService{repo: repo} }

func (s *TalentoService) GetAll(edition string) ([]domain.Talento, error) {
	return s.repo.GetAll(edition)
}
func (s *TalentoService) GetByCharacter(id uint) ([]domain.Talento, error) {
	return s.repo.GetByCharacter(id)
}
func (s *TalentoService) Add(charID, talentoID uint) error {
	return s.repo.Add(charID, talentoID)
}
func (s *TalentoService) Remove(charID, talentoID uint) error {
	return s.repo.Remove(charID, talentoID)
}