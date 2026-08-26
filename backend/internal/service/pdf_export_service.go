package service

import (
	"encoding/json"
	"fmt"

	"rpg-manager/internal/domain"
)

// abilityScoreFor converte a abreviação em pt-BR usada em domain.Pericia.Attribute
// ("FOR","DES","CON","INT","SAB","CAR") para o valor do atributo correspondente do personagem.
func abilityScoreFor(character domain.Character, attr string) int {
	switch attr {
	case "FOR":
		return character.Strength
	case "DES":
		return character.Dexterity
	case "CON":
		return character.Constitution
	case "INT":
		return character.Intelligence
	case "SAB":
		return character.Wisdom
	case "CAR":
		return character.Charisma
	default:
		return 10
	}
}

// BuildPDF5eExportPayload monta o payload já totalmente calculado que o serviço Python usa
// apenas para preencher os campos do AcroForm (nenhuma regra de D&D é reimplementada em Python).
//
// allPericias5e: catálogo completo das 18 perícias de 5e (PericiaService.GetAll("5e")).
// characterPericias: perícias em que o personagem é proficiente (PericiaService.GetByCharacter).
func BuildPDF5eExportPayload(
	character domain.Character,
	allPericias5e []domain.Pericia,
	characterPericias []domain.CharacterPericia,
	armorService *ArmorService,
) map[string]interface{} {
	profBonus := character.ProficiencyBonus
	if profBonus == 0 {
		profBonus = proficiencyBonus5e(character.Level)
	}

	proficientPericias := make(map[string]bool, len(characterPericias))
	for _, cp := range characterPericias {
		proficientPericias[cp.PericiaName] = true
	}

	var savingThrowAbilities []string
	_ = json.Unmarshal([]byte(character.Class.SavingThrows), &savingThrowAbilities)
	proficientSaves := make(map[string]bool, len(savingThrowAbilities))
	for _, ab := range savingThrowAbilities {
		proficientSaves[ab] = true
	}

	abilities := map[string]int{
		"FOR": character.Strength,
		"DES": character.Dexterity,
		"CON": character.Constitution,
		"INT": character.Intelligence,
		"SAB": character.Wisdom,
		"CAR": character.Charisma,
	}

	atributos := map[string]interface{}{}
	for attr, score := range abilities {
		atributos[attr] = map[string]interface{}{"valor": score, "mod": mod(score)}
	}

	salvaguardas := map[string]interface{}{}
	for attr, score := range abilities {
		proficiente := proficientSaves[attr]
		valor := mod(score)
		if proficiente {
			valor += profBonus
		}
		salvaguardas[attr] = map[string]interface{}{"valor": valor, "proficiente": proficiente}
	}

	pericias := map[string]interface{}{}
	for _, p := range allPericias5e {
		proficiente := proficientPericias[p.Name]
		valor := mod(abilityScoreFor(character, p.Attribute))
		if proficiente {
			valor += profBonus
		}
		pericias[p.Name] = map[string]interface{}{"valor": valor, "proficiente": proficiente}
	}

	percepcaoPassiva := 10
	if p, ok := pericias["Percepção"].(map[string]interface{}); ok {
		percepcaoPassiva = 10 + p["valor"].(int)
	}

	ca := armorService.CalculateAC(character)

	className := character.Class.Name
	classeNivel := fmt.Sprintf("%s %d", className, character.Level)

	antecedenteNome := ""
	if character.Antecedent != nil {
		antecedenteNome = character.Antecedent.Name
	}

	raca := ""
	if character.Race.Name != "" {
		raca = character.Race.Name
	}

	hitDie := character.Class.HitDie
	if hitDie == 0 {
		hitDie = 8
	}
	dadosDeVidaTotal := fmt.Sprintf("%dd%d", character.Level, hitDie)

	sucessos := character.DeathSaveSuccesses
	falhas := character.DeathSaveFailures

	return map[string]interface{}{
		"nome":                 character.Name,
		"classe_nivel":         classeNivel,
		"antecedente":          antecedenteNome,
		"raca":                 raca,
		"alinhamento":          character.Alignment,
		"xp":                   character.ExperiencePoints,
		"atributos":            atributos,
		"bonus_proficiencia":   profBonus,
		"ca":                   ca,
		"iniciativa":           mod(character.Dexterity),
		"deslocamento":         character.Speed,
		"pv_maximo":            character.MaxHP,
		"pv_atual":             character.HitPoints,
		"pv_temporario":        character.TempHP,
		"dados_de_vida_total":  dadosDeVidaTotal,
		"dados_de_vida":        "",
		"resistencia_morte":    map[string]interface{}{"sucessos": sucessos, "falhas": falhas},
		"salvaguardas":         salvaguardas,
		"pericias":             pericias,
		"percepcao_passiva":    percepcaoPassiva,
		"tracos_personalidade": character.PersonalityTraits,
		"ideais":               character.Ideals,
		"vinculos":             character.Bonds,
		"defeitos":             character.Flaws,
	}
}
