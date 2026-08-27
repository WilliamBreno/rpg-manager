import api from './api'
import type { Spell } from '../types'

export const spellService = {
  getAll: async (edition?: string): Promise<Spell[]> => {
    const { data } = await api.get('/spells', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByCharacter: async (characterId: number): Promise<Spell[]> => {
    const { data } = await api.get(`/characters/${characterId}/spells`)
    return data
  },

  add: async (characterId: number, spellId: number): Promise<void> => {
    await api.post(`/characters/${characterId}/spells/${spellId}`)
  },

  remove: async (characterId: number, spellId: number): Promise<void> => {
    await api.delete(`/characters/${characterId}/spells/${spellId}`)
  },
}
