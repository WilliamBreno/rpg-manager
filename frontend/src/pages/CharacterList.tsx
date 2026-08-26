import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { useMemo, useState } from 'react'
import { characterService } from '../services/characterService'
import { authService } from '../services/authService'
import { useAuthStore } from '../store/authStore'
import WelcomeModal from '../components/WelcomeModal'
import CharacterCard from '../components/CharacterCard'

const NAV_ITEMS = [
  { label: 'Início', icon: '🏠', path: '/' },
  { label: 'Meus Personagens', icon: '📜', path: '/characters' },
]

export default function CharacterList() {
  const navigate = useNavigate()
  const { user, logout, markWelcomeSeen } = useAuthStore()
  const [showWelcome, setShowWelcome] = useState(() => user ? !user.welcome_seen : false)
  const [search, setSearch] = useState('')

  const { data: characters, isLoading, error } = useQuery({
    queryKey: ['characters'],
    queryFn: characterService.getAll,
  })

  const filteredCharacters = useMemo(() => {
    if (!characters) return characters
    const q = search.trim().toLowerCase()
    if (!q) return characters
    return characters.filter(c => c.name.toLowerCase().includes(q))
  }, [characters, search])

  const dismissWelcome = () => {
    setShowWelcome(false)
    markWelcomeSeen()
    authService.markWelcomeSeen().catch(() => {
      // Falha silenciosa: o pior caso é o modal aparecer de novo no próximo
      // login, não é uma operação crítica o suficiente pra bloquear o usuário.
    })
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-gray-900 flex">
      {showWelcome && <WelcomeModal onClose={dismissWelcome} />}

      {/* ── Sidebar ──────────────────────────────────────────────────────── */}
      <aside className="hidden md:flex w-60 flex-shrink-0 flex-col bg-gray-950 border-r border-gray-800 py-6 px-4">
        <div className="flex items-center gap-2 px-2 mb-8">
          <span className="text-2xl">🎲</span>
          <span className="font-rpg font-bold text-lg" style={{ color: '#c9a84c' }}>RPG Manager</span>
        </div>

        <nav className="flex-1 flex flex-col gap-1">
          {NAV_ITEMS.map(item => {
            const active = item.path === '/characters'
            return (
              <button
                key={item.path}
                onClick={() => navigate(item.path)}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition text-left ${
                  active ? 'bg-rpg-gold-muted text-rpg-gold' : 'text-gray-400 hover:bg-gray-800 hover:text-gray-200'
                }`}
              >
                <span>{item.icon}</span>
                <span>{item.label}</span>
              </button>
            )
          })}
        </nav>

        <div className="pt-4 border-t border-gray-800 flex items-center gap-3 px-2">
          <div className="w-9 h-9 rounded-full bg-gray-700 flex items-center justify-center text-sm font-bold text-gray-200 flex-shrink-0">
            {user?.name?.charAt(0).toUpperCase() ?? '?'}
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-sm text-gray-200 truncate">{user?.name}</p>
            <button onClick={handleLogout} className="text-xs text-gray-500 hover:text-red-400 transition">
              Sair da conta
            </button>
          </div>
        </div>
      </aside>

      {/* ── Conteúdo principal ───────────────────────────────────────────── */}
      <div className="flex-1 px-4 py-6 sm:px-8 sm:py-8">
        <div className="max-w-5xl mx-auto">

          <div className="flex flex-col lg:flex-row gap-6">
            <div className="flex-1 min-w-0">

              {/* Header */}
              <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-3 mb-6">
                <div>
                  <h1 className="text-2xl sm:text-3xl font-bold text-white">Meus Personagens</h1>
                  <p className="text-gray-500 text-sm mt-1">Gerencie todos os seus heróis e suas histórias.</p>
                </div>
                <button
                  onClick={() => navigate('/characters/new')}
                  className="w-full sm:w-auto btn-rpg-primary px-5 py-2.5 whitespace-nowrap"
                >
                  + Criar Personagem
                </button>
              </div>

              {characters && characters.length > 0 && (
                <input
                  type="text"
                  value={search}
                  onChange={e => setSearch(e.target.value)}
                  placeholder="Buscar personagem pelo nome..."
                  className="rpg-input mb-6"
                />
              )}

              {isLoading && <p className="text-gray-400 text-center py-8">Carregando...</p>}
              {error   && <p className="text-red-400 text-center py-8">Erro ao carregar personagens.</p>}

              {characters?.length === 0 && (
                <div className="text-center py-16">
                  <p className="text-4xl mb-4">🎲</p>
                  <p className="text-gray-400">Nenhum personagem criado ainda.</p>
                  <button
                    onClick={() => navigate('/characters/new')}
                    className="mt-4 btn-rpg-primary px-6 py-2"
                  >
                    Criar meu primeiro personagem
                  </button>
                </div>
              )}

              {characters && characters.length > 0 && filteredCharacters?.length === 0 && (
                <p className="text-gray-500 text-center py-8">Nenhum personagem encontrado para "{search}".</p>
              )}

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {filteredCharacters?.map(character => (
                  <CharacterCard key={character.ID} character={character} />
                ))}
              </div>
            </div>

            {/* ── Painel lateral: resumo ─────────────────────────────────── */}
            <aside className="lg:w-56 flex-shrink-0">
              <div className="rpg-card p-4">
                <p className="text-gray-500 text-xs uppercase tracking-widest mb-3">Resumo</p>
                <div className="flex items-center justify-between">
                  <span className="text-gray-400 text-sm">Personagens</span>
                  <span className="text-2xl font-bold text-white">{characters?.length ?? 0}</span>
                </div>
              </div>
            </aside>
          </div>

        </div>
      </div>
    </div>
  )
}
