package domain

import "gorm.io/gorm"

// Campaign — mesa/campanha de um Mestre (User.Role == RoleMaster). Edição
// fixa em "5e" por enquanto (campo já pronto pra 4e no futuro, sem expor a
// opção na UI ainda — ver Etapa 1 do SISTEMA_MESTRE.md).
type Campaign struct {
	gorm.Model
	Name      string `json:"name"`
	Edition   string `json:"edition"` // "5e" por enquanto
	MainStory string `json:"main_story"`
	MasterID  uint   `json:"master_id"`
	Master    User   `json:"master" gorm:"foreignKey:MasterID"`
}

type MembershipStatus string

const (
	MembershipInvited  MembershipStatus = "invited"
	MembershipAccepted MembershipStatus = "accepted"
	MembershipDeclined MembershipStatus = "declined"
)

// CampaignMembership — relação mestre↔jogador por campanha. O documento
// original menciona "CampaignInvite/CampaignMembership" como se fossem dois
// conceitos, mas são o mesmo relacionamento em estágios diferentes (convite
// pendente → aceito/recusado) — modelado como uma linha só com Status, em vez
// de duas tabelas, pra não duplicar a mesma relação mestre-jogador-campanha.
type CampaignMembership struct {
	gorm.Model
	CampaignID uint     `json:"campaign_id"`
	Campaign   Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`
	UserID     uint     `json:"user_id"` // o jogador convidado
	User       User     `json:"user" gorm:"foreignKey:UserID"`
	// CharacterID: qual personagem esse jogador está usando nesta campanha —
	// opcional porque o convite pode ser aceito antes de vincular um
	// personagem específico.
	CharacterID *uint            `json:"character_id,omitempty"`
	Character   *Character       `json:"character,omitempty" gorm:"foreignKey:CharacterID"`
	Status      MembershipStatus `json:"status" gorm:"default:invited"`
}
