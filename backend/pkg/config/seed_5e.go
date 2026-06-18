package config

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// CLASSES 5e
// ─────────────────────────────────────────────────────────────────────────────

// seedClasses5e atualiza/cria as 12 classes do PHB 2024 com dados completos.
// Roda DEPOIS de seedClasses para sobrescrever os stubs já inseridos.
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

// ─────────────────────────────────────────────────────────────────────────────
// CARACTERÍSTICAS DE CLASSE 5e — NÍVEL 1
// ─────────────────────────────────────────────────────────────────────────────

// ─────────────────────────────────────────────────────────────────────────────
// PERÍCIAS 5e  (18 perícias do PHB 2024)
// Verifica: o nome do struct pode ser "Pericia" ou "Skill" no seu domain —
// ajuste domain.Pericia abaixo se necessário.
// ─────────────────────────────────────────────────────────────────────────────

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

// ─────────────────────────────────────────────────────────────────────────────
// CARACTERÍSTICAS DE CLASSE 5e — NÍVEL 1
// ─────────────────────────────────────────────────────────────────────────────

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

// ─────────────────────────────────────────────────────────────────────────────
// ESPÉCIES 5e  (PHB 2024)
// ─────────────────────────────────────────────────────────────────────────────

// seedRaces5e cria/atualiza as 10 espécies do PHB 2024.
func seedRaces5e(db *gorm.DB) {
	type data struct {
		Name, Description string
	}

	races := []data{
		{"Aasimar",    "Mortais com uma centelha dos Planos Superiores. Carregam herança celestial que se manifesta em resistências e poderes de luz e cura."},
		{"Anão",       "Forjados da terra por Moradin. Resilientes como montanhas, com afinidade por pedra, metal e vida subterrânea. Vivem cerca de 350 anos."},
		{"Draconato",  "Humanoides de herança dracônica, com escamas, sopro elemental e resistência ao tipo de dano do seu ancestral dragão."},
		{"Elfo",       "Seres feéricos com sentidos aguçados, resistência à magia, e a capacidade de meditar em transe no lugar de dormir. Vivem séculos."},
		{"Gnomo",      "Pequenos humanoides com astúcia mágica natural, visão no escuro e resistência a magias que afetam a mente."},
		{"Golias",     "Gigantes em miniatura com forte conexão com o povo gigante, podendo resistir a danos devastadores e assumir uma forma aumentada."},
		{"Humano",     "A espécie mais versátil do multiverso. Recebem um talento de Origem adicional e se adaptam a qualquer caminho."},
		{"Orc",        "Humanoides robustos abençoados por Gruumsh. Determinados, com visão no escuro poderosa e a capacidade de continuar lutando ao cair."},
		{"Pequenino",  "Pequenos humanoides corajosos com sorte inata, furtividade natural e agilidade para se mover entre criaturas maiores."},
		{"Tiferino",   "Humanoides com sangue ínfero. Possuem visão no escuro e um legado mágico que reflete sua linhagem — diabólica, demoníaca ou outra."},
	}

	for _, r := range races {
		var existing domain.Race
		if db.Where("name = ? AND edition = ?", r.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Race{
				Name:        r.Name,
				Edition:     "5e",
				Description: r.Description,
				Speed:       30, // PHB 2024: todas as espécies têm Deslocamento 9m (30ft)
				IsDefault:   true,
			})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description": r.Description,
				"speed":       30,
			})
		}
	}
	log.Println("  ✓ Espécies 5e: 10 espécies seedadas")
}

// seedRaceFeatures5e semeia os traços raciais de todas as espécies 5e.
func seedRaceFeatures5e(db *gorm.DB) {
	seedAasimar5e(db)
	seedAnao5e(db)
	seedDraconato5e(db)
	seedElfo5e(db)
	seedGnomo5e(db)
	seedGolias5e(db)
	seedHumano5e(db)
	seedOrc5e(db)
	seedPequenino5e(db)
	seedTiferino5e(db)
}

// getRace5e busca uma espécie 5e pelo nome.
func getRace5e(db *gorm.DB, name string) (uint, bool) {
	var race domain.Race
	if err := db.Where("name = ? AND edition = ?", name, "5e").First(&race).Error; err != nil {
		log.Printf("  ✗ Espécie 5e não encontrada: %s (rode seedRaces5e primeiro)", name)
		return 0, false
	}
	return race.ID, true
}

