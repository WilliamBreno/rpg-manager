import api from './api'
import type { Campaign } from '../types'

export const campaignService = {
  getAll: async (): Promise<Campaign[]> => {
    const { data } = await api.get('/campaigns')
    return data
  },

  getByID: async (id: number): Promise<Campaign> => {
    const { data } = await api.get(`/campaigns/${id}`)
    return data
  },

  create: async (payload: { name: string; main_story: string }): Promise<Campaign> => {
    const { data } = await api.post('/campaigns', payload)
    return data
  },

  update: async (id: number, payload: { name: string; main_story: string }): Promise<Campaign> => {
    const { data } = await api.put(`/campaigns/${id}`, payload)
    return data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/campaigns/${id}`)
  },
}
