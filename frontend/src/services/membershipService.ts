import api from './api'

export interface CampaignMembership {
  ID: number
  campaign_id: number
  campaign?: { ID: number; name: string }
  user_id: number
  user?: { ID: number; name: string; email: string }
  character_id: number | null
  character?: { ID: number; name: string }
  status: 'invited' | 'accepted' | 'declined'
}

export const membershipService = {
  invite: async (campaignId: number, email: string): Promise<CampaignMembership> => {
    const { data } = await api.post(`/campaigns/${campaignId}/invites`, { email })
    return data
  },

  getMembers: async (campaignId: number): Promise<CampaignMembership[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/members`)
    return data
  },

  getMyPending: async (): Promise<CampaignMembership[]> => {
    const { data } = await api.get('/users/me/campaign-invites')
    return data
  },

  getMyCampaigns: async (): Promise<CampaignMembership[]> => {
    const { data } = await api.get('/users/me/campaigns')
    return data
  },

  respond: async (membershipId: number, accept: boolean, characterId?: number): Promise<CampaignMembership> => {
    const { data } = await api.patch(`/campaign-memberships/${membershipId}/respond`, { accept, character_id: characterId })
    return data
  },
}
