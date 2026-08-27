package service

import "rpg-manager/internal/domain"

type spellRepo interface {
	GetAll(edition string) ([]domain.Spell, error)
	GetByCharacter(characterID uint) ([]domain.Spell, error)
	Add(characterID, spellID uint) error
	Remove(characterID, spellID uint) error
}

type SpellService struct{ repo spellRepo }

func NewSpellService(repo spellRepo) *SpellService { return &SpellService{repo: repo} }

func (s *SpellService) GetAll(edition string) ([]domain.Spell, error) {
	return s.repo.GetAll(edition)
}
func (s *SpellService) GetByCharacter(id uint) ([]domain.Spell, error) {
	return s.repo.GetByCharacter(id)
}
func (s *SpellService) Add(charID, spellID uint) error {
	return s.repo.Add(charID, spellID)
}
func (s *SpellService) Remove(charID, spellID uint) error {
	return s.repo.Remove(charID, spellID)
}
