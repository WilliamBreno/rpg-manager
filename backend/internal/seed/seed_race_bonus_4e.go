package seed

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// seedRaceBonusesComplete4e aplica bônus raciais de perícia/talento para as
// 17 raças do 4e (PHB + LJ1 + LJ2 + LJ3). Portado de pkg/config (código
// morto) — esse pacote não é mais executado desde a correção do SSL, mas
// os dados nele eram corretos. Esta função é a versão ativa.
//
// Roda DEPOIS de seedRaces4eLJ23 (para LJ2/LJ3 já existirem) e depois de
// qualquer outra função que também escreva em bonus_talentos/bonus_trained_skills,
// para garantir que estes valores sejam os finais.
func seedRaceBonusesComplete4e(db *gorm.DB) {
	type entry struct {
		Name               string
		BonusTrainedSkills int
		BonusTalentos      int
		BonusSkillValues   string
	}
	races := []entry{
		{"Humano", 1, 1, `{}`},
		{"Meio-Elfo", 1, 0, `{"Diplomacia": 2, "Insight": 2}`},
		{"Anão", 0, 0, `{"Dungeon": 2, "Endurance": 2}`},
		{"Draconato", 0, 0, `{"Blefar": 2, "Intimidação": 2}`},
		{"Eladrin", 0, 0, `{"Arcana": 2, "História": 2}`},
		{"Elfo", 0, 0, `{"Natureza": 2, "Percepção": 2}`},
		{"Halfling", 0, 0, `{"Acrobacia": 2, "Furtividade": 2}`},
		{"Tiefling", 0, 0, `{"Blefar": 2, "Furtividade": 2}`},
		{"Meio-Orc", 0, 0, `{"Intimidação": 2, "Rua": 2}`},
		{"Deva", 0, 0, `{"História": 2, "Religião": 2}`},
		{"Gnomo", 0, 0, `{"Arcana": 2, "Furtividade": 2}`},
		{"Golias", 0, 0, `{"Atletismo": 2, "Natureza": 2}`},
		{"Feral", 0, 0, `{"Acrobacia": 2, "Atletismo": 2}`},
		{"Fragmental", 0, 0, `{"Arcana": 2, "Dungeon": 2}`},
		{"Githzerai", 0, 0, `{"Acrobacia": 2, "Atletismo": 2}`},
		{"Minotauro", 0, 0, `{"Natureza": 2, "Intimidação": 2}`},
		{"Sélvio", 0, 0, `{"Natureza": 2, "Furtividade": 2}`},
	}

	encontradas := 0
	for _, e := range races {
		var race domain.Race
		if db.Where("name = ? AND edition = ?", e.Name, "4e").First(&race).Error == nil {
			db.Model(&race).Updates(map[string]interface{}{
				"bonus_trained_skills": e.BonusTrainedSkills,
				"bonus_talentos":       e.BonusTalentos,
				"bonus_skill_values":   e.BonusSkillValues,
			})
			encontradas++
		} else {
			log.Printf("  ✗ Raça %s 4e não encontrada para aplicar bônus", e.Name)
		}
	}
	log.Printf("  ✓ Bônus raciais 4e completos: %d/%d raças atualizadas (Humano: +1 talento, +1 perícia)", encontradas, len(races))
}

// fixClassTalentosCount4e garante que toda classe 4e conceda exatamente
// 1 talento no nível 1 (regra: "1 talento por nível", a menos que raça
// conceda bônus — ex: Humano +1 via bonus_talentos).
//
// Roda por último: corrige qualquer valor 2 (ou outro) que outra função
// de seed possa ter setado, garantindo a regra correta independente da
// ordem de execução do Run().
func fixClassTalentosCount4e(db *gorm.DB) {
	result := db.Model(&domain.Class{}).Where("edition = ?", "4e").Update("talentos_count", 1)
	log.Printf("  ✓ Talentos por nível corrigido: %d classes 4e ajustadas para 1 talento base", result.RowsAffected)
}