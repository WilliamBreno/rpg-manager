import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import { useForm, useWatch } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { classService } from '../services/classService'
import { raceService } from '../services/raceService'
import { skillService } from '../services/skillService'
import { characterService } from '../services/characterService'
import { armorService } from '../services/armorService'
import type { Class, Armor, Skill } from '../types'

const schema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  edition: z.string().min(1, 'Edição é obrigatória'),
  class_id: z.coerce.number().min(1, 'Selecione uma classe'),
  race_id: z.coerce.number().min(1, 'Selecione uma raça'),
  hit_points: z.coerce.number().min(1, 'HP deve ser maior que zero'),
  strength: z.coerce.number().min(1).max(20),
  dexterity: z.coerce.number().min(1).max(20),
  constitution: z.coerce.number().min(1).max(20),
  intelligence: z.coerce.number().min(1).max(20),
  wisdom: z.coerce.number().min(1).max(20),
  charisma: z.coerce.number().min(1).max(20),
  armor_id: z.coerce.number().optional(),
})

type FormData = z.infer<typeof schema>
type PowerType = 'unlimited' | 'encounter' | 'daily' | 'utility'

const SKILL_LIMITS: Record<PowerType, number> = {
  unlimited: 2,
  encounter: 1,
  daily: 1,
  utility: 1,
}

const powerConfig: Record<PowerType, { label: string; color: string; border: string; bg: string; badge: string }> = {
  unlimited: { label: '⚡ Sem Limite',    color: 'text-green-400',  border: 'border-green-800',  bg: 'bg-green-950',  badge: 'bg-green-900 text-green-300'  },
  encounter: { label: '⚔️ Por Encontro', color: 'text-yellow-400', border: 'border-yellow-800', bg: 'bg-yellow-950', badge: 'bg-yellow-900 text-yellow-300' },
  daily:     { label: '📅 Diário',        color: 'text-red-400',    border: 'border-red-800',    bg: 'bg-red-950',    badge: 'bg-red-900 text-red-300'      },
  utility:   { label: '🔧 Utilitário',    color: 'text-blue-400',   border: 'border-blue-800',   bg: 'bg-blue-950',   badge: 'bg-blue-900 text-blue-300'    },
}

const editions = ['4e', '5e']

function getConMod(con: number) {
  return Math.floor((con - 10) / 2)
}

// ── Componente de seção de poderes ──────────────────────────────────────────
interface SkillSectionProps {
  type: PowerType
  skills: Skill[]
  selected: Skill[]
  limit: number
  onToggle: (skill: Skill) => void
  isSelected: (skill: Skill) => boolean
  canSelect: (type: PowerType, skill: Skill) => boolean
}

function SkillSection({ type, skills, selected, limit, onToggle, isSelected, canSelect }: SkillSectionProps) {
  const cfg = powerConfig[type]
  if (skills.length === 0) return null

  return (
    <div>
      <div className="flex justify-between items-center mb-2">
        <h3 className={`text-sm font-bold ${cfg.color}`}>{cfg.label}</h3>
        <span className={`text-xs px-2 py-1 rounded-full ${selected.length >= limit ? cfg.badge : 'bg-gray-700 text-gray-400'}`}>
          {selected.length} / {limit}
        </span>
      </div>
      <div className="flex flex-col gap-2">
        {skills.map(skill => {
          const sel = isSelected(skill)
          const can = canSelect(type, skill)
          return (
            <button
              key={skill.ID}
              type="button"
              onClick={() => onToggle(skill)}
              disabled={!can && !sel}
              className={`text-left rounded-lg border p-3 transition-all ${
                sel
                  ? `${cfg.border} ${cfg.bg} ring-2 ring-offset-1 ring-offset-gray-800`
                  : can
                  ? 'border-gray-600 bg-gray-700 hover:border-gray-500 hover:bg-gray-600'
                  : 'border-gray-700 bg-gray-800 opacity-40 cursor-not-allowed'
              }`}
            >
              <div className="flex justify-between items-center">
                <span className="text-white text-sm font-semibold">{skill.name}</span>
                <span className={`text-xs w-5 h-5 rounded-full flex items-center justify-center border ${sel ? `${cfg.badge} border-transparent` : 'border-gray-500 text-gray-500'}`}>
                  {sel ? '✓' : ''}
                </span>
              </div>
              {skill.description && (
                <p className="text-gray-400 text-xs mt-1 line-clamp-2">{skill.description}</p>
              )}
            </button>
          )
        })}
      </div>
    </div>
  )
}

