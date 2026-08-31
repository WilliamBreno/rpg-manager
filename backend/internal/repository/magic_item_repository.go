package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type MagicItemRepository struct{ DB *gorm.DB }

func NewMagicItemRepository(db *gorm.DB) *MagicItemRepository { return &MagicItemRepository{DB: db} }

func (r *MagicItemRepository) Create(m *domain.MagicItem) error { return r.DB.Create(m).Error }

func (r *MagicItemRepository) FindByCampaign(campaignID uint) ([]domain.MagicItem, error) {
	var items []domain.MagicItem
	err := r.DB.Where("campaign_id = ?", campaignID).Order("name").Find(&items).Error
	return items, err
}

func (r *MagicItemRepository) FindByID(id uint) (domain.MagicItem, error) {
	var m domain.MagicItem
	err := r.DB.First(&m, id).Error
	return m, err
}

func (r *MagicItemRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.MagicItem{}, id).Error
}
