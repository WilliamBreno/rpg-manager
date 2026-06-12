package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type ClassRepository struct {
    DB *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
    return &ClassRepository{DB: db}
}

func (r *CharacterRepository) FindAll(userID uint) ([]domain.Character, error) {
    var characters []domain.Character
    result := r.DB.Where("user_id = ?", userID).Preload("Class").Preload("Race").Preload("Skills").Find(&characters)
    return characters, result.Error
}

func (r *ClassRepository) FindByID(id uint) (domain.Class, error) {
    var class domain.Class
    result := r.DB.First(&class, id)
    return class, result.Error
}

func (r *ClassRepository) FindByEdition(edition string) ([]domain.Class, error) {
    var classes []domain.Class
    result := r.DB.Where("edition = ?", edition).Find(&classes)
    return classes, result.Error
}

func (r *ClassRepository) Create(class *domain.Class) error {
    return r.DB.Create(class).Error
}

func (r *ClassRepository) Update(class *domain.Class) error {
    return r.DB.Save(class).Error
}

func (r *ClassRepository) Delete(id uint) error {
    return r.DB.Delete(&domain.Class{}, id).Error
}