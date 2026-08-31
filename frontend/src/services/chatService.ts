import api from './api'

export interface ChatMessage {
  ID: number
  campaign_id: number
  session_id: number | null
  sender_user_id: number
  sender?: { ID: number; name: string }
  text: string
  CreatedAt: string
}

export const chatService = {
  getHistory: async (campaignId: number): Promise<ChatMessage[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/chat`)
    return data
  },

  send: async (campaignId: number, text: string, sessionId?: number): Promise<ChatMessage> => {
    const { data } = await api.post(`/campaigns/${campaignId}/chat`, { text, session_id: sessionId })
    return data
  },
}
