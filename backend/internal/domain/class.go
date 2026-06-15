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
    BaseHP      int `json:"base_hp"`       // HP base nível 1
    HPPerLevel  int `json:"hp_per_level"`  // HP ganho por nível
    SurgesPerDay int `json:"surges_per_day"` // Pulsos base por dia

    // Bônus de defesa (4e) — somados ao FORT/REFL/VONT da classe
    FortBonus int `json:"fort_bonus"`
    ReflBonus int `json:"refl_bonus"`
    WillBonus int `json:"will_bonus"`

    IsDefault bool    `json:"is_default" gorm:"default:false"`
    Skills    []Skill `json:"skills" gorm:"foreignKey:ClassID"`
}