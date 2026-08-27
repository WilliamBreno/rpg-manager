package seed

import (
	"log"

	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

// seedClassEquipment5e povoa o "Equipamento Inicial" de cada classe 5e
// (Livro do Jogador 2024, cap. 3, tabela "Traços Básicos de X" de cada
// classe — linha "Equipamento Inicial"). Cada classe tem 2 opções lettered
// (A/B), exceto Guerreiro que tem 3 (A/B/C): ou o pacote de itens daquela
// letra, ou uma quantia fixa em PO no lugar de qualquer item.
//
// Os nomes de Item/Armor referenciados aqui têm que bater exatamente com os
// nomes seedados por seedItems5e/seedArmors5e (rodam antes desta, ver
// seed.go) — resolvidos em runtime via getItem5e/getArmor5e, não por ID fixo.
// Duas escolhas genéricas do livro (Instrumento Musical à escolha, e
// Ferramentas de Artesão à escolha — esta com 17 variantes) não têm um item
// único catalogável e ficam como texto livre (extra), mesmo padrão já usado
// em Antecedent.ToolProficiencies.
type equipComponent struct {
	ItemName  string
	ArmorName string
	Qty       int
	Extra     string
}

func item(name string, qty int) equipComponent  { return equipComponent{ItemName: name, Qty: qty} }
func armorComp(name string, qty int) equipComponent {
	return equipComponent{ArmorName: name, Qty: qty}
}
func extra(text string) equipComponent { return equipComponent{Extra: text} }

type equipOption struct {
	Label string
	Gold  int
	Comps []equipComponent
}

type classEquip struct {
	ClassName string
	Options   []equipOption
}

func classEquipmentData5e() []classEquip {
	return []classEquip{
		{"Bárbaro", []equipOption{
			{"A", 15, []equipComponent{item("Machadinha", 4), item("Machado Grande", 1), item("Kit de Aventureiro", 1)}},
			{"B", 75, nil},
		}},
		{"Bardo", []equipOption{
			{"A", 19, []equipComponent{armorComp("Armadura de Couro", 1), item("Adaga", 2), extra("Um Instrumento Musical à sua escolha"), item("Kit de Artista", 1)}},
			{"B", 90, nil},
		}},
		{"Bruxo", []equipOption{
			{"A", 15, []equipComponent{armorComp("Armadura de Couro", 1), item("Foice", 1), item("Adaga", 2), item("Foco Arcano: Orbe", 1), item("Livro", 1), item("Kit de Erudito", 1)}},
			{"B", 100, nil},
		}},
		{"Clérigo", []equipOption{
			{"A", 7, []equipComponent{armorComp("Cota de Malha Parcial", 1), armorComp("Escudo", 1), item("Maça", 1), item("Símbolo Sagrado", 1), item("Kit de Sacerdote", 1)}},
			{"B", 110, nil},
		}},
		{"Druida", []equipOption{
			{"A", 9, []equipComponent{armorComp("Armadura de Couro", 1), armorComp("Escudo", 1), item("Foice", 1), item("Foco Druídico: Cajado de Madeira", 1), item("Kit de Herbalismo", 1)}},
			{"B", 50, nil},
		}},
		{"Feiticeiro", []equipOption{
			{"A", 28, []equipComponent{item("Lança", 1), item("Adaga", 2), item("Foco Arcano: Cristal", 1), item("Kit de Explorador de Masmorras", 1)}},
			{"B", 50, nil},
		}},
		{"Guardião", []equipOption{
			{"A", 7, []equipComponent{
				armorComp("Armadura de Couro Batido", 1), item("Cimitarra", 1), item("Espada Curta", 1),
				item("Arco Longo", 1), item("Munição: Flechas (20)", 1), item("Aljava", 1),
				item("Foco Druídico: Ramo de Visco", 1), item("Kit de Aventureiro", 1),
			}},
			{"B", 150, nil},
		}},
		{"Guerreiro", []equipOption{
			{"A", 4, []equipComponent{armorComp("Cota de Malha", 1), item("Espada Grande", 1), item("Mangual", 1), item("Azagaia", 8), item("Kit de Explorador de Masmorras", 1)}},
			{"B", 11, []equipComponent{
				armorComp("Armadura de Couro Batido", 1), item("Cimitarra", 1), item("Espada Curta", 1),
				item("Arco Longo", 1), item("Munição: Flechas (20)", 1), item("Aljava", 1),
				item("Kit de Explorador de Masmorras", 1),
			}},
			{"C", 155, nil},
		}},
		{"Ladino", []equipOption{
			{"A", 8, []equipComponent{
				armorComp("Armadura de Couro", 1), item("Adaga", 2), item("Espada Curta", 1),
				item("Arco Curto", 1), item("Munição: Flechas (20)", 1), item("Aljava", 1),
				item("Ferramentas de Ladrão", 1), item("Kit de Assaltante", 1),
			}},
			{"B", 100, nil},
		}},
		{"Mago", []equipOption{
			{"A", 5, []equipComponent{item("Adaga", 2), item("Foco Arcano: Cajado", 1), item("Kit de Erudito", 1), item("Livro de Magias", 1), item("Túnica", 1)}},
			{"B", 55, nil},
		}},
		{"Monge", []equipOption{
			{"A", 11, []equipComponent{item("Lança", 1), item("Adaga", 5), extra("Ferramentas de Artesão ou Instrumento Musical à sua escolha"), item("Kit de Aventureiro", 1)}},
			{"B", 50, nil},
		}},
		{"Paladino", []equipOption{
			{"A", 9, []equipComponent{armorComp("Cota de Malha", 1), armorComp("Escudo", 1), item("Espada Longa", 1), item("Azagaia", 6), item("Símbolo Sagrado", 1), item("Kit de Sacerdote", 1)}},
			{"B", 150, nil},
		}},
	}
}

func getItem5e(db *gorm.DB, name string) (uint, bool) {
	var it domain.Item
	if err := db.Where("name = ? AND edition = ?", name, "5e").First(&it).Error; err != nil {
		log.Printf("  ✗ Item 5e não encontrado pro equipamento inicial: %s", name)
		return 0, false
	}
	return it.ID, true
}

func getArmor5e(db *gorm.DB, name string) (uint, bool) {
	var ar domain.Armor
	if err := db.Where("name = ? AND edition = ?", name, "5e").First(&ar).Error; err != nil {
		log.Printf("  ✗ Armadura 5e não encontrada pro equipamento inicial: %s", name)
		return 0, false
	}
	return ar.ID, true
}

func seedClassEquipment5e(db *gorm.DB) {
	total := 0
	for _, ce := range classEquipmentData5e() {
		classID, ok := getClass5e(db, ce.ClassName)
		if !ok {
			continue
		}
		for _, opt := range ce.Options {
			var existing domain.ClassEquipmentOption
			err := db.Where("class_id = ? AND edition = ? AND option_label = ?", classID, "5e", opt.Label).
				First(&existing).Error
			if err != nil {
				existing = domain.ClassEquipmentOption{
					ClassID: classID, Edition: "5e", OptionLabel: opt.Label, GoldPieces: opt.Gold,
				}
				db.Create(&existing)
			} else {
				db.Model(&existing).Update("gold_pieces", opt.Gold)
				// Reconstrói os componentes do zero — mais simples que diff e
				// idempotente (a opção inteira é definida por este seed).
				db.Where("option_id = ?", existing.ID).Delete(&domain.ClassEquipmentComponent{})
			}

			for _, comp := range opt.Comps {
				c := domain.ClassEquipmentComponent{OptionID: existing.ID, Quantity: comp.Qty}
				switch {
				case comp.ItemName != "":
					if id, ok := getItem5e(db, comp.ItemName); ok {
						c.ItemID = &id
					} else {
						continue
					}
				case comp.ArmorName != "":
					if id, ok := getArmor5e(db, comp.ArmorName); ok {
						c.ArmorID = &id
					} else {
						continue
					}
				case comp.Extra != "":
					c.ExtraText = comp.Extra
				default:
					continue
				}
				db.Create(&c)
			}
			total++
		}
	}
	log.Println("  ✓ Equipamento Inicial 5e seedado:", total, "opções (12 classes, 2-3 opções cada)")
}
