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

// Get retorna os campos de biografia do personagem
func (s *BackgroundService) Get(characterID uint) (map[string]interface{}, error) {
	var char domain.Character
	if err := s.db.
		Select("id", "personality_traits", "ideals", "bonds", "flaws").
		First(&char, characterID).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"personality_traits": char.PersonalityTraits,
		"ideals":             char.Ideals,
		"bonds":              char.Bonds,
		"flaws":              char.Flaws,
	}, nil
}

// Save atualiza os campos de biografia do personagem
func (s *BackgroundService) Save(characterID uint, data map[string]interface{}) error {
	allowed := []string{"personality_traits", "ideals", "bonds", "flaws"}
	updates := map[string]interface{}{}
	for _, key := range allowed {
		if val, ok := data[key]; ok {
			updates[key] = val
		}
	}
	return s.db.Model(&domain.Character{}).Where("id = ?", characterID).Updates(updates).Error
}