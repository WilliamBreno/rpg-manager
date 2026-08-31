package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type NPCRepository struct{ DB *gorm.DB }

func NewNPCRepository(db *gorm.DB) *NPCRepository { return &NPCRepository{DB: db} }

func (r *NPCRepository) Create(n *domain.NPC) error { return r.DB.Create(n).Error }

func (r *NPCRepository) FindByCampaign(campaignID uint) ([]domain.NPC, error) {
	var npcs []domain.NPC
	err := r.DB.Where("campaign_id = ?", campaignID).Order("name").Find(&npcs).Error
	return npcs, err
}

func (r *NPCRepository) FindByID(id uint) (domain.NPC, error) {
	var n domain.NPC
	err := r.DB.First(&n, id).Error
	return n, err
}

func (r *NPCRepository) Update(n *domain.NPC) error { return r.DB.Save(n).Error }

func (r *NPCRepository) Delete(id uint) error { return r.DB.Delete(&domain.NPC{}, id).Error }
