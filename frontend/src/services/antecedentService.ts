import api from './api'
import type { Background } from '../types'

// antecedentService gerencia os Antecedentes D&D 5e
// (rota: GET /antecedentes?edition=5e)
export const antecedentService = {
  getAll: async (edition: string): Promise<Background[]> => {
    const { data } = await api.get<Background[]>('/antecedentes', { params: { edition } })
    return data
  },
  getById: async (id: number): Promise<Background> => {
    const { data } = await api.get<Background>(`/antecedentes/${id}`)
    return data
  },
}