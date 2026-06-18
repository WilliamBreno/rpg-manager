import api from './api'

// backgroundService gerencia a biografia/notas do personagem
// (rota: GET/POST /characters/:id/background)
export interface BackgroundBio {
  personality_traits?: string
  ideals?:             string
  bonds?:              string
  flaws?:              string
}

export const backgroundService = {
  get: async (characterID: number): Promise<BackgroundBio> => {
    const { data } = await api.get<BackgroundBio>(`/characters/${characterID}/background`)
    return data
  },
  save: async (characterID: number, data: BackgroundBio): Promise<void> => {
    await api.post(`/characters/${characterID}/background`, data)
  },
}