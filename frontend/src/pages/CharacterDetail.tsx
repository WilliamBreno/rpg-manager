import { useState } from 'react'
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
import DeathSaves from '../components/DeathSaves'
// ── Tabelas de XP ────────────────────────────────────────────────────────────
const XP_TABLE_4E: Record<number, number> = {
  1: 0, 2: 1000, 3: 2250, 4: 3750, 5: 5500,
  6: 7500, 7: 10000, 8: 13000, 9: 16500, 10: 20500,
  11: 26000, 12: 32000, 13: 39000, 14: 47000, 15: 57000,
  16: 69000, 17: 83000, 18: 99000, 19: 119000, 20: 143000,
  21: 175000, 22: 210000, 23: 255000, 24: 310000, 25: 375000,
  26: 450000, 27: 550000, 28: 675000, 29: 825000, 30: 1000000,
}
const XP_TABLE_5E: Record<number, number> = {
  1: 0, 2: 300, 3: 900, 4: 2700, 5: 6500,
  6: 14000, 7: 23000, 8: 34000, 9: 48000, 10: 64000,
  11: 85000, 12: 100000, 13: 120000, 14: 140000, 15: 165000,
  16: 195000, 17: 225000, 18: 265000, 19: 305000, 20: 355000,
}

// ── Níveis ASI 5e ─────────────────────────────────────────────────────────────
const ASI_LEVELS_5E = new Set([4, 8, 12, 16, 19])

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

const ATTR_KEYS = ['FOR', 'DES', 'CON', 'INT', 'SAB', 'CAR'] as const
const ATTR_NAMES: Record<string, string> = {
  FOR: 'Força', DES: 'Destreza', CON: 'Constituição',
  INT: 'Inteligência', SAB: 'Sabedoria', CAR: 'Carisma',
}
const ATTR_FIELD: Record<string, string> = {
  FOR: 'strength', DES: 'dexterity', CON: 'constitution',
  INT: 'intelligence', SAB: 'wisdom', CAR: 'charisma',
}

// Mapeamento atributo → chave ASI (snake_case do backend)
const ASI_ATTRS = [
  { label: 'Força',        key: 'strength'     },
  { label: 'Destreza',     key: 'dexterity'    },
  { label: 'Constituição', key: 'constitution' },
  { label: 'Inteligência', key: 'intelligence' },
  { label: 'Sabedoria',    key: 'wisdom'       },
  { label: 'Carisma',      key: 'charisma'     },
]

const EMPTY_ASI = { strength: 0, dexterity: 0, constitution: 0, intelligence: 0, wisdom: 0, charisma: 0 }

// ── Tipos locais ──────────────────────────────────────────────────────────────
type XPFeedback = { message: string; type: 'success' | 'level' | 'error' }
type ASIChoices = typeof EMPTY_ASI

