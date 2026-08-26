package domain

import "gorm.io/gorm"

// Item é o catálogo de equipamento comprável na loja: armas e equipamento de
// aventura (Capítulo 6 do Livro do Jogador) e um catálogo curado de itens
// mágicos. Armaduras ficam no model Armor (já usado pra CA) — não duplicado
// aqui, só referenciado como uma segunda aba do catálogo da loja.
type Item struct {
    gorm.Model
    Name        string `json:"name"`
    Edition     string `json:"edition"`
    Category    string `json:"category"` // "arma", "equipamento", "item_magico"
    Description string `json:"description"`
    Weight      string `json:"weight"` // texto livre, ex: "0,5 kg" ou "Varia"
    Rarity      string `json:"rarity"` // "", "comum", "incomum", "raro", "muito_raro", "lendario" — só item_magico
    CostCopper  int    `json:"cost_copper"`
    IsDefault   bool   `json:"is_default" gorm:"default:false"`
}

// CharacterItem — inventário: quantidade de cada Item que um personagem possui.
type CharacterItem struct {
    CharacterID uint `json:"character_id" gorm:"primaryKey;not null"`
    ItemID      uint `json:"item_id" gorm:"primaryKey;not null"`
    Quantity    int  `json:"quantity"`

    Item Item `json:"item" gorm:"foreignKey:ItemID"`
}

// CharacterArmorOwned — armaduras compradas/possuídas pelo personagem
// (distinto de "equipada" — a armadura equipada continua sendo
// Character.ArmorID; isso é só o inventário de armaduras adquiridas).
type CharacterArmorOwned struct {
    CharacterID uint `json:"character_id" gorm:"primaryKey;not null"`
    ArmorID     uint `json:"armor_id" gorm:"primaryKey;not null"`
    Quantity    int  `json:"quantity"`

    Armor Armor `json:"armor" gorm:"foreignKey:ArmorID"`
}
