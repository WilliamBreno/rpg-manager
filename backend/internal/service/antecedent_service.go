package service

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type AntecedentService struct {
	db *gorm.DB
}

func NewAntecedentService(db *gorm.DB) *AntecedentService {
	return &AntecedentService{db: db}
}

func (s *AntecedentService) GetAll(edition string) ([]domain.Background, error) {
	var backgrounds []domain.Background
	result := s.db.Where("edition = ?", edition).Find(&backgrounds)
	return backgrounds, result.Error
}

func (s *AntecedentService) GetByID(id uint) (*domain.Background, error) {
	var background domain.Background
	result := s.db.First(&background, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &background, nil
}