// Parseia a Description de um Item categoria "arma" — já vem estruturada do
// seed no formato "<dado ou fixo> <Tipo de Dano> — <propriedades> (Maestria: X)",
// ex: "1d4 Perfurante — Acuidade, Arremesso (Alcance 6/18), Leve (Maestria: Ágil)".
// Não há campos estruturados separados pra isso no domain.Item (dano/propriedades
// vivem só nesse texto), então o front-end extrai o que precisa pra calcular o
// bônus de ataque/dano na pré-visualização do Equipamento Inicial.
export interface ParsedWeapon {
  dice: string        // "1d4" ou "1" (dano fixo)
  damageType: string  // "Perfurante", "Cortante", "Contundente"
  finesse: boolean     // tem "Acuidade" — pode usar DES em vez de FOR
  ranged: boolean      // tem "Munição" — usa DES por padrão (arco/besta/arma de fogo)
  mastery?: string
}

export function parseWeaponDescription(description: string): ParsedWeapon | null {
  if (!description) return null
  const [dmgPart, ...restParts] = description.split(' — ')
  const rest = restParts.join(' — ')
  const dmgMatch = dmgPart.trim().match(/^(\d+d\d+|\d+)\s+(.+)$/)
  if (!dmgMatch) return null
  const [, dice, damageType] = dmgMatch
  const masteryMatch = rest.match(/\(Maestria:\s*([^)]+)\)/)
  return {
    dice,
    damageType,
    finesse: /Acuidade/.test(rest),
    ranged: /Munição/.test(rest),
    mastery: masteryMatch?.[1],
  }
}

// Bônus de ataque no nível 1 (proficiência sempre +2 em 5e) assumindo que o
// personagem é proficiente com a arma — premissa segura aqui porque o
// Equipamento Inicial de uma classe só concede armas que ela já é proficiente.
export function attackBonusFor(weapon: ParsedWeapon, strMod: number, dexMod: number): number {
  const abilityMod = weapon.ranged ? dexMod : weapon.finesse ? Math.max(strMod, dexMod) : strMod
  return 2 + abilityMod
}

export function damageAbilityMod(weapon: ParsedWeapon, strMod: number, dexMod: number): number {
  return weapon.ranged ? dexMod : weapon.finesse ? Math.max(strMod, dexMod) : strMod
}

export function formatSigned(n: number): string {
  return n >= 0 ? `+${n}` : `${n}`
}

// CA de uma armadura para um dado modificador de Destreza (max_dex_bonus: -1
// = sem limite, 0 = sem bônus de DES, >0 = limite máximo).
export function acFor(baseAC: number, maxDexBonus: number, dexMod: number): number {
  if (maxDexBonus === 0) return baseAC
  const applied = maxDexBonus < 0 ? dexMod : Math.min(dexMod, maxDexBonus)
  return baseAC + applied
}
