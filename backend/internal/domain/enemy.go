package domain

import "gorm.io/gorm"

// EnemyKind distingue Inimigo comum, Boss e Vilão dentro da MESMA tabela —
// os três compartilham praticamente todos os campos (PV, raça, foto, som,
// classe, CA, habilidades); Boss/Vilão só acrescentam nome de destaque,
// história/vínculos/observações e falas. Uma flag IsBoss bool não bastaria
// porque Vilão é um terceiro tipo, não uma variação booleana de Boss — um
// enum de 3 valores evita duplicar ~10 campos em 3 structs separados (pedido
// explícito do SISTEMA_MESTRE.md: "sem duplicar campos à toa").
type EnemyKind string

const (
	EnemyKindEnemy   EnemyKind = "enemy"
	EnemyKindBoss    EnemyKind = "boss"
	EnemyKindVillain EnemyKind = "villain"
)

// Enemy — inimigo/Boss/Vilão de uma Campaign. Os campos narrativos
// (History/Bonds/Notes) e a lista de Lines só fazem sentido pra Boss/Vilão na
// prática, mas ficam disponíveis pra qualquer Kind em vez de exigir uma
// segunda tabela — um Enemy comum simplesmente não os preenche.
type Enemy struct {
	gorm.Model
	CampaignID uint      `json:"campaign_id"`
	Campaign   Campaign  `json:"campaign" gorm:"foreignKey:CampaignID"`
	Kind       EnemyKind `json:"kind" gorm:"default:enemy"`

	Name            string         `json:"name"`
	HP              int            `json:"hp"`
	ChallengeRating string         `json:"challenge_rating"` // ND opcional, ex: "1/4", "3" — usado só pra sugerir/avisar faixa de dano de EnemyAbility (ver service.CRDamageTable)
	Race            string         `json:"race"`
	PhotoURL        string         `json:"photo_url"`
	SoundURL        string         `json:"sound_url"` // som que o inimigo faz (não é fala — ver EnemyLine)
	Class           string         `json:"class"`     // opcional
	Armor           int            `json:"armor"`     // CA
	Abilities       []EnemyAbility `json:"abilities" gorm:"foreignKey:EnemyID"`

	// Só usados por Boss/Vilão na prática (vazios pra Enemy comum):
	History string      `json:"history,omitempty"`
	Bonds   string      `json:"bonds,omitempty"`
	Notes   string      `json:"notes,omitempty"`
	Lines   []EnemyLine `json:"lines,omitempty" gorm:"foreignKey:EnemyID"`
}

// EnemyAbility — habilidade customizada de combate. Damage segue notação
// real de dado (XdY+Z, ex: "2d6+3") — validada no service, não é um número
// solto (ver CLAUDE.md / PROGRESSAO_DE_NIVEL.md, mesmo cuidado já aplicado
// nas outras auditorias de dado real de D&D).
type EnemyAbility struct {
	gorm.Model
	EnemyID     uint   `json:"enemy_id"`
	Name        string `json:"name"`
	Damage      string `json:"damage"` // ex: "2d6+3"
	Description string `json:"description"`
}

// EnemyLineSource distingue se a fala foi gravada pelo mestre (upload) ou
// gerada por TTS — as duas vias coexistem, uma não substitui a outra (ver
// Etapa 2 do SISTEMA_MESTRE.md).
type EnemyLineSource string

const (
	EnemyLineUpload EnemyLineSource = "upload"
	EnemyLineTTS    EnemyLineSource = "tts"
)

// EnemyLine — uma fala de Boss/Vilão (texto + referência ao áudio).
type EnemyLine struct {
	gorm.Model
	EnemyID  uint            `json:"enemy_id"`
	Text     string          `json:"text"`
	AudioURL string          `json:"audio_url"`
	Source   EnemyLineSource `json:"source"`
}
