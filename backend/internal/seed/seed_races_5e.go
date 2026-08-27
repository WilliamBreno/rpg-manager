package seed

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// ESPÉCIES 5e (PHB 2024) — portado de pkg/config/seed_5e.go (código morto)
//
// ATENÇÃO: as raças 5e já existentes no banco (criadas por um seed legado)
// usam nomenclatura antiga ("Elfo da Floresta", "Alto Elfo", "Tiefling",
// "Anão das Montanhas"...). Estas aqui usam a nomenclatura unificada do
// PHB 2024 ("Elfo", "Anão", "Tiferino"...) — são raças DIFERENTES que vão
// coexistir nos dropdowns. Apenas estas (PHB 2024) têm traços raciais
// completos seedados abaixo.
// ─────────────────────────────────────────────────────────────────────────────

// seedRaces5e cria/atualiza as 10 espécies do PHB 2024.
func seedRaces5e(db *gorm.DB) {
	type data struct {
		Name, Description string
		Speed              int
	}

	// Deslocamento (PHB 2024, cap. 4): 9 metros = 30ft para a maioria das
	// espécies; Golias é a única exceção, com 10,5 metros = 35ft. Antes disso
	// todas as 10 espécies estavam hardcoded em 30 — bug real que fazia
	// personagens Golias exportarem "Deslocamento: 30" na ficha, errado.
	races := []data{
		{"Aasimar", "Mortais com uma centelha dos Planos Superiores. Carregam herança celestial que se manifesta em resistências e poderes de luz e cura.", 30},
		{"Anão", "Forjados da terra por Moradin. Resilientes como montanhas, com afinidade por pedra, metal e vida subterrânea. Vivem cerca de 350 anos.", 30},
		{"Draconato", "Humanoides de herança dracônica, com escamas, sopro elemental e resistência ao tipo de dano do seu ancestral dragão.", 30},
		{"Elfo", "Seres feéricos com sentidos aguçados, resistência à magia, e a capacidade de meditar em transe no lugar de dormir. Vivem séculos.", 30},
		{"Gnomo", "Pequenos humanoides com astúcia mágica natural, visão no escuro e resistência a magias que afetam a mente.", 30},
		{"Golias", "Gigantes descendentes de gigantes ancestrais, com forte ancestralidade gigante, capazes de resistir a danos devastadores e assumir tamanho Grande temporariamente.", 35},
		{"Humano", "A espécie mais versátil do multiverso. Recebem um talento de Origem adicional e se adaptam a qualquer caminho.", 30},
		{"Orc", "Humanoides robustos abençoados por Gruumsh. Determinados, com visão no escuro poderosa e a capacidade de continuar lutando ao cair.", 30},
		{"Pequenino", "Pequenos humanoides corajosos com sorte inata, furtividade natural e agilidade para se mover entre criaturas maiores.", 30},
		{"Tiferino", "Humanoides com sangue ínfero. Possuem visão no escuro e um legado mágico que reflete sua linhagem — diabólica, demoníaca ou outra.", 30},
	}

	for _, r := range races {
		var existing domain.Race
		if db.Where("name = ? AND edition = ?", r.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Race{
				Name:        r.Name,
				Edition:     "5e",
				Description: r.Description,
				Speed:       r.Speed,
				IsDefault:   true,
			})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description": r.Description,
				"speed":       r.Speed,
			})
		}
	}
	log.Println("  ✓ Espécies 5e (PHB 2024): 10 espécies seedadas")
}

