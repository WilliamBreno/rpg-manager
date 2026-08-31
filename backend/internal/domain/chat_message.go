package domain

import "gorm.io/gorm"

// ChatMessage — mensagem de chat de uma Campaign, entregue em tempo real via
// WebSocket (ver Etapa 5/6 do SISTEMA_MESTRE.md). Serve tanto pro chat
// mestre↔jogadores quanto jogador↔jogador dentro da mesma campanha — não há
// distinção de "canal", todo mundo na campanha vê a mesma linha do tempo.
// SessionID é opcional só pra registrar em qual sessão a mensagem ocorreu
// (referência histórica) — o chat funciona mesmo sem sessão ativa.
type ChatMessage struct {
	gorm.Model
	CampaignID   uint     `json:"campaign_id"`
	Campaign     Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`
	SessionID    *uint    `json:"session_id,omitempty"`
	Session      *Session `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	SenderUserID uint     `json:"sender_user_id"`
	Sender       User     `json:"sender" gorm:"foreignKey:SenderUserID"`
	Text         string   `json:"text"`
}
