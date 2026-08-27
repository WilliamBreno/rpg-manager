// Package pdfexport preenche a ficha oficial de personagem 5e (AcroForm) a
// partir do payload já calculado por service.BuildPDF5eExportPayload.
//
// Antes, esse preenchimento era feito por um serviço Python separado
// (ai-service/pdf_export/fill_dnd5e_sheet.py), chamado via HTTP a partir do
// backend Go. Isso funcionava em desenvolvimento local (os dois serviços
// rodando na mesma máquina), mas o ai-service nunca foi hospedado em lugar
// nenhum em produção — então em produção a exportação sempre falhava com
// 503 "Serviço de exportação de PDF não está disponível no momento",
// mesmo com o backend saudável.
//
// Este pacote resolve isso fazendo o preenchimento inteiro dentro do
// backend Go, usando github.com/pdfcpu/pdfcpu (puro Go, sem dependência de
// binário externo). O PDF-modelo e o mapa de campos ficam embutidos no
// próprio binário via go:embed — nenhum arquivo externo é necessário em
// tempo de execução, então isso funciona em qualquer lugar que o backend
// já roda, sem infraestrutura nova.
//
// Verificado manualmente (não assumido): renderizei uma ficha preenchida de
// teste com PyMuPDF e confirmei visualmente que texto E checkboxes
// aparecem corretamente com pdfcpu nesta ficha-modelo específica — há um
// issue conhecido do pdfcpu ("corrupt form field: missing entry AP") para
// checkboxes em ALGUNS PDFs, mas não se manifestou neste arquivo.
package pdfexport

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

//go:embed assets/ficha_5e_modelo.pdf
var templatePDF []byte

//go:embed assets/dnd5e_pdf_field_map.json
var fieldMapRaw []byte

//go:embed assets/form_skeleton.json
var formSkeletonRaw []byte

// ── Estruturas do mapa de campos (mesmo shape de dnd5e_pdf_field_map.json,
// gerado por introspecção do AcroForm + casamento geométrico dos Rects —
// ver ai-service/pdf_export/reference/ — não adivinhado manualmente) ──────

type fieldMap struct {
	Identidade struct {
		NomePersonagem string `json:"nome_personagem"`
		ClasseNivel    string `json:"classe_nivel"`
		Antecedente    string `json:"antecedente"`
		Raca           string `json:"raca"`
		Alinhamento    string `json:"alinhamento"`
		XP             string `json:"xp"`
	} `json:"identidade"`
	Atributos map[string]struct {
		Valor string `json:"valor"`
		Mod   string `json:"mod"`
	} `json:"atributos"`
	Combate struct {
		BonusProficiencia string `json:"bonus_proficiencia"`
		CA                string `json:"ca"`
		Iniciativa        string `json:"iniciativa"`
		Deslocamento      string `json:"deslocamento"`
		PVMaximo          string `json:"pv_maximo"`
		PVAtual           string `json:"pv_atual"`
		PVTemporario      string `json:"pv_temporario"`
		DadosDeVidaTotal  string `json:"dados_de_vida_total"`
		DadosDeVida       string `json:"dados_de_vida"`
	} `json:"combate"`
	ResistenciaContraMorte struct {
		Sucessos []string `json:"sucessos"`
		Falhas   []string `json:"falhas"`
	} `json:"resistencia_contra_morte"`
	Salvaguardas []struct {
		Atributo          string `json:"atributo"`
		CampoValor        string `json:"campo_valor"`
		CampoProficiencia string `json:"campo_proficiencia"`
	} `json:"salvaguardas"`
	Pericias []struct {
		NomePt            string `json:"nome_pt"`
		Atributo          string `json:"atributo"`
		CampoValor        string `json:"campo_valor"`
		CampoProficiencia string `json:"campo_proficiencia"`
	} `json:"pericias"`
	Passiva struct {
		PercepcaoPassiva string `json:"percepcao_passiva"`
	} `json:"passiva"`
	Personalidade struct {
		TracosPersonalidade string `json:"tracos_personalidade"`
		Ideais              string `json:"ideais"`
		Vinculos            string `json:"vinculos"`
		Defeitos            string `json:"defeitos"`
	} `json:"personalidade"`
}

