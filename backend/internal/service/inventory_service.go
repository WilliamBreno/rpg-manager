package service

import (
	"errors"

	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

// InventoryService cuida da carteira de moedas do personagem e das compras na
// loja (itens e armaduras). Não reimplementa cálculo de regras — só aritmética
// de moeda e persistência do inventário.
type InventoryService struct{ DB *gorm.DB }

func NewInventoryService(db *gorm.DB) *InventoryService { return &InventoryService{DB: db} }

// coinToCopper — conversão canônica das 5 moedas de D&D 5e.
const (
	copperPerSilver   = 10
	copperPerElectrum = 50
	copperPerGold     = 100
	copperPerPlatinum = 1000
)

func totalCopper(c domain.Character) int {
	return c.CopperPieces +
		c.SilverPieces*copperPerSilver +
		c.ElectrumPieces*copperPerElectrum +
		c.GoldPieces*copperPerGold +
		c.PlatinumPieces*copperPerPlatinum
}

// copperToPurse "refunde" um total em cobre na combinação de moedas de maior
// valor possível (mesma lógica que um jogador faria na mesa ao trocar moedas).
func copperToPurse(total int) (cp, sp, ep, gp, pp int) {
	pp = total / copperPerPlatinum
	total %= copperPerPlatinum
	gp = total / copperPerGold
	total %= copperPerGold
	ep = total / copperPerElectrum
	total %= copperPerElectrum
	sp = total / copperPerSilver
	total %= copperPerSilver
	cp = total
	return
}

func (s *InventoryService) getCharacter(characterID uint) (domain.Character, error) {
	var character domain.Character
	return character, s.DB.First(&character, characterID).Error
}

// SetCurrency — usado pelo Mestre pra "entregar" ouro: define os totais
// absolutos das 5 moedas do personagem (não é um ajuste incremental).
func (s *InventoryService) SetCurrency(characterID uint, cp, sp, ep, gp, pp int) (domain.Character, error) {
	character, err := s.getCharacter(characterID)
	if err != nil {
		return character, err
	}
	if cp < 0 || sp < 0 || ep < 0 || gp < 0 || pp < 0 {
		return character, errors.New("valores de moeda não podem ser negativos")
	}
	character.CopperPieces, character.SilverPieces = cp, sp
	character.ElectrumPieces, character.GoldPieces, character.PlatinumPieces = ep, gp, pp
	return character, s.DB.Save(&character).Error
}

// GrantCurrency soma as moedas informadas ao que o personagem já tem
// (diferente de SetCurrency, que define totais absolutos) — usado pelo
// RewardService quando o mestre dá moeda como recompensa (Etapa 8 do
// SISTEMA_MESTRE.md). Redenomina o total inteiro pelas moedas de maior valor
// possível, mesma lógica de "troco" já usada em PurchaseItem/PurchaseArmor.
func (s *InventoryService) GrantCurrency(characterID uint, cp, sp, ep, gp, pp int) (domain.Character, error) {
	character, err := s.getCharacter(characterID)
	if err != nil {
		return character, err
	}
	granted := cp + sp*copperPerSilver + ep*copperPerElectrum + gp*copperPerGold + pp*copperPerPlatinum
	newTotal := totalCopper(character) + granted
	ncp, nsp, nep, ngp, npp := copperToPurse(newTotal)
	character.CopperPieces, character.SilverPieces = ncp, nsp
	character.ElectrumPieces, character.GoldPieces, character.PlatinumPieces = nep, ngp, npp
	return character, s.DB.Save(&character).Error
}

// PurchaseItem — debita o custo (item.CostCopper * quantity) da carteira e
// credita a quantidade no inventário do personagem.
func (s *InventoryService) PurchaseItem(characterID, itemID uint, quantity int) (domain.Character, error) {
	if quantity <= 0 {
		return domain.Character{}, errors.New("quantidade deve ser maior que zero")
	}
	character, err := s.getCharacter(characterID)
	if err != nil {
		return character, err
	}
	var item domain.Item
	if err := s.DB.First(&item, itemID).Error; err != nil {
		return character, errors.New("item não encontrado")
	}

	cost := item.CostCopper * quantity
	available := totalCopper(character)
	if cost > available {
		return character, errors.New("ouro insuficiente para essa compra")
	}

	cp, sp, ep, gp, pp := copperToPurse(available - cost)
	character.CopperPieces, character.SilverPieces = cp, sp
	character.ElectrumPieces, character.GoldPieces, character.PlatinumPieces = ep, gp, pp
	if err := s.DB.Save(&character).Error; err != nil {
		return character, err
	}

	var existing domain.CharacterItem
	err = s.DB.Where("character_id = ? AND item_id = ?", characterID, itemID).First(&existing).Error
	if err != nil {
		if err := s.DB.Create(&domain.CharacterItem{CharacterID: characterID, ItemID: itemID, Quantity: quantity}).Error; err != nil {
			return character, err
		}
	} else {
		existing.Quantity += quantity
		if err := s.DB.Save(&existing).Error; err != nil {
			return character, err
		}
	}

	return character, nil
}

// PurchaseArmor — mesma lógica de PurchaseItem, mas pro catálogo de Armor.
func (s *InventoryService) PurchaseArmor(characterID, armorID uint, quantity int) (domain.Character, error) {
	if quantity <= 0 {
		return domain.Character{}, errors.New("quantidade deve ser maior que zero")
	}
	character, err := s.getCharacter(characterID)
	if err != nil {
		return character, err
	}
	var armor domain.Armor
	if err := s.DB.First(&armor, armorID).Error; err != nil {
		return character, errors.New("armadura não encontrada")
	}

	cost := armor.CostCopper * quantity
	available := totalCopper(character)
	if cost > available {
		return character, errors.New("ouro insuficiente para essa compra")
	}

	cp, sp, ep, gp, pp := copperToPurse(available - cost)
	character.CopperPieces, character.SilverPieces = cp, sp
	character.ElectrumPieces, character.GoldPieces, character.PlatinumPieces = ep, gp, pp
	if err := s.DB.Save(&character).Error; err != nil {
		return character, err
	}

	var existing domain.CharacterArmorOwned
	err = s.DB.Where("character_id = ? AND armor_id = ?", characterID, armorID).First(&existing).Error
	if err != nil {
		if err := s.DB.Create(&domain.CharacterArmorOwned{CharacterID: characterID, ArmorID: armorID, Quantity: quantity}).Error; err != nil {
			return character, err
		}
	} else {
		existing.Quantity += quantity
		if err := s.DB.Save(&existing).Error; err != nil {
			return character, err
		}
	}

	return character, nil
}

func (s *InventoryService) GetInventory(characterID uint) ([]domain.CharacterItem, []domain.CharacterArmorOwned, error) {
	var items []domain.CharacterItem
	if err := s.DB.Preload("Item").Where("character_id = ?", characterID).Find(&items).Error; err != nil {
		return nil, nil, err
	}
	var armors []domain.CharacterArmorOwned
	if err := s.DB.Preload("Armor").Where("character_id = ?", characterID).Find(&armors).Error; err != nil {
		return nil, nil, err
	}
	return items, armors, nil
}
