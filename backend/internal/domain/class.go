package domain

import "gorm.io/gorm"

type Class struct {
    gorm.Model
    Name        string  `json:"name"`
    Edition     string  `json:"edition"`
    Description string  `json:"description"`
    HitDie      int     `json:"hit_die"`      // 5e: d6, d8, d10, d12
    BaseHP      int     `json:"base_hp"`       // 4e: HP base da classe no nível 1
    HPPerLevel  int     `json:"hp_per_level"`  // 4e: HP ganho por nível
    IsDefault   bool    `json:"is_default" gorm:"default:false"`
    Skills      []Skill `json:"skills" gorm:"foreignKey:ClassID"`
}