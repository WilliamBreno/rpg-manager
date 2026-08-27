package domain

import "gorm.io/gorm"

// Spell é o catálogo de magias do PHB 2024 (capítulo "Magias"), compartilhado
// entre todas as classes conjuradoras — uma mesma magia (ex.: Bola de Fogo)
// aparece na lista de várias classes com o mesmo círculo, então não é
// duplicada por classe: Classes carrega o JSON de quais classes a conhecem e
// em que círculo cada uma a adquire (nem sempre é o mesmo círculo pra
// magias… na prática no PHB 2024 é sempre o mesmo círculo entre classes,
// mas o formato já suporta o caso de não ser).
type Spell struct {
	gorm.Model
	Name          string `json:"name"`
	Edition       string `json:"edition"`
	Level         int    `json:"level"` // 0 = truque, 1-9 = círculo
	School        string `json:"school"`
	Ritual        bool   `json:"ritual"`
	Concentration bool   `json:"concentration"`
	// Classes: JSON de {"NomeDaClasse": nivel_circulo}, ex. {"Mago":3,"Feiticeiro":3}
	Classes     string `json:"classes"`
	Description string `json:"description"`
}
