import api from './api'

export type PowerType = 'unlimited' | 'encounter' | 'daily' | 'utility'

export interface Skill {
  ID: number
  name: string
  description?: string
  power_type?: PowerType
  level?: number      // nível mínimo para desbloquear
  edition?: string
  class_id?: number
  race_id?: number
}

export const skillService = {
  getByFilter: async (classId: number, raceId?: number): Promise<Skill[]> => {
    const params = new URLSearchParams()
    params.append('class_id', String(classId))
    if (raceId) params.append('race_id', String(raceId))
    const { data } = await api.get(`/skills/filter?${params.toString()}`)
    return data
  },
}