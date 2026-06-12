package domain

import "gorm.io/gorm"

// PowerType representa o tipo do poder conforme D&D
type PowerType string

const (
    PowerUtility   PowerType = "utility"    // Utilitário
    PowerUnlimited PowerType = "unlimited"  // Sem limite de uso
    PowerEncounter PowerType = "encounter"  // Por encontro
    PowerDaily     PowerType = "daily"      // Diário
)

type Skill struct {
    gorm.Model
    Name        string    `json:"name"`
    Description string    `json:"description"`
    PowerType   PowerType `json:"power_type"`
    Level       int       `json:"level"`    // Nível mínimo para desbloquear
    Edition     string    `json:"edition"`
    ClassID     *uint     `json:"class_id"` // Pode ser nulo (habilidade racial)
    RaceID      *uint     `json:"race_id"`  // Pode ser nulo (habilidade de classe)
}