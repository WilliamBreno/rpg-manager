package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type ChatMessageRepository struct{ DB *gorm.DB }

func NewChatMessageRepository(db *gorm.DB) *ChatMessageRepository {
	return &ChatMessageRepository{DB: db}
}

func (r *ChatMessageRepository) Create(m *domain.ChatMessage) error { return r.DB.Create(m).Error }

func (r *ChatMessageRepository) FindByID(id uint) (domain.ChatMessage, error) {
	var m domain.ChatMessage
	err := r.DB.Preload("Sender").First(&m, id).Error
	return m, err
}

func (r *ChatMessageRepository) FindByCampaign(campaignID uint, limit int) ([]domain.ChatMessage, error) {
	var messages []domain.ChatMessage
	err := r.DB.Preload("Sender").Where("campaign_id = ?", campaignID).
		Order("created_at desc").Limit(limit).Find(&messages).Error
	return messages, err
}
