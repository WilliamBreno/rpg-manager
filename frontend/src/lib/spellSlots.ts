// Espaços de magia 5e (PHB 2024) — espelha backend/internal/service/character_service.go
// (fullCasterSlots5e/halfCasterSlots5e/pactMagicSlots5e/cantripsKnown5e).
// Não reintroduza uma segunda cópia dessas tabelas em outro arquivo do
// frontend — este é o único lugar, igual já vale para xpTables.ts.

// index = nível de personagem; cada valor é [espaços de 1º círculo, 2º, ...]
const FULL_CASTER_SLOTS: Record<number, number[]> = {
  1: [2], 2: [3], 3: [4, 2], 4: [4, 3], 5: [4, 3, 2],
  6: [4, 3, 3], 7: [4, 3, 3, 1], 8: [4, 3, 3, 2], 9: [4, 3, 3, 3, 1], 10: [4, 3, 3, 3, 2],
  11: [4, 3, 3, 3, 2, 1], 12: [4, 3, 3, 3, 2, 1], 13: [4, 3, 3, 3, 2, 1, 1], 14: [4, 3, 3, 3, 2, 1, 1],
  15: [4, 3, 3, 3, 2, 1, 1, 1], 16: [4, 3, 3, 3, 2, 1, 1, 1], 17: [4, 3, 3, 3, 2, 1, 1, 1, 1],
  18: [4, 3, 3, 3, 3, 1, 1, 1, 1], 19: [4, 3, 3, 3, 3, 2, 1, 1, 1], 20: [4, 3, 3, 3, 3, 2, 2, 1, 1],
}

// Paladino/Guardião — verificado contra o PDF: já ganham espaços no nível 1
// no PHB 2024 (2014 só dava a partir do nível 2).
const HALF_CASTER_SLOTS: Record<number, number[]> = {
  1: [2], 2: [2], 3: [3], 4: [3], 5: [4, 2],
  6: [4, 2], 7: [4, 3], 8: [4, 3], 9: [4, 3, 2], 10: [4, 3, 2],
  11: [4, 3, 3], 12: [4, 3, 3], 13: [4, 3, 3, 1], 14: [4, 3, 3, 1],
  15: [4, 3, 3, 2], 16: [4, 3, 3, 2], 17: [4, 3, 3, 3, 1], 18: [4, 3, 3, 3, 1],
  19: [4, 3, 3, 3, 2], 20: [4, 3, 3, 3, 2],
}

// Terço-conjurador (Cavaleiro Místico/Guerreiro, Trapaceiro Arcano/Ladino) —
// tabela idêntica pras duas subclasses, extraída verbatim das tabelas
// "Conjuração de Cavaleiro Místico"/"Conjuração de Trapaceiro Arcano" do
// PHB 2024 (não existia nenhuma tabela de terço-conjurador no sistema antes
// disso — Guerreiro/Ladino sempre voltavam vazio de spellSlotsFor).
const THIRD_CASTER_SLOTS: Record<number, number[]> = {
  1: [], 2: [], 3: [2], 4: [3], 5: [3],
  6: [3], 7: [4, 2], 8: [4, 2], 9: [4, 2], 10: [4, 3],
  11: [4, 3], 12: [4, 3], 13: [4, 3, 2], 14: [4, 3, 2], 15: [4, 3, 2],
  16: [4, 3, 3], 17: [4, 3, 3], 18: [4, 3, 3], 19: [4, 3, 3, 1], 20: [4, 3, 3, 1],
}
const THIRD_CASTER_SUBCLASSES: Record<string, string> = {
  'Cavaleiro Místico': 'Guerreiro',
  'Trapaceiro Arcano': 'Ladino',
}

