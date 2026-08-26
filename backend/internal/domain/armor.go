package domain

import "gorm.io/gorm"

type ArmorType string

const (
    ArmorNone   ArmorType = "none"
    ArmorLight  ArmorType = "light"
    ArmorMedium ArmorType = "medium"
    ArmorHeavy  ArmorType = "heavy"
    ArmorShield ArmorType = "shield"
)

type Armor struct {
    gorm.Model
    Name        string    `json:"name"`
    Edition     string    `json:"edition"`
    ArmorType   ArmorType `json:"armor_type"`
    BaseAC      int       `json:"base_ac"`
    MaxDexBonus int       `json:"max_dex_bonus"` // -1 = sem limite, 0 = sem bônus DEX
    IsDefault   bool      `json:"is_default" gorm:"default:false"`
    Description string    `json:"description"`
    Weight      string    `json:"weight"`      // texto livre, ex: "10 kg" (unidades variam no livro)
    CostCopper  int       `json:"cost_copper"` // preço canônico em peças de cobre (1 PO = 100, 1 PP = 10, 1 PL = 1000)
}