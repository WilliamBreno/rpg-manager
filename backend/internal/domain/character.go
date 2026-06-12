package domain

import "gorm.io/gorm"

type Character struct {
    gorm.Model
    Name      string  `json:"name"`
    Edition   string  `json:"edition"`
    Level     int     `json:"level"`
    HitPoints int     `json:"hit_points"`
    MaxHP     int     `json:"max_hp"`
    TempHP    int     `json:"temp_hp"`
    ArmorID   *uint   `json:"armor_id"`
    UserID    uint    `json:"user_id"`
    ClassID   uint    `json:"class_id"`
    RaceID    uint    `json:"race_id"`
    AvatarURL string  `json:"avatar_url"`

    // Atributos base
    Strength     int `json:"strength"`
    Dexterity    int `json:"dexterity"`
    Constitution int `json:"constitution"`
    Intelligence int `json:"intelligence"`
    Wisdom       int `json:"wisdom"`
    Charisma     int `json:"charisma"`

    // Relacionamentos
    Class      Class       `json:"class" gorm:"foreignKey:ClassID"`
    Race       Race        `json:"race" gorm:"foreignKey:RaceID"`
    Armor      *Armor      `json:"armor" gorm:"foreignKey:ArmorID"`
    Skills     []Skill     `json:"skills" gorm:"many2many:character_skills;"`
    Background *Background `json:"background" gorm:"foreignKey:CharacterID"`
}