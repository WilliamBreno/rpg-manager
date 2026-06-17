package seed

import (
	"log"
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

func seedTalentos(db *gorm.DB) {
	talentos := []domain.Talento{
		// ── DEFESA ───────────────────────────────────────────────────────────────
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
		// ── COMBATE ──────────────────────────────────────────────────────────────
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
		// ── PERÍCIAS ─────────────────────────────────────────────────────────────
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
		// ── MAGIA ─────────────────────────────────────────────────────────────────
		{
			Name: "Lançador de Rituais", Edition: "4e", Category: "Magia",
			Prerequisite: "Inteligência 13 ou Sabedoria 13",
			Description:  "Pode aprender e realizar rituais mágicos.",
			Tooltip:      "Permite aprender e realizar rituais com componentes adequados (ouro). Rituais possibilitam efeitos poderosos fora do combate: teletransporte, comunicação à distância, detecção mágica. Requer livro de rituais.",
		},
		// ── ARMADURA ─────────────────────────────────────────────────────────────
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
	log.Println("  ✓ Talentos seedados")
}