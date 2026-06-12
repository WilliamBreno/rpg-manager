package service

import (
    "rpg-manager/internal/domain"
    "rpg-manager/internal/repository"
)

type ArmorService struct {
    Repo *repository.ArmorRepository
}

func NewArmorService(repo *repository.ArmorRepository) *ArmorService {
    return &ArmorService{Repo: repo}
}

func (s *ArmorService) GetByEdition(edition string) ([]domain.Armor, error) {
    return s.Repo.FindByEdition(edition)
}

// CalculateAC calcula a CA baseado na edição, classe, armadura e atributos
func (s *ArmorService) CalculateAC(character domain.Character) int {
    dexMod := modifier(character.Dexterity)
    wisMod := modifier(character.Wisdom)
    conMod := modifier(character.Constitution)
    intMod := modifier(character.Intelligence)

    edition := character.Edition
    className := character.Class.Name

    // Regras especiais por classe — sem armadura
    if character.ArmorID == nil || character.Armor == nil {
        return specialClassAC(edition, className, dexMod, wisMod, conMod, intMod, character.Level)
    }

    armor := character.Armor
    baseAC := armor.BaseAC

    switch edition {
    case "1e", "2e":
        // Sistema decrescente — escudo subtrai da CA
        if armor.ArmorType == domain.ArmorShield {
            return baseAC // já é -1
        }
        return baseAC

    case "4e":
        // CA = base + metade do nível + DEX ou INT (o maior, só para armaduras leves/sem)
        halfLevel := character.Level / 2
        if armor.ArmorType == domain.ArmorNone || armor.ArmorType == domain.ArmorLight {
            bestMod := dexMod
            if intMod > bestMod {
                bestMod = intMod
            }
            return baseAC + halfLevel + bestMod
        }
        return baseAC + halfLevel

    default: // 3e, 3.5e, 5e
        if armor.ArmorType == domain.ArmorShield {
            // Escudo soma na CA atual do personagem
            return baseAC
        }
        if armor.MaxDexBonus == -1 {
            // Sem limite de DEX
            return baseAC + dexMod
        }
        if armor.MaxDexBonus == 0 {
            // Sem bônus de DEX
            return baseAC
        }
        // Com limite de DEX
        appliedDex := dexMod
        if appliedDex > armor.MaxDexBonus {
            appliedDex = armor.MaxDexBonus
        }
        return baseAC + appliedDex
    }
}

// specialClassAC retorna a CA especial para classes sem armadura
func specialClassAC(edition, className string, dexMod, wisMod, conMod, intMod, level int) int {
    switch edition {
    case "5e":
        switch className {
        case "Monge":
            return 10 + dexMod + wisMod
        case "Bárbaro":
            return 10 + dexMod + conMod
        default:
            return 10 + dexMod
        }
    case "4e":
        halfLevel := level / 2
        bestMod := dexMod
        if intMod > bestMod {
            bestMod = intMod
        }
        return 10 + halfLevel + bestMod
    default:
        return 10 + dexMod
    }
}

// modifier calcula o modificador de atributo padrão do D&D
func modifier(attr int) int {
    return (attr - 10) / 2
}