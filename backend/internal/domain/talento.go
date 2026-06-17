package domain

import "gorm.io/gorm"

type Talento struct {
	gorm.Model
	Name         string `json:"name"`
	Edition      string `json:"edition"`
	Category     string `json:"category"`     // "Combate", "Defesa", "Perícia", "Magia", "Armadura"
	Description  string `json:"description"`
	Prerequisite string `json:"prerequisite"` // "" se não tiver
	Tooltip      string `json:"tooltip"`
}