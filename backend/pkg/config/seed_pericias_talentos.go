package config

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// PERÍCIAS 4e
// ─────────────────────────────────────────────────────────────────────────────

func seedPericias(db *gorm.DB) {
	pericias := []domain.Pericia{
		{
			Name: "Acrobacia", Attribute: "Destreza", Edition: "4e",
			Description: "Manobras acrobáticas, equilíbrio e escapar de restrições.",
			Tooltip:     "Use para: equilibrar em superfícies instáveis, escapar de agarrões, rolar por debaixo de inimigos e realizar saltos acrobáticos. Dificuldade típica: 15 (moderado). Modificador: Destreza.",
		},
		{
			Name: "Arcana", Attribute: "Inteligência", Edition: "4e",
			Description: "Conhecimento sobre magia, itens mágicos e criaturas do Além.",
			Tooltip:     "Use para: identificar feitiços ativos, reconhecer itens mágicos, lembrar lendas sobre locais mágicos e navegar por planos de existência. Dificuldade típica: 15. Modificador: Inteligência.",
		},
		{
			Name: "Atletismo", Attribute: "Força", Edition: "4e",
			Description: "Escalar, nadar, saltar e resistir a esforço físico.",
			Tooltip:     "Use para: escalar paredes e penhascos, nadar contra correntezas, saltar grandes distâncias e resistir a trabalho físico prolongado. Dificuldade típica: 15. Modificador: Força.",
		},
		{
			Name: "Blefar", Attribute: "Carisma", Edition: "4e",
			Description: "Enganar criaturas, disfarçar intenções e criar distrações.",
			Tooltip:     "Use para: mentir de forma convincente, criar disfarces e enganar com linguagem corporal. Contraposto por Insight. Dificuldade típica: 10 + mod Insight do oponente. Modificador: Carisma.",
		},
		{
			Name: "Cura", Attribute: "Sabedoria", Edition: "4e",
			Description: "Estabilizar aliados, identificar doenças e venenos.",
			Tooltip:     "Use para: estabilizar aliado moribundo (DC 15), identificar doenças e venenos. Personagem treinado pode usar o pulso de cura de um aliado com Ação Padrão. Modificador: Sabedoria.",
		},
		{
			Name: "Destreza com Ladrão", Attribute: "Destreza", Edition: "4e",
			Description: "Arrombar fechaduras, desativar armadilhas e abrir bolsos.",
			Tooltip:     "Use para: arrombar cadeados (DC 20), desativar armadilhas mecânicas e abrir bolsos. Requer kit de ladrão para a maioria dos usos. Dificuldade típica: 20. Modificador: Destreza.",
		},
		{
			Name: "Diplomacia", Attribute: "Carisma", Edition: "4e",
			Description: "Negociar acordos, acalmar NPCs e influenciar atitudes.",
			Tooltip:     "Use para: fazer pedidos formais, negociar tréguas e melhorar a atitude de NPCs. Dificuldade típica: 20 para NPCs neutros; 30 para hostis. Modificador: Carisma.",
		},
		{
			Name: "Dungeon", Attribute: "Inteligência", Edition: "4e",
			Description: "Conhecimento sobre masmorras, subterrâneos e suas criaturas.",
			Tooltip:     "Use para: identificar armadilhas comuns, reconhecer arquitetura de masmorras, navegar em cavernas e conhecer ecologia de aberrações. Dificuldade típica: 15. Modificador: Inteligência.",
		},
		{
			Name: "Endurance", Attribute: "Constituição", Edition: "4e",
			Description: "Resistir a condições adversas: fome, veneno, frio e fadiga.",
			Tooltip:     "Use para: resistir a afogamento, privação de sono, venenos, doenças e extremos de temperatura. Testes são geralmente passivos ou durante descanso. Modificador: Constituição.",
		},
		{
			Name: "Furtividade", Attribute: "Destreza", Edition: "4e",
			Description: "Mover-se em silêncio e sem ser visto por inimigos.",
			Tooltip:     "Use para: passar despercebido por guardas, realizar emboscadas e escapar de perseguições. Contraposto pela Percepção passiva (10 + mod) do oponente. Modificador: Destreza.",
		},
		{
			Name: "História", Attribute: "Inteligência", Edition: "4e",
			Description: "Conhecimento sobre eventos passados, reinos e figuras históricas.",
			Tooltip:     "Use para: lembrar eventos históricos, identificar ruínas antigas, reconhecer brasões nobres e conhecer lendas de heróis do passado. Dificuldade típica: 15. Modificador: Inteligência.",
		},
		{
			Name: "Insight", Attribute: "Sabedoria", Edition: "4e",
			Description: "Detectar mentiras, ler intenções e perceber motivações ocultas.",
			Tooltip:     "Use para: detectar se alguém está mentindo, entender motivações de NPCs e perceber intenções hostis. Contraposto por Blefar. Modificador: Sabedoria.",
		},
		{
			Name: "Intimidação", Attribute: "Carisma", Edition: "4e",
			Description: "Ameaçar criaturas para extrair informações ou afastá-las.",
			Tooltip:     "Use para: intimidar NPCs para obter informações e forçar cooperação através do medo. Atenção: pode mudar a atitude do NPC para hostil após o uso. Modificador: Carisma.",
		},
		{
			Name: "Natureza", Attribute: "Sabedoria", Edition: "4e",
			Description: "Conhecimento sobre fauna, flora, clima e criaturas naturais.",
			Tooltip:     "Use para: identificar plantas medicinais e venenosas, prever clima, navegar em terrenos naturais e reconhecer criaturas selvagens. Dificuldade típica: 15. Modificador: Sabedoria.",
		},
		{
			Name: "Percepção", Attribute: "Sabedoria", Edition: "4e",
			Description: "Notar detalhes, detectar inimigos ocultos e sentir perigos.",
			Tooltip:     "Uma das perícias mais importantes! Detecta criaturas em Furtividade, armadilhas e passagens ocultas. A Percepção passiva (10 + mod) é usada constantemente pelo Mestre. Modificador: Sabedoria.",
		},
		{
			Name: "Religião", Attribute: "Inteligência", Edition: "4e",
			Description: "Conhecimento sobre divindades, rituais e criaturas divinas.",
			Tooltip:     "Use para: identificar símbolos religiosos, conhecer fraquezas de mortos-vivos e demônios, lembrar rituais sagrados e identificar artefatos divinos. Dificuldade típica: 15. Modificador: Inteligência.",
		},
		{
			Name: "Rua", Attribute: "Carisma", Edition: "4e",
			Description: "Conhecimento do submundo: gírias, rumores e como se virar nas cidades.",
			Tooltip:     "Use para: encontrar informantes, descobrir rumores urbanos, navegar em comunidades marginalizadas e identificar membros de guildas criminosas. Dificuldade típica: 15. Modificador: Carisma.",
		},
	}

	for _, p := range pericias {
		var existing domain.Pericia
		if db.Where("name = ? AND edition = ?", p.Name, p.Edition).First(&existing).Error != nil {
			db.Create(&p)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"attribute": p.Attribute, "description": p.Description, "tooltip": p.Tooltip,
			})
		}
	}
	log.Println("  ✓ Perícias 4e seedadas")
}

