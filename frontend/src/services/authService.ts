import api from './api'

interface RegisterData {
  name: string
  email: string
  password: string
  role: 'player' | 'master'
}

interface LoginData {
  email: string
  password: string
}

export const authService = {
  register: async (data: RegisterData) => {
    const { data: response } = await api.post('/auth/register', data)
    return response
  },

  login: async (data: LoginData) => {
    const { data: response } = await api.post('/auth/login', data)
    return response
  },

  markWelcomeSeen: async (): Promise<void> => {
    await api.patch('/users/me/welcome-seen')
  },
}