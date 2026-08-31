import api from './api'
import type { CRDamageStats, Enemy, EnemyAbility, EnemyKind, EnemyLine, EnemyLineSource } from '../types'

export interface EnemyAbilityInput {
  name: string
  damage: string
  description: string
}

export interface EnemyLineInput {
  text: string
  audio_url: string
  source: EnemyLineSource
}

export interface EnemyInput {
  kind: EnemyKind
  name: string
  hp: number
  challenge_rating: string
  race: string
  photo_url: string
  sound_url: string
  class: string
  armor: number
  history: string
  bonds: string
  notes: string
  abilities: EnemyAbilityInput[]
  lines: EnemyLineInput[]
}

export const enemyService = {
  getByCampaign: async (campaignId: number): Promise<Enemy[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/enemies`)
    return data
  },

  create: async (campaignId: number, payload: EnemyInput): Promise<{ enemy: Enemy; warnings: string[] | null }> => {
    const { data } = await api.post(`/campaigns/${campaignId}/enemies`, payload)
    return data
  },

  update: async (id: number, payload: EnemyInput): Promise<Enemy> => {
    const { data } = await api.put(`/enemies/${id}`, payload)
    return data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/enemies/${id}`)
  },

  addAbility: async (enemyId: number, payload: EnemyAbilityInput): Promise<{ ability: EnemyAbility; warning: string }> => {
    const { data } = await api.post(`/enemies/${enemyId}/abilities`, payload)
    return data
  },

  deleteAbility: async (id: number): Promise<void> => {
    await api.delete(`/enemy-abilities/${id}`)
  },

  addLine: async (enemyId: number, payload: EnemyLineInput): Promise<EnemyLine> => {
    const { data } = await api.post(`/enemies/${enemyId}/lines`, payload)
    return data
  },

  deleteLine: async (id: number): Promise<void> => {
    await api.delete(`/enemy-lines/${id}`)
  },

  getCRDamageTable: async (): Promise<CRDamageStats[]> => {
    const { data } = await api.get('/dnd/cr-damage-table')
    return data
  },

  playSound: async (enemyId: number): Promise<void> => {
    await api.post(`/enemies/${enemyId}/play-sound`)
  },

  playLine: async (lineId: number): Promise<void> => {
    await api.post(`/enemy-lines/${lineId}/play`)
  },
}
