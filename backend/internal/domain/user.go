package domain

import "gorm.io/gorm"

type UserRole string

const (
    RolePlayer UserRole = "player"
    RoleMaster UserRole = "master"
)

type User struct {
    gorm.Model
    Name        string      `json:"name"`
    Email       string      `json:"email" gorm:"uniqueIndex"`
    Password    string      `json:"-"`
    Role        UserRole    `json:"role" gorm:"default:player"`
    WelcomeSeen bool        `json:"welcome_seen" gorm:"default:false"`
    Characters  []Character `json:"characters" gorm:"foreignKey:UserID"`
}