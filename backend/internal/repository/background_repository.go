package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type BackgroundRepository struct {
    DB *gorm.DB
}

func NewBackgroundRepository(db *gorm.DB) *BackgroundRepository {
    return &BackgroundRepository{DB: db}
}

func (r *BackgroundRepository) FindByCharacterID(characterID uint) (domain.Background, error) {
    var background domain.Background
    result := r.DB.Where("character_id = ?", characterID).First(&background)
    return background, result.Error
}

func (r *BackgroundRepository) Create(background *domain.Background) error {
    return r.DB.Create(background).Error
}

func (r *BackgroundRepository) Update(background *domain.Background) error {
    return r.DB.Save(background).Error
}

func (r *BackgroundRepository) Delete(characterID uint) error {
    return r.DB.Where("character_id = ?", characterID).Delete(&domain.Background{}).Error
}