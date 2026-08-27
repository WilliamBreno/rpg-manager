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

	// Boatos/rumores que circulam sobre o personagem (edição-agnóstico,
	// editável via /characters/:id/background junto com os outros campos
	// de biografia).
	Rumors string `json:"rumors"`

	// Descrição física e história — edição-agnóstico, editável pelo mesmo
	// endpoint /characters/:id/background. Usado na página 2 da ficha 5e
	// exportada (campos Age/Height/Weight/Eyes/Skin/Hair/Backstory).
	Age     string `json:"age"`
	Height  string `json:"height"`
	Weight  string `json:"weight"`
	Eyes    string `json:"eyes"`
	Skin    string `json:"skin"`
	Hair    string `json:"hair"`
	History string `json:"history"`

	// AbilityBonusChoice: só usado no payload de criação (não persistido —
	// gorm:"-"). Regra 2024: o bônus de atributo vem do Antecedente, nunca
	// da raça. O jogador digita os 6 atributos BASE (sem o bônus) e escolhe
	// aqui como distribuir os pontos entre as opções de
	// Antecedent.AbilityBonusOptions: +2/+1 em duas delas, ou +1 nas três.
	// Ex: {"INT": 2, "SAB": 1}. Ver CharacterService.Create.
	AbilityBonusChoice map[string]int `json:"ability_bonus_choice" gorm:"-"`

	// OriginFeatChoiceID: só usado no payload de criação (gorm:"-"), só
	// quando o antecedente é IsLegacy (livro antigo) e por isso não tem um
	// OriginFeatName fixo — o jogador escolhe livremente qualquer Talento
	// da categoria "Origem" (RAW 2024, caixa "Antecedentes e Espécies de
	// Livros Antigos"). Para antecedentes 2024 normais isso é ignorado, o
	// talento vem de Antecedent.OriginFeatName.
	OriginFeatChoiceID *uint `json:"origin_feat_choice_id" gorm:"-"`

	// EquipmentOptionID: só usado no payload de criação (gorm:"-") — a letra
	// (A/B/C) do "Equipamento Inicial" da classe (ClassEquipmentOption) que o
	// jogador escolheu. Opcional: se não informado, o personagem simplesmente
	// não recebe nenhum item/armadura/PO de partida (mesmo comportamento de
	// antes desta feature existir). Ver CharacterService.Create.
	EquipmentOptionID *uint `json:"equipment_option_id" gorm:"-"`

	// ── Relacionamentos ───────────────────────────────────────────────────────
	Class      Class              `json:"class"       gorm:"foreignKey:ClassID"`
	Race       Race               `json:"race"        gorm:"foreignKey:RaceID"`
	Armor      *Armor             `json:"armor"       gorm:"foreignKey:ArmorID"`
	Antecedent *Antecedent        `json:"antecedent"  gorm:"foreignKey:AntecedentID"`
	Background *Background        `json:"background"  gorm:"foreignKey:BackgroundID"`
	Skills     []Skill            `json:"skills"      gorm:"many2many:character_skills;"`
	Pericias   []CharacterPericia `json:"pericias"    gorm:"foreignKey:CharacterID"`
	Talentos   []Talento          `json:"talentos"    gorm:"many2many:character_talentos;"`
	Spells     []Spell            `json:"spells"      gorm:"many2many:character_spells;"`
	
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