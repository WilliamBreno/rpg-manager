package repository

import (
	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

type ItemRepository struct{ DB *gorm.DB }

func NewItemRepository(db *gorm.DB) *ItemRepository { return &ItemRepository{DB: db} }

func (r *ItemRepository) FindAll(edition, category string) ([]domain.Item, error) {
	var items []domain.Item
	q := r.DB.Where("edition = ?", edition).Order("category, name")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	return items, q.Find(&items).Error
}

func (r *ItemRepository) FindByID(id uint) (domain.Item, error) {
	var item domain.Item
	return item, r.DB.First(&item, id).Error
}
