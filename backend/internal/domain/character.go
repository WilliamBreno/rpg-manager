package domain

import "gorm.io/gorm"

type Character struct {
	gorm.Model
	Name    string `json:"name"`
	Edition string `json:"edition"`
	Level   int    `json:"level"`

	// XP — novo campo para o sistema de progressão
	ExperiencePoints int `json:"experience_points"`

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
	AntecedentID      *uint  `json:"antecedent_id"`       // FK para Antecedent (D&D 5e)
	BackgroundID      *uint  `json:"background_id"`        // FK para Background (biografia)
	Alignment         string `json:"alignment"`
	PersonalityTraits string `json:"personality_traits"`
	Ideals            string `json:"ideals"`
	Bonds             string `json:"bonds"`
	Flaws             string `json:"flaws"`
	Speed             int    `json:"speed"`
	ProficiencyBonus  int    `json:"proficiency_bonus"`

	// ── Relacionamentos ───────────────────────────────────────────────────────
	Class      Class              `json:"class"       gorm:"foreignKey:ClassID"`
	Race       Race               `json:"race"        gorm:"foreignKey:RaceID"`
	Armor      *Armor             `json:"armor"       gorm:"foreignKey:ArmorID"`
	Antecedent *Antecedent        `json:"antecedent"  gorm:"foreignKey:AntecedentID"`
	Background *Background        `json:"background"  gorm:"foreignKey:BackgroundID"`
	Skills     []Skill            `json:"skills"      gorm:"many2many:character_skills;"`
	Pericias   []CharacterPericia `json:"pericias"    gorm:"foreignKey:CharacterID"`
	Talentos   []Talento          `json:"talentos"    gorm:"many2many:character_talentos;"`
	
	// Death Saving Throws (5e)
	DeathSaveSuccesses int `json:"death_save_successes"`
	DeathSaveFailures  int `json:"death_save_failures"`

	// ── Moedas (5e) — as 5 moedas oficiais, cada uma independente (não são
	// "trocadas" automaticamente entre si) ────────────────────────────────────
	CopperPieces   int `json:"copper_pieces"`
	SilverPieces   int `json:"silver_pieces"`
	ElectrumPieces int `json:"electrum_pieces"`
	GoldPieces     int `json:"gold_pieces"`
	PlatinumPieces int `json:"platinum_pieces"`

	// ── Inventário (5e) ──────────────────────────────────────────────────────
	Items  []CharacterItem       `json:"items"  gorm:"foreignKey:CharacterID"`
	Armors []CharacterArmorOwned `json:"armors" gorm:"foreignKey:CharacterID"`
}