export default function CharacterDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  // ── Estado local ─────────────────────────────────────────────────────────
  const [xpInput, setXpInput] = useState('')
  const [xpFeedback, setXpFeedback] = useState<XPFeedback | null>(null)
  const [showASIModal, setShowASIModal] = useState(false)
  const [asiChoices, setAsiChoices] = useState<ASIChoices>({ ...EMPTY_ASI })

  // ── Queries ──────────────────────────────────────────────────────────────
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
    enabled: !!character,
    staleTime: Infinity,
  })

  const { data: characterTalentos } = useQuery({
    queryKey: ['character-talentos', id],
    queryFn: () => talentoService.getByCharacter(Number(id)),
    enabled: !!id,
  })

  // ── Mutations ────────────────────────────────────────────────────────────

  // Adicionar XP — dispara level up automático no backend
  const addXPMutation = useMutation({
    mutationFn: (xp: number) => characterService.addXP(Number(id), xp),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['character', id] })
      setXpInput('')
      if (data.leveled_up) {
        if (data.needs_asi) {
          // Nível ASI: abre modal de melhoria de atributo
          setXpFeedback({ message: `⬆️ Subiu para o nível ${data.new_level}! Escolha uma melhoria de atributo.`, type: 'level' })
          setShowASIModal(true)
        } else {
          setXpFeedback({ message: `⬆️ Subiu para o nível ${data.new_level}!`, type: 'level' })
          setTimeout(() => setXpFeedback(null), 5000)
        }
      } else {
        setXpFeedback({ message: `✓ XP adicionado com sucesso`, type: 'success' })
        setTimeout(() => setXpFeedback(null), 3000)
      }
    },
    onError: () => {
      setXpFeedback({ message: 'Erro ao adicionar XP', type: 'error' })
      setTimeout(() => setXpFeedback(null), 3000)
    },
  })

  // Aplicar ASI — chamado quando o jogador confirma as melhorias
  const applyASIMutation = useMutation({
    mutationFn: (choices: ASIChoices) => characterService.applyASI(Number(id), choices),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['character', id] })
      setShowASIModal(false)
      setAsiChoices({ ...EMPTY_ASI })
      if (data.leveled_up) {
        if (data.needs_asi) {
          // Outro nível ASI pendente — reabre o modal após breve pausa
          setXpFeedback({ message: `⬆️ Subiu para o nível ${data.new_level}! Escolha outra melhoria.`, type: 'level' })
          setTimeout(() => setShowASIModal(true), 500)
        } else {
          setXpFeedback({ message: `⬆️ Subiu para o nível ${data.new_level}!`, type: 'level' })
          setTimeout(() => setXpFeedback(null), 5000)
        }
      } else {
        setXpFeedback({ message: `✓ Atributos melhorados`, type: 'success' })
        setTimeout(() => setXpFeedback(null), 3000)
      }
    },
    onError: () => {
      setXpFeedback({ message: 'Erro ao aplicar melhorias', type: 'error' })
      setTimeout(() => setXpFeedback(null), 3000)
    },
  })

  const levelUpMutation = useMutation({
    mutationFn: () => characterService.levelUp(Number(id)),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character', id] }),
  })

  const deleteMutation = useMutation({
    mutationFn: () => characterService.delete(Number(id)),
    onSuccess: () => navigate('/characters'),
  })

  // ── Guards ───────────────────────────────────────────────────────────────
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
  const maxLvl = is4e ? 30 : 20
  const isMaxLevel = character.level >= maxLvl

  // ── Cálculo de XP ────────────────────────────────────────────────────────
  const xpTable = is4e ? XP_TABLE_4E : XP_TABLE_5E
  const currentXP = character.experience_points ?? 0
  const currentLevelXP = xpTable[character.level] ?? 0
  const nextLevelXP = xpTable[character.level + 1] ?? 0
  const xpNeeded = isMaxLevel ? 0 : Math.max(0, nextLevelXP - currentXP)
  const xpProgress = isMaxLevel
    ? 100
    : nextLevelXP > currentLevelXP
    ? Math.min(100, Math.max(0, ((currentXP - currentLevelXP) / (nextLevelXP - currentLevelXP)) * 100))
    : 100

  // ── Dados 5e calculados ──────────────────────────────────────────────────
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

  const perceptionPericia = allPericias?.find((p: Pericia) => p.name === 'Percepção')
  const passivePerception = perceptionPericia
    ? 10 + getSkillValue('Percepção', perceptionPericia.attribute)
    : 10 + abilityMod(character.wisdom)

  // ── Dados 4e ─────────────────────────────────────────────────────────────
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

  // ── ASI helpers ──────────────────────────────────────────────────────────
  const asiTotal = Object.values(asiChoices).reduce((a, b) => a + b, 0)
  const canConfirmASI = asiTotal >= 1 && asiTotal <= 2

  // ── Render ───────────────────────────────────────────────────────────────
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

        {/* ── Progressão (XP) ─────────────────────────────────────────────── */}
        <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold uppercase tracking-widest" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Progressão
            </h2>
            <span className="text-xs text-gray-500 font-mono">
              {currentXP.toLocaleString('pt-BR')} XP total
            </span>
          </div>

          {/* Barra de progresso */}
          <div className="mb-3">
            <div className="flex justify-between text-xs text-gray-500 mb-1.5">
              <span>Nível {character.level}</span>
              {isMaxLevel
                ? <span style={{ color: '#c9a84c' }}>✦ Nível máximo atingido</span>
                : <span>{xpNeeded.toLocaleString('pt-BR')} XP para o nível {character.level + 1}</span>
              }
            </div>
            <div className="w-full h-3 rounded-full overflow-hidden" style={{ background: '#1c1c1e' }}>
              <div
                className="h-full rounded-full transition-all duration-700"
                style={{
                  width: `${xpProgress}%`,
                  background: isMaxLevel
                    ? '#c9a84c'
                    : 'linear-gradient(90deg, rgba(201,168,76,0.6) 0%, #c9a84c 100%)',
                }}
              />
            </div>
            {!isMaxLevel && (
              <div className="flex justify-between text-xs text-gray-700 mt-1 font-mono">
                <span>{currentLevelXP.toLocaleString('pt-BR')}</span>
                <span>{nextLevelXP.toLocaleString('pt-BR')}</span>
              </div>
            )}
          </div>

          {/* Input para adicionar XP */}
          {!isMaxLevel && (
            <div className="flex gap-2">
              <input
                type="number"
                min={1}
                placeholder="XP ganho na aventura..."
                value={xpInput}
                onChange={e => setXpInput(e.target.value)}
                onKeyDown={e => {
                  if (e.key === 'Enter' && xpInput && Number(xpInput) > 0)
                    addXPMutation.mutate(Number(xpInput))
                }}
                className="rpg-input flex-1 text-sm"
              />
              <button
                disabled={!xpInput || Number(xpInput) <= 0 || addXPMutation.isPending}
                onClick={() => xpInput && Number(xpInput) > 0 && addXPMutation.mutate(Number(xpInput))}
                className="px-4 py-2 rounded-lg text-sm font-semibold transition flex-shrink-0"
                style={{
                  background: !xpInput || Number(xpInput) <= 0 ? '#27272a' : '#c9a84c',
                  color: !xpInput || Number(xpInput) <= 0 ? '#52525b' : '#0a0a0a',
                  cursor: !xpInput || Number(xpInput) <= 0 ? 'not-allowed' : 'pointer',
                }}
              >
                {addXPMutation.isPending ? '...' : '+ XP'}
              </button>
            </div>
          )}

          {/* Feedback */}
          {xpFeedback && (
            <div className="mt-2 text-xs text-center font-semibold rounded-lg p-2 transition"
              style={
                xpFeedback.type === 'level'
                  ? { background: 'rgba(201,168,76,0.12)', border: '1px solid rgba(201,168,76,0.3)', color: '#c9a84c' }
                  : xpFeedback.type === 'error'
                  ? { background: 'rgba(220,38,38,0.1)', border: '1px solid rgba(220,38,38,0.2)', color: '#f87171' }
                  : { background: 'rgba(16,185,129,0.08)', border: '1px solid rgba(16,185,129,0.2)', color: '#6ee7b7' }
              }
            >
              {xpFeedback.message}
            </div>
          )}

          {/* Indicador de níveis ASI 5e próximos */}
          {is5e && !isMaxLevel && (
            <p className="text-gray-600 text-xs mt-2">
              {ASI_LEVELS_5E.has(character.level + 1)
                ? <span className="text-amber-600">✦ Próximo nível concede melhoria de atributo (ASI)</span>
                : (() => {
                    const nextASI = [4, 8, 12, 16, 19].find(l => l > character.level)
                    return nextASI ? `Próxima melhoria de atributo: nível ${nextASI}` : null
                  })()
              }
            </p>
          )}
        </div>

        {/* ── HP Manager ──────────────────────────────────────────────────── */}
        <HPManager character={character} />
        {is5e && <DeathSaves characterId={Number(id)} character={character} />}
        {/* ── Stats de Combate 5e ─────────────────────────────────────────── */}
        {is5e && (
          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Combate
            </h2>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3">
              {[
                { label: String(character.defense_ac || '—'), sub: 'CA' },
                { label: fmtMod(initiative), sub: 'Iniciativa' },
                { label: String(speed), sub: 'Deslocamento (pés)' },
                { label: `+${profBonus}`, sub: 'Bônus de Prof.' },
              ].map((stat, i) => (
                <div key={i} className="text-center bg-gray-700/60 rounded-lg p-3 border border-gray-600/50">
                  <p className="font-bold text-xl sm:text-2xl" style={{ color: '#c9a84c' }}>{stat.label}</p>
                  <p className="text-gray-500 text-xs mt-1 leading-tight">{stat.sub}</p>
                </div>
              ))}
            </div>
            <p className="text-gray-600 text-xs mt-3">
              Percepção Passiva: <span className="text-gray-400 font-semibold">{passivePerception}</span>
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
          </div>
        )}

        {/* ── Perícias 5e ─────────────────────────────────────────────────── */}
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
                    <div className="w-3 h-3 rounded-full border-2 flex-shrink-0"
                      style={proficient ? { background: '#c9a84c', borderColor: '#c9a84c' } : { borderColor: '#52525b' }} />
                    <span
                      className="text-sm font-bold w-8 flex-shrink-0 text-right"
                      style={proficient ? { color: '#c9a84c' } : { color: '#9ca3af' }}
                    >
                      {val >= 0 ? '+' : ''}{val}
                    </span>
                    <span className={`text-sm flex-1 ${proficient ? 'text-white font-medium' : 'text-gray-400'}`}>{p.name}</span>
                    <span className="text-gray-600 text-xs flex-shrink-0">({p.attribute})</span>
                    <Tooltip content={p.tooltip} />
                  </div>
                )
              })}
            </div>
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
            <div className="rounded-lg p-3 mb-3" style={{ background: 'rgba(201,168,76,0.06)', border: '1px solid rgba(201,168,76,0.2)' }}>
              <p className="text-xs font-semibold mb-1" style={{ color: '#c9a84c' }}>✦ {bg5e.feature}</p>
              <p className="text-gray-400 text-xs">{bg5e.feature_description}</p>
            </div>
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
            disabled={levelUpMutation.isPending || isMaxLevel}
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

      {/* ── Modal ASI ────────────────────────────────────────────────────── */}
      {showASIModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4"
          style={{ background: 'rgba(0,0,0,0.8)' }}>
          <div className="bg-gray-800 rounded-xl p-6 border border-gray-600 max-w-md w-full shadow-2xl">

            <h2 className="font-rpg text-xl font-bold mb-1" style={{ color: '#c9a84c' }}>
              ⬆️ Melhoria de Atributo
            </h2>
            <p className="text-gray-400 text-sm mb-4">
              Distribua <strong className="text-white">2 pontos</strong> nos atributos.<br />
              +2 em um único ou +1 em dois diferentes. Máximo 20 por atributo.
            </p>

            {/* Contador de pontos */}
            <div className="flex items-center justify-center gap-3 mb-5 rounded-lg p-3"
              style={{ background: asiTotal === 2 ? 'rgba(201,168,76,0.1)' : 'rgba(255,255,255,0.03)', border: `1px solid ${asiTotal === 2 ? 'rgba(201,168,76,0.3)' : '#3f3f46'}` }}>
              <span className="text-3xl font-bold font-rpg" style={{ color: asiTotal === 2 ? '#c9a84c' : '#71717a' }}>
                {asiTotal}
              </span>
              <span className="text-gray-500 text-sm">/ 2 pontos distribuídos</span>
            </div>

            <div className="grid grid-cols-2 gap-3 mb-5">
              {ASI_ATTRS.map(attr => {
                const currentVal = (character as any)[attr.key] ?? 10
                const bonus = asiChoices[attr.key as keyof ASIChoices] ?? 0
                const canAdd = asiTotal < 2 && currentVal + bonus < 20
                const canRemove = bonus > 0

                return (
                  <div key={attr.key}
                    className="flex items-center justify-between rounded-lg px-3 py-2.5 border transition"
                    style={{
                      background: bonus > 0 ? 'rgba(201,168,76,0.08)' : '#27272a',
                      borderColor: bonus > 0 ? 'rgba(201,168,76,0.4)' : '#3f3f46',
                    }}
                  >
                    <div>
                      <p className="text-gray-500 text-xs uppercase tracking-wider leading-none mb-0.5">{attr.label}</p>
                      <p className="text-white font-bold text-base leading-none">
                        {currentVal}
                        {bonus > 0 && <span style={{ color: '#c9a84c' }}> → {currentVal + bonus}</span>}
                      </p>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <button
                        onClick={() => canRemove && setAsiChoices(prev => ({ ...prev, [attr.key]: prev[attr.key as keyof ASIChoices] - 1 }))}
                        disabled={!canRemove}
                        className="w-7 h-7 rounded text-base font-bold flex items-center justify-center transition"
                        style={{ background: canRemove ? '#3f3f46' : 'transparent', color: canRemove ? 'white' : '#3f3f46' }}
                      >−</button>
                      <span className="w-5 text-center text-sm font-bold" style={{ color: bonus > 0 ? '#c9a84c' : '#52525b' }}>
                        {bonus > 0 ? `+${bonus}` : '·'}
                      </span>
                      <button
                        onClick={() => canAdd && setAsiChoices(prev => ({ ...prev, [attr.key]: prev[attr.key as keyof ASIChoices] + 1 }))}
                        disabled={!canAdd}
                        className="w-7 h-7 rounded text-base font-bold flex items-center justify-center transition"
                        style={{ background: canAdd ? '#3f3f46' : 'transparent', color: canAdd ? 'white' : '#3f3f46' }}
                      >+</button>
                    </div>
                  </div>
                )
              })}
            </div>

            <button
              onClick={() => canConfirmASI && applyASIMutation.mutate(asiChoices)}
              disabled={!canConfirmASI || applyASIMutation.isPending}
              className="w-full py-3 rounded-lg font-semibold text-sm transition"
              style={canConfirmASI
                ? { background: '#c9a84c', color: '#0a0a0a', cursor: 'pointer' }
                : { background: '#27272a', color: '#52525b', cursor: 'not-allowed' }
              }
            >
              {applyASIMutation.isPending
                ? 'Aplicando...'
                : `Confirmar (+${asiTotal} ponto${asiTotal !== 1 ? 's' : ''})`
              }
            </button>

          </div>
        </div>
      )}

    </div>
  )
}