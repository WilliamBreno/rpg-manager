package domain

import "gorm.io/gorm"

type Background struct {
	gorm.Model
	Name               string `json:"name"`
	Edition            string `json:"edition"`
	Description        string `json:"description"`
	SkillProficiencies string `json:"skill_proficiencies"` // JSON: ["Intuição","Religião"]
	ToolProficiencies  string `json:"tool_proficiencies"`  // ex: "Ferramentas de ladrão, 1 tipo de jogo"
	Languages          string `json:"languages"`           // ex: "Dois idiomas à sua escolha"
	Equipment          string `json:"equipment"`           // descrição do equipamento inicial
	Feature            string `json:"feature"`             // nome da característica especial
	FeatureDescription string `json:"feature_description"` // descrição completa
	IsDefault          bool   `json:"is_default" gorm:"default:false"`
}