package domain

import "gorm.io/gorm"

type Background struct {
    gorm.Model
    CharacterID       uint   `json:"character_id"`
    History           string `json:"history"`
    PersonalityTraits string `json:"personality_traits"`
    Ideals            string `json:"ideals"`
    Bonds             string `json:"bonds"`
    Flaws             string `json:"flaws"`
}