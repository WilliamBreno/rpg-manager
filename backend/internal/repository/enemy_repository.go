package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type EnemyRepository struct{ DB *gorm.DB }

func NewEnemyRepository(db *gorm.DB) *EnemyRepository { return &EnemyRepository{DB: db} }

func (r *EnemyRepository) Create(e *domain.Enemy) error { return r.DB.Create(e).Error }

func (r *EnemyRepository) FindByCampaign(campaignID uint) ([]domain.Enemy, error) {
	var enemies []domain.Enemy
	err := r.DB.Preload("Abilities").Preload("Lines").
		Where("campaign_id = ?", campaignID).Order("name").Find(&enemies).Error
	return enemies, err
}

func (r *EnemyRepository) FindByID(id uint) (domain.Enemy, error) {
	var e domain.Enemy
	err := r.DB.Preload("Abilities").Preload("Lines").First(&e, id).Error
	return e, err
}

func (r *EnemyRepository) Update(e *domain.Enemy) error { return r.DB.Save(e).Error }

func (r *EnemyRepository) Delete(id uint) error { return r.DB.Delete(&domain.Enemy{}, id).Error }

func (r *EnemyRepository) CreateAbility(a *domain.EnemyAbility) error { return r.DB.Create(a).Error }

func (r *EnemyRepository) FindAbilityByID(id uint) (domain.EnemyAbility, error) {
	var a domain.EnemyAbility
	err := r.DB.First(&a, id).Error
	return a, err
}

func (r *EnemyRepository) UpdateAbility(a *domain.EnemyAbility) error { return r.DB.Save(a).Error }

func (r *EnemyRepository) DeleteAbility(id uint) error {
	return r.DB.Delete(&domain.EnemyAbility{}, id).Error
}

func (r *EnemyRepository) CreateLine(l *domain.EnemyLine) error { return r.DB.Create(l).Error }

func (r *EnemyRepository) FindLineByID(id uint) (domain.EnemyLine, error) {
	var l domain.EnemyLine
	err := r.DB.First(&l, id).Error
	return l, err
}

func (r *EnemyRepository) UpdateLine(l *domain.EnemyLine) error { return r.DB.Save(l).Error }

func (r *EnemyRepository) DeleteLine(id uint) error {
	return r.DB.Delete(&domain.EnemyLine{}, id).Error
}
