package service

import (
	"encoding/json"
	"fmt"
	"strings"

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

// buildCaracteristicasHabilidades monta o texto livre da caixa "Features and
// Traits" (página 1) a partir dos dados JÁ estruturados do personagem —
// características de classe/raça (domain.Skill com IsClassFeature/
// IsRaceFeature) e talentos (domain.Talento) — sem inventar nada que não
// esteja registrado no personagem.
func buildCaracteristicasHabilidades(character domain.Character) string {
	var classFeatures, raceFeatures []string
	for _, s := range character.Skills {
		switch {
		case s.IsClassFeature:
			classFeatures = append(classFeatures, s.Name)
		case s.IsRaceFeature:
			raceFeatures = append(raceFeatures, s.Name)
		}
	}

	var talentos []string
	for _, t := range character.Talentos {
		talentos = append(talentos, t.Name)
	}

	var b strings.Builder
	if len(classFeatures) > 0 {
		b.WriteString("CARACTERÍSTICAS DE CLASSE: ")
		b.WriteString(strings.Join(classFeatures, "; "))
		b.WriteString("\n\n")
	}
	if len(raceFeatures) > 0 {
		b.WriteString("TRAÇOS RACIAIS: ")
		b.WriteString(strings.Join(raceFeatures, "; "))
		b.WriteString("\n\n")
	}
	if len(talentos) > 0 {
		b.WriteString("TALENTOS: ")
		b.WriteString(strings.Join(talentos, "; "))
	}

	return strings.TrimSpace(b.String())
}

var abilityLabelPT = map[string]string{
	"FOR": "Força", "DES": "Destreza", "CON": "Constituição",
	"INT": "Inteligência", "SAB": "Sabedoria", "CAR": "Carisma",
}

// buildConjuracao monta os dados da página 3 (conjuração): habilidade-chave,
// CD de resistência, bônus de ataque de magia, espaços de magia por círculo
// (0 a 9) e as magias conhecidas em cada círculo. Retorna nil se a classe do
// personagem não conjura (a página inteira fica em branco nesse caso, ver
// pdfexport.FillCharacterSheet5e).
func buildConjuracao(character domain.Character) map[string]interface{} {
	className := character.Class.Name
	abilityAbbr, casts := spellcastingAbility5e[className]

	// Terço-conjurador (Cavaleiro Místico/Trapaceiro Arcano) — só conjura se
	// tiver escolhido essa subclasse especificamente; Guerreiro/Ladino comuns
	// não estão no mapa acima e por isso não conjuram por padrão.
	subclassName := ""
	for _, s := range character.Skills {
		if s.Name == "Cavaleiro Místico" || s.Name == "Trapaceiro Arcano" {
			subclassName = s.Name
		}
	}
	if className == "Guerreiro" || className == "Ladino" {
		casts = subclassName != ""
		abilityAbbr = "INT"
	}
	if !casts {
		return nil
	}

	profBonus := character.ProficiencyBonus
	if profBonus == 0 {
		profBonus = proficiencyBonus5e(character.Level)
	}
	habilidadeMod := mod(abilityScoreFor(character, abilityAbbr))

	var slots []int
	if className == "Bruxo" {
		s, circle := pactMagic5e(character.Level)
		if s > 0 && circle > 0 {
			slots = make([]int, circle)
			slots[circle-1] = s
		}
	} else {
		slots = spellSlots5e(className, subclassName, character.Level)
	}

	spellsByCircle := map[int][]string{}
	for _, sp := range character.Spells {
		spellsByCircle[sp.Level] = append(spellsByCircle[sp.Level], sp.Name)
	}

	circulos := map[string]interface{}{}
	for c := 0; c <= 9; c++ {
		total := 0
		if c >= 1 && c <= len(slots) {
			total = slots[c-1]
		}
		circulos[fmt.Sprintf("%d", c)] = map[string]interface{}{
			"total":  total,
			"magias": spellsByCircle[c],
		}
	}

	return map[string]interface{}{
		"classe":     className,
		"habilidade": abilityLabelPT[abilityAbbr],
		"cd":         8 + profBonus + habilidadeMod,
		"ataque":     profBonus + habilidadeMod,
		"circulos":   circulos,
	}
}

// BuildPDF5eExportPayload monta o payload já totalmente calculado que
// internal/pdfexport usa apenas para preencher os campos do AcroForm
// (nenhuma regra de D&D é reimplementada ali).
//
// allPericias5e: catálogo completo das 18 perícias de 5e (PericiaService.GetAll("5e")).
// characterPericias: perícias em que o personagem é proficiente (PericiaService.GetByCharacter).
// playerName: nome do dono da conta (User.Name) — não confundir com o nome do personagem,
// vai no campo "Nome do Jogador" da ficha.
func BuildPDF5eExportPayload(
	character domain.Character,
	allPericias5e []domain.Pericia,
	characterPericias []domain.CharacterPericia,
	armorService *ArmorService,
	playerName string,
) map[string]interface{} {
	profBonus := character.ProficiencyBonus
	if profBonus == 0 {
		profBonus = proficiencyBonus5e(character.Level)
	}

	proficientPericias := make(map[string]bool, len(characterPericias))
	expertisePericias := make(map[string]bool, len(characterPericias))
	for _, cp := range characterPericias {
		proficientPericias[cp.PericiaName] = true
		if cp.Expertise {
			expertisePericias[cp.PericiaName] = true
		}
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
		especializado := expertisePericias[p.Name]
		valor := mod(abilityScoreFor(character, p.Attribute))
		switch {
		case especializado:
			valor += profBonus * 2
		case proficiente:
			valor += profBonus
		}
		pericias[p.Name] = map[string]interface{}{"valor": valor, "proficiente": proficiente, "especializado": especializado}
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
		"nome_jogador":         playerName,
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
		"caracteristicas_habilidades": buildCaracteristicasHabilidades(character),
		// ── Página 2: aparência física + história ────────────────────────
		"idade":      character.Age,
		"altura":     character.Height,
		"peso":       character.Weight,
		"olhos":      character.Eyes,
		"pele":       character.Skin,
		"cabelos":    character.Hair,
		"historia":   character.History,
		"avatar_url": character.AvatarURL,
		// ── Moedas ─────────────────────────────────────────────────────────
		"cp": character.CopperPieces,
		"sp": character.SilverPieces,
		"ep": character.ElectrumPieces,
		"gp": character.GoldPieces,
		"pp": character.PlatinumPieces,
		// ── Página 3: conjuração (nil se a classe não conjurar) ───────────
		"conjuracao": buildConjuracao(character),
	}
}
