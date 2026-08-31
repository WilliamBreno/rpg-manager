package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CRDamageStats espelha uma linha da tabela "Estatísticas de Monstro por
// Nível de Desafio" (Guia do Mestre, p.275 — extraído verbatim via PyMuPDF,
// não de memória) — é a referência real do livro pra CA/PV/bônus de
// ataque/dano por rodada/CD de resistência esperados em cada ND. Usada aqui
// só pra sugerir e avisar sobre a faixa de dano de uma EnemyAbility; os
// outros campos (AC, HPMin/Max, AtkBonus, SaveDC) ficam guardados pra
// reaproveitamento futuro (ex: sugerir CA/PV ao criar um Enemy), não usados
// ainda por nenhum handler.
type CRDamageStats struct {
	CR                string `json:"cr"`
	ProficiencyBonus  int    `json:"proficiency_bonus"`
	ArmorClass        int    `json:"armor_class"`
	HPMin             int    `json:"hp_min"`
	HPMax             int    `json:"hp_max"`
	AttackBonus       int    `json:"attack_bonus"`
	DamagePerRoundMin int    `json:"damage_per_round_min"`
	DamagePerRoundMax int    `json:"damage_per_round_max"`
	SaveDC            int    `json:"save_dc"`
}

// CRDamageTable é a tabela completa, ND 0 a 30, na ordem do livro.
var CRDamageTable = []CRDamageStats{
	{"0", 2, 13, 1, 6, 3, 0, 1, 13},
	{"1/8", 2, 13, 7, 35, 3, 2, 3, 13},
	{"1/4", 2, 13, 36, 49, 3, 4, 5, 13},
	{"1/2", 2, 13, 50, 70, 3, 6, 8, 13},
	{"1", 2, 13, 71, 85, 3, 9, 14, 13},
	{"2", 2, 13, 86, 100, 3, 15, 20, 13},
	{"3", 2, 13, 101, 115, 4, 21, 26, 13},
	{"4", 2, 14, 116, 130, 5, 27, 32, 14},
	{"5", 3, 15, 131, 145, 6, 33, 38, 15},
	{"6", 3, 15, 146, 160, 6, 39, 44, 15},
	{"7", 3, 15, 161, 175, 6, 45, 50, 15},
	{"8", 3, 16, 176, 190, 7, 51, 56, 16},
	{"9", 4, 16, 191, 205, 7, 57, 62, 16},
	{"10", 4, 17, 206, 220, 7, 63, 68, 16},
	{"11", 4, 17, 221, 235, 8, 69, 74, 17},
	{"12", 4, 17, 236, 250, 8, 75, 80, 17},
	{"13", 5, 18, 251, 265, 8, 81, 86, 18},
	{"14", 5, 18, 266, 280, 8, 87, 92, 18},
	{"15", 5, 18, 281, 295, 8, 93, 98, 18},
	{"16", 5, 18, 296, 310, 9, 99, 104, 18},
	{"17", 6, 19, 311, 325, 10, 105, 110, 19},
	{"18", 6, 19, 326, 340, 10, 111, 116, 19},
	{"19", 6, 19, 341, 355, 10, 117, 122, 19},
	{"20", 6, 19, 356, 400, 10, 123, 140, 19},
	{"21", 7, 19, 401, 445, 11, 141, 158, 20},
	{"22", 7, 19, 446, 490, 11, 159, 176, 20},
	{"23", 7, 19, 491, 535, 11, 177, 194, 20},
	{"24", 7, 19, 536, 580, 12, 195, 212, 21},
	{"25", 8, 19, 581, 625, 12, 213, 230, 21},
	{"26", 8, 19, 626, 670, 12, 231, 248, 21},
	{"27", 8, 19, 671, 715, 13, 249, 266, 22},
	{"28", 8, 19, 716, 760, 13, 267, 284, 22},
	{"29", 9, 19, 761, 805, 13, 285, 302, 22},
	{"30", 9, 19, 806, 850, 14, 303, 320, 23},
}

func FindCRDamageStats(cr string) (CRDamageStats, bool) {
	cr = strings.TrimSpace(cr)
	for _, s := range CRDamageTable {
		if s.CR == cr {
			return s, true
		}
	}
	return CRDamageStats{}, false
}

var diceNotationRe = regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)
var validDiceSides = map[int]bool{4: true, 6: true, 8: true, 10: true, 12: true, 20: true, 100: true}

// ParseDiceNotation valida e decompõe uma string XdY+Z real de D&D (ex:
// "2d6+3", "1d8", "3d10-2") — rejeita qualquer coisa que não siga essa forma,
// incluindo um número solto sem dado (pedido explícito do SISTEMA_MESTRE.md:
// "o campo de dano não deveria ser um número solto").
func ParseDiceNotation(raw string) (count, sides, mod int, err error) {
	s := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(raw), " ", ""))
	m := diceNotationRe.FindStringSubmatch(s)
	if m == nil {
		return 0, 0, 0, fmt.Errorf("dano deve seguir notação real de dado (ex: 2d6+3), recebido %q", raw)
	}
	count, _ = strconv.Atoi(m[1])
	sides, _ = strconv.Atoi(m[2])
	if m[3] != "" {
		mod, _ = strconv.Atoi(m[3])
	}
	if count < 1 || count > 20 {
		return 0, 0, 0, fmt.Errorf("quantidade de dados deve ser entre 1 e 20")
	}
	if !validDiceSides[sides] {
		return 0, 0, 0, fmt.Errorf("d%d não é um dado válido — use d4, d6, d8, d10, d12, d20 ou d100", sides)
	}
	return count, sides, mod, nil
}

// AverageDiceDamage calcula a média (arredondada pra baixo, como o Guia do
// Mestre faz nos próprios exemplos de dano de monstro) de uma expressão
// XdY+Z já decomposta.
func AverageDiceDamage(count, sides, mod int) int {
	total := float64(count) * (float64(sides) + 1) / 2
	return int(total) + mod
}

// ValidateAbilityDamage garante notação de dado real e, se um ND for
// informado, retorna um aviso (não um erro) quando a média foge muito da
// faixa "Dano/Rodada" sugerida pela tabela do Guia do Mestre pra esse ND. O
// próprio livro trata a faixa como orientação, não regra travada ("não se
// preocupe se o dano causado não combinar com o ND pretendido... você sempre
// pode ajustar depois") — e uma única habilidade normalmente é só uma parte
// do dano total por rodada de um inimigo com múltiplos ataques, então a
// margem de aviso aqui é deliberadamente larga (metade do mínimo a dobro do
// máximo), não a faixa exata.
func ValidateAbilityDamage(damage string, cr string) (avg int, warning string, err error) {
	count, sides, mod, err := ParseDiceNotation(damage)
	if err != nil {
		return 0, "", err
	}
	avg = AverageDiceDamage(count, sides, mod)
	if cr == "" {
		return avg, "", nil
	}
	stats, ok := FindCRDamageStats(cr)
	if !ok {
		return avg, "", nil
	}
	minWarn := stats.DamagePerRoundMin / 2
	maxWarn := stats.DamagePerRoundMax * 2
	if avg < minWarn || avg > maxWarn {
		warning = fmt.Sprintf(
			"dano médio %d está bem fora da faixa sugerida pro ND %s (%d–%d de dano/rodada, Guia do Mestre p.275) — confira se é intencional",
			avg, cr, stats.DamagePerRoundMin, stats.DamagePerRoundMax,
		)
	}
	return avg, warning, nil
}
