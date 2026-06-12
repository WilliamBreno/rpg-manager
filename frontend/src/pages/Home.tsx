import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

export default function Home() {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center gap-8">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-white mb-2">🎲 RPG Manager</h1>
        <p className="text-gray-400 text-xl">Bem-vindo, {user?.name}!</p>
        <p className="text-gray-500 text-sm mt-1">
          {user?.role === 'master' ? '👑 Mestre' : '🧙 Jogador'}
        </p>
      </div>

      <div className="flex gap-4">
        <button
          onClick={() => navigate('/characters')}
          className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-8 py-3 rounded-lg transition"
        >
          Ver Personagens
        </button>
        <button
          onClick={() => navigate('/characters/new')}
          className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold px-8 py-3 rounded-lg transition"
        >
          Criar Personagem
        </button>
      </div>

      <button
        onClick={handleLogout}
        className="text-gray-500 hover:text-red-400 transition text-sm"
      >
        Sair da conta
      </button>
    </div>
  )
}