// ── Aasimar ───────────────────────────────────────────────────────────────────

func seedAasimar5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Aasimar")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Resistência Celestial", Edition: "5e", RaceID: &id,
			Description: "Você tem Resistência a dano Necrótico e a dano Radiante.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:        "Resistência a dano Necrótico e Radiante.",
			PowerType:     domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Visão no Escuro", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:        "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType:     domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Mãos Curadoras", Edition: "5e", RaceID: &id,
			Description: "Você executa uma ação Usar Magia, toca uma criatura e joga um número de d4s igual ao seu Bônus de Proficiência. A criatura restaura PV igual ao total. Após usar esse traço, não pode usá-lo novamente até completar um Descanso Longo.",
			ActionType: "Ação", Range: "Toque",
			Effect:        "Cura Bônus de Proficiência × d4 PV. Recupera em Descanso Longo.",
			PowerType:     domain.PowerDaily, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Portador da Luz", Edition: "5e", RaceID: &id,
			Description: "Você conhece o truque Luz. Carisma é seu atributo de conjuração para ele.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:        "Truque Luz disponível com Carisma como atributo.",
			PowerType:     domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Revelação Celestial", Edition: "5e", RaceID: &id,
			Description: "No nível 3, como Ação Bônus você se transforma por 1 minuto. Escolha uma forma: Asas Celestiais (voo = deslocamento + dano extra radiante), Manto Necrótico (amedrontar inimigos + dano extra necrótico) ou Transfiguração Radiante (luz plena + dano radiante próximo). Recupera em Descanso Longo.",
			ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:        "Transformação de 1 minuto com bônus variados. Nível 3+. Recupera em Descanso Longo.",
			LevelScaling:  "Disponível a partir do nível 3 de personagem.",
			PowerType:     domain.PowerDaily, Level: 3, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "revelacao_aasimar",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Aasimar 5e: traços seedados")
}

// ── Anão ──────────────────────────────────────────────────────────────────────

func seedAnao5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Anão")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Anão)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 36 metros.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 36m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Resistência a Toxinas", Edition: "5e", RaceID: &id,
			Description: "Você tem Resistência a Dano Venenoso. Você também tem Vantagem nas salvaguardas que realizar para evitar ou encerrar a condição Envenenado.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Resistência a dano venenoso + vantagem vs Envenenado.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Tenacidade Anã", Edition: "5e", RaceID: &id,
			Description: "Seus Pontos de Vida máximos aumentam em 1, e novamente em 1 sempre que você alcança um nível de personagem.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "+1 PV máximo por nível de personagem.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Conhecimento de Pedras", Edition: "5e", RaceID: &id,
			Description: "Como Ação Bônus, você adquire Sismiconsciência com alcance de 18 metros por 10 minutos. Você deve estar tocando uma superfície de pedra (natural ou trabalhada). Usos = Bônus de Proficiência; recupera em Descanso Curto ou Longo.",
			ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:    "Sismiconsciência 18m por 10min. Usos = Bônus de Proficiência.",
			PowerType: domain.PowerEncounter, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Anão 5e: traços seedados")
}

// ── Draconato ─────────────────────────────────────────────────────────────────

func seedDraconato5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Draconato")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Ancestral Dracônico", Edition: "5e", RaceID: &id,
			Description: "Escolha um tipo de dragão: Ouro/Vermelho (Flamejante), Prata/Branco (Congelante), Latão/Cobre/Bronze (Ácido/Elétrico/Trovejante) ou Verde (Venenoso). Esta escolha determina o tipo do seu Sopro de Dragão e sua Resistência Elemental.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Define o tipo elemental do personagem.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "ancestral_draconato",
		},
		{
			Name: "Sopro de Dragão", Edition: "5e", RaceID: &id,
			Description: "Como Ação Bônus, você exala energia destrutiva do tipo do seu Ancestral Dracônico. A área e salvaguarda dependem do tipo: rajada contígua 4,5m ou linha 9m × 1,5m. CD = 8 + Bônus de Proficiência + modificador de Constituição. Usos = Bônus de Proficiência; recupera em Descanso Longo.",
			ActionType: "Ação Bônus", Range: "Área (veja desc.)",
			Effect:    "Dano elemental em área. CD = 8 + Bônus Prof + CON. Recupera em Descanso Longo.",
			LevelScaling: "Nível 1-4: 1d10. Nível 5-10: 2d10. Nível 11-16: 3d10. Nível 17+: 4d10.",
			PowerType: domain.PowerDaily, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Resistência Elemental", Edition: "5e", RaceID: &id,
			Description: "Você tem Resistência ao tipo de dano associado ao seu Ancestral Dracônico.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Resistência ao tipo elemental do ancestral dracônico.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Draconato 5e: traços seedados")
}

