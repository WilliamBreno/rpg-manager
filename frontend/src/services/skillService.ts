import api from './api'
import type { Skill } from '../types'

export type { Skill }

export const skillService = {
  getAll: async (): Promise<Skill[]> => {
    const { data } = await api.get('/skills')
    return data
  },

  getByFilter: async (classId: number, raceId?: number): Promise<Skill[]> => {
    const params = new URLSearchParams()
    params.append('class_id', String(classId))
    if (raceId) params.append('race_id', String(raceId))
    const { data } = await api.get(`/skills/filter?${params.toString()}`)
    return data
  },

  getByClassAndRace: async (classId: number, raceId: number): Promise<Skill[]> => {
    const { data } = await api.get(`/skills/filter?class_id=${classId}&race_id=${raceId}`)
    return data
  },
}