import api from './api'
import type { Antecedent } from '../types'

export const antecedentService = {
  getAll: async (edition: string): Promise<Antecedent[]> => {
    const { data } = await api.get<Antecedent[]>(`/antecedentes?edition=${edition}`)
    return data
  },
  getById: async (id: number): Promise<Antecedent> => {
    const { data } = await api.get<Antecedent>(`/antecedentes/${id}`)
    return data
  },
}