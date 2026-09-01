import { create } from 'zustand'

interface User {
  id: number
  name: string
  email: string
  role: 'player' | 'master'
  welcome_seen: boolean
  master_welcome_seen: boolean
}

interface AuthStore {
  token: string | null
  user: User | null
  setAuth: (token: string, user: User) => void
  logout: () => void
  isAuthenticated: () => boolean
  markWelcomeSeen: () => void
  markMasterWelcomeSeen: () => void
}

export const useAuthStore = create<AuthStore>((set, get) => ({
  token: localStorage.getItem('token'),
  user: JSON.parse(localStorage.getItem('user') ?? 'null'),

  setAuth: (token, user) => {
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    set({ token, user })
  },

  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    set({ token: null, user: null })
  },

  isAuthenticated: () => !!get().token,

  // Atualização otimista local — o backend já foi marcado por quem chamou isso
  // (authService.markWelcomeSeen); aqui só sincroniza o estado em memória/cache.
  markWelcomeSeen: () => {
    const current = get().user
    if (!current) return
    const updated = { ...current, welcome_seen: true }
    localStorage.setItem('user', JSON.stringify(updated))
    set({ user: updated })
  },

  // Mesma ideia de markWelcomeSeen, pra tela de boas-vindas específica de Mestre.
  markMasterWelcomeSeen: () => {
    const current = get().user
    if (!current) return
    const updated = { ...current, master_welcome_seen: true }
    localStorage.setItem('user', JSON.stringify(updated))
    set({ user: updated })
  },
}))