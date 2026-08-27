import api from './api'
import type { ClassEquipmentOption } from '../types'

export const classEquipmentService = {
  getByClass: async (classId: number): Promise<ClassEquipmentOption[]> => {
    const { data } = await api.get(`/classes/${classId}/equipment-options`)
    return data
  },
}
