package domain

import "gorm.io/gorm"

type RewardKind string

const (
	RewardCurrency RewardKind = "currency"
	RewardItem     RewardKind = "item"
)

// Reward (Transaction) — registro histórico de uma entrega do mestre pra um
// personagem (moeda ou MagicItem). "Dar a todos" na UI (Etapa 8) gera uma
// linha por personagem destinatário, não uma linha com múltiplos
// destinatários — mantém o modelo simples (1 linha = 1 entrega a 1
// personagem) e a soma/saldo de cada um sempre bate com o histórico dele.
type Reward struct {
	gorm.Model
	CampaignID      uint       `json:"campaign_id"`
	Campaign        Campaign   `json:"campaign" gorm:"foreignKey:CampaignID"`
	CharacterID     uint       `json:"character_id"`
	Character       Character  `json:"character" gorm:"foreignKey:CharacterID"`
	GrantedByUserID uint       `json:"granted_by_user_id"`
	GrantedBy       User       `json:"granted_by" gorm:"foreignKey:GrantedByUserID"`
	Kind            RewardKind `json:"kind"`

	// Preenchidos quando Kind == RewardCurrency (mesmo esquema de 5 moedas já
	// usado em Character — ver "Shop / equipment / currency" no CLAUDE.md):
	CopperPieces   int `json:"copper_pieces"`
	SilverPieces   int `json:"silver_pieces"`
	ElectrumPieces int `json:"electrum_pieces"`
	GoldPieces     int `json:"gold_pieces"`
	PlatinumPieces int `json:"platinum_pieces"`

	// Preenchido quando Kind == RewardItem:
	MagicItemID *uint      `json:"magic_item_id,omitempty"`
	MagicItem   *MagicItem `json:"magic_item,omitempty" gorm:"foreignKey:MagicItemID"`

	Note string `json:"note,omitempty"`
}
