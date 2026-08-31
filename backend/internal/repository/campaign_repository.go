package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type CampaignRepository struct{ DB *gorm.DB }

func NewCampaignRepository(db *gorm.DB) *CampaignRepository { return &CampaignRepository{DB: db} }

func (r *CampaignRepository) Create(c *domain.Campaign) error {
	return r.DB.Create(c).Error
}

func (r *CampaignRepository) FindByMaster(masterID uint) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.DB.Where("master_id = ?", masterID).Order("created_at desc").Find(&campaigns).Error
	return campaigns, err
}

func (r *CampaignRepository) FindByID(id uint) (domain.Campaign, error) {
	var c domain.Campaign
	err := r.DB.Preload("Master").First(&c, id).Error
	return c, err
}

func (r *CampaignRepository) Update(c *domain.Campaign) error {
	return r.DB.Save(c).Error
}

func (r *CampaignRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.Campaign{}, id).Error
}
