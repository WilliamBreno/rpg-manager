package domain

import "gorm.io/gorm"

type Race struct {
    gorm.Model
    Name        string  `json:"name"`
    Edition     string  `json:"edition"`
    Description string  `json:"description"`
    Speed       int     `json:"speed"`
    IsDefault   bool    `json:"is_default" gorm:"default:false"`
    Skills      []Skill `json:"skills" gorm:"foreignKey:RaceID"`
}