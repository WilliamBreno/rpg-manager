package service

import "rpg-manager/internal/domain"

type periciaRepo interface {
	GetAll(edition string) ([]domain.Pericia, error)
	GetByCharacter(characterID uint) ([]domain.CharacterPericia, error)
	Save(characterID uint, names []string) error
}

type PericiaService struct{ repo periciaRepo }

func NewPericiaService(repo periciaRepo) *PericiaService { return &PericiaService{repo: repo} }

func (s *PericiaService) GetAll(edition string) ([]domain.Pericia, error) {
	return s.repo.GetAll(edition)
}
func (s *PericiaService) GetByCharacter(id uint) ([]domain.CharacterPericia, error) {
	return s.repo.GetByCharacter(id)
}
func (s *PericiaService) Save(id uint, names []string) error {
	return s.repo.Save(id, names)
}