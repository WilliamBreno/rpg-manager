import api from './api'
import type { Item, ItemCategory } from '../types'

export const itemService = {
  getAll: async (edition: string, category?: ItemCategory): Promise<Item[]> => {
    const { data } = await api.get('/items', { params: { edition, category } })
    return data
  },
}