// seedRaceFeatures5e semeia os traços raciais de todas as espécies 5e PHB 2024.
func seedRaceFeatures5e(db *gorm.DB) {
	seedAasimar5e(db)
	seedAnao5e(db)
	seedDraconato5ePHB(db)
	seedElfo5ePHB(db)
	seedGnomo5ePHB(db)
	seedGolias5ePHB(db)
	seedHumano5ePHB(db)
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
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:        "Resistência a dano Necrótico e Radiante.",
			PowerType:     domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Visão no Escuro", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:        "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType:     domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Mãos Curadoras", Edition: "5e", RaceID: &id,
			Description: "Você executa uma ação Usar Magia, toca uma criatura e joga um número de d4s igual ao seu Bônus de Proficiência. A criatura restaura PV igual ao total. Após usar esse traço, não pode usá-lo novamente até completar um Descanso Longo.",
			ActionType:  "Ação", Range: "Toque",
			Effect:        "Cura Bônus de Proficiência × d4 PV. Recupera em Descanso Longo.",
			PowerType:     domain.PowerDaily, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Portador da Luz", Edition: "5e", RaceID: &id,
			Description: "Você conhece o truque Luz. Carisma é seu atributo de conjuração para ele.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:        "Truque Luz disponível com Carisma como atributo.",
			PowerType:     domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Revelação Celestial", Edition: "5e", RaceID: &id,
			Description: "No nível 3, como Ação Bônus você se transforma por 1 minuto. Escolha uma forma: Asas Celestiais (voo = deslocamento + dano extra radiante), Manto Necrótico (amedrontar inimigos + dano extra necrótico) ou Transfiguração Radiante (luz plena + dano radiante próximo). Recupera em Descanso Longo.",
			ActionType:  "Ação Bônus", Range: "Pessoal",
			Effect:         "Transformação de 1 minuto com bônus variados. Nível 3+. Recupera em Descanso Longo.",
			LevelScaling:   "Disponível a partir do nível 3 de personagem.",
			PowerType:      domain.PowerDaily, Level: 3, IsRaceFeature: true,
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
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 36m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Resistência a Toxinas", Edition: "5e", RaceID: &id,
			Description: "Você tem Resistência a Dano Venenoso. Você também tem Vantagem nas salvaguardas que realizar para evitar ou encerrar a condição Envenenado.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Resistência a dano venenoso + vantagem vs Envenenado.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Tenacidade Anã", Edition: "5e", RaceID: &id,
			Description: "Seus Pontos de Vida máximos aumentam em 1, e novamente em 1 sempre que você alcança um nível de personagem.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "+1 PV máximo por nível de personagem.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Conhecimento de Pedras", Edition: "5e", RaceID: &id,
			Description: "Como Ação Bônus, você adquire Sismiconsciência com alcance de 18 metros por 10 minutos. Você deve estar tocando uma superfície de pedra (natural ou trabalhada). Usos = Bônus de Proficiência; recupera em Descanso Curto ou Longo.",
			ActionType:  "Ação Bônus", Range: "Pessoal",
			Effect:    "Sismiconsciência 18m por 10min. Usos = Bônus de Proficiência.",
			PowerType: domain.PowerEncounter, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Anão 5e: traços seedados")
}

// ── Draconato (PHB 2024) ──────────────────────────────────────────────────────
// Sufixo "PHB" no nome da função para não colidir com qualquer outra função
// de mesmo propósito que já exista no projeto.

func seedDraconato5ePHB(db *gorm.DB) {
	id, ok := getRace5e(db, "Draconato")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Ancestral Dracônico", Edition: "5e", RaceID: &id,
			Description: "Escolha um tipo de dragão: Ouro/Vermelho (Flamejante), Prata/Branco (Congelante), Latão/Cobre/Bronze (Ácido/Elétrico/Trovejante) ou Verde (Venenoso). Esta escolha determina o tipo do seu Sopro de Dragão e sua Resistência Elemental.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Define o tipo elemental do personagem.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "ancestral_draconato_phb",
		},
		{
			Name: "Sopro de Dragão (PHB)", Edition: "5e", RaceID: &id,
			Description: "Como Ação Bônus, você exala energia destrutiva do tipo do seu Ancestral Dracônico. A área e salvaguarda dependem do tipo: rajada contígua 4,5m ou linha 9m × 1,5m. CD = 8 + Bônus de Proficiência + modificador de Constituição. Usos = Bônus de Proficiência; recupera em Descanso Longo.",
			ActionType:  "Ação Bônus", Range: "Área (veja desc.)",
			Effect:       "Dano elemental em área. CD = 8 + Bônus Prof + CON. Recupera em Descanso Longo.",
			LevelScaling: "Nível 1-4: 1d10. Nível 5-10: 2d10. Nível 11-16: 3d10. Nível 17+: 4d10.",
			PowerType:    domain.PowerDaily, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Resistência Elemental", Edition: "5e", RaceID: &id,
			Description: "Você tem Resistência ao tipo de dano associado ao seu Ancestral Dracônico.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Resistência ao tipo elemental do ancestral dracônico.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Visão no Escuro (Draconato)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Voo Dracônico", Edition: "5e", RaceID: &id,
			Description: "No nível 5 de personagem, você pode canalizar magia dracônica para voar temporariamente. Como uma Ação Bônus, você cria asas espectrais nas costas que duram 10 minutos ou até você as retrair (nenhuma ação necessária) ou ficar Incapacitado. Pela duração, você tem Deslocamento de Voo igual ao seu Deslocamento. Após usar esse traço, não pode usá-lo novamente até completar um Descanso Longo.",
			ActionType:  "Ação Bônus", Range: "Pessoal",
			Effect:       "Deslocamento de Voo igual ao Deslocamento normal por 10 minutos.",
			LevelScaling: "Disponível a partir do nível 5 de personagem.",
			PowerType:    domain.PowerDaily, Level: 5, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Draconato 5e (PHB): traços seedados")
}

// ── Elfo (PHB 2024) ───────────────────────────────────────────────────────────

