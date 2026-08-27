package repository

import (
	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

type ClassEquipmentRepository struct{ DB *gorm.DB }

func NewClassEquipmentRepository(db *gorm.DB) *ClassEquipmentRepository {
	return &ClassEquipmentRepository{DB: db}
}

func (r *ClassEquipmentRepository) GetByClass(classID uint) ([]domain.ClassEquipmentOption, error) {
	var options []domain.ClassEquipmentOption
	err := r.DB.Preload("Components.Item").Preload("Components.Armor").
		Where("class_id = ?", classID).Order("option_label").Find(&options).Error
	return options, err
}

func (r *ClassEquipmentRepository) FindByID(id uint) (domain.ClassEquipmentOption, error) {
	var option domain.ClassEquipmentOption
	err := r.DB.Preload("Components.Item").Preload("Components.Armor").First(&option, id).Error
	return option, err
}
