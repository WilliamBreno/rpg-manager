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
	seedSkills(db)
	seedSkills5e(db)
	seedRaceSkills(db)
	seedPericias(db)
    seedClassPericias(db)
	seedPericias5e(db) 
    seedRacePericias(db)
    seedTalentos(db)
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
			Name:    "Patrono Arcano",
			Edition: "5e", ClassID: &id,
			Description: "Escolha o patrono com o qual você firmou seu pacto: Arquifada, Celestial, O Grande Antigo ou Ínfero. A escolha concede magias sempre preparadas e características exclusivas a partir do nível 1.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Define a origem do poder do Bruxo.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "patrono_bruxo",
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
			Name:    "Domínio Divino",
			Edition: "5e", ClassID: &id,
			Description: "Escolha um domínio da sua divindade: Guerra, Luz, Trapaça ou Vida. O domínio concede magias sempre preparadas e características exclusivas desde o nível 1.",
			Keywords: "Divino", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Define a especialização do Clérigo e concede poderes de domínio.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "dominio_clerigo",
		},
	}
	for _, s := range skills {
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
			Description: "Você conhece Druídico, a língua secreta dos druidas. Pode falar e compreender essa língua e deixar mensagens ocultas que apenas outros druidas conseguem decifrar em superfícies naturais.",
			Keywords: "Primitivo", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Conhecimento da língua secreta Druídico.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
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
			Name:    "Origem da Feitiçaria",
			Edition: "5e", ClassID: &id,
			Description: "Escolha a fonte do seu poder mágico inato: Aberrante, Dracônica, Mecânica ou Selvagem. A origem concede características exclusivas a partir do nível 1.",
			Keywords: "Arcano", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Define a fonte do poder do Feiticeiro.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "origem_feiticeiro",
		},
	}
	for _, s := range skills {
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
			Name:    "Especialista",
			Edition: "5e", ClassID: &id,
			Description: "Escolha 2 perícias nas quais você já tenha proficiência. Você recebe Especialização nessas perícias: seu Bônus de Proficiência é dobrado nos testes que as usem.",
			Keywords: "Marcial", ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Bônus de Proficiência dobrado em 2 perícias escolhidas.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsClassFeature: true,
			RequiresChoice: true, ChoiceGroup: "especialista_guardiao",
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
			Name:    "Recuperação",
			Edition: "5e", ClassID: &id,
			Description: "Como Ação Bônus, você se recupera milagrosamente, curando 1d10 + nível de Guerreiro PV. 1 uso — recupera em Descanso Curto ou Longo.",
			Keywords: "Marcial, Cura", ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:       "Cura 1d10 + nível de Guerreiro PV. 1 uso por descanso curto.",
			LevelScaling: "Escala com o nível de Guerreiro.",
			PowerType:    domain.PowerEncounter, Level: 1, IsClassFeature: true,
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
			Name:    "Memória de Magia",
			Edition: "5e", ClassID: &id,
			Description: "Você pode conjurar magias marcadas com Ritual sem gastar um espaço de magia, adicionando 10 minutos ao tempo de conjuração. A magia deve estar no seu Grimório e possuir a marca de Ritual.",
			Keywords: "Arcano, Ritual", ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Conjura magias de Ritual do Grimório sem gastar espaços de magia.",
			PowerType: domain.PowerUnlimited, Level: 1, IsClassFeature: true,
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