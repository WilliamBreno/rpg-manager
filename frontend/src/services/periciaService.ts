import api from './api'
import type { Pericia, CharacterPericia } from '../types'

export const periciaService = {
  getAll: async (edition?: string): Promise<Pericia[]> => {
    const { data } = await api.get('/pericias', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByCharacter: async (characterId: number): Promise<CharacterPericia[]> => {
    const { data } = await api.get(`/characters/${characterId}/pericias`)
    return data
  },

  save: async (characterId: number, pericias: string[]): Promise<void> => {
    await api.post(`/characters/${characterId}/pericias`, { pericias })
  },

  setExpertise: async (characterId: number, periciaName: string, expertise: boolean): Promise<void> => {
    await api.patch(`/characters/${characterId}/pericias/expertise`, { pericia_name: periciaName, expertise })
  },
}