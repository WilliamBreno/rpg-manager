package domain

import "gorm.io/gorm"

// NPC — personagem não-jogável de apoio/narrativo (não é um inimigo de
// combate — isso é domain.Enemy). Pertence a uma Campaign.
type NPC struct {
	gorm.Model
	CampaignID  uint     `json:"campaign_id"`
	Campaign    Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`
	Name        string   `json:"name"`
	HP          int      `json:"hp"`
	History     string   `json:"history"`
	Bonds       string   `json:"bonds"`     // vínculos
	Alignment   string   `json:"alignment"` // tendência
	Personality string   `json:"personality"`
	Notes       string   `json:"notes"` // observações livres do mestre
}
