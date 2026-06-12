package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type ArmorRepository struct {
    DB *gorm.DB
}

func NewArmorRepository(db *gorm.DB) *ArmorRepository {
    return &ArmorRepository{DB: db}
}

func (r *ArmorRepository) FindAll() ([]domain.Armor, error) {
    var armors []domain.Armor
    result := r.DB.Find(&armors)
    return armors, result.Error
}

func (r *ArmorRepository) FindByEdition(edition string) ([]domain.Armor, error) {
    var armors []domain.Armor
    result := r.DB.Where("edition = ?", edition).Find(&armors)
    return armors, result.Error
}

func (r *ArmorRepository) FindByID(id uint) (domain.Armor, error) {
    var armor domain.Armor
    result := r.DB.First(&armor, id)
    return armor, result.Error
}