// ── Estruturas do formulário no formato que o pdfcpu espera/produz (ver
// github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form/export.go) — os IDs internos do
// AcroForm (campo "id") são obrigatórios e vêm do form_skeleton.json,
// extraído uma vez via api.ExportFormJSON contra o PDF-modelo real. Não dá
// pra reconstruir esses IDs à mão; por isso o esqueleto é embutido em vez
// de gerado em tempo de execução. ─────────────────────────────────────────

type pdfHeader struct {
	Source   string   `json:"source"`
	Version  string   `json:"version"`
	Creation string   `json:"creation"`
	ID       []string `json:"id,omitempty"`
	Title    string   `json:"title,omitempty"`
	Author   string   `json:"author,omitempty"`
	Creator  string   `json:"creator,omitempty"`
	Producer string   `json:"producer,omitempty"`
	Subject  string   `json:"subject,omitempty"`
	Keywords string   `json:"keywords,omitempty"`
}

type pdfTextField struct {
	Pages     []int  `json:"pages"`
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	AltName   string `json:"altname,omitempty"`
	Default   string `json:"default,omitempty"`
	Value     string `json:"value"`
	MaxLen    int    `json:"maxlen,omitempty"`
	Multiline bool   `json:"multiline"`
	Locked    bool   `json:"locked"`
}

type pdfCheckBox struct {
	Pages   []int  `json:"pages"`
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	AltName string `json:"altname,omitempty"`
	Default bool   `json:"default"`
	Value   bool   `json:"value"`
	Locked  bool   `json:"locked"`
}

type pdfForm struct {
	TextFields []*pdfTextField `json:"textfield,omitempty"`
	CheckBoxes []*pdfCheckBox  `json:"checkbox,omitempty"`
}

type pdfFormGroup struct {
	Header pdfHeader `json:"header"`
	Forms  []pdfForm `json:"forms"`
}

func signed(n int) string {
	if n >= 0 {
		return fmt.Sprintf("+%d", n)
	}
	return fmt.Sprintf("%d", n)
}

// ── Leitura tolerante do payload map[string]interface{} vindo direto de
// service.BuildPDF5eExportPayload (sem passar por JSON, os tipos concretos
// já são int/string/bool/map[string]interface{}) ──────────────────────────

