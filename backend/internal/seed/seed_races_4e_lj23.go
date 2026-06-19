package seed

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// seedRaces4eLJ23 cria os registros base das raças dos livros LJ2 e LJ3
// que nunca foram criadas no banco. As habilidades raciais (Memória de
// Mil Vidas, Desvanecer, Resistência de Pedra, etc.) já estão escritas
// em seedRaceSkills (seed.go) mas falhavam silenciosamente porque a
// raça base não existia — "✗ Raça X 4e não encontrada".
//
// Deve rodar ANTES de seedRaceSkills() no Run().
func seedRaces4eLJ23(db *gorm.DB) {
	type data struct {
		Name, Description string
		Speed             int
	}

	races := []data{
		// LJ2
		{"Deva", "Humanoide angelical reencarnado, portador de memórias fragmentadas de vidas passadas.", 30},
		{"Gnomo", "Pequeno humanoide com afinidade natural à magia ilusória e à furtividade.", 25},
		{"Golias", "Humanoide descendente das montanhas, de força e resistência extraordinárias.", 30},
		{"Meio-Orc", "Mistura de humano e orc, combinando fúria selvagem com resistência física.", 30},
		{"Feral", "Humanoide com herança bestial latente que desperta forças primais sob estresse.", 30},
		// LJ3
		{"Fragmental", "Ser composto por fragmentos elementais instáveis, capaz de se dispersar e se reformar.", 30},
		{"Githzerai", "Humanoide ascético oriundo do plano astral, com disciplina mental implacável.", 30},
		{"Minotauro", "Humanoide taurino de força brutal e instinto direto para o combate.", 30},
		{"Sélvio", "Ser feérico ligado a espíritos selvagens, capaz de manifestar diferentes aspectos.", 30},
	}

	for _, r := range races {
		var existing domain.Race
		if db.Where("name = ? AND edition = ?", r.Name, "4e").First(&existing).Error != nil {
			db.Create(&domain.Race{
				Name:        r.Name,
				Edition:     "4e",
				Description: r.Description,
				Speed:       r.Speed,
				IsDefault:   true,
			})
		}
	}
	log.Println("  ✓ Raças LJ2/LJ3 4e: 9 raças criadas (Deva, Gnomo, Golias, Meio-Orc, Feral, Fragmental, Githzerai, Minotauro, Sélvio)")
}