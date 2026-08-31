import api from './api'
import type { NPC } from '../types'

export type NPCInput = Omit<NPC, 'ID' | 'campaign_id'>

export const npcService = {
  getByCampaign: async (campaignId: number): Promise<NPC[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/npcs`)
    return data
  },

  create: async (campaignId: number, payload: NPCInput): Promise<NPC> => {
    const { data } = await api.post(`/campaigns/${campaignId}/npcs`, payload)
    return data
  },

  update: async (id: number, payload: NPCInput): Promise<NPC> => {
    const { data } = await api.put(`/npcs/${id}`, payload)
    return data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/npcs/${id}`)
  },
}