// ─────────────────────────────────────────────────────────────────────────────
// PERÍCIAS DISPONÍVEIS POR CLASSE (4e)
// ─────────────────────────────────────────────────────────────────────────────

func seedClassPericias(db *gorm.DB) {
	type entry struct {
		Name  string
		Count int
		List  string
	}
	classes := []entry{
		{"Bardo", 4, `["Arcana","Atletismo","Blefar","Diplomacia","Dungeon","História","Insight","Natureza","Percepção","Religião","Rua"]`},
		{"Bárbaro", 3, `["Acrobacia","Atletismo","Endurance","Cura","Intimidação","Natureza","Percepção"]`},
		{"Clérigo", 4, `["Arcana","Diplomacia","Cura","História","Insight","Religião"]`},
		{"Druida", 4, `["Arcana","Atletismo","Cura","Dungeon","História","Insight","Natureza","Percepção"]`},
		{"Feiticeiro", 4, `["Arcana","Blefar","Dungeon","História","Insight","Natureza"]`},
		{"Guardião", 3, `["Atletismo","Dungeon","Endurance","Cura","Intimidação","Natureza","Percepção"]`},
		{"Guerreiro", 3, `["Atletismo","Diplomacia","Endurance","Cura","Intimidação","Rua"]`},
		{"Invocador", 4, `["Arcana","Diplomacia","História","Insight","Natureza","Religião"]`},
		{"Ladino", 4, `["Acrobacia","Blefar","Destreza com Ladrão","Diplomacia","Dungeon","Furtividade","Insight","Intimidação","Percepção","Rua"]`},
		{"Mago", 4, `["Arcana","Dungeon","História","Insight","Natureza","Religião"]`},
		{"Monge", 4, `["Acrobacia","Atletismo","Endurance","Insight","Percepção","Religião","Rua"]`},
		{"Paladino", 4, `["Diplomacia","Endurance","Cura","Insight","Intimidação","Religião"]`},
		{"Patrulheiro", 4, `["Acrobacia","Atletismo","Dungeon","Endurance","Furtividade","Natureza","Percepção"]`},
		{"Psionista", 4, `["Arcana","Dungeon","História","Insight","Percepção","Religião"]`},
		{"Rastreador", 4, `["Acrobacia","Atletismo","Dungeon","Endurance","Furtividade","Natureza","Percepção"]`},
		{"Vingador", 4, `["Acrobacia","Atletismo","Endurance","Insight","Intimidação","Percepção","Religião","Rua"]`},
		{"Xamã", 4, `["Arcana","Cura","Dungeon","História","Insight","Natureza","Percepção","Religião"]`},
	}

	for _, e := range classes {
		var cls domain.Class
		if db.Where("name = ? AND edition = ?", e.Name, "4e").First(&cls).Error == nil {
			db.Model(&cls).Updates(map[string]interface{}{
				"trained_skills_count": e.Count,
				"available_skills":     e.List,
				"talentos_count":       2,
			})
		}
	}
	log.Println("  ✓ Perícias por classe 4e atualizadas")
}

