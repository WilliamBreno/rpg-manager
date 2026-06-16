package seed

import (
	"log"
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	log.Println("🌱 Rodando seed...")
	seedClasses(db)
	seedSkills(db)
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
	classes := []domain.Class{
		{Name: "Bardo", Edition: "4e", Description: "Líder arcano que usa música e conhecimento como armas.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Bárbaro", Edition: "4e", Description: "Combatente primitivo de força e fúria selvagem.", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 8, FortBonus: 2, IsDefault: true},
		{Name: "Clérigo", Edition: "4e", Description: "Líder divino que canaliza o poder de sua divindade.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 1, IsDefault: true},
		{Name: "Druida", Edition: "4e", Description: "Controlador primitivo com maestria da natureza.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Feiticeiro", Edition: "4e", Description: "Arcano que canaliza poder dracônico ou do caos.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Guardião", Edition: "4e", Description: "Defensor primitivo e protetor do mundo natural.", BaseHP: 17, HPPerLevel: 7, SurgesPerDay: 9, FortBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Guerreiro", Edition: "4e", Description: "Defensor marcial especialista em combate corpo a corpo.", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 9, FortBonus: 2, IsDefault: true},
		{Name: "Invocador", Edition: "4e", Description: "Controlador divino que invoca o poder dos deuses.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Ladino", Edition: "4e", Description: "Agressor furtivo especializado em ataques precisos.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, ReflBonus: 2, IsDefault: true},
		{Name: "Mago", Edition: "4e", Description: "Controlador arcano de grande poder mágico.", BaseHP: 10, HPPerLevel: 4, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Monge", Edition: "4e", Description: "Agressor psiônico de disciplina e artes marciais.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Paladino", Edition: "4e", Description: "Defensor divino e campeão sagrado.", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 10, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Patrulheiro", Edition: "4e", Description: "Agressor marcial e batedor das regiões ermas.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, FortBonus: 1, ReflBonus: 1, IsDefault: true},
		{Name: "Psionista", Edition: "4e", Description: "Controlador psiônico de grande poder mental.", BaseHP: 12, HPPerLevel: 4, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Rastreador", Edition: "4e", Description: "Controlador primitivo especialista em emboscadas.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Vingador", Edition: "4e", Description: "Agressor divino executor da vontade dos deuses.", BaseHP: 14, HPPerLevel: 6, SurgesPerDay: 7, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Xamã", Edition: "4e", Description: "Líder primitivo que canaliza espíritos do mundo natural.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, FortBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Bárbaro", Edition: "5e", Description: "Guerreiro furioso de força bruta.", HitDie: 12, IsDefault: true},
		{Name: "Bardo", Edition: "5e", Description: "Conjurador versátil com magia e música.", HitDie: 8, IsDefault: true},
		{Name: "Bruxo", Edition: "5e", Description: "Conjurador com poder concedido por um patrono.", HitDie: 8, IsDefault: true},
		{Name: "Clérigo", Edition: "5e", Description: "Servidor divino com poderes de cura e combate.", HitDie: 8, IsDefault: true},
		{Name: "Druida", Edition: "5e", Description: "Guardião da natureza com magia elemental.", HitDie: 8, IsDefault: true},
		{Name: "Feiticeiro", Edition: "5e", Description: "Conjurador de magia inata e poderosa.", HitDie: 6, IsDefault: true},
		{Name: "Guerreiro", Edition: "5e", Description: "Mestre do combate e das armas.", HitDie: 10, IsDefault: true},
		{Name: "Ladino", Edition: "5e", Description: "Especialista furtivo e hábil.", HitDie: 8, IsDefault: true},
		{Name: "Mago", Edition: "5e", Description: "Estudioso da magia arcana.", HitDie: 6, IsDefault: true},
		{Name: "Monge", Edition: "5e", Description: "Combatente de artes marciais e energia mística.", HitDie: 8, IsDefault: true},
		{Name: "Paladino", Edition: "5e", Description: "Guerreiro sagrado jurado a um ideal.", HitDie: 10, IsDefault: true},
		{Name: "Patrulheiro", Edition: "5e", Description: "Caçador e rastreador das regiões selvagens.", HitDie: 10, IsDefault: true},
		{Name: "Xamã", Edition: "5e", Description: "Líder espiritual conectado ao mundo natural.", HitDie: 8, IsDefault: true},
	}
	for _, c := range classes {
		var existing domain.Class
		if db.Where("name = ? AND edition = ?", c.Name, c.Edition).First(&existing).Error != nil {
			db.Create(&c)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"base_hp": c.BaseHP, "hp_per_level": c.HPPerLevel,
				"surges_per_day": c.SurgesPerDay, "fort_bonus": c.FortBonus,
				"refl_bonus": c.ReflBonus, "will_bonus": c.WillBonus,
				"hit_die": c.HitDie, "description": c.Description,
			})
		}
	}
	log.Println("  ✓ Classes seedadas")
}

func seedSkills(db *gorm.DB) {
	seedBardoSkills(db)
	seedMongeSkills(db)
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
		log.Println("  ✗ Monge 4e não encontrado"); return
	}
	id := cls.ID

	skills := []domain.Skill{
		{
			Name: "Sequência de Golpes Centrada", Edition: "4e", ClassID: &id,
			Description: "Seus punhos ficam ofuscados enquanto você desfere um ataque inicial seguido de outro, ajustando posições dos inimigos ao seu favor.",
			Keywords: "Psiônico", ActionType: "Nenhuma Ação (Especial)", Range: "Corpo a corpo 1",
			Target: "Uma ou duas criaturas",
			Effect: "Gatilho: O monge atinge uma criatura durante seu turno. O alvo sofre dano igual a 2 + mod Sabedoria. O monge conduz o alvo 1 quadrado adjacente.",
			Special: "Pode ser usado apenas uma vez por rodada. Nível 21: cada inimigo adjacente ao monge.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "sequencia_golpes_monge",
		},
		{
			Name: "Sequência de Golpes do Punho de Pedra", Edition: "4e", ClassID: &id,
			Description: "Você golpeia outro inimigo após seu primeiro ataque, um lembre casual da sua grande força.",
			Keywords: "Psiônico", ActionType: "Nenhuma Ação (Especial)", Range: "Corpo a corpo 1",
			Target: "Uma criatura",
			Effect: "Gatilho: O monge atinge uma criatura durante seu turno. O alvo sofre dano igual a 3 + mod Sabedoria. Caso não tenha sido a criatura atingida pelo gatilho, o dano aumenta em 2 (nível 11 e 21).",
			Special: "Pode ser usado apenas uma vez por rodada. Nível 21: cada inimigo adjacente ao monge.",
			PowerType: domain.PowerUnlimited, Level: 1,
			IsClassFeature: true, RequiresChoice: true, ChoiceGroup: "sequencia_golpes_monge",
		},
	}

	for _, s := range skills {
		upsertSkill(db, s, id)
	}
	log.Printf("  ✓ Monge 4e: %d características processadas", len(skills))
}