func seedElfo5ePHB(db *gorm.DB) {
	id, ok := getRace5e(db, "Elfo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Elfo)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Ancestralidade Feérica", Edition: "5e", RaceID: &id,
			Description: "Você tem Vantagem ao realizar salvaguardas para evitar ou encerrar a condição Enfeitiçado.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em salvaguardas contra Enfeitiçado.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Sentidos Aguçados", Edition: "5e", RaceID: &id,
			Description: "Você tem proficiência na perícia Intuição, Percepção ou Sobrevivência (escolha uma).",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Proficiência em 1 perícia à escolha entre Intuição, Percepção ou Sobrevivência.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "sentidos_agucados_elfo",
		},
		{
			Name: "Transe", Edition: "5e", RaceID: &id,
			Description: "Você não precisa dormir. Em vez disso, medita por 4 horas por dia. Após uma meditação de 4 horas, você obtém os mesmos benefícios de um Descanso Longo. Você pode realizar atividades leves durante a meditação.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Descanso Longo = 4h de meditação em vez de 8h de sono.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Linhagem Élfica", Edition: "5e", RaceID: &id,
			Description: "Escolha uma linhagem: Elfo Drow (truques de magia sombria), Elfo da Floresta (truques de natureza + movimento bônus) ou Alto Elfo (um truque de mago). Cada linhagem concede magias progressivas nos níveis 3 e 5.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Magias inatas baseadas na linhagem escolhida.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "linhagem_elfo_phb",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Elfo 5e (PHB): traços seedados")
}

// ── Gnomo (PHB 2024) ──────────────────────────────────────────────────────────

func seedGnomo5ePHB(db *gorm.DB) {
	id, ok := getRace5e(db, "Gnomo")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Visão no Escuro (Gnomo)", Edition: "5e", RaceID: &id,
			Description: "Você tem Visão no Escuro com alcance de 18 metros.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Astúcia Gnômica", Edition: "5e", RaceID: &id,
			Description: "Você tem Vantagem em todas as salvaguardas de Inteligência, Sabedoria e Carisma.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em salvaguardas de INT, SAB e CAR.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Linhagem Gnômica", Edition: "5e", RaceID: &id,
			Description: "Escolha uma linhagem: Gnomo da Floresta (comunicação com animais Pequenos, truques de ilusão) ou Gnomo das Rochas (proficiência com ferramentas de artesão, Conhecimento de Construção).",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Habilidades baseadas na linhagem gnômica escolhida.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "linhagem_gnomo_phb",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Gnomo 5e (PHB): traços seedados")
}

// ── Golias (PHB 2024) ─────────────────────────────────────────────────────────

func seedGolias5ePHB(db *gorm.DB) {
	id, ok := getRace5e(db, "Golias")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Ancestralidade Gigante", Edition: "5e", RaceID: &id,
			Description: "Você é descendente de Gigantes. Escolha um dos seguintes benefícios sobrenaturais: Arrepio do Gelo (Gigante do Gelo — ao causar dano com um ataque, também inflige 1d6 de dano Gélido e reduz o Deslocamento do alvo em 3m até o início do seu próximo turno), Queimadura de Fogo (Gigante de Fogo — também inflige 1d10 de dano Ígneo ao causar dano com um ataque), Resistência da Pedra (Gigante da Pedra — como Reação ao sofrer dano, jogue 1d12 + modificador de Constituição e reduza o dano por esse total), Salto da Nuvem (Gigante das Nuvens — como Ação Bônus, teleporta-se até 9m para um espaço desocupado à vista), Tombo da Colina (Gigante da Colina — ao causar dano a uma criatura Grande ou menor com um ataque, pode impor a condição Caído) ou Trovão da Tempestade (Gigante da Tempestade — como Reação ao sofrer dano de uma criatura a até 18m, causa 1d8 de dano Trovejante a ela). Usos = Bônus de Proficiência; recupera em Descanso Longo.",
			ActionType:  "Variado (veja desc.)", Range: "Pessoal",
			Effect:         "Um benefício sobrenatural à escolha entre 6 opções ligadas a um tipo de gigante ancestral.",
			PowerType:      domain.PowerDaily, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "ancestralidade_gigante_golias",
		},
		{
			Name: "Forma Grande", Edition: "5e", RaceID: &id,
			Description: "A partir do nível 5 de personagem, você pode alterar seu tamanho para Grande como uma Ação Bônus, se estiver em um espaço grande o suficiente. Essa transformação se mantém por 10 minutos ou até você a encerrar (nenhuma ação necessária). Pela duração, você tem Vantagem em testes de Força, e seu Deslocamento aumenta em 3 metros. Após usar este traço, não pode utilizá-lo novamente até completar um Descanso Longo.",
			ActionType:  "Ação Bônus", Range: "Pessoal",
			Effect:       "Tamanho Grande por 10 minutos: Vantagem em testes de Força, +3m de Deslocamento.",
			LevelScaling: "Disponível a partir do nível 5 de personagem.",
			PowerType:    domain.PowerDaily, Level: 5, IsRaceFeature: true,
		},
		{
			Name: "Porte Poderoso", Edition: "5e", RaceID: &id,
			Description: "Você tem Vantagem em qualquer teste de atributo que realizar para encerrar a condição Imobilizado. Você também conta como um tamanho maior ao determinar sua capacidade de carga.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Vantagem para encerrar Imobilizado; conta como 1 tamanho maior para capacidade de carga.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Golias 5e (PHB): traços seedados")
}

