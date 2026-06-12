import api from './api'
import type { Skill } from '../types'

export const skillService = {
  getAll: async (): Promise<Skill[]> => {
    const { data } = await api.get('/skills')
    return data
  },

  getByClassAndRace: async (classID: number, raceID: number): Promise<Skill[]> => {
    const { data } = await api.get('/skills/filter', {
      params: { class_id: classID, race_id: raceID },
    })
    return data
  },

  create: async (skill: Omit<Skill, 'ID'>): Promise<Skill> => {
    const { data } = await api.post('/skills', skill)
    return data
  },
}