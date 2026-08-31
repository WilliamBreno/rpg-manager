import api from './api'
import type { Scene, SceneToken } from '../types'

export interface SceneInput {
  name: string
  image_url: string
}

export interface TokenInput {
  label: string
  image_url: string
  x: number
  y: number
}

export const sceneService = {
  getByCampaign: async (campaignId: number): Promise<Scene[]> => {
    const { data } = await api.get(`/campaigns/${campaignId}/scenes`)
    return data
  },

  getByID: async (id: number): Promise<Scene> => {
    const { data } = await api.get(`/scenes/${id}`)
    return data
  },

  create: async (campaignId: number, payload: SceneInput): Promise<Scene> => {
    const { data } = await api.post(`/campaigns/${campaignId}/scenes`, payload)
    return data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/scenes/${id}`)
  },

  addToken: async (sceneId: number, payload: TokenInput): Promise<SceneToken> => {
    const { data } = await api.post(`/scenes/${sceneId}/tokens`, payload)
    return data
  },

  moveToken: async (tokenId: number, x: number, y: number): Promise<SceneToken> => {
    const { data } = await api.patch(`/tokens/${tokenId}/move`, { x, y })
    return data
  },

  deleteToken: async (tokenId: number): Promise<void> => {
    await api.delete(`/tokens/${tokenId}`)
  },

  setActiveScene: async (sessionId: number, sceneId: number) => {
    const { data } = await api.patch(`/sessions/${sessionId}/active-scene`, { scene_id: sceneId })
    return data
  },
}
