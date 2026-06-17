import api from './api'
import type { Talento } from '../types'

export const talentoService = {
  getAll: async (edition?: string): Promise<Talento[]> => {
    const { data } = await api.get('/talentos', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByCharacter: async (characterId: number): Promise<Talento[]> => {
    const { data } = await api.get(`/characters/${characterId}/talentos`)
    return data
  },

  add: async (characterId: number, talentoId: number): Promise<void> => {
    await api.post(`/characters/${characterId}/talentos/${talentoId}`)
  },

  remove: async (characterId: number, talentoId: number): Promise<void> => {
    await api.delete(`/characters/${characterId}/talentos/${talentoId}`)
  },
}