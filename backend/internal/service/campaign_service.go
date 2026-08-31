package service

import (
	"errors"

	"rpg-manager/internal/domain"
)

type campaignRepo interface {
	Create(c *domain.Campaign) error
	FindByMaster(masterID uint) ([]domain.Campaign, error)
	FindByID(id uint) (domain.Campaign, error)
	Update(c *domain.Campaign) error
	Delete(id uint) error
}

type CampaignService struct{ repo campaignRepo }

func NewCampaignService(repo campaignRepo) *CampaignService { return &CampaignService{repo: repo} }

func (s *CampaignService) Create(c *domain.Campaign) error {
	if c.Name == "" {
		return errors.New("nome da campanha é obrigatório")
	}
	if c.Edition == "" {
		c.Edition = "5e"
	}
	return s.repo.Create(c)
}

func (s *CampaignService) GetByMaster(masterID uint) ([]domain.Campaign, error) {
	return s.repo.FindByMaster(masterID)
}

func (s *CampaignService) GetByID(id uint) (domain.Campaign, error) {
	return s.repo.FindByID(id)
}

// Update só edita nome/história — MasterID nunca muda por essa rota (a
// checagem de dono já acontece no handler antes de chamar isso).
func (s *CampaignService) Update(id uint, name, mainStory string) (domain.Campaign, error) {
	campaign, err := s.repo.FindByID(id)
	if err != nil {
		return campaign, err
	}
	if name != "" {
		campaign.Name = name
	}
	campaign.MainStory = mainStory
	return campaign, s.repo.Update(&campaign)
}

func (s *CampaignService) Delete(id uint) error {
	return s.repo.Delete(id)
}
