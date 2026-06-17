package domain

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name        string `json:"name"`
	Edition     string `json:"edition"`
	Description string `json:"description"`

	// 5e
	HitDie int `json:"hit_die"`

	// 4e
	BaseHP       int `json:"base_hp"`
	HPPerLevel   int `json:"hp_per_level"`
	SurgesPerDay int `json:"surges_per_day"`
	FortBonus    int `json:"fort_bonus"`
	ReflBonus    int `json:"refl_bonus"`
	WillBonus    int `json:"will_bonus"`

	// Perícias treinadas (4e)
	TrainedSkillsCount int    `json:"trained_skills_count"` // quantas pode treinar
	AvailableSkills    string `json:"available_skills"`     // JSON: ["Atletismo","Percepção"]

	// Talentos (4e)
	TalentosCount int `json:"talentos_count"` // padrão 2 na sua campanha

	IsDefault bool    `json:"is_default" gorm:"default:false"`
	Skills    []Skill `json:"skills" gorm:"foreignKey:ClassID"`
}