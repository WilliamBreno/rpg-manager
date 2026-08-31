package service

import (
	"errors"

	"rpg-manager/internal/domain"
)

type npcRepo interface {
	Create(n *domain.NPC) error
	FindByCampaign(campaignID uint) ([]domain.NPC, error)
	FindByID(id uint) (domain.NPC, error)
	Update(n *domain.NPC) error
	Delete(id uint) error
}

type NPCService struct{ repo npcRepo }

func NewNPCService(repo npcRepo) *NPCService { return &NPCService{repo: repo} }

func (s *NPCService) Create(n *domain.NPC) error {
	if n.Name == "" {
		return errors.New("nome do NPC é obrigatório")
	}
	return s.repo.Create(n)
}

func (s *NPCService) GetByCampaign(campaignID uint) ([]domain.NPC, error) {
	return s.repo.FindByCampaign(campaignID)
}

func (s *NPCService) GetByID(id uint) (domain.NPC, error) { return s.repo.FindByID(id) }

func (s *NPCService) Update(n *domain.NPC) error { return s.repo.Update(n) }

func (s *NPCService) Delete(id uint) error { return s.repo.Delete(id) }
