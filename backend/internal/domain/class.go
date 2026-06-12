package domain

import "gorm.io/gorm"

type Class struct {
    gorm.Model
    Name        string  `json:"name"`
    Edition     string  `json:"edition"`
    Description string  `json:"description"`
    HitDie      int     `json:"hit_die"`
    IsDefault   bool    `json:"is_default" gorm:"default:false"`
    Skills      []Skill `json:"skills" gorm:"foreignKey:ClassID"`
}