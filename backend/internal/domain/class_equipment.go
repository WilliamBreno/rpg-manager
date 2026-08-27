package domain

import "gorm.io/gorm"

// ClassEquipmentOption representa uma opção lettered ("A", "B", ou "C" no
// caso do Guerreiro) do "Equipamento Inicial" de uma classe (Livro do
// Jogador 2024, cap. 3, tabela "Traços Básicos de X" de cada classe). Cada
// classe 5e tem 2 ou 3 opções mutuamente exclusivas: escolher a letra do
// pacote de itens, ou trocar tudo por uma quantia fixa em ouro.
type ClassEquipmentOption struct {
    gorm.Model
    ClassID     uint                      `json:"class_id"`
    Edition     string                    `json:"edition"`
    OptionLabel string                    `json:"option_label"` // "A", "B", "C"
    GoldPieces  int                       `json:"gold_pieces"`  // PO concedidas junto (troco da opção, ou o total se for a opção só-ouro)
    Components  []ClassEquipmentComponent `json:"components" gorm:"foreignKey:OptionID"`
}

// ClassEquipmentComponent é um item dentro de uma ClassEquipmentOption.
// Referencia um Item OU um Armor do catálogo existente (nunca os dois) para
// que o front-end possa mostrar os bônus reais (CA da armadura, dano/bônus
// de ataque da arma, já parseado da Description). ExtraText cobre escolhas
// genéricas que o livro deixa em aberto e que não têm um item catalogado
// único (ex.: "Instrumento Musical à sua escolha", "Ferramentas de Artesão
// à sua escolha" — esta última tem 17 variantes no livro, catalogar todas
// está fora do escopo desta feature).
type ClassEquipmentComponent struct {
    gorm.Model
    OptionID  uint   `json:"option_id"`
    ItemID    *uint  `json:"item_id,omitempty"`
    ArmorID   *uint  `json:"armor_id,omitempty"`
    Quantity  int    `json:"quantity"`
    ExtraText string `json:"extra_text,omitempty"`

    Item  *Item  `json:"item,omitempty" gorm:"foreignKey:ItemID"`
    Armor *Armor `json:"armor,omitempty" gorm:"foreignKey:ArmorID"`
}
