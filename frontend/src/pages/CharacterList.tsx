import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { characterService } from '../services/characterService'

const API_BASE = (import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api').replace(/\/api$/, '')

export default function CharacterList() {
  const navigate = useNavigate()

  const { data: characters, isLoading, error } = useQuery({
    queryKey: ['characters'],
    queryFn: characterService.getAll,
  })

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-4xl mx-auto">

        {/* Header responsivo */}
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mb-6 sm:mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold text-white">Personagens</h1>
          <button
            onClick={() => navigate('/characters/new')}
            className="w-full sm:w-auto bg-emerald-600 hover:bg-emerald-700 text-white font-semibold px-5 py-2.5 rounded-lg transition text-sm sm:text-base"
          >
            + Novo Personagem
          </button>
        </div>

        {isLoading && <p className="text-gray-400 text-center py-8">Carregando...</p>}
        {error   && <p className="text-red-400 text-center py-8">Erro ao carregar personagens.</p>}
        {characters?.length === 0 && (
          <div className="text-center py-16">
            <p className="text-4xl mb-4">🎲</p>
            <p className="text-gray-400">Nenhum personagem criado ainda.</p>
            <button
              onClick={() => navigate('/characters/new')}
              className="mt-4 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold px-6 py-2 rounded-lg transition"
            >
              Criar meu primeiro personagem
            </button>
          </div>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
          {characters?.map(character => (
            <div
              key={character.ID}
              onClick={() => navigate(`/characters/${character.ID}`)}
              className="bg-gray-800 rounded-xl cursor-pointer hover:bg-gray-700 transition border border-gray-700 overflow-hidden relative"
            >
              {/* Badge de nível no topo direito */}
              <div className="absolute top-3 right-3 bg-yellow-600 text-white text-xs font-bold px-2 py-1 rounded-full">
                Nv. {character.level}
              </div>

              <div className="flex gap-3 items-center p-4">
                {/* Avatar */}
                <div className="w-14 h-14 rounded-full overflow-hidden bg-gray-700 flex items-center justify-center flex-shrink-0 border-2 border-gray-600">
                  {character.avatar_url ? (
                    <img
                      src={`${API_BASE}${character.avatar_url}`}
                      alt={character.name}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <span className="text-2xl">🧙</span>
                  )}
                </div>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <h2 className="text-base sm:text-lg font-bold text-white truncate pr-10">
                    {character.name}
                  </h2>
                  <p className="text-gray-500 text-xs mb-2">D&D {character.edition}</p>
                  <div className="flex flex-wrap gap-1.5">
                    <span className="bg-indigo-900 text-indigo-300 px-2 py-0.5 rounded-full text-xs">
                      {character.class?.name ?? 'Sem classe'}
                    </span>
                    <span className="bg-emerald-900 text-emerald-300 px-2 py-0.5 rounded-full text-xs">
                      {character.race?.name ?? 'Sem raça'}
                    </span>
                    <span className="bg-gray-700 text-gray-300 px-2 py-0.5 rounded-full text-xs">
                      ❤️ {character.hit_points ?? '?'} HP
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        <button
          onClick={() => navigate('/')}
          className="mt-8 text-gray-400 hover:text-white transition text-sm"
        >
          ← Voltar
        </button>

      </div>
    </div>
  )
}