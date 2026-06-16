import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useParams } from 'react-router-dom'
import { characterService } from '../services/characterService'
import BackgroundForm from '../components/BackgroundForm'
import AvatarUpload from '../components/AvatarUpload'
import HPManager from '../components/HPManager'
import SkillsPanel from '../components/SkillsPanel'

export default function CharacterDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: character, isLoading } = useQuery({
    queryKey: ['character', id],
    queryFn: () => characterService.getByID(Number(id)),
  })

  const { data: acData } = useQuery({
    queryKey: ['character-ac', id],
    queryFn: () => characterService.getAC(Number(id)),
  })

  const levelUpMutation = useMutation({
    mutationFn: () => characterService.levelUp(Number(id)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character', id] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: () => characterService.delete(Number(id)),
    onSuccess: () => navigate('/characters'),
  })

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <p className="text-gray-400">Carregando...</p>
      </div>
    )
  }

  if (!character) return null

  const attributes = [
    { label: 'FOR', value: character.strength },
    { label: 'DES', value: character.dexterity },
    { label: 'CON', value: character.constitution },
    { label: 'INT', value: character.intelligence },
    { label: 'SAB', value: character.wisdom },
    { label: 'CAR', value: character.charisma },
  ]


  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-3xl mx-auto">

        <button
          onClick={() => navigate('/characters')}
          className="text-gray-400 hover:text-white transition mb-6 block"
        >
          ← Voltar
        </button>

        {/* Cabeçalho */}
        <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
          <div className="flex justify-between items-start">
            <div className="flex gap-5 items-center">
              <AvatarUpload
                characterID={Number(id)}
                avatarURL={character.avatar_url}
                characterName={character.name}
              />
              <div>
                <h1 className="text-3xl font-bold text-white">{character.name}</h1>
                <p className="text-gray-400 mt-1">Edição: {character.edition}</p>
                <div className="flex gap-3 mt-3">
                  <span className="bg-indigo-900 text-indigo-300 px-3 py-1 rounded-full text-sm">
                    {character.class?.name ?? 'Sem classe'}
                  </span>
                  <span className="bg-emerald-900 text-emerald-300 px-3 py-1 rounded-full text-sm">
                    {character.race?.name ?? 'Sem raça'}
                  </span>
                </div>
              </div>
            </div>
            <div className="text-center">
              <p className="text-gray-400 text-sm">Nível</p>
              <p className="text-5xl font-bold text-yellow-400">{character.level}</p>
              <p className="text-gray-400 text-sm mt-1">HP: {character.hit_points}</p>
              {acData && (
                <p className="text-indigo-300 text-sm mt-1">CA: {acData.ac}</p>
              )}
            </div>
          </div>
        </div>
        
        {/* HP Manager */}
        <HPManager character={character} />    

        {/* Atributos */}
        <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
          <h2 className="text-lg font-semibold text-white mb-4">Atributos</h2>
          <div className="grid grid-cols-6 gap-3">
            {attributes.map(attr => (
              <div key={attr.label} className="text-center bg-gray-700 rounded-lg p-3">
                <p className="text-gray-400 text-xs mb-1">{attr.label}</p>
                <p className="text-white font-bold text-xl">{attr.value}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Habilidades da Classe via IA */}
        <SkillsPanel
            skills={character.skills}
            edition={character.edition}
        />

        {/* Background */}
        <BackgroundForm
          characterID={Number(id)}
          background={character.background}
        />
        {/* Ações */}
        <div className="flex gap-4">
          <button
            onClick={() => navigate(`/characters/${id}/edit`)}
            className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-2 rounded-lg transition"
          >
            Editar
          </button>
          <button
            onClick={() => levelUpMutation.mutate()}
            disabled={levelUpMutation.isPending || character.level >= 20}
            className="bg-yellow-600 hover:bg-yellow-700 disabled:opacity-50 text-white font-semibold px-6 py-2 rounded-lg transition"
          >
            {levelUpMutation.isPending ? 'Subindo...' : '⬆ Level Up'}
          </button>
          <button
            onClick={() => deleteMutation.mutate()}
            disabled={deleteMutation.isPending}
            className="bg-red-700 hover:bg-red-800 disabled:opacity-50 text-white font-semibold px-6 py-2 rounded-lg transition"
          >
            {deleteMutation.isPending ? 'Deletando...' : 'Deletar'}
          </button>
        </div>

      </div>
    </div>
  )
}