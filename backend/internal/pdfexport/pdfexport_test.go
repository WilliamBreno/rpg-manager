package pdfexport

import (
	"bytes"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// TestFillCharacterSheet5e_CheckboxMapping é o teste de regressão do bug de
// "embaralhamento" documentado no CLAUDE.md: os widgets de perícia da ficha
// (texto + checkbox de proficiência) ficaram na ordem alfabética em inglês
// do template original, mas os rótulos impressos seguem a ordem alfabética
// em português — então o campo literalmente chamado "History " renderiza na
// linha impressa "Furtividade", não "História". Se o mapa de campos
// (dnd5e_pdf_field_map.json) for alterado sem reverificar isso, este teste
// falha em vez de silenciosamente exportar o número certo pra perícia
// errada.
func TestFillCharacterSheet5e_CheckboxMapping(t *testing.T) {
	payload := map[string]interface{}{
		"nome":                "Teste Regressão",
		"classe_nivel":        "Guerreiro 1",
		"antecedente":         "Soldado",
		"raca":                "Humano",
		"alinhamento":         "Leal e Bom",
		"xp":                  0,
		"bonus_proficiencia":  2,
		"ca":                  12,
		"iniciativa":          2,
		"deslocamento":        30,
		"pv_maximo":           12,
		"pv_atual":            12,
		"pv_temporario":       0,
		"dados_de_vida_total": "1d10",
		"dados_de_vida":       "",
		"percepcao_passiva":   12,
		"atributos": map[string]interface{}{
			"FOR": map[string]interface{}{"valor": 16, "mod": 3},
			"DES": map[string]interface{}{"valor": 14, "mod": 2},
			"CON": map[string]interface{}{"valor": 14, "mod": 2},
			"INT": map[string]interface{}{"valor": 10, "mod": 0},
			"SAB": map[string]interface{}{"valor": 10, "mod": 0},
			"CAR": map[string]interface{}{"valor": 10, "mod": 0},
		},
		"resistencia_morte": map[string]interface{}{"sucessos": 0, "falhas": 0},
		"salvaguardas": map[string]interface{}{
			"FOR": map[string]interface{}{"valor": 5, "proficiente": true},
			"DES": map[string]interface{}{"valor": 2, "proficiente": false},
			"CON": map[string]interface{}{"valor": 4, "proficiente": true},
			"INT": map[string]interface{}{"valor": 0, "proficiente": false},
			"SAB": map[string]interface{}{"valor": 0, "proficiente": false},
			"CAR": map[string]interface{}{"valor": 0, "proficiente": false},
		},
		// Só estas 3 treinadas — mesma combinação já verificada manualmente
		// (ver histórico do projeto): Furtividade, Percepção e Atletismo.
		"pericias": map[string]interface{}{
			"Furtividade": map[string]interface{}{"valor": 4, "proficiente": true},
			"Percepção":   map[string]interface{}{"valor": 2, "proficiente": true},
			"Atletismo":   map[string]interface{}{"valor": 5, "proficiente": true},
			"Acrobacia":   map[string]interface{}{"valor": 2, "proficiente": false},
			"Persuasão":   map[string]interface{}{"valor": 0, "proficiente": false},
		},
	}

	pdfBytes, err := FillCharacterSheet5e(payload)
	if err != nil {
		t.Fatalf("FillCharacterSheet5e retornou erro: %v", err)
	}

	group, err := api.ExportForm(bytes.NewReader(pdfBytes), "", nil)
	if err != nil {
		t.Fatalf("falha ao reexportar o formulário do PDF gerado: %v", err)
	}
	if len(group.Forms) == 0 {
		t.Fatal("PDF gerado não tem nenhum formulário")
	}

	textByName := map[string]string{}
	for _, tf := range group.Forms[0].TextFields {
		textByName[tf.Name] = tf.Value
	}
	checkByName := map[string]bool{}
	for _, cb := range group.Forms[0].CheckBoxes {
		checkByName[cb.Name] = cb.Value
	}

	// Campo "History " deve conter o valor de Furtividade (+4) e "Check Box
	// 28" deve estar marcado — não "História" nem qualquer outra perícia.
	cases := []struct {
		campoValor string
		valorEsp   string
		campoCheck string
		marcadoEsp bool
		descricao  string
	}{
		{"History ", "+4", "Check Box 28", true, "Furtividade"},
		{"Persuasion", "+2", "Check Box 36", true, "Percepção"},
		{"Arcana", "+5", "Check Box 25", true, "Atletismo"},
		{"Acrobatics", "+2", "Check Box 23", false, "Acrobacia (não treinada)"},
		{"Religion", "+0", "Check Box 37", false, "Persuasão (não treinada)"},
	}

	for _, tc := range cases {
		if got := textByName[tc.campoValor]; got != tc.valorEsp {
			t.Errorf("%s: campo %q = %q, esperado %q", tc.descricao, tc.campoValor, got, tc.valorEsp)
		}
		if got := checkByName[tc.campoCheck]; got != tc.marcadoEsp {
			t.Errorf("%s: checkbox %q marcado=%v, esperado %v", tc.descricao, tc.campoCheck, got, tc.marcadoEsp)
		}
	}

	if textByName["CharacterName"] != "Teste Regressão" {
		t.Errorf("CharacterName = %q, esperado %q", textByName["CharacterName"], "Teste Regressão")
	}
	if textByName["ClassLevel"] != "Guerreiro 1" {
		t.Errorf("ClassLevel = %q, esperado %q", textByName["ClassLevel"], "Guerreiro 1")
	}
}
