package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type SpellRepository struct{ db *gorm.DB }

func NewSpellRepository(db *gorm.DB) *SpellRepository { return &SpellRepository{db: db} }

func (r *SpellRepository) FindByID(id uint) (domain.Spell, error) {
	var spell domain.Spell
	err := r.db.First(&spell, id).Error
	return spell, err
}

func (r *SpellRepository) FindByName(name, edition string) (domain.Spell, error) {
	var spell domain.Spell
	err := r.db.Where("name = ? AND edition = ?", name, edition).First(&spell).Error
	return spell, err
}

// GetAll retorna magias filtradas por edição (sempre) e opcionalmente por
// nível exato — o filtro por classe é feito em memória pelo service, já que
// Classes é um JSON livre (uma coluna por classe explodiria o schema à toa
// pra uma lista que não muda com frequência).
func (r *SpellRepository) GetAll(edition string) ([]domain.Spell, error) {
	var spells []domain.Spell
	q := r.db.Order("level, name")
	if edition != "" {
		q = q.Where("edition = ?", edition)
	}
	return spells, q.Find(&spells).Error
}

func (r *SpellRepository) GetByCharacter(characterID uint) ([]domain.Spell, error) {
	var spells []domain.Spell
	return spells, r.db.
		Joins("JOIN character_spells ON character_spells.spell_id = spells.id").
		Where("character_spells.character_id = ?", characterID).
		Order("spells.level, spells.name").
		Find(&spells).Error
}

func (r *SpellRepository) Add(characterID, spellID uint) error {
	return r.db.Exec(
		"INSERT INTO character_spells (character_id, spell_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		characterID, spellID,
	).Error
}

func (r *SpellRepository) Remove(characterID, spellID uint) error {
	return r.db.Exec(
		"DELETE FROM character_spells WHERE character_id = ? AND spell_id = ?",
		characterID, spellID,
	).Error
}
