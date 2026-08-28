package domain

import "gorm.io/gorm"

// Language — catálogo de idiomas 5e (Livro do Jogador 2024, cap. 2, seção
// "Escolha Idiomas"). RAW 2024: idioma NÃO vem da raça/espécie — todo
// personagem sabe Comum e escolhe livremente mais 2 da tabela "Idiomas
// Comuns"; a tabela "Idiomas Raros" só é concedida por classe/característica
// específica, não por escolha livre na criação (ver seed_languages_5e.go).
type Language struct {
	gorm.Model
	Name     string `json:"name"`
	Edition  string `json:"edition"`
	Category string `json:"category"` // "comum" ou "raro"
	Origin   string `json:"origin"`   // ex: "Elfos", "Dragões" — flavor, não mecânico
}