func str(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func intVal(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

func boolVal(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func subMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if sm, ok := v.(map[string]interface{}); ok {
			return sm
		}
	}
	return map[string]interface{}{}
}

// FillCharacterSheet5e preenche a ficha oficial 5e com os dados já calculados
// em `character` (mesmo shape que BuildPDF5eExportPayload retorna) e devolve
// os bytes do PDF preenchido.
func FillCharacterSheet5e(character map[string]interface{}) ([]byte, error) {
	var fm fieldMap
	if err := json.Unmarshal(fieldMapRaw, &fm); err != nil {
		return nil, fmt.Errorf("carregar mapa de campos: %w", err)
	}

	var group pdfFormGroup
	if err := json.Unmarshal(formSkeletonRaw, &group); err != nil {
		return nil, fmt.Errorf("carregar esqueleto do formulário: %w", err)
	}
	if len(group.Forms) == 0 {
		return nil, fmt.Errorf("esqueleto do formulário não tem nenhum form")
	}
	form := &group.Forms[0]

	textByName := make(map[string]*pdfTextField, len(form.TextFields))
	for _, tf := range form.TextFields {
		textByName[tf.Name] = tf
	}
	checkByName := make(map[string]*pdfCheckBox, len(form.CheckBoxes))
	for _, cb := range form.CheckBoxes {
		checkByName[cb.Name] = cb
	}

	setText := func(fieldName, value string) {
		if tf, ok := textByName[fieldName]; ok {
			tf.Value = value
		}
	}
	setCheck := func(fieldName string, checked bool) {
		if cb, ok := checkByName[fieldName]; ok {
			cb.Value = checked
		}
	}

	// Identidade
	setText(fm.Identidade.NomePersonagem, str(character, "nome"))
	setText(fm.Identidade.ClasseNivel, str(character, "classe_nivel"))
	setText(fm.Identidade.Antecedente, str(character, "antecedente"))
	setText(fm.Identidade.Raca, str(character, "raca"))
	setText(fm.Identidade.Alinhamento, str(character, "alinhamento"))
	setText(fm.Identidade.XP, str(character, "xp"))

	// Atributos
	atributos := subMap(character, "atributos")
	for attr, campos := range fm.Atributos {
		dados := subMap(atributos, attr)
		setText(campos.Valor, str(dados, "valor"))
		if _, ok := dados["mod"]; ok {
			setText(campos.Mod, signed(intVal(dados, "mod")))
		} else {
			setText(campos.Mod, "")
		}
	}

	// Combate
	setText(fm.Combate.BonusProficiencia, signed(intVal(character, "bonus_proficiencia")))
	setText(fm.Combate.CA, str(character, "ca"))
	setText(fm.Combate.Iniciativa, signed(intVal(character, "iniciativa")))
	setText(fm.Combate.Deslocamento, str(character, "deslocamento"))
	setText(fm.Combate.PVMaximo, str(character, "pv_maximo"))
	// PV atual: nunca sobrescrever com valor calculado — usar exatamente o
	// que o jogador tem registrado (bug conhecido #2 do projeto: PV manual
	// não pode ser sobrescrito por cálculo automático).
	setText(fm.Combate.PVAtual, str(character, "pv_atual"))
	setText(fm.Combate.PVTemporario, str(character, "pv_temporario"))
	setText(fm.Combate.DadosDeVidaTotal, str(character, "dados_de_vida_total"))
	setText(fm.Combate.DadosDeVida, str(character, "dados_de_vida"))

	// Testes contra a morte
	morte := subMap(character, "resistencia_morte")
	sucessos := intVal(morte, "sucessos")
	falhas := intVal(morte, "falhas")
	for i, campo := range fm.ResistenciaContraMorte.Sucessos {
		setCheck(campo, i < sucessos)
	}
	for i, campo := range fm.ResistenciaContraMorte.Falhas {
		setCheck(campo, i < falhas)
	}

	// Salvaguardas
	salvaguardas := subMap(character, "salvaguardas")
	for _, save := range fm.Salvaguardas {
		dados := subMap(salvaguardas, save.Atributo)
		setText(save.CampoValor, signed(intVal(dados, "valor")))
		setCheck(save.CampoProficiencia, boolVal(dados, "proficiente"))
	}

	// Perícias
	pericias := subMap(character, "pericias")
	for _, pericia := range fm.Pericias {
		dados := subMap(pericias, pericia.NomePt)
		setText(pericia.CampoValor, signed(intVal(dados, "valor")))
		setCheck(pericia.CampoProficiencia, boolVal(dados, "proficiente"))
	}

	// Passiva
	setText(fm.Passiva.PercepcaoPassiva, str(character, "percepcao_passiva"))

	// Personalidade
	setText(fm.Personalidade.TracosPersonalidade, str(character, "tracos_personalidade"))
	setText(fm.Personalidade.Ideais, str(character, "ideais"))
	setText(fm.Personalidade.Vinculos, str(character, "vinculos"))
	setText(fm.Personalidade.Defeitos, str(character, "defeitos"))

	filledJSON, err := json.Marshal(group)
	if err != nil {
		return nil, fmt.Errorf("montar JSON de preenchimento: %w", err)
	}

	var out bytes.Buffer
	if err := api.FillForm(bytes.NewReader(templatePDF), bytes.NewReader(filledJSON), &out, nil); err != nil {
		return nil, fmt.Errorf("preencher PDF: %w", err)
	}

	return out.Bytes(), nil
}
