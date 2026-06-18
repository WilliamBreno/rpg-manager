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
import type { Pericia, Talento, Antecedent } from '../types'

const CATEGORY_CONFIG: Record<string, { color: string; icon: string }> = {
  'Combate':  { color: 'text-red-400',    icon: '⚔️' },
  'Defesa':   { color: 'text-blue-400',   icon: '🛡️' },
  'Perícia':  { color: 'text-yellow-400', icon: '📚' },
  'Magia':    { color: 'text-purple-400', icon: '✨' },
  'Armadura': { color: 'text-gray-300',   icon: '🪖' },
}

// ── Helpers 5e ────────────────────────────────────────────────────────────────
function profBonusFor(level: number): number {
  return Math.floor((level - 1) / 4) + 2
}
function abilityMod(score: number): number {
  return Math.floor((score - 10) / 2)
}
function fmtMod(n: number): string {
  return n >= 0 ? `+${n}` : `${n}`
}

const ATTR_KEYS = ['FOR','DES','CON','INT','SAB','CAR'] as const
const ATTR_NAMES: Record<string, string> = {
  FOR: 'Força', DES: 'Destreza', CON: 'Constituição',
  INT: 'Inteligência', SAB: 'Sabedoria', CAR: 'Carisma',
}
const ATTR_FIELD: Record<string, string> = {
  FOR: 'strength', DES: 'dexterity', CON: 'constitution',
  INT: 'intelligence', SAB: 'wisdom', CAR: 'charisma',
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

  // Carrega todas as perícias — 4e OU 5e
  const { data: allPericias } = useQuery({
    queryKey: ['pericias', character?.edition],
    queryFn: () => periciaService.getAll(character?.edition),
    enabled: !!character,
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
  const is5e = character.edition === '5e'

  // ── Dados 5e calculados ────────────────────────────────────────────────────
  const profBonus = profBonusFor(character.level)
  const initiative = abilityMod(character.dexterity)
  const speed = character.speed || (character.race as any)?.speed || 30

  const saveProficiencies: string[] = (() => {
    if (!character.class?.saving_throws) return []
    try { return JSON.parse(character.class.saving_throws) } catch { return [] }
  })()

  const trainedPericiaNames = (characterPericias ?? []).map((cp: any) => cp.pericia_name)

  const getAttrScore = (abbr: string): number => {
    const field = ATTR_FIELD[abbr]
    return (character as any)[field] ?? 10
  }

  const getSaveValue = (abbr: string): number => {
    const mod = abilityMod(getAttrScore(abbr))
    return mod + (saveProficiencies.includes(abbr) ? profBonus : 0)
  }

  const getSkillValue = (periciaName: string, attribute: string): number => {
    const attrKey = Object.keys(ATTR_FIELD).find(k => ATTR_FIELD[k] === attribute.toLowerCase()) ?? attribute
    const mod = abilityMod(getAttrScore(attrKey))
    const proficient = trainedPericiaNames.includes(periciaName)
    return mod + (proficient ? profBonus : 0)
  }

  // Percepção Passiva 5e
  const perceptionPericia = allPericias?.find((p: Pericia) => p.name === 'Percepção')
  const passivePerception = perceptionPericia
    ? 10 + getSkillValue('Percepção', perceptionPericia.attribute)
    : 10 + abilityMod(character.wisdom)

  // ── Dados 4e ──────────────────────────────────────────────────────────────
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
    (acc, t) => { const cat = t.category ?? 'Outros'; acc[cat] = [...(acc[cat] ?? []), t]; return acc }, {}
  )

  const hasPericias = (characterPericias ?? []).length > 0
  const hasTalentos = (characterTalentos ?? []).length > 0

 const bg5e: Antecedent | undefined = is5e ? character.antecedent : undefined

  // ── Render ─────────────────────────────────────────────────────────────────
  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-3xl mx-auto">

        <button onClick={() => navigate('/characters')} className="transition mb-4 sm:mb-6 block text-sm"
          style={{ color: 'rgba(201,168,76,0.5)' }}
          onMouseEnter={e => (e.currentTarget.style.color = '#c9a84c')}
          onMouseLeave={e => (e.currentTarget.style.color = 'rgba(201,168,76,0.5)')}
        >← Voltar</button>

        {/* ── Cabeçalho ──────────────────────────────────────────────────── */}
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
                  <span className="px-3 py-1 rounded-full text-xs font-medium"
                    style={{ background: 'rgba(201,168,76,0.12)', border: '1px solid rgba(201,168,76,0.3)', color: '#c9a84c' }}>
                    {character.class?.name ?? 'Sem classe'}
                  </span>
                  <span className="bg-emerald-900/60 text-emerald-300 border border-emerald-700/50 px-3 py-1 rounded-full text-xs font-medium">
                    {character.race?.name ?? 'Sem raça'}
                  </span>
                  {is5e && bg5e && (
                    <span className="bg-violet-900/60 text-violet-300 border border-violet-700/50 px-3 py-1 rounded-full text-xs font-medium">
                      📜 {bg5e.name}
                    </span>
                  )}
                  {is5e && character.alignment && (
                    <span className="bg-sky-900/60 text-sky-300 border border-sky-700/50 px-3 py-1 rounded-full text-xs font-medium">
                      ⚖️ {character.alignment}
                    </span>
                  )}
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
                {is5e && (
                  <p>Prof: <span className="font-semibold" style={{ color: '#c9a84c' }}>+{profBonus}</span></p>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* ── HP Manager ──────────────────────────────────────────────────── */}
        <HPManager character={character} />

        {/* ── Stats de Combate 5e ─────────────────────────────────────────── */}
        {is5e && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Combate
            </h2>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3">
              {[
                { label: 'CA',               value: character.defense_ac || '—', sub: 'Armadura'         },
                { label: fmtMod(initiative),  value: null, sub: 'Iniciativa', raw: initiative             },
                { label: `${speed}`,          value: null, sub: 'Deslocamento (pés)'                      },
                { label: `+${profBonus}`,     value: null, sub: 'Bônus de Proficiência'                   },
              ].map((stat, i) => (
                <div key={i} className="text-center bg-gray-700/60 rounded-lg p-3 border border-gray-600/50">
                  <p className="font-bold text-xl sm:text-2xl" style={{ color: '#c9a84c' }}>
                    {stat.label}
                  </p>
                  <p className="text-gray-500 text-xs mt-1 leading-tight">{stat.sub}</p>
                </div>
              ))}
            </div>
            <p className="text-gray-600 text-xs mt-3">
              Percepção Passiva (Sab): <span className="text-gray-400 font-semibold">{passivePerception}</span>
              <span className="ml-3">Dados de Vida: <span className="text-gray-400 font-semibold">{character.level}d{character.class?.hit_die ?? '?'}</span></span>
            </p>
          </div>
        )}

        {/* ── Atributos ───────────────────────────────────────────────────── */}
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

        {/* ── Defesas (4e) ────────────────────────────────────────────────── */}
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
            {character.surges_per_day > 0 && (
              <p className="text-gray-500 text-xs mt-3">
                Pulsos de Cura: <span className="text-white font-semibold">{character.surges_per_day}/dia</span>
                {' '}(valor: <span className="text-white font-semibold">{character.surge_value} PV</span>)
              </p>
            )}
          </div>
        )}

        {/* ── Testes de Resistência 5e ─────────────────────────────────────── */}
        {is5e && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Testes de Resistência
            </h2>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
              {ATTR_KEYS.map(abbr => {
                const proficient = saveProficiencies.includes(abbr)
                const val = getSaveValue(abbr)
                return (
                  <div key={abbr}
                    className="flex items-center gap-3 rounded-lg px-3 py-2.5 border"
                    style={proficient
                      ? { background: 'rgba(201,168,76,0.07)', borderColor: 'rgba(201,168,76,0.3)' }
                      : { background: '#27272a', borderColor: '#3f3f46' }
                    }
                  >
                    {/* Marcador de proficiência */}
                    <div className="w-3.5 h-3.5 rounded-full border-2 flex-shrink-0 flex items-center justify-center"
                      style={proficient ? { background: '#c9a84c', borderColor: '#c9a84c' } : { borderColor: '#52525b' }}>
                      {proficient && <span className="text-black font-bold" style={{ fontSize: '8px' }}>✓</span>}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-gray-400 text-xs">{ATTR_NAMES[abbr]}</p>
                    </div>
                    <span className={`text-sm font-bold flex-shrink-0 ${val >= 0 ? 'text-white' : 'text-red-400'}`}>
                      {val >= 0 ? '+' : ''}{val}
                    </span>
                  </div>
                )
              })}
            </div>
            {saveProficiencies.length > 0 && (
              <p className="text-gray-600 text-xs mt-2">
                Proficiências: {saveProficiencies.map(s => ATTR_NAMES[s]).join(', ')}
              </p>
            )}
          </div>
        )}

        {/* ── 18 Perícias 5e ──────────────────────────────────────────────── */}
        {is5e && allPericias && allPericias.length > 0 && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <div className="flex items-center justify-between mb-3">
              <h2 className="text-sm font-semibold uppercase tracking-widest" style={{ color: 'rgba(201,168,76,0.7)' }}>
                Perícias
              </h2>
              <span className="text-xs text-gray-500 bg-gray-700/60 px-2 py-1 rounded-full">
                {trainedPericiaNames.length} treinadas
              </span>
            </div>
            <div className="flex flex-col gap-1.5">
              {allPericias.map((p: Pericia) => {
                const proficient = trainedPericiaNames.includes(p.name)
                const val = getSkillValue(p.name, p.attribute)
                return (
                  <div key={p.ID}
                    className="flex items-center gap-3 rounded-lg px-3 py-2 border"
                    style={proficient
                      ? { background: 'rgba(201,168,76,0.07)', borderColor: 'rgba(201,168,76,0.25)' }
                      : { background: 'transparent', borderColor: 'transparent' }
                    }
                  >
                    {/* Proficiency dot */}
                    <div className="w-3 h-3 rounded-full border-2 flex-shrink-0"
                      style={proficient ? { background: '#c9a84c', borderColor: '#c9a84c' } : { borderColor: '#52525b' }} />

                    {/* Value */}
                    <span className={`text-sm font-bold w-8 flex-shrink-0 text-right ${val >= 0 ? (proficient ? '' : 'text-gray-400') : 'text-red-400'}`}
                      style={proficient ? { color: '#c9a84c' } : {}}>
                      {val >= 0 ? '+' : ''}{val}
                    </span>

                    {/* Name */}
                    <span className={`text-sm flex-1 ${proficient ? 'text-white font-medium' : 'text-gray-400'}`}>{p.name}</span>

                    {/* Attribute badge */}
                    <span className="text-gray-600 text-xs flex-shrink-0">({p.attribute})</span>

                    <Tooltip content={p.tooltip} />
                  </div>
                )
              })}
            </div>

            {/* Percepção Passiva */}
            <div className="mt-3 pt-3 border-t border-gray-700">
              <p className="text-gray-500 text-xs">
                Sabedoria Passiva (Percepção):{' '}
                <span className="text-white font-semibold text-sm">{passivePerception}</span>
              </p>
            </div>
          </div>
        )}

        {/* ── Perícias Treinadas (4e) ─────────────────────────────────────── */}
        {is4e && hasPericias && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              📚 Perícias Treinadas
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {(characterPericias ?? []).map((cp: any) => {
                const info: Pericia | undefined = allPericias?.find((p: Pericia) => p.name === cp.pericia_name)
                return (
                  <div key={cp.pericia_name}
                    className="flex items-center justify-between bg-gray-700/60 rounded-lg px-3 py-2.5 border border-gray-600/50">
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

        {/* ── Talentos (4e) ───────────────────────────────────────────────── */}
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
                        <div key={t.ID}
                          className="flex items-start justify-between bg-gray-700/60 rounded-lg px-3 py-2.5 border border-gray-600/50">
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
                          <div className="ml-3 flex-shrink-0 mt-0.5"><Tooltip content={t.tooltip} /></div>
                        </div>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* ── Antecedente 5e ──────────────────────────────────────────────── */}
        {is5e && bg5e && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              📜 Antecedente — {bg5e.name}
            </h2>
            <p className="text-gray-400 text-sm mb-3">{bg5e.description}</p>

            {/* Feature */}
            <div className="rounded-lg p-3 mb-3" style={{ background: 'rgba(201,168,76,0.06)', border: '1px solid rgba(201,168,76,0.2)' }}>
              <p className="text-xs font-semibold mb-1" style={{ color: '#c9a84c' }}>✦ {bg5e.feature}</p>
              <p className="text-gray-400 text-xs">{bg5e.feature_description}</p>
            </div>

            {/* Proficiências + idiomas */}
            <div className="flex flex-wrap gap-2 mb-2">
              {(() => { try { return JSON.parse(bg5e.skill_proficiencies) } catch { return [] } })()
                .map((s: string) => (
                  <span key={s} className="text-xs bg-indigo-900/60 text-indigo-300 px-2 py-0.5 rounded-full">📚 {s}</span>
                ))}
              {bg5e.tool_proficiencies && (
                <span className="text-xs bg-orange-900/60 text-orange-300 px-2 py-0.5 rounded-full">🔧 {bg5e.tool_proficiencies}</span>
              )}
              {bg5e.languages && (
                <span className="text-xs bg-teal-900/60 text-teal-300 px-2 py-0.5 rounded-full">🗣 {bg5e.languages}</span>
              )}
            </div>

            {bg5e.equipment && (
              <p className="text-gray-500 text-xs">
                <span className="text-gray-400 font-medium">Equipamento: </span>{bg5e.equipment}
              </p>
            )}
          </div>
        )}

        {/* ── Personalidade 5e ────────────────────────────────────────────── */}
        {is5e && (character.personality_traits || character.ideals || character.bonds || character.flaws || character.alignment) && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-4" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Personalidade
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {character.alignment && (
                <div className="col-span-full">
                  <p className="text-gray-500 text-xs uppercase tracking-wider mb-1">Tendência</p>
                  <span className="text-sm bg-sky-900/60 text-sky-300 border border-sky-700/50 px-3 py-1.5 rounded-lg inline-block">
                    ⚖️ {character.alignment}
                  </span>
                </div>
              )}
              {[
                { key: 'personality_traits', label: 'Traços de Personalidade', icon: '💬' },
                { key: 'ideals',             label: 'Ideais',   icon: '⭐' },
                { key: 'bonds',              label: 'Ligações', icon: '🔗' },
                { key: 'flaws',              label: 'Defeitos', icon: '⚠️' },
              ].map(field => {
                const val = (character as any)[field.key]
                if (!val) return null
                return (
                  <div key={field.key} className="rounded-lg p-3 border border-gray-700/60 bg-gray-700/30">
                    <p className="text-gray-500 text-xs uppercase tracking-wider mb-1">{field.icon} {field.label}</p>
                    <p className="text-gray-300 text-sm leading-relaxed">{val}</p>
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* ── Habilidades / Skills ─────────────────────────────────────────── */}
        <SkillsPanel skills={character.skills ?? []} edition={character.edition} />

        {/* ── Background Form (biografia/notas) ───────────────────────────── */}
        <BackgroundForm characterID={Number(id)} />

        {/* ── Ações ───────────────────────────────────────────────────────── */}
        <div className="flex flex-wrap gap-2 pb-6 pt-2">
          <button onClick={() => navigate(`/characters/${id}/edit`)} className="btn-rpg-outline flex-1 sm:flex-none">
            Editar
          </button>
          <button
            onClick={() => levelUpMutation.mutate()}
            disabled={levelUpMutation.isPending || character.level >= 20}
            className="btn-rpg-primary flex-1 sm:flex-none"
          >
            {levelUpMutation.isPending ? 'Subindo...' : '▲ Level Up'}
          </button>
          <button
            onClick={() => { if (confirm('Tem certeza que deseja deletar este personagem?')) deleteMutation.mutate() }}
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