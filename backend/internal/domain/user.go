package domain

import "gorm.io/gorm"

type UserRole string

const (
	RolePlayer UserRole = "player"
	RoleMaster UserRole = "master"
)

type User struct {
	gorm.Model
	Name        string   `json:"name"`
	Email       string   `json:"email" gorm:"uniqueIndex"`
	Password    string   `json:"-"`
	Role        UserRole `json:"role" gorm:"default:player"`
	WelcomeSeen bool     `json:"welcome_seen" gorm:"default:false"`
	// MasterWelcomeSeen — tela de boas-vindas específica de Mestre (Sistema do
	// Mestre), mostrada uma única vez na primeira visita a "Minhas Campanhas",
	// independente de WelcomeSeen (que é a boas-vindas geral de jogador).
	MasterWelcomeSeen bool        `json:"master_welcome_seen" gorm:"default:false"`
	Characters        []Character `json:"characters" gorm:"foreignKey:UserID"`
}
