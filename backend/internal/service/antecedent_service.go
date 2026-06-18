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

func (s *AntecedentService) GetAll(edition string) ([]domain.Antecedent, error) {
    var antecedents []domain.Antecedent
    result := s.db.Where("edition = ?", edition).Find(&antecedents)
    return antecedents, result.Error
}

func (s *AntecedentService) GetByID(id uint) (*domain.Antecedent, error) {
    var antecedent domain.Antecedent
    result := s.db.First(&antecedent, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &antecedent, nil
}