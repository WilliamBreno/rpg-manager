package domain

import (
	"time"

	"gorm.io/gorm"
)

// Session — uma sessão de jogo de uma Campaign. EndedAt nil enquanto a
// sessão está em andamento. ActiveSceneID é o cenário que o mestre definiu
// como ativo agora (o que os jogadores conectados veem — ver Etapa 5).
type Session struct {
	gorm.Model
	CampaignID    uint       `json:"campaign_id"`
	Campaign      Campaign   `json:"campaign" gorm:"foreignKey:CampaignID"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	Summary       string     `json:"summary"`             // diário de bordo, editável pelo mestre
	MusicURL      string     `json:"music_url,omitempty"` // música de fundo da sessão (Etapa 9) — URL simples, upload real depende de Cloudinary
	ActiveSceneID *uint      `json:"active_scene_id,omitempty"`
	ActiveScene   *Scene     `json:"active_scene,omitempty" gorm:"foreignKey:ActiveSceneID"`
}
