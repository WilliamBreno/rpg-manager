import api from './api'
import type { Character, CreateCharacterDTO, Currency, CharacterItem, CharacterArmorOwned } from '../types'

export const characterService = {
  getAll: async (): Promise<Character[]> => {
    const { data } = await api.get('/characters')
    return data
  },

  getByID: async (id: number): Promise<Character> => {
    const { data } = await api.get(`/characters/${id}`)
    return data
  },

  create: async (character: CreateCharacterDTO): Promise<Character> => {
    const { data } = await api.post('/characters', character)
    return data
  },

  update: async (id: number, character: Partial<CreateCharacterDTO>): Promise<Character> => {
    const { data } = await api.put(`/characters/${id}`, character)
    return data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/characters/${id}`)
  },

  levelUp: async (id: number): Promise<Character> => {
    const { data } = await api.patch(`/characters/${id}/level-up`)
    return data
  },

  addSkill: async (characterID: number, skillID: number): Promise<void> => {
    await api.post(`/characters/${characterID}/skills/${skillID}`)
  },

  removeSkill: async (characterID: number, skillID: number): Promise<void> => {
    await api.delete(`/characters/${characterID}/skills/${skillID}`)
  },
  
  uploadAvatar: async (characterID: number, file: File): Promise<{ avatar_url: string }> => {
    const formData = new FormData()
    formData.append('avatar', file)
    const { data } = await api.post(`/characters/${characterID}/avatar`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return data
  },
  getAC: async (characterID: number): Promise<{ ac: number }> => {
    const { data } = await api.get(`/characters/${characterID}/ac`)
    return data
  },

  takeDamage: async (characterID: number, damage: number): Promise<Character> => {
    const { data } = await api.patch(`/characters/${characterID}/take-damage`, { damage })
    return data
  },

  heal: async (characterID: number, amount: number): Promise<Character> => {
    const { data } = await api.patch(`/characters/${characterID}/heal`, { amount })
    return data
  },

  addTempHP: async (characterID: number, amount: number): Promise<Character> => {
    const { data } = await api.patch(`/characters/${characterID}/temp-hp`, { amount })
    return data
  },

  addXP: async (id: number, xp: number) => {
    const { data } = await api.patch(`/characters/${id}/add-xp`, { xp })
    return data
  },

  applyASI: async (id: number, choices: Record<string, number>) => {
    const { data } = await api.patch(`/characters/${id}/apply-asi`, choices)
    return data
  },
  deathSave: async (id: number, body: { success: boolean; critical: boolean }) => {
    const { data } = await api.patch(`/characters/${id}/death-save`, body)
    return data
  },
  resetDeathSaves: async (id: number) => {
    const { data } = await api.patch(`/characters/${id}/reset-death-saves`, {})
    return data
  },

  exportPdf: async (id: number): Promise<Blob> => {
    const { data } = await api.get(`/characters/${id}/export/pdf`, {
      responseType: 'blob',
    })
    return data
  },

  getInventory: async (id: number): Promise<{ items: CharacterItem[]; armors: CharacterArmorOwned[] }> => {
    const { data } = await api.get(`/characters/${id}/inventory`)
    return data
  },

  buyItem: async (characterId: number, itemId: number, quantity = 1): Promise<Character> => {
    const { data } = await api.post(`/characters/${characterId}/shop/items/${itemId}`, { quantity })
    return data
  },

  buyArmor: async (characterId: number, armorId: number, quantity = 1): Promise<Character> => {
    const { data } = await api.post(`/characters/${characterId}/shop/armors/${armorId}`, { quantity })
    return data
  },

  setCurrency: async (characterId: number, currency: Currency): Promise<Character> => {
    const { data } = await api.patch(`/characters/${characterId}/currency`, currency)
    return data
  },
}
