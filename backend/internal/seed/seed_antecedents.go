package seed

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// seedAntecedents5e insere os 16 Antecedentes oficiais do PHB 2024 (capítulo 4,
// "Origens dos Personagens", extraídos diretamente do PDF via PyMuPDF — não
// inventados) mais "Herói do Povo" (antecedente de livro antigo, 2014, mantido
// via a regra de "Antecedentes e Espécies de Livros Antigos", cap. 2 p.38).
// Usa domain.Antecedent (tabela: antecedents) — separado da biografia do
// personagem (Background).
//
// Nenhum dos 16 antecedentes de 2024 tem uma "característica" separada além
// do próprio Talento de Origem que concede (diferente do modelo de 2014, que
// tinha uma feature narrativa distinta) — por isso Feature/FeatureDescription
// ficam vazios para eles; o frontend só exibe esse bloco quando preenchido.
// "Herói do Povo" é a exceção: mantém sua feature 2014 original ("Hospitalidade
// Rústica") porque é um antecedente legado, não substituído por um equivalente
// de 2024.
func seedAntecedents5e(db *gorm.DB) {
	antecedents := []domain.Antecedent{
		{
			Name:                "Acólito",
			Edition:             "5e",
			Description:         "Você passou sua vida a serviço de um templo dedicado a um deus ou panteão de deuses. Você age como intermediário entre o reino dos sagrados e o mundo mortal.",
			SkillProficiencies:  `["Intuição","Religião"]`,
			ToolProficiencies:   "Suprimentos de Calígrafo",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Suprimentos de Calígrafo, Livro (orações), Símbolo Sagrado, Pergaminho (10 folhas), Túnica, 8 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["INT","SAB","CAR"]`,
			OriginFeatName:      "Iniciado em Magia",
		},
		{
			Name:                "Andarilho",
			Edition:             "5e",
			Description:         "Você cresceu nas ruas, cercado por outros igualmente desafortunados. Sobreviveu fazendo bicos e, quando a fome apertava, recorria ao furto — mas nunca perdeu o orgulho nem abandonou a esperança.",
			SkillProficiencies:  `["Furtividade","Intuição"]`,
			ToolProficiencies:   "Ferramentas de Ladrão",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) 2 Adagas, Ferramentas de Ladrão, Kit de Jogos (qualquer um), 2 Algibeiras, Roupas de Viagem, Saco de Dormir, 16 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["DES","SAB","CAR"]`,
			OriginFeatName:      "Sortudo",
		},
		{
			Name:                "Artesão",
			Edition:             "5e",
			Description:         "Você começou limpando pisos e balcões na oficina de um artesão até ficar forte o bastante para ser útil. Como aprendiz, aprendeu a fabricar artesanato básico, lidar com clientes exigentes e desenvolveu um olhar afiado para detalhes.",
			SkillProficiencies:  `["Investigação","Persuasão"]`,
			ToolProficiencies:   "Escolha um tipo de Ferramentas de Artesão",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Ferramentas de Artesão (a mesma escolhida), 2 Algibeiras, Roupas de Viagem, 32 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","DES","INT"]`,
			OriginFeatName:      "Artifista",
		},
		{
			Name:                "Artista",
			Edition:             "5e",
			Description:         "Você passou a juventude em feiras e festivais itinerantes, fazendo bicos para músicos e acrobatas em troca de aulas. Aprendeu a andar na corda bamba, tocar instrumentos e recitar poesia, e até hoje prospera com aplausos.",
			SkillProficiencies:  `["Acrobacia","Atuação"]`,
			ToolProficiencies:   "Escolha um tipo de Instrumento Musical",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Instrumento Musical (o mesmo escolhido), Espelho, 2 Fantasias, Perfume, Roupas de Viagem, 11 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","DES","CAR"]`,
			OriginFeatName:      "Músico",
		},
		{
			Name:                "Charlatão",
			Edition:             "5e",
			Description:         "Você percorreu o circuito de bares e botequins desde jovem, aprendendo a lidar com pessoas em busca de mentiras reconfortantes — talvez uma poção falsa ou registros de ancestralidade forjados.",
			SkillProficiencies:  `["Enganação","Prestidigitação"]`,
			ToolProficiencies:   "Kit de Falsificação",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Kit de Falsificação, Fantasia, Roupas Finas, 15 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["DES","CON","CAR"]`,
			OriginFeatName:      "Habilidoso",
		},
		{
			Name:                "Criminoso",
			Edition:             "5e",
			Description:         "Você tem um histórico de crimes e mantém contatos no submundo. É experiente em furtividade e prestidigitação, e sobreviveu à margem da lei nos becos escuros da cidade.",
			SkillProficiencies:  `["Furtividade","Prestidigitação"]`,
			ToolProficiencies:   "Ferramentas de Ladrão",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) 2 Adagas, Ferramentas de Ladrão, 2 Algibeiras, Pé de Cabra, Roupas de Viagem, 16 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["DES","CON","INT"]`,
			OriginFeatName:      "Alerta",
		},
		{
			Name:                "Eremita",
			Edition:             "5e",
			Description:         "Você passou seus primeiros anos isolado em uma cabana ou mosteiro, muito além do povoado mais próximo, tendo como companhia apenas as criaturas da floresta. A solidão te permitiu ponderar longamente os mistérios da criação.",
			SkillProficiencies:  `["Medicina","Religião"]`,
			ToolProficiencies:   "Kit de Herbalismo",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Cajado, Kit de Herbalismo, Lâmpada, Livro (filosofia), Óleo (3 frascos), Roupas de Viagem, Saco de Dormir, 16 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["CON","SAB","CAR"]`,
			OriginFeatName:      "Curandeiro",
		},
		{
			Name:                "Escriba",
			Edition:             "5e",
			Description:         "Você passou anos de formação em um scriptorium, mosteiro ou agência governamental, aprendendo a escrever com mão firme e produzir textos cuidadosamente redigidos, com atenção meticulosa aos detalhes.",
			SkillProficiencies:  `["Investigação","Percepção"]`,
			ToolProficiencies:   "Suprimentos de Calígrafo",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Suprimentos de Calígrafo, Lâmpada, Óleo (3 frascos), Pergaminho (12 folhas), Roupas Finas, 23 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["DES","INT","SAB"]`,
			OriginFeatName:      "Habilidoso",
		},
		{
			Name:                "Fazendeiro",
			Edition:             "5e",
			Description:         "Você cresceu perto da terra. Anos cuidando de animais e cultivando o solo te recompensaram com paciência e boa saúde, além de um grande apreço pela generosidade — e pela ira — da natureza.",
			SkillProficiencies:  `["Lidar com Animais","Natureza"]`,
			ToolProficiencies:   "Ferramentas de Carpinteiro",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Foice, Ferramentas de Carpinteiro, Kit de Curandeiro, Balde de Ferro, Pá, 30 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","CON","SAB"]`,
			OriginFeatName:      "Vigoroso",
		},
		{
			Name:                "Guarda",
			Edition:             "5e",
			Description:         "Você foi treinado para manter vigília constante em seu posto na torre, um olho voltado para saqueadores do lado de fora da muralha e outro para assaltantes e encrenqueiros do lado de dentro.",
			SkillProficiencies:  `["Atletismo","Percepção"]`,
			ToolProficiencies:   "Escolha um tipo de Kit de Jogos",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Lança, Besta Leve, 20 Virotes, Kit de Jogo (o mesmo escolhido), Aljava, Grilhões, Lanterna Coberta, Roupas de Viagem, 12 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","INT","SAB"]`,
			OriginFeatName:      "Alerta",
		},
		{
			Name:                "Guia",
			Edition:             "5e",
			Description:         "Você cresceu ao ar livre, longe de terras povoadas, aprendendo a se defender entre monstros estranhos e ruínas esquecidas. De tempos em tempos, guiava sacerdotes da natureza que o instruíam nos fundamentos da magia natural.",
			SkillProficiencies:  `["Furtividade","Sobrevivência"]`,
			ToolProficiencies:   "Ferramentas de Cartógrafo",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Arco Curto, 20 Flechas, Ferramentas de Cartógrafo, Aljava, Roupas de Viagem, Saco de Dormir, Tenda, 3 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["DES","CON","SAB"]`,
			OriginFeatName:      "Iniciado em Magia",
		},
		{
			Name:                "Marinheiro",
			Edition:             "5e",
			Description:         "Você viveu com o vento nas costas e os conveses balançando sob os pés, sentando-se em bares de mais portos do que consegue lembrar, enfrentando tempestades e trocando histórias com quem vive sob as ondas.",
			SkillProficiencies:  `["Acrobacia","Percepção"]`,
			ToolProficiencies:   "Ferramentas de Navegador",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Adaga, Ferramentas de Navegador, Corda, Roupas de Viagem, 20 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","DES","SAB"]`,
			OriginFeatName:      "Valentão de Taverna",
		},
		{
			Name:                "Mercador",
			Edition:             "5e",
			Description:         "Você foi aprendiz de um comerciante, mestre de caravanas ou lojista, aprendendo os fundamentos do comércio — viajando, comprando e vendendo matérias-primas e mercadorias entre artesãos e clientes.",
			SkillProficiencies:  `["Lidar com Animais","Persuasão"]`,
			ToolProficiencies:   "Ferramentas de Navegador",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Ferramentas de Navegador, 2 Algibeiras, Roupas de Viagem, 22 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["CON","INT","CAR"]`,
			OriginFeatName:      "Sortudo",
		},
		{
			Name:                "Nobre",
			Edition:             "5e",
			Description:         "Você foi criado em um castelo, cercado por riqueza, poder e privilégio, recebendo uma educação de primeira categoria e aprendendo muito sobre liderança ao observar sua família na corte.",
			SkillProficiencies:  `["História","Persuasão"]`,
			ToolProficiencies:   "Escolha um tipo de Kit de Jogos",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Kit de Jogos (o mesmo escolhido), Perfume, Roupas Finas, 29 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","INT","CAR"]`,
			OriginFeatName:      "Habilidoso",
		},
		{
			Name:                "Sábio",
			Edition:             "5e",
			Description:         "Você passou anos de formação viajando entre mansões e mosteiros, realizando serviços em troca de acesso a bibliotecas, estudando livros e pergaminhos até os rudimentos da magia — e sua mente ainda anseia por mais.",
			SkillProficiencies:  `["Arcanismo","História"]`,
			ToolProficiencies:   "Suprimentos de Calígrafo",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Cajado, Suprimentos de Calígrafo, Livro (história), Pergaminho (8 folhas), Túnica, 8 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["CON","INT","SAB"]`,
			OriginFeatName:      "Iniciado em Magia",
		},
		{
			Name:                "Soldado",
			Edition:             "5e",
			Description:         "A guerra foi sua vida desde jovem. Você treinou com armas e armaduras, aprendendo as artes militares assim que atingiu a idade adulta. A batalha está no seu sangue.",
			SkillProficiencies:  `["Atletismo","Intimidação"]`,
			ToolProficiencies:   "Escolha um tipo de Kit de Jogos",
			Languages:           "",
			Equipment:           "Escolha A ou B: (A) Lança, Arco Curto, 20 Flechas, Kit de Curandeiro, Kit de Jogo (o mesmo escolhido), Aljava, Roupas de Viagem, 14 PO; ou (B) 50 PO",
			Feature:             "",
			FeatureDescription:  "",
			IsDefault:           true,
			AbilityBonusOptions: `["FOR","DES","CON"]`,
			OriginFeatName:      "Atacante Selvagem",
		},
		{
			Name:               "Herói do Povo",
			Edition:            "5e",
			Description:        "Você veio das classes trabalhadoras e está destinado a algo maior. As pessoas comuns te respeitam como um de seus próprios.",
			SkillProficiencies: `["Lidar com Animais","Sobrevivência"]`,
			ToolProficiencies:  "Um tipo de ferramenta de artesão, veículos terrestres",
			Languages:          "",
			Equipment:          "Ferramentas de artesão, uma pá, um pote de ferro, roupas comuns e 10 po",
			Feature:            "Hospitalidade Rústica",
			FeatureDescription: "As pessoas comuns te recebem de bom grado. Consegues comida, abrigo e proteção das pessoas do povo, que te escondem de perseguições se necessário.",
			IsDefault:          true,
			// "Herói do Povo" (Folk Hero) é um antecedente só de 2014 — não foi
			// reimpresso no PHB 2024 com esse nome. Em vez de remapear pra um
			// antecedente 2024 parecido, a regra oficial (caixa "Antecedentes e
			// Espécies de Livros Antigos", PHB 2024 cap. 2, p.38) é manter o
			// antecedente antigo e liberar as duas escolhas: bônus de atributo
			// entre os 6 atributos (em vez de 3 fixos) e talento de Origem à
			// escolha (em vez de um fixo) — daí IsLegacy em vez de
			// AbilityBonusOptions/OriginFeatName.
			IsLegacy: true,
		},
	}

	for _, a := range antecedents {
		var existing domain.Antecedent
		if db.Where("name = ? AND edition = ?", a.Name, "5e").First(&existing).Error != nil {
			db.Create(&a)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description":           a.Description,
				"skill_proficiencies":   a.SkillProficiencies,
				"tool_proficiencies":    a.ToolProficiencies,
				"languages":             a.Languages,
				"equipment":             a.Equipment,
				"feature":               a.Feature,
				"feature_description":   a.FeatureDescription,
				"ability_bonus_options": a.AbilityBonusOptions,
				"origin_feat_name":      a.OriginFeatName,
				"is_legacy":             a.IsLegacy,
			})
		}
	}
	log.Println("  ✓ Antecedentes 5e: 17 seedados (16 do PHB 2024 + Herói do Povo, legado 2014)")
}
