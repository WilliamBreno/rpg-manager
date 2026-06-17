import api from './api'
import type { Pericia } from '../types'

export const periciaService = {
  getAll: async (edition?: string): Promise<Pericia[]> => {
    const { data } = await api.get('/pericias', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByCharacter: async (characterId: number): Promise<{ pericia_name: string }[]> => {
    const { data } = await api.get(`/characters/${characterId}/pericias`)
    return data
  },

  save: async (characterId: number, pericias: string[]): Promise<void> => {
    await api.post(`/characters/${characterId}/pericias`, { pericias })
  },
}