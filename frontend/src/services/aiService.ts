import api from './api'

export interface AISkill {
  name: string
  power_type: 'at-will' | 'encounter' | 'daily' | 'utility'
  action_type: string
  range: string
  target: string
  attack: string
  hit: string
  miss: string
  effect: string
  min_level: number
}

export const aiService = {
  getSkills: async (className: string, edition: string, level: number): Promise<AISkill[]> => {
    const { data } = await api.post('/ai/skills', { class_name: className, edition, level })
    return data
  },
}