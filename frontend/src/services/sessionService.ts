import api from './api'

export interface GameSession {
  ID: number
  campaign_id: number
  started_at: string
  ended_at: string | null
  summary: string
  music_url?: string
  active_scene_id: number | null
}

export const sessionService = {
  getByCampaign: async (campaignId: number): Promise<GameSession[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/sessions`)
    return data
  },

  start: async (campaignId: number): Promise<GameSession> => {
    const { data } = await api.post(`/campaigns/${campaignId}/sessions`)
    return data
  },

  updateSummary: async (id: number, summary: string): Promise<GameSession> => {
    const { data } = await api.patch(`/sessions/${id}/summary`, { summary })
    return data
  },

  end: async (id: number, summary?: string): Promise<GameSession> => {
    const { data } = await api.patch(`/sessions/${id}/end`, { summary })
    return data
  },

  setMusic: async (id: number, musicUrl: string, playing: boolean): Promise<GameSession> => {
    const { data } = await api.patch(`/sessions/${id}/music`, { music_url: musicUrl, playing })
    return data
  },
}
