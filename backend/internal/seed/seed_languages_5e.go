package seed

import (
	"log"

	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

// seedLanguages5e povoa o catálogo de idiomas 5e — extraído verbatim do
// Livro do Jogador 2024, cap. 2, seção "Escolha Idiomas" (tabelas "Idiomas
// Comuns" e "Idiomas Raros"). RAW 2024, confirmado com o usuário antes de
// implementar (decisão explícita, não suposição): idioma NÃO vem da
// raça/espécie — todo personagem sabe Comum (concedido automaticamente, ver
// CharacterService.Create) e escolhe livremente mais 2 da tabela "Idiomas
// Comuns" na criação (ver CharacterCreate.tsx). A tabela "Idiomas Raros" só é
// concedida por classe/característica específica (ex: Druida ganha Druídico
// automaticamente via a característica "Idioma Druídico" já seedada como
// Skill) — não é uma escolha livre na criação, por isso não aparece no
// seletor, só existe no catálogo pra granting futuro por feature.
func seedLanguages5e(db *gorm.DB) {
	type l struct {
		Name, Category, Origin string
	}
	languages := []l{
		// Idiomas Comuns (Comum é automático; os outros 9 são a escolha livre)
		{"Comum", "comum", "Sigil"},
		{"Língua de Sinais Comum", "comum", "Sigil"},
		{"Dracônico", "comum", "Dragões"},
		{"Anão", "comum", "Anões"},
		{"Élfico", "comum", "Elfos"},
		{"Gigante", "comum", "Gigantes"},
		{"Gnômico", "comum", "Gnomos"},
		{"Goblin", "comum", "Goblinoides"},
		{"Pequenino", "comum", "Pequeninos"},
		{"Orc", "comum", "Orcs"},
		// Idiomas Raros (só por feature, não por escolha livre na criação)
		{"Abissal", "raro", "Demônios do Abismo"},
		{"Celestial", "raro", "Celestiais"},
		{"Dialeto Obscuro", "raro", "Aberrações"},
		{"Druídico", "raro", "Círculos druídicos"},
		{"Gíria dos Ladrões", "raro", "Várias guildas criminosas"},
		{"Infernal", "raro", "Diabos dos Nove Infernos"},
		{"Primordial", "raro", "Elementais"},
		{"Silvestre", "raro", "A Faéria"},
		{"Subcomum", "raro", "A Umbraeterna"},
	}
	for _, lang := range languages {
		var existing domain.Language
		if db.Where("name = ? AND edition = ?", lang.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Language{Name: lang.Name, Edition: "5e", Category: lang.Category, Origin: lang.Origin})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"category": lang.Category, "origin": lang.Origin,
			})
		}
	}
	log.Println("  ✓ Idiomas 5e seedados:", len(languages))
}