// ─────────────────────────────────────────────────────────────────────────────
// BÔNUS DE PERÍCIAS POR RAÇA (4e)
// ─────────────────────────────────────────────────────────────────────────────

func seedRacePericias(db *gorm.DB) {
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

	for _, e := range races {
		var race domain.Race
		if db.Where("name = ? AND edition = ?", e.Name, "4e").First(&race).Error == nil {
			db.Model(&race).Updates(map[string]interface{}{
				"bonus_trained_skills": e.BonusTrainedSkills,
				"bonus_talentos":       e.BonusTalentos,
				"bonus_skill_values":   e.BonusSkillValues,
			})
		}
	}
	log.Println("  ✓ Bônus de perícias por raça 4e atualizados")
}

// ─────────────────────────────────────────────────────────────────────────────
// TALENTOS 4e
// ─────────────────────────────────────────────────────────────────────────────

func seedTalentos(db *gorm.DB) {
	talentos := []domain.Talento{
		// ── DEFESA ──────────────────────────────────────────────────────────
		{
			Name: "Vigor", Edition: "4e", Category: "Defesa", Prerequisite: "",
			Description: "+5 PV no nível heroico.",
			Tooltip:     "+5 PV (nível 1-10), +10 PV (nível 11-20), +15 PV (nível 21+). Um dos talentos mais versáteis do jogo — útil para qualquer classe.",
		},
		{
			Name: "Vontade de Ferro", Edition: "4e", Category: "Defesa", Prerequisite: "",
			Description: "+2 na defesa de Vontade.",
			Tooltip:     "Protege contra ataques mentais, encantamentos e poderes psíquicos. Essencial para classes sem bônus natural de Vontade como Guerreiro e Bárbaro.",
		},
		{
			Name: "Grande Fortitude", Edition: "4e", Category: "Defesa", Prerequisite: "",
			Description: "+2 na defesa de Fortitude.",
			Tooltip:     "Protege contra ataques físicos pesados, venenos e doenças. Ótimo para conjuradores com Constituição baixa.",
		},
		{
			Name: "Reflexos Relâmpago", Edition: "4e", Category: "Defesa", Prerequisite: "",
			Description: "+2 na defesa de Reflexos.",
			Tooltip:     "Protege contra explosões, armadilhas e ataques de área. Ótimo para personagens com Destreza baixa como Guerreiros com armadura pesada.",
		},
		{
			Name: "Mobilidade Defensiva", Edition: "4e", Category: "Defesa", Prerequisite: "",
			Description: "+2 CA contra ataques de oportunidade.",
			Tooltip:     "Ideal para personagens que se movem muito em combate ou que precisam alcançar alvos distantes. Permite reposicionamento mais seguro.",
		},
		{
			Name: "Durável", Edition: "4e", Category: "Defesa", Prerequisite: "",
			Description: "+2 pulsos de cura por dia.",
			Tooltip:     "Você recebe 2 pulsos de cura adicionais por dia. Excelente para personagens que sobrevivem a muitos encontros consecutivos sem descanso prolongado.",
		},
		// ── COMBATE ─────────────────────────────────────────────────────────
		{
			Name: "Foco em Arma", Edition: "4e", Category: "Combate",
			Prerequisite: "Proficiência com o tipo de arma",
			Description:  "+1 no dano com um tipo de arma escolhido.",
			Tooltip:      "Escolha um grupo de armas (espadas, machados, arcos...). +1 dano (nível 1), +2 (nível 11), +3 (nível 21). Útil para qualquer classe que ataca frequentemente.",
		},
		{
			Name: "Especialização em Arma", Edition: "4e", Category: "Combate",
			Prerequisite: "Proficiência com o tipo de arma",
			Description:  "+1 nas jogadas de ataque com um tipo de arma.",
			Tooltip:      "Escolha um grupo de armas. +1 ataque (nível 1), +2 (nível 5), +3 (nível 15). Aumenta a confiabilidade dos ataques — muito valioso em todos os níveis.",
		},
		{
			Name: "Foco em Implemento", Edition: "4e", Category: "Combate",
			Prerequisite: "Proficiência com o implemento",
			Description:  "+1 no dano com um implemento escolhido.",
			Tooltip:      "Escolha um tipo de implemento (varinha, báculo, orbe...). +1 dano (nível 1), +2 (nível 11), +3 (nível 21). Equivalente ao Foco em Arma para conjuradores.",
		},
		{
			Name: "Especialização em Implemento", Edition: "4e", Category: "Combate",
			Prerequisite: "Proficiência com o implemento",
			Description:  "+1 nas jogadas de ataque com um implemento.",
			Tooltip:      "Escolha um implemento. +1 ataque (nível 1), +2 (nível 5), +3 (nível 15). Equivalente à Especialização em Arma para conjuradores — muito recomendado.",
		},
		{
			Name: "Ataque Poderoso", Edition: "4e", Category: "Combate",
			Prerequisite: "Força 13",
			Description:  "-2 ataque, +4 dano (+6 com arma de duas mãos).",
			Tooltip:      "Aceite -2 na jogada de ataque para receber +4 dano (+6 com armas de duas mãos). Melhor usado quando você tem alta CA sobre o inimigo e dificilmente erra.",
		},
		{
			Name: "Combate com Duas Armas", Edition: "4e", Category: "Combate",
			Prerequisite: "Destreza 13",
			Description:  "Reduz a penalidade da mão inábil de -4 para -2.",
			Tooltip:      "A penalidade na mão inábil cai de -4 para -2 (ao usar arma leve). Essencial para Patrulheiros e outros combatentes de duas armas — pré-requisito para Defesa com Duas Armas.",
		},
		{
			Name: "Defesa com Duas Armas", Edition: "4e", Category: "Combate",
			Prerequisite: "Combate com Duas Armas",
			Description:  "+1 CA e Reflexos ao empunhar duas armas.",
			Tooltip:      "Enquanto você empunha uma arma em cada mão, recebe +1 CA e +1 Reflexos. Combina diretamente com Combate com Duas Armas — pegue os dois juntos.",
		},
		{
			Name: "Combate às Cegas", Edition: "4e", Category: "Combate",
			Prerequisite: "",
			Description:  "Relança ataques que falham por ocultação.",
			Tooltip:      "Quando um ataque falha por ocultação (leve ou total), você relança a jogada uma vez. Excelente contra inimigos Furtivos, Invisíveis ou em combates no escuro.",
		},
		{
			Name: "Iniciativa Aprimorada", Edition: "4e", Category: "Combate",
			Prerequisite: "",
			Description:  "+4 nas jogadas de iniciativa.",
			Tooltip:      "+4 de bônus nas jogadas de iniciativa. Agir antes dos inimigos pode mudar o resultado de um encontro — especialmente poderoso para classes de burst e controle.",
		},
		{
			Name: "Reflexos em Combate", Edition: "4e", Category: "Combate",
			Prerequisite: "Destreza 13",
			Description:  "+2 nas jogadas de ataque de oportunidade.",
			Tooltip:      "+2 nos ataques de oportunidade. Ideal para defensores que patrulham o campo de batalha protegendo aliados e punindo inimigos que tentam se mover livremente.",
		},
		{
			Name: "Crítico Devastador", Edition: "4e", Category: "Combate",
			Prerequisite: "Nível 15",
			Description:  "+1d6 de dano extra em acertos críticos.",
			Tooltip:      "Em acertos críticos (resultado máximo no dado de ataque), causa +1d6 do mesmo tipo de dano do ataque. O dano adicional escala com itens mágicos de alto nível.",
		},
		// ── PERÍCIAS ────────────────────────────────────────────────────────
		{
			Name: "Foco em Perícia", Edition: "4e", Category: "Perícia",
			Prerequisite: "Treinado na perícia",
			Description:  "+3 em uma perícia treinada escolhida.",
			Tooltip:      "Escolha uma perícia treinada e receba +3 permanente nela. Empilha com o bônus de treinamento (+5) e de atributo. Ótimo para criar especialistas extremamente eficientes.",
		},
		{
			Name: "Treino em Perícia", Edition: "4e", Category: "Perícia",
			Prerequisite: "",
			Description:  "Treine qualquer perícia fora da lista da sua classe.",
			Tooltip:      "Fica treinado em qualquer perícia do jogo, recebendo o bônus de treinamento (+5). Muito versátil para cobrir lacunas do grupo — ex: Guerreiro treinando Arcana.",
		},
		{
			Name: "Sentido de Perigo", Edition: "4e", Category: "Perícia",
			Prerequisite: "",
			Description:  "+2 iniciativa e +5 Percepção passiva contra emboscadas.",
			Tooltip:      "+2 nas jogadas de iniciativa e +5 na Percepção passiva para detectar emboscadas. Se qualquer aliado adjacente não estiver surpreso, você também não fica. Excelente para batedores.",
		},
		// ── MAGIA ───────────────────────────────────────────────────────────
		{
			Name: "Lançador de Rituais", Edition: "4e", Category: "Magia",
			Prerequisite: "Inteligência 13 ou Sabedoria 13",
			Description:  "Pode aprender e realizar rituais mágicos.",
			Tooltip:      "Permite aprender e realizar rituais com componentes adequados (ouro). Rituais possibilitam efeitos poderosos fora do combate: teletransporte, comunicação à distância, detecção mágica. Requer livro de rituais.",
		},
		// ── ARMADURA ────────────────────────────────────────────────────────
		{
			Name: "Proficiência em Armadura — Couro", Edition: "4e", Category: "Armadura",
			Prerequisite: "",
			Description:  "Proficiência em armadura de couro.",
			Tooltip:      "Sem proficiência, usar uma armadura impõe -2 nas jogadas de ataque e testes de atributo. Couro oferece CA moderada com boa mobilidade — ótimo ponto de entrada para conjuradores.",
		},
		{
			Name: "Proficiência em Armadura — Malha", Edition: "4e", Category: "Armadura",
			Prerequisite: "",
			Description:  "Proficiência em armadura de malha.",
			Tooltip:      "Malha oferece boa CA com penalidade moderada. Recomendado para classes de suporte ou controle que ficam no front mas não têm proficiência nativa.",
		},
		{
			Name: "Proficiência em Armadura — Escamas", Edition: "4e", Category: "Armadura",
			Prerequisite: "Prof. em armadura de malha",
			Description:  "Proficiência em armadura de escamas.",
			Tooltip:      "Escamas oferecem excelente CA com penalidade significativa na velocidade. Para personagens que priorizam defesa sobre mobilidade.",
		},
		{
			Name: "Proficiência em Armadura — Placa", Edition: "4e", Category: "Armadura",
			Prerequisite: "Prof. em armadura de escamas",
			Description:  "Proficiência em armadura de placa completa.",
			Tooltip:      "A melhor proteção disponível. Reduz significativamente velocidade e testes de Destreza. Ideal para Guerreiros e Paladinos focados em absorver dano.",
		},
	}

	for _, t := range talentos {
		var existing domain.Talento
		if db.Where("name = ? AND edition = ?", t.Name, t.Edition).First(&existing).Error != nil {
			db.Create(&t)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description":  t.Description,
				"prerequisite": t.Prerequisite,
				"category":     t.Category,
				"tooltip":      t.Tooltip,
			})
		}
	}
	log.Println("  ✓ Talentos 4e seedados")
}