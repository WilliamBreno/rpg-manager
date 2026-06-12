package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type CharacterRepository struct {
    DB *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *CharacterRepository {
    return &CharacterRepository{DB: db}
}

func (r *ClassRepository) FindAll() ([]domain.Class, error) {
    var classes []domain.Class
    result := r.DB.Find(&classes)
    return classes, result.Error
}

func (r *CharacterRepository) FindByID(id uint) (domain.Character, error) {
    var character domain.Character
    result := r.DB.Preload("Class").Preload("Race").Preload("Skills").First(&character, id)
    return character, result.Error
}

func (r *CharacterRepository) Create(character *domain.Character) error {
    return r.DB.Create(character).Error
}

func (r *CharacterRepository) Update(character *domain.Character) error {
    return r.DB.Save(character).Error
}

func (r *CharacterRepository) Delete(id uint) error {
    return r.DB.Delete(&domain.Character{}, id).Error
}

func (r *CharacterRepository) AddSkill(character *domain.Character, skill *domain.Skill) error {
    return r.DB.Model(character).Association("Skills").Append(skill)
}

func (r *CharacterRepository) RemoveSkill(character *domain.Character, skill *domain.Skill) error {
    return r.DB.Model(character).Association("Skills").Delete(skill)
}