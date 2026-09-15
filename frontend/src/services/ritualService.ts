import api from './api'
import type { Ritual, RitualAccessInfo } from '../types'

export const ritualService = {
  getAll: async (edition?: string): Promise<Ritual[]> => {
    const { data } = await api.get('/rituals', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByCharacter: async (characterId: number): Promise<Ritual[]> => {
    const { data } = await api.get(`/characters/${characterId}/rituals`)
    return data
  },

  getAccess: async (characterId: number): Promise<RitualAccessInfo> => {
    const { data } = await api.get(`/characters/${characterId}/rituals/access`)
    return data
  },

  add: async (characterId: number, ritualId: number): Promise<void> => {
    await api.post(`/characters/${characterId}/rituals/${ritualId}`)
  },

  remove: async (characterId: number, ritualId: number): Promise<void> => {
    await api.delete(`/characters/${characterId}/rituals/${ritualId}`)
  },
}
