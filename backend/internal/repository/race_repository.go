package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type RaceRepository struct {
    DB *gorm.DB
}

func NewRaceRepository(db *gorm.DB) *RaceRepository {
    return &RaceRepository{DB: db}
}

func (r *RaceRepository) FindAll() ([]domain.Race, error) {
    var races []domain.Race
    result := r.DB.Find(&races)
    return races, result.Error
}

func (r *RaceRepository) FindByID(id uint) (domain.Race, error) {
    var race domain.Race
    result := r.DB.First(&race, id)
    return race, result.Error
}

func (r *RaceRepository) FindByEdition(edition string) ([]domain.Race, error) {
    var races []domain.Race
    result := r.DB.Where("edition = ?", edition).Find(&races)
    return races, result.Error
}

func (r *RaceRepository) Create(race *domain.Race) error {
    return r.DB.Create(race).Error
}

func (r *RaceRepository) Update(race *domain.Race) error {
    return r.DB.Save(race).Error
}

func (r *RaceRepository) Delete(id uint) error {
    return r.DB.Delete(&domain.Race{}, id).Error
}