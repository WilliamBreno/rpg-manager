package seed

import (
	"log"
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	log.Println("🌱 Rodando seed...")
	seedClasses(db)
	seedClasses5e(db) 
	seedAntecedents5e(db)
	seedRaces4eLJ23(db)
	seedRaces5e(db)
	seedSkills(db)
	seedSkills5e(db)
	seedRaceSkills(db)
	seedRaceFeatures5e(db) 
	seedPericias(db)
    seedClassPericias(db)
	seedPericias5e(db) 
    seedRacePericias(db)
	seedRaceBonusesComplete4e(db)
	fixClassTalentosCount4e(db)
    seedTalentos(db)
    seedTalentos5e(db)
	seedArmors5e(db)
	seedItems5e(db)
	log.Println("✅ Seed concluído!")
}

func upsertSkill(db *gorm.DB, s domain.Skill, classID uint) {
	var existing domain.Skill
	if db.Where("name = ? AND edition = ? AND class_id = ?", s.Name, s.Edition, classID).First(&existing).Error != nil {
		if err := db.Create(&s).Error; err != nil {
			log.Printf("  Erro ao criar skill %s: %v", s.Name, err)
		}
	} else {
		db.Model(&existing).Updates(map[string]interface{}{
			"description": s.Description, "keywords": s.Keywords,
			"action_type": s.ActionType, "range": s.Range,
			"target": s.Target, "attack": s.Attack,
			"hit": s.Hit, "miss": s.Miss, "effect": s.Effect,
			"special": s.Special, "level_scaling": s.LevelScaling,
			"is_class_feature": s.IsClassFeature,
			"requires_choice": s.RequiresChoice,
			"choice_group": s.ChoiceGroup,
		})
	}
}

func seedClasses(db *gorm.DB) {
    type classData struct {
        Name              string
        Description       string
        BaseHP, HPPerLevel, SurgesPerDay int
        FortBonus, ReflBonus, WillBonus  int
        TrainedSkillsCount               int
        AvailableSkills                  string
        AutomaticPericias                string
    }

    classes := []classData{
        // ⚠️ Verificar cada linha contra o PDF LJ1/LJ2/LJ3
        {
            Name: "Bardo", Description: "Líder arcano que usa música e conhecimento como armas.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 4,
            AutomaticPericias:  `["Arcanismo"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Blefe","Diplomacia","Exploração","História","Intimidação","Intuição","Manha","Natureza","Percepção","Religião","Socorro"]`,
        },
        {
            Name: "Bárbaro", Description: "Combatente primitivo de força e fúria selvagem.",
            BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 8, FortBonus: 2,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `[]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Intimidação","Natureza","Percepção","Socorro","Tolerância"]`,
        },
        {
            Name: "Clérigo", Description: "Líder divino que canaliza o poder de sua divindade.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Religião"]`,
            AvailableSkills:    `["Diplomacia","Dungeon","História","Intuição","Natureza","Percepção","Socorro"]`,
        },
        {
            Name: "Druida", Description: "Controlador primitivo com maestria da natureza.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Natureza"]`,
            AvailableSkills:    `["Acrobacia","Arcanismo","Atletismo","Dungeon","Endurance","Exploração","História","Intuição","Percepção","Religião","Socorro"]`,
        },
        {
            Name: "Feiticeiro", Description: "Arcano que canaliza poder dracônico ou do caos.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 2,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Arcanismo"]`,
            AvailableSkills:    `["Blefe","Diplomacia","Dungeon","Endurance","História","Intuição","Natureza","Percepção"]`,
        },
        {
            Name: "Guardião", Description: "Defensor primitivo e protetor do mundo natural.",
            BaseHP: 17, HPPerLevel: 7, SurgesPerDay: 9, FortBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Natureza"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Dungeon","Endurance","Exploração","Intimidação","Percepção","Socorro","Tolerância"]`,
        },
        {
            Name: "Guerreiro", Description: "Defensor marcial especialista em combate corpo a corpo.",
            BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 9, FortBonus: 2,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `[]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Blefe","Dungeon","Endurance","Intimidação","Percepção","Rua","Tolerância"]`,
        },
        {
            Name: "Invocador", Description: "Controlador divino que invoca o poder dos deuses.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 2,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Religião"]`,
            AvailableSkills:    `["Arcanismo","Diplomacia","Dungeon","História","Intuição","Natureza","Percepção"]`,
        },
        {
            Name: "Ladino", Description: "Agressor furtivo especializado em ataques precisos.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, ReflBonus: 2,
            TrainedSkillsCount: 4,
            AutomaticPericias:  `["Destreza com Ladrão","Furtividade"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Blefe","Diplomacia","Dungeon","Intuição","Intimidação","Manha","Percepção","Rua"]`,
        },
        {
            Name: "Mago", Description: "Controlador arcano de grande poder mágico.",
            BaseHP: 10, HPPerLevel: 4, SurgesPerDay: 6, WillBonus: 2,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Arcanismo"]`,
            AvailableSkills:    `["Diplomacia","Dungeon","Endurance","História","Intuição","Natureza","Percepção","Religião"]`,
        },
        {
            Name: "Monge", Description: "Agressor psiônico de disciplina e artes marciais.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, FortBonus: 1, ReflBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 4,
            AutomaticPericias:  `["Religião"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Diplomacia","Dungeon","Endurance","Furtividade","Intuição","Percepção","Tolerância"]`,
        },
        {
            Name: "Paladino", Description: "Defensor divino e campeão sagrado.",
            BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 10, FortBonus: 1, ReflBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 2,
            AutomaticPericias:  `["Religião"]`,
            AvailableSkills:    `["Diplomacia","Dungeon","Endurance","História","Intuição","Percepção"]`,
        },
        {
            Name: "Patrulheiro", Description: "Agressor marcial e batedor das regiões ermas.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, FortBonus: 1, ReflBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Natureza","Exploração"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Dungeon","Endurance","Furtividade","Intimidação","Percepção","Socorro"]`,
        },
        {
            Name: "Psionista", Description: "Controlador psiônico de grande poder mental.",
            BaseHP: 12, HPPerLevel: 4, SurgesPerDay: 6, WillBonus: 2,
            TrainedSkillsCount: 4,
            AutomaticPericias:  `[]`,
            AvailableSkills:    `["Arcanismo","Atletismo","Blefe","Diplomacia","Dungeon","Endurance","História","Intuição","Percepção","Religião","Tolerância"]`,
        },
        {
            Name: "Rastreador", Description: "Controlador primitivo especialista em emboscadas.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Natureza"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Dungeon","Endurance","Exploração","Furtividade","Percepção","Socorro","Tolerância"]`,
        },
        {
            Name: "Vingador", Description: "Agressor divino executor da vontade dos deuses.",
            BaseHP: 14, HPPerLevel: 6, SurgesPerDay: 7, FortBonus: 1, ReflBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Religião"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Diplomacia","Dungeon","Endurance","Intimidação","Intuição","Percepção","Rua"]`,
        },
        {
            Name: "Xamã", Description: "Líder primitivo que canaliza espíritos do mundo natural.",
            BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, FortBonus: 1, WillBonus: 1,
            TrainedSkillsCount: 3,
            AutomaticPericias:  `["Natureza"]`,
            AvailableSkills:    `["Acrobacia","Atletismo","Diplomacia","Dungeon","Endurance","História","Intuição","Percepção","Religião","Socorro"]`,
        },
    }

    for _, c := range classes {
        var existing domain.Class
        if db.Where("name = ? AND edition = ?", c.Name, "4e").First(&existing).Error != nil {
            db.Create(&domain.Class{
                Name: c.Name, Edition: "4e", Description: c.Description,
                BaseHP: c.BaseHP, HPPerLevel: c.HPPerLevel, SurgesPerDay: c.SurgesPerDay,
                FortBonus: c.FortBonus, ReflBonus: c.ReflBonus, WillBonus: c.WillBonus,
                TrainedSkillsCount: c.TrainedSkillsCount,
                AvailableSkills:    c.AvailableSkills,
                AutomaticPericias:  c.AutomaticPericias,
                IsDefault: true,
            })
        } else {
            db.Model(&existing).Updates(map[string]interface{}{
                "base_hp": c.BaseHP, "hp_per_level": c.HPPerLevel,
                "surges_per_day": c.SurgesPerDay, "fort_bonus": c.FortBonus,
                "refl_bonus": c.ReflBonus, "will_bonus": c.WillBonus,
                "description":          c.Description,
                "trained_skills_count": c.TrainedSkillsCount,
                "available_skills":     c.AvailableSkills,
                "automatic_pericias":   c.AutomaticPericias, // ← NOVO
            })
        }
    }
    log.Println("  ✓ Classes seedadas")
}

func seedSkills(db *gorm.DB) {
	seedBardoSkills(db)
	seedMongeSkills(db)
	seedClerigoSkills(db)
	seedGuerreiroSkills(db)
	seedLadinoSkills(db)
	seedMagoSkills(db)
	seedPaladinoSkills(db)
	seedPatrulheiroSkills(db)
	seedBarbaroSkills(db)
	seedDruidaSkills(db)
	seedFeiticeiroSkills(db)
	seedGuardiaoSkills(db)
	seedInvocadorSkills(db)
	seedVingadorSkills(db)
	seedXamaSkills(db)
}
func seedPericias5e(db *gorm.DB) {
	type p struct {
		Name, Attribute, Description, Tooltip string
	}
 
	pericias := []p{
		{"Acrobacia",         "DES", "Manter equilíbrio, rolar, cair com estilo.",                        "Teste de Destreza para acrobacias e equilíbrio."},
		{"Arcanismo",         "INT", "Recordar conhecimentos sobre magia, itens e planos.",               "Teste de Inteligência para lembrar fatos arcanos."},
		{"Atletismo",         "FOR", "Saltar, escalar, nadar e superar obstáculos físicos.",              "Teste de Força para proezas atléticas."},
		{"Atuação",           "CAR", "Cantar, dançar, tocar, contar histórias.",                          "Teste de Carisma para se apresentar a uma audiência."},
		{"Enganação",         "CAR", "Mentir de forma convincente e usar disfarces.",                     "Teste de Carisma para iludir ou enganar."},
		{"Furtividade",       "DES", "Mover-se sem ser notado e se esconder.",                            "Teste de Destreza para se mover silenciosamente."},
		{"História",          "INT", "Relembrar eventos históricos, culturas e povos.",                   "Teste de Inteligência para lembrar fatos históricos."},
		{"Intimidação",       "CAR", "Aterrorizar ou coagir alguém pela força da personalidade.",         "Teste de Carisma para intimidar ou ameaçar."},
		{"Intuição",          "SAB", "Perceber o humor e as intenções de uma pessoa.",                    "Teste de Sabedoria para ler pessoas e situações."},
		{"Investigação",      "INT", "Encontrar pistas e deduzir como algo funciona.",                    "Teste de Inteligência para investigar e deduzir."},
		{"Lidar com Animais", "SAB", "Acalmar, treinar ou comandar animais.",                             "Teste de Sabedoria para interagir com animais."},
		{"Medicina",          "SAB", "Diagnosticar doenças e estabilizar ferimentos.",                    "Teste de Sabedoria para tratar ferimentos."},
		{"Natureza",          "INT", "Reconhecer plantas, animais, terrenos e fenômenos naturais.",       "Teste de Inteligência para lembrar fatos naturais."},
		{"Percepção",         "SAB", "Notar detalhes usando todos os sentidos.",                          "Teste de Sabedoria para detectar coisas ao redor."},
		{"Persuasão",         "CAR", "Convencer alguém de algo de forma honesta e graciosa.",             "Teste de Carisma para persuadir com tato."},
		{"Prestidigitação",   "DES", "Furtar bolsos, esconder objetos, truques de mão.",                  "Teste de Destreza para truques sutis com as mãos."},
		{"Religião",          "INT", "Relembrar fatos sobre deuses, rituais e símbolos sagrados.",        "Teste de Inteligência para lembrar conhecimento religioso."},
		{"Sobrevivência",     "SAB", "Rastrear, caçar, encontrar trilhas e evitar perigos naturais.",     "Teste de Sabedoria para sobreviver em ambientes hostis."},
	}
 
	for _, p := range pericias {
		var existing domain.Pericia
		if db.Where("name = ? AND edition = ?", p.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Pericia{
				Name:        p.Name,
				Attribute:   p.Attribute,
				Edition:     "5e",
				Description: p.Description,
				Tooltip:     p.Tooltip,
			})
		}
	}
	log.Println("  ✓ Perícias 5e: 18 perícias seedadas")
}
func seedClasses5e(db *gorm.DB) {
	type data struct {
		Name               string
		Description        string
		HitDie             int
		SavingThrows       string
		TrainedSkillsCount int
		AvailableSkills    string
	}

	todas := `["Acrobacia","Arcanismo","Atletismo","Atuação","Enganação","Furtividade","História","Intimidação","Intuição","Investigação","Lidar com Animais","Medicina","Natureza","Percepção","Persuasão","Prestidigitação","Religião","Sobrevivência"]`

	classes := []data{
		{
			Name:               "Bárbaro",
			Description:        "Guerreiro furioso que entra em Fúria para devastar inimigos e resistir a ferimentos.",
			HitDie:             12,
			SavingThrows:       `["FOR","CON"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Atletismo","Intimidação","Lidar com Animais","Natureza","Percepção","Sobrevivência"]`,
		},
		{
			Name:               "Bardo",
			Description:        "Conjurador versátil que usa música e palavras para inspirar aliados e confundir inimigos.",
			HitDie:             8,
			SavingThrows:       `["DES","CAR"]`,
			TrainedSkillsCount: 3,
			AvailableSkills:    todas,
		},
		{
			Name:               "Bruxo",
			Description:        "Conjurador que extrai poder arcano de um patrono misterioso.",
			HitDie:             8,
			SavingThrows:       `["SAB","CAR"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Arcanismo","Enganação","História","Intimidação","Investigação","Natureza","Religião"]`,
		},
		{
			Name:               "Clérigo",
			Description:        "Servo divino que canaliza poder sagrado para curar aliados e punir inimigos.",
			HitDie:             8,
			SavingThrows:       `["SAB","CAR"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["História","Intuição","Medicina","Persuasão","Religião"]`,
		},
		{
			Name:               "Druida",
			Description:        "Guardião da natureza que molda os elementos e assume formas animais.",
			HitDie:             8,
			SavingThrows:       `["INT","SAB"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Arcanismo","Lidar com Animais","Intuição","Medicina","Natureza","Percepção","Religião","Sobrevivência"]`,
		},
		{
			Name:               "Feiticeiro",
			Description:        "Conjurador de magia inata que molda o poder à sua própria natureza.",
			HitDie:             6,
			SavingThrows:       `["CON","CAR"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Arcanismo","Enganação","Intimidação","Intuição","Persuasão","Religião"]`,
		},
		{
			// PHB 2024: Guardião é o Ranger (não "Patrulheiro").
			// O stub "Patrulheiro 5e" já existente fica intacto; este é novo.
			Name:               "Guardião",
			Description:        "Guerreiro ágil que une proezas marciais, magia da natureza e sobrevivência.",
			HitDie:             10,
			SavingThrows:       `["FOR","DES"]`,
			TrainedSkillsCount: 3,
			AvailableSkills:    `["Atletismo","Furtividade","Intuição","Investigação","Lidar com Animais","Natureza","Percepção","Sobrevivência"]`,
		},
		{
			Name:               "Guerreiro",
			Description:        "Mestre do combate que domina todas as armas e armaduras.",
			HitDie:             10,
			SavingThrows:       `["FOR","CON"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Acrobacia","Atletismo","Furtividade","História","Intimidação","Intuição","Percepção","Sobrevivência"]`,
		},
		{
			Name:               "Ladino",
			Description:        "Especialista furtivo que desfere ataques mortais e resolve problemas com astúcia.",
			HitDie:             8,
			SavingThrows:       `["DES","INT"]`,
			TrainedSkillsCount: 4,
			AvailableSkills:    `["Acrobacia","Atletismo","Enganação","Furtividade","Intimidação","Investigação","Percepção","Atuação","Persuasão","Prestidigitação"]`,
		},
		{
			Name:               "Mago",
			Description:        "Estudioso da magia arcana que domina feitiços de todos os propósitos.",
			HitDie:             6,
			SavingThrows:       `["INT","SAB"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Arcanismo","História","Intuição","Investigação","Medicina","Religião"]`,
		},
		{
			Name:               "Monge",
			Description:        "Combatente de artes marciais que canaliza energia mística pelo corpo.",
			HitDie:             8,
			SavingThrows:       `["FOR","DES"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Acrobacia","Atletismo","Furtividade","História","Intuição","Religião"]`,
		},
		{
			Name:               "Paladino",
			Description:        "Guerreiro sagrado jurado a um ideal, combinando poder marcial e divino.",
			HitDie:             10,
			SavingThrows:       `["SAB","CAR"]`,
			TrainedSkillsCount: 2,
			AvailableSkills:    `["Atletismo","Intuição","Intimidação","Medicina","Persuasão","Religião"]`,
		},
	}

	for _, c := range classes {
		var existing domain.Class
		if db.Where("name = ? AND edition = ?", c.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Class{
				Name:               c.Name,
				Edition:            "5e",
				Description:        c.Description,
				HitDie:             c.HitDie,
				SavingThrows:       c.SavingThrows,
				TrainedSkillsCount: c.TrainedSkillsCount,
				AvailableSkills:    c.AvailableSkills,
				IsDefault:          true,
			})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description":          c.Description,
				"hit_die":              c.HitDie,
				"saving_throws":        c.SavingThrows,
				"trained_skills_count": c.TrainedSkillsCount,
				"available_skills":     c.AvailableSkills,
			})
		}
	}
	log.Println("  ✓ Classes 5e: 12 classes atualizadas com dados completos")
}

// seedSkills5e semeia as características de classe do nível 1 para o PHB 2024.
func seedSkills5e(db *gorm.DB) {
	seedBarbaro5e(db)
	seedBardo5e(db)
	seedBruxo5e(db)
	seedClerigo5e(db)
	seedDruida5e(db)
	seedFeiticeiro5e(db)
	seedGuardiao5e(db)
	seedGuerreiro5e(db)
	seedLadino5e(db)
	seedMago5e(db)
	seedMonge5e(db)
	seedPaladino5e(db)
}

// getClass5e busca uma classe 5e pelo nome e retorna o ID.
func getClass5e(db *gorm.DB, name string) (uint, bool) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", name, "5e").First(&cls).Error; err != nil {
		log.Printf("  ✗ Classe 5e não encontrada: %s (rode seedClasses5e primeiro)", name)
		return 0, false
	}
	return cls.ID, true
}

// ── Bárbaro ──────────────────────────────────────────────────────────────────

func seedBarbaro5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Bárbaro")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Fúria",
			Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus você entra em Fúria. Enquanto em Fúria: +2 de bônus nas jogadas de dano de Força, resistência a dano contundente/perfurante/cortante, vantagem em testes de Força. Dura 1 minuto, encerrando se você ficar Incapacitado ou parar de atacar. Usos: 2 — recupera em Descanso Curto ou Longo.",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:       "Fúria ativa por 1 minuto: +2 dano FOR, resistência a dano físico, vantagem em testes de FOR.",
			LevelScaling: "Nível 3: 3 usos. Nível 5: +2→+3 dano. Nível 6: 4 usos. Nível 9: 5 usos. Nível 12: 6 usos. Nível 17: 7 usos.",
			PowerType:    domain.PowerEncounter, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Defesa Sem Armadura",
			Edition: "5e", ClassID: &id,
			Description: "Enquanto não usar armadura, sua CA = 10 + modificador de Destreza + modificador de Constituição. Você ainda pode usar um escudo.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "CA = 10 + DES + CON sem armadura.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Maestria em Armas",
			Edition: "5e", ClassID: &id,
			Description: "Você pode usar as propriedades de Maestria de 3 tipos de armas Simples ou Marciais à sua escolha. Sempre que completar um Descanso Longo pode substituir uma das escolhas.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Desbloqueia propriedades de Maestria em 3 armas escolhidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todas as trilhas) ──
		{
			Name:    "Ataque Imprudente", Edition: "5e", ClassID: &id,
			Description: "Ao realizar sua primeira jogada de ataque no turno, pode decidir atacar de forma imprudente: Vantagem em jogadas de ataque usando Força até o início do seu próximo turno, mas jogadas de ataque contra você também têm Vantagem nesse período.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Troca Vantagem em ataques de Força por Vantagem também para inimigos atacarem você.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Sentido de Perigo", Edition: "5e", ClassID: &id,
			Description: "Você tem Vantagem em salvaguardas de Destreza, a menos que tenha a condição Incapacitado.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em salvaguardas de Destreza.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Conhecimento Primordial", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência em outra perícia disponível para Bárbaros. Enquanto em Fúria, pode realizar testes de Acrobacia, Furtividade, Intimidação, Percepção ou Sobrevivência como teste de Força em vez do atributo normal.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+1 perícia; usa Força no lugar do atributo normal em 5 perícias durante a Fúria.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true,
		},
		{
			Name:    "Ataque Extra", Edition: "5e", ClassID: &id,
			Description: "Você pode atacar duas vezes, em vez de uma, sempre que executar a ação Atacar no seu turno.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 ataques por ação Atacar.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Movimento Rápido", Edition: "5e", ClassID: &id,
			Description: "Seu Deslocamento aumenta em 3 metros enquanto você não estiver usando Armadura Pesada.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+3m de Deslocamento sem Armadura Pesada.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Bote Instintivo", Edition: "5e", ClassID: &id,
			Description: "Como parte da Ação Bônus que você realiza para entrar em Fúria, pode se mover até metade do seu Deslocamento.",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:    "Move metade do Deslocamento ao entrar em Fúria.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Instintos Primitivos", Edition: "5e", ClassID: &id,
			Description: "Você tem Vantagem nas jogadas de Iniciativa.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em Iniciativa.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Golpe Brutal", Edition: "5e", ClassID: &id,
			Description: "Ao usar Ataque Imprudente, pode renunciar à Vantagem em uma jogada de ataque com Força à sua escolha para causar 1d10 de dano adicional se acertar, e aplicar um efeito de Golpe Brutal: Golpe Debilitador (reduz o Deslocamento do alvo em 4,5m até o início do seu próximo turno) ou Golpe Poderoso (empurra o alvo 4,5m e você pode se mover em direção a ele sem provocar Ataques de Oportunidade).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "1d10 de dano extra + efeito de controle, custando a Vantagem de um ataque.",
			LevelScaling: "Nível 13: ganha os efeitos Golpe Atordoante e Golpe Destruidor. Nível 17: dano extra sobe para 2d10 e pode aplicar 2 efeitos de Golpe Brutal ao mesmo tempo.",
			PowerType:    domain.PowerUnlimited, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Fúria Implacável", Edition: "5e", ClassID: &id,
			Description: "Se cair a 0 PV enquanto sua Fúria estiver ativa e não morrer imediatamente, pode realizar uma salvaguarda de Constituição CD 10 para ficar com PV igual a duas vezes seu nível de Bárbaro. A CD aumenta em 5 a cada uso repetido, voltando a 10 num Descanso Curto ou Longo.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Chance de não cair a 0 PV enquanto em Fúria, via salvaguarda de Constituição.",
			PowerType: domain.PowerUnlimited, Level: 11, IsClassFeature: true,
		},
		{
			Name:    "Fúria Persistente", Edition: "5e", ClassID: &id,
			Description: "Ao jogar Iniciativa, pode recuperar todos os usos gastos de Fúria (1x por Descanso Longo). Sua Fúria também passa a durar 10 minutos sem precisar ser estendida rodada a rodada, encerrando apenas se ficar Inconsciente ou vestir Armadura Pesada.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recupera usos de Fúria ao rolar Iniciativa; Fúria dura 10 minutos sem precisar estender.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true,
		},
		{
			Name:    "Força Indomável", Edition: "5e", ClassID: &id,
			Description: "Se o total de um teste de Força ou salvaguarda de Força for menor que seu valor de Força, pode usar esse valor no lugar do resultado total.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Piso mínimo em testes/salvaguardas de Força igual ao seu valor de Força.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Campeão Primitivo", Edition: "5e", ClassID: &id,
			Description: "Seus valores de Força e Constituição aumentam em 4, até um máximo de 25.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+4 Força e +4 Constituição (máx. 25).",
			PowerType: domain.PowerUnlimited, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		{
			Name: "Trilha do Berserker", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Trilhas Primitivas de Bárbaro, escolhida no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "trilha_barbaro",
		},
		{
			Name: "Trilha do Coração Selvagem", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Trilhas Primitivas de Bárbaro, escolhida no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "trilha_barbaro",
		},
		{
			Name: "Trilha da Árvore do Mundo", Edition: "5e", ClassID: &id,
			Description: "Entrelace as Raízes e Ramos do Multiverso — bárbaros que seguem esta trilha conectam-se à árvore cósmica Yggdrasil por meio de sua Fúria, extraindo vitalidade e a capacidade de viajar entre dimensões.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "trilha_barbaro",
		},
		{
			Name: "Trilha do Fanático", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Trilhas Primitivas de Bárbaro, escolhida no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "trilha_barbaro",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/10/14 — ChoiceGroup aqui é o
		// Name exato da trilha escolhida, não o slug "trilha_barbaro"; ver
		// hasChosenSubclass em character_service.go) ───────────────────────
		{
			Name: "Vitalidade da Árvore", Edition: "5e", ClassID: &id,
			Description: "Enquanto em Fúria: no início de cada turno pode dar PV Temporários (Xd6, X = bônus de Dano da Fúria) a uma criatura a até 3m; e ao ativar a Fúria recebe PV Temporários igual ao seu nível de Bárbaro (Surto de Vitalidade).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "3 metros",
			Effect: "PV Temporários próprios e para aliados ao ativar/manter a Fúria.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trilha da Árvore do Mundo",
		},
		{
			Name: "Ramos da Árvore", Edition: "5e", ClassID: &id,
			Description: "Como Reação, quando uma criatura a até 9m começa o turno à vista enquanto em Fúria, convoca ramos espectrais: salvaguarda de Força (CD 8 + Bônus Prof + FOR) ou é teleportada para perto de você e pode ter o Deslocamento reduzido a 0 até o fim do turno.",
			Keywords: "Primitivo", ActionType: "Reação", Range: "9 metros",
			Effect: "Puxa e imobiliza um inimigo próximo, ou falha a salvaguarda.",
			PowerType: domain.PowerEncounter, Level: 6, IsClassFeature: true, ChoiceGroup: "Trilha da Árvore do Mundo",
		},
		{
			Name: "Raízes Devastadoras", Edition: "5e", ClassID: &id,
			Description: "No seu turno, armas corpo a corpo Pesadas ou Versáteis ganham +3m de alcance. Ao acertar com elas, pode ativar a maestria Derrubar ou Empurrar além de outra maestria que já esteja usando.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+3m de alcance corpo a corpo; maestria extra de arma.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Trilha da Árvore do Mundo",
		},
		{
			Name: "Percorrer a Árvore", Edition: "5e", ClassID: &id,
			Description: "Ao ativar a Fúria (ou como Ação Bônus enquanto ativa), teleporta-se até 18m para um espaço à vista. Uma vez por Fúria pode estender o alcance para 45m e levar até 6 criaturas voluntárias a até 3m consigo.",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "18-45 metros",
			Effect: "Teleporte pessoal (e em grupo, 1x por Fúria) durante a Fúria.",
			PowerType: domain.PowerEncounter, Level: 14, IsClassFeature: true, ChoiceGroup: "Trilha da Árvore do Mundo",
		},
		{
			Name: "Frenesi", Edition: "5e", ClassID: &id,
			Description: "Se usar Ataque Imprudente enquanto em Fúria, o primeiro alvo atingido no turno com um ataque de Força sofre dano adicional (Xd6, X = bônus de Dano da Fúria, mesmo tipo da arma/Ataque Desarmado).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano extra no primeiro ataque de Força do turno, combinado com Ataque Imprudente.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trilha do Berserker",
		},
		{
			Name: "Fúria Irracional", Edition: "5e", ClassID: &id,
			Description: "Enquanto em Fúria, tem Imunidade a Amedrontado e Enfeitiçado; se já estiver sob uma dessas condições ao entrar em Fúria, ela encerra.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Imunidade a Amedrontado/Enfeitiçado durante a Fúria.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Trilha do Berserker",
		},
		{
			Name: "Retaliação", Edition: "5e", ClassID: &id,
			Description: "Ao sofrer dano de uma criatura a até 1,5m, pode usar uma Reação para realizar um ataque corpo a corpo (arma ou Desarmado) contra ela.",
			Keywords: "Primitivo", ActionType: "Reação", Range: "1,5 metro",
			Effect: "Contra-ataque corpo a corpo como Reação.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Trilha do Berserker",
		},
		{
			Name: "Presença Intimidante", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, cada criatura à escolha numa Emanação de 9m faz salvaguarda de Sabedoria (CD 8 + FOR + Bônus Prof) ou fica Amedrontada por 1 minuto (repete a salvaguarda a cada fim de turno). 1 uso por Descanso Longo, ou gaste um uso de Fúria para restaurar.",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "Emanação 9 metros",
			Effect: "Amedronta em área; recarrega com Descanso Longo ou gastando um uso de Fúria.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true, ChoiceGroup: "Trilha do Berserker",
		},
		{
			Name: "Arauto da Fauna", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar Falar com Animais e Sentido Feral, apenas como Rituais, usando Sabedoria como atributo de conjuração.",
			Keywords: "Primitivo, Conjuração", ActionType: "Ritual", Range: "Pessoal",
			Effect: "2 magias rituais fixas.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trilha do Coração Selvagem",
		},
		{
			Name: "Fúria dos Selvagens", Edition: "5e", ClassID: &id,
			Description: "Ao ativar a Fúria, escolhe: Águia (pode executar Correr e Desengajar como parte da Ação Bônus de ativar a Fúria, e repetir ambas como Ação Bônus enquanto ativa), Lobo (aliados têm Vantagem contra inimigos a até 1,5m de você) ou Urso (Resistência a quase todos os tipos de dano, exceto Energético/Necrótico/Psíquico/Radiante).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Escolha de bônus (Águia/Lobo/Urso) toda vez que ativa a Fúria.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trilha do Coração Selvagem",
		},
		{
			Name: "Aspecto dos Selvagens", Edition: "5e", ClassID: &id,
			Description: "Escolhe um: Coruja (Visão no Escuro 18m, ou +18m se já tiver), Pantera (Deslocamento de Escalada igual ao normal) ou Salmão (Deslocamento de Natação igual ao normal). Pode trocar a escolha a cada Descanso Longo.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Traço sensorial/de movimento à escolha, trocável por Descanso Longo.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Trilha do Coração Selvagem",
		},
		{
			Name: "Arauto da Natureza", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar Comunhão com a Natureza, apenas como Ritual, usando Sabedoria como atributo de conjuração.",
			Keywords: "Primitivo, Conjuração", ActionType: "Ritual", Range: "Pessoal",
			Effect: "Mais 1 magia ritual fixa.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Trilha do Coração Selvagem",
		},
		{
			Name: "Poder dos Selvagens", Edition: "5e", ClassID: &id,
			Description: "Ao ativar a Fúria, escolhe: Carneiro (impõe Caído em criaturas Grandes ou menores atingidas corpo a corpo), Falcão (Deslocamento de Voo igual ao normal, sem armadura) ou Leão (inimigos a até 1,5m têm Desvantagem para atacar quem não seja você ou outro Bárbaro com essa opção ativa).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Segunda escolha de bônus (Carneiro/Falcão/Leão) ao ativar a Fúria.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Trilha do Coração Selvagem",
		},
		{
			Name: "Campeão dos Deuses", Edition: "5e", ClassID: &id,
			Description: "Reserva de 4 d12 (5 no nível 6, 6 no nível 12, 7 no nível 17) que pode gastar como Ação Bônus, jogando os dados para recuperar PV igual ao total. Restaura em Descanso Longo.",
			Keywords: "Primitivo, Cura", ActionType: "Ação Bônus", Range: "Pessoal",
			LevelScaling: "Nível 6: 5 dados. Nível 12: 6 dados. Nível 17: 7 dados.",
			Effect: "Cura pessoal com uma reserva de dados que cresce por nível.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Trilha do Fanático",
		},
		{
			Name: "Fúria Divina", Edition: "5e", ClassID: &id,
			Description: "Enquanto em Fúria, o primeiro alvo atingido a cada turno com arma ou Ataque Desarmado sofre 1d6 + metade do seu nível de Bárbaro (arred. p/ baixo) de dano adicional Necrótico ou Radiante, à sua escolha a cada vez.",
			Keywords: "Primitivo, Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano extra Necrótico/Radiante no primeiro acerto do turno em Fúria.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trilha do Fanático",
		},
		{
			Name: "Concentração Fanática", Edition: "5e", ClassID: &id,
			Description: "Uma vez por Fúria ativa, ao falhar numa salvaguarda pode rejogá-la com bônus igual ao seu bônus de Dano da Fúria, usando o novo resultado.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Rejogada de salvaguarda com bônus, 1x por Fúria.",
			PowerType: domain.PowerEncounter, Level: 6, IsClassFeature: true, ChoiceGroup: "Trilha do Fanático",
		},
		{
			Name: "Presença Zelosa", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, até 10 criaturas à escolha a até 18m recebem Vantagem em ataques e salvaguardas até o início do seu próximo turno. 1 uso por Descanso Longo, ou gaste um uso de Fúria para restaurar.",
			Keywords: "Primitivo, Divino", ActionType: "Ação Bônus", Range: "18 metros",
			Effect: "Vantagem em ataques/salvaguardas para até 10 aliados.",
			PowerType: domain.PowerDaily, Level: 10, IsClassFeature: true, ChoiceGroup: "Trilha do Fanático",
		},
		{
			Name: "Fúria dos Deuses", Edition: "5e", ClassID: &id,
			Description: "Ao ativar a Fúria, pode assumir forma de combatente divino por 1 minuto ou até cair a 0 PV (1 uso por Descanso Longo): Resistência a dano Necrótico/Psíquico/Radiante, e criaturas reduzidas a 0 PV perto de você podem ser revividas com 1 PV em vez de morrer (uma vez).",
			Keywords: "Primitivo, Divino", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Forma divina temporária: resistências extras e uma revivificação por uso.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true, ChoiceGroup: "Trilha do Fanático",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Bárbaro 5e: características seedadas")
}

// ── Bardo ─────────────────────────────────────────────────────────────────────

func seedBardo5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Bardo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Inspiração de Bardo",
			Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, você inspira uma criatura a até 18 metros que possa vê-lo ou ouvi-lo. Ela recebe um dado de Inspiração de Bardo (d6). Uma vez na próxima 1 hora, após falhar em um Teste de D20, pode jogar esse dado e somar ao resultado. Usos = modificador de Carisma (mínimo 1); recupera em Descanso Longo.",
			Keywords: "Arcano", ActionType: "Ação Bônus", Range: "18 metros",
			Effect:       "Concede dado de Inspiração d6 a um aliado.",
			LevelScaling: "Nível 5: d8 e recupera em Descanso Curto. Nível 10: d10. Nível 15: d12.",
			PowerType:    domain.PowerEncounter, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Conjuração de Bardo",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Carisma. Nível 1: 2 truques, 4 magias preparadas, 2 espaços de 1° círculo. Magias preparadas aumentam com o nível conforme a tabela de Características de Bardo.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a truques e espaços de magia de Bardo.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todos os colégios) ──
		{
			Name:    "Especialista (Bardo)", Edition: "5e", ClassID: &id,
			Description: "Ganha Especialização (dobra o Bônus de Proficiência) em 2 perícias em que já é proficiente.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Dobra o bônus de proficiência em 2 perícias.",
			LevelScaling: "Nível 9: Especialização em mais 2 perícias em que já é proficiente.",
			PowerType:    domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Pau pra Toda Obra", Edition: "5e", ClassID: &id,
			Description: "Pode somar metade do seu Bônus de Proficiência (arredondado para baixo) a qualquer teste de atributo que use uma perícia em que não é proficiente e que não use o Bônus de Proficiência.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+metade do Bônus de Proficiência em testes de perícias não-proficientes.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Fonte de Inspiração", Edition: "5e", ClassID: &id,
			Description: "Passa a restaurar todos os usos de Inspiração de Bardo em Descanso Curto ou Longo. Também pode gastar um espaço de magia (sem ação) para recuperar um uso de Inspiração de Bardo.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Inspiração de Bardo recarrega em Descanso Curto; conversível a partir de espaços de magia.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Contra-Encantamento", Edition: "5e", ClassID: &id,
			Description: "Se você ou uma criatura a até 9m falhar numa salvaguarda contra um efeito que aplica Amedrontado ou Enfeitiçado, pode usar uma Reação para rejogar a salvaguarda com Vantagem.",
			Keywords: "Arcano", ActionType: "Reação", Range: "9 metros",
			Effect:    "Rejogada com Vantagem contra Amedrontado/Enfeitiçado, para si ou aliados.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Segredos Mágicos", Edition: "5e", ClassID: &id,
			Description: "Sempre que o número de Magias Preparadas aumentar (a partir deste nível), pode escolher a nova magia das listas de Bardo, Clérigo, Druida ou Mago — contando como magia de Bardo para você.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Magias preparadas podem vir de 4 listas de classe, não só a de Bardo.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true,
		},
		{
			Name:    "Inspiração Superior", Edition: "5e", ClassID: &id,
			Description: "Ao rolar Iniciativa, recupera usos gastos de Inspiração de Bardo até ter pelo menos 2, se tiver menos que isso.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recarrega Inspiração de Bardo (até 2 usos) ao rolar Iniciativa.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Palavras de Criação", Edition: "5e", ClassID: &id,
			Description: "Sempre tem preparadas Palavra de Poder: Matar e Palavra de Poder: Salvar. Ao conjurá-las, pode escolher um segundo alvo a até 3m do primeiro.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 magias de 9º círculo sempre preparadas, com alvo extra.",
			PowerType: domain.PowerUnlimited, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		{
			Name: "Colégio da Bravura", Edition: "5e", ClassID: &id,
			Description: "Os Bardos do Colégio da Bravura são narradores ousados cujas histórias preservam a memória dos grandes heróis do passado, cantando seus feitos em salões suntuosos ou junto a fogueiras.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "colegio_bardo",
		},
		{
			Name: "Colégio da Dança", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Colégios de Bardo, escolhido no nível 3.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "colegio_bardo",
		},
		{
			Name: "Colégio do Conhecimento", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Colégios de Bardo, escolhido no nível 3.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "colegio_bardo",
		},
		{
			Name: "Colégio do Glamour", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Colégios de Bardo, escolhido no nível 3.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "colegio_bardo",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/14) ─────────────────────────
		{
			Name: "Inspiração em Combate", Edition: "5e", ClassID: &id,
			Description: "Quem tiver um dado de Inspiração de Bardo pode usá-lo de forma Defensiva (Reação ao ser atingido: soma o dado à própria CA contra aquele ataque) ou Ofensiva (soma o dado ao dano após acertar um ataque).",
			Keywords: "Arcano, Marcial", ActionType: "Reação", Range: "Pessoal",
			Effect: "Uso em combate do dado de Inspiração de Bardo: defesa ou dano extra.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio da Bravura",
		},
		{
			Name: "Treinamento Marcial", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência com armas Marciais, armaduras Médias e Escudos. Pode usar uma arma Simples ou Marcial como Foco de Conjuração para magias de Bardo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Proficiências marciais completas + arma como foco de conjuração.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio da Bravura",
		},
		{
			Name: "Ataque Extra (Bardo)", Edition: "5e", ClassID: &id,
			Description: "Pode atacar duas vezes ao executar a ação Atacar. Pode substituir um desses ataques por um truque de tempo de conjuração de 1 ação.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "2 ataques por turno; um pode virar um truque.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Colégio da Bravura",
		},
		{
			Name: "Magia de Batalha", Edition: "5e", ClassID: &id,
			Description: "Após conjurar uma magia de tempo de conjuração de 1 ação, pode realizar um ataque com arma como Ação Bônus.",
			Keywords: "Marcial, Arcano", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Ataque bônus após conjurar uma magia de ação.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Colégio da Bravura",
		},
		{
			Name: "Ginga Fascinante", Edition: "5e", ClassID: &id,
			Description: "Sem armadura/escudo: Vantagem em Atuação (dança); usa Destreza para Ataques Desarmados, causando dano igual ao dado de Inspiração de Bardo + Destreza (sem gastar o dado); CA = 10 + DES + CAR; gastar Inspiração de Bardo permite um Ataque Desarmado extra na mesma ação.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Combo de combate desarmado baseado em Destreza/Carisma sem armadura.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio da Dança",
		},
		{
			Name: "Gingado Coordenado", Edition: "5e", ClassID: &id,
			Description: "Ao rolar Iniciativa, pode gastar um uso de Inspiração de Bardo: você e aliados a até 9m que possam vê-lo/ouvi-lo ganham bônus na Iniciativa igual ao dado jogado.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "9 metros",
			Effect: "Bônus de Iniciativa em grupo ao custo de um uso de Inspiração.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Colégio da Dança",
		},
		{
			Name: "Movimento Inspirador", Edition: "5e", ClassID: &id,
			Description: "Quando um inimigo à vista termina o turno a até 1,5m de você, pode gastar um uso de Inspiração de Bardo (Reação) para se mover metade do Deslocamento; um aliado a até 9m também pode se mover metade do dele. Sem provocar Ataques de Oportunidade.",
			Keywords: "Arcano", ActionType: "Reação", Range: "9 metros",
			Effect: "Reposicionamento próprio + de um aliado, sem provocar Oportunidade.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Colégio da Dança",
		},
		{
			Name: "Evasão Liderada", Edition: "5e", ClassID: &id,
			Description: "Em salvaguardas de Destreza para metade do dano: sem dano se passar, metade se falhar. Criaturas a até 1,5m que também salvarem podem compartilhar esse benefício. Não funciona se estiver Incapacitado.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "1,5 metro",
			Effect: "Evasão (estilo Ladino) compartilhável com aliados próximos.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Colégio da Dança",
		},
		{
			Name: "Palavras de Interrupção", Edition: "5e", ClassID: &id,
			Description: "Quando uma criatura à vista a até 18m tiver sucesso num ataque, teste ou dano, pode gastar Inspiração de Bardo (Reação) e subtrair o dado do resultado dela, podendo virar sucesso em fracasso.",
			Keywords: "Arcano", ActionType: "Reação", Range: "18 metros",
			Effect: "Reduz retroativamente o sucesso/dano de um inimigo.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio do Conhecimento",
		},
		{
			Name: "Proficiências Bônus (Conhecimento)", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência em 3 perícias à escolha.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+3 perícias.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio do Conhecimento",
		},
		{
			Name: "Descobertas Mágicas", Edition: "5e", ClassID: &id,
			Description: "Aprende 2 magias (truque ou de círculo disponível) das listas de Clérigo, Druida ou Mago, sempre preparadas; pode trocá-las a cada novo nível de Bardo.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect: "2 magias extras sempre preparadas de outras listas de classe.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Colégio do Conhecimento",
		},
		{
			Name: "Perícia Inigualável", Edition: "5e", ClassID: &id,
			Description: "Ao falhar num teste de atributo ou ataque, pode gastar um uso de Inspiração de Bardo e somar o dado ao d20 — se ainda assim falhar, o uso não é gasto.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Converte falhas em sucesso sem custo garantido de uso.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Colégio do Conhecimento",
		},
		{
			Name: "Magia Fascinante", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Enfeitiçar Pessoa e Reflexos preparadas. Ao conjurar uma magia de Encantamento ou Ilusão com espaço de magia, pode forçar uma criatura à vista a até 18m a salvar (Sabedoria) ou ficar Amedrontada/Enfeitiçada (à escolha) por 1 minuto. 1 uso por Descanso Longo, ou gaste Inspiração de Bardo para restaurar.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "18 metros",
			Effect: "2 magias fixas + efeito extra de medo/encantamento após conjurar Encantamento/Ilusão.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio do Glamour",
		},
		{
			Name: "Manto de Inspiração", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, gasta Inspiração de Bardo: escolhe um número de criaturas a até 18m igual ao seu mod. de Carisma (mín. 1) para receberem PV Temporários (2x o dado rolado) e poderem se mover até o Deslocamento máximo sem provocar Oportunidade.",
			Keywords: "Arcano", ActionType: "Ação Bônus", Range: "18 metros",
			Effect: "PV Temporários + reposicionamento livre para várias criaturas.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Colégio do Glamour",
		},
		{
			Name: "Manto de Majestade", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Comando preparada e pode conjurá-la sem espaço de magia como Ação Bônus, assumindo aparência sobrenatural por 1 minuto (durante a qual pode repetir Comando de graça); alvos Enfeitiçados por você falham automaticamente a salvaguarda contra esse Comando. 1 uso por Descanso Longo, ou gaste espaço de 3º círculo+ para restaurar.",
			Keywords: "Arcano, Magia", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Comando gratuito e repetível por 1 minuto, com falha automática contra Enfeitiçados.",
			PowerType: domain.PowerDaily, Level: 6, IsClassFeature: true, ChoiceGroup: "Colégio do Glamour",
		},
		{
			Name: "Majestade Inquebrável", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, assume presença majestosa por 1 minuto ou até ficar Incapacitado: todo atacante que o acerta pela primeira vez no turno deve salvar (Carisma) ou o ataque falha. 1 uso por Descanso Curto ou Longo.",
			Keywords: "Arcano", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Chance de anular o primeiro acerto de cada atacante por turno, por 1 minuto.",
			PowerType: domain.PowerEncounter, Level: 14, IsClassFeature: true, ChoiceGroup: "Colégio do Glamour",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Bardo 5e: características seedadas")
}

// ── Bruxo ─────────────────────────────────────────────────────────────────────

func seedBruxo5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Bruxo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Magia de Pacto",
			Edition: "5e", ClassID: &id,
			Description: "Você obtém espaços de magia por meio de um pacto com uma entidade mística. Nível 1: 1 espaço de 1° círculo, 2 truques, 2 magias preparadas. Diferente de outras classes, seus espaços de magia recuperam em Descanso Curto ou Longo. Atributo de conjuração: Carisma.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Espaços de magia que recuperam em Descanso Curto. Atributo: Carisma.",
			LevelScaling: "Nível 3: 2 espaços de 2° círculo. Nível 5: 2 espaços de 3° círculo.",
			PowerType:    domain.PowerEncounter, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Invocações Místicas", Edition: "5e", ClassID: &id,
			Description: "Recebe uma Invocação Mística à escolha (fragmentos de conhecimento proibido que concedem uma habilidade mágica permanente, ex.: Pacto do Tomo). Ganha invocações extras em níveis superiores; ao subir de nível pode trocar uma invocação por outra elegível.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "1 Invocação Mística à escolha.",
			LevelScaling: "Total de invocações conhecidas: 1 (nível 1), 2 (nível 2), 3 (nível 5), 4 (nível 7), 5 (nível 9), 6 (nível 12), 7 (nível 15), 8 (nível 18).",
			PowerType:    domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name: "Patrono Arquifada", Edition: "5e", ClassID: &id,
			Description: "Faça Acordos com Feéricos Excêntricos — seu pacto é fundamentado no poder de Faéria, com uma Arquifada como o Príncipe do Gelo ou Titânia, ou um espectro Feérico.",
			Keywords: "Arcano, Feérico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "patrono_bruxo",
		},
		{
			Name: "Patrono Celestial", Edition: "5e", ClassID: &id,
			Description: "Invoque o Poder dos Céus — seu pacto é fundamentado nos Planos Superiores, com um empiriano, couatl, esfinge, unicórnio ou outra entidade celestial.",
			Keywords: "Arcano, Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "patrono_bruxo",
		},
		{
			Name: "Patrono O Grande Antigo", Edition: "5e", ClassID: &id,
			Description: "Descubra o Conhecimento Proibido de Seres Inefáveis — conecta-se a uma entidade indescritível do Reino Distante ou a um deus ancestral como Tharizdun ou Grande Cthulhu.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "patrono_bruxo",
		},
		{
			Name: "Patrono Ínfero", Edition: "5e", ClassID: &id,
			Description: "Realize um Pacto com os Planos Inferiores — negocia com um lorde demônio, um arquidiabo, ou outra entidade maligna dos Planos Inferiores.",
			Keywords: "Arcano, Ínfero", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "patrono_bruxo",
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todos os patronos) ──
		{
			Name:    "Astúcia Mágica", Edition: "5e", ClassID: &id,
			Description: "Ao final de um rito de 1 minuto, recupera espaços de Magia de Pacto gastos em número igual à metade do seu máximo (arred. para cima). 1 uso por Descanso Longo.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recupera metade dos espaços de magia gastos, uma vez por Descanso Longo.",
			PowerType: domain.PowerDaily, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Contatar Patrono", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Contato Extraplanar preparada e pode conjurá-la sem espaço de magia para falar com seu patrono, tendo sucesso automático na salvaguarda da magia. 1 uso por Descanso Longo.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Contato Extraplanar gratuito e garantido com o patrono.",
			PowerType: domain.PowerDaily, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Arcana Mística", Edition: "5e", ClassID: &id,
			Description: "Escolhe uma magia de Bruxo de 6º círculo como arcanum, conjurável 1x sem espaço de magia (recarrega em Descanso Longo). Ganha mais um arcanum nos níveis 13 (7º círculo), 15 (8º círculo) e 17 (9º círculo).",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Magias de alto círculo conjuráveis de graça, 1x por Descanso Longo cada.",
			LevelScaling: "Nível 13: +magia de 7º círculo. Nível 15: +8º círculo. Nível 17: +9º círculo.",
			PowerType:    domain.PowerDaily, Level: 11, IsClassFeature: true,
		},
		{
			Name:    "Mestre Místico", Edition: "5e", ClassID: &id,
			Description: "Ao usar Astúcia Mágica, agora recupera TODOS os espaços de Magia de Pacto gastos, não apenas metade.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Astúcia Mágica passa a recuperar 100% dos espaços de magia.",
			PowerType: domain.PowerUnlimited, Level: 20, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/10/14) ──────────────────────
		{
			Name: "Passos Feéricos", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar Passo Nebuloso sem espaço de magia (usos = mod. Carisma, mín. 1, recarrega em Descanso Longo), com efeito extra à escolha: Passo Provocante (impõe Desvantagem a inimigos perto do espaço que deixou) ou Passo Revigorante (dá PV Temporários após teleportar).",
			Keywords: "Arcano, Feérico", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 6: pode conjurar como Reação ao sofrer dano; ganha Passo Desvanecedor (Invisibilidade) e Passo Terrível (dano Psíquico em área) como novas opções.",
			Effect: "Passo Nebuloso gratuito com efeitos extras de controle/cura.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Patrono Arquifada",
		},
		{
			Name: "Defesas Sedutoras", Edition: "5e", ClassID: &id,
			Description: "Imune a Enfeitiçado. Ao ser atingido por um ataque, pode usar Reação para reduzir o dano pela metade e forçar salvaguarda de Sabedoria no atacante — se falhar, sofre o mesmo dano Psíquico. 1 uso por Descanso Longo, ou gaste um espaço de Magia de Pacto para restaurar.",
			Keywords: "Arcano, Feérico", ActionType: "Reação", Range: "Pessoal",
			Effect: "Imunidade a Enfeitiçado + retaliação psíquica reduzindo dano recebido.",
			PowerType: domain.PowerDaily, Level: 10, IsClassFeature: true, ChoiceGroup: "Patrono Arquifada",
		},
		{
			Name: "Magia Sedutora", Edition: "5e", ClassID: &id,
			Description: "Após conjurar uma magia de Encantamento ou Ilusão com ação e espaço de magia, pode conjurar Passo Nebuloso na mesma ação sem gastar espaço.",
			Keywords: "Arcano, Feérico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Teleporte gratuito combinado com Encantamento/Ilusão.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Patrono Arquifada",
		},
		{
			Name: "Luz Medicinal", Edition: "5e", ClassID: &id,
			Description: "Reserva de d6s (1 + nível de Bruxo). Como Ação Bônus, gasta até seu mod. de Carisma (mín. 1) desses dados para curar a si ou uma criatura à vista a até 18m. Restaura em Descanso Longo.",
			Keywords: "Divino, Cura", ActionType: "Ação Bônus", Range: "18 metros",
			Effect: "Reserva de cura em d6, escalando com o nível.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Patrono Celestial",
		},
		{
			Name: "Alma Radiante", Edition: "5e", ClassID: &id,
			Description: "Resistência a Dano Radiante. Uma vez por turno, ao conjurar uma magia que cause dano Ígneo ou Radiante, pode somar seu mod. de Carisma ao dano contra um alvo.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência a Radiante + dano bônus de Carisma em magias de fogo/luz.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Patrono Celestial",
		},
		{
			Name: "Resiliência Celestial", Edition: "5e", ClassID: &id,
			Description: "Ao usar Astúcia Mágica ou completar um Descanso, recebe PV Temporários (nível de Bruxo + mod. Carisma) e pode dar metade disso a até 5 criaturas à vista.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "PV Temporários para si e aliados em descansos/Astúcia Mágica.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Patrono Celestial",
		},
		{
			Name: "Vingança Calcinante", Edition: "5e", ClassID: &id,
			Description: "Quando você ou aliado a até 18m for fazer Salvaguarda Contra a Morte, pode gastar essa característica: a criatura recupera metade do PV máximo e se levanta; inimigos a até 9m dela sofrem 2d8+CAR de dano Radiante e ficam Cegos até o fim do turno. Recarrega em Descanso Longo.",
			Keywords: "Divino", ActionType: "Reação", Range: "18 metros",
			Effect: "Cura de emergência em área com dano/cegueira em inimigos próximos.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true, ChoiceGroup: "Patrono Celestial",
		},
		{
			Name: "Mente Desperta", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, estabelece conexão telepática com uma criatura à vista a até 9m, comunicando-se a até 1,5km × mod. de Carisma, por um número de minutos igual ao seu nível de Bruxo.",
			Keywords: "Arcano, Aberrante", ActionType: "Ação Bônus", Range: "9 metros",
			Effect: "Telepatia de longo alcance com uma criatura.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Patrono O Grande Antigo",
		},
		{
			Name: "Combatente Clarividente", Edition: "5e", ClassID: &id,
			Description: "Ao formar o vínculo de Mente Desperta, pode forçar salvaguarda de Sabedoria: se falhar, a criatura tem Desvantagem para atacá-lo e você tem Vantagem contra ela pela duração. Recarrega em Descanso Curto/Longo, ou gastando um espaço de Magia de Pacto.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Vantagem/Desvantagem de combate contra o alvo do vínculo telepático.",
			PowerType: domain.PowerEncounter, Level: 6, IsClassFeature: true, ChoiceGroup: "Patrono O Grande Antigo",
		},
		{
			Name: "Danação Mística", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Danação preparada; ao conjurá-la e escolher um atributo, o alvo também tem Desvantagem nas salvaguardas desse atributo pela duração.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Danação sempre disponível, com Desvantagem extra de salvaguarda.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Patrono O Grande Antigo",
		},
		{
			Name: "Escudo Mental", Edition: "5e", ClassID: &id,
			Description: "Seus pensamentos são protegidos contra telepatia/leitura de mente. Resistência a Dano Psíquico; quem causar Dano Psíquico em você sofre o mesmo dano de volta.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Proteção mental + retaliação em espelho de dano Psíquico.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Patrono O Grande Antigo",
		},
		{
			Name: "Criar Servo", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar Invocar Aberração, pode dispensar a Concentração (duração vira 1 minuto), dando à Aberração PV Temporários (nível de Bruxo + CAR); ela causa dano Psíquico extra contra alvos sob sua Danação.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Invocar Aberração sem Concentração + dano extra combinado com Danação.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Patrono O Grande Antigo",
		},
		{
			Name: "Magias Psíquicas", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de Bruxo que causa dano, pode mudar o tipo para Psíquico. Pode conjurar magias de Bruxo de Encantamento ou Ilusão sem componentes Verbais ou Somáticos.",
			Keywords: "Arcano, Ínfero", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Flexibilidade de tipo de dano e componentes em magias.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Patrono Ínfero",
		},
		{
			Name: "Bênção do Tenebroso", Edition: "5e", ClassID: &id,
			Description: "Ao reduzir um inimigo a 0 PV (ou outra pessoa fazer isso a até 3m de você), recebe PV Temporários iguais ao seu mod. de Carisma + nível de Bruxo (mín. 1).",
			Keywords: "Arcano, Ínfero", ActionType: "Passiva", Range: "3 metros",
			Effect: "PV Temporários ao abater inimigos perto de você.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Patrono Ínfero",
		},
		{
			Name: "A Sorte do Próprio Tenebroso", Edition: "5e", ClassID: &id,
			Description: "Após ver o resultado de um teste de atributo ou salvaguarda, pode somar 1d10 ao resultado. Usos = mod. de Carisma (mín. 1), 1x por jogada, recarrega em Descanso Longo.",
			Keywords: "Arcano, Ínfero", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+1d10 retroativo em um teste/salvaguarda.",
			PowerType: domain.PowerDaily, Level: 6, IsClassFeature: true, ChoiceGroup: "Patrono Ínfero",
		},
		{
			Name: "Resistência Ínfera", Edition: "5e", ClassID: &id,
			Description: "Ao completar um Descanso Curto ou Longo, escolhe um tipo de dano (exceto Energético) e ganha Resistência a ele até escolher outro.",
			Keywords: "Arcano, Ínfero", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência a um tipo de dano, trocável por descanso.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Patrono Ínfero",
		},
		{
			Name: "Lançar no Inferno", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao acertar um ataque, pode forçar salvaguarda de Carisma no alvo ou ele desaparece para os Planos Inferiores, sofrendo 8d10 de dano Psíquico (se não for Ínfero) e ficando Incapacitado até seu próximo turno. Recarrega em Descanso Longo, ou gaste um espaço de Magia de Pacto.",
			Keywords: "Arcano, Ínfero", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Banimento temporário com dano pesado, 1x por turno.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true, ChoiceGroup: "Patrono Ínfero",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Bruxo 5e: características seedadas")
}

// ── Clérigo ───────────────────────────────────────────────────────────────────

func seedClerigo5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Clérigo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Conjuração de Clérigo",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Sabedoria. Nível 1: 3 truques, 2 espaços de 1° círculo. Magias preparadas = nível de Clérigo + modificador de Sabedoria (sempre preparadas as magias do Domínio).",
			Keywords: "Divino, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a truques e espaços de magia de Clérigo.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Protetor", Edition: "5e", ClassID: &id,
			Description: "Treinado para a batalha: ganha proficiência com armas Marciais e treinamento com Armadura Pesada.",
			Keywords: "Divino, Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Proficiência com armas Marciais e Armadura Pesada.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "ordem_divina_clerigo",
		},
		{
			Name:    "Taumaturgo", Edition: "5e", ClassID: &id,
			Description: "Conhece um truque adicional de Clérigo, e ganha bônus (mod. de Sabedoria, mín. +1) em testes de Inteligência (Arcanismo ou Religião).",
			Keywords: "Divino, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "+1 truque e bônus em Arcanismo/Religião.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "ordem_divina_clerigo",
		},
		{
			Name:    "Canalizar Divindade (Clérigo)", Edition: "5e", ClassID: &id,
			Description: "2 usos (recupera 1 em Descanso Curto, todos em Longo) para ativar Centelha Divina (cura ou dano Necrótico/Radiante, 1d8 + Sabedoria, com salvaguarda de Constituição) ou Expulsar Mortos-Vivos (Mortos-Vivos a até 9m salvam Sabedoria ou ficam Amedrontados/Incapacitados por 1 minuto, fugindo de você).",
			Keywords: "Divino", ActionType: "Passiva", Range: "9 metros",
			Effect:       "2 usos de Canalizar Divindade: cura/dano radiante-necrótico, ou repelir mortos-vivos.",
			LevelScaling: "Centelha Divina escala para 2d8 (nível 7), 3d8 (nível 13), 4d8 (nível 18). Mais usos de Canalizar Divindade em níveis altos.",
			PowerType:    domain.PowerEncounter, Level: 2, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todos os domínios) ──
		{
			Name:    "Fulminar Mortos-Vivos", Edition: "5e", ClassID: &id,
			Description: "Ao usar Expulsar Mortos-Vivos, joga d8s (= mod. de Sabedoria, mín. 1) e cada Morto-Vivo que falhar na salvaguarda sofre esse total como dano Radiante, sem encerrar o efeito de expulsão.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Dano Radiante extra ao usar Expulsar Mortos-Vivos.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Golpes Abençoados", Edition: "5e", ClassID: &id,
			Description: "Escolhe: Conjuração Poderosa (soma mod. de Sabedoria ao dano de truques de Clérigo) ou Golpe Divino (1x por turno, +1d8 de dano Necrótico ou Radiante ao acertar com arma).",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Dano extra em truques ou em ataques com arma, à escolha.",
			LevelScaling: "Nível 14 (Golpes Abençoados Aprimorados): Conjuração Poderosa também dá PV Temporários (2x Sabedoria); Golpe Divino sobe para 2d8.",
			PowerType:    domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Intervenção Divina", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, conjura qualquer magia de Clérigo de 5º círculo ou menor (sem Reação) sem gastar espaço de magia nem componentes Materiais. 1 uso por Descanso Longo.",
			Keywords: "Divino, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Conjura uma magia de até 5º círculo de graça, 1x por Descanso Longo.",
			LevelScaling: "Nível 20 (Intervenção Divina Maior): pode escolher Desejo; se fizer isso, só pode reusar após 2d4 Descansos Longos.",
			PowerType:    domain.PowerDaily, Level: 10, IsClassFeature: true,
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}

	dominios := []domain.Skill{
		{
			Name: "Domínio da Guerra", Edition: "5e", ClassID: &id,
			Description: "Inspire Bravura e Derrote Inimigos — Clérigos do Domínio da Guerra se destacam em batalhas, inspirando outros a lutar pelo bem ou convertendo violência em oração.",
			Keywords: "Divino, Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "dominio_clerigo",
		},
		{
			Name: "Domínio da Luz", Edition: "5e", ClassID: &id,
			Description: "Traga a Luz para Banir a Escuridão — Clérigos do Domínio da Luz possuem a visão clara de suas divindades, encarregados de afastar mentiras e dissipar as trevas.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "dominio_clerigo",
		},
		{
			Name: "Domínio da Trapaça", Edition: "5e", ClassID: &id,
			Description: "Pregue Peças e Desafie as Autoridades — Clérigos do Domínio da Trapaça usam magias de enganação, ilusão e furtividade, preferindo estratagemas ao confronto direto.",
			Keywords: "Divino, Ilusão", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "dominio_clerigo",
		},
		{
			Name: "Domínio da Vida", Edition: "5e", ClassID: &id,
			Description: "Alivie as Feridas do Mundo — Clérigos do Domínio da Vida são mestres da cura, canalizando energia positiva para restaurar os feridos.",
			Keywords: "Divino, Cura", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "dominio_clerigo",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/17) ─────────────────────────
		{
			Name: "Ataque Direcionado", Edition: "5e", ClassID: &id,
			Description: "Quando você ou uma criatura a até 9m erra um ataque, pode gastar um uso de Canalizar Divindade (Reação) para dar +10 a essa jogada de ataque.",
			Keywords: "Divino, Marcial", ActionType: "Reação", Range: "9 metros",
			Effect: "+10 retroativo numa jogada de ataque, gastando Canalizar Divindade.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Guerra",
		},
		{
			Name: "Sacerdote da Guerra", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, realiza um ataque com arma ou Desarmado. Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Curto ou Longo.",
			Keywords: "Divino, Marcial", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Ataque extra como Ação Bônus, várias vezes por descanso.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Guerra",
		},
		{
			Name: "Bênção do Deus da Guerra", Edition: "5e", ClassID: &id,
			Description: "Pode gastar um uso de Canalizar Divindade para conjurar Arma Espiritual ou Escudo da Fé sem espaço de magia e sem precisar de Concentração (dura 1 minuto, encerra se reconjurar, ficar Incapacitado ou morrer).",
			Keywords: "Divino, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect: "2 magias de proteção/ataque gratuitas via Canalizar Divindade, sem Concentração.",
			PowerType: domain.PowerEncounter, Level: 6, IsClassFeature: true, ChoiceGroup: "Domínio da Guerra",
		},
		{
			Name: "Avatar da Guerra", Edition: "5e", ClassID: &id,
			Description: "Ganha Resistência a dano Contundente, Cortante e Perfurante.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência aos 3 tipos de dano físico.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Domínio da Guerra",
		},
		{
			Name: "Brilho do Amanhecer", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta Canalizar Divindade para emitir luz numa Emanação de 9m, dissipando Escuridão mágica; criaturas à escolha na área salvam Constituição ou sofrem 2d10 + nível de Clérigo de dano Radiante (metade se passar).",
			Keywords: "Divino", ActionType: "Passiva", Range: "Emanação 9 metros",
			Effect: "Explosão de luz radiante em área, via Canalizar Divindade.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Luz",
		},
		{
			Name: "Labareda Protetora", Edition: "5e", ClassID: &id,
			Description: "Quando uma criatura à vista a até 9m faz um ataque, pode usar Reação para impor Desvantagem nele. Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Divino", ActionType: "Reação", Range: "9 metros",
			LevelScaling: "Nível 6 (Labareda Protetora Aprimorada): recarrega em Descanso Curto ou Longo, e concede 2d6 + Sabedoria de PV Temporários ao alvo defendido.",
			Effect: "Impõe Desvantagem em um ataque contra alguém próximo.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Luz",
		},
		{
			Name: "Coroa de Luz", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, emite aura de Luz Plena (18m) + Meia-luz (mais 9m) por 1 minuto; inimigos na Luz Plena têm Desvantagem em salvaguardas contra Brilho do Amanhecer e magias de dano Ígneo/Radiante. Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Divino", ActionType: "Passiva", Range: "18 metros",
			Effect: "Aura de luz que enfraquece salvaguardas inimigas contra fogo/radiante.",
			PowerType: domain.PowerDaily, Level: 17, IsClassFeature: true, ChoiceGroup: "Domínio da Luz",
		},
		{
			Name: "Bênção do Trapaceiro", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, dá a si ou a um voluntário a até 9m Vantagem em testes de Destreza (Furtividade) até Descanso Longo ou até reusar a característica.",
			Keywords: "Divino, Ilusão", ActionType: "Passiva", Range: "9 metros",
			Effect: "Vantagem duradoura em Furtividade para si ou um aliado.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Trapaça",
		},
		{
			Name: "Invocar Duplicidade", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, gasta Canalizar Divindade para criar uma ilusão idêntica a si a até 9m por 1 minuto: pode conjurar magias como se estivesse nela, tem Vantagem em ataques contra quem está perto dela, e pode movê-la até 9m (até 36m de você) como Ação Bônus.",
			Keywords: "Divino, Ilusão", ActionType: "Ação Bônus", Range: "9-36 metros",
			LevelScaling: "Nível 6: pode trocar de lugar (teleportar) com a ilusão ao criá-la/movê-la. Nível 17: Vantagem compartilhada com aliados perto da ilusão, e cura ao dissipá-la.",
			Effect: "Duplo ilusório para conjurar/atacar à distância e enganar inimigos.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Trapaça",
		},
		{
			Name: "Discípulo da Vida", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia com espaço de magia que restaura PV, a criatura curada recebe PV adicionais iguais a 2 + o nível do espaço usado.",
			Keywords: "Divino, Cura", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+2 e mais PV de cura por círculo do espaço de magia usado.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Vida",
		},
		{
			Name: "Preservar a Vida", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta Canalizar Divindade para dividir PV iguais a 5x seu nível de Clérigo entre criaturas Sangrando a até 9m (incluindo você), sem passar de metade do PV máximo de cada uma.",
			Keywords: "Divino, Cura", ActionType: "Passiva", Range: "9 metros",
			Effect: "Cura em massa distribuída, via Canalizar Divindade.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Domínio da Vida",
		},
		{
			Name: "Curandeiro Abençoado", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia com espaço de magia que cura outra criatura, você também recupera PV iguais a 2 + o círculo do espaço usado.",
			Keywords: "Divino, Cura", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Autocura sempre que cura outra pessoa com espaço de magia.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Domínio da Vida",
		},
		{
			Name: "Cura Suprema", Edition: "5e", ClassID: &id,
			Description: "Ao restaurar PV com magia ou Canalizar Divindade, usa o valor máximo de cada dado em vez de rolar.",
			Keywords: "Divino, Cura", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Toda cura de dados vira o valor máximo possível.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Domínio da Vida",
		},
	}
	for _, s := range dominios {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Clérigo 5e: características seedadas")
}

// ── Druida ────────────────────────────────────────────────────────────────────

func seedDruida5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Druida")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Conjuração de Druida",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Sabedoria. Nível 1: 2 truques, 2 espaços de 1° círculo. Magias preparadas = nível de Druida + modificador de Sabedoria.",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a truques e espaços de magia de Druida.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Druídico",
			Edition: "5e", ClassID: &id,
			Description: "Você conhece Druídico, a língua secreta dos druidas, e sempre tem Falar com Animais preparada. Pode deixar mensagens ocultas que apenas outros druidas conseguem decifrar em superfícies naturais.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Conhecimento da língua secreta Druídico + Falar com Animais sempre preparada.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Protetor (Druida)", Edition: "5e", ClassID: &id,
			Description: "Treinado para a batalha: ganha proficiência com armas Marciais e treinamento com Armadura Média.",
			Keywords: "Primitivo, Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Proficiência com armas Marciais e Armadura Média.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "ordem_primal_druida",
		},
		{
			Name:    "Xamã", Edition: "5e", ClassID: &id,
			Description: "Conhece um truque adicional de Druida, e ganha bônus (mod. de Sabedoria, mín. +1) em testes de Inteligência (Arcanismo ou Natureza).",
			Keywords: "Primitivo, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "+1 truque e bônus em Arcanismo/Natureza.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "ordem_primal_druida",
		},
		{
			Name:    "Companheiro Selvagem", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta um espaço de magia ou um uso de Forma Selvagem para conjurar Convocar Familiar sem componentes Materiais; o familiar é uma criatura Feérica e some ao Descanso Longo.",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Convocar Familiar de graça, usando magia ou Forma Selvagem.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Forma Selvagem", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, multimorfa para uma forma Animal conhecida, por um número de horas igual a metade do seu nível de Druida (ou até reusar, ficar Incapacitado ou morrer).",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:       "Transformação em animal, com duração baseada no nível.",
			LevelScaling: "Nível 5 (Ressurgimento Selvagem): pode recuperar 1 uso gastando um espaço de magia, 1x por turno; e pode converter 1 uso em um espaço de 1º círculo, 1x por Descanso Longo.",
			PowerType:    domain.PowerEncounter, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Fúria Elemental", Edition: "5e", ClassID: &id,
			Description: "Escolhe: Ataque Primal (1x por turno, +1d8 de dano elemental à escolha ao acertar com arma ou ataque de Forma Selvagem) ou Conjuração Poderosa (soma mod. de Sabedoria ao dano de truques de Druida).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Dano extra em ataques ou truques, à escolha.",
			LevelScaling: "Nível 15 (Fúria Elemental Aprimorada): Ataque Primal sobe para 2d8; Conjuração Poderosa aumenta o alcance de truques de 3m+ para 90m.",
			PowerType:    domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Magias Bestiais", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar magias enquanto em Forma Selvagem, exceto as que exigem componente Material com custo ou que seja consumido.",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Conjuração de magias liberada durante a Forma Selvagem.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Arquidruida", Edition: "5e", ClassID: &id,
			Description: "Forma Selvagem Eterna: recupera 1 uso ao rolar Iniciativa se não tiver mais nenhum. Natureza Xamânica: converte usos de Forma Selvagem em espaço de magia (cada uso = 2º círculo), 1x por Descanso Longo. Longevidade: envelhece 1 ano a cada 10.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Forma Selvagem quase ilimitada + conversão em magia + longevidade.",
			PowerType: domain.PowerUnlimited, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024 — moveu do nível 2 pro 3) ──
		{
			Name: "Círculo da Lua", Edition: "5e", ClassID: &id,
			Description: "Assuma Formas Animais para Proteger a Vida Selvagem — Druidas do Círculo da Lua canalizam a magia lunar para se transformarem, espreitando como um grande felino, sobrevoando como uma águia ou atravessando a vegetação como um urso.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "circulo_druida",
		},
		{
			Name: "Círculo da Terra", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Círculos Druídicos, escolhido no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "circulo_druida",
		},
		{
			Name: "Círculo das Estrelas", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Círculos Druídicos, escolhido no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "circulo_druida",
		},
		{
			Name: "Círculo do Mar", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Círculos Druídicos, escolhido no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "circulo_druida",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/10/14) ──────────────────────
		{
			Name: "Formas Animais dos Círculos Druídicos", Edition: "5e", ClassID: &id,
			Description: "Em Forma Selvagem: Nível de Desafio máximo = nível de Druida ÷ 3 (arredondado p/ baixo); CA vira 13 + Sabedoria se maior que a da Fera; ganha PV Temporários = 3x nível de Druida.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 6: ataques na forma podem causar dano Radiante à escolha, e soma Sabedoria em salvaguardas de Constituição enquanto transformado.",
			Effect: "Formas selvagens mais fortes (CD, CA, PV Temp) que o padrão da classe.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Círculo da Lua",
		},
		{
			Name: "Passo Lunar", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, teleporta-se até 9m para um espaço à vista e ganha Vantagem no próximo ataque do turno. Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo (ou gastando espaço de 2º círculo+ por uso).",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "9 metros",
			LevelScaling: "Nível 14 (Forma Lunar): 1x por turno causa 2d10 de dano Radiante extra na Forma Selvagem; Passo Lunar pode levar uma criatura voluntária junto.",
			Effect: "Teleporte + Vantagem em ataque.",
			PowerType: domain.PowerDaily, Level: 10, IsClassFeature: true, ChoiceGroup: "Círculo da Lua",
		},
		{
			Name: "Auxílio da Terra", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta um uso de Forma Selvagem: numa Esfera de 3m a até 18m, inimigos salvam Constituição ou sofrem 2d6 de dano Necrótico (metade se passar), e uma criatura na área recupera 2d6 PV.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "18 metros",
			LevelScaling: "Dano e cura sobem para 3d6 no nível 10 e 4d6 no nível 14.",
			Effect: "Dano em área + cura simultânea, gastando Forma Selvagem.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Círculo da Terra",
		},
		{
			Name: "Recuperação Natural", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar uma magia preparada via Magias do Círculo sem espaço de magia (1x por Descanso Longo). Em Descanso Curto, pode recuperar espaços de magia cuja soma de círculos seja até metade do seu nível (arred. p/ cima), nenhum de 6º círculo ou mais — 1x por Descanso Longo.",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Magia grátis + recuperação parcial de espaços de magia em Descanso Curto.",
			PowerType: domain.PowerDaily, Level: 6, IsClassFeature: true, ChoiceGroup: "Círculo da Terra",
		},
		{
			Name: "Proteção Natural", Edition: "5e", ClassID: &id,
			Description: "Imune a Envenenado; tem Resistência ao tipo de dano do terreno escolhido (Árido=Ígneo, Polar=Gélido, Temperado=Elétrico, Tropical=Venenoso).",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Imunidade a veneno + resistência elemental ligada ao terreno.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Círculo da Terra",
		},
		{
			Name: "Santuário Natural", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta um uso de Forma Selvagem para criar um Cubo de 4,5m de árvores/vinhas a até 36m, por 1 minuto: você e aliados ali têm Cobertura Parcial e a Resistência de Proteção Natural. Pode mover o Cubo 18m como Ação Bônus.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "36 metros",
			Effect: "Zona de cobertura + resistência compartilhada em área móvel.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Círculo da Terra",
		},
		{
			Name: "Forma Estrelada", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, gasta um uso de Forma Selvagem para virar uma forma luminosa (mantém suas estatísticas) por 10 minutos, emitindo luz. Escolhe uma constelação: Arqueiro (ataque à distância 1d8+Sabedoria Radiante), Dragão (trata 9- como 10 em INT/SAB/Concentração) ou Taça (cura extra 1d8+Sabedoria ao curar com magia).",
			Keywords: "Primitivo, Radiante", ActionType: "Ação Bônus", Range: "Pessoal",
			LevelScaling: "Nível 10: Arqueiro e Taça sobem para 2d8; Dragão ganha Voo 6m. Nível 14: Resistência a dano físico na Forma Estrelada.",
			Effect: "Forma alternativa com 3 builds à escolha (dano, utilidade ou cura).",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Círculo das Estrelas",
		},
		{
			Name: "Mapa Estelar", Edition: "5e", ClassID: &id,
			Description: "Cria um Foco de Conjuração especial; enquanto o segura, tem Orientação e Raio Guia sempre preparadas, podendo conjurar Raio Guia sem espaço de magia (usos = mod. Sabedoria, mín. 1, recarrega em Descanso Longo).",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 6 (Presságio Cósmico): 1x por Descanso Longo, ganha uma Reação de sorte (+ ou −1d6 num d20 de uma criatura à vista) até o próximo Descanso Longo.",
			Effect: "Foco de conjuração com magias de suporte sempre disponíveis.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Círculo das Estrelas",
		},
		{
			Name: "Ira do Mar", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, gasta um uso de Forma Selvagem para criar uma Emanação de 1,5m por 10 minutos; como Ação Bônus em turnos seguintes, escolhe um alvo na área para salvar Constituição ou sofrer dano Gélido (Xd6, X = mod. Sabedoria) e ser empurrado 4,5m (se Grande ou menor).",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "1,5 metro",
			LevelScaling: "Nível 6: Emanação cresce para 3m e ganha Deslocamento de Natação. Nível 10 (Filho da Tempestade): ganha Voo e Resistência a Elétrico/Gélido/Trovejante enquanto ativa.",
			Effect: "Aura de dano/controle marítimo, escalando em alcance e resistências.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Círculo do Mar",
		},
		{
			Name: "Manifestação Oceânica", Edition: "5e", ClassID: &id,
			Description: "Pode manifestar a Emanação de Ira do Mar ao redor de uma criatura voluntária a até 18m em vez de si mesmo; gastando 2 usos de Forma Selvagem, pode manifestá-la ao redor de ambos ao mesmo tempo.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "18 metros",
			Effect: "Compartilha a aura de Ira do Mar com um aliado à distância.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Círculo do Mar",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Druida 5e: características seedadas")
}

// ── Feiticeiro ────────────────────────────────────────────────────────────────

func seedFeiticeiro5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Feiticeiro")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Conjuração de Feiticeiro",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Carisma. Nível 1: 4 truques, 2 magias preparadas, 2 espaços de 1° círculo.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a truques e espaços de magia de Feiticeiro.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Feitiçaria Inata", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, libera sua magia latente por 1 minuto: CD de suas magias +1 e Vantagem nas jogadas de ataque de magia de Feiticeiro. 2 usos, recarrega em Descanso Longo.",
			Keywords: "Arcano", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:       "+1 CD e Vantagem em ataques mágicos por 1 minuto, 2x por Descanso Longo.",
			LevelScaling: "Nível 7 (Feitiçaria Encarnada): pode reativar gastando 2 Pontos de Feitiçaria quando sem usos; enquanto ativa, pode usar 2 opções de Metamagia por magia. Nível 20 (Apoteose Arcana): 1 opção de Metamagia por turno de graça enquanto ativa.",
			PowerType:    domain.PowerDaily, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Origem da Feitiçaria",
			Edition: "5e", ClassID: &id,
			Description: "Escolha a fonte do seu poder mágico inato: Aberrante, Dracônica, Mecânica ou Selvagem. A origem concede características exclusivas a partir do nível 1.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Define a fonte do poder do Feiticeiro.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "origem_feiticeiro",
		},
		{
			Name: "Feitiçaria Aberrante", Edition: "5e", ClassID: &id,
			Description: "Exerça o Sobrenatural Poder Psiônico — uma influência alienígena (Plano Astral, Reino Distante ou um girino devorador de mentes) concedeu poder psiônico à sua mente.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "origem_feiticeiro",
		},
		{
			Name: "Feitiçaria Dracônica", Edition: "5e", ClassID: &id,
			Description: "Respire a Magia dos Dragões — sua magia inata provém da dádiva de um dragão ancestral, um local impregnado de poder dracônico, ou um antepassado dragão.",
			Keywords: "Arcano, Dracônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "origem_feiticeiro",
		},
		{
			Name: "Feitiçaria Mecânica", Edition: "5e", ClassID: &id,
			Description: "Canalize as Forças Cósmicas da Ordem — seu poder vem de Mecanos, o plano moldado pela eficiência de um relógio, habitado pelos modrons.",
			Keywords: "Arcano, Ordem", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "origem_feiticeiro",
		},
		{
			Name: "Feitiçaria Selvagem", Edition: "5e", ClassID: &id,
			Description: "Liberte a Magia Caótica — sua magia inata provém das forças do caos, seja por exposição a magia bruta, uma bênção feérica, uma marca demoníaca ou puro acaso.",
			Keywords: "Arcano, Caótico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "origem_feiticeiro",
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todas as origens) ──
		{
			Name:    "Fonte de Magia", Edition: "5e", ClassID: &id,
			Description: "Ganha Pontos de Feitiçaria (2 no nível 2, mais em níveis altos) usados para criar efeitos mágicos e alimentar a Metamagia.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recurso de Pontos de Feitiçaria.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Metamagia", Edition: "5e", ClassID: &id,
			Description: "Ganha 2 Opções de Metamagia à escolha (ex.: Magia Acelerada) para modificar magias conjuradas, gastando Pontos de Feitiçaria. Apenas 1 opção por magia, salvo indicação contrária.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "2 modificadores de magia à escolha, custando Pontos de Feitiçaria.",
			LevelScaling: "Nível 10: +2 opções de Metamagia. Nível 17: +2 opções de Metamagia.",
			PowerType:    domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Restauração Feiticeira", Edition: "5e", ClassID: &id,
			Description: "Ao completar um Descanso Curto, recupera Pontos de Feitiçaria gastos, até metade do seu nível de Feiticeiro (arred. p/ baixo). 1 uso por Descanso Longo.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recupera parte dos Pontos de Feitiçaria em Descanso Curto.",
			PowerType: domain.PowerDaily, Level: 5, IsClassFeature: true,
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}

	origens := []domain.Skill{
		{
			Name: "Fala Telepática", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, estabelece telepatia com uma criatura à vista a até 9m, comunicando-se a 1,5km × mod. de Carisma de distância, por minutos iguais ao seu nível de Feiticeiro.",
			Keywords: "Arcano, Aberrante", ActionType: "Ação Bônus", Range: "9 metros",
			Effect: "Telepatia de longo alcance com uma criatura.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Feitiçaria Aberrante",
		},
		{
			Name: "Defesas Psíquicas", Edition: "5e", ClassID: &id,
			Description: "Resistência a Dano Psíquico e Vantagem em salvaguardas para evitar ou encerrar Amedrontado ou Enfeitiçado.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência a Psíquico + Vantagem vs Amedrontado/Enfeitiçado.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Feitiçaria Aberrante",
		},
		{
			Name: "Feitiçaria Psiônica", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar magias da lista de Magias Psiônicas gastando Pontos de Feitiçaria (igual ao círculo) em vez de um espaço de magia; ao fazer isso, dispensa componentes Verbais, Somáticos e Materiais não-consumíveis/sem custo.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Conjuração silenciosa e sem gestos usando Pontos de Feitiçaria.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Feitiçaria Aberrante",
		},
		{
			Name: "Revelação em Carne", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, gasta 1+ Pontos de Feitiçaria para se transformar por 10 minutos, ganhando 1 benefício por ponto gasto: Adaptação Aquática, Movimento Vermiforme, Ver o Invisível ou Voo Reluzente.",
			Keywords: "Arcano, Aberrante", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Transformação corporal com benefícios escaláveis por Ponto de Feitiçaria.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Feitiçaria Aberrante",
		},
		{
			Name: "Implosão de Distorção", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, teleporta-se até 36m; criaturas a até 9m do espaço deixado salvam Força ou sofrem 3d10 de dano Energético e são puxadas para lá (metade do dano se passar). Recarrega em Descanso Longo ou gastando 5 Pontos de Feitiçaria.",
			Keywords: "Arcano, Aberrante", ActionType: "Passiva", Range: "36 metros",
			Effect: "Teleporte pessoal + dano/puxão em área no ponto de origem.",
			PowerType: domain.PowerDaily, Level: 18, IsClassFeature: true, ChoiceGroup: "Feitiçaria Aberrante",
		},
		{
			Name: "Resiliência Dracônica", Edition: "5e", ClassID: &id,
			Description: "PV máximos +3 (e +1 a cada nível de Feiticeiro seguinte). Sem armadura, CA = 10 + Destreza + Carisma.",
			Keywords: "Arcano, Dracônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Mais PV máximo + CA alternativa sem armadura.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Feitiçaria Dracônica",
		},
		{
			Name: "Afinidade Elemental", Edition: "5e", ClassID: &id,
			Description: "Escolhe um tipo de dano dracônico (Ácido/Elétrico/Gélido/Ígneo/Venenoso): Resistência a ele, e soma o mod. de Carisma ao dano de uma magia que cause esse tipo.",
			Keywords: "Arcano, Dracônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência elemental + dano bônus de Carisma no tipo escolhido.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Feitiçaria Dracônica",
		},
		{
			Name: "Asas de Dragão", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, ganha asas dracônicas por 1 hora: Deslocamento de Voo de 18m. Recarrega em Descanso Longo, ou gastando 3 Pontos de Feitiçaria.",
			Keywords: "Arcano, Dracônico", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Voo temporário de 18m.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true, ChoiceGroup: "Feitiçaria Dracônica",
		},
		{
			Name: "Companheiro Dracônico", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar Invocar Dragão sem componente Material, e 1x sem espaço de magia (recarrega em Descanso Longo); pode dispensar a Concentração dela (duração vira 1 minuto).",
			Keywords: "Arcano, Dracônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Invocar Dragão gratuito e sem Concentração.",
			PowerType: domain.PowerDaily, Level: 18, IsClassFeature: true, ChoiceGroup: "Feitiçaria Dracônica",
		},
		{
			Name: "Restaurar Equilíbrio", Edition: "5e", ClassID: &id,
			Description: "Como Reação, quando uma criatura à vista a até 18m vai rolar um d20 com Vantagem ou Desvantagem, pode anular esse efeito. Usos = mod. de Carisma (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Arcano, Ordem", ActionType: "Reação", Range: "18 metros",
			Effect: "Remove Vantagem/Desvantagem de uma jogada de d20 alheia.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Feitiçaria Mecânica",
		},
		{
			Name: "Bastião da Lei", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta 1-5 Pontos de Feitiçaria para dar a si ou a uma criatura à vista a até 9m uma proteção de Xd8 (X = pontos gastos); ao sofrer dano, a criatura pode gastar dados dessa reserva para reduzir o dano. Dura até Descanso Longo ou até reusar.",
			Keywords: "Arcano, Ordem", ActionType: "Passiva", Range: "9 metros",
			Effect: "Reserva de redução de dano escalável por Pontos de Feitiçaria.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Feitiçaria Mecânica",
		},
		{
			Name: "Transe da Ordem", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, entra em transe por 1 minuto: ataques contra você não podem ter Vantagem, e trata 9 ou menos no d20 como 10. Recarrega em Descanso Longo, ou gastando 5 Pontos de Feitiçaria.",
			Keywords: "Arcano, Ordem", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Anula Vantagem contra você + piso de 10 em testes de d20, por 1 minuto.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true, ChoiceGroup: "Feitiçaria Mecânica",
		},
		{
			Name: "Cavalgada Mecânica", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, convoca espíritos de ordem num Cubo de 9m: podem curar até 100 PV divididos entre criaturas, dissipar magias de 6º círculo ou menor, e reparar objetos danificados na área. Recarrega em Descanso Longo, ou gastando 7 Pontos de Feitiçaria.",
			Keywords: "Arcano, Ordem", ActionType: "Passiva", Range: "9 metros",
			Effect: "Cura em massa + dissipar magia + reparo de objetos em área.",
			PowerType: domain.PowerDaily, Level: 18, IsClassFeature: true, ChoiceGroup: "Feitiçaria Mecânica",
		},
		{
			Name: "Marés do Caos", Edition: "5e", ClassID: &id,
			Description: "Antes de rolar, garante Vantagem num Teste de D20 à escolha. Recarrega ao conjurar uma magia de Feiticeiro com espaço de magia (ou Descanso Longo) — o que também dispara um Surto de Magia Selvagem automaticamente.",
			Keywords: "Arcano, Caótico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Vantagem garantida, recarregando ao conjurar com espaço de magia.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Feitiçaria Selvagem",
		},
		{
			Name: "Surto de Magia Selvagem", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao conjurar magia de Feiticeiro com espaço de magia, pode jogar 1d20 — em 20, rola na tabela de Surtos de Magia Selvagem para um efeito extra imprevisível (não afetado por Metamagia se for uma magia).",
			Keywords: "Arcano, Caótico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Chance de efeito mágico caótico extra ao conjurar.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Feitiçaria Selvagem",
		},
		{
			Name: "Distorcer a Sorte", Edition: "5e", ClassID: &id,
			Description: "Como Reação, após ver outra criatura à vista rolar um d20, gasta 1 Ponto de Feitiçaria para jogar 1d4 e somar ou subtrair do resultado dela, à escolha.",
			Keywords: "Arcano, Caótico", ActionType: "Reação", Range: "Vista",
			Effect: "Altera retroativamente o resultado de um d20 alheio.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Feitiçaria Selvagem",
		},
		{
			Name: "Caos Controlado", Edition: "5e", ClassID: &id,
			Description: "Ao rolar na tabela de Surto de Magia Selvagem, rola duas vezes e escolhe qual resultado usar.",
			Keywords: "Arcano, Caótico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Escolhe entre 2 rolagens na tabela de Surto de Magia Selvagem.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Feitiçaria Selvagem",
		},
	}
	for _, s := range origens {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Feiticeiro 5e: características seedadas")
}

// ── Guardião ──────────────────────────────────────────────────────────────────

func seedGuardiao5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Guardião")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Maestria em Armas",
			Edition: "5e", ClassID: &id,
			Description: "Você pode usar as propriedades de Maestria de 3 tipos de armas Simples ou Marciais à sua escolha. Sempre que completar um Descanso Longo pode substituir uma das escolhas.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Desbloqueia propriedades de Maestria em 3 armas escolhidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Conjuração de Guardião",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Sabedoria. Nível 1: 2 truques, 2 espaços de 1° círculo. Magias preparadas = nível de Guardião + modificador de Sabedoria.",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a truques e espaços de magia de Guardião.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Inimigo Favorito", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Marca do Predador preparada, conjurável 2x sem espaço de magia (recarrega em Descanso Longo).",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Marca do Predador gratuita, 2x por Descanso Longo.",
			LevelScaling: "Usos sem espaço de magia aumentam em níveis mais altos de Guardião.",
			PowerType:    domain.PowerDaily, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Explorador Hábil", Edition: "5e", ClassID: &id,
			Description: "Escolhe uma perícia em que já é proficiente (mas não Especialista) e ganha Especialização nela. Também aprende 2 idiomas à escolha.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Especialização em 1 perícia + 2 idiomas novos.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todas as subclasses) ──
		{
			Name:    "Ataque Extra (Guardião)", Edition: "5e", ClassID: &id,
			Description: "Você pode atacar duas vezes, em vez de uma, sempre que executar a ação Atacar no seu turno.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 ataques por ação Atacar.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Errante", Edition: "5e", ClassID: &id,
			Description: "Deslocamento +3m sem Armadura Pesada. Ganha Deslocamento de Escalada e de Natação iguais ao seu Deslocamento.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+3m de Deslocamento + Escalada/Natação iguais ao Deslocamento.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true,
		},
		{
			Name:    "Especialista (Guardião)", Edition: "5e", ClassID: &id,
			Description: "Escolhe 2 perícias em que já é proficiente (mas não Especialista) e ganha Especialização nelas.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Especialização em mais 2 perícias.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Incansável", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, ganha 1d8 + mod. de Sabedoria (mín. 1) de PV Temporários; usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo. Além disso, reduz seu nível de Exaustão em 1 a cada Descanso Curto.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "PV Temporários sob demanda + recuperação de Exaustão em Descanso Curto.",
			PowerType: domain.PowerDaily, Level: 10, IsClassFeature: true,
		},
		{
			Name:    "Predador Implacável", Edition: "5e", ClassID: &id,
			Description: "Sofrer dano não quebra mais sua Concentração na magia Marca do Predador.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Marca do Predador imune a interrupção por dano.",
			PowerType: domain.PowerUnlimited, Level: 13, IsClassFeature: true,
		},
		{
			Name:    "Véu da Natureza", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, fica Invisível até o final do próximo turno. Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:    "Invisibilidade rápida e reutilizável.",
			PowerType: domain.PowerDaily, Level: 14, IsClassFeature: true,
		},
		{
			Name:    "Caçador Preciso", Edition: "5e", ClassID: &id,
			Description: "Vantagem em jogadas de ataque contra a criatura marcada pela sua Marca do Predador.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem contra o alvo marcado.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true,
		},
		{
			Name:    "Sentidos Selvagens", Edition: "5e", ClassID: &id,
			Description: "Ganha Visão às Cegas com alcance de 9 metros.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Visão às Cegas 9m.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Matador de Inimigos Favoritos", Edition: "5e", ClassID: &id,
			Description: "O dado de dano da sua Marca do Predador vira d10 em vez de d6.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Dado de dano de Marca do Predador sobe para d10.",
			PowerType: domain.PowerUnlimited, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		{
			Name: "Andarilho Feérico", Edition: "5e", ClassID: &id,
			Description: "Empunhe o Deleite e a Fúria Feérica — uma das 4 subclasses de Guardião, escolhida no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "subclasse_guardiao",
		},
		{
			Name: "Caçador", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 subclasses de Guardião, escolhida no nível 3.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "subclasse_guardiao",
		},
		{
			Name: "Senhor das Feras", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 subclasses de Guardião, escolhida no nível 3.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "subclasse_guardiao",
		},
		{
			Name: "Vigilante das Sombras", Edition: "5e", ClassID: &id,
			Description: "Aproveite a Magia das Sombras para Lutar contra Seus Inimigos — Vigilantes das Sombras empunham magia extraída do Sombral para combater inimigos escondidos na escuridão.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "subclasse_guardiao",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/7/11/15) ──────────────────────
		{
			Name: "Glamour Transcendental", Edition: "5e", ClassID: &id,
			Description: "Soma o mod. de Sabedoria (mín. +1) em testes de Carisma. Ganha proficiência em Atuação, Enganação ou Persuasão, à escolha.",
			Keywords: "Primitivo, Feérico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Bônus de Carisma via Sabedoria + 1 perícia social.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Andarilho Feérico",
		},
		{
			Name: "Golpes Terríveis", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao acertar com arma, causa 1d4 de dano Psíquico adicional (sobe para 1d6 no nível 11).",
			Keywords: "Primitivo, Feérico", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 11: dano extra sobe para 1d6.",
			Effect: "Dano Psíquico extra 1x por turno.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Andarilho Feérico",
		},
		{
			Name: "Detalhe Sedutor", Edition: "5e", ClassID: &id,
			Description: "Vantagem contra Amedrontado/Enfeitiçado. Quando você ou alguém à vista a até 36m resiste a uma dessas condições, pode usar Reação para forçar outra criatura a salvar (Sabedoria) ou ficar Amedrontada/Enfeitiçada por 1 minuto.",
			Keywords: "Primitivo, Feérico", ActionType: "Reação", Range: "36 metros",
			Effect: "Redireciona resistência a medo/encanto como ataque em outro alvo.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Andarilho Feérico",
		},
		{
			Name: "Reforços Feéricos", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar Convocar Feérico sem componente Material, e 1x sem espaço de magia (recarrega em Descanso Longo), podendo dispensar a Concentração (duração vira 1 minuto).",
			Keywords: "Primitivo, Feérico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Convocar Feérico gratuito e sem Concentração.",
			PowerType: domain.PowerDaily, Level: 11, IsClassFeature: true, ChoiceGroup: "Andarilho Feérico",
		},
		{
			Name: "Andarilho Nebuloso", Edition: "5e", ClassID: &id,
			Description: "Pode conjurar Passo Nebuloso sem espaço de magia (usos = mod. Sabedoria, mín. 1, recarrega em Descanso Longo), podendo levar uma criatura voluntária junto a até 1,5m.",
			Keywords: "Primitivo, Feérico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Passo Nebuloso gratuito, em dupla.",
			PowerType: domain.PowerDaily, Level: 15, IsClassFeature: true, ChoiceGroup: "Andarilho Feérico",
		},
		{
			Name: "Conhecimento do Caçador", Edition: "5e", ClassID: &id,
			Description: "Enquanto uma criatura está marcada pela Marca do Predador, sabe se ela tem Imunidades, Resistências ou Vulnerabilidades, e quais são.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Revela imunidades/resistências/vulnerabilidades do alvo marcado.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Caçador",
		},
		{
			Name: "Presa do Caçador", Edition: "5e", ClassID: &id,
			Description: "Escolhe: Assassino de Colossos (1x por turno, +1d8 de dano ao acertar criatura abaixo do PV máximo) ou Destruidor de Hordas (1x por turno, ataque extra contra um alvo diferente a 1,5m do original). Troca a escolha em Descanso Curto ou Longo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano extra contra feridos, ou ataque extra em outro alvo.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Caçador",
		},
		{
			Name: "Táticas Defensivas", Edition: "5e", ClassID: &id,
			Description: "Escolhe: Defesa Contra Ataques Múltiplos (quem te acerta tem Desvantagem nos outros ataques contra você naquele turno) ou Escapar de Hordas (Ataques de Oportunidade contra você têm Desvantagem). Troca em Descanso Curto ou Longo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Mitigação contra múltiplos atacantes ou fuga segura.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Caçador",
		},
		{
			Name: "Presa do Caçador Superior", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao causar dano ao alvo marcado pela Marca do Predador, também causa esse dano a outra criatura à vista a até 9m do alvo original.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "9 metros",
			Effect: "Dano da Marca do Predador se espalha para um segundo alvo.",
			PowerType: domain.PowerUnlimited, Level: 11, IsClassFeature: true, ChoiceGroup: "Caçador",
		},
		{
			Name: "Defesa do Caçador Superior", Edition: "5e", ClassID: &id,
			Description: "Ao sofrer dano, pode usar Reação para ganhar Resistência a esse tipo de dano (e a qualquer outro do mesmo tipo) até o final do turno.",
			Keywords: "Marcial", ActionType: "Reação", Range: "Pessoal",
			Effect: "Resistência reativa ao tipo de dano sofrido.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Caçador",
		},
		{
			Name: "Companheiro Primal", Edition: "5e", ClassID: &id,
			Description: "Invoca uma fera primal aliada (Fera da Terra, do Céu ou do Mar, à escolha), amigável e obediente, que age no seu turno (Ação Bônus para comandá-la; sacrifica um dos seus ataques para ordenar 'Golpe da Fera'). Some seu Bônus de Proficiência aos testes/salvaguardas dela; pode restaurá-la ou trocá-la a cada Descanso Longo.",
			Keywords: "Primitivo", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Companheiro animal invocado, com estatísticas escalando pelo seu nível.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Senhor das Feras",
		},
		{
			Name: "Treinamento Excepcional", Edition: "5e", ClassID: &id,
			Description: "Pode ordenar sua Fera a executar Ajudar, Correr, Desengajar ou Esquivar como Ação Bônus dela. A fera também pode trocar o dano do Golpe da Fera por dano Energético à escolha.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Mais opções táticas e flexibilidade de dano para a Fera.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Senhor das Feras",
		},
		{
			Name: "Fúria Bestial", Edition: "5e", ClassID: &id,
			Description: "Ao ordenar Golpe da Fera, ela pode usá-lo duas vezes. A cada turno, na primeira vez que acertar um alvo marcado por Marca do Predador, causa dano Energético extra igual ao bônus dessa magia.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "2 Golpes da Fera por ordem + dano extra em alvo marcado.",
			PowerType: domain.PowerUnlimited, Level: 11, IsClassFeature: true, ChoiceGroup: "Senhor das Feras",
		},
		{
			Name: "Compartilhar Magias", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia em si mesmo, também pode afetar sua Fera com ela se estiver a até 9m.",
			Keywords: "Primitivo, Magia", ActionType: "Passiva", Range: "9 metros",
			Effect: "Magias pessoais passam a afetar a Fera também.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Senhor das Feras",
		},
		{
			Name: "Emboscador das Sombras", Edition: "5e", ClassID: &id,
			Description: "Soma Sabedoria à Iniciativa. Golpe Terrível: 1x por turno, +2d6 de dano Psíquico ao acertar com arma (usos = mod. Sabedoria, mín. 1, recarrega em Descanso Longo). Impulso do Emboscador: +3m de Deslocamento no primeiro turno de cada combate.",
			Keywords: "Arcano, Sombrio", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Iniciativa melhorada + dano psíquico de emboscada + arrancada inicial.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Vigilante das Sombras",
		},
		{
			Name: "Visão Umbrosa", Edition: "5e", ClassID: &id,
			Description: "Ganha Visão no Escuro 18m (ou +18m se já tiver). Enquanto totalmente na Escuridão, fica Invisível para quem depende de Visão no Escuro para vê-lo.",
			Keywords: "Arcano, Sombrio", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Visão no escuro melhorada + invisibilidade na escuridão total.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Vigilante das Sombras",
		},
		{
			Name: "Mente de Ferro", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência em salvaguardas de Sabedoria; se já tiver, ganha em Carisma ou Inteligência à escolha.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+1 proficiência de salvaguarda mental.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Vigilante das Sombras",
		},
		{
			Name: "Torrente do Vigilante", Edition: "5e", ClassID: &id,
			Description: "Golpe Terrível sobe para 2d8 de dano Psíquico. Ao usá-lo, pode escolher um efeito extra: Golpe Repentino (ataque extra em alvo próximo) ou Medo em Massa (área ao redor do alvo salva Sabedoria ou fica Amedrontada).",
			Keywords: "Arcano, Sombrio", ActionType: "Passiva", Range: "3 metros",
			Effect: "Golpe Terrível mais forte, com efeito extra à escolha.",
			PowerType: domain.PowerUnlimited, Level: 11, IsClassFeature: true, ChoiceGroup: "Vigilante das Sombras",
		},
		{
			Name: "Esquiva Sombria", Edition: "5e", ClassID: &id,
			Description: "Como Reação a um ataque contra você, impõe Desvantagem nele; acerte ou erre, pode se teleportar até 9m para um espaço à vista.",
			Keywords: "Arcano, Sombrio", ActionType: "Reação", Range: "9 metros",
			Effect: "Desvantagem no ataque inimigo + teleporte defensivo.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Vigilante das Sombras",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Guardião 5e: características seedadas")
}

// ── Guerreiro ─────────────────────────────────────────────────────────────────

func seedGuerreiro5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Guerreiro")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Maestria em Armas",
			Edition: "5e", ClassID: &id,
			Description: "Você pode usar as propriedades de Maestria de 3 tipos de armas Simples ou Marciais à sua escolha. Sempre que completar um Descanso Longo pode substituir uma das escolhas.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Desbloqueia propriedades de Maestria em 3 armas escolhidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Estilo de Luta",
			Edition: "5e", ClassID: &id,
			Description: "Você adota um estilo de luta como especialidade. Escolha um talento de Estilo de Luta: Arquearia (+2 em ataques à distância), Duelismo (+2 dano com uma mão), Combate com Armas Grandes (trate 1 ou 2 como 3 no dado de dano), Defensivo (+1 CA com armadura), Interceptação, Protetivo, entre outros.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Concede bônus passivo conforme o estilo escolhido.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "estilo_luta_guerreiro",
		},
		{
			Name:    "Recuperar Fôlego",
			Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, cura 1d10 + nível de Guerreiro PV. 2 usos — recupera 1 em Descanso Curto, todos em Descanso Longo.",
			Keywords: "Marcial, Cura", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:       "Cura 1d10 + nível de Guerreiro PV, 2x por Descanso Curto/Longo.",
			LevelScaling: "Mais usos em níveis altos de Guerreiro (coluna Recuperar Fôlego da tabela).",
			PowerType:    domain.PowerEncounter, Level: 1, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todos os arquétipos) ──
		{
			Name:    "Mente Tática", Edition: "5e", ClassID: &id,
			Description: "Ao falhar num teste de atributo, pode gastar um uso de Recuperar Fôlego: joga 1d10 e soma ao teste (em vez de curar). Se ainda assim falhar, o uso não é gasto.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Converte um uso de Recuperar Fôlego em bônus de teste de atributo.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Surto de Ação", Edition: "5e", ClassID: &id,
			Description: "No seu turno, pode executar uma ação adicional (exceto Usar Magia). 1 uso por Descanso Curto ou Longo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "1 ação extra no turno.",
			LevelScaling: "Nível 17: 2 usos antes de descansar (só 1 por turno).",
			PowerType:    domain.PowerEncounter, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Ajuste Tático", Edition: "5e", ClassID: &id,
			Description: "Ao usar a Ação Bônus de Recuperar Fôlego, também pode se mover até metade do Deslocamento sem provocar Ataques de Oportunidade.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Movimento livre combinado com Recuperar Fôlego.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Ataque Extra (Guerreiro)", Edition: "5e", ClassID: &id,
			Description: "Você pode atacar duas vezes, em vez de uma, sempre que executar a ação Atacar no seu turno.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "2 ataques por ação Atacar.",
			LevelScaling: "Nível 11 (Dois Ataques Extras): 3 ataques. Nível 20 (Três Ataques Extras): 4 ataques.",
			PowerType:    domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Indomável", Edition: "5e", ClassID: &id,
			Description: "Ao falhar numa salvaguarda, pode rejogá-la somando um bônus igual ao seu nível de Guerreiro, usando o novo resultado. 1 uso por Descanso Longo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Rejogada garantida de salvaguarda com bônus de nível.",
			LevelScaling: "Nível 13: 2 usos por Descanso Longo. Nível 17: 3 usos.",
			PowerType:    domain.PowerDaily, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Mestre Tático", Edition: "5e", ClassID: &id,
			Description: "Ao atacar com uma arma cuja maestria você pode usar, pode substituir a propriedade dela por Empurrar, Drenar ou Lentidão nesse ataque.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Troca a propriedade de maestria da arma por outra à escolha, por ataque.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Ataques Estudados", Edition: "5e", ClassID: &id,
			Description: "Se errar um ataque contra uma criatura, ganha Vantagem no próximo ataque contra ela antes do final do seu próximo turno.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem garantida após errar um ataque no mesmo alvo.",
			PowerType: domain.PowerUnlimited, Level: 13, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		{
			Name: "Campeão", Edition: "5e", ClassID: &id,
			Description: "Busque a Excelência Física em Combate — o Campeão foca no desenvolvimento de habilidades marciais em sua busca incessante pela vitória, combinando treino rigoroso com excelência física.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_guerreiro",
		},
		{
			Name: "Cavaleiro Místico", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Arquétipos Marciais de Guerreiro, escolhido no nível 3.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_guerreiro",
		},
		{
			Name: "Combatente Psíquico", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Arquétipos Marciais de Guerreiro, escolhido no nível 3.",
			Keywords: "Marcial, Psiônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_guerreiro",
		},
		{
			Name: "Mestre da Batalha", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Arquétipos Marciais de Guerreiro, escolhido no nível 3.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_guerreiro",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/7/10/15/18) ───────────────────
		{
			Name: "Atleta Extraordinário", Edition: "5e", ClassID: &id,
			Description: "Vantagem em Iniciativa e em testes de Força (Atletismo). Após um Acerto Crítico, pode se mover até metade do Deslocamento sem provocar Ataques de Oportunidade.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Vantagem em Iniciativa/Atletismo + movimento livre após crítico.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Campeão",
		},
		{
			Name: "Crítico Aprimorado", Edition: "5e", ClassID: &id,
			Description: "Seus ataques com arma ou Desarmados causam Acerto Crítico com 19 ou 20 no d20.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 15 (Crítico Superior): faixa de crítico expande para 18-20.",
			Effect: "Faixa de crítico ampliada para 19-20.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Campeão",
		},
		{
			Name: "Estilo de Luta Adicional", Edition: "5e", ClassID: &id,
			Description: "Ganha um segundo talento de Estilo de Luta à escolha.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+1 Estilo de Luta.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Campeão",
		},
		{
			Name: "Combatente Heroico", Edition: "5e", ClassID: &id,
			Description: "Em combate, pode se conceder Inspiração Heroica sempre que começar seu turno sem ela.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Inspiração Heroica praticamente garantida em combate.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Campeão",
		},
		{
			Name: "Sobrevivente", Edition: "5e", ClassID: &id,
			Description: "Desafie a Morte: Vantagem em Salvaguardas Contra a Morte, e 18-20 nelas conta como 20. Regeneração Heroica: no início do turno, recupera 5 + Constituição de PV se estiver Sangrando com ao menos 1 PV.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência forte contra a morte + regeneração passiva quando ferido.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true, ChoiceGroup: "Campeão",
		},
		{
			Name: "Conjuração de Cavaleiro Místico", Edition: "5e", ClassID: &id,
			Description: "Conjura magias de Mago (Inteligência é o atributo). Começa com 2 truques e 3 magias de 1º círculo preparadas, ganhando mais em níveis altos, com espaços de magia próprios (máx. 4º círculo).",
			Keywords: "Marcial, Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Meia-conjuração de Mago integrada ao Guerreiro.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Cavaleiro Místico",
		},
		{
			Name: "Vínculo com Arma", Edition: "5e", ClassID: &id,
			Description: "Ritual de 1 hora vincula você a uma arma (até 2 simultâneas): ela não pode ser desarmada de você, e pode ser invocada à mão como Ação Bônus se estiver no mesmo plano.",
			Keywords: "Marcial, Arcano", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Arma pessoal teleportável à mão e impossível de desarmar.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Cavaleiro Místico",
		},
		{
			Name: "Magia de Guerra", Edition: "5e", ClassID: &id,
			Description: "Ao executar a ação Atacar, pode substituir um dos ataques por um truque de Mago de 1 ação.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 18 (Magia de Guerra Aprimorada): pode substituir 2 ataques por uma magia de 1º ou 2º círculo de 1 ação.",
			Effect: "Combina ataques com truques na mesma ação Atacar.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Cavaleiro Místico",
		},
		{
			Name: "Golpe Místico", Edition: "5e", ClassID: &id,
			Description: "Ao acertar um ataque com arma, o alvo tem Desvantagem na próxima salvaguarda contra uma magia sua até o final do seu próximo turno.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Enfraquece a resistência mágica do alvo após um golpe físico.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Cavaleiro Místico",
		},
		{
			Name: "Investida Mística", Edition: "5e", ClassID: &id,
			Description: "Ao usar Surto de Ação, pode se teleportar até 9m para um espaço à vista, antes ou depois da ação extra.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "9 metros",
			Effect: "Teleporte combinado com Surto de Ação.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Cavaleiro Místico",
		},
		{
			Name: "Poder Psiônico", Edition: "5e", ClassID: &id,
			Description: "Reserva de Dados de Energia Psiônica (d6 a d12, crescendo com o nível). Golpe Psiônico: 1x por turno, gasta um dado para +dano Energético num ataque. Movimento Telecinético: move objeto/criatura 9m como ação. Vínculo Protetivo: Reação para reduzir dano sofrido por si ou aliado próximo.",
			Keywords: "Marcial, Psiônico", ActionType: "Passiva", Range: "9 metros",
			Effect: "Recurso de dados psiônicos alimentando dano extra, telecinesia e proteção.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente Psíquico",
		},
		{
			Name: "Adepto Telecinético", Edition: "5e", ClassID: &id,
			Description: "Estocada Telecinética: ao acertar com Golpe Psiônico, força salvaguarda de Força ou empurra/derruba o alvo. Salto com Impulsão Psíquica: Ação Bônus dá Voo igual ao dobro do Deslocamento até o fim do turno.",
			Keywords: "Marcial, Psiônico", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Controle de posição (empurrar/derrubar) + voo curto de escape/perseguição.",
			PowerType: domain.PowerEncounter, Level: 7, IsClassFeature: true, ChoiceGroup: "Combatente Psíquico",
		},
		{
			Name: "Resguardo Mental", Edition: "5e", ClassID: &id,
			Description: "Resistência a Dano Psíquico. Ao começar o turno Amedrontado ou Enfeitiçado, pode gastar um Dado de Energia Psiônica para encerrar o efeito.",
			Keywords: "Marcial, Psiônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Resistência psíquica + cura de condições mentais sob demanda.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Combatente Psíquico",
		},
		{
			Name: "Baluarte de Energia", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, dá Cobertura Parcial por 1 minuto a criaturas a até 9m (número = mod. Inteligência, mín. 1). 1 uso por Descanso Longo, ou gaste um Dado de Energia Psiônica para restaurar.",
			Keywords: "Marcial, Psiônico", ActionType: "Ação Bônus", Range: "9 metros",
			Effect: "Cobertura em grupo via telecinese.",
			PowerType: domain.PowerDaily, Level: 15, IsClassFeature: true, ChoiceGroup: "Combatente Psíquico",
		},
		{
			Name: "Mestre Telecinético", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Telecinese preparada, conjurável sem espaço de magia ou componentes (Inteligência como atributo); enquanto a mantém, pode atacar com arma como Ação Bônus todo turno. 1 uso desse modo por Descanso Longo, ou gaste um Dado de Energia Psiônica.",
			Keywords: "Marcial, Psiônico, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Telecinese gratuita combinada com ataques extras.",
			PowerType: domain.PowerDaily, Level: 18, IsClassFeature: true, ChoiceGroup: "Combatente Psíquico",
		},
		{
			Name: "Estudioso da Guerra", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência com um tipo de Ferramentas de Artesão e em uma perícia disponível para Guerreiros no nível 1.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "+1 ferramenta + 1 perícia.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Mestre da Batalha",
		},
		{
			Name: "Superioridade em Combate", Edition: "5e", ClassID: &id,
			Description: "Aprende 3 manobras (de um catálogo do capítulo, ex.: Aparar, Ataque Ameaçador, Ataque de Varredura) alimentadas por 4 Dados de Superioridade d8 (recarregam em Descanso Curto ou Longo). Aprende +2 manobras nos níveis 7, 10 e 15.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Catálogo de manobras táticas alimentado por uma reserva de dados.",
			LevelScaling: "Dado de Superioridade: d8 (nível 3) → d10 (nível 10) → d12 (nível 18). Reserva: 4 dados → 5 (nível 7) → 6 (nível 15).",
			PowerType:    domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Mestre da Batalha",
		},
		{
			Name: "Conheça Seu Inimigo", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, revela Imunidades, Resistências e Vulnerabilidades de uma criatura à vista a até 9m. 1 uso por Descanso Longo, ou gaste um Dado de Superioridade para restaurar.",
			Keywords: "Marcial", ActionType: "Ação Bônus", Range: "9 metros",
			Effect: "Revela pontos fracos/fortes de um inimigo.",
			PowerType: domain.PowerDaily, Level: 7, IsClassFeature: true, ChoiceGroup: "Mestre da Batalha",
		},
		{
			Name: "Implacável", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao usar uma manobra, pode jogar 1d8 em vez de gastar um Dado de Superioridade.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "1 manobra grátis por turno.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Mestre da Batalha",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Guerreiro 5e: características seedadas")
}

// ── Ladino ────────────────────────────────────────────────────────────────────

func seedLadino5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Ladino")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Especialista",
			Edition: "5e", ClassID: &id,
			Description: "Escolha 2 perícias nas quais você já tenha proficiência. Você recebe Especialização nessas perícias: seu Bônus de Proficiência é dobrado nos testes que as usem.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Bônus de Proficiência dobrado em 2 perícias escolhidas.",
			LevelScaling:   "Nível 6: Especialização em mais 2 perícias em que já é proficiente.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "especialista_ladino",
		},
		{
			Name:    "Ataque Furtivo",
			Edition: "5e", ClassID: &id,
			Description: "Uma vez por turno, você pode causar 1d6 de dano extra ao atingir uma criatura com um ataque usando uma arma de Finesse ou à Distância — desde que tenha Vantagem na jogada ou um aliado esteja adjacente ao alvo e você não tenha Desvantagem.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Dano extra de 1d6 com vantagem ou aliado adjacente ao alvo.",
			LevelScaling: "Nível 3: 2d6. Nível 5: 3d6. Nível 7: 4d6. Nível 9: 5d6. Nível 11: 6d6.",
			PowerType:    domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Linguagem dos Ladrões",
			Edition: "5e", ClassID: &id,
			Description: "Você conhece a Linguagem dos Ladrões — um código secreto de sinais e mensagens que criminosos e ladinos usam para se comunicar discretamente sem que outros percebam.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Compreende e usa a Linguagem dos Ladrões.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Maestria em Armas",
			Edition: "5e", ClassID: &id,
			Description: "Você pode usar as propriedades de Maestria de 2 tipos de armas Simples ou Marciais à sua escolha.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Desbloqueia propriedades de Maestria em 2 armas escolhidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todas as subclasses) ──
		{
			Name:    "Ação Ardilosa", Edition: "5e", ClassID: &id,
			Description: "No seu turno, pode executar Correr, Desengajar ou Esconder como Ação Bônus.",
			Keywords: "Marcial", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:    "1 dessas 3 ações extra como Ação Bônus todo turno.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Mira Firme", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus (só se não tiver se movido no turno), ganha Vantagem no próximo ataque; seu Deslocamento vira 0 até o fim do turno.",
			Keywords: "Marcial", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:    "Vantagem garantida num ataque, ao custo de não se mover.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true,
		},
		{
			Name:    "Golpe Astuto", Edition: "5e", ClassID: &id,
			Description: "Ao causar dano com Ataque Furtivo, pode trocar parte dos dados por um efeito: Envenenar (1d6, Constituição ou Envenenado 1 min), Retirada (1d6, move metade do Deslocamento sem Oportunidade) ou Tropeço (1d6, Destreza ou Caído).",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "Converte dados de Ataque Furtivo em efeitos de controle/utilidade.",
			LevelScaling: "Nível 11 (Golpe Astuto Aprimorado): até 2 efeitos por ataque. Nível 14 (Golpes Sujos): +Aturdir, Nocaute e Obscurecer como novas opções.",
			PowerType:    domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Esquiva Sobrenatural", Edition: "5e", ClassID: &id,
			Description: "Quando um atacante à vista te acerta, pode usar Reação para reduzir o dano pela metade.",
			Keywords: "Marcial", ActionType: "Reação", Range: "Pessoal",
			Effect:    "Reduz pela metade um dano recebido.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Evasão", Edition: "5e", ClassID: &id,
			Description: "Em salvaguardas de Destreza para metade do dano: sem dano se passar, metade se falhar (não funciona se Incapacitado).",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Anula ou reduz dano em área evitável por Destreza.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Talento Confiável", Edition: "5e", ClassID: &id,
			Description: "Em testes de perícia/ferramenta com proficiência, trata um resultado de 9 ou menos no d20 como 10.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Piso de 10 em testes proficientes.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Mente Escorregadia", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência em salvaguardas de Sabedoria e Carisma.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+2 proficiências de salvaguarda mental.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true,
		},
		{
			Name:    "Elusivo", Edition: "5e", ClassID: &id,
			Description: "Nenhuma jogada de ataque pode ter Vantagem contra você, a menos que esteja Incapacitado.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Imune a Vantagem de ataques inimigos.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Golpe de Sorte", Edition: "5e", ClassID: &id,
			Description: "Ao falhar num Teste de D20, pode transformar o resultado em 20. 1 uso por Descanso Curto ou Longo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Converte uma falha garantida em sucesso natural.",
			PowerType: domain.PowerEncounter, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		// Nome "Adaga Espiritual" extraído do PDF via RAG (rag_5e.py) — nome
		// pouco comum pra uma tradução oficial, mas apareceu de forma
		// idêntica em duas buscas semânticas independentes; se estiver
		// errado, é erro de extração do PDF-fonte, não invenção.
		{
			Name: "Adaga Espiritual", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 subclasses de Ladino, escolhida no nível 3.",
			Keywords: "Psiônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_ladino",
		},
		{
			Name: "Assassino", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 subclasses de Ladino, escolhida no nível 3.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_ladino",
		},
		{
			Name: "Ladrão", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 subclasses de Ladino, escolhida no nível 3.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_ladino",
		},
		{
			Name: "Trapaceiro Arcano", Edition: "5e", ClassID: &id,
			Description: "Combina astúcia e agilidade com magia. Ao atingir o nível 3, aprende a conjurar magias de Mago (3 truques: Mãos Mágicas e mais 2 à escolha) usando Inteligência como atributo de conjuração.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3, incluindo truques de Mago.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "arquetipo_ladino",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/9/13/17) ──────────────────────
		{
			Name: "Lâminas Psíquicas", Edition: "5e", ClassID: &id,
			Description: "Ao Atacar (ou em Ataque de Oportunidade), manifesta uma lâmina de energia psíquica na mão livre: arma Simples Corpo a Corpo, Acuidade, Arremesso 18/36m, 1d6 Psíquico + mod. de atributo. Pode manifestar uma 2ª (1d4) como Ação Bônus no mesmo turno se a outra mão estiver livre.",
			Keywords: "Psiônico", ActionType: "Passiva", Range: "18-36 metros",
			Effect: "Arma psíquica invocável, com um segundo ataque bônus.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Adaga Espiritual",
		},
		{
			Name: "Poder Psiônico (Ladino)", Edition: "5e", ClassID: &id,
			Description: "Reserva de Dados de Energia Psiônica (d6 a d12, cresce com o nível). Aptidão Reforçada Psiquicamente: joga um dado para transformar falha em sucesso num teste proficiente. Sussurros Psíquicos: telepatia por horas com várias criaturas.",
			Keywords: "Psiônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Recurso de dados psiônicos para testes e telepatia.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true, ChoiceGroup: "Adaga Espiritual",
		},
		{
			Name: "Lâminas da Alma", Edition: "5e", ClassID: &id,
			Description: "Golpes Teleguiados: ao errar com a Lâmina Psíquica, pode gastar um dado para somar ao ataque. Teleporte Psíquico: Ação Bônus, arremessa a lâmina e se teleporta até 3x o resultado do dado em metros.",
			Keywords: "Psiônico", ActionType: "Ação Bônus", Range: "Variável",
			Effect: "Correção de ataque + teleporte via lâmina arremessada.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true, ChoiceGroup: "Adaga Espiritual",
		},
		{
			Name: "Véu Psíquico", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, fica Invisível por 1 hora (encerra ao causar dano ou forçar salvaguarda). Recarrega gastando um Dado de Energia Psiônica, ou em Descanso Longo.",
			Keywords: "Psiônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Invisibilidade sob demanda, quebrada por ação ofensiva.",
			PowerType: domain.PowerDaily, Level: 13, IsClassFeature: true, ChoiceGroup: "Adaga Espiritual",
		},
		{
			Name: "Rasgar Mente", Edition: "5e", ClassID: &id,
			Description: "Ao causar dano de Ataque Furtivo com a Lâmina Psíquica, força salvaguarda de Sabedoria ou o alvo fica Atordoado por 1 minuto. Recarrega gastando 3 Dados de Energia Psiônica, ou em Descanso Longo.",
			Keywords: "Psiônico", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Atordoa o alvo do Ataque Furtivo.",
			PowerType: domain.PowerDaily, Level: 17, IsClassFeature: true, ChoiceGroup: "Adaga Espiritual",
		},
		{
			Name: "Assassinar", Edition: "5e", ClassID: &id,
			Description: "Vantagem em Iniciativa. Na primeira rodada de combate, Vantagem em ataques contra quem ainda não agiu; se o Ataque Furtivo acertar nessa rodada, causa dano extra igual ao seu nível de Ladino.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Iniciativa alta + emboscada devastadora na 1ª rodada.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Assassino",
		},
		{
			Name: "Ferramentas de Assassino", Edition: "5e", ClassID: &id,
			Description: "Ganha um Kit de Disfarce e um Kit de Veneno, com proficiência em ambos.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "2 kits + proficiências.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Assassino",
		},
		{
			Name: "Especialista em Infiltração", Edition: "5e", ClassID: &id,
			Description: "Mimetismo Magistral: imita fala/caligrafia de alguém estudado por 1h. Mira Móvel: usar Mira Firme não zera mais o Deslocamento.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Disfarce social + Mira Firme sem perder mobilidade.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true, ChoiceGroup: "Assassino",
		},
		{
			Name: "Armas Venenosas", Edition: "5e", ClassID: &id,
			Description: "Ao usar Envenenar do Golpe Astuto, o alvo também sofre 2d6 de dano Venenoso ao falhar na salvaguarda, ignorando Resistência a esse dano.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano de veneno extra que ignora resistência.",
			PowerType: domain.PowerUnlimited, Level: 13, IsClassFeature: true, ChoiceGroup: "Assassino",
		},
		{
			Name: "Golpe Mortal", Edition: "5e", ClassID: &id,
			Description: "Ao acertar Ataque Furtivo na 1ª rodada de combate, o alvo salva Constituição ou o dano do ataque é dobrado.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano dobrado em emboscada bem-sucedida.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Assassino",
		},
		{
			Name: "Andarilho de Telhados", Edition: "5e", ClassID: &id,
			Description: "Ganha Deslocamento de Escalada igual ao normal, e usa Destreza (em vez de Força) para calcular distância de salto.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Escalada completa + saltos baseados em Destreza.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Ladrão",
		},
		{
			Name: "Mão Leve", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, realiza Prestidigitação (abrir fechadura, desarmar armadilha, roubar bolso) ou usa um item mágico que exija ação Usar Objeto/Usar Magia.",
			Keywords: "Marcial", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Furtos e uso de itens mágicos como Ação Bônus.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Ladrão",
		},
		{
			Name: "Furtividade Suprema", Edition: "5e", ClassID: &id,
			Description: "Nova opção de Golpe Astuto — Ataque Escondido (1d6): atacar Invisível (por Esconder) não encerra a condição se terminar o turno com Cobertura de Três Quartos ou Total.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Mantém-se escondido mesmo após atacar, com cobertura.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true, ChoiceGroup: "Ladrão",
		},
		{
			Name: "Usar Dispositivo Mágico", Edition: "5e", ClassID: &id,
			Description: "Sintoniza até 4 itens mágicos. 1-em-6 de chance de usar uma propriedade de carga sem gastar cargas. Pode usar Pergaminhos Mágicos com Inteligência (truques/1º círculo automaticamente; círculos maiores exigem teste de Arcanismo).",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Domínio avançado de itens mágicos e pergaminhos.",
			PowerType: domain.PowerUnlimited, Level: 13, IsClassFeature: true, ChoiceGroup: "Ladrão",
		},
		{
			Name: "Reflexos de Ladrão", Edition: "5e", ClassID: &id,
			Description: "Na primeira rodada de combate, age duas vezes: uma na Iniciativa normal, outra na Iniciativa −10.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Turno extra na primeira rodada de cada combate.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Ladrão",
		},
		{
			Name: "Conjuração de Trapaceiro Arcano", Edition: "5e", ClassID: &id,
			Description: "Conjura magias de Mago (Inteligência): começa com Mãos Mágicas + 2 truques e 3 magias de 1º círculo preparadas, crescendo com o nível (espaços até 4º círculo).",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Meia-conjuração de Mago integrada ao Ladino.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trapaceiro Arcano",
		},
		{
			Name: "Mãos Mágicas Ligeiras", Edition: "5e", ClassID: &id,
			Description: "Conjura Mãos Mágicas como Ação Bônus e pode torná-la Invisível; controla a mão como Ação Bônus, podendo fazer testes de Destreza (Prestidigitação) com ela.",
			Keywords: "Arcano", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Mãos Mágicas invisível e mais rápida, usável para furtos remotos.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Trapaceiro Arcano",
		},
		{
			Name: "Emboscada Mágica", Edition: "5e", ClassID: &id,
			Description: "Se estiver Invisível ao conjurar uma magia num alvo, ele tem Desvantagem em qualquer salvaguarda contra ela naquele turno.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Desvantagem de salvaguarda contra magias lançadas enquanto invisível.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true, ChoiceGroup: "Trapaceiro Arcano",
		},
		{
			Name: "Trapaceiro Versátil", Edition: "5e", ClassID: &id,
			Description: "Ao usar uma opção de Golpe Astuto, também pode aplicá-la a uma segunda criatura a até 1,5m da Mão Mágica.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "1,5 metro",
			Effect: "Golpe Astuto duplicado via Mãos Mágicas.",
			PowerType: domain.PowerUnlimited, Level: 13, IsClassFeature: true, ChoiceGroup: "Trapaceiro Arcano",
		},
		{
			Name: "Ladrão de Magias", Edition: "5e", ClassID: &id,
			Description: "Como Reação a ser alvo/incluído numa magia, força salvaguarda de Inteligência no conjurador; se falhar, nega o efeito e rouba a magia (se de círculo conjurável por você), tendo-a preparada por 8 horas enquanto o alvo não pode conjurá-la.",
			Keywords: "Arcano", ActionType: "Reação", Range: "Pessoal",
			Effect: "Anula e rouba temporariamente uma magia lançada contra você.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Trapaceiro Arcano",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Ladino 5e: características seedadas")
}

// ── Mago ──────────────────────────────────────────────────────────────────────

func seedMago5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Mago")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Conjuração de Mago",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Inteligência. Nível 1: 3 truques, 2 espaços de 1° círculo. Magias preparadas = nível de Mago + modificador de Inteligência. Você possui um Grimório com 6 magias de 1° círculo no nível 1.",
			Keywords: "Arcano, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a truques e espaços de magia. Grimório armazena magias conhecidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Adepto de Ritual",
			Edition: "5e", ClassID: &id,
			Description: "Pode conjurar qualquer magia com o marcador Ritual sem gastar espaço de magia, desde que esteja no Grimório e ele esteja em mãos (não precisa estar preparada).",
			Keywords: "Arcano, Ritual", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Conjura magias de Ritual do Grimório sem gastar espaços de magia.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Recuperação Arcana", Edition: "5e", ClassID: &id,
			Description: "Ao completar um Descanso Curto, pode recuperar espaços de magia gastos cuja soma de círculos seja até metade do seu nível de Mago (arred. p/ cima), nenhum de 6º círculo ou mais. 1 uso por Descanso Longo.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recupera parte dos espaços de magia em Descanso Curto.",
			PowerType: domain.PowerDaily, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Acadêmico", Edition: "5e", ClassID: &id,
			Description: "Escolhe uma perícia em que já é proficiente entre Arcanismo, História, Investigação, Medicina, Natureza ou Religião, e ganha Especialização nela.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Especialização em 1 perícia de conhecimento.",
			PowerType: domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Memorizar Magia", Edition: "5e", ClassID: &id,
			Description: "Ao completar um Descanso Curto, pode substituir uma das magias preparadas de 1º círculo ou mais por outra do Grimório do mesmo círculo ou inferior.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Troca magias preparadas em Descanso Curto, não só em Longo.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Maestria de Magias", Edition: "5e", ClassID: &id,
			Description: "Escolhe uma magia de 1º e uma de 2º círculo (tempo de conjuração de 1 ação) do Grimório: ficam sempre preparadas e conjuráveis no círculo mais baixo sem gastar espaço de magia (espaço necessário só para conjurar em círculo maior). Trocáveis em Descanso Longo.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 magias fixas conjuráveis de graça no círculo base.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Assinatura Mágica", Edition: "5e", ClassID: &id,
			Description: "Escolhe 2 magias de 3º círculo do Grimório como assinaturas: sempre preparadas, cada uma conjurável 1x no 3º círculo sem espaço de magia (recarrega em Descanso Curto ou Longo). Círculo maior exige espaço de magia normal.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 magias de 3º círculo grátis, 1x cada por descanso.",
			PowerType: domain.PowerEncounter, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024 — moveu do nível 2 pro 3) ──
		{
			Name: "Abjurador", Edition: "5e", ClassID: &id,
			Description: "Proteja seus Companheiros e Bana Inimigos — o estudo do Abjurador concentra-se em magias de bloqueio, banimento e proteção, eliminando efeitos nocivos e repelindo influências malignas.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "escola_mago",
		},
		{
			Name: "Adivinhador", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Escolas de Magia de Mago, escolhida no nível 3.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "escola_mago",
		},
		{
			Name: "Evocador", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Escolas de Magia de Mago, escolhida no nível 3.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "escola_mago",
		},
		{
			Name: "Ilusionista", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Escolas de Magia de Mago, escolhida no nível 3.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "escola_mago",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/10/14) ──────────────────────
		{
			Name: "Proteção Arcana", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de Abjuração com espaço de magia, cria uma proteção com PV = 2x nível de Mago + Inteligência, que absorve dano até um Descanso Longo. Conjurar outra magia de Abjuração (ou gastar um espaço como Ação Bônus) restaura 2x o círculo do espaço em PV da proteção.",
			Keywords: "Arcano, Abjuração", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Escudo de PV renovável, alimentado por magias de Abjuração.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Abjurador",
		},
		{
			Name: "Versado em Abjuração", Edition: "5e", ClassID: &id,
			Description: "Ganha 2 magias de Abjuração (até 2º círculo) de graça no Grimório; ao acessar um novo círculo de espaços, ganha mais 1 magia de Abjuração de graça.",
			Keywords: "Arcano, Abjuração", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Magias de Abjuração extras gratuitas no Grimório.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Abjurador",
		},
		{
			Name: "Proteção Projetada", Edition: "5e", ClassID: &id,
			Description: "Como Reação, quando uma criatura à vista a até 9m sofre dano, sua Proteção Arcana pode absorvê-lo em vez dela.",
			Keywords: "Arcano, Abjuração", ActionType: "Reação", Range: "9 metros",
			Effect: "Compartilha o escudo de Proteção Arcana com aliados próximos.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Abjurador",
		},
		{
			Name: "Rompe-Magia", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Contramagia e Dissipar Magia preparadas; pode conjurar Dissipar Magia como Ação Bônus e soma o Bônus de Proficiência ao teste. Se a magia falhar em interromper, o espaço não é gasto.",
			Keywords: "Arcano, Abjuração", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Contramagia/Dissipar Magia sempre disponíveis e mais confiáveis.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Abjurador",
		},
		{
			Name: "Resistência à Magia", Edition: "5e", ClassID: &id,
			Description: "Vantagem em salvaguardas contra magias e Resistência a dano de origem mágica.",
			Keywords: "Arcano, Abjuração", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Defesa forte contra magia inimiga.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Abjurador",
		},
		{
			Name: "Prodígio", Edition: "5e", ClassID: &id,
			Description: "Ao completar um Descanso Longo, joga 2d20 e guarda os resultados: pode substituir qualquer Teste de D20 (seu ou de alguém à vista) por uma dessas jogadas, 1x por turno, decidindo antes de rolar. Perde as não usadas no próximo Descanso Longo.",
			Keywords: "Arcano, Adivinhação", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 14 (Prodígio Maior): joga 3d20 em vez de 2.",
			Effect: "Reserva de resultados de d20 pré-rolados para substituir testes.",
			PowerType: domain.PowerDaily, Level: 3, IsClassFeature: true, ChoiceGroup: "Adivinhador",
		},
		{
			Name: "Versado em Adivinhação", Edition: "5e", ClassID: &id,
			Description: "Ganha 2 magias de Adivinhação (até 2º círculo) de graça no Grimório; ao acessar um novo círculo de espaços, ganha mais 1 magia de Adivinhação de graça.",
			Keywords: "Arcano, Adivinhação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Magias de Adivinhação extras gratuitas no Grimório.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Adivinhador",
		},
		{
			Name: "Perito em Adivinhação", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de Adivinhação com espaço de 2º círculo ou mais, recupera um espaço gasto de círculo inferior (até 5º círculo).",
			Keywords: "Arcano, Adivinhação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Magias de Adivinhação praticamente de graça (recuperam espaço).",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Adivinhador",
		},
		{
			Name: "O Terceiro Olho", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, escolhe até Descanso Curto/Longo: Compreensão Superior (lê qualquer idioma), Ver o Invisível (conjura de graça) ou Visão no Escuro 36m. 1 uso por Descanso Curto ou Longo.",
			Keywords: "Arcano, Adivinhação", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "3 benefícios sensoriais à escolha, um por descanso.",
			PowerType: domain.PowerEncounter, Level: 10, IsClassFeature: true, ChoiceGroup: "Adivinhador",
		},
		{
			Name: "Truque Potente", Edition: "5e", ClassID: &id,
			Description: "Truques de dano afetam mesmo quem resiste: ao errar o ataque ou o alvo passar na salvaguarda contra um truque seu, ele ainda sofre metade do dano (sem efeitos adicionais).",
			Keywords: "Arcano, Evocação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Truques nunca erram completamente.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Evocador",
		},
		{
			Name: "Versado em Evocação", Edition: "5e", ClassID: &id,
			Description: "Ganha 2 magias de Evocação (até 2º círculo) de graça no Grimório; ao acessar um novo círculo de espaços, ganha mais 1 magia de Evocação de graça.",
			Keywords: "Arcano, Evocação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Magias de Evocação extras gratuitas no Grimório.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Evocador",
		},
		{
			Name: "Esculpir Magias", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de Evocação em área, escolhe 1 + o círculo da magia de criaturas à vista que automaticamente passam na salvaguarda e não sofrem dano (se seria metade em sucesso normal).",
			Keywords: "Arcano, Evocação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Protege aliados do próprio dano em área.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Evocador",
		},
		{
			Name: "Evocação Potencializada", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de Evocação, soma o mod. de Inteligência a uma jogada de dano dela.",
			Keywords: "Arcano, Evocação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano extra de Inteligência em magias de Evocação.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true, ChoiceGroup: "Evocador",
		},
		{
			Name: "Sobrecarga", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de dano com espaço de 1º-5º círculo, pode causar dano máximo com ela. Sem custo na primeira vez desde o último Descanso Longo; usos seguintes causam 2d12 de dano Necrótico por círculo do espaço (ignora Resistência/Imunidade), aumentando 1d12 a cada uso extra.",
			Keywords: "Arcano, Evocação", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Dano máximo garantido, com custo crescente de autodano.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Evocador",
		},
		{
			Name: "Ilusões Aprimoradas", Edition: "5e", ClassID: &id,
			Description: "Magias de Ilusão dispensam componentes Verbais, e alcance de 3m+ aumenta em 18m. Também aprende (ou troca) Ilusão Menor, que passa a criar som E imagem juntos e pode ser conjurada como Ação Bônus.",
			Keywords: "Arcano, Ilusão", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Ilusões silenciosas, de maior alcance, e Ilusão Menor aprimorada.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Ilusionista",
		},
		{
			Name: "Versado em Ilusão", Edition: "5e", ClassID: &id,
			Description: "Ganha 2 magias de Ilusão (até 2º círculo) de graça no Grimório; ao acessar um novo círculo de espaços, ganha mais 1 magia de Ilusão de graça.",
			Keywords: "Arcano, Ilusão", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Magias de Ilusão extras gratuitas no Grimório.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Ilusionista",
		},
		{
			Name: "Criaturas Espectrais", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Convocar Feérico e Invocar Fera preparadas, podendo conjurá-las como Ilusão (criatura espectral) sem espaço de magia — mas com metade dos PV. 1 uso sem espaço por Descanso Longo.",
			Keywords: "Arcano, Ilusão", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Invocações gratuitas em versão ilusória (menos PV).",
			PowerType: domain.PowerDaily, Level: 6, IsClassFeature: true, ChoiceGroup: "Ilusionista",
		},
		{
			Name: "Autoimagem Ilusória", Edition: "5e", ClassID: &id,
			Description: "Ao ser atingido por um ataque, usa Reação para criar uma duplicata ilusória entre você e o atacante: o ataque erra automaticamente. Recarrega em Descanso Curto/Longo, ou gastando um espaço de 2º círculo+.",
			Keywords: "Arcano, Ilusão", ActionType: "Reação", Range: "Pessoal",
			Effect: "Anula um ataque certeiro criando um duplo ilusório.",
			PowerType: domain.PowerEncounter, Level: 10, IsClassFeature: true, ChoiceGroup: "Ilusionista",
		},
		{
			Name: "Realidade Ilusória", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar uma magia de Ilusão com espaço de magia, escolhe um objeto inanimado não-mágico dentro dela e o torna real por um período, utilizável fisicamente.",
			Keywords: "Arcano, Ilusão", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Torna um objeto ilusório temporariamente real e utilizável.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true, ChoiceGroup: "Ilusionista",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Mago 5e: características seedadas")
}

// ── Monge ─────────────────────────────────────────────────────────────────────

func seedMonge5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Monge")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Artes Marciais",
			Edition: "5e", ClassID: &id,
			Description: "Sua prática de artes marciais lhe dá domínio em combate desarmado e com armas de Monge. Você pode usar Destreza em vez de Força para ataques e jogadas de dano. Dado de dano desarmado = 1d6 (níveis 1–4). Depois de atacar com uma ação, pode realizar um Ataque Desarmado como Ação Bônus.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "DES em ataques; dado 1d6 desarmado; Ataque Desarmado bônus.",
			LevelScaling: "Nível 5: 1d8. Nível 11: 1d10. Nível 17: 1d12.",
			PowerType:    domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Defesa Sem Armadura",
			Edition: "5e", ClassID: &id,
			Description: "Enquanto não usar armadura nem escudo, sua CA = 10 + modificador de Destreza + modificador de Sabedoria.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "CA = 10 + DES + SAB sem armadura e sem escudo.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Maestria em Armas",
			Edition: "5e", ClassID: &id,
			Description: "Você pode usar as propriedades de Maestria de 2 tipos de armas Simples ou Marciais à sua escolha.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Desbloqueia propriedades de Maestria em 2 armas escolhidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todas as tradições) ──
		{
			Name:    "Foco do Monge", Edition: "5e", ClassID: &id,
			Description: "Reserva de Pontos de Foco (recarrega em Descanso Curto ou Longo) alimentando: Defesa Paciente (Ação Bônus Desengajar, ou 1 ponto para Desengajar+Esquivar), Passo do Vento (Ação Bônus Correr, ou 1 ponto para Desengajar+Correr com salto dobrado) e Torrente de Golpes (1 ponto: 2 Ataques Desarmados como Ação Bônus).",
			Keywords: "Marcial", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:       "Recurso central do Monge alimentando 3 técnicas de combate/mobilidade.",
			LevelScaling: "Nível 10 (Foco Aprimorado): Defesa Paciente dá PV Temporários; Passo do Vento pode levar um aliado; Torrente de Golpes vira 3 ataques.",
			PowerType:    domain.PowerEncounter, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Metabolismo Incomum", Edition: "5e", ClassID: &id,
			Description: "Ao rolar Iniciativa, pode restaurar todos os Pontos de Foco gastos e curar (dado de Artes Marciais + nível de Monge) PV. 1 uso por Descanso Longo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Recarrega Foco + cura ao entrar em combate.",
			PowerType: domain.PowerDaily, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Movimento sem Armadura", Edition: "5e", ClassID: &id,
			Description: "Deslocamento +3m sem armadura nem escudo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:       "+3m de Deslocamento sem armadura/escudo.",
			LevelScaling: "Bônus cresce com o nível: +4,5m (6), +6m (10), +7,5m (14), +9m (18).",
			PowerType:    domain.PowerUnlimited, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Defletir Ataques", Edition: "5e", ClassID: &id,
			Description: "Como Reação a um ataque com dano Contundente/Cortante/Perfurante, reduz o dano em 1d10 + Destreza + nível de Monge. Se zerar o dano, pode gastar 1 Ponto de Foco para redirecioná-lo a outra criatura próxima (salvaguarda de Destreza ou sofre 2x o dado de Artes Marciais + Destreza).",
			Keywords: "Marcial", ActionType: "Reação", Range: "Pessoal",
			Effect:       "Reduz e pode redirecionar dano físico recebido.",
			LevelScaling: "Nível 13 (Defletir Energia): funciona contra qualquer tipo de dano, não só físico.",
			PowerType:    domain.PowerUnlimited, Level: 3, IsClassFeature: true,
		},
		{
			Name:    "Queda Lenta", Edition: "5e", ClassID: &id,
			Description: "Como Reação ao cair, reduz o dano de queda em 5x seu nível de Monge.",
			Keywords: "Marcial", ActionType: "Reação", Range: "Pessoal",
			Effect:    "Reduz muito o dano de queda.",
			PowerType: domain.PowerUnlimited, Level: 4, IsClassFeature: true,
		},
		{
			Name:    "Ataque Extra (Monge)", Edition: "5e", ClassID: &id,
			Description: "Você pode atacar duas vezes, em vez de uma, sempre que executar a ação Atacar no seu turno.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 ataques por ação Atacar.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Golpe Atordoante", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao acertar com arma de Monge ou Desarmado, gasta 1 Ponto de Foco: o alvo salva Constituição ou fica Atordoado até seu próximo turno; se passar, tem o Deslocamento reduzido à metade e sofre Vantagem no próximo ataque contra ele.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Atordoa ou debilita o alvo, gastando Foco.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Golpes Potencializados", Edition: "5e", ClassID: &id,
			Description: "Ao causar dano com Ataque Desarmado, pode escolher causar dano Energético em vez do tipo normal.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Dano Desarmado pode virar Energético (ignora várias resistências físicas).",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true,
		},
		{
			Name:    "Evasão (Monge)", Edition: "5e", ClassID: &id,
			Description: "Em salvaguardas de Destreza para metade do dano: sem dano se passar, metade se falhar (não funciona se Incapacitado).",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Anula ou reduz dano em área evitável por Destreza.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true,
		},
		{
			Name:    "Movimento Acrobático", Edition: "5e", ClassID: &id,
			Description: "Sem armadura/escudo, pode se mover por superfícies verticais e líquidos no seu turno sem cair.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Anda em paredes/água sem cair, sem armadura.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Restauro Pessoal", Edition: "5e", ClassID: &id,
			Description: "No final de cada turno, pode remover de si Amedrontado, Enfeitiçado ou Envenenado. Não sofre Exaustão por falta de comida/água.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Auto-cura de condições mentais/veneno a cada turno + imunidade a privação.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true,
		},
		{
			Name:    "Sobrevivente Disciplinado", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência em TODAS as salvaguardas. Ao falhar uma, pode gastar 1 Ponto de Foco para rejogar, usando o novo resultado.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Proficiência total em salvaguardas + rejogada via Foco.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true,
		},
		{
			Name:    "Foco Perfeito", Edition: "5e", ClassID: &id,
			Description: "Ao rolar Iniciativa sem usar Metabolismo Incomum, recupera Pontos de Foco até ter 4, se tiver 3 ou menos.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Garante um mínimo de 4 Pontos de Foco ao entrar em combate.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true,
		},
		{
			Name:    "Defesa Superior", Edition: "5e", ClassID: &id,
			Description: "No início do turno, gasta 3 Pontos de Foco para ter Resistência a todos os tipos de dano (exceto Energético) por 1 minuto ou até ficar Incapacitado.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Resistência quase universal a dano, sob demanda.",
			PowerType: domain.PowerUnlimited, Level: 18, IsClassFeature: true,
		},
		{
			Name:    "Corpo e Mente", Edition: "5e", ClassID: &id,
			Description: "Seus valores de Destreza e Sabedoria aumentam em 4, até um máximo de 25.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+4 Destreza e +4 Sabedoria (máx. 25).",
			PowerType: domain.PowerUnlimited, Level: 20, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		{
			Name: "Combatente da Mão Espalmada", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Tradições Monásticas de Monge, escolhida no nível 3.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "tradicao_monge",
		},
		{
			Name: "Combatente da Misericórdia", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Tradições Monásticas de Monge, escolhida no nível 3.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "tradicao_monge",
		},
		{
			Name: "Combatente das Sombras", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Tradições Monásticas de Monge, escolhida no nível 3.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "tradicao_monge",
		},
		{
			Name: "Combatente dos Elementos", Edition: "5e", ClassID: &id,
			Description: "Uma das 4 Tradições Monásticas de Monge, escolhida no nível 3.",
			Keywords: "Marcial, Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "tradicao_monge",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/6/11/17) ──────────────────────
		{
			Name: "Técnica da Mão Espalmada", Edition: "5e", ClassID: &id,
			Description: "Ao acertar com um ataque da Torrente de Golpes, aplica um efeito à escolha: Derrubar (Destreza ou Caído), Desorientar (sem Ataques de Oportunidade até o próximo turno) ou Empurrar (Força ou empurra 4,5m).",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Efeito de controle extra nos ataques da Torrente de Golpes.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente da Mão Espalmada",
		},
		{
			Name: "Integridade Corporal", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, joga o dado de Artes Marciais e cura esse valor + Sabedoria (mín. 1) PV. Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Marcial", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Autocura sob demanda.",
			PowerType: domain.PowerDaily, Level: 6, IsClassFeature: true, ChoiceGroup: "Combatente da Mão Espalmada",
		},
		{
			Name: "Passo Veloz", Edition: "5e", ClassID: &id,
			Description: "Ao usar uma Ação Bônus diferente de Passo do Vento, também pode usar Passo do Vento logo em seguida.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Combina Passo do Vento com outra Ação Bônus no mesmo turno.",
			PowerType: domain.PowerUnlimited, Level: 11, IsClassFeature: true, ChoiceGroup: "Combatente da Mão Espalmada",
		},
		{
			Name: "Palma Vibrante", Edition: "5e", ClassID: &id,
			Description: "Ao acertar Ataque Desarmado, gasta 4 Pontos de Foco para plantar vibrações latentes (duram dias = nível de Monge); pode detoná-las depois (renunciando a um ataque) forçando salvaguarda de Constituição: 10d12 de dano Energético, metade se passar.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Bomba de dano com gatilho retardado.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Combatente da Mão Espalmada",
		},
		{
			Name: "Implementos de Misericórdia", Edition: "5e", ClassID: &id,
			Description: "Ganha proficiência em Intuição e Medicina, e com o Kit de Herbalismo.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect: "2 perícias + 1 ferramenta.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente da Misericórdia",
		},
		{
			Name: "Mão de Cura", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta 1 Ponto de Foco para tocar e curar (dado de Artes Marciais + Sabedoria) PV. Numa Torrente de Golpes, pode substituir um ataque por isso de graça.",
			Keywords: "Marcial, Cura", ActionType: "Passiva", Range: "Toque",
			LevelScaling: "Nível 6 (Toque de Médico): também remove Atordoado/Cego/Envenenado/Paralisado/Surdo.",
			Effect: "Cura por toque, integrável à Torrente de Golpes.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente da Misericórdia",
		},
		{
			Name: "Mão de Dolo", Edition: "5e", ClassID: &id,
			Description: "1x por turno, ao causar dano com Ataque Desarmado, gasta 1 Ponto de Foco para +dano Necrótico (dado de Artes Marciais + Sabedoria).",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 6 (Toque de Médico): também impõe Envenenado até o próximo turno do alvo.",
			Effect: "Dano Necrótico extra 1x por turno.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente da Misericórdia",
		},
		{
			Name: "Torrente de Cura e Dolo", Edition: "5e", ClassID: &id,
			Description: "Na Torrente de Golpes, pode usar Mão de Cura em cada ataque substituído sem gastar Foco para a cura, e usar Mão de Dolo sem gastar Foco (ainda 1x por turno). Usos = mod. de Sabedoria (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Marcial, Cura", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Mão de Cura/Dolo de graça durante a Torrente de Golpes.",
			PowerType: domain.PowerDaily, Level: 11, IsClassFeature: true, ChoiceGroup: "Combatente da Misericórdia",
		},
		{
			Name: "Mão da Misericórdia Final", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, toca um cadáver morto há até 24h e gasta 5 Pontos de Foco: a criatura revive com 4d10 + Sabedoria PV, sem as condições Atordoado/Cego/Envenenado/Paralisado/Surdo. 1 uso por Descanso Longo.",
			Keywords: "Marcial, Cura", ActionType: "Passiva", Range: "Toque",
			Effect: "Ressurreição limitada via Pontos de Foco.",
			PowerType: domain.PowerDaily, Level: 17, IsClassFeature: true, ChoiceGroup: "Combatente da Misericórdia",
		},
		{
			Name: "Artes das Sombras", Edition: "5e", ClassID: &id,
			Description: "Escuridão: gasta 1 Ponto de Foco para conjurar Escuridão sem componentes, enxergando dentro dela e podendo movê-la 18m por turno. Ilusão Sombria: conhece Ilusão Menor (Sabedoria como atributo). Também ganha Visão no Escuro 18m (ou +18m se já tiver).",
			Keywords: "Arcano, Sombrio", ActionType: "Passiva", Range: "18 metros",
			Effect: "Controle de escuridão + ilusão + visão no escuro.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente das Sombras",
		},
		{
			Name: "Passo da Sombra", Edition: "5e", ClassID: &id,
			Description: "Em Meia-luz ou Escuridão total, como Ação Bônus teleporta-se até 18m para outro ponto sob Meia-luz/Escuridão, ganhando Vantagem no próximo ataque corpo a corpo do turno.",
			Keywords: "Arcano, Sombrio", ActionType: "Ação Bônus", Range: "18 metros",
			LevelScaling: "Nível 11 (Aprimorado): gastando 1 Ponto de Foco, dispensa o requisito de luz e pode atacar Desarmado logo após teleportar.",
			Effect: "Teleporte furtivo com Vantagem de ataque.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Combatente das Sombras",
		},
		{
			Name: "Manto da Sombra", Edition: "5e", ClassID: &id,
			Description: "Em Meia-luz/Escuridão, gasta 3 Pontos de Foco para virar sombra por 1 minuto: Invisível, parcialmente incorpóreo (atravessa espaços ocupados como Terreno Difícil) e usa Torrente de Golpes sem gastar Foco.",
			Keywords: "Arcano, Sombrio", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Forma sombria: invisibilidade + atravessar criaturas + Foco grátis.",
			PowerType: domain.PowerUnlimited, Level: 17, IsClassFeature: true, ChoiceGroup: "Combatente das Sombras",
		},
		{
			Name: "Manipular Elementos", Edition: "5e", ClassID: &id,
			Description: "Conhece a magia Elementalismo, usando Sabedoria como atributo de conjuração.",
			Keywords: "Arcano, Elemental", ActionType: "Passiva", Range: "Pessoal",
			Effect: "1 magia utilitária elemental fixa.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente dos Elementos",
		},
		{
			Name: "Sintonia Elemental", Edition: "5e", ClassID: &id,
			Description: "No início do turno, gasta 1 Ponto de Foco para se imbuir de energia elemental por 10 minutos: Ataques Desarmados podem causar Ácido/Elétrico/Gélido/Ígneo/Trovejante à escolha (com chance de empurrar/puxar o alvo), e o alcance desarmado aumenta em 3m.",
			Keywords: "Arcano, Elemental", ActionType: "Passiva", Range: "Pessoal",
			LevelScaling: "Nível 11 (Passo dos Elementos): ganha Natação e Voo iguais ao Deslocamento enquanto ativa. Nível 17 (Ápice Elemental): +dano extra 1x por turno e Passo do Vento com +6m.",
			Effect: "Buff elemental central da subclasse: dano, alcance e depois mobilidade.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Combatente dos Elementos",
		},
		{
			Name: "Explosão Elemental", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta 2 Pontos de Foco para detonar uma Esfera de 6m de raio a até 36m com um tipo de dano elemental à escolha: salvaguarda de Destreza, 3x o dado de Artes Marciais de dano (metade se passar).",
			Keywords: "Arcano, Elemental", ActionType: "Passiva", Range: "36 metros",
			Effect: "Dano elemental em área, alimentado por Pontos de Foco.",
			PowerType: domain.PowerUnlimited, Level: 6, IsClassFeature: true, ChoiceGroup: "Combatente dos Elementos",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Monge 5e: características seedadas")
}

// ── Paladino ──────────────────────────────────────────────────────────────────

func seedPaladino5e(db *gorm.DB) {
	id, ok := getClass5e(db, "Paladino")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name:    "Maestria em Armas",
			Edition: "5e", ClassID: &id,
			Description: "Você pode usar as propriedades de Maestria de 2 tipos de armas Simples ou Marciais à sua escolha.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Desbloqueia propriedades de Maestria em 2 armas escolhidas.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Mãos Curadoras",
			Edition: "5e", ClassID: &id,
			Description: "Como Ação, você toca uma criatura e restaura PV de uma reserva de cura igual a 5 × seu nível de Paladino. Você pode distribuir essa cura entre múltiplas ações durante o dia. Recupera em Descanso Longo.",
			Keywords: "Divino, Cura", ActionType: "Ação", Range: "Toque",
			Effect:       "Reserva de cura = Nível × 5 PV. Recupera em Descanso Longo.",
			LevelScaling: "Nível 1: 5 PV. Nível 5: 25 PV. Nível 10: 50 PV.",
			PowerType:    domain.PowerDaily, Level: 1, IsClassFeature: true,
		},
		{
			Name:    "Conjuração de Paladino",
			Edition: "5e", ClassID: &id,
			Description: "Atributo de conjuração: Carisma. Nível 1: sem truques, 2 magias preparadas, 2 espaços de 1° círculo. Magias preparadas = metade do nível de Paladino (arredondado para cima) + modificador de Carisma.",
			Keywords: "Divino, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Acesso a espaços de magia de Paladino. Sem truques.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
		},
		// ── PROGRESSÃO DE NÍVEL (características base, todos os juramentos) ──
		{
			Name:    "Destruição do Paladino", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Destruição Divina preparada, conjurável 1x sem espaço de magia. Recarrega em Descanso Longo.",
			Keywords: "Divino, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Destruição Divina gratuita 1x por Descanso Longo.",
			PowerType: domain.PowerDaily, Level: 2, IsClassFeature: true,
		},
		{
			Name:    "Estilo de Luta (Paladino)", Edition: "5e", ClassID: &id,
			Description: "Ganha um talento de Estilo de Luta à escolha (ou a opção Combatente Abençoado: 2 truques de Clérigo, usando Carisma como atributo de conjuração).",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Bônus passivo de combate ou 2 truques divinos, à escolha.",
			PowerType:      domain.PowerUnlimited, Level: 2, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "estilo_luta_paladino",
		},
		{
			Name:    "Canalizar Divindade (Paladino)", Edition: "5e", ClassID: &id,
			Description: "2 usos (3 a partir do nível 11; recupera 1 em Descanso Curto, todos em Longo) para ativar Sentido Divino: Ação Bônus revelando Celestiais/Ínferos/Mortos-Vivos e locais consagrados/profanados a até 18m por 10 minutos.",
			Keywords: "Divino", ActionType: "Ação Bônus", Range: "18 metros",
			Effect:    "Detecção sobrenatural via Canalizar Divindade.",
			PowerType: domain.PowerEncounter, Level: 3, IsClassFeature: true,
		},
		{
			Name:    "Ataque Extra (Paladino)", Edition: "5e", ClassID: &id,
			Description: "Você pode atacar duas vezes, em vez de uma, sempre que executar a ação Atacar no seu turno.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "2 ataques por ação Atacar.",
			PowerType: domain.PowerUnlimited, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Montaria Fiel", Edition: "5e", ClassID: &id,
			Description: "Sempre tem Convocar Montaria preparada, conjurável 1x sem espaço de magia. Recarrega em Descanso Longo.",
			Keywords: "Divino, Magia", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Convocar Montaria gratuita 1x por Descanso Longo.",
			PowerType: domain.PowerDaily, Level: 5, IsClassFeature: true,
		},
		{
			Name:    "Aura de Proteção", Edition: "5e", ClassID: &id,
			Description: "Emanação de 3m (9m a partir do nível 18): você e aliados na área ganham bônus em salvaguardas igual ao mod. de Carisma (mín. +1). Inativa se você ficar Incapacitado.",
			Keywords: "Divino", ActionType: "Passiva", Range: "3 metros",
			LevelScaling: "Nível 18 (Aura Expandida): raio aumenta para 9 metros.",
			Effect:       "Bônus de salvaguarda em área para o grupo.",
			PowerType:    domain.PowerUnlimited, Level: 6, IsClassFeature: true,
		},
		{
			Name:    "Repudiar Inimigos", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta um uso de Canalizar Divindade: escolhe criaturas à vista a até 18m (número = mod. Carisma, mín. 1) que salvam Sabedoria ou ficam Amedrontadas por 1 minuto (ou até sofrer dano), restritas a mover-se OU agir OU usar Ação Bônus por turno.",
			Keywords: "Divino", ActionType: "Passiva", Range: "18 metros",
			Effect:    "Amedronta e restringe várias criaturas em área.",
			PowerType: domain.PowerUnlimited, Level: 9, IsClassFeature: true,
		},
		{
			Name:    "Aura de Coragem", Edition: "5e", ClassID: &id,
			Description: "Você e aliados na Aura de Proteção têm Imunidade a Amedrontado (a condição fica suprimida enquanto permanecerem na aura).",
			Keywords: "Divino", ActionType: "Passiva", Range: "3 metros",
			Effect:    "Imunidade a medo em área.",
			PowerType: domain.PowerUnlimited, Level: 10, IsClassFeature: true,
		},
		{
			Name:    "Golpes Radiantes", Edition: "5e", ClassID: &id,
			Description: "Ao acertar com arma Corpo a Corpo ou Ataque Desarmado, causa 1d8 de dano Radiante adicional.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+1d8 Radiante em todo ataque corpo a corpo.",
			PowerType: domain.PowerUnlimited, Level: 11, IsClassFeature: true,
		},
		{
			Name:    "Toque Restaurador", Edition: "5e", ClassID: &id,
			Description: "Ao usar Mãos Consagradas, pode gastar 5 PV da reserva por condição para remover Amedrontado, Atordoado, Cego, Enfeitiçado, Paralisado ou Surdo (não restaura PV da criatura).",
			Keywords: "Divino, Cura", ActionType: "Passiva", Range: "Toque",
			Effect:    "Mãos Consagradas passa a remover condições, além de curar.",
			PowerType: domain.PowerUnlimited, Level: 14, IsClassFeature: true,
		},
		// ── SUBCLASSE (nível 3, PHB 2024) ───────────────────────────
		{
			Name: "Juramento da Devoção", Edition: "5e", ClassID: &id,
			Description: "Defenda os Ideais da Justiça — um dos 4 Juramentos Sagrados de Paladino, firmado no nível 3 como o ápice do treinamento inicial.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "juramento_paladino",
		},
		{
			Name: "Juramento da Glória", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Juramentos Sagrados de Paladino, firmado no nível 3.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "juramento_paladino",
		},
		{
			Name: "Juramento de Vingança", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Juramentos Sagrados de Paladino, firmado no nível 3.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "juramento_paladino",
		},
		{
			Name: "Juramento dos Anciões", Edition: "5e", ClassID: &id,
			Description: "Um dos 4 Juramentos Sagrados de Paladino, firmado no nível 3.",
			Keywords: "Divino, Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Concede características exclusivas a partir do nível 3.",
			PowerType: domain.PowerUnlimited, Level: 3,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "juramento_paladino",
		},
		// ── PROGRESSÃO DE SUBCLASSE (níveis 3/7/15/20) ──────────────────────
		{
			Name: "Arma Sagrada", Edition: "5e", ClassID: &id,
			Description: "Ao Atacar, gasta um uso de Canalizar Divindade para imbuir uma arma corpo a corpo por 10 min (ou até reusar): soma Carisma (mín. +1) nos ataques com ela, dano pode virar Radiante, e ela emite Luz Plena 6m + Meia-luz 6m.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Arma imbuída com bônus de ataque, dano radiante e luz.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Juramento da Devoção",
		},
		{
			Name: "Aura de Devoção", Edition: "5e", ClassID: &id,
			Description: "Você e aliados na Aura de Proteção têm Imunidade a Enfeitiçado (suprimida enquanto permanecerem na aura).",
			Keywords: "Divino", ActionType: "Passiva", Range: "3 metros",
			Effect: "Imunidade a encantamento em área.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Juramento da Devoção",
		},
		{
			Name: "Destruição Protetora", Edition: "5e", ClassID: &id,
			Description: "Ao conjurar Destruição Divina, você e aliados na Aura de Proteção ganham Cobertura Parcial até o início do seu próximo turno.",
			Keywords: "Divino", ActionType: "Passiva", Range: "3 metros",
			Effect: "Cobertura em área ao usar a magia de assinatura.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Juramento da Devoção",
		},
		{
			Name: "Resplendor Sagrado", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, ativa por 10 minutos: inimigos que começam o turno na sua Aura sofrem dano Radiante (Carisma + Bônus de Proficiência), a aura vira Luz Plena solar, e você tem Vantagem em salvaguardas forçadas por Ínferos/Mortos-Vivos. 1 uso por Descanso Longo, ou gaste um espaço de 5º círculo.",
			Keywords: "Divino", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Aura ofensiva solar de alto nível.",
			PowerType: domain.PowerDaily, Level: 20, IsClassFeature: true, ChoiceGroup: "Juramento da Devoção",
		},
		{
			Name: "Atleta Inigualável", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, gasta um uso de Canalizar Divindade: por 1 hora, Vantagem em Força (Atletismo) e Destreza (Acrobacia), e +3m em Saltos.",
			Keywords: "Divino", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Buff atlético duradouro via Canalizar Divindade.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Juramento da Glória",
		},
		{
			Name: "Destruição Inspiradora", Edition: "5e", ClassID: &id,
			Description: "Após conjurar Destruição Divina, gasta um uso de Canalizar Divindade para distribuir 2d8 + nível de Paladino em PV Temporários entre criaturas à escolha a até 9m (incluindo você).",
			Keywords: "Divino", ActionType: "Passiva", Range: "9 metros",
			Effect: "PV Temporários em grupo combinados com a magia de assinatura.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Juramento da Glória",
		},
		{
			Name: "Aura de Vivacidade", Edition: "5e", ClassID: &id,
			Description: "Deslocamento +3m. Aliados que entram na Aura de Proteção (ou começam o turno nela) ganham +3m de Deslocamento até o fim do próximo turno deles.",
			Keywords: "Divino", ActionType: "Passiva", Range: "3 metros",
			Effect: "Mobilidade extra para si e para o grupo.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Juramento da Glória",
		},
		{
			Name: "Defesa Gloriosa", Edition: "5e", ClassID: &id,
			Description: "Como Reação a um ataque contra você ou alguém a até 3m, soma o mod. de Carisma (mín. +1) à CA do alvo contra aquele ataque; se errar, pode contra-atacar o agressor. Usos = mod. de Carisma (mín. 1), recarrega em Descanso Longo.",
			Keywords: "Divino", ActionType: "Reação", Range: "3 metros",
			Effect: "Defesa reativa + contra-ataque para si ou aliado próximo.",
			PowerType: domain.PowerDaily, Level: 15, IsClassFeature: true, ChoiceGroup: "Juramento da Glória",
		},
		{
			Name: "Lenda Viva", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, ativa por 10 minutos: Vantagem em todos os testes de Carisma, 1x por turno transforma um ataque errado em acerto, e pode rejogar uma salvaguarda falha como Reação. 1 uso por Descanso Longo, ou gaste um espaço de 5º círculo.",
			Keywords: "Divino", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Pacote de bônus heroicos: Carisma, acerto garantido e rejogada.",
			PowerType: domain.PowerDaily, Level: 20, IsClassFeature: true, ChoiceGroup: "Juramento da Glória",
		},
		{
			Name: "Voto de Inimizade", Edition: "5e", ClassID: &id,
			Description: "Ao Atacar, gasta um uso de Canalizar Divindade para jurar inimizade contra um alvo à vista a até 9m: Vantagem em ataques contra ele por 1 minuto ou até reusar. Se o alvo cair a 0 PV antes disso, pode transferir o voto para outro alvo a até 9m.",
			Keywords: "Divino", ActionType: "Passiva", Range: "9 metros",
			Effect: "Vantagem de ataque contra um alvo marcado, transferível.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Juramento de Vingança",
		},
		{
			Name: "Vingador Implacável", Edition: "5e", ClassID: &id,
			Description: "Ao acertar um Ataque de Oportunidade, pode reduzir o Deslocamento do alvo a 0 até o fim do turno e se mover metade do seu Deslocamento como parte da mesma Reação, sem provocar Oportunidade.",
			Keywords: "Divino", ActionType: "Reação", Range: "Pessoal",
			Effect: "Imobiliza e persegue quem tenta fugir.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Juramento de Vingança",
		},
		{
			Name: "Alma Vingativa", Edition: "5e", ClassID: &id,
			Description: "Quando o alvo do seu Voto de Inimizade acerta ou erra um ataque, pode usar Reação para atacá-lo corpo a corpo se estiver ao alcance.",
			Keywords: "Divino", ActionType: "Reação", Range: "Pessoal",
			Effect: "Contra-ataque automático contra o alvo do Voto.",
			PowerType: domain.PowerUnlimited, Level: 15, IsClassFeature: true, ChoiceGroup: "Juramento de Vingança",
		},
		{
			Name: "Anjo Vingador", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, ativa por 10 minutos: inimigos que começam o turno na Aura salvam Sabedoria ou ficam Amedrontados (com Vantagem de ataque contra eles), e você ganha asas espectrais com Voo 18m (paira). 1 uso por Descanso Longo, ou gaste um espaço de 5º círculo.",
			Keywords: "Divino", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Forma angelical: medo em área + voo de combate.",
			PowerType: domain.PowerDaily, Level: 20, IsClassFeature: true, ChoiceGroup: "Juramento de Vingança",
		},
		{
			Name: "A Ira da Natureza", Edition: "5e", ClassID: &id,
			Description: "Como ação Usar Magia, gasta um uso de Canalizar Divindade: criaturas à escolha à vista a até 4,5m salvam Força ou ficam Contidas por 1 minuto (repetem a salvaguarda a cada turno).",
			Keywords: "Divino, Primitivo", ActionType: "Passiva", Range: "4,5 metros",
			Effect: "Imobiliza várias criaturas próximas com videiras espectrais.",
			PowerType: domain.PowerUnlimited, Level: 3, IsClassFeature: true, ChoiceGroup: "Juramento dos Anciões",
		},
		{
			Name: "Aura de Resistência", Edition: "5e", ClassID: &id,
			Description: "Você e aliados na Aura de Proteção têm Resistência a dano Necrótico, Psíquico e Radiante.",
			Keywords: "Divino, Primitivo", ActionType: "Passiva", Range: "3 metros",
			Effect: "Resistência a 3 tipos de dano em área.",
			PowerType: domain.PowerUnlimited, Level: 7, IsClassFeature: true, ChoiceGroup: "Juramento dos Anciões",
		},
		{
			Name: "Sentinela Imortal", Edition: "5e", ClassID: &id,
			Description: "Ao cair a 0 PV sem morrer, fica com 1 PV e recupera 3x seu nível de Paladino em PV. 1 uso por Descanso Longo. Também para de envelhecer magicamente.",
			Keywords: "Divino, Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect: "Evita a morte uma vez por descanso + imunidade a envelhecimento.",
			PowerType: domain.PowerDaily, Level: 15, IsClassFeature: true, ChoiceGroup: "Juramento dos Anciões",
		},
		{
			Name: "Campeão Ancestral", Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, ativa por 1 minuto: inimigos na Aura de Proteção têm Desvantagem em salvaguardas contra suas magias/Canalizar Divindade (Aliviar Desafio); magias de 1 ação podem ser conjuradas como Ação Bônus (Magias Ágeis); recupera 10 PV no início de cada turno (Regeneração). 1 uso por Descanso Longo, ou gaste um espaço de 5º círculo.",
			Keywords: "Divino, Primitivo", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect: "Ativação final de aura: Desvantagem inimiga + conjuração rápida + regeneração.",
			PowerType: domain.PowerDaily, Level: 20, IsClassFeature: true, ChoiceGroup: "Juramento dos Anciões",
		},
	}
	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Println("  ✓ Paladino 5e: características seedadas")
}
func seedBardoSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Bardo", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Bardo 4e não encontrado"); return
	}
	id := cls.ID

	skills := []domain.Skill{
		// ── CARACTERÍSTICAS AUTOMÁTICAS ─────────────────────────────
		{
			Name: "Palavras de Amizade", Edition: "4e", ClassID: &id,
			Description: "Suas palavras estão repletas de poder arcano e transformam até um simples discurso em uma oratória convincente.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O bardo recebe +5 de bônus de poder no próximo teste de Diplomacia até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1, IsClassFeature: true,
		},
		{
			Name: "Palavra Majestosa", Edition: "4e", ClassID: &id,
			Description: "Ao proferir palavras de inspiração sobrenatural, você restaura o vigor de um aliado.",
			Keywords: "Arcano, Cura", ActionType: "Ação Mínima", Range: "Explosão contígua 5 (10 no nível 11, 15 no nível 21)",
			Target: "O bardo ou um aliado dentro da explosão",
			Effect: "O alvo pode gastar um pulso de cura e recupera PV adicionais iguais ao mod Carisma do bardo. O alvo também é conduzido 1 quadrado.",
			Special: "Pode ser usado duas vezes por encontro, uma por rodada. No nível 16: três vezes por encontro.",
			LevelScaling: "Nível 6: +1d6. Nível 11: +2d6. Nível 16: +3d6. Nível 21: +4d6. Nível 26: +5d6.",
			PowerType: domain.PowerEncounter, Level: 1, IsClassFeature: true,
		},
		// ── VIRTUDE (escolha) ────────────────────────────────────────
		{
			Name: "Virtude da Astúcia", Edition: "4e", ClassID: &id,
			Description: "Quando o ataque de um inimigo fracassa contra um aliado, você pode conduzi-lo 1 quadrado.",
			Keywords: "Arcano",
			Effect: "Uma vez por rodada: quando o ataque de um inimigo fracassa contra um aliado a até 5 + mod INT quadrados, use ação livre para conduzir o aliado 1 quadrado.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "virtude_bardo",
		},
		{
			Name: "Virtude da Bravura", Edition: "4e", ClassID: &id,
			Description: "Quando um aliado próximo reduz um inimigo a 0 PV, você concede PV temporários a esse aliado.",
			Keywords: "Arcano",
			Effect: "Uma vez por rodada: quando aliado a até 5 quadrados reduz inimigo a 0 PV ou deixa sangrando, conceda PV temporários iguais a 1 + mod CON. (Nível 11: 3 + mod CON.)",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "virtude_bardo",
		},
		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Condutor", Edition: "4e", ClassID: &id,
			Description: "O golpe de sua arma guia seus aliados, mostrando onde devem focalizar seus ataques.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "1[A] + mod Carisma de dano. O alvo sofre -2 em uma defesa à escolha do bardo até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma de dano.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe da Canção de Guerra", Edition: "4e", ClassID: &id,
			Description: "Com sua canção de vitória e de guerra, seus aliados se sentem revigorados a cada ataque.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "1[A] + mod Carisma de dano. Os aliados que atingirem o alvo até o final do próximo turno recebem PV temporários iguais ao mod CON.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma de dano.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Marca Indireta", Edition: "4e", ClassID: &id,
			Description: "Ao ocultar seu ataque arcano, você engana o adversário fazendo-o crer que veio de um aliado.",
			Keywords: "Arcano, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Reflexos",
			Hit: "1d8 + mod Carisma de dano. O alvo fica marcado por um aliado do bardo a até 5 quadrados até o final do próximo turno.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma de dano.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Zombaria Malévola", Edition: "4e", ClassID: &id,
			Description: "Emitindo impropérios contra seu adversário, você o coloca num estado de fúria cega.",
			Keywords: "Arcano, Encanto, Implemento, Psíquico", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "1d6 + mod Carisma de dano psíquico. O alvo sofre -2 nas jogadas de ataque até o final do próximo turno.",
			LevelScaling: "Nível 21: 2d6 + mod Carisma de dano psíquico.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Grito do Triunfo", Edition: "4e", ClassID: &id,
			Description: "Você emite um poderoso chamado ao combate, espalhando inimigos e instigando aliados.",
			Keywords: "Arcano, Implemento, Trovejante", ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target: "Os inimigos dentro da rajada", Attack: "Carisma vs. Fortitude",
			Hit: "1d6 + mod Carisma de dano trovejante. O alvo é empurrado 1 quadrado.",
			Effect: "Os aliados dentro da rajada são conduzidos 1 quadrado.",
			Special: "Virtude da Bravura: empurra e conduz quadrados iguais ao mod CON.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Refrão Inspirador", Edition: "4e", ClassID: &id,
			Description: "Sua arma entoa uma canção arcana que guia seus aliados para a vitória.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "2[A] + mod Carisma de dano. Os aliados a até 5 quadrados recebem +1 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Amigos Rápidos", Edition: "4e", ClassID: &id,
			Description: "Entoando uma melodia de falsa amizade, você leva o alvo a um estado de devaneio.",
			Keywords: "Arcano, Encanto, Implemento", ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "O alvo não pode atacar a criatura escolhida pelo bardo até o final do próximo turno ou até que o bardo ou aliado ataque o alvo.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Gafe", Edition: "4e", ClassID: &id,
			Description: "Suas palavras certeiras confundem e perturbam o inimigo, deixando-o vulnerável.",
			Keywords: "Arcano, Encanto, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "O alvo fica atordoado até o final do próximo turno do bardo.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Ecos do Guardião", Edition: "4e", ClassID: &id,
			Description: "Sons arcanos envolvem o alvo, dificultando seus movimentos e ataques.",
			Keywords: "Arcano, Encanto, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "O alvo fica imobilizado e sofre -2 nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Grito Inspirador", Edition: "4e", ClassID: &id,
			Description: "Seu brado furioso apunhala a mente do inimigo, negando-lhe vigor a cada golpe aliado.",
			Keywords: "Arcano, Cura, Implemento, Psíquico", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "3d6 + mod Carisma de dano psíquico. Até o final do encontro, sempre que um aliado atingir o alvo, ele pode gastar um pulso de cura.",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Canção do Matador", Edition: "4e", ClassID: &id,
			Description: "Sua canção de batalha enfraquece as defesas do inimigo a cada golpe.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "2[A] + mod Carisma de dano. O alvo concede vantagem de combate ao bardo e aliados (TR encerra).",
			Miss: "Metade do dano.",
			Effect: "Até o final do encontro, sempre que o bardo atingir um inimigo, este concede vantagem de combate.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Verso do Triunfo", Edition: "4e", ClassID: &id,
			Description: "Com palavras inspiradoras, você incentiva seus aliados ao ataque.",
			Keywords: "Arcano, Implemento", ActionType: "Ação Padrão", Range: "Rajada contígua 5",
			Target: "Os inimigos na rajada", Attack: "Carisma vs. Vontade",
			Hit: "2d6 + mod Carisma de dano psíquico.",
			Miss: "Metade do dano.",
			Effect: "Os aliados dentro da rajada recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Melodia do Caçador", Edition: "4e", ClassID: &id,
			Description: "Ao moldar uma corrente sonora, você cria uma área de silêncio absoluto.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "Rajada contígua 5",
			Effect: "Cria zona de silêncio. Criaturas dentro não podem usar poderes Trovejantes ou Sônicos. Sustentação Mínima: a zona persiste.",
			PowerType: domain.PowerDaily, Level: 2,
		},
		{
			Name: "Canção da Coragem", Edition: "4e", ClassID: &id,
			Description: "Uma canção de encorajamento fortalece aliados para enfrentar os desafios.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O bardo e aliados a até 5 quadrados recebem +1 de bônus de poder nas jogadas de salvamento até o final do encontro.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Canção da Defesa", Edition: "4e", ClassID: &id,
			Description: "Uma melodia defensiva protege você e seus aliados de ataques iminentes.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O bardo e aliados a até 5 quadrados recebem +1 de bônus de poder na CA até o início do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Inspirar Competência", Edition: "4e", ClassID: &id,
			Description: "Sua música inspira um aliado a realizar tarefas além de suas capacidades normais.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O alvo recebe +4 de bônus de poder no próximo teste de perícia até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		// ── NÍVEL 3 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Energia Impulsora", Edition: "4e", ClassID: &id,
			Description: "A magia cria um ataque de energia que afasta um inimigo de um aliado.",
			Keywords: "Arcano, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Fortitude",
			Hit: "1d8 + mod Carisma de dano. O alvo é empurrado 2 quadrados.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma de dano.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},
		{
			Name: "Estrofe Dissonante", Edition: "4e", ClassID: &id,
			Description: "Sons discordantes perturbam o inimigo, reduzindo sua eficácia em combate.",
			Keywords: "Arcano, Implemento, Sônico", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Fortitude",
			Hit: "1d6 + mod Carisma de dano sônico. O alvo sofre -2 em todas as defesas até o final do próximo turno.",
			LevelScaling: "Nível 21: 2d6 + mod Carisma de dano sônico.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},
		{
			Name: "Ferocidade Astuta", Edition: "4e", ClassID: &id,
			Description: "Combinando astúcia e ferocidade, você lança um ataque desconcertante.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "1[A] + mod Carisma de dano. Se o alvo tentar atacar você antes do seu próximo turno, sofre -2 nas jogadas de ataque.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma de dano.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},
		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Chamado aos Cavalos de Batalha", Edition: "4e", ClassID: &id,
			Description: "Seu chamado concede velocidade e poder de ataque renovados a seus aliados.",
			Keywords: "Arcano, Implemento", ActionType: "Ação Padrão", Range: "Rajada contígua 5",
			Target: "Os aliados dentro da rajada",
			Effect: "Os alvos podem se deslocar sua velocidade como ação livre e recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 3,
		},
		// ── NÍVEL 5 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Palavra de Proteção Mística", Edition: "4e", ClassID: &id,
			Description: "Uma palavra de poder cria uma barreira mística protetora ao redor de um aliado.",
			Keywords: "Arcano, Implemento", ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O aliado recebe resistência 5 a todos os danos até o final do próximo turno.",
			LevelScaling: "Nível 21: resistência 10.",
			PowerType: domain.PowerUnlimited, Level: 5,
		},
		{
			Name: "Sátira da Bravura", Edition: "4e", ClassID: &id,
			Description: "Sua sátira desmoraliza os inimigos enquanto inspira bravura em seus aliados.",
			Keywords: "Arcano, Encanto, Implemento, Psíquico", ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target: "Os inimigos dentro da rajada", Attack: "Carisma vs. Vontade",
			Hit: "1d8 + mod Carisma de dano psíquico. O alvo sofre -2 nas jogadas de ataque até o final do próximo turno.",
			Effect: "Os aliados dentro da rajada recebem +1 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma de dano psíquico.",
			PowerType: domain.PowerUnlimited, Level: 5,
		},
		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Melodia do Gelo e do Vento", Edition: "4e", ClassID: &id,
			Description: "Uma melodia gélida lança ventos cortantes contra inimigos, congelando-os.",
			Keywords: "Arcano, Implemento, Congelante", ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão", Attack: "Carisma vs. Fortitude",
			Hit: "2d6 + mod Carisma de dano congelante. O alvo fica lento até o final do próximo turno.",
			Miss: "Metade do dano.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Canção da Discórdia", Edition: "4e", ClassID: &id,
			Description: "Enchendo um inimigo de desconfiança, você o obriga a atacar um aliado.",
			Keywords: "Arcano, Encanto, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "O alvo realiza um ataque básico corpo a corpo contra um aliado à escolha do bardo como ação livre.",
			Effect: "Até o final do encontro, uma vez por rodada quando o alvo se mover, o bardo pode fazê-lo atacar o aliado mais próximo.",
			Miss: "O efeito não persiste.",
			PowerType: domain.PowerDaily, Level: 5,
		},
		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Alegro", Edition: "4e", ClassID: &id,
			Description: "Você cria um ritmo apressado que concede velocidade a você e seus aliados.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "Rajada contígua 5",
			Effect: "O bardo e aliados dentro da rajada recebem +2 de bônus de poder no deslocamento até o final do encontro.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Canção da Conquista", Edition: "4e", ClassID: &id,
			Description: "Uma canção triunfante motiva aliados a avançar com determinação.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "Rajada contígua 5",
			Effect: "Os aliados dentro da rajada recebem +2 de bônus de poder nas jogadas de dano até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Cura do Trapaceiro", Edition: "4e", ClassID: &id,
			Description: "Uma melodia enganosa restaura a vitalidade de um aliado ferido.",
			Keywords: "Arcano, Cura", ActionType: "Ação Mínima", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O aliado pode gastar um pulso de cura e recupera PV adicionais iguais ao mod Carisma do bardo.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Ode ao Sacrifício", Edition: "4e", ClassID: &id,
			Description: "Uma ode ao sacrifício heroico inspira um aliado a agir com abnegação.",
			Keywords: "Arcano", ActionType: "Ação Mínima", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "Até o final do próximo turno, quando o aliado sofrer dano, o bardo sofre metade desse dano em seu lugar.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		// ── NÍVEL 7 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Desviar a Atenção", Edition: "4e", ClassID: &id,
			Description: "Suas palavras desviam a atenção do inimigo, permitindo que aliados se movam livremente.",
			Keywords: "Arcano, Encanto, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "1d8 + mod Carisma de dano psíquico. Os aliados adjacentes ao alvo não provocam ataques de oportunidade até o final do próximo turno.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma de dano psíquico.",
			PowerType: domain.PowerUnlimited, Level: 7,
		},
		{
			Name: "Golpe da Garra do Escorpião", Edition: "4e", ClassID: &id,
			Description: "Sua distração permite que um aliado se desloque ao redor do inimigo.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "1[A] + mod Carisma de dano. Um aliado adjacente ao alvo pode se deslocar 2 quadrados sem provocar ataques de oportunidade.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma de dano.",
			PowerType: domain.PowerUnlimited, Level: 7,
		},
		{
			Name: "Lâmina do Trovão", Edition: "4e", ClassID: &id,
			Description: "Sua arma retumba como um trovão, atingindo um inimigo e perturbando sua concentração.",
			Keywords: "Arcano, Arma, Trovejante", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "1[A] + mod Carisma de dano trovejante. O alvo fica atordoado até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma de dano trovejante.",
			PowerType: domain.PowerUnlimited, Level: 7,
		},
		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Grito de Distração", Edition: "4e", ClassID: &id,
			Description: "Um grito ensurdecedor distrai inimigos, deixando-os vulneráveis.",
			Keywords: "Arcano, Implemento, Trovejante", ActionType: "Ação Padrão", Range: "Rajada contígua 5",
			Target: "Os inimigos dentro da rajada", Attack: "Carisma vs. Fortitude",
			Hit: "3d6 + mod Carisma de dano trovejante. O alvo fica surdo e sofre -2 nas defesas (TR encerra).",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 7,
		},
		{
			Name: "Má-Sorte", Edition: "4e", ClassID: &id,
			Description: "O que parecia uma ode ao destino se torna uma maldição que persegue o inimigo.",
			Keywords: "Arcano, Encanto, Implemento", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "2d8 + mod Carisma de dano psíquico. Até o final do encontro, o alvo joga dois dados em qualquer jogada de ataque e usa o menor resultado.",
			Miss: "Metade do dano. O efeito não persiste.",
			PowerType: domain.PowerDaily, Level: 7,
		},
		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Hino do Resgate Audacioso", Edition: "4e", ClassID: &id,
			Description: "Seu ataque ressoa uma canção que permite a um aliado se mover para posição segura.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "2[A] + mod Carisma de dano. Um aliado a até 5 quadrados pode se deslocar sua velocidade como reação.",
			PowerType: domain.PowerEncounter, Level: 9,
		},
		{
			Name: "Riso Repugnante", Edition: "4e", ClassID: &id,
			Description: "Um terrível ataque de riso convulsivo assola o alvo, incapacitando-o temporariamente.",
			Keywords: "Arcano, Encanto, Implemento, Psíquico", ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Carisma vs. Vontade",
			Hit: "2d8 + mod Carisma de dano psíquico. O alvo fica atordoado até o final do próximo turno do bardo.",
			PowerType: domain.PowerEncounter, Level: 9,
		},
		{
			Name: "Condutor Poderoso", Edition: "4e", ClassID: &id,
			Description: "Como um maestro poderoso, você dirige seus aliados em um ataque coordenado devastador.",
			Keywords: "Arcano, Arma", ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "3[A] + mod Carisma de dano. Os aliados a até 3 quadrados do alvo podem realizar ataques básicos contra ele como reação.",
			PowerType: domain.PowerEncounter, Level: 9,
		},
		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Canção da Recuperação", Edition: "4e", ClassID: &id,
			Description: "Com uma canção inspiradora, aliados recuperam vigor e perseverança.",
			Keywords: "Arcano, Cura", ActionType: "Ação Mínima", Range: "Rajada contígua 5",
			Effect: "Os aliados dentro da rajada podem gastar um pulso de cura.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
		{
			Name: "Palavra da Vida", Edition: "4e", ClassID: &id,
			Description: "Uma palavra de poder restaura a vitalidade de um aliado gravemente ferido.",
			Keywords: "Arcano, Cura", ActionType: "Ação Mínima", Range: "À distância 5",
			Target: "Um aliado inconsciente ou sangrando",
			Effect: "O aliado recupera PV como se tivesse gasto um pulso de cura, mais PV adicionais iguais ao mod Carisma. O aliado também pode ficar de pé como ação livre.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Rasura Ilusória", Edition: "4e", ClassID: &id,
			Description: "Uma ilusão sonora confunde inimigos sobre a posição real de seus aliados.",
			Keywords: "Arcano, Ilusão", ActionType: "Ação Mínima", Range: "Rajada contígua 10",
			Effect: "Os aliados dentro da rajada ficam ocultos para os inimigos até o início do próximo turno do bardo.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
		{
			Name: "Véu", Edition: "4e", ClassID: &id,
			Description: "Um véu de ilusão sonora oculta você e seus aliados dos sentidos inimigos.",
			Keywords: "Arcano, Ilusão", ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O bardo e aliados a até 2 quadrados ficam invisíveis até o início do próximo turno.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Bardo 4e: %d habilidades processadas", len(skills))
}

func seedMongeSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Monge", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Monge 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE (escolha obrigatória — OU) ────
		{
			Name: "Sequência de Golpes Centrada", Edition: "4e", ClassID: &id,
			Description: "Seus punhos ficam ofuscados enquanto você desfere um ataque inicial seguido de outro, ajustando as posições dos inimigos ao seu favor.",
			Keywords:   "Psiônico",
			ActionType: "Nenhuma Ação (Especial)", Range: "Corpo a corpo 1",
			Target: "Uma criatura",
			Effect: "Gatilho: O monge atinge uma criatura durante seu turno. O alvo sofre dano igual a 2 + o modificador de Sabedoria do monge, e o monge conduz o alvo 1 quadrado para 1 quadrado adjacente ao monge ou 1 quadrado em qualquer direção caso o alvo não tenha sido a criatura atingida pelo ataque que ativou o gatilho.",
			Special: "O monge pode usar este poder apenas uma vez por rodada. Nível 11: Uma ou duas criaturas. Nível 21: Cada inimigo adjacente ao monge.",
			LevelScaling: "Nível 21: Cada inimigo adjacente ao monge.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "sequencia_golpes_monge",
		},
		{
			Name: "Sequência de Golpes do Punho de Pedra", Edition: "4e", ClassID: &id,
			Description: "Você golpeia outro inimigo após seu primeiro ataque, um lembre casual da sua grande força.",
			Keywords:   "Psiônico",
			ActionType: "Nenhuma Ação (Especial)", Range: "Corpo a corpo 1",
			Target: "Uma criatura",
			Effect: "Gatilho: O monge atinge uma criatura durante seu turno. O alvo sofre dano igual a 3 + o modificador de Sabedoria do monge. Caso o alvo não tenha sido a criatura atingida pelo ataque que ativou o gatilho, o dano aumenta em 2 (no nível 11 e 6 no nível 21).",
			Special: "O monge pode usar este poder apenas uma vez por rodada. Nível 11: Uma ou duas criaturas. Nível 21: Cada inimigo adjacente ao monge.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "sequencia_golpes_monge",
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Asas do Grou", Edition: "4e", ClassID: &id,
			Description: "Você salta cruzando o campo de batalha e chuta seu adversário, fazendo-o recuar desnorteado.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. Fortitude",
			Hit:          "1d10 + mod Destreza de dano. O monge empurra o alvo 1 quadrado.",
			Special:      "Técnica de Movimento (Ação de Movimento, Pessoal): O monge realiza um teste de Atletismo para saltar com +5 de bônus de poder. A distância do salto não fica limitada pelo deslocamento.",
			LevelScaling: "Nível 21: 2d10 + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Cauda do Dragão", Edition: "4e", ClassID: &id,
			Description: "Sua mão ondula como a cauda de um dragão e com um leve toque libera-se poder que derruba seu adversário ao chão.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. Fortitude",
			Hit:          "1d6 + mod Destreza de dano e o alvo fica derrubado.",
			Special:      "Técnica de Movimento (Ação de Movimento, Corpo a corpo 1): Um aliado ou inimigo derrubado. O monge troca de lugar com o alvo.",
			LevelScaling: "Nível 21: 2d6 + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Cinco Tempestades", Edition: "4e", ClassID: &id,
			Description: "Você se movimenta como um tornado, rodopiando enquanto desfere uma sucessão de chutes e socos que atingem seus adversários.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target:       "Cada inimigo que o monge consiga enxergar dentro da explosão",
			Attack:       "Destreza vs. Reflexos",
			Hit:          "1d8 + mod Destreza de dano.",
			Special:      "Técnica de Movimento (Ação de Movimento, Pessoal): O monge ajusta 2 quadrados.",
			LevelScaling: "Nível 21: 2d8 + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Cobra Dançarina", Edition: "4e", ClassID: &id,
			Description: "Você se esquiva e se movimenta como uma cobra, confundindo seu inimigo e voltando seus ataques contra ele mesmo.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. Reflexos",
			Hit:          "1d10 + mod Destreza de dano. Caso o alvo tenha realizado um ataque de oportunidade contra o monge durante este turno, o alvo sofre dano adicional igual ao modificador de Sabedoria do monge.",
			Special:      "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento +2.",
			LevelScaling: "Nível 21: 2d10 + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Abrir os Portões da Batalha", Edition: "4e", ClassID: &id,
			Description: "Seu movimento repentino pega seu adversário desprevenido e você parte para o ataque.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Reflexos",
			Hit:     "2d10 + mod Destreza de dano. O alvo sofre 1d10 de dano adicional caso ele esteja com todos os pontos de vida no momento em que é atingido pelo monge com este ataque.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento +2. Durante este movimento, o monge não provoca ataques de oportunidade do primeiro inimigo que o monge se afasta.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Despertar da Dor Entorpecida", Edition: "4e", ClassID: &id,
			Description: "Os ferimentos do seu adversário o permitem se esquivar nos ângulos certos. Quando você ataca, você alveja os ferimentos de um único inimigo e encontra o ponto perfeito para atacar.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "2d8 + mod Destreza de dano. Se o alvo estiver sangrando, ele sofre dano adicional neste ataque e no próximo ataque do monge contra ele antes do final do próximo turno. O dano adicional é igual ao modificador de Força do monge.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento. Durante este movimento, os inimigos sangrando não podem atacar o monge com ações de oportunidade ou ação imediatas.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Macaco Bêbado", Edition: "4e", ClassID: &id,
			Description: "Você cambaleia aparentemente fora de controle, seus inimigos ficam espantados por não conseguirem atingir seu corpo desgonçado, após um soco habilidoso, você faz seu adversário atacar um companheiro dele.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Um inimigo",
			Attack:  "Destreza vs. Vontade",
			Hit:     "1d8 + mod Destreza de dano. O monge conduz o alvo 1 quadrado. O alvo então realiza um ataque básico corpo a corpo usando uma ação livre contra um inimigo escolhido pelo monge. O alvo recebe um bônus na jogada de ataque igual ao modificador de Sabedoria do monge.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento +2. Durante este movimento, o monge ignora terreno acidentado e recebe um bônus de poder em todas as defesas contra ataques de oportunidade igual ao seu modificador de Sabedoria.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Tempestade Ascendente", Edition: "4e", ClassID: &id,
			Description: "O ar ao se redor fica carregado de poder enquanto você focaliza sua energia interior em um rugido de trovões.",
			Keywords:   "Disciplina Total, Implemento, Psiônico, Trovejante",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "2d8 + mod Destreza de dano trovejante. Cada inimigo adjacente ao alvo sofre dano trovejante igual ao modificador de Força do monge.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge voa seu deslocamento. Se o monge não pousar no final deste movimento, ele cai.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Chute Giratório do Louva-Deus", Edition: "4e", ClassID: &id,
			Description: "Com passos rápidos e com um impulso assustador, você empurra de lado seus adversários e os aleijam com chutes cruéis.",
			Keywords:   "Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Um, duas ou três criaturas",
			Attack:  "Destreza vs. Fortitude",
			Effect:  "O monge ajusta seu deslocamento. Se o monge entrar num quadrado adjacente a qualquer inimigo durante este ajuste, o monge conduz aquele inimigo 1 quadrado. O monge pode conduzir cada inimigo apenas uma única vez durante o ajuste.",
			Hit:     "2d10 + mod Destreza de dano e o alvo fica lento (TR encerra).",
			Miss:    "Metade do dano e o alvo fica lento até o final do próximo turno do monge.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Espiral Exímia", Edition: "4e", ClassID: &id,
			Description: "Com uma explosão de movimento repentina, você desfere chutes devastadores de energia psiônica nos inimigos próximos.",
			Keywords:   "Energético, Implemento, Postura, Psiônico",
			ActionType: "Ação Padrão", Range: "Explosão contígua 2",
			Target: "Cada inimigo dentro da explosão",
			Attack: "Destreza vs. Reflexos",
			Hit:    "3d8 + mod Destreza de dano energético.",
			Miss:   "Metade do dano.",
			Effect: "O monge pode assumir a postura da espiral. Enquanto durar a postura, o alcance do monge com ataques de toque corpo a corpo aumenta em 1.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Manobra Volteadora do Leopardo", Edition: "4e", ClassID: &id,
			Description: "Mantendo seu equilíbrio perfeito, você percorre um caminho mortal através do embate, desferindo chutes e socos em cada adversário que passar.",
			Keywords:   "Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Corpo a corpo 1",
			Effect:  "O monge ajusta seu deslocamento e pode realizar o seguinte ataque uma vez contra cada inimigo que o monge passar adjacente durante o ajuste.",
			Target:  "Um inimigo",
			Attack:  "Destreza vs. Reflexos",
			Hit:     "2d6 + mod Destreza de dano.",
			Miss:    "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Trovão Harmonioso", Edition: "4e", ClassID: &id,
			Description: "Você soca um adversário, então gira e desfere um chute em outro. Trovões troam à distância, aproximando-se e explodindo entre seus dois adversários.",
			Keywords:   "Implemento, Psiônico, Trovejante",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma ou duas criaturas",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "3d6 + mod Destreza de dano trovejante.",
			Miss:    "Metade do dano.",
			Effect:  "Na primeira vez em que cada um dos alvos sofre dano durante um turno, o outro alvo sofre dano trovejante igual ao modificador de Força do monge. Este efeito dura até o final do encontro ou até um dos alvos cair a 0 pontos de vida.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Agarrar o Vento", Edition: "4e", ClassID: &id,
			Description: "Antes que seu inimigo o obrigue a recuar, você salta girando, usando o poder do ataque do inimigo para lançar-se para onde você queira.",
			Keywords:   "Psiônico",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O monge é conduzido, empurrado ou puxado.",
			Effect:  "Em vez de ser afetado pelo movimento forçado, o monge ajusta o número de quadrados que ele deveria ser movido.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Disciplina Harmoniosa", Edition: "4e", ClassID: &id,
			Description: "Uma sequência específica de respirações disciplinadas aprimoram tanto sua defesa quanto sua ofensiva.",
			Keywords:   "Psiônico",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O monge adquire pontos de vida adicionais igual ao seu modificador de Sabedoria. Quando o monge não possuir nenhum ponto de vida temporário sobrando, ele recebe um bônus na jogada de dano do seu próximo ataque corpo a corpo antes do final do seu próximo turno. O bônus é igual ao modificador de Sabedoria do Monge.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Marcha Cuidadosa", Edition: "4e", ClassID: &id,
			Description: "Você caminha com tamanha precisão e controle que aquele piso quebrado e mesmo corpos de água não conseguem impedi-lo.",
			Keywords:   "Psiônico",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "Até o final deste turno, o monge ignora terreno acidentado e pode pisar em terrenos líquidos sem afundar e permanecer parado acima de líquidos como se fosse um piso sólido. Além disso, o monge move seu deslocamento.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Sequência Suprema", Edition: "4e", ClassID: &id,
			Description: "Seu movimento fica ofuscado. Onde um golpe termina e o outro começa? Isso não importa enquanto os golpes atingirem.",
			Keywords:   "Psiônico",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special: "Gatilho: O monge usa seu poder Sequência de Golpes e resolve os efeitos do poder que ativou o poder.",
			Effect:  "O monge ajusta metade do seu deslocamento e utiliza seu poder Sequência de Golpes novamente.",
			PowerType: domain.PowerDaily, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Campeão Duradouro", Edition: "4e", ClassID: &id,
			Description: "Você concentra sua dor em um ponto na extremidade do seu punho. Ao atacar, você transmite sua dor para o seu inimigo.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "2d10 + mod Destreza de dano. O monge pode realizar um teste de resistência contra um efeito que possa ser encerrado com um sucesso em um teste de resistência, com um bônus igual ao modificador de Sabedoria do monge. Se o monge obter sucesso no teste de resistência, o efeito não apenas se encerra mas o alvo sofre dano igual ao modificador de Sabedoria do monge.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento +2. Sempre que o monge for atacado durante este movimento, o monge recebe +1 de bônus no deslocamento até o final do seu próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Dança das Espadas", Edition: "4e", ClassID: &id,
			Description: "Enquanto seus adversários cercam você, você salta entre eles, voltando seus números contra eles mesmos.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Reflexos",
			Hit:     "2d8 + mod Destreza de dano e o alvo sofre dano adicional igual a duas vezes o número de inimigos adjacentes ao monge.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento +2. Se os inimigos realizarem ataques de oportunidade contra o monge durante este movimento e fracassarem, o monge adquire vantagem de combate contra os inimigos que fracassaram até o final do turno do monge.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Montanha Eterna", Edition: "4e", ClassID: &id,
			Description: "Você concentra sua mente, convocando sua disciplina de ferro para caminhar, combater e resistir aos ataques com um espírito duradouro da montanha.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target:  "Cada inimigo dentro da explosão",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "2d8 + mod Destreza de dano e o monge derruba o alvo.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge adquire resistência contra todos os tipos de dano igual ao seu modificador de Força até o final do seu próximo turno. Além disso, o monge ajusta 2 quadrados.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Trovões Gêmeos", Edition: "4e", ClassID: &id,
			Description: "Você se movimenta ofuscado e desfere um chute giratório com tamanha ferocidade que uma energia trovejante atinge tanto seu adversário quanto os companheiros dele.",
			Keywords:   "Disciplina Total, Implemento, Psiônico, Trovejante",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "2d10 + mod Destreza de dano trovejante e um inimigo adjacente ao alvo sofre 1d10 de dano trovejante.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge move seu deslocamento +2 e escolhe um inimigo adjacente ao monge no início do movimento. Durante este movimento, o monge não provoca ataques de oportunidade do inimigo escolhido.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Cem Folhas", Edition: "4e", ClassID: &id,
			Description: "Você desfere uma sequência de ataques, golpeando com tamanha velocidade que as criaturas se dispersam diante de você como folhas em um furacão.",
			Keywords:   "Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target: "Cada criatura dentro da rajada",
			Attack: "Destreza vs. Reflexos",
			Hit:    "3d8 + mod Destreza de dano e o monge empurra o alvo 2 quadrados.",
			Miss:   "Metade do dano e o monge empurra o alvo 1 quadrado.",
			Effect: "Até o final do próximo turno do monge, o monge pode alvejar uma criatura adicional dentro do alcance com o seu poder Sequência de Golpes.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Mente em Branco", Edition: "4e", ClassID: &id,
			Description: "Você esvazia sua mente de todos os pensamentos, tornando-se difícil de atingir com poderes que afetam a mente.",
			Keywords:   "Psiônico",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do próximo turno do monge, ele recebe resistência 5 a dano psíquico e +2 de bônus de poder nas defesas de Vontade.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Descanso do Guerreiro", Edition: "4e", ClassID: &id,
			Description: "Com uma série de respirações profundas, você restaura sua vitalidade em plena batalha.",
			Keywords:   "Cura, Psiônico",
			ActionType: "Ação Padrão", Range: "Pessoal",
			Effect: "O monge pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Sabedoria.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Voo do Pássaro de Ferro", Edition: "4e", ClassID: &id,
			Description: "Você voa em direção ao adversário com a velocidade e precisão de um pássaro de ferro.",
			Keywords:   "Disciplina Total, Implemento, Psiônico",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Effect:  "O monge voa até sua velocidade antes de realizar o ataque. Esse voo não provoca ataques de oportunidade.",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "3d8 + mod Destreza de dano. O alvo fica imobilizado até o final do próximo turno.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge voa seu deslocamento +2.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Vento Cortante", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma do vento cortante, tornando-se imparável no campo de batalha.",
			Keywords:   "Implemento, Postura, Psiônico",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura do Vento Cortante. Enquanto nessa postura: o monge recebe +2 de bônus de poder na velocidade, pode ignorar terreno difícil durante o movimento e seus ataques de toque corpo a corpo causam dano adicional igual ao modificador de Sabedoria.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Punho da Tempestade Perfeita", Edition: "4e", ClassID: &id,
			Description: "Seu punho concentra toda a energia da tempestade num único golpe devastador.",
			Keywords:   "Disciplina Total, Implemento, Psiônico, Trovejante",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "3d10 + mod Destreza de dano trovejante. O alvo fica atordoado até o final do próximo turno. Todos os inimigos adjacentes ao alvo sofrem 1d10 de dano trovejante.",
			Special: "Técnica de Movimento (Ação de Movimento, Pessoal): O monge ajusta seu deslocamento +2 e recebe +2 de bônus de poder nas defesas até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Corpo de Diamante", Edition: "4e", ClassID: &id,
			Description: "Você fortalece seu corpo até a dureza de um diamante, tornando-se temporariamente imune a danos.",
			Keywords:   "Psiônico",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do próximo turno do monge, ele recebe resistência 10 a todos os danos.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
		{
			Name: "Transcendência Marcial", Edition: "4e", ClassID: &id,
			Description: "Você atinge um estado de perfeição marcial que potencializa todos os seus ataques.",
			Keywords:   "Postura, Psiônico",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura da Transcendência Marcial. Enquanto nessa postura: o monge recebe +2 de bônus de poder em todas as jogadas de ataque e dano, e sua Sequência de Golpes pode ser usada duas vezes por rodada.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Monge 4e: %d habilidades processadas", len(skills))
}

func seedClerigoSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Clérigo", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Clérigo 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE (automáticas) ─────────────────
		{
			Name: "Canalizar Divindade: Força Divina", Edition: "4e", ClassID: &id,
			Description: "Você invoca o poder da sua divindade para fortalecê-lo durante um ataque.",
			Keywords:    "Divino",
			ActionType:  "Ação Livre",
			Range:       "Pessoal",
			Effect:      "Até o final do turno, o clérigo adiciona o modificador de Força a uma jogada de dano.",
			PowerType: domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Canalizar Divindade: Fortuna Divina", Edition: "4e", ClassID: &id,
			Description: "Diante do perigo, você confia na sua fé e recebe uma vantagem especial.",
			Keywords:    "Divino",
			ActionType:  "Ação Livre",
			Range:       "Pessoal",
			Effect:      "O clérigo recebe +1 de bônus na sua próxima jogada de ataque ou teste de resistência antes do final do seu próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Chama Sagrada", Edition: "4e", ClassID: &id,
			Description: "Uma luz sagrada brilha do alto, queimando um inimigo com sua radiância e restaurando um companheiro com seu poder.",
			Keywords: "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura", Attack: "Sabedoria vs. Reflexos",
			Hit: "1d6 + mod Sabedoria de dano radiante. Um aliado na linha de visão escolhe entre receber PV temporários iguais ao mod Carisma + metade do nível ou realizar um teste de resistência.",
			LevelScaling: "Nível 21: 2d6 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Escudo do Sacerdote", Edition: "4e", ClassID: &id,
			Description: "Você profere uma curta oração defensiva enquanto ataca com sua arma.",
			Keywords: "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força de dano. O personagem e um aliado adjacente recebem +1 de bônus de poder na CA até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Lança da Fé", Edition: "4e", ClassID: &id,
			Description: "Um raio de luz brilhante incinera o inimigo com radiação dourada, guiando os ataques de um aliado.",
			Keywords: "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura", Attack: "Sabedoria vs. Reflexos",
			Hit: "1d8 + mod Sabedoria de dano radiante. Um aliado na linha de visão recebe +2 de bônus de poder na próxima jogada de ataque contra o alvo.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Marca da Integridade", Edition: "4e", ClassID: &id,
			Description: "Você golpeia o inimigo inscrevendo nele o símbolo luminescente da fúria da sua divindade, concedendo poder ao aliado escolhido.",
			Keywords: "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força de dano. Um aliado escolhido recebe bônus nas jogadas de ataque contra o alvo igual ao modificador de Força do clérigo até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Causar Medo", Edition: "4e", ClassID: &id,
			Description: "Seu símbolo sagrado irrompe com as chamas da fúria da sua divindade, obrigando o inimigo a recuar imediatamente.",
			Keywords: "Divino, Implemento, Medo",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Sabedoria vs. Vontade",
			Hit: "O alvo se afasta da personagem usando seu próprio deslocamento + mod Carisma em fuga, evitando espaços perigosos quando possível.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe Restaurador", Edition: "4e", ClassID: &id,
			Description: "Uma radiação divina brilha na sua arma. Quando você golpeia um inimigo, sua divindade concede uma pequena bênção em forma de cura a você ou a um aliado.",
			Keywords: "Arma, Cura, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "2[A] + mod Força de dano radiante. O alvo fica marcado até o final do próximo turno. O clérigo ou um aliado a até 5 quadrados pode gastar um pulso de cura.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Luminescência Divina", Edition: "4e", ClassID: &id,
			Description: "Murmurando uma prece à sua divindade, você invoca uma rajada de luminosidade branca com seu símbolo sagrado.",
			Keywords: "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target: "Os inimigos dentro da rajada", Attack: "Sabedoria vs. Reflexos",
			Hit:    "1d8 + mod Sabedoria de dano radiante.",
			Effect: "Os aliados do clérigo dentro da área recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Trovão Colérico", Edition: "4e", ClassID: &id,
			Description: "Seu braço é fortalecido pelo poder da sua divindade. Uma trovoada implacável esmaga e atordoa o oponente.",
			Keywords: "Arma, Divino, Trovejante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força de dano trovejante. O alvo fica pasmo até o final do próximo turno do clérigo.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Cascata de Luz", Edition: "4e", ClassID: &id,
			Description: "Uma explosão de energia divina incinera o oponente com radiação avassaladora.",
			Keywords: "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Sabedoria vs. Vontade",
			Hit:  "3d8 + mod Sabedoria de dano radiante. O alvo adquire vulnerabilidade 5 a todos os ataques do clérigo (TR encerra).",
			Miss: "Metade do dano e o alvo não adquire a vulnerabilidade.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Chama Vingadora", Edition: "4e", ClassID: &id,
			Description: "Você ataca o inimigo com sua arma e ele irrompe em chamas que vingarão qualquer ataque que ele ousar desferir.",
			Keywords: "Arma, Divino, Flamejante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:     "2[A] + mod Força de dano e 5 de dano flamejante contínuo (TR encerra).",
			Miss:    "Metade do dano e nenhum dano contínuo.",
			Special: "Se atacar durante o seu turno, o alvo não pode realizar um teste de resistência contra o dano contínuo.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Farol da Esperança", Edition: "4e", ClassID: &id,
			Description: "Uma explosão de energia divina fere seus oponentes e cura seus aliados durante todo o encontro.",
			Keywords: "Cura, Divino, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão", Attack: "Sabedoria vs. Vontade",
			Hit:    "O alvo fica enfraquecido até o final do próximo turno do clérigo.",
			Effect: "O clérigo e seus aliados dentro da explosão recuperam +5 PV sempre que usarem um pulso de cura pelo resto do encontro.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Auxílio Divino", Edition: "4e", ClassID: &id,
			Description: "Rezando à sua divindade, você pede que ela conceda a força para superar uma dificuldade.",
			Keywords: "Divino",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "O personagem ou um aliado",
			Effect: "O alvo realiza um teste de resistência com um bônus igual ao modificador de Carisma do clérigo.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Bênção", Edition: "4e", ClassID: &id,
			Description: "Suplicando à sua divindade, você pede bênçãos para você e seus aliados.",
			Keywords: "Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 20",
			Effect: "Até o final do encontro, todos os aliados dentro da explosão recebem +1 de bônus de poder nas jogadas de ataque.",
			PowerType: domain.PowerDaily, Level: 2,
		},
		{
			Name: "Curar Ferimentos Leves", Edition: "4e", ClassID: &id,
			Description: "Você recita uma prece curta e adquire o poder de curar ferimentos instantaneamente.",
			Keywords: "Cura, Divino",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target: "O personagem ou outra criatura",
			Effect: "O alvo recupera pontos de vida como se tivesse gasto um pulso de cura.",
			PowerType: domain.PowerDaily, Level: 2,
		},
		{
			Name: "Escudo da Fé", Edition: "4e", ClassID: &id,
			Description: "Um escudo brilhante de energia divina surge no ar, protegendo você e seus aliados mais próximos.",
			Keywords: "Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 5",
			Target: "O personagem e aliados dentro da explosão",
			Effect: "Os alvos recebem +2 de bônus de poder na CA até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 2,
		},
		{
			Name: "Santuário", Edition: "4e", ClassID: &id,
			Description: "Você invoca um campo protetor sobre uma criatura que diminui a eficácia dos ataques dos inimigos.",
			Keywords: "Divino",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "O personagem ou outra criatura",
			Effect: "O alvo recebe +5 de bônus em todas as defesas. O efeito persiste até o alvo atacar ou até o final do próximo turno do clérigo.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Comando", Edition: "4e", ClassID: &id,
			Description: "Você pronuncia uma única palavra de comando ao seu adversário, exigindo obediência.",
			Keywords: "Encanto, Divino, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Sabedoria vs. Vontade",
			Hit: "O alvo fica pasmo até o final do próximo turno do clérigo. Além disso, fica derrubado ou é conduzido um número de quadrados igual a 3 + o modificador de Carisma do clérigo (à escolha).",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Farol Escaldante", Edition: "4e", ClassID: &id,
			Description: "Você invoca o nome da sua divindade e uma luz sagrada envolve sua arma, guiando os ataques dos aliados à distância.",
			Keywords: "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força de dano radiante. Todos os aliados recebem +4 de bônus de poder nas jogadas de ataque à distância contra o alvo até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Arma dos Deuses", Edition: "4e", ClassID: &id,
			Description: "Sua arma brilha com radiação divina, aprimorando seus ataques durante todo o encontro.",
			Keywords: "Arma, Divino, Radiante",
			ActionType: "Ação Mínima", Range: "Toque corpo a corpo",
			Target: "Uma arma na mão do personagem",
			Effect: "Até o final do encontro, os ataques com essa arma causam +1d6 de dano radiante. Sempre que atingir um inimigo, ele sofre -2 de penalidade na CA até o final do próximo turno do atacante.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Fascinante", Edition: "4e", ClassID: &id,
			Description: "A fascinação e o pavor sobrenatural que irradiam conforme você brande sua arma deixam os oponentes congelados de terror.",
			Keywords: "Arma, Divino, Medo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. Vontade",
			Hit: "1[A] + mod Força de dano. O alvo fica imobilizado até o final do próximo turno do clérigo.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Luz Abrasadora", Edition: "4e", ClassID: &id,
			Description: "Clamando pelo poder da sua divindade, seu símbolo emite uma luz tão intensa que queima os inimigos.",
			Keywords: "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão", Attack: "Sabedoria vs. Reflexos",
			Hit: "2d6 + mod Sabedoria de dano radiante.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Coluna de Chamas", Edition: "4e", ClassID: &id,
			Description: "Uma coluna de chamas surge acima do alvo e engolfa seus oponentes com um estrondo.",
			Keywords: "Divino, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão", Attack: "Sabedoria vs. Reflexos",
			Hit:  "2d10 + mod Sabedoria de dano flamejante e 5 + mod Sabedoria de dano flamejante contínuo (TR encerra).",
			Miss: "Metade do dano e nenhum dano contínuo.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Baluarte da Saúde", Edition: "4e", ClassID: &id,
			Description: "Você entoa uma oração que fortalece instantaneamente um aliado ferido.",
			Keywords: "Cura, Divino",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O alvo pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Sabedoria do clérigo.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Defensores Astrais", Edition: "4e", ClassID: &id,
			Description: "Você conjura dois soldados fantasmagóricos com armas luminescentes que investem contra inimigos próximos.",
			Keywords: "Conjuração, Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Effect: "Você conjura dois defensores astrais adjacentes a um quadrado dentro do alcance. Cada defensor tem CA 20, Fortitude 18, Reflexos 14, Vontade 14. Quando um inimigo adjacente a um defensor atacar uma criatura diferente do defensor, o defensor ataca: Sabedoria vs. CA — 1d8 + mod Sabedoria radiante.",
			Special: "Sustentação Mínima: os defensores persistem.",
			PowerType: domain.PowerDaily, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe do Guardião", Edition: "4e", ClassID: &id,
			Description: "Você golpeia seu inimigo com sua arma e ele fica marcado e seus aliados recebem proteção divina.",
			Keywords: "Arma, Cura, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O clérigo e seus aliados adjacentes ao alvo podem gastar um pulso de cura. Adicione o mod Carisma do clérigo nos PV recuperados.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Escudo do Sacerdote Melhorado", Edition: "4e", ClassID: &id,
			Description: "Uma aura de proteção divina envolve você e seus aliados, reduzindo o dano que sofremos.",
			Keywords: "Divino",
			ActionType: "Ação Padrão", Range: "Pessoal",
			Effect: "Até o final do próximo turno do clérigo, o personagem e aliados adjacentes recebem resistência 5 a todos os danos.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Cura Abençoada", Edition: "4e", ClassID: &id,
			Description: "Uma onda de energia curativa emana de você, restaurando a saúde de todos os aliados próximos.",
			Keywords: "Cura, Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 5",
			Target: "Cada aliado dentro da explosão",
			Effect: "Cada aliado pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Sabedoria do clérigo.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Palavra Sagrada", Edition: "4e", ClassID: &id,
			Description: "Você profere uma palavra sagrada que debilita os inimigos e fortalece seus aliados.",
			Keywords: "Divino, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão", Attack: "Sabedoria vs. Vontade",
			Hit:    "Os inimigos ficam atordoados até o final do próximo turno.",
			Effect: "Os aliados dentro da explosão recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Clérigo 4e: %d habilidades processadas", len(skills))
}

func seedGuerreiroSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Guerreiro", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Guerreiro 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Certeiro", Edition: "4e", ClassID: &id,
			Description: "Você substitui potência por precisão, mirando cuidadosamente antes de atacar.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força + 2 vs. CA",
			Hit:          "1[A] de dano.",
			LevelScaling: "Nível 21: 2[A] de dano.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Fulminante", Edition: "4e", ClassID: &id,
			Description: "Você combina ataques precisos com golpes inesperados que atravessam as defesas do inimigo.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano.",
			Miss:         "Metade do mod Força de dano. Com arma de duas mãos, causa todo o modificador de Força ao fracassar.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Maré de Ferro", Edition: "4e", ClassID: &id,
			Description: "Depois de cada golpe vigoroso, você ergue seu escudo e o utiliza para empurrar o inimigo para trás.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Força vs. CA",
			Hit:     "1[A] + mod Força de dano. O alvo é empurrado 1 quadrado se for do mesmo tamanho ou menor que o guerreiro.",
			Special: "Condição: o guerreiro deve empunhar um escudo.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Trespassar", Edition: "4e", ClassID: &id,
			Description: "Você atinge um inimigo e seu golpe trespassa para atingir o próximo.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. Um inimigo adjacente ao guerreiro (exceto o alvo inicial) sofre dano igual ao modificador de Força do guerreiro.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Ataque de Cobertura", Edition: "4e", ClassID: &id,
			Description: "Você desfere uma série de ataques estonteantes, obrigando o inimigo a prestar atenção em você enquanto um aliado recua com segurança.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "2[A] + mod Força de dano. Um aliado do guerreiro adjacente ao alvo pode ajustar 2 quadrados.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Ataque de Passagem", Edition: "4e", ClassID: &id,
			Description: "Você ataca um inimigo e permite que o impulso o carregue adiante para golpear um segundo adversário.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura (primário)", Attack: "Força vs. CA",
			Hit:    "1[A] + mod Força de dano. O guerreiro pode ajustar 1 quadrado. Realize um ataque secundário contra uma criatura diferente: Força + 2 vs. CA — 1[A] + mod Força de dano.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Lâmina da Serpente de Aço", Edition: "4e", ClassID: &id,
			Description: "Você golpeia violentamente os joelhos ou as pernas do adversário para atrasá-lo.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "2[A] + mod Força de dano. O alvo fica lento e não pode ajustar até o final do próximo turno do guerreiro.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Rasteira Giratória", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe longo e poderoso na parte inferior da guarda do adversário e o derruba no solo.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força de dano. O alvo fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Ameaçar o Vilão", Edition: "4e", ClassID: &id,
			Description: "Você golpeia arduamente, intimidando o inimigo com bloqueios técnicos e contra-ataques severos.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:  "2[A] + mod Força de dano. O guerreiro recebe +2 de bônus de poder nas jogadas de ataque e +4 de bônus de poder no dano contra esse alvo até o final do encontro.",
			Miss: "O guerreiro recebe +1 de bônus de poder nas jogadas de ataque e +2 de bônus de poder no dano contra esse alvo até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Golpe Brutal", Edition: "4e", ClassID: &id,
			Description: "Você despedaça a armadura e os ossos do inimigo com um golpe ressoante.",
			Keywords: "Arma, Confiável, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "3[A] + mod Força de dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Golpe de Revinda", Edition: "4e", ClassID: &id,
			Description: "Um golpe certeiro contra um inimigo odiado revigoroa você e concede a força e a determinação para continuar lutando.",
			Keywords: "Arma, Confiável, Cura, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "2[A] + mod Força de dano. O guerreiro pode gastar um pulso de cura.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Comigo, Agora!", Edition: "4e", ClassID: &id,
			Description: "Você puxa um aliado para uma posição mais vantajosa no campo de batalha.",
			Keywords: "Marcial",
			ActionType: "Ação de Movimento", Range: "Corpo a corpo 1",
			Target: "Um aliado voluntário adjacente",
			Effect: "O alvo é conduzido 2 quadrados para um quadrado adjacente ao guerreiro.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Fechar a Guarda", Edition: "4e", ClassID: &id,
			Description: "Você ergue sua arma ou escudo para bloquear uma brecha nas suas defesas quando um inimigo tenta explorá-la.",
			Keywords: "Marcial",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: Um inimigo com vantagem de combate contra o guerreiro o ataca.",
			Effect:  "Cancele a vantagem de combate para esse ataque.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Irrefreável", Edition: "4e", ClassID: &id,
			Description: "Você permite que um surto de adrenalina o controle durante a batalha.",
			Keywords: "Cura, Marcial",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O guerreiro recebe um número de pontos de vida temporários igual a 2d6 + seu modificador de Constituição.",
			PowerType: domain.PowerDaily, Level: 2,
		},

		// ── NÍVEL 3 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Poderoso", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe com toda a sua força, empurrando o inimigo e abrindo espaço para os aliados.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O alvo é empurrado 2 quadrados.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},
		{
			Name: "Golpe Marcante", Edition: "4e", ClassID: &id,
			Description: "Seu golpe marca o inimigo de forma clara, fazendo com que ele sofra consequências se atacar seus aliados.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O alvo fica marcado pelo guerreiro até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Derrubador", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe tão devastador que derruba o inimigo e o mantém no chão.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. Fortitude",
			Hit: "1[A] + mod Força de dano. O alvo fica derrubado e não pode se levantar até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Impulso do Escudo", Edition: "4e", ClassID: &id,
			Description: "Você usa seu escudo como arma, empurrando o inimigo e criando espaço para manobrar.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Força vs. Fortitude",
			Hit:     "1[A] + mod Força de dano. O alvo é empurrado 3 quadrados.",
			Special: "Condição: deve empunhar um escudo.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Postura do Mestre das Armas", Edition: "4e", ClassID: &id,
			Description: "Você assume uma postura de combate avançada que potencializa todos os seus ataques.",
			Keywords: "Marcial, Postura",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura do Mestre das Armas. Enquanto nessa postura: +2 de bônus de poder nas jogadas de ataque com armas corpo a corpo e +2 de bônus de poder no dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Esfacelador", Edition: "4e", ClassID: &id,
			Description: "Você ataca com tamanha força que enfraquece as defesas do inimigo.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O alvo sofre -2 de penalidade na CA até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 5,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Tumultuoso", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe que perturba todos os inimigos ao redor do alvo.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão", Attack: "Força vs. CA",
			Hit: "2[A] + mod Força de dano.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Manter a Linha", Edition: "4e", ClassID: &id,
			Description: "Você estabelece uma posição defensiva formidável que protege seus aliados.",
			Keywords: "Marcial, Postura",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura de Manter a Linha. Enquanto nessa postura: você e aliados adjacentes recebem +2 de bônus de poder na CA.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Desafio do Guerreiro", Edition: "4e", ClassID: &id,
			Description: "Você desafia todos os inimigos próximos, marcando-os e ameaçando retaliar.",
			Keywords: "Marcial",
			ActionType: "Ação Mínima", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Effect: "Todos os inimigos dentro da explosão ficam marcados pelo guerreiro até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Segundo Fôlego", Edition: "4e", ClassID: &id,
			Description: "Você encontra reservas ocultas de energia para continuar lutando.",
			Keywords: "Cura, Marcial",
			ActionType: "Ação Padrão", Range: "Pessoal",
			Effect: "O guerreiro gasta um pulso de cura e recupera PV adicionais iguais a 1d10 + modificador de Constituição.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Vingativo", Edition: "4e", ClassID: &id,
			Description: "Quando você ou um aliado sofre dano, você responde com um ataque poderoso.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. Se o alvo atacou você ou um aliado no turno anterior, causa dano adicional igual ao mod Constituição.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 7,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe da Tempestade de Ferro", Edition: "4e", ClassID: &id,
			Description: "Uma série de golpes rápidos e furiosos abate múltiplos inimigos ao redor de você.",
			Keywords: "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força de dano por ataque.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Golpe da Montanha", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe com a força de uma montanha, esmagando as defesas do inimigo.",
			Keywords: "Arma, Confiável, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. Fortitude",
			Hit:  "4[A] + mod Força de dano. O alvo fica imobilizado e sofre -4 em todas as defesas (TR encerra).",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Terrível", Edition: "4e", ClassID: &id,
			Description: "Seu ataque é tão devastador que faz os inimigos próximos recuarem de medo.",
			Keywords: "Arma, Marcial, Medo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:    "3[A] + mod Força de dano. Os inimigos adjacentes ao alvo ficam amedrontados até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Postura Invencível", Edition: "4e", ClassID: &id,
			Description: "Você entra em uma postura defensiva formidável que torna você quase impossível de abater.",
			Keywords: "Marcial, Postura",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura Invencível. Enquanto nessa postura: você recebe resistência 5 a todos os danos e +2 de bônus de poder em todas as defesas.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Inspirar o Batalhão", Edition: "4e", ClassID: &id,
			Description: "Sua presença imponente no campo de batalha inspira todos os aliados a lutarem mais bravamente.",
			Keywords: "Marcial",
			ActionType: "Ação Mínima", Range: "Explosão contígua 5",
			Effect: "Os aliados dentro da explosão recebem +2 de bônus de poder nas jogadas de ataque e +2 de bônus de poder no dano até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Guerreiro 4e: %d habilidades processadas", len(skills))
}

func seedLadinoSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Ladino", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Ladino 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Floreio Ardiloso", Edition: "4e", ClassID: &id,
			Description: "Seus floreios causam uma distração e o inimigo quase se esquece da lâmina perto da própria garganta.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza + mod Carisma de dano.",
			Special:      "Condição: o ladino deve empunhar uma lâmina leve, besta ou funda.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Decidido", Edition: "4e", ClassID: &id,
			Description: "Um golpe bem-colocado lhe permite adquirir uma posição mais vantajosa.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano.",
			Special:      "Condição: lâmina leve, besta ou funda. O ladino pode se mover 2 quadrados antes do ataque.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe em Resposta", Edition: "4e", ClassID: &id,
			Description: "Com um golpe calculado, você deixa o oponente vulnerável a um contragolpe súbito caso ele ouse atacá-lo.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano. Se o alvo atacar o personagem antes do começo do próximo turno do ladino, o ladino realiza um contragolpe usando uma interrupção imediata: Força vs. CA — 1[A] + mod Força de dano.",
			Special:      "Condição: o ladino deve empunhar uma lâmina leve.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza no ataque e no contragolpe.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Perfurante", Edition: "4e", ClassID: &id,
			Description: "Uma ponta afiada atravessa a armadura do alvo, penetrando na carne.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. Reflexos",
			Hit:          "1[A] + mod Destreza de dano.",
			Special:      "Condição: o ladino deve empunhar uma lâmina leve.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Castelo Real", Edition: "4e", ClassID: &id,
			Description: "É difícil atingir um baixinho quando ele se esconde atrás de um aliado capaz de esmagar placas de aço com os dentes.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura", Attack: "Destreza vs. Reflexos",
			Hit:    "2[A] + mod Destreza de dano.",
			Effect: "O ladino troca de lugar com um aliado adjacente e voluntário.",
			Special: "Condição: o ladino deve empunhar uma lâmina leve, besta ou funda.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe de Posicionamento", Edition: "4e", ClassID: &id,
			Description: "Um tropeço em falso e um impulso súbito posicionam o inimigo exatamente onde você queria.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. Vontade",
			Hit:     "1[A] + mod Destreza de dano. O alvo é conduzido 1 quadrado.",
			Special: "Esquivo Hábil: o alvo é conduzido um número de quadrados igual ao modificador de Carisma do ladino.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe Pasmoso", Edition: "4e", ClassID: &id,
			Description: "Um golpe tático apanha seu adversário de surpresa e o deixa gemendo de dor.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "1[A] + mod Destreza de dano. O alvo fica pasmo até o final do próximo turno do ladino.",
			Special: "Condição: o ladino deve empunhar uma lâmina leve.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe Tortuoso", Edition: "4e", ClassID: &id,
			Description: "Se você girar a lâmina depois de atingir um golpe, pode fazer com que seu inimigo urre de dor.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "2[A] + mod Destreza de dano.",
			Special: "Vigarista Brutal: o ladino recebe um bônus na jogada de dano igual ao modificador de Força.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Golpe Ardiloso", Edition: "4e", ClassID: &id,
			Description: "Você finge um ataque e então desfere um golpe devastador que marca o inimigo para o resto do encontro.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:  "2[A] + mod Destreza de dano. Até o final do encontro, o alvo concede vantagem de combate ao ladino.",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Veneno da Mente Precisa", Edition: "4e", ClassID: &id,
			Description: "Você atinge o oponente num ponto vital, enfraquecendo-o com veneno que atua diretamente na mente.",
			Keywords:   "Arma, Marcial, Veneno",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:  "2[A] + mod Destreza de dano venenoso. O alvo fica enfraquecido (TR encerra).",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Esquiva", Edition: "4e", ClassID: &id,
			Description: "Você desvia o ataque no último segundo com sua agilidade natural.",
			Keywords:   "Marcial",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special: "Gatilho: O ladino é alvo de um ataque.",
			Effect:  "O ladino recebe +2 de bônus em uma defesa à sua escolha contra esse ataque.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Preparar Armadilha", Edition: "4e", ClassID: &id,
			Description: "Você prepara uma armadilha mortal para o inimigo que se aventurar a se aproximar.",
			Keywords:   "Marcial",
			ActionType: "Ação Padrão", Range: "Pessoal",
			Effect: "O ladino prepara uma emboscada. Até o final do próximo turno, o primeiro inimigo que se mover para um quadrado adjacente sofre 2d6 + mod Destreza de dano.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Deslizante", Edition: "4e", ClassID: &id,
			Description: "Com um movimento fluido, você ataca e desliza para uma posição mais segura.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano. O ladino pode se deslocar 1 quadrado sem provocar ataques de oportunidade.",
			Special:      "Condição: lâmina leve.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},
		{
			Name: "Golpe Traiçoeiro", Edition: "4e", ClassID: &id,
			Description: "Você engana o inimigo com um ataque falso e então golpeia onde ele menos espera.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano. Se o alvo tiver vantagem de combate, causa dano adicional igual ao mod Carisma.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Faca nas Sombras", Edition: "4e", ClassID: &id,
			Description: "Você ataca das sombras, surpreendendo o inimigo e desaparecendo antes que ele possa reagir.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "2[A] + mod Destreza de dano. Se o ladino estiver oculto, o alvo fica cego até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Manobra de Distração", Edition: "4e", ClassID: &id,
			Description: "Você distrai o inimigo enquanto um aliado o flanqueia.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "1[A] + mod Destreza de dano. Até o final do próximo turno, todos os aliados têm vantagem de combate contra o alvo.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Golpe do Assassino", Edition: "4e", ClassID: &id,
			Description: "Um golpe certeiro e letal que explora a abertura criada pela distração.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:  "3[A] + mod Destreza de dano. Se o alvo tiver vantagem de combate, causa dano adicional igual a 2d6.",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Sombrio", Edition: "4e", ClassID: &id,
			Description: "Você canaliza as sombras ao seu redor para potencializar seu ataque.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano. Se não houver luz brilhante no quadrado do alvo, o dano aumenta em 1d6.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 5,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Dança das Lâminas", Edition: "4e", ClassID: &id,
			Description: "Você executa uma série de ataques rápidos que mantêm o inimigo constantemente na defensiva.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma ou duas criaturas",
			Attack:  "Destreza vs. CA (dois ataques)",
			Hit:     "1[A] + mod Destreza de dano por ataque.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Veneno Paralisante", Edition: "4e", ClassID: &id,
			Description: "Você aplica um veneno potente que gradualmente paralisa o inimigo.",
			Keywords:   "Arma, Marcial, Veneno",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. Fortitude",
			Hit:  "2[A] + mod Destreza de dano venenoso. O alvo fica lento (TR encerra). Se fracassar no primeiro TR, fica imobilizado (TR encerra).",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Esconder nas Sombras", Edition: "4e", ClassID: &id,
			Description: "Você se funde com as sombras, tornando-se praticamente invisível.",
			Keywords:   "Marcial",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O ladino fica oculto até o início do próximo turno, mesmo se não houver cobertura adequada.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Velocidade das Sombras", Edition: "4e", ClassID: &id,
			Description: "Você se move com a velocidade e a graciosidade das sombras.",
			Keywords:   "Marcial",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O ladino se move até o dobro de sua velocidade e não provoca ataques de oportunidade durante esse movimento.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Toque Envenenado", Edition: "4e", ClassID: &id,
			Description: "Cada golpe seu deposita uma dose de veneno que enfraquece gradualmente o inimigo.",
			Keywords:   "Arma, Marcial, Veneno",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano venenoso. O alvo sofre -1 cumulativo em todas as defesas até o final do encontro (máximo -3).",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 7,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Fantasma", Edition: "4e", ClassID: &id,
			Description: "Você ataca com tal velocidade que parece um fantasma, deixando o inimigo confuso e desorientado.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "2[A] + mod Destreza de dano. O alvo fica desorientado e trata todos os aliados do ladino como se tivessem vantagem de combate até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Rastro de Sangue", Edition: "4e", ClassID: &id,
			Description: "Você inflige uma ferida que continua sangrando, marcando o inimigo para a morte.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Destreza vs. CA",
			Hit:  "3[A] + mod Destreza de dano. O alvo sofre 10 de dano contínuo (TR encerra) e fica sangrando, tornando mais fácil rastreá-lo.",
			Miss: "Metade do dano e 5 de dano contínuo.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe do Escorpião", Edition: "4e", ClassID: &id,
			Description: "Você imita o golpe mortal do escorpião, atingindo um ponto vital do inimigo.",
			Keywords:   "Arma, Marcial, Veneno",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. Fortitude",
			Hit:     "2[A] + mod Destreza de dano venenoso. O alvo fica imobilizado até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Arte da Fuga", Edition: "4e", ClassID: &id,
			Description: "Você usa sua habilidade de fuga para escapar de qualquer situação difícil.",
			Keywords:   "Marcial",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O ladino escapa de qualquer efeito que o esteja imobilizando, lentificando ou prendendo, e se move até sua velocidade.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
		{
			Name: "Névoa das Sombras", Edition: "4e", ClassID: &id,
			Description: "Você lança uma névoa de sombras que oculta você e seus aliados dos olhos inimigos.",
			Keywords:   "Marcial",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Effect: "Até o final do próximo turno do ladino, você e seus aliados dentro da explosão ficam ocultos para os inimigos.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Ladino 4e: %d habilidades processadas", len(skills))
}

func seedMagoSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Mago", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Mago 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── TRUQUES (SEM LIMITE) ─────────────────────────────────────
		{
			Name: "Prestidigitação", Edition: "4e", ClassID: &id,
			Description: "Você realiza um truque de mágica como diversão, criando um filete de luz que dança, rejuvenescendo uma flor murcha ou fazendo uma moeda desaparecer.",
			Keywords:   "Arcano",
			ActionType: "Ação Padrão", Range: "À distância 2",
			Effect: "Simula um dos efeitos: mover 0,5 kg de matéria; criar efeito sensorial inofensivo; colorir, limpar ou manchar itens; acender ou apagar chamas; criar pequena marca numa superfície; criar pequeno item ilusório até o final do próximo turno do mago.",
			Special: "É possível manter somente três truques ativos simultaneamente.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Som Fantasma", Edition: "4e", ClassID: &id,
			Description: "Com um gesto, você cria um som ilusório que emana de qualquer local próximo.",
			Keywords:   "Arcano, Ilusão",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Um objeto ou espaço desocupado",
			Effect: "O personagem cria um som — tão sutil como um sussurro ou tão elevado quanto gritos — que emana do alvo escolhido. O som pode ser qualquer som que o mago conheça.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Explosão Incandescente", Edition: "4e", ClassID: &id,
			Description: "Uma coluna vertical de chamas douradas incinera seus adversários.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Reflexos",
			Hit:          "1d6 + mod Inteligência de dano flamejante.",
			LevelScaling: "Nível 21: 2d6 + mod Inteligência.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Mísseis Mágicos", Edition: "4e", ClassID: &id,
			Description: "Você dispara uma rajada de energia prateada contra o inimigo.",
			Keywords:   "Arcano, Energético, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 20",
			Target: "Uma criatura", Attack: "Inteligência vs. Reflexos",
			Hit:          "2d4 + mod Inteligência de dano energético.",
			Special:      "Este poder é um ataque básico à distância. Quando um poder permitir que o mago realize um ataque básico à distância, ele pode disparar um míssil mágico.",
			LevelScaling: "Nível 21: 4d4 + mod Inteligência.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Nuvem de Adagas", Edition: "4e", ClassID: &id,
			Description: "Você cria um pequeno turbilhão de adagas de energia que atacam implacavelmente as criaturas na área.",
			Keywords:   "Arcano, Energético, Implemento",
			ActionType: "Ação Padrão", Range: "Área de 1 quadrado a até 10 quadrados",
			Target: "As criaturas dentro do quadrado", Attack: "Inteligência vs. Reflexos",
			Hit:          "1d6 + mod Inteligência de dano energético.",
			Effect:       "A área fica repleta de adagas compostas de energia. Qualquer criatura que ingressar ou começar seu turno dentro da área sofre dano energético igual ao modificador de Sabedoria do mago (mínimo 1). A nuvem permanece até o final do próximo turno.",
			LevelScaling: "Nível 21: 2d6 + mod Inteligência.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Onda Trovejante", Edition: "4e", ClassID: &id,
			Description: "Você cria uma ondulação de poder sônico que emerge do solo e empurra os inimigos.",
			Keywords:   "Arcano, Implemento, Trovejante",
			ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target: "As criaturas dentro da rajada", Attack: "Inteligência vs. Fortitude",
			Hit:          "1d6 + mod Inteligência de dano trovejante. O alvo é empurrado um número de quadrados igual ao modificador de Sabedoria do mago.",
			LevelScaling: "Nível 21: 2d6 + mod Inteligência.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Raio Álgido", Edition: "4e", ClassID: &id,
			Description: "Um raio brilhante de gelo esbranquiçado é disparado contra o alvo, congelando seus movimentos.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Inteligência vs. Fortitude",
			Hit:          "1d6 + mod Inteligência de dano congelante. O alvo fica lento até o final do próximo turno do mago.",
			LevelScaling: "Nível 21: 2d6 + mod Inteligência.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Resfriante", Edition: "4e", ClassID: &id,
			Description: "Você cria um raio de energia púrpura e frígida em torno da sua mão e o arremessa diretamente contra o inimigo.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Inteligência vs. Fortitude",
			Hit: "2d8 + mod Inteligência de dano congelante. O alvo fica pasmo até o final do próximo turno do mago.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Mãos Ardentes", Edition: "4e", ClassID: &id,
			Description: "Uma explosão de chamas terríveis irrompe de suas mãos e incinera os inimigos mais próximos.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "Rajada contígua 5",
			Target: "As criaturas dentro da rajada", Attack: "Inteligência vs. Reflexos",
			Hit: "2d6 + mod Inteligência de dano flamejante.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Orbe de Energia", Edition: "4e", ClassID: &id,
			Description: "Você arremessa um orbe de energia que explode no alvo, arremessando estilhaços de energia afiados como navalhas.",
			Keywords:   "Arcano, Energético, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 20",
			Target: "Uma criatura ou objeto (primário)", Attack: "Inteligência vs. Reflexos",
			Hit: "2d8 + mod Inteligência de dano energético. Realize um ataque secundário: um inimigo adjacente ao alvo primário — Inteligência vs. Reflexos — 1d10 + mod Inteligência de dano energético.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Raio de Enfraquecimento", Edition: "4e", ClassID: &id,
			Description: "Uma névoa esverdeada e macabra brota da carne do inimigo, carregando consigo a força do alvo.",
			Keywords:   "Arcano, Implemento, Necrótico",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Inteligência vs. Fortitude",
			Hit: "1d10 + mod Inteligência de dano necrótico. O alvo fica enfraquecido até o final do próximo turno do mago.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Terreno Gélido", Edition: "4e", ClassID: &id,
			Description: "Com um sopro congelante, você murmura uma única palavra arcana que cria uma trilha traiçoeira de gelo no chão.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Reflexos",
			Hit:    "1d6 + mod Inteligência de dano congelante. O alvo fica derrubado.",
			Effect: "A área afetada se torna terreno acidentado até o final do próximo turno do mago.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Esfera Flamejante", Edition: "4e", ClassID: &id,
			Description: "Você conjura uma esfera de fogo em movimento e controla seu deslocamento pelo campo de batalha.",
			Keywords:   "Arcano, Conjuração, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Effect:  "O mago conjura uma esfera flamejante Média. Qualquer criatura que começar o turno adjacente à esfera sofre 1d4 + mod Inteligência de dano flamejante. Com ação de movimento, o mago pode deslocar a esfera até 6 quadrados.",
			Special: "Sustentação Mínima: o mago pode realizar um novo ataque com a esfera (Inteligência vs. Reflexos contra criatura adjacente — 2d6 + mod INT de dano flamejante).",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Flecha Ácida", Edition: "4e", ClassID: &id,
			Description: "Uma flecha bruxuleante, composta de um líquido verde e brilhante, é disparada contra o alvo e explode numa rajada de ácido corrosivo.",
			Keywords:   "Ácido, Arcano, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 20",
			Target: "Uma criatura (primário)", Attack: "Inteligência vs. Reflexos",
			Hit:  "2d8 + mod Inteligência de dano ácido e 5 de dano ácido contínuo (TR encerra). Realize um ataque secundário: criaturas adjacentes ao alvo primário — Inteligência vs. Reflexos — 1d8 + mod INT de dano ácido e 5 de dano contínuo.",
			Miss: "Metade do dano e 2 de dano ácido contínuo no alvo primário, sem ataque secundário.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Nuvem Congelante", Edition: "4e", ClassID: &id,
			Description: "Você dispara um projétil com a mão que explode numa nuvem de névoa congelada no ponto de impacto.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Fortitude",
			Hit:    "1d8 + mod Inteligência de dano congelante.",
			Miss:   "Metade do dano.",
			Effect: "A nuvem permanece ativa até o final do próximo turno do mago. Qualquer criatura que ingressar ou começar seu turno dentro da nuvem é alvo de outro ataque do mago.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Sono", Edition: "4e", ClassID: &id,
			Description: "Você impõe sua vontade contra os adversários, tentando sobreujá-los com uma onda de cansaço mágico.",
			Keywords:   "Arcano, Implemento, Sono",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 20 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Vontade",
			Hit:  "O alvo fica lento (TR encerra). Se fracassar no primeiro teste de resistência contra esse poder, o alvo fica inconsciente (TR encerra).",
			Miss: "O alvo fica lento (TR encerra).",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Escudo Arcano", Edition: "4e", ClassID: &id,
			Description: "Um escudo de energia arcana surge ao seu redor no momento em que você é atacado.",
			Keywords:   "Arcano",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O mago é atingido por um ataque.",
			Effect:  "O mago recebe +4 de bônus de poder na CA e Reflexos contra esse ataque.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Salto Arcano", Edition: "4e", ClassID: &id,
			Description: "Com um gesto, você dobra o espaço ao seu redor e se teletransporta a uma curta distância.",
			Keywords:   "Arcano, Teletransporte",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O mago se teletransporta até 3 quadrados.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Fogo Deslizante", Edition: "4e", ClassID: &id,
			Description: "Uma corrente de fogo escorrega pelo solo, queimando tudo que encontra em seu caminho.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Inteligência vs. Reflexos",
			Hit:    "2d6 + mod Inteligência de dano flamejante. O mago pode deslizar o alvo 2 quadrados.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Prisão de Gelo", Edition: "4e", ClassID: &id,
			Description: "Você envolve o inimigo em uma prisão de gelo que o mantém imobilizado.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura", Attack: "Inteligência vs. Fortitude",
			Hit: "2d8 + mod Inteligência de dano congelante. O alvo fica imobilizado até o final do próximo turno do mago.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Bola de Fogo", Edition: "4e", ClassID: &id,
			Description: "Uma bola de fogo explode com força devastadora, queimando todos ao redor.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Reflexos",
			Hit:  "3d6 + mod Inteligência de dano flamejante.",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Raio Flamejante", Edition: "4e", ClassID: &id,
			Description: "Um raio concentrado de fogo puro atinge o alvo com precisão devastadora.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Inteligência vs. Reflexos",
			Hit:          "1d10 + mod Inteligência de dano flamejante.",
			LevelScaling: "Nível 21: 2d10 + mod Inteligência.",
			PowerType: domain.PowerUnlimited, Level: 5,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Tempestade de Gelo", Edition: "4e", ClassID: &id,
			Description: "Uma tempestade de fragmentos de gelo aflige seus inimigos, desacelerando seus movimentos.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Fortitude",
			Hit:    "2d6 + mod Inteligência de dano congelante. Os alvos ficam lentos (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Tempestade Relampejante", Edition: "4e", ClassID: &id,
			Description: "Uma tempestade de raios surge ao seu redor, atingindo múltiplos inimigos com descargas elétricas.",
			Keywords:   "Arcano, Elétrico, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão contígua 2",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Reflexos",
			Hit:    "3d6 + mod Inteligência de dano elétrico.",
			Miss:   "Metade do dano.",
			Effect: "Até o final do encontro, qualquer criatura que iniciar o turno dentro da zona sofre 1d6 de dano elétrico.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Invisibilidade", Edition: "4e", ClassID: &id,
			Description: "Você lança um feitiço que torna você temporariamente invisível.",
			Keywords:   "Arcano, Ilusão",
			ActionType: "Ação Padrão", Range: "Pessoal",
			Effect:  "O mago fica invisível até o início do próximo turno ou até realizar um ataque.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Levitação", Edition: "4e", ClassID: &id,
			Description: "Você se eleva pelo ar com o poder da magia arcana.",
			Keywords:   "Arcano",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O mago voa até 4 quadrados na vertical ou horizontal. Permanece no ar até o início do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Raio em Cadeia", Edition: "4e", ClassID: &id,
			Description: "Um raio de eletricidade salta de inimigo em inimigo, atingindo vários alvos.",
			Keywords:   "Arcano, Elétrico, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura (primário)", Attack: "Inteligência vs. Reflexos",
			Hit: "2d6 + mod Inteligência de dano elétrico. Realize ataques secundários contra até 2 criaturas dentro de 10 quadrados do alvo primário: Inteligência vs. Reflexos — 1d6 + mod INT de dano elétrico cada.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Nuvem Mortal", Edition: "4e", ClassID: &id,
			Description: "Você conjura uma nuvem de névoa venenosa que persiste no campo de batalha, envenenando todos que entrarem.",
			Keywords:   "Arcano, Implemento, Venenoso, Zona",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Fortitude",
			Hit:    "2d10 + mod Inteligência de dano venenoso.",
			Miss:   "Metade do dano.",
			Effect: "A zona persiste até o final do encontro. Criaturas que iniciarem o turno dentro sofrem 2d10 de dano venenoso.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Explosão Arcana", Edition: "4e", ClassID: &id,
			Description: "Uma explosão maciça de energia arcana pura devastadoras seus inimigos.",
			Keywords:   "Arcano, Energético, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "As criaturas dentro da explosão", Attack: "Inteligência vs. Reflexos",
			Hit: "3d8 + mod Inteligência de dano energético.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Teletransporte", Edition: "4e", ClassID: &id,
			Description: "Você dominou o poder do teletransporte e pode se mover instantaneamente a grandes distâncias.",
			Keywords:   "Arcano, Teletransporte",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O mago se teletransporta até 8 quadrados.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
		{
			Name: "Contingência", Edition: "4e", ClassID: &id,
			Description: "Você prepara um feitiço para ser ativado automaticamente quando certas condições forem atendidas.",
			Keywords:   "Arcano",
			ActionType: "Ação Padrão", Range: "Pessoal",
			Effect: "Você prepara um feitiço de nível 1 sem limite. Quando você for reduzido a 0 PV ou menos, esse feitiço é ativado automaticamente contra o inimigo mais próximo, sem custar uma ação.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Mago 4e: %d habilidades processadas", len(skills))
}

func seedPaladinoSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Paladino", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Paladino 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Canalizar Divindade: Entusiasmo Divino", Edition: "4e", ClassID: &id,
			Description: "Você invoca o poder da sua divindade para inflamar seus aliados com fervor sagrado.",
			Keywords:   "Divino",
			ActionType: "Ação Livre", Range: "Pessoal",
			Effect:         "O paladino e todos os aliados a até 5 quadrados recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType:      domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Canalizar Divindade: Força Divina", Edition: "4e", ClassID: &id,
			Description: "Você invoca o poder da sua divindade para fortalecê-lo em combate.",
			Keywords:   "Divino",
			ActionType: "Ação Livre", Range: "Pessoal",
			Effect:         "Até o final do turno, o paladino adiciona o modificador de Força a uma jogada de dano.",
			PowerType:      domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe Debilitante", Edition: "4e", ClassID: &id,
			Description: "O ataque brutal da sua arma deixa seu oponente enfraquecido e menos capaz de ameaçar seus aliados.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit:          "1[A] + mod Carisma de dano. Se o paladino marcou o alvo, ele sofre -2 de penalidade nas jogadas de ataque até o final do próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Estimulante", Edition: "4e", ClassID: &id,
			Description: "Você ataca sem misericórdia e sua precisão é recompensada com uma dádiva de vigor divino.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit:          "1[A] + mod Carisma de dano. O paladino recebe um número de pontos de vida temporários igual ao seu modificador de Sabedoria.",
			LevelScaling: "Nível 21: 2[A] + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Sagrado", Edition: "4e", ClassID: &id,
			Description: "Você ataca um inimigo com sua arma, envolvendo-a com uma explosão de luz sagrada purificadora.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:          "1[A] + mod Força de dano radiante. Se tiver marcado o alvo, o paladino recebe um bônus na jogada de dano igual ao seu modificador de Sabedoria.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Valente", Edition: "4e", ClassID: &id,
			Description: "Conforme ergue sua arma para o golpe, sua desvantagem numérica se torna uma vantagem e fortalece seus ataques.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força + 1 por inimigo adjacente vs. CA",
			Hit:          "1[A] + mod Força de dano.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Punição Perfurante", Edition: "4e", ClassID: &id,
			Description: "Fagulhas prateadas recobrem a sua arma, atravessando a armadura do oponente e marcando os aliados dele.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. Reflexos",
			Hit: "2[A] + mod Força de dano. Um número de inimigos adjacentes ao personagem igual ao modificador de Sabedoria ficam marcados até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Punição Protetora", Edition: "4e", ClassID: &id,
			Description: "Um escudo dourado e translúcido se materializa diante de um aliado próximo conforme você ataca com sua arma.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "2[A] + mod Carisma de dano. Até o final do próximo turno do paladino, um aliado a até 5 quadrados recebe um bônus de poder na CA igual ao modificador de Sabedoria do paladino.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Punição Radiante", Edition: "4e", ClassID: &id,
			Description: "Sua arma brilha com luminescência perolada. Os inimigos se encolhem diante da pureza dessa luz.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "1[A] + mod Força + mod Sabedoria de dano radiante.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Punição Temível", Edition: "4e", ClassID: &id,
			Description: "Quando você ataca com sua arma, a força do golpe faz com que o adversário se estremeça e revele as próprias táticas.",
			Keywords:   "Arma, Divino, Medo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. CA",
			Hit: "2[A] + mod Carisma de dano. Até o final do próximo turno do paladino, o alvo sofre uma penalidade nas jogadas de ataque igual ao modificador de Sabedoria do personagem.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Delírio Radiante", Edition: "4e", ClassID: &id,
			Description: "Você engolfa seu inimigo em faixas incandescentes de radiação que o cegam e enfraquecem.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura", Attack: "Carisma vs. Reflexos",
			Hit:  "3d8 + mod Carisma de dano radiante. O alvo fica pasmo até o final do próximo turno do paladino. Além disso, ele sofre -2 de penalidade na CA (TR encerra).",
			Miss: "Metade do dano. O alvo fica pasmo até o final do próximo turno do paladino.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Julgamento do Paladino", Edition: "4e", ClassID: &id,
			Description: "Seus ataques corpo a corpo esmagam seu inimigo e revigoram um aliado com energia divina.",
			Keywords:   "Arma, Cura, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:  "3[A] + mod Força de dano. Um aliado do paladino a até 5 quadrados pode gastar um pulso de cura.",
			Miss: "Um aliado do paladino a até 5 quadrados pode gastar um pulso de cura.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Cura das Mãos Leigas", Edition: "4e", ClassID: &id,
			Description: "Com um toque suave e uma prece fervorosa, você restaura instantaneamente a vitalidade de um aliado ferido.",
			Keywords:   "Cura, Divino",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target: "O personagem ou um aliado adjacente",
			Effect: "O alvo recupera PV como se tivesse gasto um pulso de cura, mais PV adicionais iguais ao modificador de Sabedoria do paladino.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Sagrada Promessa", Edition: "4e", ClassID: &id,
			Description: "Você faz uma promessa sagrada de proteção a um aliado ferido, garantindo que ele possa se curar rapidamente.",
			Keywords:   "Divino",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "Até o final do próximo turno do paladino, quando o aliado sofrer dano, ele pode gastar um pulso de cura como resposta imediata.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Bênção de Batalha", Edition: "4e", ClassID: &id,
			Description: "Você clama pela bênção da sua divindade para enrijecer você e seus aliados contra os golpes dos inimigos.",
			Keywords:   "Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Effect: "Você e os aliados dentro da explosão recebem +2 de bônus de poder em todas as defesas até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Golpe da Fé Inabalável", Edition: "4e", ClassID: &id,
			Description: "Sua fé inabalável potencializa seu golpe, causando dano sagrado que queima o mal.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "2[A] + mod Força + mod Sabedoria de dano radiante. Se o alvo for um morto-vivo ou demônio, causa dano adicional igual a 1d6.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Martelo da Justiça", Edition: "4e", ClassID: &id,
			Description: "Um martelo de luz divina desce sobre o inimigo com força esmagadora, quebrando suas defesas.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura", Attack: "Carisma vs. Fortitude",
			Hit:  "3d8 + mod Carisma de dano radiante. O alvo sofre -4 em todas as defesas até o final do próximo turno.",
			Miss: "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe do Servo Sagrado", Edition: "4e", ClassID: &id,
			Description: "O poder da sua divindade flui através de você, transformando seu golpe numa explosão de energia sagrada.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano radiante. Todos os aliados a até 5 quadrados recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Punição Arrebatadora", Edition: "4e", ClassID: &id,
			Description: "Seu golpe é tão poderoso que lança o inimigo pelos ares e o mantém desorientado.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Carisma vs. Fortitude",
			Hit: "2[A] + mod Carisma de dano. O alvo é empurrado 3 quadrados e fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Aura de Proteção", Edition: "4e", ClassID: &id,
			Description: "Você emana uma aura de proteção divina que escuda você e seus aliados de dano.",
			Keywords:   "Divino",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você emite uma aura 2 até o final do encontro. Enquanto dentro da aura, você e seus aliados recebem resistência 5 a todos os danos.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Palavra de Cura", Edition: "4e", ClassID: &id,
			Description: "Você profere uma palavra sagrada que restaura a saúde de um aliado próximo.",
			Keywords:   "Cura, Divino",
			ActionType: "Ação Mínima", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O alvo pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Carisma do paladino.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Purificar", Edition: "4e", ClassID: &id,
			Description: "Você invoca a pureza divina para remover condições negativas de um aliado.",
			Keywords:   "Cura, Divino",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target: "Um aliado",
			Effect: "Remova um dos seguintes efeitos do alvo: lento, imobilizado, enfraquecido ou amedrontado.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Devastador Sagrado", Edition: "4e", ClassID: &id,
			Description: "Você descarrega toda a força da sua fé num único golpe avassalador.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "3[A] + mod Força + mod Sabedoria de dano radiante. O alvo fica cego até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Escudo da Fé Divina", Edition: "4e", ClassID: &id,
			Description: "Você cria um campo de energia divina que protege você e seus aliados contra todos os males.",
			Keywords:   "Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Effect: "Até o final do encontro, você e os aliados dentro da explosão recebem +2 de bônus de poder em todas as defesas e resistência 5 a dano necrótico e venenoso.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Punição Celestial", Edition: "4e", ClassID: &id,
			Description: "A ira dos céus desce sobre o inimigo através do seu golpe, causando dano sagrado avassalador.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura", Attack: "Força vs. CA",
			Hit: "3[A] + mod Força + mod Sabedoria de dano radiante. O alvo fica imobilizado e sofre dano radiante contínuo igual ao mod Carisma (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Ressurgimento Divino", Edition: "4e", ClassID: &id,
			Description: "Quando um aliado cai, você invoca o poder divino para ressuscitá-lo imediatamente.",
			Keywords:   "Cura, Divino",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Um aliado com 0 PV ou menos",
			Effect: "O alvo recupera PV iguais ao máximo de um pulso de cura mais PV adicionais iguais ao mod Carisma do paladino. O alvo pode realizar uma ação imediata.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Aura do Campeão", Edition: "4e", ClassID: &id,
			Description: "Sua presença no campo de batalha irradia um poder tão divino que fortalece todos os aliados ao redor.",
			Keywords:   "Divino",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você emite uma aura 3 até o final do encontro. Aliados dentro da aura recebem +2 de bônus de poder nas jogadas de ataque e nos danos.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Paladino 4e: %d habilidades processadas", len(skills))
}

func seedPatrulheiroSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Patrulheiro", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Patrulheiro 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Ataque Cauteloso", Edition: "4e", ClassID: &id,
			Description: "Você estuda o inimigo, procurando uma brecha nas defesas dele. Somente quando encontra é que você ataca.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target:       "Uma criatura",
			Attack:       "Força + 2 vs. CA (corpo a corpo) ou Destreza + 2 vs. CA (à distância)",
			Hit:          "1[A] de dano.",
			Special:      "Condição: o patrulheiro deve empunhar duas armas de combate corpo a corpo ou uma arma de combate à distância.",
			LevelScaling: "Nível 21: 2[A] de dano.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Bater e Correr", Edition: "4e", ClassID: &id,
			Description: "Permitindo que o guerreiro encare o monstro de frente, você prefere desferir seu ataque e recuar para um local mais seguro.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano.",
			Effect:       "Se o patrulheiro se mover nesse turno, após esse ataque, abandonar o primeiro quadrado adjacente ao alvo não provoca ataques de oportunidade em favor do alvo.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Lépido", Edition: "4e", ClassID: &id,
			Description: "Você escapa para a lateral da guarda do seu inimigo para desferir seu ataque, ou ataca e então recua para uma posição mais vantajosa.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma à distância",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza de dano.",
			Special:      "O patrulheiro ajusta um quadrado antes ou depois de realizar esse ataque.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpes Gêmeos", Edition: "4e", ClassID: &id,
			Description: "Se o primeiro ataque não matá-lo, o segundo deve fazê-lo.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target:       "Uma ou duas criaturas",
			Attack:       "Força vs. CA (corpo a corpo; arma principal e arma da mão inábil) ou Destreza vs. CA (à distância), dois ataques",
			Hit:          "1[A] de dano por ataque.",
			Special:      "Condição: o patrulheiro deve empunhar duas armas de combate corpo a corpo ou uma arma de combate à distância.",
			LevelScaling: "Nível 21: 2[A] de dano.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Astúcia da Raposa", Edition: "4e", ClassID: &id,
			Description: "Usando o impulso do golpe do seu inimigo para recuar ou escapar para o lado, você desfere um ataque súbito de retaliação conforme ele tropeça.",
			Keywords:   "Arma, Marcial",
			ActionType: "Reação Imediata", Range: "Arma corpo a corpo ou à distância",
			Special: "Gatilho: Um inimigo realiza um ataque corpo a corpo contra o personagem.",
			Effect:  "O patrulheiro pode ajustar 1 quadrado, então realizar um ataque básico contra o inimigo. O patrulheiro recebe um bônus de poder no ataque igual ao modificador de Sabedoria.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe das Duas Presas", Edition: "4e", ClassID: &id,
			Description: "Você enfia duas flechas ou lâminas na carne do seu inimigo, fazendo-o uivar de dor.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura",
			Attack: "Força vs. CA (corpo a corpo; arma principal e arma da mão inábil) ou Destreza vs. CA (à distância), dois ataques",
			Hit:    "1[A] + mod Força (corpo a corpo) ou 1[A] + mod Destreza (à distância) de dano por ataque. Se os dois ataques obtiverem sucesso, o patrulheiro causa dano adicional igual ao modificador de Sabedoria.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe do Carcajú Atroz", Edition: "4e", ClassID: &id,
			Description: "Os inimigos o cercam — para o desgosto deles, você os corta em pedaços com a ferocidade de um carcajú atroz ferido.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target:  "Os inimigos dentro da explosão e na linha de visão",
			Attack:  "Força vs. CA",
			Hit:     "1[A] + mod Força de dano.",
			Special: "Condição: o patrulheiro deve empunhar duas armas de combate corpo a corpo.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe Evasivo", Edition: "4e", ClassID: &id,
			Description: "Você confunde os inimigos serpenteando ileso pelo campo de batalha enquanto desfere seus ataques.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target:  "Uma criatura",
			Attack:  "Força vs. CA (corpo a corpo) ou Destreza vs. CA (à distância)",
			Hit:     "2[A] + mod Força (corpo a corpo) ou 2[A] + mod Destreza (à distância) de dano.",
			Special: "O patrulheiro pode ajustar um número de quadrados igual a 1 + seu modificador de Sabedoria antes ou depois de desferir esse ataque.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Armadilha do Caçador de Ursos", Edition: "4e", ClassID: &id,
			Description: "Um tiro ou golpe bem colocado na perna deixa seu inimigo mancando e sangrando.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura",
			Attack: "Força vs. CA (corpo a corpo) ou Destreza vs. CA (à distância)",
			Hit:    "2[A] + mod Força (corpo a corpo) ou 2[A] + mod Destreza (à distância) de dano. O alvo fica lento e sofre 5 de dano contínuo (TR encerra).",
			Miss:   "Metade do dano, sem dano contínuo. O alvo fica lento até o final do próximo turno do patrulheiro.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Golpe Súbito", Edition: "4e", ClassID: &id,
			Description: "Você empunha suas armas com a lâmina para baixo e corta a face do oponente com uma delas. Conforme a vítima se afasta e baixa a guarda, você rola de lado, se ergue e crava a outra lâmina nas costas do inimigo.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Força vs. CA (arma da mão inábil)",
			Hit:     "1[A] de dano (arma da mão inábil). O patrulheiro ajusta 1 quadrado e desfere um ataque secundário: Força vs. CA (arma principal) — Sucesso: 2[A] + mod Força de dano e o alvo fica enfraquecido até o final do próximo turno do patrulheiro.",
			Miss:    "Metade do dano.",
			Special: "Condição: o patrulheiro deve empunhar duas armas de combate corpo a corpo.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Mandíbulas do Lobo", Edition: "4e", ClassID: &id,
			Description: "Você usa suas armas para acuar seu inimigo e enganá-lo até que ele exponha um ponto fraco; nesse momento, você ataca.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Força vs. CA (arma principal e arma da mão inábil), dois ataques",
			Hit:     "2[A] + mod Força de dano por ataque.",
			Miss:    "Metade do dano em cada ataque.",
			Special: "Condição: o patrulheiro deve empunhar duas armas de combate corpo a corpo.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Repartir Disparo", Edition: "4e", ClassID: &id,
			Description: "Você dispara duas flechas de uma vez, que se separam no meio do voo para atingir dois alvos diferentes.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma à distância",
			Target:  "Duas criaturas a até 3 uma da outra",
			Attack:  "Destreza vs. CA. Realize duas jogadas de ataque e aplique o melhor resultado contra os dois alvos.",
			Hit:     "2[A] + mod Destreza de dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Passo do Lobo", Edition: "4e", ClassID: &id,
			Description: "Você se move com a graça e a velocidade de um lobo na floresta, ignorando o terreno que atrapalharia outros.",
			Keywords:   "Marcial",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O patrulheiro se desloca até a sua velocidade. Durante esse deslocamento, ignora terreno difícil e pode se mover através de espaços de criaturas aliadas sem penalidade.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Sentido da Caça", Edition: "4e", ClassID: &id,
			Description: "Você afina seus sentidos de predador, tornando-se mais eficaz ao atacar inimigos que não possam contar com o auxílio de aliados.",
			Keywords:   "Marcial",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do encontro, o patrulheiro ignora a penalidade de -2 nas jogadas de ataque quando um aliado não está adjacente ao alvo.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe do Falcão", Edition: "4e", ClassID: &id,
			Description: "Como um falcão em mergulho, você se lança sobre o inimigo com velocidade e precisão mortais.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma à distância",
			Target:       "Uma criatura",
			Attack:       "Destreza vs. CA",
			Hit:          "1[A] + mod Destreza + mod Sabedoria de dano.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Ataque Rápido", Edition: "4e", ClassID: &id,
			Description: "Você desfere dois ataques em rápida sucessão antes que o inimigo possa reagir.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura",
			Attack: "Força vs. CA (corpo a corpo) ou Destreza vs. CA (à distância), dois ataques",
			Hit:    "1[A] + mod Força ou mod Destreza de dano por ataque.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Armadilha do Predador", Edition: "4e", ClassID: &id,
			Description: "Você posiciona o inimigo exatamente onde quer que ele esteja para atacar com vantagem.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target:  "Uma criatura",
			Attack:  "Força vs. CA (corpo a corpo) ou Destreza vs. CA (à distância)",
			Hit:     "1[A] + mod Força ou mod Destreza de dano. O alvo fica imobilizado até o final do próximo turno do patrulheiro.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Frenesi do Lobo", Edition: "4e", ClassID: &id,
			Description: "Como um lobo em frenesi, você lança uma série de ataques implacáveis que deixam o inimigo sangrando.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA, três ataques",
			Hit:    "1[A] + mod Força de dano por ataque. Se todos os três ataques acertarem, o alvo sofre 5 de dano contínuo (TR encerra).",
			Miss:   "Metade do dano por ataque que acertar.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Flecha Certeira", Edition: "4e", ClassID: &id,
			Description: "Você mira com extrema precisão, encontrando a fresta na armadura do inimigo.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma à distância",
			Target:       "Uma criatura",
			Attack:       "Destreza + 2 vs. CA",
			Hit:          "1[A] + mod Destreza + mod Sabedoria de dano.",
			LevelScaling: "Nível 21: 2[A] + mod Destreza + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 5,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Chuva de Flechas", Edition: "4e", ClassID: &id,
			Description: "Você dispara uma saraivada de flechas que abate múltiplos inimigos de uma vez.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Destreza vs. CA",
			Hit:    "1[A] + mod Destreza de dano por alvo.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Perseguição Implacável", Edition: "4e", ClassID: &id,
			Description: "Você marca o inimigo como sua presa e o persegue implacavelmente pelo campo de batalha.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo ou à distância",
			Target: "Uma criatura",
			Attack: "Força vs. CA (corpo a corpo) ou Destreza vs. CA (à distância)",
			Hit:    "2[A] + mod Força ou mod Destreza de dano. Até o final do encontro, o patrulheiro pode usar a Ação de Movimento para se mover até sua velocidade em direção ao alvo sem provocar ataques de oportunidade.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Instinto de Sobrevivência", Edition: "4e", ClassID: &id,
			Description: "Seus instintos aguçados permitem que você reaja antes mesmo de perceber o perigo conscientemente.",
			Keywords:   "Marcial",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O patrulheiro é surpreendido ou alvo de um ataque de emboscada.",
			Effect:  "O patrulheiro não fica surpreendido e pode agir normalmente. Além disso, recebe +2 de bônus de poder na iniciativa até o final do encontro.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Terreno Favorável", Edition: "4e", ClassID: &id,
			Description: "Você usa o conhecimento do terreno ao seu redor para se posicionar estrategicamente.",
			Keywords:   "Marcial",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O patrulheiro se move até sua velocidade. Durante esse movimento, recebe +4 de bônus de poder na CA e Reflexos contra ataques de oportunidade.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe do Javali", Edition: "4e", ClassID: &id,
			Description: "Você investe como um javali furioso, derrubando o inimigo no chão.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Força vs. Fortitude",
			Hit:     "2[A] + mod Força de dano. O alvo fica derrubado. O patrulheiro pode ajustar 1 quadrado após o ataque.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Tempestade de Lâminas", Edition: "4e", ClassID: &id,
			Description: "Você gira num torvelinho de aço, atingindo todos os inimigos ao redor com ataques imparáveis.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Tiro do Atirador de Elite", Edition: "4e", ClassID: &id,
			Description: "Você mira com uma precisão sobre-humana, enviando um projétil que atravessa qualquer defesa.",
			Keywords:   "Arma, Marcial",
			ActionType: "Ação Padrão", Range: "Arma à distância",
			Target:  "Uma criatura",
			Attack:  "Destreza vs. CA",
			Hit:     "3[A] + mod Destreza + mod Sabedoria de dano. Esse ataque ignora cobertura e camuflagem.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Marca do Predador", Edition: "4e", ClassID: &id,
			Description: "Você designa um inimigo como sua presa preferencial, tornando impossível para ele escapar.",
			Keywords:   "Marcial",
			ActionType: "Ação Mínima", Range: "À distância 5",
			Target: "Um inimigo",
			Effect: "Até o final do encontro, o patrulheiro recebe +2 de bônus de poder em todas as jogadas de ataque e dano contra o alvo. Além disso, sempre que o alvo se mover, o patrulheiro pode se mover até sua velocidade como reação imediata.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Postura do Predador", Edition: "4e", ClassID: &id,
			Description: "Você assume uma postura de combate que maximiza sua mobilidade e eficácia com duas armas ou arco.",
			Keywords:   "Marcial, Postura",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura do Predador. Enquanto nessa postura: +1 de bônus de poder nas jogadas de ataque, pode ajustar 1 quadrado após cada ataque bem-sucedido e não provoca ataques de oportunidade ao usar poderes à distância.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Patrulheiro 4e: %d habilidades processadas", len(skills))
}

func seedBarbaroSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Bárbaro", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Bárbaro 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Golpe em Fúria", Edition: "4e", ClassID: &id,
			Description: "Você canaliza sua fúria primitiva num ataque devastador que cresce em poder conforme você avança de nível.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Força vs. CA",
			Hit:     "Nível 1: 3[A] + mod Força. Nível 5: 4[A]. Nível 9: 5[A]. Nível 15: 6[A]. Nível 19: 7[A]. Nível 20: 7[A]. Nível 25: 8[A]. Nível 29: 9[A] + mod Força.",
			Miss:    "Metade do dano.",
			Special: "Diário (Especial): pode ser usado duas vezes por dia. Condição: deve estar em fúria e gastar um poder de bárbaro com palavra-chave fúria não utilizado.",
			PowerType:      domain.PowerDaily, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Investida Célere", Edition: "4e", ClassID: &id,
			Description: "Conforme seu inimigo cai, você corre em direção à sua próxima vítima.",
			Keywords:   "Primitivo",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special:        "Gatilho: Um ataque do bárbaro reduz um inimigo a 0 PV ou menos.",
			Effect:         "O bárbaro realiza uma investida contra um inimigo.",
			PowerType:      domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Rugido de Triunfo", Edition: "4e", ClassID: &id,
			Description: "Seu uivo de vitória abala os inimigos no âmago, pois eles sabem que sua ânsia por sangue ainda não foi saciada.",
			Keywords:   "Medo, Primitivo",
			ActionType: "Ação Livre", Range: "Explosão contígua 5",
			Special:        "Gatilho: Um ataque do bárbaro reduz um inimigo a 0 PV ou menos.",
			Target:         "Os inimigos dentro da explosão",
			Effect:         "O alvo sofre -2 de penalidade em todas as defesas até o final do próximo turno.",
			PowerType:      domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe de Recuperação", Edition: "4e", ClassID: &id,
			Description: "Nada restaura sua vontade de lutar como um golpe bem dado no oponente. Com cada balanço esmagador, você tem mais vontade de pressioná-lo.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O bárbaro recebe PV temporários iguais ao seu modificador de Constituição. Se estiver em fúria, recebe PV temp iguais a 5 + seu modificador de Constituição.",
			Special:      "Condição: o bárbaro deve empunhar uma arma de combate corpo a corpo com as duas mãos.",
			LevelScaling: "Nível 11: 1[A] + 1d6 + mod Força. Nível 21: 2[A] + 2d6 + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Devastador", Edition: "4e", ClassID: &id,
			Description: "Você golpeia com um poder terrível, mais interessado com a força ofensiva que com o posicionamento defensivo.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano.",
			Effect:       "Até o começo do próximo turno do bárbaro, qualquer atacante recebe +2 de bônus nas jogadas de ataque contra ele. Se o personagem estiver em fúria, os atacantes não recebem esse bônus.",
			Special:      "Condição: o bárbaro deve empunhar uma arma de combate corpo a corpo com as duas mãos.",
			LevelScaling: "Nível 11: 1[A] + 1d8 + mod Força. Nível 21: 2[A] + 3d8 + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Insistente", Edition: "4e", ClassID: &id,
			Description: "Você não recua, não hesita — apenas pressiona o ataque sem parar até o inimigo cair.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Saltar o Caído", Edition: "4e", ClassID: &id,
			Description: "Você salta sobre um inimigo caído para desferir-lhe um golpe poderoso que o impede de se recuperar.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O alvo fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Sangria", Edition: "4e", ClassID: &id,
			Description: "Você inflige um corte profundo que continua sangrando, enfraquecendo gradualmente o inimigo.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O alvo sofre 5 de dano contínuo (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Golpe Uivante", Edition: "4e", ClassID: &id,
			Description: "Com um grito selvagem, você ataca e entra num estado de fúria primitiva que potencializa seus ataques seguintes.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "1[A] + mod Força de dano. O bárbaro entra em fúria até o final do encontro.",
			Effect: "Enquanto em fúria: recebe +2 de bônus de poder nas jogadas de ataque contra o alvo deste poder.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Fúria da Caçada de Sangue", Edition: "4e", ClassID: &id,
			Description: "O sangue do seu inimigo inflamou seus instintos primitivos de caça, tornando cada morte um combustível para sua fúria.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O bárbaro entra em fúria da caçada de sangue.",
			Effect: "Enquanto em fúria: sempre que um inimigo adjacente ao bárbaro for reduzido a 0 PV ou menos, o personagem pode se deslocar até a sua velocidade como ação livre.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Golpe Avalanche", Edition: "4e", ClassID: &id,
			Description: "Como uma avalanche imparável, cada golpe empurra o inimigo e cria espaço para você avançar ainda mais.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O bárbaro entra em fúria da avalanche.",
			Effect: "Enquanto em fúria: sempre que o bárbaro atingir um inimigo com um ataque, empurra esse inimigo 1 quadrado.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Fúria da Pantera Célere", Edition: "4e", ClassID: &id,
			Description: "Seus movimentos tornam-se imprevisíveis e furiosos como os de uma pantera selvagem.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O bárbaro entra em fúria da pantera célere.",
			Effect: "Enquanto em fúria: o bárbaro recebe +2 de bônus de poder na velocidade e pode ajustar 1 quadrado como ação livre uma vez por turno.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Rugido de Desafio", Edition: "4e", ClassID: &id,
			Description: "Seu rugido desafiador chama todos os inimigos próximos para lutar com você, marcando-os como suas presas.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Enquanto em fúria, até o final do encontro, todos os inimigos a até 3 quadrados ficam marcados pelo bárbaro.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Espírito Totem do Urso", Edition: "4e", ClassID: &id,
			Description: "O espírito do urso lhe concede a capacidade de se recuperar de ferimentos durante a batalha.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O bárbaro recebe regeneração 2 enquanto estiver em fúria e com menos de metade dos seus PV máximos.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Grande Trespassar", Edition: "4e", ClassID: &id,
			Description: "Seu golpe avassalador atravessa o alvo e atinge o inimigo por trás dele.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. Um inimigo adjacente ao alvo (exceto o alvo inicial) sofre dano igual ao modificador de Força do bárbaro.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe do Cauda-de-Maça", Edition: "4e", ClassID: &id,
			Description: "Você balança sua arma num arco devastador que derruba tudo que estiver na frente.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. CA",
			Hit:    "1[A] + mod Força de dano. Se o bárbaro estiver em fúria, o alvo fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Brado do Terror", Edition: "4e", ClassID: &id,
			Description: "Você solta um brado aterrorizante que faz os inimigos recuarem de medo.",
			Keywords:   "Medo, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. Vontade",
			Hit:    "Os alvos ficam amedrontados e recuam usando seu deslocamento. Se o bárbaro estiver em fúria, os alvos ficam atordoados até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Fúria do Dragonete Raivoso", Edition: "4e", ClassID: &id,
			Description: "A ferocidade de um dragonete furioso toma conta de você, tornando seus ataques imparáveis.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "3[A] + mod Força de dano. O bárbaro entra em fúria do dragonete raivoso.",
			Effect: "Enquanto em fúria: os ataques corpo a corpo do bárbaro causam dano adicional igual ao seu modificador de Constituição.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Quebrar o Crânio", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe violento na cabeça do inimigo, deixando-o atordoado.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. Fortitude",
			Hit:    "2[A] + mod Força de dano. O alvo fica atordoado até o final do próximo turno do bárbaro.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Turbilhão de Lâminas", Edition: "4e", ClassID: &id,
			Description: "Você gira como um turbilhão de aço, atingindo todos os inimigos ao redor.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. CA",
			Hit:    "1[A] + mod Força de dano. Se o bárbaro estiver em fúria, causa dano adicional igual ao seu modificador de Constituição.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Fúria do Coração Flamejante", Edition: "4e", ClassID: &id,
			Description: "Sua fúria se torna tão intensa que chamas ardentes emana do seu corpo, queimando tudo ao redor.",
			Keywords:   "Arma, Flamejante, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano flamejante. O bárbaro entra em fúria do coração flamejante.",
			Effect: "Enquanto em fúria: os inimigos que iniciarem o turno adjacentes ao bárbaro sofrem 5 de dano flamejante.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Estimular o Ciclo", Edition: "4e", ClassID: &id,
			Description: "Da mesma forma que no mundo natural a morte leva a nova vida, matar seus inimigos estimula você a uma nova ação.",
			Keywords:   "Primitivo",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special: "Gatilho: O bárbaro reduz um inimigo a 0 PV ou menos durante o seu turno.",
			Effect:  "O personagem ganha uma ação padrão.",
			PowerType: domain.PowerDaily, Level: 6,
		},
		{
			Name: "Resistência Primitiva", Edition: "4e", ClassID: &id,
			Description: "Você permanece intocado pelas energias mágicas de seus inimigos.",
			Keywords:   "Postura, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você assume a Postura de Resistência Primitiva. Enquanto nessa postura: o bárbaro adquire resistência 10 a um tipo de dano à sua escolha: ácido, congelante, flamejante, elétrico ou trovejante.",
			PowerType: domain.PowerDaily, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe Vigoroso", Edition: "4e", ClassID: &id,
			Description: "Você desfere um golpe com tamanha força que recupera seu vigor enquanto derruba o inimigo.",
			Keywords:   "Arma, Cura, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. O bárbaro pode gastar um pulso de cura.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Queda do Trovão", Edition: "4e", ClassID: &id,
			Description: "Seu golpe cai com a força de um trovão, esmagando o inimigo e fazendo o solo tremer.",
			Keywords:   "Arma, Fúria, Primitivo, Trovejante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. Fortitude",
			Hit:    "3[A] + mod Força de dano trovejante. O alvo fica atordoado (TR encerra). O bárbaro entra em fúria do trovão.",
			Effect: "Enquanto em fúria: os ataques corpo a corpo do bárbaro causam dano trovejante adicional igual ao seu modificador de Constituição.",
			Miss:   "Metade do dano e o alvo fica atordoado até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Tempestade de Lâminas", Edition: "4e", ClassID: &id,
			Description: "Você libera uma tempestade de golpes furiosos que varrem todos os inimigos ao redor.",
			Keywords:   "Arma, Fúria, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. CA",
			Hit:    "2[A] + mod Força de dano. Se o bárbaro estiver em fúria, o alvo fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Grande Pisão", Edition: "4e", ClassID: &id,
			Description: "Você bate o pé no chão e energias primitivas fluem através de você, afundando o solo devido a tanto poder.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Explosão contígua 5",
			Effect: "Os quadrados dentro da explosão se tornam terreno acidentado até o final do próximo turno do bárbaro.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Alimentar o Fogo", Edition: "4e", ClassID: &id,
			Description: "Você baixa a arma, permitindo que o oponente acerte um golpe fácil, mas a dor só alimenta sua fúria e dá força a seus ataques.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Os inimigos adjacentes ao bárbaro podem realizar um ataque de oportunidade contra ele. Até o final do próximo turno, o bárbaro recebe +2 de bônus de poder nas jogadas de ataque para cada inimigo que realizou um ataque de oportunidade contra ele.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Bárbaro 4e: %d habilidades processadas", len(skills))
}

func seedDruidaSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Druida", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Druida 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Forma Selvagem", Edition: "4e", ClassID: &id,
			Description: "Você assume um dos aspectos da Fera Primitiva ou retorna à sua forma humanoide.",
			Keywords:   "Metamorfose, Primitivo",
			ActionType: "Ação Mínima (Especial)", Range: "Pessoal",
			Effect:         "O druida altera sua forma para a animal ou humanoide. Enquanto na forma animal: não pode usar poderes de ataque sem a palavra-chave Forma Animal, utilitários ou poderes de talentos sem a palavra-chave animal, embora possa sustentar tais poderes. O druida pode usar este poder uma vez por rodada.",
			Special:        "Escolha uma forma específica. A forma animal é do mesmo tamanho do druida. O equipamento se torna parte da forma animal, mas o personagem continua recebendo os benefícios dos equipamentos que estiver usando.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Acometida", Edition: "4e", ClassID: &id,
			Description: "Com um salto repentino, como o de um tigre, seu inimigo baixa a guarda, surpreso.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano. O alvo concede vantagem de combate à próxima criatura que o atacar antes do final do próximo turno do druida.",
			Special:      "Pode ser usado como ataque básico corpo a corpo.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Chamado da Fera", Edition: "4e", ClassID: &id,
			Description: "Ao invocar o lado selvagem inerente de cada criatura, você compele seus inimigos a lutar sem tática ou planejamento.",
			Keywords:   "Encanto, Implemento, Primitivo, Psíquico",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target:       "Os inimigos dentro da explosão",
			Attack:       "Sabedoria vs. Vontade",
			Hit:          "O alvo não pode adquirir vantagem de combate até o final do próximo turno do druida. Além disso, em seu próximo turno, o alvo sofre dano psíquico igual a 5 + o modificador de Sabedoria sempre que realizar um ataque que não inclua o aliado mais próximo do druida.",
			LevelScaling: "Nível 21: 10 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Chicote de Espinhos", Edition: "4e", ClassID: &id,
			Description: "Vinhas espinhosas brotam da madeira de seu implemento, açoitando sua presa e trazendo-a para mais perto.",
			Keywords:   "Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Fortitude",
			Hit:          "1d8 + mod Sabedoria de dano. O alvo é puxado 2 quadrados.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Dilaceramento Selvagem", Edition: "4e", ClassID: &id,
			Description: "Você se prepara para o golpe final, rasgando o inimigo com suas garras afiadas.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano. O alvo é conduzido 1 quadrado.",
			Special:      "Pode ser usado como ataque básico corpo a corpo.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Garras Aderentes", Edition: "4e", ClassID: &id,
			Description: "Dilacerando, seu inimigo é incapaz de escapar do seu próximo ataque.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano. O alvo fica lento até o final do próximo turno do druida.",
			Special:      "Pode ser usado como ataque básico corpo a corpo.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Semente Flamejante", Edition: "4e", ClassID: &id,
			Description: "Você arremessa uma semente repleta de energia primitiva contra os adversários, que atinge o solo e explode num claro flamejante.",
			Keywords:   "Flamejante, Implemento, Primitivo, Zona",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d6 de dano flamejante. Os quadrados adjacentes ao alvo se transformam numa zona flamejante até o final do próximo turno. Inimigos que ingressarem ou começarem seu turno dentro da zona sofrem dano flamejante igual ao modificador de Sabedoria.",
			LevelScaling: "Nível 21: 2d6 de dano flamejante.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Tempestade de Espigos", Edition: "4e", ClassID: &id,
			Description: "Você atinge o inimigo com um relâmpago, deixando o ar ao redor dele carregado de eletricidade. Ele deve se mover ou será atingido novamente.",
			Keywords:   "Elétrico, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano elétrico. Se o alvo não se mover pelo menos 2 quadrados em seu próximo turno, ele sofre dano elétrico adicional igual ao modificador de Sabedoria do druida.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Vento Resfriante", Edition: "4e", ClassID: &id,
			Description: "Uma rajada de vento resfriante abate seus inimigos, separando-os e afastando-os.",
			Keywords:   "Congelante, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target:       "As criaturas dentro da explosão",
			Attack:       "Sabedoria vs. Fortitude",
			Hit:          "1d6 de dano congelante. O alvo é conduzido 1 quadrado.",
			LevelScaling: "Nível 21: 2d6 de dano congelante.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Lampejo Álgido", Edition: "4e", ClassID: &id,
			Description: "Você atinge o inimigo com um ataque congelante que o paralisa onde quer que ele esteja.",
			Keywords:   "Congelante, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. Fortitude",
			Hit:     "1d6 + mod Sabedoria de dano congelante. O alvo fica imobilizado até o final do próximo turno do druida.",
			Special: "Guardião Primitivo: o ataque causa dano adicional igual ao modificador de Constituição.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Mordida Brusca", Edition: "4e", ClassID: &id,
			Description: "Usando rapidez e astúcia, você morde seus inimigos enquanto se esquiva para evitar contra-ataques.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target:  "Uma ou duas criaturas",
			Attack:  "Sabedoria vs. Reflexos",
			Hit:     "2[A] + mod Sabedoria de dano. Se obtiver sucesso em pelo menos um ataque, o druida ajusta 2 quadrados.",
			Special: "Predador Primitivo: o número de quadrados que pode ajustar é igual ao modificador de Destreza.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Separar o Rebanho", Edition: "4e", ClassID: &id,
			Description: "Seu olhar de fera assola a mente do alvo, criando um sentimento de derrota e trazendo-o para mais perto de suas garras.",
			Keywords:   "Encanto, Forma Animal, Implemento, Primitivo, Psíquico",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "2d8 + mod Sabedoria de dano psíquico. O alvo é puxado 3 quadrados.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma da Besta Agressora", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma de uma besta agressora e selvagem cujos ataques são mais letais.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d6 + mod Sabedoria de dano. O druida entra na Forma da Besta Agressora.",
			Effect: "Enquanto nessa forma: os poderes de forma animal do druida causam +1d6 de dano adicional. Além disso, pode usar a Acometida como ataque básico à distância com alcance de 10 quadrados.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Forma da Serpente Deslizante", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma de uma serpente ágil que pode deslizar através de qualquer abertura.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d6 + mod Sabedoria de dano. O alvo fica imobilizado (TR encerra). O druida entra na Forma da Serpente Deslizante.",
			Effect: "Enquanto nessa forma: o druida pode ignorar terreno difícil e se mover através de espaços de criaturas de qualquer tamanho sem penalidade.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Fogos da Vida", Edition: "4e", ClassID: &id,
			Description: "Uma explosão de luz dourada queima seus inimigos e fortalece seus aliados com energia curativa.",
			Keywords:   "Cura, Implemento, Primitivo, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d6 + mod Sabedoria de dano radiante.",
			Effect: "Cada aliado dentro da explosão recupera PV iguais ao modificador de Sabedoria do druida.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Pele de Pedra", Edition: "4e", ClassID: &id,
			Description: "Sua pele endurece como pedra para se proteger de ataques iminentes.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O druida recebe +2 de bônus de poder na CA até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Salto Selvagem", Edition: "4e", ClassID: &id,
			Description: "Você se move com o vigor e a agilidade de uma fera selvagem.",
			Keywords:   "Primitivo",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O druida salta até 3 quadrados, ignorando obstáculos e terreno difícil durante o movimento.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Frenesi Selvagem", Edition: "4e", ClassID: &id,
			Description: "Sua fera interior toma conta de você e você ataca com uma série de golpes ferozes.",
			Keywords:   "Forma Animal, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Toque corpo a corpo",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2[A] + mod Sabedoria de dano. O druida pode realizar esse ataque novamente contra o mesmo alvo ou um alvo diferente.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Chamado da Tempestade", Edition: "4e", ClassID: &id,
			Description: "Você invoca o poder da tempestade para atingir múltiplos inimigos com raios e ventos furiosos.",
			Keywords:   "Elétrico, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "1d10 + mod Sabedoria de dano elétrico. O alvo é empurrado 2 quadrados.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Urso Poderoso", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma maciça de um urso poderoso, tornando-se resistente a danos e devastador em combate.",
			Keywords:   "Forma Animal, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O druida entra na Forma do Urso Poderoso. Enquanto nessa forma: recebe +2 na CA, resistência 5 a todos os danos e seus ataques de forma animal causam dano adicional igual ao modificador de Constituição.",
			Special: "Uma vez durante o encontro: Ação Padrão — Toque corpo a corpo — Sabedoria vs. Fortitude — 2d10 + mod SAB de dano e alvo fica derrubado.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Presas da Natureza", Edition: "4e", ClassID: &id,
			Description: "Você invoca os espíritos da natureza para morder e prender o inimigo no lugar.",
			Keywords:   "Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d8 + mod Sabedoria de dano. O alvo fica imobilizado até o final do próximo turno do druida.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Enxame de Insetos", Edition: "4e", ClassID: &id,
			Description: "Você invoca um enxame de insetos espirituais que atormenta seus inimigos.",
			Keywords:   "Implemento, Primitivo, Zona",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "1d10 + mod Sabedoria de dano. O alvo fica lento até o final do próximo turno.",
			Effect: "A zona persiste até o final do próximo turno. Criaturas que iniciarem o turno dentro recebem -2 em todas as jogadas de ataque.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Lobo Caçador", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma ágil de um lobo caçador, tornando-se um predador implacável.",
			Keywords:   "Forma Animal, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O druida entra na Forma do Lobo Caçador. Enquanto nessa forma: recebe +2 de bônus de poder na velocidade e +2 de bônus nas jogadas de ataque contra inimigos adjacentes a um aliado do druida.",
			Special: "Uma vez: Ação Padrão — Toque corpo a corpo — Sabedoria vs. Reflexos — 2d8 + mod SAB de dano e alvo derrubado. O druida ajusta até 3 quadrados antes do ataque.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Cura das Ervas", Edition: "4e", ClassID: &id,
			Description: "Você canaliza o poder curativo da natureza para restaurar a saúde de um aliado.",
			Keywords:   "Cura, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O alvo pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Sabedoria do druida.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Comunhão com a Natureza", Edition: "4e", ClassID: &id,
			Description: "Você se funde brevemente com a natureza ao redor, tornando-se difícil de localizar.",
			Keywords:   "Primitivo",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O druida fica oculto até o final do próximo turno se houver vegetação, rochas ou terreno natural no quadrado em que está.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Rugido da Tempestade", Edition: "4e", ClassID: &id,
			Description: "Um rugido tempestuoso emana de você, atordoando os inimigos com seu poder primitivo.",
			Keywords:   "Implemento, Primitivo, Trovejante",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d8 + mod Sabedoria de dano trovejante. O alvo fica atordoado até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Gavião das Tempestades", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma de um gavião colossal que domina os céus e comanda a tempestade.",
			Keywords:   "Elétrico, Forma Animal, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O druida entra na Forma do Gavião das Tempestades. Enquanto nessa forma: recebe velocidade de voo 6 e seus ataques de forma animal causam dano elétrico adicional igual ao modificador de Sabedoria.",
			Special: "Uma vez: Ação Padrão — Rajada contígua 5 — Sabedoria vs. Reflexos — 2d10 + mod SAB de dano elétrico. Os alvos ficam cegos até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Fúria da Natureza", Edition: "4e", ClassID: &id,
			Description: "A própria natureza responde à sua fúria, desencadeando um ataque avassalador contra os inimigos.",
			Keywords:   "Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d10 + mod Sabedoria de dano. O alvo fica derrubado e lento (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Regeneração Natural", Edition: "4e", ClassID: &id,
			Description: "A força curativa da natureza flui continuamente através de você durante a batalha.",
			Keywords:   "Cura, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do encontro, o druida recebe regeneração 3. Enquanto na forma animal, a regeneração aumenta para 5.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Invocar o Espírito da Floresta", Edition: "4e", ClassID: &id,
			Description: "Você invoca um poderoso espírito da floresta para auxiliá-lo no campo de batalha.",
			Keywords:   "Conjuração, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "Um espírito da floresta aparece num quadrado desocupado. O espírito tem CA 20, Fortitude 18, Reflexos 16, Vontade 18 e PV iguais ao modificador de Sabedoria do druida × 5. O espírito pode realizar ataques básicos e, no início de cada turno do druida, pode curar um aliado adjacente de PV iguais ao modificador de Sabedoria.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Druida 4e: %d habilidades processadas", len(skills))
}

func seedFeiticeiroSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Feiticeiro", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Feiticeiro 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Algidez Dracônica", Edition: "4e", ClassID: &id,
			Description: "Com uma lufada de vento, seu inimigo é golpeado com força e jogado para trás.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Carisma vs. Fortitude",
			Hit:          "1d8 + mod Carisma de dano congelante. O alvo é empurrado 1 quadrado.",
			Special:      "Pode ser usado como ataque básico à distância.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Caminhar da Tempestade", Edition: "4e", ClassID: &id,
			Description: "Seus passos trovejantes abalam o inimigo.",
			Keywords:   "Arcano, Implemento, Trovejante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Carisma vs. Fortitude",
			Hit:          "1d8 + mod Carisma de dano trovejante.",
			Special:      "O feiticeiro ajusta 1 quadrado antes ou depois do ataque.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Jorro Ardente", Edition: "4e", ClassID: &id,
			Description: "Ao movimentar seu braço num arco, você conjura fogo líquido contra seu inimigo.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target:       "As criaturas dentro da rajada",
			Attack:       "Carisma vs. Reflexos",
			Hit:          "1d8 + mod Carisma de dano flamejante.",
			Special:      "Magia Dracônica: o próximo inimigo que atingir o feiticeiro com um ataque corpo a corpo sofre dano flamejante igual ao modificador de Força do feiticeiro.",
			LevelScaling: "Nível 21: 2d8 + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Orbe Ácido", Edition: "4e", ClassID: &id,
			Description: "Um globo de ácido aparece em sua mão e rapidamente você o arremessa contra o alvo.",
			Keywords:   "Ácido, Arcano, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 20",
			Target:       "Uma criatura",
			Attack:       "Carisma vs. Reflexos",
			Hit:          "1d10 + mod Carisma de dano ácido.",
			Special:      "Pode ser usado como ataque básico à distância.",
			LevelScaling: "Nível 21: 2d10 + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Raio do Caos", Edition: "4e", ClassID: &id,
			Description: "Um raio de luzes multicoloridas salta de sua mão e segue gritando em direção à cabeça do inimigo.",
			Keywords:   "Arcano, Implemento, Psíquico",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Carisma vs. Vontade",
			Hit:          "1d10 + mod Carisma de dano psíquico.",
			Special:      "Magia Selvagem: se obter um número par na jogada de ataque, realize um ataque secundário contra uma criatura a até 5 quadrados do alvo — Carisma vs. Vontade — 1d6 de dano psíquico.",
			LevelScaling: "Nível 21: 2d10 + mod Carisma.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Elo Álgido", Edition: "4e", ClassID: &id,
			Description: "Um gelo rangente envolve o inimigo, retardando seus movimentos.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Carisma vs. Fortitude",
			Hit:    "3d6 + mod Carisma de dano congelante. O alvo sofre -2 de penalidade na defesa de Reflexos até o final do próximo turno do feiticeiro.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Explosão Molestadora", Edition: "4e", ClassID: &id,
			Description: "Criando uma erupção de energia mental, você atinge o oponente, debilitando-o.",
			Keywords:   "Arcano, Implemento, Psíquico",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target:  "Uma ou duas criaturas dentro da explosão",
			Attack:  "Carisma vs. Vontade",
			Hit:     "1d10 + mod Carisma de dano psíquico. O alvo é empurrado um número de quadrados igual ao modificador de Destreza.",
			Special: "Magia Selvagem: se acertar, conduz o alvo em vez de empurrá-lo.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Pancada Trovejante", Edition: "4e", ClassID: &id,
			Description: "Atingido por uma onda sonora, seu inimigo é arremessado para trás.",
			Keywords:   "Arcano, Implemento, Trovejante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Carisma vs. Fortitude",
			Hit:    "2d10 + mod Carisma de dano trovejante. O alvo é empurrado 3 quadrados.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Pira Explosiva", Edition: "4e", ClassID: &id,
			Description: "Seu alvo percebe tarde demais que está no centro de uma verdadeira conflagração criada por você.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Carisma vs. Reflexos",
			Hit:    "2d8 + mod Carisma de dano flamejante. Até o começo do próximo turno do feiticeiro, qualquer inimigo que ingressar ou começar seu turno num quadrado adjacente ao alvo sofre 1d6 de dano flamejante.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Sopro da Tempestade", Edition: "4e", ClassID: &id,
			Description: "Ao expirar uma névoa causticante, você queima seus inimigos, desorientando-os.",
			Keywords:   "Ácido, Arcano, Implemento",
			ActionType: "Ação Padrão", Range: "Rajada contígua 3",
			Target:  "As criaturas dentro da rajada",
			Attack:  "Carisma vs. Reflexos",
			Hit:     "2d6 + mod Carisma de dano ácido. O alvo não pode obter vantagem de combate contra qualquer criatura até o final do próximo turno do feiticeiro.",
			Special: "Magia Dracônica: o feiticeiro adquire ocultação até o final do seu próximo turno.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Orbe Cromático", Edition: "4e", ClassID: &id,
			Description: "Um orbe de energia arcana e de cores alternantes é lançado contra o alvo, liberando uma energia contida de efeito variável.",
			Keywords:   "Arcano, Implemento, Variável",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:  "Uma criatura",
			Attack:  "Carisma vs. Reflexos",
			Hit:     "3d10 + mod Carisma de dano. Jogue 1d6: 1-Amarelo (Radiante, alvo pasmo TR enc); 2-Vermelho (Flamejante, adj sofrem dano flamejante = mod DES); 3-Verde (Venenoso e 5 dano venoso contínuo TR enc); 4-Turquesa (Elétrico, alvo conduzido = mod DES); 5-Azul (Congelante, alvo imobilizado TR enc); 6-Violeta (Psíquico, -2 na CA do alvo TR enc).",
			Miss:    "1d10 de dano. Jogue 1d6 para determinar o tipo e o efeito do ataque.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Raio da Presa Dracônica", Edition: "4e", ClassID: &id,
			Description: "Presas venenosas se lançam contra os inimigos, rasgando a carne deles e envenenando seus corpos.",
			Keywords:   "Arcano, Implemento, Venenoso",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma ou duas criaturas",
			Attack: "Carisma vs. Fortitude",
			Hit:    "2d8 + mod Carisma de dano e 5 de dano venenoso contínuo (TR encerra).",
			Miss:   "2d8 + mod Carisma de dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Asas Arcanas", Edition: "4e", ClassID: &id,
			Description: "Você manifesta asas de energia arcana e se eleva brevemente pelos ares.",
			Keywords:   "Arcano",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O feiticeiro voa até 4 quadrados.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Esquiva Arcana", Edition: "4e", ClassID: &id,
			Description: "Você canaliza o caos arcano para teleportar-se instantaneamente para fora do caminho de um ataque.",
			Keywords:   "Arcano",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O feiticeiro é alvo de um ataque.",
			Effect:  "O feiticeiro se teletransporta até 3 quadrados e recebe +2 nas defesas contra esse ataque.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Raio Deslumbrante", Edition: "4e", ClassID: &id,
			Description: "Um raio de luz cega seu inimigo, deixando-o vulnerável a ataques.",
			Keywords:   "Arcano, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Carisma vs. Reflexos",
			Hit:    "2d8 + mod Carisma de dano radiante. O alvo fica cego até o final do próximo turno do feiticeiro.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Corrente Elétrica", Edition: "4e", ClassID: &id,
			Description: "Uma corrente de eletricidade salta entre múltiplos inimigos, atingindo vários de uma vez.",
			Keywords:   "Arcano, Elétrico, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura (primário)",
			Attack: "Carisma vs. Reflexos",
			Hit:    "2d6 + mod Carisma de dano elétrico. Realize um ataque secundário contra um inimigo a até 5 quadrados do alvo primário: Carisma vs. Reflexos — 1d6 + mod CAR de dano elétrico.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Raio Deslumbrante", Edition: "4e", ClassID: &id,
			Description: "Um potente raio de luz solar concentrada cega permanentemente o inimigo durante o encontro.",
			Keywords:   "Arcano, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Carisma vs. Vontade",
			Hit:    "3d8 + mod Carisma de dano radiante. O alvo fica cego até o final do encontro.",
			Miss:   "Metade do dano. O alvo fica cego até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Tempestade de Gelo", Edition: "4e", ClassID: &id,
			Description: "Uma tempestade de fragmentos de gelo aflige seus inimigos, desacelerando seus movimentos.",
			Keywords:   "Arcano, Congelante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "As criaturas dentro da explosão",
			Attack: "Carisma vs. Fortitude",
			Hit:    "2d6 + mod Carisma de dano congelante. O alvo fica lento (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Explosão Ígnea", Edition: "4e", ClassID: &id,
			Description: "Uma enorme explosão de fogo devasta tudo numa grande área ao redor do alvo.",
			Keywords:   "Arcano, Flamejante, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "As criaturas dentro da explosão",
			Attack: "Carisma vs. Reflexos",
			Hit:    "2d8 + mod Carisma de dano flamejante.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Fogo Trovejante", Edition: "4e", ClassID: &id,
			Description: "Uma combinação devastadora de fogo e trovão destrói os inimigos e paralisa os sobreviventes.",
			Keywords:   "Arcano, Flamejante, Implemento, Trovejante",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "As criaturas dentro da explosão",
			Attack: "Carisma vs. Reflexos",
			Hit:    "3d6 + mod Carisma de dano flamejante e trovejante. O alvo fica atordoado até o final do próximo turno.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Escudo de Chamas", Edition: "4e", ClassID: &id,
			Description: "Você envolve seu corpo em chamas que queimam quem ousar atacar você.",
			Keywords:   "Arcano, Flamejante",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do próximo turno do feiticeiro, qualquer inimigo que atingir o feiticeiro com um ataque corpo a corpo sofre 1d6 + mod Carisma de dano flamejante.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Teleporte Arcano", Edition: "4e", ClassID: &id,
			Description: "Com um estalo de dedos, você se teletransporta a uma distância considerável.",
			Keywords:   "Arcano, Teletransporte",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O feiticeiro se teletransporta até 6 quadrados.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Explosão de Caos", Edition: "4e", ClassID: &id,
			Description: "Uma explosão de magia caótica atinge aleatoriamente vários inimigos ao redor.",
			Keywords:   "Arcano, Implemento, Variável",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target:  "Os inimigos dentro da explosão",
			Attack:  "Carisma vs. Reflexos",
			Hit:     "2d8 + mod Carisma de dano. Jogue 1d4: 1-Flamejante; 2-Congelante; 3-Elétrico; 4-Ácido.",
			Special: "Magia Selvagem: se acertar, o feiticeiro se teletransporta até 3 quadrados.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Sopro Dracônico", Edition: "4e", ClassID: &id,
			Description: "Você libera um devastador sopro de energia dracônica que varre tudo à sua frente.",
			Keywords:   "Arcano, Implemento, Variável",
			ActionType: "Ação Padrão", Range: "Explosão contígua 5",
			Target:  "Os inimigos dentro da explosão",
			Attack:  "Carisma vs. Reflexos",
			Hit:     "4d6 + mod Carisma de dano. Jogue 1d4 para o tipo: 1-Flamejante; 2-Congelante; 3-Elétrico; 4-Ácido.",
			Miss:    "Metade do dano.",
			Special: "Magia Dracônica: o feiticeiro escolhe o tipo de dano em vez de jogar.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Devastação Arcana", Edition: "4e", ClassID: &id,
			Description: "Você libera toda a energia arcana acumulada em um único ataque devastador.",
			Keywords:   "Arcano, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Carisma vs. Vontade",
			Hit:    "4d8 + mod Carisma de dano. O alvo fica atordoado até o final do próximo turno. Todos os inimigos adjacentes ao alvo sofrem 2d8 + mod CAR de dano do mesmo tipo.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Absorver Magia", Edition: "4e", ClassID: &id,
			Description: "Você absorve a energia de um ataque mágico e a converte em poder arcano.",
			Keywords:   "Arcano",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O feiticeiro é alvo de um poder de ataque arcano ou divino.",
			Effect:  "O feiticeiro recebe resistência 15 contra esse ataque e recupera a utilização de um poder por encontro que já tenha usado.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Forma Elemental", Edition: "4e", ClassID: &id,
			Description: "Você transforma seu corpo em energia elemental pura, tornando-se temporariamente invulnerável a certos danos.",
			Keywords:   "Arcano, Variável",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do encontro, o feiticeiro recebe resistência 10 a um tipo de dano à sua escolha (flamejante, congelante, elétrico ou ácido) e seu tipo de dano em todos os poderes muda para o tipo escolhido.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Feiticeiro 4e: %d habilidades processadas", len(skills))
}

func seedGuardiaoSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Guardião", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Guardião 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Aperto do Guardião", Edition: "4e", ClassID: &id,
			Description: "Vinhas espectrais se fixam ao inimigo que atacou seu aliado, impedindo seus movimentos e punindo sua ousadia.",
			Keywords:   "Primitivo",
			ActionType: "Reação Imediata", Range: "Explosão contígua 5",
			Target:         "O inimigo que ativou o gatilho",
			Special:        "Gatilho: Um inimigo marcado pelo guardião a até 5 quadrados realiza um ataque que não o inclui como alvo.",
			Effect:         "O alvo é conduzido 1 quadrado, fica lento e não pode ajustar até o final do seu próximo turno.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Fúria do Guardião", Edition: "4e", ClassID: &id,
			Description: "Canalizando a fúria da natureza, você ataca um inimigo que atacou seu aliado, minando suas defesas.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Interrupção Imediata", Range: "Arma corpo a corpo",
			Target:         "O inimigo que ativou o gatilho",
			Special:        "Gatilho: Um inimigo marcado realiza um ataque que não inclui o guardião como alvo.",
			Attack:         "Força vs. Fortitude",
			Hit:            "1[A] + mod Força de dano. O alvo concede vantagem de combate ao guardião e aos aliados dele até o final do próximo turno do personagem.",
			LevelScaling:   "Nível 21: 2[A] + mod Força.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Força da Pedra", Edition: "4e", ClassID: &id,
			Description: "Extraindo o poder da terra, você se revigora contra os ataques enquanto golpeia o inimigo com força implacável.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O guardião recebe PV temporários iguais ao seu modificador de Constituição.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe de Espinhos", Edition: "4e", ClassID: &id,
			Description: "Espinhos espectrais brotam de sua arma e se fixam no inimigo, trazendo-o para perto de você.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Corpo a corpo 2",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O alvo é puxado 1 quadrado.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe do Escudo de Terra", Edition: "4e", ClassID: &id,
			Description: "Através do solo, o poder das rochas concede a você o peso necessário para fortalecer sua próxima defesa.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O guardião recebe +1 de bônus de poder na CA até o final do seu próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Peso da Terra", Edition: "4e", ClassID: &id,
			Description: "Seu ataque envia a energia da terra contra o inimigo, retardando seus movimentos.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. O alvo fica lento até o final do próximo turno do guardião.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Assalto do Carneiro Trovejante", Edition: "4e", ClassID: &id,
			Description: "Ao atingir seu inimigo, você canaliza o espírito do carneiro do trovão para atacar os aliados do inimigo com uma rajada sônica.",
			Keywords:   "Arma, Primitivo, Trovejante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura (primário)",
			Attack:  "Força vs. CA",
			Hit:     "1[A] + mod Força de dano trovejante. Realize um ataque secundário em rajada contígua 3: Força vs. Fortitude — 1d6 de dano trovejante.",
			Special: "Força Térrea: o alvo primário é empurrado um número de quadrados igual ao modificador de Constituição do guardião.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Espigos da Terra", Edition: "4e", ClassID: &id,
			Description: "A própria terra se ergue em resposta ao seu ataque, arremessando pontas afiadas de madeira e pedra contra os inimigos próximos.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. CA",
			Hit:    "1[A] + mod Força de dano. Até o final do próximo turno do guardião, o espaço do alvo e os quadrados adjacentes ficam repletos de espigos. Qualquer inimigo que ingressar num quadrado com espigos sofre 5 de dano.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Frenesi do Sangue Selvagem", Edition: "4e", ClassID: &id,
			Description: "O poder primitivo corre em suas veias e você é tomado por um frenesi, desferindo dois ataques poderosos.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura (primário)",
			Attack:  "Força vs. CA",
			Hit:     "1[A] + mod Força de dano.",
			Special: "Sangue Selvagem: o ataque causa dano adicional igual ao modificador de Sabedoria do guardião.",
			Effect:  "O guardião realiza esse ataque novamente contra o mesmo alvo ou um alvo diferente.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Terra Faminta", Edition: "4e", ClassID: &id,
			Description: "Batendo no solo com sua arma, você faz emergir a energia primitiva que ataca os inimigos por bem debaixo dos pés deles.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. Fortitude",
			Effect: "Até o final do próximo turno do guardião, os quadrados dentro da explosão se tornam terreno acidentado para os inimigos do personagem.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma da Pantera Incansável", Edition: "4e", ClassID: &id,
			Description: "Você adquire as patas afiadas e a graça caçadora de uma pantera. Sempre que desejar, pode desferir um ataque ágil que persiste no inimigo.",
			Keywords:   "Metamorfose, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O personagem assume a Forma da Pantera Incansável. Enquanto nessa forma: recebe +2 de bônus nas jogadas de ataque contra inimigos marcados, +1 de bônus nos ataques de oportunidade e pode ajustar 2 quadrados com uma ação de movimento.",
			Special: "Uma vez durante o encontro, o guardião pode usar: Ação Padrão — Arma corpo a corpo — Força vs. Reflexos — 2[A] + mod Força de dano e 5 de dano contínuo (TR encerra). Fracasso: metade do dano e 2 de dano contínuo.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Forma do Arauto do Inverno", Edition: "4e", ClassID: &id,
			Description: "Gelo tão forte quanto o aço se forma sobre sua armadura enquanto o arauto do inverno toma forma, tornando-o um bastião de defesa.",
			Keywords:   "Congelante, Metamorfose, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O personagem assume a Forma do Arauto do Inverno. Enquanto nessa forma: recebe +1 de bônus na CA, resistência 5 ao dano congelante e os quadrados a até 2 onde o personagem estiver são considerados terreno acidentado para os inimigos.",
			Special: "Uma vez durante o encontro, pode usar: Ação Padrão — Explosão contígua 1 — Força vs. CA — 1[A] + mod Força de dano congelante e alvo fica imobilizado (TR encerra). Fracasso: metade do dano e alvo imobilizado até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Marca do Guardião", Edition: "4e", ClassID: &id,
			Description: "Você marca um inimigo para garantir que ele não poderá atacar seus aliados sem consequências.",
			Keywords:   "Primitivo",
			ActionType: "Ação Livre", Range: "À distância 5",
			Target: "Um inimigo",
			Effect: "O guardião marca o alvo até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Escudo Primitivo", Edition: "4e", ClassID: &id,
			Description: "O poder primitivo da terra cria uma barreira momentânea ao redor do guardião quando ele é atacado.",
			Keywords:   "Primitivo",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O guardião é alvo de um ataque.",
			Effect:  "O guardião recebe +2 de bônus de poder na CA e Fortitude contra esse ataque.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Golpe das Raízes", Edition: "4e", ClassID: &id,
			Description: "Raízes espectrais brotam do chão e prendem o inimigo no lugar após seu golpe.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Força vs. CA",
			Hit:          "1[A] + mod Força de dano. Se o alvo tentar se mover voluntariamente até o final do próximo turno do guardião, ele fica imobilizado até o final do seu próximo turno.",
			LevelScaling: "Nível 21: 2[A] + mod Força.",
			PowerType: domain.PowerUnlimited, Level: 3,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Ira da Pedra", Edition: "4e", ClassID: &id,
			Description: "O poder da pedra concentra-se em seu golpe, criando uma explosão de fragmentos rochosos.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 1",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. Fortitude",
			Hit:    "1[A] + mod Força de dano. Os alvos ficam derrubados.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Chamado da Besta Guardiã", Edition: "4e", ClassID: &id,
			Description: "Você convoca o espírito de uma besta guardiã que protege seus aliados com sua presença feroz.",
			Keywords:   "Conjuração, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "O guardião conjura uma besta espiritual num quadrado desocupado. Inimigos adjacentes à besta ficam marcados pelo guardião. A besta persiste até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Touro Colossal", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma massiva de um touro colossal, tornando-se uma força imparável no campo de batalha.",
			Keywords:   "Metamorfose, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O guardião entra na Forma do Touro Colossal. Enquanto nessa forma: aumenta de tamanho para Grande, recebe +4 de bônus de poder na Fortitude e seus ataques empurram o alvo 1 quadrado.",
			Special: "Uma vez: Ação Padrão — Arma corpo a corpo — Força vs. Fortitude — 3[A] + mod Força de dano. O alvo é empurrado 3 quadrados e fica derrubado. Fracasso: metade do dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Barreira de Pedra", Edition: "4e", ClassID: &id,
			Description: "Você cria uma barreira temporária de pedra que protege seus aliados e bloqueia o avanço inimigo.",
			Keywords:   "Conjuração, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "Uma barreira de pedra de até 4 quadrados de comprimento surge no alcance. A barreira fornece cobertura superior e bloqueia o movimento. Persiste até o final do próximo turno do guardião.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Abraço da Terra", Edition: "4e", ClassID: &id,
			Description: "A terra se abre e abraça o inimigo, prendendo-o no lugar com força irresistível.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. Fortitude",
			Hit:    "2[A] + mod Força de dano. O alvo fica imobilizado (TR encerra) e sofre -2 em todas as defesas enquanto imobilizado.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Paquiderme Blindado", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma robusta de um paquiderme blindado, tornando-se praticamente indestrutível.",
			Keywords:   "Metamorfose, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O guardião entra na Forma do Paquiderme Blindado. Enquanto nessa forma: recebe +4 de bônus na CA, resistência 5 a todos os danos e não pode ser empurrado, conduzido ou derrubado.",
			Special: "Uma vez: Ação Padrão — Arma corpo a corpo — Força vs. Fortitude — 2[A] + mod Força de dano. O alvo é empurrado 4 quadrados. Cada criatura no caminho do alvo sofre 1d8 de dano.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Aura Protetora da Terra", Edition: "4e", ClassID: &id,
			Description: "Uma aura de energia primitiva irradia de você, protegendo todos os aliados próximos.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Você emite uma aura 2 até o final do encontro. Seus aliados dentro da aura recebem +1 de bônus de poder na CA.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Presença Intimidadora", Edition: "4e", ClassID: &id,
			Description: "Sua presença no campo de batalha amedronta os inimigos, fazendo-os hesitar antes de atacar seus aliados.",
			Keywords:   "Medo, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. Vontade",
			Hit:    "Os inimigos ficam amedrontados até o final do próximo turno do guardião.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Punho da Montanha", Edition: "4e", ClassID: &id,
			Description: "Seu golpe carrega o peso de uma montanha, esmagando o inimigo com força devastadora.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Força vs. Fortitude",
			Hit:    "3[A] + mod Força de dano. O alvo fica atordoado e derrubado até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Forma do Leviatã das Profundezas", Edition: "4e", ClassID: &id,
			Description: "Você assume a forma terrível de um leviatã das profundezas, tornando-se um monstro de combate.",
			Keywords:   "Metamorfose, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:  "O guardião entra na Forma do Leviatã. Enquanto nessa forma: recebe velocidade de nado 8, +2 na CA e resistência 10 a dano ácido e congelante.",
			Special: "Uma vez: Ação Padrão — Explosão contígua 2 — Força vs. Fortitude — 3[A] + mod Força de dano ácido. Os alvos ficam imobilizados (TR encerra).",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Tempestade Primitiva", Edition: "4e", ClassID: &id,
			Description: "Você desencadeia o poder bruto da natureza, varrendo tudo ao redor com uma tempestade de energia primitiva.",
			Keywords:   "Arma, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 2",
			Target: "Os inimigos dentro da explosão",
			Attack: "Força vs. Fortitude",
			Hit:    "2[A] + mod Força de dano. Os alvos ficam derrubados e lentos (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Couraça da Pedra Viva", Edition: "4e", ClassID: &id,
			Description: "A pedra viva se funde com sua pele, tornando-o praticamente impenetrável por um momento.",
			Keywords:   "Primitivo",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O guardião sofre dano.",
			Effect:  "Reduza o dano sofrido em 10. Além disso, o guardião recebe resistência 5 a todos os danos até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Chamado dos Ancestrais da Terra", Edition: "4e", ClassID: &id,
			Description: "Você invoca os espíritos ancestrais da terra para lutar ao seu lado.",
			Keywords:   "Conjuração, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "Dois guardiões espirituais de pedra surgem em quadrados desocupados. Cada guardião tem CA 22, Fortitude 20, Reflexos 16, Vontade 18. No início do seu turno, cada guardião adjacente a um inimigo marcado realiza um ataque: Força vs. CA — 1d10 + mod Força de dano. Os guardiões persistem até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Guardião 4e: %d habilidades processadas", len(skills))
}

func seedInvocadorSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Invocador", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Invocador 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Canalizar Divindade: Repreensão do Preservador", Edition: "4e", ClassID: &id,
			Description: "Você clama aos deuses para punir um inimigo que ousou ferir aqueles sob seus cuidados.",
			Keywords:   "Divino",
			ActionType: "Reação Imediata", Range: "Pessoal",
			Special:        "Gatilho: Um aliado do invocador a até 10 quadrados é atingido por um inimigo.",
			Effect:         "Até antes do final do seu próximo turno, o invocador recebe um bônus na próxima jogada de ataque contra o inimigo que ativou o gatilho, igual ao seu modificador de Inteligência.",
			PowerType:      domain.PowerEncounter, Level: 1,
			IsClassFeature: true,
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Estilhaços Aderentes", Edition: "4e", ClassID: &id,
			Description: "Você arremessa uma esfera mágica contra os inimigos. No impacto, ela se fragmenta em centenas de minúsculas lâminas radiantes que cortam os inimigos e atraem seus movimentos.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target:       "As criaturas dentro da explosão",
			Attack:       "Sabedoria vs. Fortitude",
			Hit:          "Modificador de Sabedoria de dano radiante. O alvo fica lento até o final do próximo turno do invocador.",
			LevelScaling: "Nível 21: 1d10 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Solar", Edition: "4e", ClassID: &id,
			Description: "Um raio de energia radiante se laça de suas mãos e banha um inimigo com uma luz causticante, forçando-o a se mover.",
			Keywords:   "Divino, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano radiante. O alvo é conduzido 1 quadrado.",
			Special:      "Pode ser usado como ataque básico à distância.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Luz Vingadora", Edition: "4e", ClassID: &id,
			Description: "Você castiga o inimigo com um orbe causticante de luz que queima com o fogo da vingança se seus aliados sofrerem algum mal.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Fortitude",
			Hit:          "1d10 + mod Sabedoria de dano radiante. Se um aliado sangrando do invocador estiver adjacente ao alvo, o ataque causa dano radiante adicional igual ao modificador de Constituição do invocador.",
			Special:      "Pode ser usado como ataque básico à distância.",
			LevelScaling: "Nível 21: 2d10 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Raios Divinos", Edition: "4e", ClassID: &id,
			Description: "Raios divinos de eletricidade atingem os inimigos com a fúria dos deuses.",
			Keywords:   "Divino, Elétrico, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma ou duas criaturas",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d6 + mod Sabedoria de dano elétrico.",
			LevelScaling: "Nível 21: 2d6 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Relâmpago de Vanguarda", Edition: "4e", ClassID: &id,
			Description: "Raios de luz divina arquejam de suas mãos, queimando uma área e permanecendo ali por algum tempo, prontos para retribuir um ataque inimigo.",
			Keywords:   "Divino, Elétrico, Implemento",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target:       "As criaturas dentro da explosão",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d6 + mod Sabedoria de dano elétrico. Sempre que o alvo realizar um ataque de oportunidade contra o invocador antes do final do próximo turno, o alvo sofre dano elétrico adicional igual ao modificador de Inteligência.",
			LevelScaling: "Nível 21: 2d6 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Lâminas de Fogo Astral", Edition: "4e", ClassID: &id,
			Description: "Lâminas cintilantes de energia radiante surgem e golpeiam os inimigos. Em seguida, elas se transformam em escudos espectrais para proteger seus aliados.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 1 a até 10 quadrados",
			Target:  "Os inimigos dentro da explosão",
			Attack:  "Sabedoria vs. Reflexos",
			Hit:     "1d6 + mod Sabedoria de dano radiante.",
			Effect:  "Os aliados do invocador dentro da explosão recebem +2 de bônus de poder na CA até o final do próximo turno.",
			Special: "Contrato da Preservação: o bônus é igual a 1 + o modificador de Inteligência do invocador.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Lança do Inquisidor", Edition: "4e", ClassID: &id,
			Description: "Uma lança de energia cintilante corta o ar em direção ao inimigo, queimando-o com o poder dos deuses e prendendo-o ao chão.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "1d10 + mod Sabedoria de dano radiante. O alvo fica imobilizado até o final do próximo turno do invocador.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Terror Astral", Edition: "4e", ClassID: &id,
			Description: "A energia astral flui através de você, transformando-o num farol de terror divino que assola os inimigos.",
			Keywords:   "Divino, Implemento, Medo, Psíquico",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "1d6 + mod Sabedoria de dano psíquico. O alvo é empurrado 2 quadrados.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Mão do Destino", Edition: "4e", ClassID: &id,
			Description: "A mão do destino se abate sobre o inimigo com força devastadora, imobilizando-o com energia divina.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d10 + mod Sabedoria de dano radiante. O alvo fica imobilizado (TR encerra).",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Juramento de Destruição", Edition: "4e", ClassID: &id,
			Description: "Você pronuncia um juramento divino de destruição sobre o inimigo, enfraquecendo permanentemente todas as suas defesas.",
			Keywords:   "Divino, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "O alvo sofre -4 de penalidade em todas as defesas até o final do encontro.",
			Miss:   "O alvo sofre -2 de penalidade em todas as defesas até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Punição Astral", Edition: "4e", ClassID: &id,
			Description: "Uma explosão de energia astral pune o inimigo por seus crimes contra os aliados do invocador.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "3d8 + mod Sabedoria de dano radiante. O alvo fica cego até o final do próximo turno. Se algum aliado do invocador estiver sangrando, o alvo também fica atordoado até o final do próximo turno.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Presença Divina", Edition: "4e", ClassID: &id,
			Description: "Você manifesta a presença divina ao seu redor, inspirando seus aliados a lutar com mais determinação.",
			Keywords:   "Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 5",
			Effect: "Os aliados dentro da explosão recebem +1 de bônus de poder nas jogadas de ataque até o final do encontro.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Muralha de Luz", Edition: "4e", ClassID: &id,
			Description: "Você cria uma barreira de luz divina que queima os inimigos que tentarem cruzá-la.",
			Keywords:   "Conjuração, Divino",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "Uma muralha de luz de até 4 quadrados surge no alcance. Inimigos que tentarem cruzá-la sofrem 1d8 + mod Sabedoria de dano radiante e ficam cegos até o final do próximo turno. A muralha persiste até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Barreira Astral", Edition: "4e", ClassID: &id,
			Description: "Uma barreira de energia astral surge entre você e seus inimigos, protegendo seus aliados.",
			Keywords:   "Conjuração, Divino, Implemento",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Effect: "Uma barreira astral de até 5 quadrados surge. Aliados que passarem pela barreira recebem +4 de bônus na CA até o final do próximo turno. Inimigos que passarem pela barreira sofrem 1d6 + mod Sabedoria de dano radiante.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Flagelo Divino", Edition: "4e", ClassID: &id,
			Description: "Você desce o flagelo divino sobre um grupo de inimigos, queimando-os com energia sagrada.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d6 + mod Sabedoria de dano radiante.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Guardião Astral", Edition: "4e", ClassID: &id,
			Description: "Você conjura um guardião astral que protege seus aliados e pune os inimigos que se aproximam.",
			Keywords:   "Conjuração, Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Effect: "Um guardião astral aparece num quadrado desocupado. Ele tem CA 20, Fort 18, Refl 14, Vont 18 e PV iguais ao seu nível. Inimigos adjacentes ao guardião ficam marcados. Ao início de cada turno do invocador, o guardião ataca um inimigo adjacente marcado: Sabedoria vs. CA — 1d8 + mod SAB de dano radiante. O guardião persiste até o final do encontro ou até ser destruído.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Chuva de Estrelas", Edition: "4e", ClassID: &id,
			Description: "Uma chuva de estrelas divinas cai sobre seus inimigos, queimando-os com energia celestial.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d8 + mod Sabedoria de dano radiante. O alvo fica cego até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Corrente da Condenação", Edition: "4e", ClassID: &id,
			Description: "Uma corrente de energia divina salta de inimigo em inimigo, condenando todos que tocar.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura (primário)",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d6 + mod Sabedoria de dano radiante. Realize ataques secundários contra até 2 inimigos a até 5 quadrados do alvo primário: Sabedoria vs. Reflexos — 1d6 + mod SAB de dano radiante cada.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Tempestade Divina", Edition: "4e", ClassID: &id,
			Description: "Você desencadeia o poder pleno da sua divindade, criando uma tempestade de energia sagrada que varre o campo de batalha.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "3d8 + mod Sabedoria de dano radiante. O alvo fica cego e atordoado até o final do próximo turno.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Escudo da Fé", Edition: "4e", ClassID: &id,
			Description: "Você invoca um escudo de fé divina que protege um aliado de ataques.",
			Keywords:   "Divino",
			ActionType: "Interrupção Imediata", Range: "À distância 5",
			Target:  "Um aliado",
			Special: "Gatilho: Um aliado é alvo de um ataque.",
			Effect:  "O aliado recebe +4 de bônus de poder na CA e Reflexos contra esse ataque.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Palavra de Cura Divina", Edition: "4e", ClassID: &id,
			Description: "Com uma palavra sagrada, você restaura a saúde de um aliado ferido.",
			Keywords:   "Cura, Divino",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Um aliado",
			Effect: "O alvo pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Sabedoria do invocador.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Cólera dos Céus", Edition: "4e", ClassID: &id,
			Description: "A cólera dos céus desce sobre seus inimigos numa explosão de energia divina devastadora.",
			Keywords:   "Divino, Elétrico, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 2 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d10 + mod Sabedoria de dano radiante e elétrico. O alvo fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Anjo da Destruição", Edition: "4e", ClassID: &id,
			Description: "Você conjura um anjo da destruição que varre o campo de batalha com poder devastador.",
			Keywords:   "Conjuração, Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Effect: "Um anjo da destruição aparece num quadrado desocupado dentro do alcance. Ele tem CA 24, Fort 22, Refl 18, Vont 20. No início de cada turno do invocador, o anjo pode se mover até 6 quadrados e realizar: Explosão contígua 1 — Sabedoria vs. CA — 2d8 + mod SAB de dano radiante. O anjo persiste até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Juízo Divino", Edition: "4e", ClassID: &id,
			Description: "O juízo divino cai sobre todos os inimigos ao mesmo tempo, punindo-os por sua impiedade.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "3d6 + mod Sabedoria de dano radiante. O alvo fica cego e imobilizado até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Bênção Divina", Edition: "4e", ClassID: &id,
			Description: "Você canaliza o poder pleno da sua divindade para abençoar todos os seus aliados.",
			Keywords:   "Divino",
			ActionType: "Ação Padrão", Range: "Explosão contígua 10",
			Effect: "Todos os aliados dentro da explosão recebem +2 de bônus de poder em todas as jogadas de ataque, dano e defesas até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Símbolo Divino", Edition: "4e", ClassID: &id,
			Description: "Você erge seu símbolo sagrado e o poder da sua divindade irradia em todas as direções.",
			Keywords:   "Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "Os inimigos ficam cegos e recebem -4 em todas as defesas até o final do próximo turno.",
			Effect: "Os aliados dentro da explosão recebem +2 de bônus de poder na CA e +2 de bônus nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Invocador 4e: %d habilidades processadas", len(skills))
}

func seedVingadorSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Vingador", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Vingador 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Juramento de Inimizade", Edition: "4e", ClassID: &id,
			Description: "Você pronuncia um juramento sagrado contra um inimigo, comprometendo-se a eliminá-lo acima de tudo.",
			Keywords:   "Divino",
			ActionType: "Ação Menor", Range: "À distância 10",
			Target:         "Um inimigo",
			Effect:         "O vingador escolhe um inimigo como alvo do Juramento de Inimizade. Sempre que o vingador realizar uma jogada de ataque contra esse alvo com um poder de vingador e o resultado for desfavorável, o vingador pode rolar novamente. Apenas um inimigo pode ser o alvo do Juramento de Inimizade por vez, e o efeito persiste até o final do encontro ou até o alvo ser reduzido a 0 PV.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Censura da Perseguição", Edition: "4e", ClassID: &id,
			Description: "Quando seu alvo jurado tenta escapar, você o persegue com velocidade sobrenatural.",
			Keywords:   "Divino",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special:        "Gatilho: O alvo do Juramento de Inimizade se move voluntariamente.",
			Effect:         "O vingador pode se mover um número de quadrados igual a 1 + seu modificador de Destreza em direção ao alvo como reação imediata.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
			RequiresChoice: true,
			ChoiceGroup:    "censura_vingador",
		},
		{
			Name: "Censura da Retribuição", Edition: "4e", ClassID: &id,
			Description: "Quando outros ousam atacar você enquanto persegue seu alvo, eles pagam um preço.",
			Keywords:   "Divino, Radiante",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special:        "Gatilho: Uma criatura diferente do alvo do Juramento de Inimizade ataca o vingador.",
			Effect:         "O alvo do Juramento de Inimizade sofre dano radiante igual a 5 + o modificador de Inteligência do vingador.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
			RequiresChoice: true,
			ChoiceGroup:    "censura_vingador",
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Elo da Perseguição", Edition: "4e", ClassID: &id,
			Description: "Com um ataque, você promete perseguir o inimigo se ele tentar fugir.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. CA",
			Hit:          "1[A] + mod Sabedoria de dano. Se o alvo não terminar seu próximo turno adjacente ao personagem, o vingador pode ajustar um número de quadrados igual a 1 + seu modificador de Destreza usando ação livre, mas deve terminar o ajuste mais próximo do alvo.",
			LevelScaling: "Nível 21: 2[A] + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Elo da Retribuição", Edition: "4e", ClassID: &id,
			Description: "Uma energia divina rodopiante é a promessa de vingança rápida contra os companheiros do inimigo que decidiram atacá-lo.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. CA",
			Hit:          "1[A] + mod Sabedoria de dano. Na primeira vez que um inimigo que não seja o alvo atacar o vingador até o final do próximo turno, o alvo sofre dano radiante igual ao modificador de Inteligência do vingador.",
			LevelScaling: "Nível 21: 2[A] + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Avassalador", Edition: "4e", ClassID: &id,
			Description: "Ao atacar, você manobra ao redor do inimigo, forçando-o a acompanhar seu ritmo.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. CA",
			Hit:          "1[A] + mod Sabedoria de dano. O vingador ajusta 1 quadrado e o alvo é conduzido 1 quadrado para o espaço que o personagem ocupava.",
			LevelScaling: "Nível 21: 2[A] + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Vingança Radiante", Edition: "4e", ClassID: &id,
			Description: "Ao convocar os poderes da sua divindade, você transfere a dor de seus ferimentos para um inimigo.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano radiante. O vingador recebe PV temporários iguais ao seu modificador de Sabedoria.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Eco Vingador", Edition: "4e", ClassID: &id,
			Description: "Ao brandir sua arma num arco mortal, ela deixa em seu caminho um rastro de energia radiante que obriga o inimigo a manter distância.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. CA",
			Hit:     "1[A] + mod Sabedoria de dano radiante. Até o final do próximo turno do vingador, qualquer inimigo que terminar seu turno adjacente ao personagem ou atacá-lo sofre 5 de dano radiante.",
			Special: "Censura da Retribuição: o dano radiante é igual a 5 + o modificador de Inteligência.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Fulgor Angelical", Edition: "4e", ClassID: &id,
			Description: "Você focaliza sua energia através do seu corpo para adquirir velocidade sobrenatural ao atacar um inimigo.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Effect:  "Antes do ataque, o vingador ajusta 2 quadrados.",
			Attack:  "Sabedoria vs. CA",
			Hit:     "2[A] + mod Sabedoria de dano.",
			Special: "Censura da Perseguição: o vingador ajusta um número de quadrados igual a 1 + seu modificador de Destreza.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Investida em Turbilhão", Edition: "4e", ClassID: &id,
			Description: "Ao investir contra um inimigo, uma luz divina rodeia seu corpo como uma nuvem protetora e, então, também ataca.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. CA",
			Hit:     "2[A] + mod Sabedoria de dano.",
			Special: "Ao realizar uma investida, o vingador pode utilizar este poder em lugar de um ataque básico corpo a corpo. Se o fizer, recebe +4 de bônus na CA contra ataques de oportunidade durante o movimento da investida.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Loucura Compartilhada", Edition: "4e", ClassID: &id,
			Description: "A ira do seu deus assola a mente de um inimigo e ecoa para atacar ainda outro adversário próximo.",
			Keywords:   "Divino, Implemento, Psíquico",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "1d10 + mod Sabedoria de dano psíquico. Uma segunda criatura na linha de visão do vingador sofre o mesmo dano.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Castigo Divino", Edition: "4e", ClassID: &id,
			Description: "A justiça divina cai sobre o inimigo, cegando-o com uma luz avassaladora.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. CA",
			Hit:    "3[A] + mod Sabedoria de dano radiante. O alvo fica cego até o final do encontro.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Juramento de Ânimo", Edition: "4e", ClassID: &id,
			Description: "Você pronuncia um juramento de destruição total, tornando o inimigo vulnerável a todos os seus ataques.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "3d8 + mod Sabedoria de dano radiante. O alvo fica vulnerável 5 a todos os ataques do vingador até o final do encontro.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Passo do Vingador", Edition: "4e", ClassID: &id,
			Description: "A morte de um inimigo libera uma explosão de poder divino que o impulsiona em direção ao próximo alvo.",
			Keywords:   "Divino",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special: "Gatilho: O vingador reduz um inimigo a 0 PV ou menos.",
			Effect:  "O vingador se desloca até 3 quadrados.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Escudo da Fé Vingadora", Edition: "4e", ClassID: &id,
			Description: "Sua fé inabalável cria um escudo ao redor de você quando enfrenta o inimigo sozinho.",
			Keywords:   "Divino",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special: "Gatilho: O vingador é atacado enquanto não houver aliados adjacentes.",
			Effect:  "O vingador recebe +4 de bônus de poder na CA e Reflexos contra esse ataque.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Raio do Juízo", Edition: "4e", ClassID: &id,
			Description: "Um raio de energia divina cai sobre o inimigo como o julgamento implacável do seu deus.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d8 + mod Sabedoria de dano radiante. Se o alvo for o alvo do Juramento de Inimizade, fica cego até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Salto do Vingador", Edition: "4e", ClassID: &id,
			Description: "Você se lança sobre o inimigo com velocidade divina, cobrindo grandes distâncias em um instante.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Effect:  "O vingador pode se mover até sua velocidade antes de realizar esse ataque.",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. CA",
			Hit:     "2[A] + mod Sabedoria de dano. Esse movimento não provoca ataques de oportunidade.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Execução Divina", Edition: "4e", ClassID: &id,
			Description: "Você executa seu inimigo jurado com um golpe devastador ungido pelo poder da sua divindade.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. CA",
			Hit:    "4[A] + mod Sabedoria de dano radiante. Se o alvo for o alvo do Juramento de Inimizade, causa dano adicional igual ao seu nível.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Toque do Destino", Edition: "4e", ClassID: &id,
			Description: "Você marca o inimigo com o toque do destino, selando sua morte com poder divino.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. CA",
			Hit:     "2[A] + mod Sabedoria de dano radiante. O alvo sofre dano radiante contínuo igual ao modificador de Inteligência do vingador (TR encerra).",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Fúria Divina", Edition: "4e", ClassID: &id,
			Description: "A fúria da sua divindade flui através de você, tornando seus ataques imparáveis.",
			Keywords:   "Arma, Divino",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. CA",
			Hit:     "3[A] + mod Sabedoria de dano. O alvo é empurrado 3 quadrados e fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Ira dos Céus", Edition: "4e", ClassID: &id,
			Description: "A ira plena dos céus se abate sobre o inimigo através de você.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "4d8 + mod Sabedoria de dano radiante. O alvo fica cego, atordoado e imobilizado até o final do próximo turno.",
			Miss:   "Metade do dano e o alvo fica cego até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Graça Divina", Edition: "4e", ClassID: &id,
			Description: "Você se move com a graça sobrenatural concedida pela sua divindade.",
			Keywords:   "Divino",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O vingador se move até o dobro de sua velocidade. Durante esse movimento, não provoca ataques de oportunidade e ignora terreno difícil.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Olho do Vingador", Edition: "4e", ClassID: &id,
			Description: "Você concentra sua percepção divina no alvo do seu juramento, tornando impossível para ele se esconder.",
			Keywords:   "Divino",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do encontro, o vingador pode sempre localizar o alvo do Juramento de Inimizade, ignorando ocultação, invisibilidade e escuridão.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Golpe do Anjo Vingador", Edition: "4e", ClassID: &id,
			Description: "Você golpeia com a força de um anjo vingador, causando dano devastador.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. CA",
			Hit:    "3[A] + mod Sabedoria + mod Inteligência de dano radiante. Se o alvo for o alvo do Juramento de Inimizade, causa dano adicional igual ao seu modificador de Sabedoria.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Sentença Divina", Edition: "4e", ClassID: &id,
			Description: "Você pronuncia a sentença divina sobre o inimigo, condenando-o à destruição.",
			Keywords:   "Divino, Implemento, Radiante",
			ActionType: "Ação Padrão", Range: "À distância 10",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "O alvo sofre -6 em todas as defesas até o final do encontro e fica vulnerável 10 a dano radiante.",
			Miss:   "O alvo sofre -2 em todas as defesas até o final do próximo turno.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Punição Final", Edition: "4e", ClassID: &id,
			Description: "Você desfere o golpe final com toda a força da sua fé, destruindo o inimigo com energia divina.",
			Keywords:   "Arma, Divino, Radiante",
			ActionType: "Ação Padrão", Range: "Arma corpo a corpo",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. CA",
			Hit:     "4[A] + mod Sabedoria + mod Inteligência de dano radiante. Se o alvo for o alvo do Juramento de Inimizade e tiver menos da metade dos PV máximos, o ataque causa dano máximo.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Forma do Anjo", Edition: "4e", ClassID: &id,
			Description: "Você assume brevemente a forma de um anjo divino, tornando-se uma presença aterrorizante no campo de batalha.",
			Keywords:   "Divino, Radiante",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do próximo turno, o vingador recebe velocidade de voo 6, resistência 10 a dano necrótico e radiante, e todos os inimigos a até 3 quadrados ficam cegos.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Perseguição Implacável", Edition: "4e", ClassID: &id,
			Description: "Nada pode impedir você de alcançar o alvo do seu juramento sagrado.",
			Keywords:   "Divino",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect: "O vingador se move até sua velocidade em direção ao alvo do Juramento de Inimizade. Durante esse movimento: ignora terreno difícil, não provoca ataques de oportunidade e pode se mover através de espaços de inimigos.",
			PowerType: domain.PowerEncounter, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Vingador 4e: %d habilidades processadas", len(skills))
}

func seedXamaSkills(db *gorm.DB) {
	var cls domain.Class
	if err := db.Where("name = ? AND edition = ?", "Xamã", "4e").First(&cls).Error; err != nil {
		log.Println("  ✗ Xamã 4e não encontrado")
		return
	}
	id := cls.ID

	skills := []domain.Skill{

		// ── CARACTERÍSTICAS DE CLASSE ────────────────────────────────
		{
			Name: "Presas do Espírito", Edition: "4e", ClassID: &id,
			Description: "Quando o inimigo baixa a guarda, seu companheiro espiritual salta sobre ele, mordendo e arranhando.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação de Oportunidade", Range: "Espírito corpo a corpo 1",
			Special:        "Gatilho: Um inimigo deixa um quadrado adjacente ao companheiro espiritual sem ajustar.",
			Target:         "O inimigo que ativou o gatilho",
			Attack:         "Sabedoria vs. Reflexos",
			Hit:            "1d10 + mod Sabedoria de dano.",
			LevelScaling:   "Nível 21: 2d10 + mod Sabedoria.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
		},
		{
			Name: "Espírito Perseguidor", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual é especialmente eficaz contra inimigos que já estão feridos.",
			Keywords:   "Espírito, Primitivo",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special:        "Sempre que o xamã acertar um ataque, se o alvo estiver sangrando, recebe um bônus de ataque igual à metade do modificador de Inteligência.",
			Effect:         "O companheiro espiritual pode flanquear com o personagem ou aliados do xamã.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
			RequiresChoice: true,
			ChoiceGroup:    "espirito_xama",
		},
		{
			Name: "Espírito Protetor", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual protege os aliados próximos, curando-os quando necessário.",
			Keywords:   "Cura, Espírito, Primitivo",
			ActionType: "Ação Livre", Range: "À distância 5",
			Special:        "Gatilho: Um aliado a até 5 quadrados do companheiro espiritual gasta um pulso de cura.",
			Effect:         "O aliado recupera PV adicionais iguais ao modificador de Constituição do xamã.",
			PowerType:      domain.PowerUnlimited, Level: 1,
			IsClassFeature: true,
			RequiresChoice: true,
			ChoiceGroup:    "espirito_xama",
		},

		// ── NÍVEL 1 — SEM LIMITE ────────────────────────────────────
		{
			Name: "Cólera do Inverno", Edition: "4e", ClassID: &id,
			Description: "Os espíritos do inverno cercam o inimigo, rasgando-o com suas presas e garras e chamando seu companheiro espiritual para se juntar ao ataque.",
			Keywords:   "Congelante, Espírito, Implemento, Primitivo, Teletransporte",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Fortitude",
			Hit:          "1d10 + mod Sabedoria de dano congelante. O xamã pode teleportar seu companheiro espiritual para um espaço adjacente ao alvo.",
			LevelScaling: "Nível 21: 2d10 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Espíritos Assombrosos", Edition: "4e", ClassID: &id,
			Description: "Espíritos lamentadores surgem ao redor do inimigo, distraindo-o dos ataques de seus aliados.",
			Keywords:   "Implemento, Primitivo, Psíquico",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Vontade",
			Hit:          "1d6 + mod Sabedoria de dano psíquico. Até o final do próximo turno do xamã, o alvo concede vantagem de combate à próxima criatura que o atacar.",
			LevelScaling: "Nível 21: 2d6 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Defensor", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual ataca o inimigo, drenando energia dele para usar como um escudo protetor para seus aliados.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano. Até o final do próximo turno do xamã, o personagem e seus aliados recebem +1 de bônus de poder na CA enquanto estiverem adjacentes ao companheiro espiritual.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe do Perseguidor", Edition: "4e", ClassID: &id,
			Description: "Conforme o ataque do companheiro atinge o inimigo, esse espírito se enche de uma fúria predatória, tornando-se uma ameaça ainda maior para os adversários.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Fortitude",
			Hit:          "1d10 + mod Sabedoria de dano. Se o alvo estiver sangrando, o xamã recebe um bônus na jogada de ataque igual à metade do modificador de Inteligência. Até o final do próximo turno, seu companheiro espiritual pode flanquear com o personagem ou aliados.",
			LevelScaling: "Nível 21: 2d10 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe do Vigilante", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual encurrala os inimigos, distraindo-os e criando aberturas para ataques. Você e seus aliados também compartilham os sentidos afiados do espírito.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Reflexos",
			Hit:          "1d8 + mod Sabedoria de dano. Até o final do próximo turno, o xamã e seus aliados recebem +1 de bônus nas jogadas de ataque e +5 nos testes de Percepção enquanto estiverem adjacentes ao companheiro espiritual.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},
		{
			Name: "Golpe Protetor", Edition: "4e", ClassID: &id,
			Description: "Ecos estrondosos de antigas cavernas e grutas seguem o ataque de seu companheiro, enchendo seus aliados com vitalidade.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target:       "Uma criatura",
			Attack:       "Sabedoria vs. Vontade",
			Hit:          "1d8 + mod Sabedoria de dano. Os aliados adjacentes ao companheiro espiritual recebem PV temporários iguais ao modificador de Constituição do xamã.",
			LevelScaling: "Nível 21: 2d8 + mod Sabedoria.",
			PowerType: domain.PowerUnlimited, Level: 1,
		},

		// ── NÍVEL 1 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Chamado ao Combatente Ancestral", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual canaliza um poderoso espírito ancestral para atacar seus adversários e reforçar as defesas de seus aliados.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "1d10 + mod Sabedoria de dano. Até o final do próximo turno do xamã, ele e seus aliados recebem +2 de bônus em todas as defesas enquanto adjacentes ao companheiro espiritual.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Chamado ao Defensor Ancestral", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual canaliza o espírito de um guerreiro ancestral que protege a retirada segura de seus aliados.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d8 + mod Sabedoria de dano. Até o final do próximo turno do xamã, ele e seus aliados recebem +5 de bônus em todas as defesas contra ataques de oportunidade enquanto adjacentes ao companheiro espiritual.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Panteras Gêmeas", Edition: "4e", ClassID: &id,
			Description: "Dois espíritos de panteras saltam sobre o inimigo, canalizando seus instintos predatórios através do companheiro espiritual.",
			Keywords:   "Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. Reflexos",
			Hit:     "1d8 + mod Sabedoria de dano.",
			Effect:  "O xamã realiza esse ataque novamente contra o mesmo alvo ou um alvo diferente.",
			Special: "Espírito Perseguidor: se o alvo estiver sangrando, recebe bônus de ataque igual ao mod Inteligência.",
			PowerType: domain.PowerEncounter, Level: 1,
		},
		{
			Name: "Proteção do Urso do Trovão", Edition: "4e", ClassID: &id,
			Description: "Um espírito urso ancião ruge como o trovão e canaliza sua força para dar suporte a seus aliados.",
			Keywords:   "Implemento, Primitivo, Trovejante",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. Fortitude",
			Hit:     "1d6 + mod Sabedoria de dano trovejante. Até o final do próximo turno, o xamã e seus aliados adjacentes ao companheiro espiritual adquirem resistência a todos os danos igual ao modificador de Constituição.",
			Special: "Espírito Protetor: o xamã ou um aliado a até 5 recebe PV temporários iguais ao mod CON.",
			PowerType: domain.PowerEncounter, Level: 1,
		},

		// ── NÍVEL 1 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Bênção dos Sete Ventos", Edition: "4e", ClassID: &id,
			Description: "Você convoca os espíritos dos sete ventos que surgem no campo de batalha, derrubando um inimigo e os arremessando para longe dos demais.",
			Keywords:   "Implemento, Primitivo, Zona",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "2d10 + mod Sabedoria de dano. O alvo é conduzido 2 quadrados.",
			Miss:   "Metade do dano.",
			Effect: "O ataque cria uma zona de redemoinhos centralizada no alvo. A zona persiste até o final do encontro. Com ação de movimento o xamã pode mover a zona 5 quadrados. Com ação mínima, pode fazer qualquer alvo dentro da zona ser conduzido 1 quadrado.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Cólera do Mundo Espiritual", Edition: "4e", ClassID: &id,
			Description: "Os espíritos furiosos assaltam a mente dos inimigos ao seu redor e do seu companheiro espiritual.",
			Keywords:   "Implemento, Primitivo, Psíquico",
			ActionType: "Ação Padrão", Range: "Explosão contígua 2",
			Target: "Os inimigos dentro da explosão ou adjacentes ao companheiro espiritual",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "3d6 + mod Sabedoria de dano psíquico. O alvo fica derrubado.",
			Miss:   "Metade do dano.",
			PowerType: domain.PowerDaily, Level: 1,
		},
		{
			Name: "Espírito da Inundação de Cura", Edition: "4e", ClassID: &id,
			Description: "O espírito de uma grande inundação aparece na forma de uma criatura formada por águas agitadas, sustentando seus aliados e afogando os inimigos.",
			Keywords:   "Cura, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 5",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Fortitude",
			Hit:    "1d8 + mod Sabedoria de dano.",
			Miss:   "Metade do dano.",
			Effect: "Até o final do encontro, o xamã e seus aliados dentro da explosão adquirem regeneração 2 enquanto estiverem sangrando. Com ação mínima, qualquer alvo pode encerrar esse efeito em si mesmo para recuperar 10 PV.",
			PowerType: domain.PowerDaily, Level: 1,
		},

		// ── NÍVEL 2 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Mover o Companheiro", Edition: "4e", ClassID: &id,
			Description: "Com um simples pensamento, você reposiciona seu companheiro espiritual no campo de batalha.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "O xamã move o companheiro espiritual até 6 quadrados.",
			PowerType: domain.PowerEncounter, Level: 2,
		},
		{
			Name: "Cura do Companheiro", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual flui até o lado do aliado ferido e canaliza energia curativa através do vínculo espiritual.",
			Keywords:   "Cura, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Target: "O personagem ou um aliado",
			Effect: "O alvo pode gastar um pulso de cura e recupera PV adicionais iguais ao modificador de Sabedoria do xamã. O companheiro espiritual se teletransporta para um espaço adjacente ao alvo.",
			PowerType: domain.PowerEncounter, Level: 2,
		},

		// ── NÍVEL 3 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Uivo do Lobo da Tempestade", Edition: "4e", ClassID: &id,
			Description: "Seu companheiro espiritual solta um uivo que aterroriza os inimigos e fortalece seus aliados.",
			Keywords:   "Espírito, Implemento, Medo, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 3 (centrada no companheiro)",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "1d8 + mod Sabedoria de dano psíquico. O alvo recua 2 quadrados.",
			Effect: "Os aliados dentro da explosão recebem +2 de bônus de poder nas jogadas de ataque até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},
		{
			Name: "Garras do Espírito Predador", Edition: "4e", ClassID: &id,
			Description: "As garras do espírito predador rasgam o inimigo, marcando-o como presa.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target: "Uma criatura",
			Attack: "Sabedoria vs. CA",
			Hit:    "2d8 + mod Sabedoria de dano. O alvo fica marcado e sofre -2 em todas as defesas até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 3,
		},

		// ── NÍVEL 3 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Invocação do Grande Espírito", Edition: "4e", ClassID: &id,
			Description: "Você invoca um Grande Espírito que fortalece todos os seus aliados durante todo o encontro.",
			Keywords:   "Conjuração, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "O Grande Espírito aparece num quadrado desocupado. Enquanto ativo, todos os aliados a até 5 quadrados do Grande Espírito recebem +1 de bônus de poder nas jogadas de ataque e nos danos. O Grande Espírito persiste até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 3,
		},

		// ── NÍVEL 5 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Fúria dos Animais Espirituais", Edition: "4e", ClassID: &id,
			Description: "Uma horda de espíritos animais furiosos ataca os inimigos ao redor do companheiro espiritual.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 2 (centrada no companheiro)",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Reflexos",
			Hit:    "2d6 + mod Sabedoria de dano. O alvo fica derrubado.",
			PowerType: domain.PowerEncounter, Level: 5,
		},
		{
			Name: "Espírito Curador das Marés", Edition: "4e", ClassID: &id,
			Description: "As marés espirituais fluem através do campo de batalha, curando seus aliados e ferindo seus inimigos.",
			Keywords:   "Cura, Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. Fortitude",
			Hit:     "2d8 + mod Sabedoria de dano.",
			Effect:  "Um aliado adjacente ao companheiro espiritual pode gastar um pulso de cura.",
			PowerType: domain.PowerEncounter, Level: 5,
		},

		// ── NÍVEL 5 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Tempestade Espiritual", Edition: "4e", ClassID: &id,
			Description: "Uma tempestade de espíritos furiosos varre o campo de batalha, devastando seus inimigos.",
			Keywords:   "Implemento, Primitivo, Zona",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "3d6 + mod Sabedoria de dano psíquico. O alvo fica atordoado até o final do próximo turno.",
			Miss:   "Metade do dano.",
			Effect: "A zona persiste até o final do encontro. Criaturas que começarem o turno dentro da zona sofrem 1d6 de dano psíquico.",
			PowerType: domain.PowerDaily, Level: 5,
		},

		// ── NÍVEL 6 — UTILITÁRIO ────────────────────────────────────
		{
			Name: "Vínculo Espiritual", Edition: "4e", ClassID: &id,
			Description: "Você fortalece o vínculo espiritual com um aliado, permitindo que o companheiro espiritual cure ele à distância.",
			Keywords:   "Cura, Primitivo",
			ActionType: "Ação Mínima", Range: "À distância 10",
			Target: "Um aliado",
			Effect: "Até o final do encontro, sempre que o alvo gastar um pulso de cura, recupera PV adicionais iguais ao modificador de Sabedoria do xamã, desde que esteja a até 10 quadrados do companheiro espiritual.",
			PowerType: domain.PowerEncounter, Level: 6,
		},
		{
			Name: "Forma Espiritual", Edition: "4e", ClassID: &id,
			Description: "Você assume brevemente uma forma espiritual translúcida, tornando-se difícil de atingir.",
			Keywords:   "Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do próximo turno, o xamã recebe ocultação e pode mover-se através de espaços de criaturas sem penalidade.",
			PowerType: domain.PowerEncounter, Level: 6,
		},

		// ── NÍVEL 7 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Rugido do Urso Ancestral", Edition: "4e", ClassID: &id,
			Description: "O espírito de um urso ancestral ruge através do seu companheiro, fortalecendo todos os aliados próximos.",
			Keywords:   "Espírito, Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Espírito corpo a corpo 1",
			Target:  "Uma criatura",
			Attack:  "Sabedoria vs. Fortitude",
			Hit:     "2d10 + mod Sabedoria de dano. O alvo fica derrubado.",
			Effect:  "O xamã e seus aliados adjacentes ao companheiro espiritual recebem +2 de bônus de poder na CA e regeneração 3 até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 7,
		},

		// ── NÍVEL 7 — DIÁRIO ────────────────────────────────────────
		{
			Name: "Grande Espírito da Cura", Edition: "4e", ClassID: &id,
			Description: "Você invoca um Grande Espírito da Cura que percorre o campo de batalha restaurando a saúde de seus aliados.",
			Keywords:   "Conjuração, Cura, Primitivo",
			ActionType: "Ação Padrão", Range: "À distância 5",
			Effect: "Um Grande Espírito da Cura aparece. No início de cada turno do xamã, o espírito pode se mover até 4 quadrados e curar um aliado adjacente: o aliado recupera PV iguais ao modificador de Sabedoria do xamã sem gastar um pulso de cura. O espírito persiste até o final do encontro.",
			PowerType: domain.PowerDaily, Level: 7,
		},

		// ── NÍVEL 9 — POR ENCONTRO ──────────────────────────────────
		{
			Name: "Horda de Espíritos", Edition: "4e", ClassID: &id,
			Description: "Uma horda de espíritos selvagens devasta o campo de batalha, atacando todos os inimigos.",
			Keywords:   "Implemento, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão de área 3 a até 10 quadrados",
			Target: "Os inimigos dentro da explosão",
			Attack: "Sabedoria vs. Vontade",
			Hit:    "3d8 + mod Sabedoria de dano. O alvo fica imobilizado e atordoado até o final do próximo turno.",
			PowerType: domain.PowerEncounter, Level: 9,
		},

		// ── NÍVEL 10 — UTILITÁRIO ───────────────────────────────────
		{
			Name: "Fusão com o Espírito", Edition: "4e", ClassID: &id,
			Description: "Você se funde completamente com seu companheiro espiritual, tornando-se uma força da natureza imparável.",
			Keywords:   "Metamorfose, Primitivo",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect: "Até o final do encontro, o xamã e seu companheiro espiritual se fundem. O xamã recebe +2 de bônus de poder em todas as defesas, velocidade de voo 6 e seus poderes de espírito passam a ter alcance igual ao seu deslocamento.",
			PowerType: domain.PowerDaily, Level: 10,
		},
		{
			Name: "Invocação dos Ancestrais", Edition: "4e", ClassID: &id,
			Description: "Você invoca todos os espíritos ancestrais para proteger e fortalecer seus aliados.",
			Keywords:   "Conjuração, Primitivo",
			ActionType: "Ação Padrão", Range: "Explosão contígua 5",
			Effect: "Até o final do encontro, todos os aliados dentro da explosão recebem +2 de bônus de poder em todas as jogadas de ataque e danos, e regeneração 2 enquanto estiverem adjacentes ao companheiro espiritual do xamã.",
			PowerType: domain.PowerDaily, Level: 10,
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Xamã 4e: %d habilidades processadas", len(skills))
}
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
 
func seedRaceSkills(db *gorm.DB) {
	// LJ1
	seedDraconatoRaceSkills(db)
	seedEladrinRaceSkills(db)
	seedElfoRaceSkills(db)
	seedHalflingRaceSkills(db)
	seedTieflingRaceSkills(db)
	// Anão, Humano e Meio-Elfo: sem poder racial com card próprio
 
	// LJ2
	seedDevaRaceSkills(db)
	seedGnomoRaceSkills(db)
	seedGoliasRaceSkills(db)
	seedMeioOrcRaceSkills(db)
	seedFeralRaceSkills(db)
 
	// LJ3
	seedFragmentalRaceSkills(db)
	seedGithzeraiRaceSkills(db)
	seedMinotauroRaceSkills(db)
	seedSelvioRaceSkills(db)
}
 
// ── LJ1 ────────────────────────────────────────────────────────────────────────
 
// Draconato — Sopro de Dragão (escolha de tipo de dano e atributo na criação)
func seedDraconatoRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Draconato", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Draconato 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Sopro de Dragão", Edition: "4e", RaceID: &id,
			Description: "Escancarando sua boca com um rugido, o poder mortífero dos seus ancestrais dracônicos explode numa rajada, engolfando seus inimigos.",
			Keywords:   "Ácido, Congelante, Flamejante, Elétrico ou Venenoso",
			ActionType: "Ação Mínima", Range: "Rajada contígua 3",
			Target:  "As criaturas dentro da área",
			Attack:  "Força +2 vs. Reflexos, Constituição +2 vs. Reflexos, ou Destreza +2 vs. Reflexos",
			Hit:     "1d6 + modificador de Constituição de dano.",
			Special: "Ao criar o personagem, escolha Força, Constituição ou Destreza como o valor de atributo relevante para as jogadas de ataque e dano desse poder. Você também deve escolher o tipo de dano: ácido, congelante, flamejante, elétrico ou venenoso. Essas duas opções permanecem durante toda a vida do personagem.",
			LevelScaling: "Nível 11: +4 de bônus e 2d6 + mod CON. Nível 21: +6 de bônus e 3d6 + mod CON.",
			PowerType: domain.PowerEncounter, Level: 1,
			IsRaceFeature: true, RequiresChoice: true, ChoiceGroup: "sopro_draconato",
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Draconato 4e: habilidades raciais seedadas")
}
 
// Eladrin — Passo Féerico (automático)
func seedEladrinRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Eladrin", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Eladrin 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Passo Féerico", Edition: "4e", RaceID: &id,
			Description: "Com um passo, você desaparece de um lugar e ressurge em outro.",
			Keywords:   "Teleporte",
			ActionType: "Ação de Movimento", Range: "Pessoal",
			Effect:        "O eladrin se teleporta até 5 quadrados.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Eladrin 4e: habilidades raciais seedadas")
}
 
// Elfo — Precisão Élfica (automático)
func seedElfoRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Elfo", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Elfo 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Precisão Élfica", Edition: "4e", RaceID: &id,
			Description: "Com um momento de concentração, você mira cuidadosamente no seu adversário e dispara com a precisão lendária dos elfos.",
			ActionType: "Ação Livre", Range: "Pessoal",
			Effect:        "O elfo refaz sua jogada de ataque. Utilize o segundo resultado, mesmo que seja inferior.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Elfo 4e: habilidades raciais seedadas")
}
 
// Halfling — Segunda Chance (automático)
func seedHalflingRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Halfling", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Halfling 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Segunda Chance", Edition: "4e", RaceID: &id,
			Description: "A sorte e o tamanho reduzido se combinam para favorecê-lo conforme você se desvia do ataque do inimigo.",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special:       "Quando um ataque atingir o personagem, o inimigo é obrigado a refazer a jogada. O atacante deve utilizar o segundo resultado, mesmo que seja inferior.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Halfling 4e: habilidades raciais seedadas")
}
 
// Tiefling — Cólera Infernal (automático)
func seedTieflingRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Tiefling", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Tiefling 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Cólera Infernal", Edition: "4e", RaceID: &id,
			Description: "Você convoca sua natureza furiosa para aumentar sua capacidade de ferir seus inimigos.",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:        "O tiefling canaliza sua fúria para receber +1 de bônus racial na próxima jogada de ataque contra um inimigo que o tenha atingido desde seu último turno. Se obtiver sucesso nesse ataque e ele causar dano, adicione o modificador de Carisma do tiefling como dano adicional.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Tiefling 4e: habilidades raciais seedadas")
}
 
// ── LJ2 ────────────────────────────────────────────────────────────────────────
 
// Deva — Memória de Mil Vidas (automático)
func seedDevaRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Deva", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Deva 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Memória de Mil Vidas", Edition: "4e", RaceID: &id,
			Description: "Como sonhos, memórias de vidas passadas retornam para ajudá-lo.",
			ActionType: "Nenhuma Ação", Range: "Pessoal",
			Special:       "Gatilho: O deva realiza uma jogada de ataque, teste de resistência ou teste de atributo e não gosta do resultado.",
			Effect:        "O personagem adiciona 1d6 ao resultado da jogada que ativou o gatilho.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Deva 4e: habilidades raciais seedadas")
}
 
// Gnomo — Desvanecer (automático)
func seedGnomoRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Gnomo", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Gnomo 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Desvanecer", Edition: "4e", RaceID: &id,
			Description: "Você fica invisível em resposta ao perigo.",
			Keywords:   "Ilusão",
			ActionType: "Reação Imediata", Range: "Pessoal",
			Special:       "Gatilho: O gnomo sofre dano.",
			Effect:        "O gnomo fica invisível até atacar ou até o final do seu próximo turno.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Gnomo 4e: habilidades raciais seedadas")
}
 
// Golias — Resistência de Pedra (automático)
func seedGoliasRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Golias", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Golias 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Resistência de Pedra", Edition: "4e", RaceID: &id,
			Description: "Os ataques de seus inimigos resvalam em seu corpo pedregoso.",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Effect:        "O golias adquire resistência 5 contra todos os tipos de dano até o final do seu próximo turno.",
			LevelScaling: "Nível 11: Resistência 10. Nível 21: Resistência 15.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Golias 4e: habilidades raciais seedadas")
}
 
// Meio-Orc — Assalto Furioso (automático)
func seedMeioOrcRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Meio-Orc", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Meio-Orc 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Assalto Furioso", Edition: "4e", RaceID: &id,
			Description: "Sua ira monstruosa queima dentro de seu corpo, aumentando a força de seu ataque.",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special:       "Gatilho: O personagem atinge um inimigo.",
			Effect:        "O ataque causa 1[A] de dano adicional se for um ataque com arma ou 1d8 de dano adicional se não for.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Meio-Orc 4e: habilidades raciais seedadas")
}
 
// Feral — escolha obrigatória entre duas mutações (Dente-Alongado OU Garra-Navalha)
// A escolha determina qual subtipo de Feral o personagem é
func seedFeralRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Feral", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Feral 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		// Feral Dente-Alongado
		{
			Name: "Mutação do Dente-Alongado", Edition: "4e", RaceID: &id,
			Description: "Você libera a fera primitiva em seu interior e assume um aspecto mais selvagem. Seus dentes crescem, aumentando sua ferocidade em combate.",
			Keywords:   "Cura",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Special:       "Condição: O personagem deve estar sangrando.",
			Effect:        "Até o final do encontro, o personagem recebe +2 de bônus nas jogadas de dano. Além disso, ele adquire regeneração 2 enquanto estiver sangrando.",
			LevelScaling: "Nível 11: Regeneração 4. Nível 21: Regeneração 6.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true, RequiresChoice: true, ChoiceGroup: "mutacao_feral",
		},
		// Feral Garra-Navalha
		{
			Name: "Mutação do Garra-Navalha", Edition: "4e", RaceID: &id,
			Description: "Você libera a fera primitiva em seu interior e assume um aspecto mais selvagem. Suas garras se alongam como navalhas, tornando-o mais ágil e difícil de atingir.",
			ActionType: "Ação Mínima", Range: "Pessoal",
			Special:       "Condição: O personagem deve estar sangrando.",
			Effect:        "Até o final do encontro, o deslocamento do personagem aumenta em 2 quadrados e ele recebe +1 de bônus na CA e na defesa de Reflexos.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true, RequiresChoice: true, ChoiceGroup: "mutacao_feral",
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Feral 4e: habilidades raciais seedadas")
}
 
// ── LJ3 ────────────────────────────────────────────────────────────────────────
 
// Fragmental — Enxame de Fragmentos (automático)
func seedFragmentalRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Fragmental", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Fragmental 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Enxame de Fragmentos", Edition: "4e", RaceID: &id,
			Description: "Você libera o controle mental sobre sua forma física, distraindo seus inimigos com uma nuvem de fragmentos. Você então reforma seu corpo em outro lugar.",
			Keywords:   "Teleporte",
			ActionType: "Ação de Movimento", Range: "Explosão contígua 1",
			Target:        "Cada inimigo dentro da explosão",
			Effect:        "Cada alvo concede vantagem de combate ao fragmental até o final do próximo turno do personagem. O fragmental então teleporta metade de seu deslocamento.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Fragmental 4e: habilidades raciais seedadas")
}
 
// Githzerai — Mente de Ferro (automático)
func seedGithzeraiRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Githzerai", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Githzerai 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Mente de Ferro", Edition: "4e", RaceID: &id,
			Description: "Sob o peso de um ataque, você usa o poder de sua mente para fortalecer-se contra danos.",
			ActionType: "Interrupção Imediata", Range: "Pessoal",
			Special:       "Gatilho: O githzerai é atingido por um ataque.",
			Effect:        "O githzerai recebe +2 de bônus em todas as defesas até o final do seu próximo turno.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Githzerai 4e: habilidades raciais seedadas")
}
 
// Minotauro — Investida com Chifres (automático)
func seedMinotauroRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Minotauro", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Minotauro 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		{
			Name: "Investida com Chifres", Edition: "4e", RaceID: &id,
			Description: "Você investe contra o inimigo e o perfura com seus chifres.",
			ActionType: "Ação Padrão", Range: "Corpo a corpo 1",
			Effect:        "O minotauro realiza uma investida e realiza o seguinte ataque no lugar do ataque básico corpo a corpo.",
			Target:        "Uma criatura",
			Attack:        "Força, Constituição ou Destreza +4 vs. CA",
			Hit:           "1d6 + modificador de Força, Constituição ou Destreza de dano e o alvo fica derrubado.",
			LevelScaling: "Nível 11: 2d6 + modificador. Nível 21: 3d6 + modificador.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true,
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Minotauro 4e: habilidades raciais seedadas")
}
 
// Sélvio — escolha obrigatória de Aspecto (3 opções, muda a cada descanso prolongado)
// Cada aspecto concede um poder racial diferente
func seedSelvioRaceSkills(db *gorm.DB) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", "Sélvio", "4e").First(&race).Error; err != nil {
		log.Println("  ✗ Raça Sélvio 4e não encontrada"); return
	}
	id := race.ID
 
	skills := []domain.Skill{
		// Aspecto do Destruidor
		{
			Name: "Ira do Destruidor", Edition: "4e", RaceID: &id,
			Description: "Seu aspecto destrutivo responde a um ataque com força mortal. Quando um inimigo ferido ousa atacar você ou seus aliados, você reage com violência imediata.",
			ActionType: "Reação Imediata", Range: "Pessoal",
			Special:       "Gatilho: Um inimigo sangrando ataca o sélvio ou um aliado adjacente ao sélvio. Condição: O sélvio deve estar manifestando o Aspecto do Destruidor.",
			Effect:        "O sélvio decide realizar um ataque básico corpo a corpo ou uma investida contra o inimigo que ativou o gatilho. Se o ataque for bem-sucedido, o inimigo fica pasmo até o final do próximo turno do sélvio.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true, RequiresChoice: true, ChoiceGroup: "aspecto_selvio",
		},
		// Aspecto do Caçador
		{
			Name: "Perseguição do Caçador", Edition: "4e", RaceID: &id,
			Description: "Sua presa tenta se distanciar de você, mas não há escapatória. O instinto caçador toma conta de você quando o inimigo tenta recuar.",
			ActionType: "Reação Imediata", Range: "Pessoal",
			Special:       "Gatilho: Um inimigo a até 2 quadrados do sélvio se desloca no turno dele. Condição: O sélvio deve estar manifestando o Aspecto do Caçador.",
			Effect:        "O sélvio ajusta 3 quadrados. Até o final do próximo turno do sélvio, ele causa 1d6 de dano adicional ao inimigo que ativou o gatilho quando o sélvio o atinge. O sélvio ignora a penalidade de -2 nas jogadas de ataque ao atacar o inimigo que tiver cobertura ou ocultação.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true, RequiresChoice: true, ChoiceGroup: "aspecto_selvio",
		},
		// Aspecto dos Anciões
		{
			Name: "Viagem dos Anciões", Edition: "4e", RaceID: &id,
			Description: "Você desaparece e deixa um inimigo perplexo em seu rastro. O conhecimento ancestral fae permite que você se teleporte após atingir um grupo de inimigos.",
			Keywords:   "Teleporte",
			ActionType: "Ação Livre", Range: "Pessoal",
			Special:       "Gatilho: O sélvio atinge um inimigo com um ataque de área ou contíguo. Condição: O sélvio deve estar manifestando o Aspecto dos Anciões.",
			Effect:        "O sélvio teleporta 3 quadrados. Escolha um inimigo atingido pelo ataque: o sélvio e um aliado dentro da linha de visão do sélvio adquirem vantagem de combate contra aquele inimigo até o final do próximo turno do sélvio.",
			PowerType:     domain.PowerEncounter, Level: 1,
			IsRaceFeature: true, RequiresChoice: true, ChoiceGroup: "aspecto_selvio",
		},
	}
	for _, s := range skills { upsertRaceSkill(db, s, id) }
	log.Println("  ✓ Sélvio 4e: habilidades raciais seedadas")
}