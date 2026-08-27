package service

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// BackgroundService gerencia a biografia/notas do personagem
// (rota: GET/POST /characters/:id/background)
type BackgroundService struct {
	db *gorm.DB
}

func NewBackgroundService(db *gorm.DB) *BackgroundService {
	return &BackgroundService{db: db}
}

// allowedBackgroundFields: história/personalidade + descrição física —
// tudo que /characters/:id/background pode ler/escrever. Também usado pela
// exportação de PDF 5e (página 2 da ficha: idade/altura/peso/olhos/pele/
// cabelos/história).
var allowedBackgroundFields = []string{
	"personality_traits", "ideals", "bonds", "flaws", "rumors",
	"age", "height", "weight", "eyes", "skin", "hair", "history",
}

// Get retorna os campos de biografia do personagem
func (s *BackgroundService) Get(characterID uint) (map[string]interface{}, error) {
	var char domain.Character
	if err := s.db.
		Select("id", "personality_traits", "ideals", "bonds", "flaws", "rumors",
			"age", "height", "weight", "eyes", "skin", "hair", "history").
		First(&char, characterID).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"personality_traits": char.PersonalityTraits,
		"ideals":             char.Ideals,
		"bonds":              char.Bonds,
		"flaws":              char.Flaws,
		"rumors":             char.Rumors,
		"age":                char.Age,
		"height":             char.Height,
		"weight":             char.Weight,
		"eyes":               char.Eyes,
		"skin":               char.Skin,
		"hair":               char.Hair,
		"history":            char.History,
	}, nil
}

// Save atualiza os campos de biografia do personagem
func (s *BackgroundService) Save(characterID uint, data map[string]interface{}) error {
	updates := map[string]interface{}{}
	for _, key := range allowedBackgroundFields {
		if val, ok := data[key]; ok {
			updates[key] = val
		}
	}
	return s.db.Model(&domain.Character{}).Where("id = ?", characterID).Updates(updates).Error
}
