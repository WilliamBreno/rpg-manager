package domain

import "gorm.io/gorm"

type Character struct {
	gorm.Model
	Name    string `json:"name"`
	Edition string `json:"edition"`
	Level   int    `json:"level"`

	HitPoints int `json:"hit_points"`
	MaxHP     int `json:"max_hp"`
	TempHP    int `json:"temp_hp"`

	SurgeValue   int `json:"surge_value"`
	SurgesPerDay int `json:"surges_per_day"`

	// Defesas 4e
	Defense_AC   int `json:"defense_ac"`
	Defense_Fort int `json:"defense_fort"`
	Defense_Refl int `json:"defense_refl"`
	Defense_Will int `json:"defense_will"`

	ArmorID   *uint  `json:"armor_id"`
	UserID    uint   `json:"user_id"`
	ClassID   uint   `json:"class_id"`
	RaceID    uint   `json:"race_id"`
	AvatarURL string `json:"avatar_url"`

	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`

	// ── Campos exclusivos do D&D 5e ──────────────────────────────────────────
	BackgroundID      *uint  `json:"background_id"`       // FK para Background
	Alignment         string `json:"alignment"`            // ex: "Leal e Bom"
	PersonalityTraits string `json:"personality_traits"`   // traços de personalidade
	Ideals            string `json:"ideals"`               // ideais
	Bonds             string `json:"bonds"`                // ligações
	Flaws             string `json:"flaws"`                // defeitos
	Speed             int    `json:"speed"`                // deslocamento (pés)
	ProficiencyBonus  int    `json:"proficiency_bonus"`    // bônus de proficiência calculado

	// ── Relacionamentos ───────────────────────────────────────────────────────
	Class      Class              `json:"class"      gorm:"foreignKey:ClassID"`
	Race       Race               `json:"race"       gorm:"foreignKey:RaceID"`
	Armor      *Armor             `json:"armor"      gorm:"foreignKey:ArmorID"`
	Background *Background        `json:"background" gorm:"foreignKey:BackgroundID"`  // ← corrigido
	Skills     []Skill            `json:"skills"     gorm:"many2many:character_skills;"`
	Pericias   []CharacterPericia `json:"pericias"   gorm:"foreignKey:CharacterID"`
	Talentos   []Talento          `json:"talentos"   gorm:"many2many:character_talentos;"`
}