// ── Elfo ──────────────────────────────────────────────────────────────────────

func seedElfo5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Elfo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Elfo)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Ancestral das Fadas", Edition: "5e", RaceID: &id,
			Description: "Você tem Vantagem nas salvaguardas para evitar ou encerrar a condição Enfeitiçado. Magia não pode colocá-lo para dormir.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem vs Enfeitiçado. Imune a sono mágico.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Transe", Edition: "5e", RaceID: &id,
			Description: "Você não precisa dormir. Em vez disso, medita por 4 horas por dia. Após uma meditação de 4 horas, você obtém os mesmos benefícios de um Descanso Longo. Você pode realizar atividades leves durante a meditação.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Descanso Longo = 4h de meditação em vez de 8h de sono.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Linhagem Élfica", Edition: "5e", RaceID: &id,
			Description: "Escolha uma linhagem: Elfo Drow (truques de magia sombria), Elfo da Floresta (truques de natureza + movimento bônus) ou Alto Elfo (um truque de mago). Cada linhagem concede magias progressivas nos níveis 3 e 5.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Magias inatas baseadas na linhagem escolhida.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "linhagem_elfo",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Elfo 5e: traços seedados")
}

// ── Gnomo ─────────────────────────────────────────────────────────────────────

func seedGnomo5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Gnomo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Gnomo)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Astúcia Gnômica", Edition: "5e", RaceID: &id,
			Description: "Você tem Vantagem em todas as salvaguardas de Inteligência, Sabedoria e Carisma.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em salvaguardas de INT, SAB e CAR.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Linhagem Gnômica", Edition: "5e", RaceID: &id,
			Description: "Escolha uma linhagem: Gnomo da Floresta (comunicação com animais Pequenos, truques de ilusão) ou Gnomo das Rochas (proficiência com ferramentas de artesão, Conhecimento de Construção).",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Habilidades baseadas na linhagem gnômica escolhida.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "linhagem_gnomo",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Gnomo 5e: traços seedados")
}

// ── Golias ────────────────────────────────────────────────────────────────────

func seedGolias5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Golias")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Gigante Pequeno", Edition: "5e", RaceID: &id,
			Description: "Você conta como uma criatura Grande quando se trata de sua capacidade de carga e do peso que pode empurrar, arrastar ou erguer.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Conta como Grande para carga e empurrar/arrastar.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Potência da Pedra", Edition: "5e", RaceID: &id,
			Description: "Quando você sofre dano, pode usar sua Reação para reduzir esse dano em 1d12 + modificador de Constituição. Usos = Bônus de Proficiência; recupera em Descanso Longo.",
			ActionType: "Reação", Range: "Pessoal",
			Effect:    "Reduz dano em 1d12 + CON como Reação. Recupera em Descanso Longo.",
			PowerType: domain.PowerDaily, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Forma do Gigante", Edition: "5e", RaceID: &id,
			Description: "No nível 3, como Ação Bônus você assume uma forma de Gigante (Pedra, Gelo ou Tempestade, à escolha ao criar o personagem) por 1 minuto. Cada forma concede benefícios distintos de combate. Recupera em Descanso Longo.",
			ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:         "Transformação de Gigante por 1 minuto. Nível 3+. Recupera em Descanso Longo.",
			LevelScaling:   "Disponível a partir do nível 3 de personagem.",
			PowerType:      domain.PowerDaily, Level: 3, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "forma_gigante_golias",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Golias 5e: traços seedados")
}

