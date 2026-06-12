import api from './api'
import type { Background } from '../types'

export const backgroundService = {
  getByCharacterID: async (characterID: number): Promise<Background> => {
    const { data } = await api.get(`/characters/${characterID}/background`)
    return data
  },

  save: async (characterID: number, background: Background): Promise<Background> => {
    const { data } = await api.post(`/characters/${characterID}/background`, background)
    return data
  },
}