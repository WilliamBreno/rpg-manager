package domain

import "gorm.io/gorm"

// MagicItem — item/objeto de campanha que o mestre pode dar de recompensa
// (ver domain.Reward). Não é o mesmo catálogo de domain.Item (loja/inventário
// padrão 5e) — este é conteúdo específico da campanha do mestre (ex: um
// artefato único da história), por isso pertence a uma Campaign.
type MagicItem struct {
	gorm.Model
	CampaignID  uint     `json:"campaign_id"`
	Campaign    Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Effect      string   `json:"effect"`
}
