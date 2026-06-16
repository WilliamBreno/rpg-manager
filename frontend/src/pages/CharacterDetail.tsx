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

  const levelUpMutation = useMutation({
    mutationFn: () => characterService.levelUp(Number(id)),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character', id] }),
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

  const is4e = character.edition === '4e'

  const attributes = [
    { label: 'FOR', value: character.strength },
    { label: 'DES', value: character.dexterity },
    { label: 'CON', value: character.constitution },
    { label: 'INT', value: character.intelligence },
    { label: 'SAB', value: character.wisdom },
    { label: 'CAR', value: character.charisma },
  ]

  const defenses = is4e ? [
    { label: 'CA',   value: character.defense_ac },
    { label: 'FORT', value: character.defense_fort },
    { label: 'REFL', value: character.defense_refl },
    { label: 'VONT', value: character.defense_will },
  ] : []

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-3xl mx-auto">

        <button
          onClick={() => navigate('/characters')}
          className="text-gray-400 hover:text-white transition mb-4 sm:mb-6 block text-sm"
        >
          ← Voltar
        </button>

        {/* Cabeçalho */}
        <div className="bg-gray-800 rounded-lg p-4 sm:p-6 mb-4 border border-gray-700">

          {/* Mobile: empilhado | Desktop: lado a lado */}
          <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4">

            {/* Avatar + info */}
            <div className="flex gap-4 items-center">
              <AvatarUpload
                characterID={Number(id)}
                avatarURL={character.avatar_url}
                characterName={character.name}
              />
              <div>
                <h1 className="text-2xl sm:text-3xl font-bold text-white">{character.name}</h1>
                <p className="text-gray-400 text-sm mt-0.5">Edição: {character.edition}</p>
                <div className="flex flex-wrap gap-2 mt-2">
                  <span className="bg-indigo-900 text-indigo-300 px-3 py-1 rounded-full text-xs sm:text-sm">
                    {character.class?.name ?? 'Sem classe'}
                  </span>
                  <span className="bg-emerald-900 text-emerald-300 px-3 py-1 rounded-full text-xs sm:text-sm">
                    {character.race?.name ?? 'Sem raça'}
                  </span>
                </div>
              </div>
            </div>

            {/* Nível + stats rápidos */}
            <div className="flex sm:flex-col items-center sm:items-end gap-4 sm:gap-1 sm:text-right">
              <div className="text-center sm:text-right">
                <p className="text-gray-400 text-xs">Nível</p>
                <p className="text-4xl sm:text-5xl font-bold text-yellow-400 leading-none">{character.level}</p>
              </div>
              <div className="text-sm text-gray-300 space-y-0.5">
                <p>HP: <span className="text-white font-semibold">{character.hit_points}/{character.max_hp}</span></p>
                {is4e && character.surge_value > 0 && (
                  <p>Pulso: <span className="text-white font-semibold">{character.surge_value}</span></p>
                )}
                {is4e && character.defense_ac > 0 && (
                  <p>CA: <span className="text-indigo-300 font-semibold">{character.defense_ac}</span></p>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* HP Manager */}
        <HPManager character={character} />

        {/* Atributos */}
        <div className="bg-gray-800 rounded-lg p-4 sm:p-6 mb-4 border border-gray-700">
          <h2 className="text-base sm:text-lg font-semibold text-white mb-3">Atributos</h2>
          <div className="grid grid-cols-3 sm:grid-cols-6 gap-2 sm:gap-3">
            {attributes.map(attr => {
              const modVal = Math.floor((attr.value - 10) / 2)
              return (
                <div key={attr.label} className="text-center bg-gray-700 rounded-lg p-2 sm:p-3">
                  <p className="text-gray-400 text-xs mb-1">{attr.label}</p>
                  <p className="text-white font-bold text-lg sm:text-xl">{attr.value}</p>
                  <p className={`text-xs font-semibold ${modVal >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                    {modVal >= 0 ? '+' : ''}{modVal}
                  </p>
                </div>
              )
            })}
          </div>
        </div>

        {/* Defesas (4e) */}
        {is4e && defenses.some(d => d.value > 0) && (
          <div className="bg-gray-800 rounded-lg p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-base sm:text-lg font-semibold text-white mb-3">Defesas</h2>
            <div className="grid grid-cols-4 gap-2 sm:gap-3">
              {defenses.map(d => (
                <div key={d.label} className="text-center bg-gray-700 rounded-lg p-2 sm:p-3">
                  <p className="text-gray-400 text-xs mb-1">{d.label}</p>
                  <p className="text-indigo-300 font-bold text-lg sm:text-xl">{d.value}</p>
                </div>
              ))}
            </div>
            {is4e && character.surges_per_day > 0 && (
              <p className="text-gray-400 text-xs mt-3">
                Pulsos de Cura: <span className="text-white font-semibold">{character.surges_per_day}/dia</span>
                {' '}(valor: <span className="text-white font-semibold">{character.surge_value} PV</span>)
              </p>
            )}
          </div>
        )}

        {/* Habilidades */}
        <SkillsPanel
          skills={character.skills ?? []}
          edition={character.edition}
        />

        {/* Background */}
        <BackgroundForm
          characterID={Number(id)}
          background={character.background}
        />

        {/* Ações */}
        <div className="flex flex-wrap gap-3 pb-6">
          <button
            onClick={() => navigate(`/characters/${id}/edit`)}
            className="flex-1 sm:flex-none bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-5 py-2.5 rounded-lg transition text-sm"
          >
            ✏️ Editar
          </button>
          <button
            onClick={() => levelUpMutation.mutate()}
            disabled={levelUpMutation.isPending || character.level >= 20}
            className="flex-1 sm:flex-none bg-yellow-600 hover:bg-yellow-700 disabled:opacity-50 text-white font-semibold px-5 py-2.5 rounded-lg transition text-sm"
          >
            {levelUpMutation.isPending ? 'Subindo...' : '⬆ Level Up'}
          </button>
          <button
            onClick={() => {
              if (confirm('Tem certeza que deseja deletar este personagem?')) deleteMutation.mutate()
            }}
            disabled={deleteMutation.isPending}
            className="flex-1 sm:flex-none bg-red-700 hover:bg-red-800 disabled:opacity-50 text-white font-semibold px-5 py-2.5 rounded-lg transition text-sm"
          >
            {deleteMutation.isPending ? 'Deletando...' : '🗑️ Deletar'}
          </button>
        </div>

      </div>
    </div>
  )
}