// ── Humano (PHB 2024) ─────────────────────────────────────────────────────────

func seedHumano5ePHB(db *gorm.DB) {
	id, ok := getRace5e(db, "Humano")
	if !ok {
		return
	}

	skills := []domain.Skill{
		{
			Name: "Eficiente", Edition: "5e", RaceID: &id,
			Description: "Você adquire Inspiração Heroica sempre que completar um Descanso Longo.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Ganha Inspiração Heroica a cada Descanso Longo.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Hábil", Edition: "5e", RaceID: &id,
			Description: "Você adquire proficiência em uma perícia à sua escolha.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Proficiência em 1 perícia à escolha.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "habil_humano_phb",
		},
		{
			Name: "Versátil", Edition: "5e", RaceID: &id,
			Description: "Você adquire um talento de Origem à sua escolha (veja o capítulo 5). Habilidoso é recomendado.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Ganhe um talento de Origem adicional à escolha.",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "versatil_humano_phb",
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Humano 5e (PHB): traços seedados")
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
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 36m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Pico de Adrenalina", Edition: "5e", RaceID: &id,
			Description: "Você pode executar a ação Correr como uma Ação Bônus. Ao fazê-lo, você ganha PV Temporários iguais ao seu Bônus de Proficiência. Usos = Bônus de Proficiência; recupera em Descanso Curto ou Longo.",
			ActionType:  "Ação Bônus", Range: "Pessoal",
			Effect:    "Ação Correr + PV Temporários = Bônus de Proficiência. Recupera em Descanso Curto.",
			PowerType: domain.PowerEncounter, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Vigor Implacável", Edition: "5e", RaceID: &id,
			Description: "Quando você é reduzido a 0 Pontos de Vida, mas não morre imediatamente, você pode usar sua Reação para ficar com 1 Ponto de Vida. Após usar este traço, você não pode fazê-lo novamente até completar um Descanso Longo.",
			ActionType:  "Reação", Range: "Pessoal",
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
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Vantagem em salvaguardas contra Amedrontado.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Agilidade Pequenina", Edition: "5e", RaceID: &id,
			Description: "Você pode se mover pelo espaço de qualquer criatura que seja um tamanho maior que você, mas não pode parar no mesmo espaço.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Pode se mover pelo espaço de criaturas Médias ou maiores.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Sorte", Edition: "5e", RaceID: &id,
			Description: "Ao tirar 1 no dado D20 de um Teste de D20, você pode jogar novamente o dado e deve usar a nova jogada.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Rolar 1 no d20 permite rejogada obrigatória.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Furtividade Natural", Edition: "5e", RaceID: &id,
			Description: "Você pode executar a ação Esconder mesmo quando estiver encoberto apenas por uma criatura que seja pelo menos um tamanho maior que você.",
			ActionType:  "Passiva", Range: "Pessoal",
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
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Enxerga em escuridão até 18m como se fosse meia-luz.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
		{
			Name: "Legado Ínfero", Edition: "5e", RaceID: &id,
			Description: "Escolha um legado que determina sua linhagem ínfera: Infernal (diabólico), Abissal (demoníaco) ou Ctônico (sombrio). Cada legado concede um truque diferente no nível 1, e magias progressivas nos níveis 3 e 5, conjuradas uma vez por dia sem espaço de magia.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:         "Magias inatas por linhagem ínfera. Atributo de conjuração: INT, SAB ou CAR (escolha ao criar).",
			PowerType:      domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
			RequiresChoice: true, ChoiceGroup: "legado_tiferino",
		},
		{
			Name: "Presença Sobrenatural", Edition: "5e", RaceID: &id,
			Description: "Você conhece o truque Taumaturgia. Ao conjurar com este traço, use o mesmo atributo de conjuração do seu Legado Ínfero.",
			ActionType:  "Passiva", Range: "Pessoal",
			Effect:    "Truque Taumaturgia disponível.",
			PowerType: domain.PowerUnlimited, Level: 1, IsRaceFeature: true,
		},
	}
	for _, s := range skills {
		upsertRaceSkill(db, s, id)
	}
	log.Println("  ✓ Tiferino 5e: traços seedados")
}