const PACT_MAGIC_SLOTS: Record<number, number> = {
  1: 1, 2: 2, 3: 2, 4: 2, 5: 2, 6: 2, 7: 2, 8: 2, 9: 2, 10: 2,
  11: 3, 12: 3, 13: 3, 14: 3, 15: 3, 16: 3, 17: 4, 18: 4, 19: 4, 20: 4,
}
const PACT_MAGIC_CIRCLE: Record<number, number> = {
  1: 1, 2: 1, 3: 2, 4: 2, 5: 3, 6: 3, 7: 4, 8: 4, 9: 5, 10: 5,
  11: 5, 12: 5, 13: 5, 14: 5, 15: 5, 16: 5, 17: 5, 18: 5, 19: 5, 20: 5,
}

const CANTRIP_BASE: Record<string, number> = {
  Bardo: 2, Bruxo: 2, Clérigo: 3, Druida: 2, Feiticeiro: 4, Mago: 3,
}

// subclassName: nome da subclasse escolhida (ex: "Cavaleiro Místico"), só
// relevante pra saber se um Guerreiro/Ladino é terço-conjurador — ignorado
// pras demais classes, que já sabem conjurar (ou não) só pelo nome da classe.
export function spellSlotsFor(className: string, level: number, subclassName?: string): number[] {
  if (['Bardo', 'Clérigo', 'Druida', 'Feiticeiro', 'Mago'].includes(className)) {
    return FULL_CASTER_SLOTS[level] ?? []
  }
  if (['Paladino', 'Guardião'].includes(className)) {
    return HALF_CASTER_SLOTS[level] ?? []
  }
  if (subclassName && THIRD_CASTER_SUBCLASSES[subclassName] === className) {
    return THIRD_CASTER_SLOTS[level] ?? []
  }
  return []
}

export function pactMagicFor(level: number): { slots: number; circle: number } {
  return { slots: PACT_MAGIC_SLOTS[level] ?? 0, circle: PACT_MAGIC_CIRCLE[level] ?? 0 }
}

export function cantripsKnownFor(className: string, level: number): number {
  const base = CANTRIP_BASE[className]
  if (!base) return 0
  let n = base
  if (level >= 4) n++
  if (level >= 10) n++
  return n
}

export const SPELLCASTING_CLASSES = ['Bardo', 'Bruxo', 'Clérigo', 'Druida', 'Feiticeiro', 'Guardião', 'Mago', 'Paladino']

// Magias/truques conhecidos no nível 1 — usado na tela de criação, não é uma
// fórmula geral por nível (Bardo/Bruxo/Feiticeiro têm listas fixas por nível
// mais adiante que não seguem "nível + atributo"; isso fica para uma fase
// futura de level-up de magias).
export function spellsKnownAtLevel1(className: string, chaOrWisOrIntMod: number): number {
  switch (className) {
    case 'Bardo': return 4
    case 'Bruxo': return 2
    case 'Feiticeiro': return 2
    case 'Clérigo':
    case 'Druida':
    case 'Paladino':
    case 'Guardião':
      // conjurador "preparado": nível (=1) + modificador do atributo-chave, mínimo 1
      return Math.max(1, 1 + chaOrWisOrIntMod)
    case 'Mago':
      // Mago não "prepara" a partir de fórmula pura no nível 1 nesta versão —
      // grimório começa com 6 magias de 1º círculo (Conjuração de Mago), e o
      // número de magias preparadas também é nível + INT — mesma fórmula.
      return Math.max(1, 1 + chaOrWisOrIntMod)
    default:
      return 0
  }
}

// Atributo-chave de conjuração por classe (usado para localizar o modificador
// certo do personagem ao calcular spellsKnownAtLevel1/CD/bônus de ataque).
export const SPELLCASTING_ABILITY: Record<string, 'intelligence' | 'wisdom' | 'charisma'> = {
  Mago: 'intelligence',
  Clérigo: 'wisdom', Druida: 'wisdom', Guardião: 'wisdom',
  Bardo: 'charisma', Bruxo: 'charisma', Feiticeiro: 'charisma', Paladino: 'charisma',
}
