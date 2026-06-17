package domain

import "gorm.io/gorm"

type Race struct {
	gorm.Model
	Name        string `json:"name"`
	Edition     string `json:"edition"`
	Description string `json:"description"`
	Speed       int    `json:"speed"`

	// Perícias (4e)
	BonusTrainedSkills int    `json:"bonus_trained_skills"`
	BonusSkillValues   string `json:"bonus_skill_values"`

	// Talentos (4e)
	BonusTalentos int `json:"bonus_talentos"` 

	IsDefault bool    `json:"is_default" gorm:"default:false"`
	Skills    []Skill `json:"skills" gorm:"foreignKey:RaceID"`
}