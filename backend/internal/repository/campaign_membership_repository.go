package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type CampaignMembershipRepository struct{ DB *gorm.DB }

func NewCampaignMembershipRepository(db *gorm.DB) *CampaignMembershipRepository {
	return &CampaignMembershipRepository{DB: db}
}

func (r *CampaignMembershipRepository) Create(m *domain.CampaignMembership) error {
	return r.DB.Create(m).Error
}

func (r *CampaignMembershipRepository) FindByCampaign(campaignID uint) ([]domain.CampaignMembership, error) {
	var memberships []domain.CampaignMembership
	err := r.DB.Preload("User").Preload("Character").Where("campaign_id = ?", campaignID).Find(&memberships).Error
	return memberships, err
}

func (r *CampaignMembershipRepository) FindByCampaignAndUser(campaignID, userID uint) (domain.CampaignMembership, error) {
	var m domain.CampaignMembership
	err := r.DB.Where("campaign_id = ? AND user_id = ?", campaignID, userID).First(&m).Error
	return m, err
}

func (r *CampaignMembershipRepository) FindByUser(userID uint) ([]domain.CampaignMembership, error) {
	var memberships []domain.CampaignMembership
	err := r.DB.Preload("Campaign").Where("user_id = ?", userID).Find(&memberships).Error
	return memberships, err
}

func (r *CampaignMembershipRepository) FindByID(id uint) (domain.CampaignMembership, error) {
	var m domain.CampaignMembership
	err := r.DB.First(&m, id).Error
	return m, err
}

func (r *CampaignMembershipRepository) Update(m *domain.CampaignMembership) error {
	return r.DB.Save(m).Error
}
