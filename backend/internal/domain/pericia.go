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
	// Expertise (Especialização) — Ladino (nível 1 e 6), Bardo (nível 2 e 9) e
	// Guardião (nível 9) dobram o Bônus de Proficiência em perícias já
	// proficientes à escolha. Só pode ser true numa perícia que já esteja
	// nesta mesma tabela (proficiente) — ver ExpertiseSlotsFor no service.
	Expertise bool `json:"expertise"`
}