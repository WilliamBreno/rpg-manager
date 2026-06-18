package seed

import (
	"log"
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// seedAntecedents5e insere os Antecedentes das Regras Básicas do D&D 5e.
// Usa domain.Antecedent (tabela: antecedents) — separado da biografia do personagem (Background).
func seedAntecedents5e(db *gorm.DB) {
	antecedents := []domain.Antecedent{
		{
			Name:               "Acólito",
			Edition:            "5e",
			Description:        "Você passou sua vida a serviço de um templo dedicado a um deus ou panteão de deuses. Você age como intermediário entre o reino dos sagrados e o mundo mortal.",
			SkillProficiencies: `["Intuição","Religião"]`,
			ToolProficiencies:  "",
			Languages:          "Dois idiomas à sua escolha",
			Equipment:          "Símbolo sagrado, livro de orações, 5 velas, robes, roupas comuns e 15 po",
			Feature:            "Abrigo dos Fiéis",
			FeatureDescription: "Como acólito, você comanda o respeito de outros que compartilham sua fé. Pode realizar cerimônias religiosas e receber abrigo, refeições e cuidados em qualquer templo da sua divindade.",
			IsDefault:          true,
		},
		{
			Name:               "Criminoso",
			Edition:            "5e",
			Description:        "Você tem um histórico de crimes e mantém contatos no submundo. É experiente em furtividade e enganação, e sobreviveu à margem da lei.",
			SkillProficiencies: `["Furtividade","Enganação"]`,
			ToolProficiencies:  "Ferramentas de ladrão, um tipo de jogo",
			Languages:          "",
			Equipment:          "Pé de cabra, roupas escuras com capuz e 15 po",
			Feature:            "Contato no Crime",
			FeatureDescription: "Você tem um contato confiável no submundo criminal. Esse contato pode fornecer informações sobre pessoas, locais e crimes locais, além de intermediar serviços ilícitos.",
			IsDefault:          true,
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
		},
		{
			Name:               "Sábio",
			Edition:            "5e",
			Description:        "Você passou anos absorvendo tomos e pergaminhos em uma biblioteca ou academia. Sua compreensão do mundo é teórica mas profunda.",
			SkillProficiencies: `["Arcanismo","História"]`,
			ToolProficiencies:  "",
			Languages:          "Dois idiomas à sua escolha",
			Equipment:          "Garrafa de tinta, pena, faca pequena, carta de um colega com perguntas sem resposta, roupas comuns e 10 po",
			Feature:            "Pesquisador",
			FeatureDescription: "Quando tenta aprender ou lembrar de uma informação, mesmo que não saiba a resposta, frequentemente sabe onde e de quem pode obtê-la. O Mestre pode decidir que a informação é obscura demais para ser encontrada.",
			IsDefault:          true,
		},
		{
			Name:               "Soldado",
			Edition:            "5e",
			Description:        "A guerra foi sua vida desde jovem. Você treinou com armas e armaduras, aprendendo as artes militares. Seu passado de combate moldou quem você é.",
			SkillProficiencies: `["Atletismo","Intimidação"]`,
			ToolProficiencies:  "Um tipo de jogo, veículos terrestres",
			Languages:          "",
			Equipment:          "Insígnia de patente, troféu de inimigo abatido, roupas comuns e 10 po",
			Feature:            "Hierarquia Militar",
			FeatureDescription: "Você tem uma patente reconhecida de sua carreira militar. Soldados leais ao seu antigo exército te reconhecem e prestam deferência, e você tem acesso a acampamentos e suprimentos militares.",
			IsDefault:          true,
		},
	}

	for _, a := range antecedents {
		var existing domain.Antecedent
		if db.Where("name = ? AND edition = ?", a.Name, "5e").First(&existing).Error != nil {
			db.Create(&a)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description":         a.Description,
				"skill_proficiencies": a.SkillProficiencies,
				"tool_proficiencies":  a.ToolProficiencies,
				"languages":           a.Languages,
				"equipment":           a.Equipment,
				"feature":             a.Feature,
				"feature_description": a.FeatureDescription,
			})
		}
	}
	log.Println("  ✓ Antecedentes 5e: 5 seedados (Acólito, Criminoso, Herói do Povo, Sábio, Soldado)")
}