package seed

import (
	"log"

	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

// seedArmors5e povoa o catálogo de armaduras 5e (Capítulo 6 do Livro do
// Jogador) — a tabela estava vazia até aqui, então isso também corrige a
// ausência total de armaduras equipáveis no sistema, não só a loja.
//
// "Sem Armadura" (CA base 10, sem limite de DEX) precisava entrar aqui
// explicitamente: já existia uma linha manual com esse nome no banco de dev
// (sem nenhum código que a criasse ou mantivesse), então em qualquer ambiente
// sem esse registro à mão (produção, um banco recriado do zero) o dropdown de
// armadura na criação de personagem simplesmente não tinha opção nenhuma de
// "sem armadura" real — só o placeholder vazio do <select>, que não é um
// armor_id válido e travava a criação silenciosamente se o jogador não
// escolhesse manualmente uma armadura de verdade. Ver CharacterCreate.tsx:
// o form agora pré-seleciona esta linha (armor_type == "none") como padrão.
func seedArmors5e(db *gorm.DB) {
	type a struct {
		Name        string
		ArmorType   domain.ArmorType
		BaseAC      int
		MaxDexBonus int
		Weight      string
		CostCopper  int
	}
	armors := []a{
		{"Sem Armadura", domain.ArmorNone, 10, -1, "", 0},
		{"Armadura Acolchoada", domain.ArmorLight, 11, -1, "4 kg", 500},
		{"Armadura de Couro", domain.ArmorLight, 11, -1, "5 kg", 1000},
		{"Armadura de Couro Batido", domain.ArmorLight, 12, -1, "6,5 kg", 4500},
		{"Gibão de Peles", domain.ArmorMedium, 12, 2, "6 kg", 1000},
		{"Cota de Malha Parcial", domain.ArmorMedium, 13, 2, "10 kg", 5000},
		{"Loriga de Escamas", domain.ArmorMedium, 14, 2, "22 kg", 5000},
		{"Couraça Peitoral", domain.ArmorMedium, 14, 2, "10 kg", 40000},
		{"Armadura de Placas Parcial", domain.ArmorMedium, 15, 2, "20 kg", 75000},
		{"Cota de Anéis", domain.ArmorHeavy, 14, 0, "20 kg", 3000},
		{"Cota de Malha", domain.ArmorHeavy, 16, 0, "27 kg", 7500},
		{"Armadura de Talas", domain.ArmorHeavy, 17, 0, "30 kg", 20000},
		{"Armadura de Placas", domain.ArmorHeavy, 18, 0, "32 kg", 150000},
		{"Escudo", domain.ArmorShield, 2, 0, "3 kg", 1000},
	}
	for _, ar := range armors {
		var existing domain.Armor
		if db.Where("name = ? AND edition = ?", ar.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Armor{
				Name: ar.Name, Edition: "5e", ArmorType: ar.ArmorType,
				BaseAC: ar.BaseAC, MaxDexBonus: ar.MaxDexBonus,
				CostCopper: ar.CostCopper, Weight: ar.Weight,
			})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"cost_copper": ar.CostCopper, "weight": ar.Weight,
			})
		}
	}
	log.Println("  ✓ Armaduras 5e seedadas:", len(armors))
}

