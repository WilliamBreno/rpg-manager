package repository

import (
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

type PericiaRepository struct{ db *gorm.DB }

func NewPericiaRepository(db *gorm.DB) *PericiaRepository { return &PericiaRepository{db: db} }

func (r *PericiaRepository) GetAll(edition string) ([]domain.Pericia, error) {
	var pericias []domain.Pericia
	q := r.db.Order("name")
	if edition != "" {
		q = q.Where("edition = ?", edition)
	}
	return pericias, q.Find(&pericias).Error
}

func (r *PericiaRepository) GetByCharacter(characterID uint) ([]domain.CharacterPericia, error) {
	var pericias []domain.CharacterPericia
	return pericias, r.db.Where("character_id = ?", characterID).Find(&pericias).Error
}

func (r *PericiaRepository) Save(characterID uint, names []string) error {
	tx := r.db.Begin()
	if err := tx.Where("character_id = ?", characterID).Delete(&domain.CharacterPericia{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, name := range names {
		cp := domain.CharacterPericia{CharacterID: characterID, PericiaName: name}
		if err := tx.Create(&cp).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}