import api from './api'
import type { Class } from '../types'

export const classService = {
  getAll: async (edition?: string): Promise<Class[]> => {
    const { data } = await api.get('/classes', {
      params: edition ? { edition } : {},
    })
    return data
  },

  getByID: async (id: number): Promise<Class> => {
    const { data } = await api.get(`/classes/${id}`)
    return data
  },

  create: async (cls: Omit<Class, 'ID'>): Promise<Class> => {
    const { data } = await api.post('/classes', cls)
    return data
  },
}