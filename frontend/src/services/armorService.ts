import api from './api'
import type { Armor } from '../types'

export const armorService = {
  getByEdition: async (edition: string): Promise<Armor[]> => {
    const { data } = await api.get('/armors', { params: { edition } })
    return data
  },
}