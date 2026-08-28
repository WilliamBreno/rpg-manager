import api from './api'
import type { Language } from '../types'

export const languageService = {
  getAll: async (edition?: string): Promise<Language[]> => {
    const { data } = await api.get('/languages', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByCharacter: async (characterId: number): Promise<Language[]> => {
    const { data } = await api.get(`/characters/${characterId}/languages`)
    return data
  },

  add: async (characterId: number, languageId: number): Promise<void> => {
    await api.post(`/characters/${characterId}/languages/${languageId}`)
  },

  remove: async (characterId: number, languageId: number): Promise<void> => {
    await api.delete(`/characters/${characterId}/languages/${languageId}`)
  },
}
