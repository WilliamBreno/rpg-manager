package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type RitualRepository struct{ db *gorm.DB }

func NewRitualRepository(db *gorm.DB) *RitualRepository { return &RitualRepository{db: db} }

func (r *RitualRepository) FindByID(id uint) (domain.Ritual, error) {
	var ritual domain.Ritual
	err := r.db.First(&ritual, id).Error
	return ritual, err
}

func (r *RitualRepository) FindByName(name, edition string) (domain.Ritual, error) {
	var ritual domain.Ritual
	err := r.db.Where("name = ? AND edition = ?", name, edition).First(&ritual).Error
	return ritual, err
}

func (r *RitualRepository) GetAll(edition string) ([]domain.Ritual, error) {
	var rituals []domain.Ritual
	q := r.db.Order("level, name")
	if edition != "" {
		q = q.Where("edition = ?", edition)
	}
	return rituals, q.Find(&rituals).Error
}

func (r *RitualRepository) GetByCharacter(characterID uint) ([]domain.Ritual, error) {
	var rituals []domain.Ritual
	return rituals, r.db.
		Joins("JOIN character_rituals ON character_rituals.ritual_id = rituals.id").
		Where("character_rituals.character_id = ?", characterID).
		Order("rituals.level, rituals.name").
		Find(&rituals).Error
}

func (r *RitualRepository) Add(characterID, ritualID uint) error {
	return r.db.Exec(
		"INSERT INTO character_rituals (character_id, ritual_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		characterID, ritualID,
	).Error
}

func (r *RitualRepository) Remove(characterID, ritualID uint) error {
	return r.db.Exec(
		"DELETE FROM character_rituals WHERE character_id = ? AND ritual_id = ?",
		characterID, ritualID,
	).Error
}
