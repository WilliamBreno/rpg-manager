package domain

import "gorm.io/gorm"

// Scene (Battleground) — um cenário de mapa da biblioteca de uma Campaign.
// Existe independente de qualquer Session, pra dar pro mestre montar a
// biblioteca inteira antes da sessão começar (pedido explícito do usuário).
type Scene struct {
	gorm.Model
	CampaignID uint     `json:"campaign_id"`
	Campaign   Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`
	Name       string   `json:"name"`
	ImageURL   string   `json:"image_url"`
	Tokens     []Token  `json:"tokens,omitempty" gorm:"foreignKey:SceneID"`
}

// Token — posição de um elemento decorativo/manual sobre uma Scene. Escopo
// desta rodada é só o modelo de dados (posição no grid) pra não travar o
// trabalho futuro de token 2D jogável com movimento por turno — a lógica de
// mover/agir NÃO é construída agora (limite de escopo já combinado). EnemyID/
// NPCID ficam prontos pra quando a importação automática de inimigo/NPC pro
// mapa for implementada (também adiada) — por enquanto um Token pode existir
// sem nenhum dos dois (só imagem+label, colocado manualmente pelo mestre).
type Token struct {
	gorm.Model
	SceneID  uint    `json:"scene_id"`
	Scene    Scene   `json:"scene" gorm:"foreignKey:SceneID"`
	Label    string  `json:"label"`
	ImageURL string  `json:"image_url"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	EnemyID  *uint   `json:"enemy_id,omitempty"`
	NPCID    *uint   `json:"npc_id,omitempty"`
}
