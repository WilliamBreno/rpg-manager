package repository

import (
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

type TalentoRepository struct{ db *gorm.DB }

func NewTalentoRepository(db *gorm.DB) *TalentoRepository { return &TalentoRepository{db: db} }

func (r *TalentoRepository) FindByID(id uint) (domain.Talento, error) {
	var talento domain.Talento
	err := r.db.First(&talento, id).Error
	return talento, err
}

func (r *TalentoRepository) FindByName(name, edition string) (domain.Talento, error) {
	var talento domain.Talento
	err := r.db.Where("name = ? AND edition = ?", name, edition).First(&talento).Error
	return talento, err
}

func (r *TalentoRepository) GetAll(edition string) ([]domain.Talento, error) {
	var talentos []domain.Talento
	q := r.db.Order("category, name")
	if edition != "" {
		q = q.Where("edition = ?", edition)
	}
	return talentos, q.Find(&talentos).Error
}

func (r *TalentoRepository) GetByCharacter(characterID uint) ([]domain.Talento, error) {
	var talentos []domain.Talento
	return talentos, r.db.
		Joins("JOIN character_talentos ON character_talentos.talento_id = talentos.id").
		Where("character_talentos.character_id = ?", characterID).
		Find(&talentos).Error
}

func (r *TalentoRepository) Add(characterID, talentoID uint) error {
	return r.db.Exec(
		"INSERT INTO character_talentos (character_id, talento_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		characterID, talentoID,
	).Error
}

func (r *TalentoRepository) Remove(characterID, talentoID uint) error {
	return r.db.Exec(
		"DELETE FROM character_talentos WHERE character_id = ? AND talento_id = ?",
		characterID, talentoID,
	).Error
}