// ── Página principal ─────────────────────────────────────────────────────────
export default function CharacterCreate() {
  const navigate = useNavigate()
  const [selectedEdition, setSelectedEdition] = useState<string | null>(null)
  const [selectedClass, setSelectedClass] = useState<number | null>(null)
  const [selectedRace, setSelectedRace] = useState<number | null>(null)
  const [selectedClassData, setSelectedClassData] = useState<Class | null>(null)
  const [selectedArmorData, setSelectedArmorData] = useState<Armor | null>(null)
  const [selectedSkills, setSelectedSkills] = useState<Record<PowerType, Skill[]>>({
    unlimited: [], encounter: [], daily: [], utility: [],
  })

  const { register, handleSubmit, setValue, reset, control, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema) as any,
    defaultValues: { strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10, hit_points: 10 },
  })

  const constitution = useWatch({ control, name: 'constitution' })

  useEffect(() => {
    if (!selectedClassData || !selectedEdition) return
    const conValue = Number(constitution)
    const newHP = selectedEdition === '4e'
      ? (selectedClassData.base_hp ?? 0) + conValue
      : (selectedClassData.hit_die ?? 0) + getConMod(conValue)
    setValue('hit_points', Math.max(1, newHP))
  }, [selectedClassData, constitution, selectedEdition, setValue])

  const { data: classes } = useQuery({ queryKey: ['classes', selectedEdition], queryFn: () => classService.getAll(selectedEdition!), enabled: !!selectedEdition })
  const { data: armors }  = useQuery({ queryKey: ['armors', selectedEdition],  queryFn: () => armorService.getByEdition(selectedEdition!), enabled: !!selectedEdition })
  const { data: races }   = useQuery({ queryKey: ['races', selectedEdition],   queryFn: () => raceService.getAll(selectedEdition!), enabled: !!selectedEdition })
  const { data: allSkills } = useQuery({
    queryKey: ['skills-filter', selectedClass, selectedRace],
    queryFn: () => skillService.getByFilter(selectedClass!, selectedRace ?? undefined),
    enabled: !!selectedClass,
  })

  const createMutation = useMutation({
    mutationFn: async (data: FormData) => {
      const character = await characterService.create(data)
      const allSelected = Object.values(selectedSkills).flat()
      for (const skill of allSelected) {
        await characterService.addSkill(character.ID, skill.ID)
      }
      return character
    },
    onSuccess: () => navigate('/characters'),
  })

  const handleEditionChange = (edition: string) => {
    setSelectedEdition(edition)
    setSelectedClass(null); setSelectedRace(null); setSelectedClassData(null)
    setSelectedSkills({ unlimited: [], encounter: [], daily: [], utility: [] })
    reset({ name: '', edition, class_id: 0, race_id: 0, hit_points: 10, strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10 })
    setValue('edition', edition)
  }

  const handleClassChange = (classId: number) => {
    setSelectedClass(classId)
    setValue('class_id', classId)
    setSelectedClassData(classes?.find(c => c.ID === classId) ?? null)
    setSelectedSkills({ unlimited: [], encounter: [], daily: [], utility: [] })
  }

  const canSelect = (type: PowerType, skill: Skill) => {
    const sel = selectedSkills[type]
    return !!sel.find(s => s.ID === skill.ID) || sel.length < SKILL_LIMITS[type]
  }

  const toggleSkill = (skill: Skill) => {
    const type = (skill.power_type as PowerType) ?? 'unlimited'
    setSelectedSkills(prev => {
      const cur = prev[type]
      if (cur.find(s => s.ID === skill.ID)) return { ...prev, [type]: cur.filter(s => s.ID !== skill.ID) }
      if (cur.length >= SKILL_LIMITS[type]) return prev
      return { ...prev, [type]: [...cur, skill] }
    })
  }

  const isSelected = (skill: Skill) => !!selectedSkills[(skill.power_type as PowerType) ?? 'unlimited'].find(s => s.ID === skill.ID)
  const totalSelected = Object.values(selectedSkills).flat().length

  // Filtra por tipo e nível compatível com a criação (nível 1)
  const skillsByType = (type: PowerType): Skill[] =>
    (allSkills ?? []).filter(s => s.power_type === type && (!s.level || s.level <= 1))

  const utilitySkills = (allSkills ?? []).filter(s => s.power_type === 'utility')

  const classInfo = () => {
    if (!selectedClassData) return null
    return selectedEdition === '4e'
      ? `PV base: ${selectedClassData.base_hp ?? '?'} + CON — ${selectedClassData.description}`
      : `Hit Die: d${selectedClassData.hit_die} — ${selectedClassData.description}`
  }

  const hpLabel = () => {
    if (!selectedClassData) return null
    const conValue = Number(constitution)
    const conMod = getConMod(conValue)
    return selectedEdition === '4e'
      ? `${selectedClassData.base_hp ?? '?'} + ${conValue} (CON) = ${Math.max(1, (selectedClassData.base_hp ?? 0) + conValue)} PV`
      : `d${selectedClassData.hit_die} + mod CON (${conMod >= 0 ? '+' : ''}${conMod})`
  }

  const attributes = [
    { label: 'Força', key: 'strength' }, { label: 'Destreza', key: 'dexterity' },
    { label: 'Constituição', key: 'constitution' }, { label: 'Inteligência', key: 'intelligence' },
    { label: 'Sabedoria', key: 'wisdom' }, { label: 'Carisma', key: 'charisma' },
  ] as const

  const is4e = selectedEdition === '4e'
  const step = (n: number) => is4e ? `Passo ${n}` : `Passo ${n - 1}`

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-2xl mx-auto">
        <button onClick={() => navigate('/characters')} className="text-gray-400 hover:text-white transition mb-6 block">← Voltar</button>
        <h1 className="text-3xl font-bold text-white mb-8">Criar Personagem</h1>

        {/* Passo 1 — Edição */}
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 mb-6">
          <h2 className="text-lg font-semibold text-white mb-2">Passo 1 — Escolha a Edição</h2>
          <p className="text-gray-400 text-sm mb-4">Selecione qual edição do D&D você vai jogar.</p>
          <div className="grid grid-cols-2 gap-3">
            {editions.map(edition => (
              <button key={edition} type="button" onClick={() => handleEditionChange(edition)}
                className={`py-4 rounded-lg font-semibold text-base transition border ${selectedEdition === edition ? 'bg-indigo-600 border-indigo-500 text-white' : 'bg-gray-700 border-gray-600 text-gray-300 hover:bg-gray-600'}`}>
                D&D {edition}
              </button>
            ))}
          </div>
        </div>

        {selectedEdition && (
          <form onSubmit={handleSubmit(data => createMutation.mutate(data))} className="flex flex-col gap-6">

            {/* Passo 2 — Nome */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <h2 className="text-lg font-semibold text-white mb-4">Passo 2 — Nome do Personagem</h2>
              <input {...register('name')} className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500" placeholder="Ex: Aragorn" />
              {errors.name && <p className="text-red-400 text-xs mt-1">{errors.name.message}</p>}
            </div>

            {/* Passo 3 — Classe e Raça */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 flex flex-col gap-4">
              <h2 className="text-lg font-semibold text-white">Passo 3 — Classe e Raça</h2>
              <div>
                <label className="text-gray-400 text-sm mb-1 block">Classe</label>
                <select {...register('class_id', { valueAsNumber: true })} onChange={e => handleClassChange(Number(e.target.value))}
                  className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500">
                  <option value="">Selecione a classe</option>
                  {classes?.map(c => <option key={c.ID} value={c.ID}>{c.name}</option>)}
                </select>
                {selectedClassData && <p className="text-gray-400 text-xs mt-1">{classInfo()}</p>}
                {errors.class_id && <p className="text-red-400 text-xs mt-1">Selecione uma classe</p>}
              </div>
              <div>
                <label className="text-gray-400 text-sm mb-1 block">Raça</label>
                <select {...register('race_id', { valueAsNumber: true })} onChange={e => { const v = Number(e.target.value); setValue('race_id', v); setSelectedRace(v) }}
                  className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500">
                  <option value="">Selecione a raça</option>
                  {races?.map(r => <option key={r.ID} value={r.ID}>{r.name}</option>)}
                </select>
                {errors.race_id && <p className="text-red-400 text-xs mt-1">Selecione uma raça</p>}
              </div>
            </div>

            {/* Passo 4 — Poderes (somente 4e) */}
            {is4e && selectedClass && allSkills && allSkills.length > 0 && (
              <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
                <div className="flex justify-between items-center mb-1">
                  <h2 className="text-lg font-semibold text-white">Passo 4 — Poderes</h2>
                  <span className="text-xs text-gray-400 bg-gray-700 px-3 py-1 rounded-full">{totalSelected} selecionado{totalSelected !== 1 ? 's' : ''}</span>
                </div>
                <p className="text-gray-400 text-sm mb-5">Escolha seus poderes iniciais conforme os limites de cada tipo.</p>
                <div className="flex flex-col gap-6">
                  <SkillSection type="unlimited" skills={skillsByType('unlimited')} selected={selectedSkills.unlimited} limit={SKILL_LIMITS.unlimited} onToggle={toggleSkill} isSelected={isSelected} canSelect={canSelect} />
                  <SkillSection type="encounter" skills={skillsByType('encounter')} selected={selectedSkills.encounter} limit={SKILL_LIMITS.encounter} onToggle={toggleSkill} isSelected={isSelected} canSelect={canSelect} />
                  <SkillSection type="daily"     skills={skillsByType('daily')}     selected={selectedSkills.daily}     limit={SKILL_LIMITS.daily}     onToggle={toggleSkill} isSelected={isSelected} canSelect={canSelect} />

                  {/* Utilitário — disponível apenas no nível 2+ */}
                  {utilitySkills.length > 0 && (
                    <div className="opacity-50">
                      <div className="flex justify-between items-center mb-1">
                        <h3 className="text-sm font-bold text-blue-400">🔧 Utilitário</h3>
                        <span className="text-xs bg-gray-700 text-gray-500 px-2 py-1 rounded-full">Disponível no nível 2</span>
                      </div>
                      <p className="text-gray-500 text-xs">Poderes utilitários são desbloqueados ao atingir o nível 2.</p>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Armadura */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <h2 className="text-lg font-semibold text-white mb-4">{step(5)} — Armadura</h2>
              <select {...register('armor_id')} className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                onChange={e => { const v = Number(e.target.value); setValue('armor_id', v); setSelectedArmorData(armors?.find((a: Armor) => a.ID === v) ?? null) }}>
                <option value="">Sem armadura</option>
                {armors?.filter((a: Armor) => a.armor_type !== 'shield').map((a: Armor) => (
                  <option key={a.ID} value={a.ID}>{a.name} (CA base {a.base_ac})</option>
                ))}
              </select>
              {selectedArmorData && <p className="text-gray-400 text-xs mt-2">{selectedArmorData.description}</p>}
            </div>

            {/* Atributos */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <h2 className="text-lg font-semibold text-white mb-1">{step(6)} — Atributos</h2>
              <p className="text-gray-400 text-sm mb-4">Alterar a <span className="text-indigo-400">Constituição</span> recalcula o HP automaticamente.</p>
              <div className="grid grid-cols-2 gap-4">
                {attributes.map(attr => (
                  <div key={attr.key}>
                    <label className={`text-sm mb-1 block ${attr.key === 'constitution' ? 'text-indigo-400 font-semibold' : 'text-gray-400'}`}>{attr.label}</label>
                    <input type="number" {...register(attr.key)} min={1} max={20}
                      className={`w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 ${attr.key === 'constitution' ? 'focus:ring-indigo-400 ring-1 ring-indigo-800' : 'focus:ring-indigo-500'}`} />
                  </div>
                ))}
              </div>
            </div>

            {/* HP */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <div className="flex justify-between items-center mb-2">
                <h2 className="text-lg font-semibold text-white">{step(7)} — Hit Points</h2>
                {selectedClassData && <span className="text-xs text-gray-400 bg-gray-700 px-3 py-1 rounded-full">{hpLabel()}</span>}
              </div>
              {selectedClassData && <p className="text-gray-400 text-xs mb-3">Calculado automaticamente. Você pode ajustar manualmente.</p>}
              <input type="number" {...register('hit_points')} min={1} className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500" />
              {errors.hit_points && <p className="text-red-400 text-xs mt-1">{errors.hit_points.message}</p>}
            </div>

            {/* Resumo de poderes */}
            {is4e && totalSelected > 0 && (
              <div className="bg-gray-800 rounded-lg p-4 border border-indigo-700">
                <p className="text-indigo-300 text-sm font-semibold mb-2">⚡ Poderes selecionados ({totalSelected})</p>
                <div className="flex flex-wrap gap-2">
                  {Object.values(selectedSkills).flat().map(skill => (
                    <span key={skill.ID} className={`text-xs px-2 py-1 rounded-full ${powerConfig[(skill.power_type as PowerType) ?? 'unlimited'].badge}`}>
                      {skill.name}
                    </span>
                  ))}
                </div>
              </div>
            )}

            <button type="submit" disabled={createMutation.isPending}
              className="bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white font-bold py-3 rounded-lg transition">
              {createMutation.isPending ? 'Criando...' : '✨ Criar Personagem'}
            </button>

          </form>
        )}
      </div>
    </div>
  )
}