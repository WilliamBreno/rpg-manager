package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type RewardRepository struct{ DB *gorm.DB }

func NewRewardRepository(db *gorm.DB) *RewardRepository { return &RewardRepository{DB: db} }

func (r *RewardRepository) Create(reward *domain.Reward) error { return r.DB.Create(reward).Error }

func (r *RewardRepository) FindByCampaign(campaignID uint) ([]domain.Reward, error) {
	var rewards []domain.Reward
	err := r.DB.Preload("Character").Preload("GrantedBy").Preload("MagicItem").
		Where("campaign_id = ?", campaignID).Order("created_at desc").Find(&rewards).Error
	return rewards, err
}

func (r *RewardRepository) FindByCharacter(characterID uint) ([]domain.Reward, error) {
	var rewards []domain.Reward
	err := r.DB.Preload("MagicItem").Where("character_id = ?", characterID).Order("created_at desc").Find(&rewards).Error
	return rewards, err
}
