import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useParams } from 'react-router-dom'
import { characterService } from '../services/characterService'
import { periciaService } from '../services/periciaService'
import { talentoService } from '../services/talentoService'
import BackgroundForm from '../components/BackgroundForm'
import AvatarUpload from '../components/AvatarUpload'
import HPManager from '../components/HPManager'
import SkillsPanel from '../components/SkillsPanel'
import { Tooltip } from '../components/Tooltip'
import type { Pericia, Talento } from '../types'

const CATEGORY_CONFIG: Record<string, { color: string; icon: string }> = {
  'Combate':  { color: 'text-red-400',    icon: '⚔️' },
  'Defesa':   { color: 'text-blue-400',   icon: '🛡️' },
  'Perícia':  { color: 'text-yellow-400', icon: '📚' },
  'Magia':    { color: 'text-purple-400', icon: '✨' },
  'Armadura': { color: 'text-gray-300',   icon: '🪖' },
}

export default function CharacterDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: character, isLoading } = useQuery({
    queryKey: ['character', id],
    queryFn: () => characterService.getByID(Number(id)),
  })

  const { data: characterPericias } = useQuery({
    queryKey: ['character-pericias', id],
    queryFn: () => periciaService.getByCharacter(Number(id)),
    enabled: !!id,
  })

  const { data: allPericias } = useQuery({
    queryKey: ['pericias', character?.edition],
    queryFn: () => periciaService.getAll(character?.edition),
    enabled: !!character && character.edition === '4e',
    staleTime: Infinity,
  })

  const { data: characterTalentos } = useQuery({
    queryKey: ['character-talentos', id],
    queryFn: () => talentoService.getByCharacter(Number(id)),
    enabled: !!id,
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
        <p className="text-gray-500">Carregando...</p>
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

  const talentosByCategory = (characterTalentos ?? []).reduce<Record<string, Talento[]>>(
    (acc, t) => {
      const cat = t.category ?? 'Outros'
      acc[cat] = [...(acc[cat] ?? []), t]
      return acc
    }, {}
  )

  const hasPericias = (characterPericias ?? []).length > 0
  const hasTalentos = (characterTalentos ?? []).length > 0

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-3xl mx-auto">

        {/* Voltar */}
        <button
          onClick={() => navigate('/characters')}
          className="transition mb-4 sm:mb-6 block text-sm"
          style={{ color: 'rgba(201,168,76,0.5)' }}
          onMouseEnter={e => (e.currentTarget.style.color = '#c9a84c')}
          onMouseLeave={e => (e.currentTarget.style.color = 'rgba(201,168,76,0.5)')}
        >
          ← Voltar
        </button>

        {/* ── Cabeçalho ──────────────────────────────────────────────────────── */}
        <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
          <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4">
            <div className="flex gap-4 items-center">
              <AvatarUpload
                characterID={Number(id)}
                avatarURL={character.avatar_url}
                characterName={character.name}
              />
              <div>
                <h1 className="text-2xl sm:text-3xl font-bold text-white font-rpg">{character.name}</h1>
                <p className="text-gray-500 text-xs mt-0.5 uppercase tracking-wider">D&D {character.edition}</p>
                <div className="flex flex-wrap gap-2 mt-2">
                  {/* Badge classe — dourado */}
                  <span
                    className="px-3 py-1 rounded-full text-xs font-medium"
                    style={{ background: 'rgba(201,168,76,0.12)', border: '1px solid rgba(201,168,76,0.3)', color: '#c9a84c' }}
                  >
                    {character.class?.name ?? 'Sem classe'}
                  </span>
                  {/* Badge raça — esmeralda */}
                  <span className="bg-emerald-900/60 text-emerald-300 border border-emerald-700/50 px-3 py-1 rounded-full text-xs font-medium">
                    {character.race?.name ?? 'Sem raça'}
                  </span>
                </div>
              </div>
            </div>

            <div className="flex sm:flex-col items-center sm:items-end gap-4 sm:gap-1 sm:text-right">
              <div className="text-center sm:text-right">
                <p className="text-gray-500 text-xs uppercase tracking-widest">Nível</p>
                <p className="text-4xl sm:text-5xl font-bold leading-none font-rpg" style={{ color: '#c9a84c' }}>
                  {character.level}
                </p>
              </div>
              <div className="text-sm text-gray-400 space-y-0.5">
                <p>HP: <span className="text-white font-semibold">{character.hit_points}/{character.max_hp}</span></p>
                {is4e && character.surge_value > 0 && (
                  <p>Pulso: <span className="text-white font-semibold">{character.surge_value}</span></p>
                )}
                {is4e && character.defense_ac > 0 && (
                  <p>CA: <span className="font-semibold" style={{ color: '#c9a84c' }}>{character.defense_ac}</span></p>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* ── HP Manager ─────────────────────────────────────────────────────── */}
        <HPManager character={character} />

        {/* ── Atributos ──────────────────────────────────────────────────────── */}
        <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
          <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
            Atributos
          </h2>
          <div className="grid grid-cols-3 sm:grid-cols-6 gap-2 sm:gap-3">
            {attributes.map(attr => {
              const modVal = Math.floor((attr.value - 10) / 2)
              return (
                <div key={attr.label} className="text-center bg-gray-700/60 rounded-lg p-2 sm:p-3 border border-gray-600/50">
                  <p className="text-gray-500 text-xs mb-1 uppercase tracking-widest">{attr.label}</p>
                  <p className="text-white font-bold text-lg sm:text-xl">{attr.value}</p>
                  <p className={`text-xs font-semibold ${modVal >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                    {modVal >= 0 ? '+' : ''}{modVal}
                  </p>
                </div>
              )
            })}
          </div>
        </div>

        {/* ── Defesas (4e) ───────────────────────────────────────────────────── */}
        {is4e && defenses.some(d => d.value > 0) && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Defesas
            </h2>
            <div className="grid grid-cols-4 gap-2 sm:gap-3">
              {defenses.map(d => (
                <div key={d.label} className="text-center bg-gray-700/60 rounded-lg p-2 sm:p-3 border border-gray-600/50">
                  <p className="text-gray-500 text-xs mb-1 uppercase tracking-widest">{d.label}</p>
                  <p className="font-bold text-lg sm:text-xl" style={{ color: '#c9a84c' }}>{d.value}</p>
                </div>
              ))}
            </div>
            {is4e && character.surges_per_day > 0 && (
              <p className="text-gray-500 text-xs mt-3">
                Pulsos de Cura:{' '}
                <span className="text-white font-semibold">{character.surges_per_day}/dia</span>
                {' '}(valor:{' '}
                <span className="text-white font-semibold">{character.surge_value} PV</span>)
              </p>
            )}
          </div>
        )}

        {/* ── Péricias Treinadas (4e) ─────────────────────────────────────────── */}
        {is4e && hasPericias && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              📚 Perícias Treinadas
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {(characterPericias ?? []).map(cp => {
                const info: Pericia | undefined = allPericias?.find(p => p.name === cp.pericia_name)
                return (
                  <div
                    key={cp.pericia_name}
                    className="flex items-center justify-between bg-gray-700/60 rounded-lg px-3 py-2.5 border border-gray-600/50"
                  >
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-white text-sm font-medium">{cp.pericia_name}</span>
                      {info && <span className="text-gray-500 text-xs">({info.attribute})</span>}
                      <span className="text-xs font-semibold" style={{ color: '#5eead4' }}>+5</span>
                    </div>
                    {info && <Tooltip content={info.tooltip} />}
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* ── Talentos (4e) ──────────────────────────────────────────────────── */}
        {is4e && hasTalentos && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-4" style={{ color: 'rgba(201,168,76,0.7)' }}>
              🏆 Talentos
            </h2>
            <div className="flex flex-col gap-5">
              {Object.entries(talentosByCategory).map(([category, talentos]) => {
                const cfg = CATEGORY_CONFIG[category] ?? { color: 'text-gray-300', icon: '📌' }
                return (
                  <div key={category}>
                    <h3 className={`text-xs font-bold uppercase tracking-wider mb-2 ${cfg.color}`}>
                      {cfg.icon} {category}
                    </h3>
                    <div className="flex flex-col gap-2">
                      {talentos.map((t: Talento) => (
                        <div
                          key={t.ID}
                          className="flex items-start justify-between bg-gray-700/60 rounded-lg px-3 py-2.5 border border-gray-600/50"
                        >
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 flex-wrap mb-0.5">
                              <span className="text-white text-sm font-medium">{t.name}</span>
                              {t.prerequisite && (
                                <span className="text-xs bg-orange-900/60 text-orange-300 px-1.5 py-0.5 rounded border border-orange-700/40">
                                  Req: {t.prerequisite}
                                </span>
                              )}
                            </div>
                            <p className="text-gray-400 text-xs">{t.description}</p>
                          </div>
                          <div className="ml-3 flex-shrink-0 mt-0.5">
                            <Tooltip content={t.tooltip} />
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* ── Habilidades ────────────────────────────────────────────────────── */}
        <SkillsPanel skills={character.skills ?? []} edition={character.edition} />

        {/* ── Background ─────────────────────────────────────────────────────── */}
        <BackgroundForm characterID={Number(id)} background={character.background} />

        {/* ── Ações ──────────────────────────────────────────────────────────── */}
        <div className="flex flex-wrap gap-2 pb-6 pt-2">

          {/* Editar — outline dourado, sem emoji */}
          <button
            onClick={() => navigate(`/characters/${id}/edit`)}
            className="btn-rpg-outline flex-1 sm:flex-none"
          >
            Editar
          </button>

          {/* Level Up — primário dourado */}
          <button
            onClick={() => levelUpMutation.mutate()}
            disabled={levelUpMutation.isPending || character.level >= 20}
            className="btn-rpg-primary flex-1 sm:flex-none"
          >
            {levelUpMutation.isPending ? 'Subindo...' : '▲ Level Up'}
          </button>

          {/* Deletar — destrutivo, sem emoji */}
          <button
            onClick={() => {
              if (confirm('Tem certeza que deseja deletar este personagem?')) deleteMutation.mutate()
            }}
            disabled={deleteMutation.isPending}
            className="btn-rpg-danger flex-1 sm:flex-none"
          >
            {deleteMutation.isPending ? 'Deletando...' : 'Deletar'}
          </button>

        </div>

      </div>
    </div>
  )
}