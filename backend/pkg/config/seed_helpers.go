package config

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// upsertSkill cria ou atualiza uma Skill de classe no banco.
func upsertSkill(db *gorm.DB, s domain.Skill, classID uint) {
	var existing domain.Skill
	if db.Where("name = ? AND edition = ? AND class_id = ?", s.Name, s.Edition, classID).First(&existing).Error != nil {
		if err := db.Create(&s).Error; err != nil {
			log.Printf("  Erro ao criar skill %s: %v", s.Name, err)
		}
	} else {
		db.Model(&existing).Updates(map[string]interface{}{
			"description":      s.Description,
			"keywords":         s.Keywords,
			"action_type":      s.ActionType,
			"range":            s.Range,
			"target":           s.Target,
			"attack":           s.Attack,
			"hit":              s.Hit,
			"miss":             s.Miss,
			"effect":           s.Effect,
			"special":          s.Special,
			"level_scaling":    s.LevelScaling,
			"is_class_feature": s.IsClassFeature,
			"requires_choice":  s.RequiresChoice,
			"choice_group":     s.ChoiceGroup,
		})
	}
}

// upsertRaceSkill cria ou atualiza uma Skill racial no banco.
func upsertRaceSkill(db *gorm.DB, s domain.Skill, raceID uint) {
	var existing domain.Skill
	if db.Where("name = ? AND edition = ? AND race_id = ?", s.Name, s.Edition, raceID).First(&existing).Error != nil {
		if err := db.Create(&s).Error; err != nil {
			log.Printf("  Erro ao criar racial skill %s: %v", s.Name, err)
		}
	} else {
		db.Model(&existing).Updates(map[string]interface{}{
			"description":     s.Description,
			"keywords":        s.Keywords,
			"action_type":     s.ActionType,
			"range":           s.Range,
			"target":          s.Target,
			"attack":          s.Attack,
			"hit":             s.Hit,
			"miss":            s.Miss,
			"effect":          s.Effect,
			"special":         s.Special,
			"level_scaling":   s.LevelScaling,
			"power_type":      s.PowerType,
			"is_race_feature": s.IsRaceFeature,
			"requires_choice": s.RequiresChoice,
			"choice_group":    s.ChoiceGroup,
		})
	}
}