package domain

import "gorm.io/gorm"

// Ritual é o catálogo de rituais de D&D 4e (Livro do Jogador 1 cap. 10 +
// rituais novos do Livro do Jogador 2 cap. 3 — o Livro do Jogador 3 não
// acrescenta nenhum ritual, confirmado por varredura completa do PDF).
// Diferente de Spell (5e), um ritual não é "por classe" — qualquer
// personagem com a característica Conjuração Ritual pode aprender qualquer
// ritual de nível igual ou inferior ao seu, exceto quando o próprio livro
// lista um Pré-requisito de classe (ex: vários rituais do LJ2 exigem
// "Bardo") — daí o campo Prerequisite ser uma classe específica, não um
// mapa de classes como em Spell.
type Ritual struct {
	gorm.Model
	Name    string `json:"name"`
	Edition string `json:"edition"` // sempre "4e"

	Level    int    `json:"level"`    // nível do ritual (1-30), não o nível do personagem
	Category string `json:"category"` // Contenção, Criação, Enganação, Adivinhação, Exploração, Restauração, Sondagem, Viagem, Proteção
	KeySkill string `json:"key_skill"` // perícia-chave, ex: "Arcanismo", "Arcanismo ou Natureza"

	ComponentCost string `json:"component_cost"` // custo dos componentes gasto ao executar, ex: "25 PO"
	MarketPrice   string `json:"market_price"`   // preço pra dominar/copiar o ritual, ex: "125 PO"
	Prerequisite  string `json:"prerequisite"`   // nome da classe exigida, "" se não houver

	SourceBook  string `json:"source_book"`  // "Livro do Jogador 1" ou "Livro do Jogador 2"
	Description string `json:"description"`  // efeito resumido
}