// seedItems5e povoa armas, equipamento de aventura e um catálogo curado de
// itens mágicos (preço = média da faixa sugerida por raridade no Guia de
// Xanathar — o D&D 5e não define preço fixo de tabela pra itens mágicos,
// isso é uma simplificação deliberada pra viabilizar compra direta na loja).
func seedItems5e(db *gorm.DB) {
	type it struct {
		Name, Category, Description, Weight, Rarity string
		CostCopper                                   int
	}
	items := []it{}
	items = append(items, it{"Adaga", "arma", "1d4 Perfurante — Acuidade, Arremesso (Alcance 6/18), Leve (Maestria: Ágil)", "0,5 kg", "", 200})
	items = append(items, it{"Azagaia", "arma", "1d6 Perfurante — Arremesso (Alcance 9/36) (Maestria: Lentidão)", "1 kg", "", 50})
	items = append(items, it{"Cajado", "arma", "1d6 Contundente — Versátil (1d8) (Maestria: Derrubar)", "2 kg", "", 20})
	items = append(items, it{"Clava", "arma", "1d4 Contundente — Leve (Maestria: Lentidão)", "1 kg", "", 10})
	items = append(items, it{"Clava Grande", "arma", "1d8 Contundente — Duas Mãos (Maestria: Empurrar)", "5 kg", "", 20})
	items = append(items, it{"Foice", "arma", "1d4 Cortante — Leve (Maestria: Ágil)", "1 kg", "", 100})
	items = append(items, it{"Lança", "arma", "1d6 Perfurante — Arremesso (Alcance 6/18), Versátil (1d8) (Maestria: Drenar)", "1,5 kg", "", 100})
	items = append(items, it{"Maça", "arma", "1d6 Contundente — — (Maestria: Drenar)", "2 kg", "", 500})
	items = append(items, it{"Machadinha", "arma", "1d6 Cortante — Arremesso (Alcance 6/18), Leve (Maestria: Afligir)", "1 kg", "", 500})
	items = append(items, it{"Martelo Leve", "arma", "1d4 Contundente — Arremesso (Alcance 6/18), Leve (Maestria: Ágil)", "1 kg", "", 200})
	items = append(items, it{"Arco Curto", "arma", "1d6 Perfurante — Duas Mãos, Munição (Alcance 24/96; Flecha) (Maestria: Afligir)", "1 kg", "", 2500})
	items = append(items, it{"Besta Leve", "arma", "1d8 Perfurante — Duas Mãos, Munição (Alcance 24/96; Virote), Recarga (Maestria: Lentidão)", "2,5 kg", "", 2500})
	items = append(items, it{"Dardo", "arma", "1d4 Perfurante — Acuidade, Arremesso (Alcance 6/18) (Maestria: Afligir)", "150 g", "", 5})
	items = append(items, it{"Funda", "arma", "1d4 Contundente — Munição (Alcance 9/36; Bala) (Maestria: Lentidão)", "—", "", 10})
	items = append(items, it{"Alabarda", "arma", "1d10 Cortante — Duas Mãos, Extensão, Pesada (Maestria: Trespassar)", "3 kg", "", 2000})
	items = append(items, it{"Chicote", "arma", "1d4 Cortante — Acuidade, Extensão (Maestria: Lentidão)", "1,5 kg", "", 200})
	items = append(items, it{"Cimitarra", "arma", "1d6 Cortante — Acuidade, Leve (Maestria: Ágil)", "1,5 kg", "", 2500})
	items = append(items, it{"Espada Curta", "arma", "1d6 Perfurante — Acuidade, Leve (Maestria: Afligir)", "1 kg", "", 1000})
	items = append(items, it{"Espada Grande", "arma", "2d6 Cortante — Duas Mãos, Pesada (Maestria: Garantido)", "3 kg", "", 5000})
	items = append(items, it{"Espada Longa", "arma", "1d8 Cortante — Versátil (1d10) (Maestria: Drenar)", "1,5 kg", "", 1500})
	items = append(items, it{"Glaive", "arma", "1d10 Cortante — Duas Mãos, Extensão, Pesada (Maestria: Garantido)", "3 kg", "", 2000})
	items = append(items, it{"Lança de Montaria", "arma", "1d10 Perfurante — Duas Mãos (a menos que montado), Extensão, Pesada (Maestria: Derrubar)", "3 kg", "", 1000})
	items = append(items, it{"Lança Longa", "arma", "1d10 Perfurante — Duas Mãos, Extensão, Pesada (Maestria: Empurrar)", "9 kg", "", 500})
	items = append(items, it{"Maça Estrela", "arma", "1d8 Perfurante — — (Maestria: Drenar)", "2 kg", "", 1500})
	items = append(items, it{"Machado de Batalha", "arma", "1d8 Cortante — Versátil (1d10) (Maestria: Derrubar)", "2,5 kg", "", 1000})
	items = append(items, it{"Machado Grande", "arma", "1d12 Cortante — Duas Mãos, Pesada (Maestria: Trespassar)", "3,5 kg", "", 3000})
	items = append(items, it{"Malho", "arma", "2d6 Contundente — Duas Mãos, Pesada (Maestria: Derrubar)", "5 kg", "", 1000})
	items = append(items, it{"Mangual", "arma", "1d8 Contundente — — (Maestria: Drenar)", "1 kg", "", 1000})
	items = append(items, it{"Martelo de Guerra", "arma", "1d8 Contundente — Versátil (1d10) (Maestria: Empurrar)", "1 kg", "", 1500})
	items = append(items, it{"Picareta de Guerra", "arma", "1d8 Perfurante — Versátil (1d10) (Maestria: Drenar)", "1 kg", "", 500})
	items = append(items, it{"Rapieira", "arma", "1d8 Perfurante — Acuidade (Maestria: Afligir)", "1 kg", "", 2500})
	items = append(items, it{"Tridente", "arma", "1d8 Perfurante — Arremesso (Alcance 6/18), Versátil (1d10) (Maestria: Derrubar)", "2 kg", "", 500})
	items = append(items, it{"Arco Longo", "arma", "1d8 Perfurante — Duas Mãos, Munição (Alcance 45/180; Flecha), Pesada (Maestria: Lentidão)", "1 kg", "", 5000})
	items = append(items, it{"Besta de Mão", "arma", "1d6 Perfurante — Leve, Munição (Alcance 9/36; Virote), Recarga (Maestria: Afligir)", "1,5 kg", "", 7500})
	items = append(items, it{"Besta Pesada", "arma", "1d10 Perfurante — Duas Mãos, Munição (Alcance 30/120; Virote), Pesada, Recarga (Maestria: Empurrar)", "9 kg", "", 5000})
	items = append(items, it{"Mosquete", "arma", "1d12 Perfurante — Duas Mãos, Munição (Alcance 12/36; Bala), Recarga (Maestria: Lentidão)", "5 kg", "", 50000})
	items = append(items, it{"Pistola", "arma", "1d10 Perfurante — Munição (Alcance 9/27; Bala), Recarga (Maestria: Afligir)", "1,5 kg", "", 25000})
	items = append(items, it{"Zarabatana", "arma", "1 Perfurante — Munição (Alcance 7,5/30; Agulha), Recarga (Maestria: Afligir)", "0,5 kg", "", 1000})
	items = append(items, it{"Ácido", "equipamento", "", "0,5 kg", "", 2500})
	items = append(items, it{"Água Benta", "equipamento", "", "0,5 kg", "", 2500})
	items = append(items, it{"Algibeira", "equipamento", "", "0,5 kg", "", 50})
	items = append(items, it{"Aljava", "equipamento", "", "0,5 kg", "", 100})
	items = append(items, it{"Antitoxina", "equipamento", "", "—", "", 5000})
	items = append(items, it{"Apito Sinalizador", "equipamento", "", "—", "", 5})
	items = append(items, it{"Aríete Portável", "equipamento", "", "16,5 kg", "", 400})
	items = append(items, it{"Armadilha de Caça", "equipamento", "", "12,5 kg", "", 500})
	items = append(items, it{"Arpéu", "equipamento", "", "2 kg", "", 200})
	items = append(items, it{"Balde", "equipamento", "", "1 kg", "", 5})
	items = append(items, it{"Baliza", "equipamento", "", "3,5 kg", "", 5})
	items = append(items, it{"Barril", "equipamento", "", "35 kg", "", 200})
	items = append(items, it{"Baú", "equipamento", "", "12,5 kg", "", 500})
	items = append(items, it{"Bolsa de Componentes", "equipamento", "", "1 kg", "", 2500})
	items = append(items, it{"Cadeado", "equipamento", "", "0,5 kg", "", 1000})
	items = append(items, it{"Caixa para Fogo", "equipamento", "", "0,5 kg", "", 50})
	items = append(items, it{"Caneta Tinteiro", "equipamento", "", "—", "", 2})
	items = append(items, it{"Cantil (cheio)", "equipamento", "", "2,5 kg", "", 20})
	items = append(items, it{"Cesta", "equipamento", "", "1 kg", "", 40})
	items = append(items, it{"Cobertor", "equipamento", "", "1,5 kg", "", 50})
	items = append(items, it{"Corda", "equipamento", "", "2,5 kg", "", 100})
	items = append(items, it{"Cordão", "equipamento", "", "—", "", 10})
	items = append(items, it{"Corrente", "equipamento", "", "5 kg", "", 500})
	items = append(items, it{"Escada", "equipamento", "", "12,5 kg", "", 10})
	items = append(items, it{"Esferas de Metal", "equipamento", "", "1 kg (saco)", "", 100})
	items = append(items, it{"Espelho", "equipamento", "", "250 g", "", 500})
	items = append(items, it{"Estacas de Ferro", "equipamento", "", "2,5 kg", "", 100})
	items = append(items, it{"Estojo, Mapa ou Pergaminho", "equipamento", "", "0,5 kg", "", 100})
	items = append(items, it{"Estojo, Virote de Besta", "equipamento", "", "0,5 kg", "", 100})
	items = append(items, it{"Estrepes", "equipamento", "", "1 kg (saco)", "", 100})
	items = append(items, it{"Fogo Alquímico", "equipamento", "", "0,5 kg", "", 5000})
	items = append(items, it{"Frasco", "equipamento", "", "—", "", 100})
	items = append(items, it{"Garrafa de Vidro (1 litro)", "equipamento", "", "1 kg", "", 200})
	items = append(items, it{"Grilhões", "equipamento", "", "3 kg", "", 200})
	items = append(items, it{"Jarro (4 litros)", "equipamento", "", "2 kg", "", 2})
	items = append(items, it{"Kit de Artista", "equipamento", "", "29 kg", "", 4000})
	items = append(items, it{"Kit de Assaltante", "equipamento", "", "21 kg", "", 1600})
	items = append(items, it{"Kit de Aventureiro", "equipamento", "", "27,5 kg", "", 1000})
	items = append(items, it{"Kit de Curandeiro", "equipamento", "", "1,5 kg", "", 500})
	items = append(items, it{"Kit de Diplomata", "equipamento", "", "19,5 kg", "", 3900})
	items = append(items, it{"Kit de Erudito", "equipamento", "", "11 kg", "", 4000})
	items = append(items, it{"Kit de Escalada", "equipamento", "", "6 kg", "", 2500})
	items = append(items, it{"Kit de Explorador de Masmorras", "equipamento", "", "27,5 kg", "", 1200})
	items = append(items, it{"Kit de Sacerdote", "equipamento", "", "14,5 kg", "", 3300})
	items = append(items, it{"Lâmpada", "equipamento", "", "0,5 kg", "", 50})
	items = append(items, it{"Lanterna Coberta", "equipamento", "", "1 kg", "", 500})
	items = append(items, it{"Lanterna Foca-facho", "equipamento", "", "1 kg", "", 1000})
	items = append(items, it{"Livro", "equipamento", "", "2,5 kg", "", 2500})
	items = append(items, it{"Luneta", "equipamento", "", "0,5 kg", "", 100000})
	items = append(items, it{"Lupa", "equipamento", "", "—", "", 10000})
	items = append(items, it{"Mapa", "equipamento", "", "—", "", 100})
	items = append(items, it{"Mochila", "equipamento", "", "2,5 kg", "", 200})
	items = append(items, it{"Óleo", "equipamento", "", "0,5 kg", "", 10})
	items = append(items, it{"Pá", "equipamento", "", "2,5 kg", "", 200})
	items = append(items, it{"Papel", "equipamento", "", "—", "", 20})
	items = append(items, it{"Pé de Cabra", "equipamento", "", "2,5 kg", "", 200})
	items = append(items, it{"Perfume", "equipamento", "", "—", "", 500})
	items = append(items, it{"Pergaminho", "equipamento", "", "—", "", 10})
	items = append(items, it{"Pergaminho Mágico (1º Círculo)", "equipamento", "", "—", "", 5000})
	items = append(items, it{"Pergaminho Mágico (Truque)", "equipamento", "", "—", "", 3000})
	items = append(items, it{"Poção de Cura", "equipamento", "", "250 g", "", 5000})
	items = append(items, it{"Pote", "equipamento", "", "0,5 kg", "", 2})
	items = append(items, it{"Pote de Ferro", "equipamento", "", "5 kg", "", 200})
	items = append(items, it{"Rações (1 dia)", "equipamento", "", "1 kg", "", 50})
	items = append(items, it{"Rede", "equipamento", "", "1,5 kg", "", 100})
	items = append(items, it{"Roldana e Polias", "equipamento", "", "2,5 kg", "", 100})
	items = append(items, it{"Roupas, Fantasia", "equipamento", "", "2 kg", "", 500})
	items = append(items, it{"Roupas, Finas", "equipamento", "", "3 kg", "", 1500})
	items = append(items, it{"Roupas, Viagem", "equipamento", "", "2 kg", "", 200})
	items = append(items, it{"Saca", "equipamento", "", "250 g", "", 1})
	items = append(items, it{"Saco de Dormir", "equipamento", "", "3,5 kg", "", 100})
	items = append(items, it{"Sino", "equipamento", "", "—", "", 100})
	items = append(items, it{"Tenda", "equipamento", "", "10 kg", "", 200})
	items = append(items, it{"Tinta", "equipamento", "", "—", "", 1000})
	items = append(items, it{"Tocha", "equipamento", "", "0,5 kg", "", 1})
	items = append(items, it{"Vela", "equipamento", "", "—", "", 1})
	items = append(items, it{"Túnica", "equipamento", "", "2 kg", "", 100})
	items = append(items, it{"Veneno Básico", "equipamento", "", "—", "", 10000})
	items = append(items, it{"Foco Arcano: Cajado", "equipamento", "", "2 kg", "", 500})
	items = append(items, it{"Foco Arcano: Cetro", "equipamento", "", "1 kg", "", 1000})
	items = append(items, it{"Foco Arcano: Cristal", "equipamento", "", "0,5 kg", "", 1000})
	items = append(items, it{"Foco Arcano: Orbe", "equipamento", "", "1,5 kg", "", 2000})
	items = append(items, it{"Foco Arcano: Varinha", "equipamento", "", "0,5 kg", "", 1000})
	items = append(items, it{"Foco Druídico: Cajado de Madeira", "equipamento", "", "2 kg", "", 500})
	items = append(items, it{"Foco Druídico: Ramo de Visco", "equipamento", "", "—", "", 100})
	items = append(items, it{"Foco Druídico: Varinha de Teixo", "equipamento", "", "0,5 kg", "", 1000})
	items = append(items, it{"Símbolo Sagrado", "equipamento", "", "0,5 kg", "", 500})
	items = append(items, it{"Munição: Agulhas (50)", "equipamento", "", "0,5 kg", "", 100})
	items = append(items, it{"Munição: Balas de Arma de Fogo (10)", "equipamento", "", "1 kg", "", 300})
	items = append(items, it{"Munição: Balas de Funda (20)", "equipamento", "", "750 g", "", 4})
	items = append(items, it{"Munição: Flechas (20)", "equipamento", "", "0,5 kg", "", 100})
	items = append(items, it{"Munição: Virotes (20)", "equipamento", "", "750 g", "", 100})
	items = append(items, it{"Poção de Cura Superior", "item_magico", "Restaura 4d4+4 Pontos de Vida ao beber.", "", "incomum", 35000})
	items = append(items, it{"Poção de Resistência ao Fogo", "item_magico", "Resistência a dano de fogo por 1 hora.", "", "incomum", 35000})
	items = append(items, it{"Pergaminho de Proteção", "item_magico", "Proteção contra um tipo de criatura por 5 minutos.", "", "incomum", 35000})
	items = append(items, it{"Adaga +1", "item_magico", "Arma mágica: +1 nas jogadas de ataque e dano.", "", "incomum", 35000})
	items = append(items, it{"Cota de Malha +1", "item_magico", "Armadura mágica: +1 na CA.", "", "raro", 400000})
	items = append(items, it{"Anel de Proteção", "item_magico", "+1 na CA e em salvaguardas.", "", "raro", 400000})
	items = append(items, it{"Botas de Elfo", "item_magico", "Passos silenciosos; Vantagem em Furtividade (Destreza) relacionados a movimento.", "", "incomum", 35000})
	items = append(items, it{"Capa de Proteção", "item_magico", "+1 na CA e em salvaguardas.", "", "incomum", 35000})
	items = append(items, it{"Poção de Voo", "item_magico", "Voo por 1 hora ao beber.", "", "raro", 400000})
	items = append(items, it{"Amuleto de Saúde", "item_magico", "Constituição se torna 19 enquanto usado, se já não for maior.", "", "raro", 400000})
	items = append(items, it{"Luvas de Destreza do Ladrão", "item_magico", "Bônus de +5 em Prestidigitação, Acrobacia e Furtividade.", "", "raro", 400000})
	items = append(items, it{"Espada Longa +1", "item_magico", "Arma mágica: +1 nas jogadas de ataque e dano.", "", "raro", 400000})
	items = append(items, it{"Bastão de Cura", "item_magico", "10 cargas; gasta cargas pra conjurar magias de cura.", "", "raro", 400000})
	items = append(items, it{"Manto de Resistência Elemental", "item_magico", "Resistência a um tipo de dano escolhido.", "", "raro", 400000})
	items = append(items, it{"Armadura de Placas +1", "item_magico", "Armadura mágica: +1 na CA.", "", "muito_raro", 1150000})
	items = append(items, it{"Espada Vorpal", "item_magico", "Arma lendária: chance de decapitar em acerto crítico.", "", "lendario", 15000000})
	items = append(items, it{"Bastão da Ressurreição", "item_magico", "Conjura ressurreição verdadeira.", "", "lendario", 15000000})
	items = append(items, it{"Manto das Asas de Morcego", "item_magico", "Transforma-se em morcego; voo.", "", "raro", 400000})
	items = append(items, it{"Orbe da Tempestade", "item_magico", "Foco de conjuração que amplia magias de tempestade/relâmpago.", "", "muito_raro", 1150000})
	items = append(items, it{"Anel de Regeneração", "item_magico", "Recupera Pontos de Vida a cada rodada.", "", "muito_raro", 1150000})
	items = append(items, it{"Elmo da Telepatia", "item_magico", "Permite comunicação telepática e ler pensamentos superficiais.", "", "incomum", 35000})
	items = append(items, it{"Bolsa de Contenção", "item_magico", "Espaço interdimensional para guardar itens.", "", "raro", 400000})
	// Faltavam do Equipamento Inicial de classe (capítulo 3) — não estão nas
	// tabelas de compra do capítulo 6, mas são referenciados por nome exato
	// pelo pacote de equipamento de Ladino/Druida/Mago respectivamente.
	items = append(items, it{"Ferramentas de Ladrão", "equipamento", "", "0,5 kg", "", 2500})
	items = append(items, it{"Kit de Herbalismo", "equipamento", "", "1,5 kg", "", 500})
	items = append(items, it{"Livro de Magias", "equipamento", "Grimório em branco com 100 páginas para um Mago registrar magias.", "1,5 kg", "", 5000})

	for _, i := range items {
		var existing domain.Item
		if db.Where("name = ? AND edition = ?", i.Name, "5e").First(&existing).Error != nil {
			db.Create(&domain.Item{
				Name: i.Name, Edition: "5e", Category: i.Category,
				Description: i.Description, Weight: i.Weight,
				Rarity: i.Rarity, CostCopper: i.CostCopper,
			})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"cost_copper": i.CostCopper, "description": i.Description,
				"weight": i.Weight, "rarity": i.Rarity,
			})
		}
	}
	log.Println("  ✓ Itens 5e seedados:", len(items))
}
