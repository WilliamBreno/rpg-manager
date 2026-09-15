// Regras de acesso a Rituais (Conjuração Ritual, D&D 4e) — espelha
// ritualAccessRules4e em backend/internal/service/ritual_service.go. Duplicado
// aqui pelo mesmo motivo já documentado em spellSlots.ts: front-end e
// back-end não compartilham runtime. Se as regras mudarem, corrija os dois.
//
// Usado só na CRIAÇÃO de personagem (sempre nível 1) — por isso só o valor
// de nível 1 de cada classe importa aqui; o crescimento por nível (ex: Mago
// +2 em 5/11/15/21/25) só é relevante depois, via GET
// /characters/:id/rituals/access (usado na ficha, não na criação).
export interface RitualAccessRule {
  level1Slots: number
  fixedRituals: string[]
  requiredPrerequisiteClass?: string
  restrictedOptions?: string[]
}

export const RITUAL_ACCESS_RULES_4E: Record<string, RitualAccessRule> = {
  Mago: { level1Slots: 3, fixedRituals: [] },
  Clérigo: { level1Slots: 2, fixedRituals: ['Repouso Tranquilo'] },
  Bardo: { level1Slots: 2, fixedRituals: [], requiredPrerequisiteClass: 'Bardo' },
  Druida: { level1Slots: 2, fixedRituals: ['Mensageiro Animal'] },
  Invocador: { level1Slots: 2, fixedRituals: ['Mão do Destino'] },
  Psionista: { level1Slots: 2, fixedRituals: [], restrictedOptions: ['Disco Flutuante de Tenser', 'Enviar Mensagem'] },
}

export function isRitualCastingClass(className?: string | null): boolean {
  return !!className && className in RITUAL_ACCESS_RULES_4E
}

export function ritualRuleFor(className?: string | null): RitualAccessRule | null {
  if (!className) return null
  return RITUAL_ACCESS_RULES_4E[className] ?? null
}
