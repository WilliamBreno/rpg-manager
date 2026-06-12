import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { characterService } from '../services/characterService'

export default function CharacterList() {
  const navigate = useNavigate()

  const { data: characters, isLoading, error } = useQuery({
    queryKey: ['characters'],
    queryFn: characterService.getAll,
  })

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white">Personagens</h1>
          <button
            onClick={() => navigate('/characters/new')}
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold px-6 py-2 rounded-lg transition"
          >
            + Novo Personagem
          </button>
        </div>

        {isLoading && (
          <p className="text-gray-400 text-center">Carregando...</p>
        )}

        {error && (
          <p className="text-red-400 text-center">Erro ao carregar personagens.</p>
        )}

        {characters && characters.length === 0 && (
          <p className="text-gray-400 text-center">Nenhum personagem criado ainda.</p>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {characters?.map(character => (
            <div
              key={character.ID}
              onClick={() => navigate(`/characters/${character.ID}`)}
              className="bg-gray-800 rounded-lg p-6 cursor-pointer hover:bg-gray-700 transition border border-gray-700 flex gap-4 items-center"
            >
              <div className="w-14 h-14 rounded-full overflow-hidden bg-gray-700 flex items-center justify-center flex-shrink-0">
                {character.avatar_url ? (
                  <img
                    src={`http://localhost:8080${character.avatar_url}`}
                    alt={character.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <span className="text-2xl">🧙</span>
                )}
              </div>
              <div>
                <h2 className="text-xl font-bold text-white mb-1">{character.name}</h2>
                <p className="text-gray-400 text-sm mb-3">Edição: {character.edition}</p>
                <div className="flex gap-4 text-sm">
                  <span className="bg-indigo-900 text-indigo-300 px-3 py-1 rounded-full">
                    {character.class?.name ?? 'Sem classe'}
                  </span>
                  <span className="bg-emerald-900 text-emerald-300 px-3 py-1 rounded-full">
                    {character.race?.name ?? 'Sem raça'}
                  </span>
                  <span className="bg-yellow-900 text-yellow-300 px-3 py-1 rounded-full">
                    Nível {character.level}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>

        <button
          onClick={() => navigate('/')}
          className="mt-8 text-gray-400 hover:text-white transition"
        >
          ← Voltar
        </button>
      </div>
    </div>
  )
}