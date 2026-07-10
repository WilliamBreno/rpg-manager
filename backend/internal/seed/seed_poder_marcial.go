package seed

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// Poderes adicionais extraídos do sourcebook "Poder Marcial" (Martial Power) para
// as classes marciais: Guerreiro, Ladino e Patrulheiro. O capítulo do Comandante
// de Guerra (Warlord) foi ignorado pois essa classe ainda não existe no banco.
// Este arquivo é independente de seed.go e não é chamado automaticamente
// pelo seed.Run — a conexão deve ser feita manualmente.

func seedGuerreiroSkillsPoderMarcial(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Guerreiro", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Guerreiro 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Guerreiro 4e (Poder Marcial): poderes seedados")
}

func seedLadinoSkillsPoderMarcial(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Ladino", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Ladino 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Ladino 4e (Poder Marcial): poderes seedados")
}

func seedPatrulheiroSkillsPoderMarcial(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Patrulheiro", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Patrulheiro 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Patrulheiro 4e (Poder Marcial): poderes seedados")
}
