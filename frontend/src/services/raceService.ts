import api from './api'
import type { Race } from '../types'

export const raceService = {
  getAll: async (edition?: string): Promise<Race[]> => {
    const { data } = await api.get('/races', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByID: async (id: number): Promise<Race> => {
    const { data } = await api.get(`/races/${id}`)
    return data
  },

  create: async (race: Omit<Race, 'ID'>): Promise<Race> => {
    const { data } = await api.post('/races', race)
    return data
  },
}