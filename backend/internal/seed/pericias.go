package seed

import (
	"log"
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

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
	log.Println("  ✓ Perícias seedadas")
}

// ── Perícias disponíveis por classe ──────────────────────────────────────────

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
				"talentos_count":       2, // sua campanha usa 2
			})
		}
	}
	log.Println("  ✓ Perícias por classe atualizadas")
}

// ── Bônus de perícias por raça ────────────────────────────────────────────────

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
	log.Println("  ✓ Bônus de perícias por raça atualizados")
}