package domain

import "gorm.io/gorm"

type Antecedent struct {
    gorm.Model
    Name               string `json:"name"`
    Edition            string `json:"edition"`
    Description        string `json:"description"`
    SkillProficiencies string `json:"skill_proficiencies"`
    ToolProficiencies  string `json:"tool_proficiencies"`
    Languages          string `json:"languages"`
    Equipment          string `json:"equipment"`
    Feature            string `json:"feature"`
    FeatureDescription string `json:"feature_description"`
    IsDefault          bool   `json:"is_default" gorm:"default:false"`
}