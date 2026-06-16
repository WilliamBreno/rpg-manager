package domain

import "gorm.io/gorm"

type PowerType string

const (
	PowerUtility   PowerType = "utility"
	PowerUnlimited PowerType = "unlimited"
	PowerEncounter PowerType = "encounter"
	PowerDaily     PowerType = "daily"
)

type Skill struct {
	gorm.Model

	Name    string    `json:"name"`
	Edition string    `json:"edition"`
	ClassID *uint     `json:"class_id"`
	RaceID  *uint     `json:"race_id"`

	PowerType PowerType `json:"power_type"`
	Level     int       `json:"level"`

	// Palavras-chave secundárias ex: "Arcano, Arma"
	Keywords string `json:"keywords"`

	// Mecânicas
	ActionType string `json:"action_type"`
	Range      string `json:"range"`
	Target     string `json:"target"`
	Attack     string `json:"attack"`

	// Efeitos
	Description  string `json:"description"`
	Hit          string `json:"hit"`
	Miss         string `json:"miss"`
	Effect       string `json:"effect"`
	Special      string `json:"special"`
	LevelScaling string `json:"level_scaling"`

	// Características de classe
	IsClassFeature bool   `json:"is_class_feature"`
	RequiresChoice bool   `json:"requires_choice"`
	ChoiceGroup    string `json:"choice_group"`
}