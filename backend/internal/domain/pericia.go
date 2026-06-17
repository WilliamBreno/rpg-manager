package domain

type Pericia struct {
	ID          uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name"`
	Attribute   string `json:"attribute"`   // "Força", "Destreza", etc.
	Description string `json:"description"`
	Tooltip     string `json:"tooltip"`
	Edition     string `json:"edition"`
}

type CharacterPericia struct {
	CharacterID uint   `json:"character_id" gorm:"primaryKey;not null"`
	PericiaName string `json:"pericia_name" gorm:"primaryKey;size:100;not null"`
}