// ── Humano ────────────────────────────────────────────────────────────────────

func seedHumano5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Humano")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Versátil", Edition: "5e", RaceID: &id,
			Description: "Você recebe um talento de Origem à sua escolha (veja o capítulo 5). Esse talento reflete sua versatilidade inata e experiência diversificada antes de sua vida de aventura. Exemplos: Curandeiro, Atacante Selvagem, Habilidoso, Músico, Sortudo.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Ganhe um talento de Origem adicional à escolha.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "versatil_humano",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Humano 5e: traços seedados")
}

// ── Orc ───────────────────────────────────────────────────────────────────────

func seedOrc5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Orc")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Orc)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 36 metros.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 36m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Pico de Adrenalina", Edition: "5e", RaceID: &id,
			Description: "Você pode executar a ação Correr como uma Ação Bônus. Ao fazê-lo, você ganha PV Temporários iguais ao seu Bônus de Proficiência. Usos = Bônus de Proficiência; recupera em Descanso Curto ou Longo.",
			ActionType: "Ação Bônus", Range: "Pessoal",
			Effect:    "Ação Correr + PV Temporários = Bônus de Proficiência. Recupera em Descanso Curto.",
			PowerType: domain.PowerEncounter, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Vigor Implacável", Edition: "5e", RaceID: &id,
			Description: "Quando você é reduzido a 0 Pontos de Vida, mas não morre imediatamente, você pode usar sua Reação para ficar com 1 Ponto de Vida. Após usar este traço, você não pode fazê-lo novamente até completar um Descanso Longo.",
			ActionType: "Reação", Range: "Pessoal",
			Effect:    "Ao cair a 0 PV, fica com 1 PV. Recupera em Descanso Longo.",
			PowerType: domain.PowerDaily, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Orc 5e: traços seedados")
}

// ── Pequenino ─────────────────────────────────────────────────────────────────

func seedPequenino5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Pequenino")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Corajoso", Edition: "5e", RaceID: &id,
			Description: "Você tem Vantagem nas salvaguardas que realizar para evitar ou encerrar a condição Amedrontado.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em salvaguardas contra Amedrontado.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Agilidade Pequenina", Edition: "5e", RaceID: &id,
			Description: "Você pode se mover pelo espaço de qualquer criatura que seja um tamanho maior que você, mas não pode parar no mesmo espaço.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Pode se mover pelo espaço de criaturas Médias ou maiores.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Sorte", Edition: "5e", RaceID: &id,
			Description: "Ao tirar 1 no dado D20 de um Teste de D20, você pode jogar novamente o dado e deve usar a nova jogada.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Rolar 1 no d20 permite rejogada obrigatória.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Furtividade Natural", Edition: "5e", RaceID: &id,
			Description: "Você pode executar a ação Esconder mesmo quando estiver encoberto apenas por uma criatura que seja pelo menos um tamanho maior que você.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Pode se Esconder usando uma criatura maior como cobertura.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Pequenino 5e: traços seedados")
}

// ── Tiferino ──────────────────────────────────────────────────────────────────

func seedTiferino5e(db *gorm.DB) {
	id, ok := getRace5e(db, "Tiferino")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Tiferino)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Legado Ínfero", Edition: "5e", RaceID: &id,
			Description: "Escolha um legado que determina sua linhagem ínfera: Infernal (diabólico), Abissal (demoníaco) ou Ctônico (sombrio). Cada legado concede um truque diferente no nível 1, e magias progressivas nos níveis 3 e 5, conjuradas uma vez por dia sem espaço de magia.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:         "Magias inatas por linhagem ínfera. Atributo de conjuração: INT, SAB ou CAR (escolha ao criar).",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "legado_tiferino",
		},
		{
			Name: "Presença Sobrenatural", Edition: "5e", RaceID: &id,
			Description: "Você conhece o truque Taumaturgia. Ao conjurar com este traço, use o mesmo atributo de conjuração do seu Legado Ínfero.",
			ActionType: "Passiva", Range: "Pessoal",
			Effect:    "Truque Taumaturgia disponível.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Tiferino 5e: traços seedados")
}