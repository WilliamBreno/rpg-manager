export const XP_TABLE_4E: Record<number, number> = {
  1: 0, 2: 1000, 3: 2250, 4: 3750, 5: 5500,
  6: 7500, 7: 10000, 8: 13000, 9: 16500, 10: 20500,
  11: 26000, 12: 32000, 13: 39000, 14: 47000, 15: 57000,
  16: 69000, 17: 83000, 18: 99000, 19: 119000, 20: 143000,
  21: 175000, 22: 210000, 23: 255000, 24: 310000, 25: 375000,
  26: 450000, 27: 550000, 28: 675000, 29: 825000, 30: 1000000,
}

export const XP_TABLE_5E: Record<number, number> = {
  1: 0, 2: 300, 3: 900, 4: 2700, 5: 6500,
  6: 14000, 7: 23000, 8: 34000, 9: 48000, 10: 64000,
  11: 85000, 12: 100000, 13: 120000, 14: 140000, 15: 165000,
  16: 195000, 17: 225000, 18: 265000, 19: 305000, 20: 355000,
}

export function maxLevelFor(edition: string): number {
  return edition === '4e' ? 30 : 20
}

export interface XPProgress {
  currentXP: number
  currentLevelXP: number
  nextLevelXP: number
  xpNeeded: number
  progressPercent: number
  isMaxLevel: boolean
}

export function xpProgressFor(edition: string, level: number, experiencePoints: number): XPProgress {
  const xpTable = edition === '4e' ? XP_TABLE_4E : XP_TABLE_5E
  const isMaxLevel = level >= maxLevelFor(edition)
  const currentXP = experiencePoints ?? 0
  const currentLevelXP = xpTable[level] ?? 0
  const nextLevelXP = xpTable[level + 1] ?? 0
  const xpNeeded = isMaxLevel ? 0 : Math.max(0, nextLevelXP - currentXP)
  const progressPercent = isMaxLevel
    ? 100
    : nextLevelXP > currentLevelXP
    ? Math.min(100, Math.max(0, ((currentXP - currentLevelXP) / (nextLevelXP - currentLevelXP)) * 100))
    : 100

  return { currentXP, currentLevelXP, nextLevelXP, xpNeeded, progressPercent, isMaxLevel }
}
