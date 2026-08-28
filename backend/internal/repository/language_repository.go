package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type LanguageRepository struct{ db *gorm.DB }

func NewLanguageRepository(db *gorm.DB) *LanguageRepository { return &LanguageRepository{db: db} }

func (r *LanguageRepository) FindByID(id uint) (domain.Language, error) {
	var language domain.Language
	err := r.db.First(&language, id).Error
	return language, err
}

func (r *LanguageRepository) FindByName(name, edition string) (domain.Language, error) {
	var language domain.Language
	err := r.db.Where("name = ? AND edition = ?", name, edition).First(&language).Error
	return language, err
}

func (r *LanguageRepository) GetAll(edition string) ([]domain.Language, error) {
	var languages []domain.Language
	q := r.db.Order("category, name")
	if edition != "" {
		q = q.Where("edition = ?", edition)
	}
	return languages, q.Find(&languages).Error
}

func (r *LanguageRepository) GetByCharacter(characterID uint) ([]domain.Language, error) {
	var languages []domain.Language
	return languages, r.db.
		Joins("JOIN character_languages ON character_languages.language_id = languages.id").
		Where("character_languages.character_id = ?", characterID).
		Order("languages.category, languages.name").
		Find(&languages).Error
}

func (r *LanguageRepository) Add(characterID, languageID uint) error {
	return r.db.Exec(
		"INSERT INTO character_languages (character_id, language_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		characterID, languageID,
	).Error
}

func (r *LanguageRepository) Remove(characterID, languageID uint) error {
	return r.db.Exec(
		"DELETE FROM character_languages WHERE character_id = ? AND language_id = ?",
		characterID, languageID,
	).Error
}
