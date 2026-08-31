import api from './api'

export interface MagicItem {
  ID: number
  campaign_id: number
  name: string
  description: string
  effect: string
}

export interface Reward {
  ID: number
  CreatedAt: string
  campaign_id: number
  character_id: number
  character?: { ID: number; name: string }
  granted_by_user_id: number
  granted_by?: { ID: number; name: string }
  kind: 'currency' | 'item'
  copper_pieces: number
  silver_pieces: number
  electrum_pieces: number
  gold_pieces: number
  platinum_pieces: number
  magic_item_id: number | null
  magic_item?: MagicItem
  note: string
}

export interface GrantCurrencyInput {
  character_id?: number
  all?: boolean
  copper_pieces?: number
  silver_pieces?: number
  electrum_pieces?: number
  gold_pieces?: number
  platinum_pieces?: number
  note?: string
}

export interface GrantItemInput {
  character_id?: number
  all?: boolean
  magic_item_id: number
  note?: string
}

export const rewardService = {
  createMagicItem: async (campaignId: number, payload: { name: string; description: string; effect: string }): Promise<MagicItem> => {
    const { data } = await api.post(`/campaigns/${campaignId}/magic-items`, payload)
    return data
  },

  getMagicItems: async (campaignId: number): Promise<MagicItem[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/magic-items`)
    return data
  },

  deleteMagicItem: async (id: number): Promise<void> => {
    await api.delete(`/magic-items/${id}`)
  },

  getHistory: async (campaignId: number): Promise<Reward[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/rewards`)
    return data
  },

  grantCurrency: async (campaignId: number, payload: GrantCurrencyInput): Promise<Reward[]> => {
    const { data } = await api.post(`/campaigns/${campaignId}/rewards/currency`, payload)
    return data
  },

  grantItem: async (campaignId: number, payload: GrantItemInput): Promise<Reward[]> => {
    const { data } = await api.post(`/campaigns/${campaignId}/rewards/item`, payload)
    return data
  },
}
