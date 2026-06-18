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
import { periciaService } from '../services/periciaService'
import { talentoService } from '../services/talentoService'
import { antecedentService } from '../services/antecedentService'
import { SkillCard, powerConfig } from '../components/SkillCard'
import { Tooltip } from '../components/Tooltip'
// FIX 1: import único, Background removido, Antecedent adicionado
import type { Class, Armor, Skill, PowerType, Race, Talento, Antecedent } from '../types'

const schema = z.object({
  name:               z.string().min(1, 'Nome é obrigatório'),
  edition:            z.string().min(1, 'Edição é obrigatória'),
  class_id:           z.coerce.number().min(1, 'Selecione uma classe'),
  race_id:            z.coerce.number().min(1, 'Selecione uma raça'),
  hit_points:         z.coerce.number().min(1),
  strength:           z.coerce.number().min(1).max(20),
  dexterity:          z.coerce.number().min(1).max(20),
  constitution:       z.coerce.number().min(1).max(20),
  intelligence:       z.coerce.number().min(1).max(20),
  wisdom:             z.coerce.number().min(1).max(20),
  charisma:           z.coerce.number().min(1).max(20),
  armor_id:           z.coerce.number().optional(),
  // FIX 4a: background_id → antecedent_id (bate com json:"antecedent_id" no Go)
  antecedent_id:      z.coerce.number().optional(),
  alignment:          z.string().optional(),
  personality_traits: z.string().optional(),
  ideals:             z.string().optional(),
  bonds:              z.string().optional(),
  flaws:              z.string().optional(),
})
type FormData = z.infer<typeof schema>

const SKILL_LIMITS: Record<PowerType, number> = {
  unlimited: 2, encounter: 1, daily: 1, utility: 1,
}

const EDITIONS = [
  { value: '4e', label: 'D&D 4e', disabled: false },
  { value: '5e', label: 'D&D 5e', disabled: false },
]

const CATEGORY_CONFIG: Record<string, { color: string; icon: string }> = {
  'Combate':  { color: 'text-red-400',    icon: '⚔️' },
  'Defesa':   { color: 'text-blue-400',   icon: '🛡️' },
  'Perícia':  { color: 'text-yellow-400', icon: '📚' },
  'Magia':    { color: 'text-purple-400', icon: '✨' },
  'Armadura': { color: 'text-gray-300',   icon: '🪖' },
}

const ALIGNMENTS = [
  'Leal e Bom',    'Neutro e Bom',    'Caótico e Bom',
  'Leal e Neutro', 'Neutro',          'Caótico e Neutro',
  'Leal e Mau',    'Neutro e Mau',    'Caótico e Mau',
]

function getConMod(con: number) { return Math.floor((con - 10) / 2) }

export default function CharacterCreate() {
  const navigate = useNavigate()

  const [selectedEdition, setSelectedEdition] = useState<string | null>(() =>
    sessionStorage.getItem('rpg_selected_edition')
  )
  const [selectedClass,     setSelectedClass]     = useState<number | null>(null)
  const [selectedRace,      setSelectedRace]      = useState<number | null>(null)
  const [selectedClassData, setSelectedClassData] = useState<Class | null>(null)
  const [selectedRaceData,  setSelectedRaceData]  = useState<Race | null>(null)
  const [selectedArmorData, setSelectedArmorData] = useState<Armor | null>(null)
  const [selectedBackground, setSelectedBackground] = useState<Antecedent | null>(null)
  const [selectedAlignment,  setSelectedAlignment]  = useState<string>('')

  const [selectedSkills, setSelectedSkills] = useState<Record<PowerType, Skill[]>>({
    unlimited: [], encounter: [], daily: [], utility: [],
  })
  const [choiceSelections, setChoiceSelections] = useState<Record<string, Skill>>({})
  const [selectedPericias, setSelectedPericias] = useState<string[]>([])
  const [selectedTalentos, setSelectedTalentos] = useState<Talento[]>([])

  const {
    register, handleSubmit, setValue, reset, control, formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema) as any,
    defaultValues: {
      strength: 10, dexterity: 10, constitution: 10,
      intelligence: 10, wisdom: 10, charisma: 10, hit_points: 10,
    },
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

  const { data: classes }   = useQuery({ queryKey: ['classes', selectedEdition],  queryFn: () => classService.getAll(selectedEdition!),      enabled: !!selectedEdition })
  const { data: races }     = useQuery({ queryKey: ['races', selectedEdition],    queryFn: () => raceService.getAll(selectedEdition!),        enabled: !!selectedEdition })
  const { data: armors }    = useQuery({ queryKey: ['armors', selectedEdition],   queryFn: () => armorService.getByEdition(selectedEdition!), enabled: !!selectedEdition })
  const { data: allSkills } = useQuery({
    queryKey: ['skills-filter', selectedClass, selectedRace],
    queryFn:  () => skillService.getByFilter(selectedClass!, selectedRace ?? undefined),
    enabled:  !!selectedClass, staleTime: 0,
  })
  const { data: allPericias } = useQuery({
    queryKey: ['pericias', selectedEdition],
    queryFn:  () => periciaService.getAll(selectedEdition!),
    enabled:  !!selectedEdition, staleTime: Infinity,
  })
  const { data: allTalentos } = useQuery({
    queryKey: ['talentos', selectedEdition],
    queryFn:  () => talentoService.getAll(selectedEdition!),
    enabled:  !!selectedEdition && selectedEdition === '4e', staleTime: Infinity,
  })
  const { data: allBackgrounds } = useQuery({
    queryKey: ['antecedentes', selectedEdition],
    queryFn:  () => antecedentService.getAll(selectedEdition!),
    enabled:  !!selectedEdition && selectedEdition === '5e', staleTime: Infinity,
  })

  const is4e = selectedEdition === '4e'
  const is5e = selectedEdition === '5e'
  const hasSkills = allSkills && allSkills.length > 0

  const availableSkills: string[] = (() => {
    if (!selectedClassData?.available_skills) return []
    try { return JSON.parse(selectedClassData.available_skills) } catch { return [] }
  })()

  const savingThrows: string[] = (() => {
    if (!selectedClassData?.saving_throws) return []
    try { return JSON.parse(selectedClassData.saving_throws) } catch { return [] }
  })()

  const backgroundSkills: string[] = (() => {
    if (!selectedBackground?.skill_proficiencies) return []
    try { return JSON.parse(selectedBackground.skill_proficiencies) } catch { return [] }
  })()

  const totalTrainable = (selectedClassData?.trained_skills_count ?? 0) + (selectedRaceData?.bonus_trained_skills ?? 0)

  const bonusSkillValues: Record<string, number> = (() => {
    if (!selectedRaceData?.bonus_skill_values) return {}
    try { return JSON.parse(selectedRaceData.bonus_skill_values) } catch { return {} }
  })()

  const totalTalentos = (selectedClassData?.talentos_count ?? 2) + (selectedRaceData?.bonus_talentos ?? 0)

  const classFeatures   = (allSkills ?? []).filter(s => s.is_class_feature && !s.requires_choice)
  const choiceFeatures  = (allSkills ?? []).filter(s => s.is_class_feature && s.requires_choice)
  const choiceGroups    = choiceFeatures.reduce<Record<string, Skill[]>>((acc, s) => {
    const g = s.choice_group ?? 'default'; acc[g] = [...(acc[g] ?? []), s]; return acc
  }, {})
  const raceFeatures       = (allSkills ?? []).filter(s => s.is_race_feature && !s.requires_choice)
  const raceChoiceFeatures = (allSkills ?? []).filter(s => s.is_race_feature && s.requires_choice)
  const raceChoiceGroups   = raceChoiceFeatures.reduce<Record<string, Skill[]>>((acc, s) => {
    const g = s.choice_group ?? 'default'; acc[g] = [...(acc[g] ?? []), s]; return acc
  }, {})

  const hasAnyChoiceGroups = Object.keys(choiceGroups).length > 0 || Object.keys(raceChoiceGroups).length > 0
  const allChoicesMade =
    Object.keys(choiceGroups).every(g => !!choiceSelections[g]) &&
    Object.keys(raceChoiceGroups).every(g => !!choiceSelections[g])
  const hasAnyFeatures =
    classFeatures.length > 0 || Object.keys(choiceGroups).length > 0 ||
    raceFeatures.length > 0  || Object.keys(raceChoiceGroups).length > 0

  const normalByType = (type: PowerType): Skill[] =>
    (allSkills ?? []).filter(s =>
      !s.is_class_feature && !s.is_race_feature && s.power_type === type && (!s.level || s.level <= 1)
    )
  const utilitySkills = (allSkills ?? []).filter(s =>
    !s.is_class_feature && !s.is_race_feature && s.power_type === 'utility'
  )

  const totalNormal = Object.values(selectedSkills).flat().length
  const totalSelected = totalNormal + Object.keys(choiceSelections).length
  const pericias5eDisponiveis = (allPericias ?? []).filter(p =>
    availableSkills.includes(p.name) && !backgroundSkills.includes(p.name)
  )

  const pendingItems: string[] = []
  if (is4e) {
    if (!selectedClass)  pendingItems.push('Selecione uma classe')
    if (!selectedRace)   pendingItems.push('Selecione uma raça')
    if (selectedClass && hasAnyChoiceGroups && !allChoicesMade)
      pendingItems.push('Complete as escolhas de características')
    if (selectedClass && totalTrainable > 0 && selectedPericias.length < totalTrainable)
      pendingItems.push(`Escolha mais ${totalTrainable - selectedPericias.length} perícia(s) treinada(s)`)
    if (selectedClass && totalTalentos > 0 && selectedTalentos.length < totalTalentos)
      pendingItems.push(`Escolha mais ${totalTalentos - selectedTalentos.length} talento(s)`)
  }
  if (is5e) {
    if (!selectedClass)      pendingItems.push('Selecione uma classe')
    if (!selectedRace)       pendingItems.push('Selecione uma raça')
    if (!selectedBackground) pendingItems.push('Selecione um antecedente')
    if (selectedClass && hasAnyChoiceGroups && !allChoicesMade)
      pendingItems.push('Complete as escolhas de características')
    if (selectedClass && pericias5eDisponiveis.length > 0) {
      const classChoicesMade = selectedPericias.filter(p => !backgroundSkills.includes(p)).length
      if (classChoicesMade < totalTrainable)
        pendingItems.push(`Escolha mais ${totalTrainable - classChoicesMade} perícia(s) da classe`)
    }
  }

  const createMutation = useMutation({
    mutationFn: async (data: FormData) => {
      const characterData = {
        ...data,
        // FIX 4b: background_id → antecedent_id para bater com o backend Go
        antecedent_id: selectedBackground?.ID ?? undefined,
        alignment:     selectedAlignment || undefined,
      }
      const character = await characterService.create(characterData)
      for (const skill of Object.values(selectedSkills).flat())
        await characterService.addSkill(character.ID, skill.ID)
      for (const skill of Object.values(choiceSelections))
        await characterService.addSkill(character.ID, skill.ID)
      for (const skill of classFeatures)
        await characterService.addSkill(character.ID, skill.ID)
      for (const skill of raceFeatures)
        await characterService.addSkill(character.ID, skill.ID)
      if (selectedPericias.length > 0)
        await periciaService.save(character.ID, selectedPericias)
      for (const talento of selectedTalentos)
        await talentoService.add(character.ID, talento.ID)
      return character
    },
    onSuccess: () => navigate('/characters'),
  })

  const canSubmit = !createMutation.isPending && pendingItems.length === 0

  const handleEditionChange = (edition: string) => {
    sessionStorage.setItem('rpg_selected_edition', edition)
    setSelectedEdition(edition)
    setSelectedClass(null); setSelectedRace(null)
    setSelectedClassData(null); setSelectedRaceData(null); setSelectedArmorData(null)
    setSelectedBackground(null); setSelectedAlignment('')
    setSelectedSkills({ unlimited: [], encounter: [], daily: [], utility: [] })
    setChoiceSelections({})
    setSelectedPericias([]); setSelectedTalentos([])
    reset({ name: '', edition, class_id: 0, race_id: 0, hit_points: 10,
            strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10 })
    setValue('edition', edition)
  }
  const handleClassChange = (classId: number) => {
    setSelectedClass(classId)
    setValue('class_id', classId)
    setSelectedClassData(classes?.find(c => c.ID === classId) ?? null)
    setSelectedSkills({ unlimited: [], encounter: [], daily: [], utility: [] })
    setChoiceSelections({})
    setSelectedPericias(backgroundSkills)
    setSelectedTalentos([])
  }
  const handleRaceChange = (raceId: number) => {
    setValue('race_id', raceId)
    setSelectedRace(raceId || null)
    setSelectedRaceData(races?.find(r => r.ID === raceId) ?? null)
  }
  // FIX 4c: parâmetro tipado como Antecedent (não Background)
  const handleBackgroundChange = (bg: Antecedent) => {
    const newBgSkills: string[] = (() => {
      try { return JSON.parse(bg.skill_proficiencies) } catch { return [] }
    })()
    setSelectedBackground(bg)
    // FIX 4d: setValue usa antecedent_id (não background_id)
    setValue('antecedent_id', bg.ID)
    setSelectedPericias(newBgSkills)
  }

  const canSelect    = (type: PowerType, skill: Skill) => {
    const sel = selectedSkills[type]
    return !!sel.find(s => s.ID === skill.ID) || sel.length < SKILL_LIMITS[type]
  }
  const toggleSkill  = (skill: Skill) => {
    const type = (skill.power_type as PowerType) ?? 'unlimited'
    setSelectedSkills(prev => {
      const cur = prev[type]
      if (cur.find(s => s.ID === skill.ID)) return { ...prev, [type]: cur.filter(s => s.ID !== skill.ID) }
      if (cur.length >= SKILL_LIMITS[type]) return prev
      return { ...prev, [type]: [...cur, skill] }
    })
  }
  const isSelected       = (skill: Skill) => !!selectedSkills[(skill.power_type as PowerType) ?? 'unlimited'].find(s => s.ID === skill.ID)
  const selectChoice     = (skill: Skill) => setChoiceSelections(prev => ({ ...prev, [skill.choice_group ?? 'default']: skill }))
  const isChoiceSelected = (skill: Skill) => choiceSelections[skill.choice_group ?? 'default']?.ID === skill.ID

  const togglePericia = (name: string) => {
    if (backgroundSkills.includes(name)) return
    setSelectedPericias(prev => {
      if (prev.includes(name)) return prev.filter(p => p !== name)
      const classChoices = prev.filter(p => !backgroundSkills.includes(p))
      if (classChoices.length >= totalTrainable) return prev
      return [...prev, name]
    })
  }
  const toggleTalento = (talento: Talento) => {
    setSelectedTalentos(prev => {
      if (prev.some(t => t.ID === talento.ID)) return prev.filter(t => t.ID !== talento.ID)
      if (prev.length >= totalTalentos) return prev
      return [...prev, talento]
    })
  }

  // FIX 2: classInfo agora tem a chave } de fechamento correta
  const classInfo = () => {
    if (!selectedClassData) return null
    if (is4e) return `PV base: ${selectedClassData.base_hp ?? '?'} + CON — ${selectedClassData.description}`
    const saves = savingThrows.length > 0 ? savingThrows.join(', ') : '—'
    return `Hit Die: d${selectedClassData.hit_die ?? '?'} · Salvaguardas: ${saves} — ${selectedClassData.description}`
  } // ← CHAVE QUE ESTAVA FALTANDO

  // FIX 3: hpLabel com template literals corretos (sem \`)
  const hpLabel = () => {
    if (!selectedClassData) return null
    const conValue = Number(constitution)
    const conMod = getConMod(conValue)
    return is4e
      ? `${selectedClassData.base_hp ?? '?'} + ${conValue} CON = ${Math.max(1, (selectedClassData.base_hp ?? 0) + conValue)} PV`
      : `d${selectedClassData.hit_die ?? '?'} + mod CON (${conMod >= 0 ? '+' : ''}${conMod})`
  }

  const attrFields = [
    { label: 'Força',        key: 'strength'     },
    { label: 'Destreza',     key: 'dexterity'    },
    { label: 'Constituição', key: 'constitution' },
    { label: 'Inteligência', key: 'intelligence' },
    { label: 'Sabedoria',    key: 'wisdom'       },
    { label: 'Carisma',      key: 'charisma'     },
  ] as const

  let step = 1
  const S = () => step++

  const sectionHeader = (label: string) => (
    <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
      {label}
    </h2>
  )

  const counterBadge = (current: number, total: number) => {
    const done = current >= total && total > 0
    return (
      <span className="text-xs px-2 py-1 rounded-full transition"
        style={done
          ? { background: 'rgba(201,168,76,0.15)', border: '1px solid rgba(201,168,76,0.4)', color: '#c9a84c' }
          : { background: 'rgba(255,255,255,0.05)', border: '1px solid #3f3f46', color: '#71717a' }
        }
      >{current} / {total}</span>
    )
  }
  // Sincroniza a edition do sessionStorage com o react-hook-form no mount
  useEffect(() => {
    if (selectedEdition) {
      setValue('edition', selectedEdition)
    }
  }, [selectedEdition, setValue])
  
  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-2xl mx-auto">
        <button onClick={() => navigate('/characters')} className="transition mb-6 block text-sm"
          style={{ color: 'rgba(201,168,76,0.5)' }}
          onMouseEnter={e => (e.currentTarget.style.color = '#c9a84c')}
          onMouseLeave={e => (e.currentTarget.style.color = 'rgba(201,168,76,0.5)')}
        >← Voltar</button>

        <h1 className="font-rpg text-2xl sm:text-3xl font-bold mb-6 sm:mb-8" style={{ color: '#c9a84c' }}>
          Criar Personagem
        </h1>

        {/* PASSO 1 — Edição */}
        <div className="bg-gray-800 rounded-xl p-5 border border-gray-700 mb-5">
          {sectionHeader(`Passo ${S()} — Edição`)}
          <div className="grid grid-cols-2 gap-3">
            {EDITIONS.map(e => (
              <button key={e.value} type="button" disabled={e.disabled}
                onClick={() => !e.disabled && handleEditionChange(e.value)}
                className="py-4 rounded-lg font-semibold text-sm transition border"
                style={selectedEdition === e.value
                  ? { background: '#c9a84c', borderColor: '#c9a84c', color: '#0a0a0a' }
                  : e.disabled
                  ? { background: 'transparent', borderColor: '#3f3f46', color: '#52525b', cursor: 'not-allowed' }
                  : { background: '#3f3f46', borderColor: '#52525b', color: '#d4d4d8' }
                }
              >{e.label}</button>
            ))}
          </div>
        </div>

        {selectedEdition && (
          <form onSubmit={handleSubmit(
            data => createMutation.mutate(data),
            (errors) => console.log('Erros de validação:', errors)
          )} className="flex flex-col gap-5">

            {/* PASSO 2 — Nome */}
            <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
              {sectionHeader(`Passo ${S()} — Nome do Personagem`)}
              <input {...register('name')} className="rpg-input" placeholder="Ex: Aragorn" />
              {errors.name && <p className="text-red-400 text-xs mt-1">{errors.name.message}</p>}
            </div>

            {/* PASSO 3 — Classe e Raça */}
            <div className="bg-gray-800 rounded-xl p-5 border border-gray-700 flex flex-col gap-4">
              {sectionHeader(`Passo ${S()} — Classe e Raça`)}
              <div>
                <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">Classe</label>
                <select {...register('class_id', { valueAsNumber: true })}
                  onChange={e => handleClassChange(Number(e.target.value))} className="rpg-select">
                  <option value="">Selecione a classe</option>
                  {classes?.map(c => <option key={c.ID} value={c.ID}>{c.name}</option>)}
                </select>
                {selectedClassData && <p className="text-gray-500 text-xs mt-1">{classInfo()}</p>}
                {errors.class_id && <p className="text-red-400 text-xs mt-1">Selecione uma classe</p>}
              </div>
              <div>
                <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">Raça / Espécie</label>
                <select {...register('race_id', { valueAsNumber: true })}
                  onChange={e => handleRaceChange(Number(e.target.value))} className="rpg-select">
                  <option value="">Selecione a raça</option>
                  {races?.map(r => <option key={r.ID} value={r.ID}>{r.name}</option>)}
                </select>
                {errors.race_id && <p className="text-red-400 text-xs mt-1">Selecione uma raça</p>}
              </div>
              {is5e && savingThrows.length > 0 && (
                <div className="flex items-center gap-2 flex-wrap pt-1">
                  <span className="text-gray-500 text-xs">Salvaguardas:</span>
                  {savingThrows.map(s => (
                    <span key={s} className="text-xs px-2 py-0.5 rounded-full"
                      style={{ background: 'rgba(201,168,76,0.1)', border: '1px solid rgba(201,168,76,0.3)', color: '#c9a84c' }}>
                      {s}
                    </span>
                  ))}
                </div>
              )}
            </div>

            {/* PASSO 4 — Antecedente (só 5e) */}
            {is5e && allBackgrounds && allBackgrounds.length > 0 && (
              <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                {sectionHeader(`Passo ${S()} — Antecedente`)}
                <p className="text-gray-500 text-xs mb-4">
                  O antecedente revela de onde você vem. Concede{' '}
                  <span style={{ color: '#c9a84c' }}>2 proficiências em perícias automáticas</span>,
                  equipamento inicial e uma característica especial narrativa.
                </p>
                <div className="flex flex-col gap-3">
                  {allBackgrounds.map(bg => {
                    const bgSkills: string[] = (() => { try { return JSON.parse(bg.skill_proficiencies) } catch { return [] } })()
                    const isSel = selectedBackground?.ID === bg.ID
                    return (
                      <button key={bg.ID} type="button" onClick={() => handleBackgroundChange(bg)}
                        className="text-left rounded-xl p-4 border transition"
                        style={isSel
                          ? { background: 'rgba(201,168,76,0.08)', borderColor: 'rgba(201,168,76,0.5)' }
                          : { background: '#1a1a1a', borderColor: '#3f3f46' }
                        }
                      >
                        <div className="flex items-start justify-between gap-3">
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-1 flex-wrap">
                              <span className="font-semibold text-sm" style={{ color: isSel ? '#c9a84c' : '#e4e4e7' }}>{bg.name}</span>
                              {isSel && <span className="text-xs bg-yellow-900/60 text-yellow-300 px-2 py-0.5 rounded-full">✓ Selecionado</span>}
                            </div>
                            <p className="text-gray-400 text-xs mb-2">{bg.description}</p>
                            <div className="flex flex-wrap gap-1.5">
                              {bgSkills.map(s => <span key={s} className="text-xs bg-indigo-900/60 text-indigo-300 px-2 py-0.5 rounded-full">📚 {s}</span>)}
                              {bg.tool_proficiencies && <span className="text-xs bg-orange-900/60 text-orange-300 px-2 py-0.5 rounded-full">🔧 {bg.tool_proficiencies}</span>}
                              {bg.languages && <span className="text-xs bg-teal-900/60 text-teal-300 px-2 py-0.5 rounded-full">🗣 {bg.languages}</span>}
                            </div>
                          </div>
                          <div className="w-5 h-5 rounded-full border-2 flex-shrink-0 mt-0.5 flex items-center justify-center transition"
                            style={isSel ? { background: '#c9a84c', borderColor: '#c9a84c' } : { borderColor: '#52525b' }}>
                            {isSel && <span className="text-black text-xs font-bold">✓</span>}
                          </div>
                        </div>
                        {isSel && (
                          <div className="mt-3 pt-3 border-t" style={{ borderColor: 'rgba(201,168,76,0.2)' }}>
                            <p className="text-xs font-semibold mb-1" style={{ color: '#c9a84c' }}>✦ {bg.feature}</p>
                            <p className="text-gray-400 text-xs mb-2">{bg.feature_description}</p>
                            <p className="text-gray-500 text-xs"><span className="text-gray-400 font-medium">Equipamento: </span>{bg.equipment}</p>
                          </div>
                        )}
                      </button>
                    )
                  })}
                </div>
              </div>
            )}

            {/* PASSO 5 — Personalidade (só 5e) */}
            {is5e && selectedBackground && (
              <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                {sectionHeader(`Passo ${S()} — Tendência e Personalidade`)}
                <p className="text-gray-500 text-xs mb-4">Campos opcionais que definem como seu personagem pensa e age.</p>
                <div className="mb-5">
                  <label className="text-gray-500 text-xs mb-2 block uppercase tracking-wider">Tendência</label>
                  <div className="grid grid-cols-3 gap-1.5">
                    {ALIGNMENTS.map(al => (
                      <button key={al} type="button" onClick={() => setSelectedAlignment(al === selectedAlignment ? '' : al)}
                        className="py-2 rounded-lg text-xs font-medium transition border"
                        style={selectedAlignment === al
                          ? { background: 'rgba(201,168,76,0.15)', borderColor: 'rgba(201,168,76,0.5)', color: '#c9a84c' }
                          : { background: '#27272a', borderColor: '#3f3f46', color: '#a1a1aa' }
                        }
                      >{al}</button>
                    ))}
                  </div>
                </div>
                {[
                  { label: 'Traços de Personalidade', name: 'personality_traits' as const, placeholder: 'Como seu personagem age e se comporta no cotidiano?' },
                  { label: 'Ideais',   name: 'ideals' as const,  placeholder: 'O que move seu personagem? Quais princípios ele nunca abandona?' },
                  { label: 'Ligações', name: 'bonds' as const,   placeholder: 'Quem ou o que seu personagem protegeria com a própria vida?' },
                  { label: 'Defeitos', name: 'flaws' as const,   placeholder: 'Qual é a fraqueza, vício ou medo que pode atrapalhar seu personagem?' },
                ].map(field => (
                  <div key={field.name} className="mb-3">
                    <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">{field.label}</label>
                    <textarea {...register(field.name)} placeholder={field.placeholder} rows={2}
                      className="rpg-input resize-none text-sm w-full" style={{ minHeight: '56px' }} />
                  </div>
                ))}
              </div>
            )}

            {/* PASSO 6 — Características (4e e 5e) */}
            {selectedClass && hasSkills && hasAnyFeatures && (
              <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                {sectionHeader(`Passo ${S()} — Características`)}
                <div className="flex flex-col gap-6">
                  {classFeatures.length > 0 && (
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-sm font-bold text-indigo-400">📖 Classe — Automáticas</span>
                        <span className="text-xs bg-indigo-900/60 text-indigo-300 px-2 py-0.5 rounded-full">Todos possuem</span>
                      </div>
                      <p className="text-gray-500 text-xs mb-3">Concedidas automaticamente pela classe.</p>
                      <div className="flex flex-col gap-2">
                        {classFeatures.map(s => <SkillCard key={s.ID} skill={s} informative defaultExpanded />)}
                      </div>
                    </div>
                  )}
                  {Object.entries(choiceGroups).map(([group, options]) => {
                    const chosen = choiceSelections[group]
                    return (
                      <div key={group}>
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-bold text-purple-400">🎯 Classe — Escolha Obrigatória</span>
                          {chosen ? <span className="text-xs bg-purple-900/60 text-purple-300 px-2 py-0.5 rounded-full">✓ {chosen.name}</span>
                                  : <span className="text-xs bg-red-900/60 text-red-300 px-2 py-0.5 rounded-full animate-pulse">Escolha uma opção</span>}
                        </div>
                        <p className="text-gray-500 text-xs mb-3">Escolha <strong className="text-white">uma</strong> das opções — permanente.</p>
                        <div className="flex flex-col gap-2">
                          {options.map(s => <SkillCard key={s.ID} skill={s} selectable selected={isChoiceSelected(s)} disabled={false} onToggle={selectChoice} defaultExpanded />)}
                        </div>
                      </div>
                    )
                  })}
                  {(classFeatures.length > 0 || Object.keys(choiceGroups).length > 0) &&
                   (raceFeatures.length > 0 || Object.keys(raceChoiceGroups).length > 0) && (
                    <hr style={{ borderColor: '#3f3f46' }} />
                  )}
                  {raceFeatures.length > 0 && (
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-sm font-bold text-emerald-400">🐉 Raça — Automáticas</span>
                        <span className="text-xs bg-emerald-900/60 text-emerald-300 px-2 py-0.5 rounded-full">Todos possuem</span>
                      </div>
                      <p className="text-gray-500 text-xs mb-3">Concedidas automaticamente pela raça.</p>
                      <div className="flex flex-col gap-2">
                        {raceFeatures.map(s => <SkillCard key={s.ID} skill={s} informative defaultExpanded />)}
                      </div>
                    </div>
                  )}
                  {Object.entries(raceChoiceGroups).map(([group, options]) => {
                    const chosen = choiceSelections[group]
                    return (
                      <div key={group}>
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-bold text-amber-400">🎯 Raça — Escolha Obrigatória</span>
                          {chosen ? <span className="text-xs bg-amber-900/60 text-amber-300 px-2 py-0.5 rounded-full">✓ {chosen.name}</span>
                                  : <span className="text-xs bg-red-900/60 text-red-300 px-2 py-0.5 rounded-full animate-pulse">Escolha uma opção</span>}
                        </div>
                        <p className="text-gray-500 text-xs mb-3">Escolha <strong className="text-white">uma</strong> das opções — permanente.</p>
                        <div className="flex flex-col gap-2">
                          {options.map(s => <SkillCard key={s.ID} skill={s} selectable selected={isChoiceSelected(s)} disabled={false} onToggle={selectChoice} defaultExpanded />)}
                        </div>
                      </div>
                    )
                  })}
                </div>
              </div>
            )}

            {/* PASSO 7 — Perícias (4e e 5e) */}
            {(is4e || is5e) && selectedClass && allPericias && allPericias.length > 0 && (
              <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                <div className="flex items-center justify-between mb-1">
                  {sectionHeader(`Passo ${S()} — Perícias Treinadas`)}
                  {is5e
                    ? counterBadge(selectedPericias.filter(p => !backgroundSkills.includes(p)).length, totalTrainable)
                    : counterBadge(selectedPericias.length, totalTrainable)
                  }
                </div>
                <p className="text-gray-500 text-xs mb-4">
                  {is4e
                    ? <><span className="text-yellow-400">⭐</span> = disponíveis pela sua classe. Perícias treinadas recebem <span className="text-white">+5</span> nas jogadas.</>
                    : <>Perícias <span className="text-indigo-400">🔒 fixas</span> vêm do antecedente. Escolha mais <span className="text-white">{totalTrainable}</span> da lista da classe.</>
                  }
                </p>

                {is5e && backgroundSkills.length > 0 && (
                  <div className="mb-4 rounded-lg p-3" style={{ background: 'rgba(99,102,241,0.08)', border: '1px solid rgba(99,102,241,0.25)' }}>
                    <p className="text-indigo-400 text-xs font-semibold mb-2">🔒 Antecedente — Automáticas</p>
                    <div className="flex flex-wrap gap-2">
                      {backgroundSkills.map(s => <span key={s} className="text-xs bg-indigo-900/60 text-indigo-300 px-2 py-1 rounded-full">📚 {s}</span>)}
                    </div>
                  </div>
                )}

                {Object.keys(bonusSkillValues).length > 0 && (
                  <div className="mb-4 rounded-lg p-3" style={{ background: 'rgba(16,185,129,0.08)', border: '1px solid rgba(16,185,129,0.25)' }}>
                    <p className="text-emerald-400 text-xs font-semibold mb-2">🐉 Bônus Racial Automático</p>
                    <div className="flex flex-wrap gap-2">
                      {Object.entries(bonusSkillValues).map(([skill, bonus]) => (
                        <span key={skill} className="text-xs bg-emerald-900/60 text-emerald-300 px-2 py-0.5 rounded-full">{skill} +{bonus}</span>
                      ))}
                    </div>
                  </div>
                )}

                <div className="flex flex-col gap-2 mb-4">
                  {allPericias
                    .filter(p => availableSkills.includes(p.name) && !backgroundSkills.includes(p.name))
                    .map(p => {
                      const isPerSelected = selectedPericias.includes(p.name)
                      const classChoicesMade = selectedPericias.filter(x => !backgroundSkills.includes(x)).length
                      const isDisabled = !isPerSelected && classChoicesMade >= totalTrainable
                      const racialBonus = bonusSkillValues[p.name]
                      return (
                        <button key={p.ID} type="button" disabled={isDisabled} onClick={() => togglePericia(p.name)}
                          className="flex items-center justify-between rounded-lg px-4 py-3 border text-left transition"
                          style={isPerSelected
                            ? { background: 'rgba(201,168,76,0.1)', borderColor: 'rgba(201,168,76,0.5)' }
                            : isDisabled
                            ? { background: '#1a1a1a', borderColor: '#3f3f46', opacity: 0.4, cursor: 'not-allowed' }
                            : { background: '#27272a', borderColor: '#3f3f46' }
                          }
                        >
                          <div className="flex items-center gap-2 flex-1 min-w-0">
                            <span className="text-yellow-400 text-xs flex-shrink-0">⭐</span>
                            <div className="min-w-0">
                              <div className="flex items-center gap-2 flex-wrap">
                                <span className="font-medium text-sm text-white">{p.name}</span>
                                <span className="text-xs text-gray-500">({p.attribute})</span>
                                {racialBonus && <span className="text-xs bg-emerald-900/60 text-emerald-300 px-1.5 py-0.5 rounded">+{racialBonus} racial</span>}
                              </div>
                              <p className="text-gray-500 text-xs mt-0.5 truncate">{p.description}</p>
                            </div>
                          </div>
                          <div className="flex items-center gap-2 ml-3 flex-shrink-0">
                            <Tooltip content={p.tooltip} />
                            <div className="w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 transition"
                              style={isPerSelected ? { background: '#c9a84c', borderColor: '#c9a84c' } : { borderColor: '#52525b' }}>
                              {isPerSelected && <span className="text-black text-xs leading-none font-bold">✓</span>}
                            </div>
                          </div>
                        </button>
                      )
                    })}
                </div>

                {allPericias.filter(p => !availableSkills.includes(p.name) && !backgroundSkills.includes(p.name)).length > 0 && (
                  <div>
                    <p className="text-gray-600 text-xs mb-2">🔒 Não disponíveis para esta classe:</p>
                    <div className="flex flex-wrap gap-2">
                      {allPericias.filter(p => !availableSkills.includes(p.name) && !backgroundSkills.includes(p.name)).map(p => (
                        <div key={p.ID} className="flex items-center gap-1">
                          <span className="text-xs bg-gray-800 text-gray-600 px-2 py-1 rounded border border-gray-700">{p.name}</span>
                          <Tooltip content={p.tooltip} />
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}

            {/* PASSO 8 — Talentos (só 4e) */}
            {is4e && selectedClass && allTalentos && allTalentos.length > 0 && (
              <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                <div className="flex items-center justify-between mb-1">
                  {sectionHeader(`Passo ${S()} — Talentos`)}
                  {counterBadge(selectedTalentos.length, totalTalentos)}
                </div>
                <p className="text-gray-500 text-xs mb-5">
                  Escolha {totalTalentos} talento{totalTalentos !== 1 ? 's' : ''}.
                  {selectedRaceData?.bonus_talentos ? ` Sua raça concede +${selectedRaceData.bonus_talentos} extra.` : ''}
                </p>
                <div className="flex flex-col gap-6">
                  {Object.entries(CATEGORY_CONFIG).map(([category, cfg]) => {
                    const items = allTalentos.filter(t => t.category === category)
                    if (items.length === 0) return null
                    return (
                      <div key={category}>
                        <h3 className={`text-xs font-bold uppercase tracking-wider mb-2 ${cfg.color}`}>{cfg.icon} {category}</h3>
                        <div className="flex flex-col gap-2">
                          {items.map(t => {
                            const isTalSelected = selectedTalentos.some(st => st.ID === t.ID)
                            const isDisabled    = !isTalSelected && selectedTalentos.length >= totalTalentos
                            return (
                              <button key={t.ID} type="button" disabled={isDisabled} onClick={() => toggleTalento(t)}
                                className="flex items-start justify-between rounded-lg px-4 py-3 border text-left transition"
                                style={isTalSelected
                                  ? { background: 'rgba(201,168,76,0.1)', borderColor: 'rgba(201,168,76,0.5)' }
                                  : isDisabled
                                  ? { background: '#1a1a1a', borderColor: '#3f3f46', opacity: 0.4, cursor: 'not-allowed' }
                                  : { background: '#27272a', borderColor: '#3f3f46' }
                                }
                              >
                                <div className="flex-1 min-w-0">
                                  <div className="flex items-center gap-2 flex-wrap mb-0.5">
                                    <span className="text-sm font-medium text-white">{t.name}</span>
                                    {t.prerequisite && (
                                      <span className="text-xs bg-orange-900/60 text-orange-300 px-1.5 py-0.5 rounded border border-orange-700/40">Req: {t.prerequisite}</span>
                                    )}
                                  </div>
                                  <p className="text-gray-400 text-xs">{t.description}</p>
                                </div>
                                <div className="flex items-center gap-2 ml-3 mt-0.5 flex-shrink-0">
                                  <Tooltip content={t.tooltip} />
                                  <div className="w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 transition"
                                    style={isTalSelected ? { background: '#c9a84c', borderColor: '#c9a84c' } : { borderColor: '#52525b' }}>
                                    {isTalSelected && <span className="text-black text-xs leading-none font-bold">✓</span>}
                                  </div>
                                </div>
                              </button>
                            )
                          })}
                        </div>
                      </div>
                    )
                  })}
                </div>
              </div>
            )}

            {selectedClass && hasSkills && hasAnyChoiceGroups && !allChoicesMade && (
              <div className="rounded-xl p-4" style={{ background: 'rgba(234,179,8,0.08)', border: '1px solid rgba(234,179,8,0.25)' }}>
                <p className="text-yellow-400 text-sm font-semibold">⚠️ Complete todas as escolhas acima para continuar</p>
              </div>
            )}

            {/* PASSO 9 — Poderes (só 4e) */}
            {is4e && selectedClass && hasSkills && allChoicesMade && (
              <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                <div className="flex justify-between items-center mb-1">
                  {sectionHeader(`Passo ${S()} — Poderes`)}
                  <span className="text-xs text-gray-500 bg-gray-700/60 px-3 py-1 rounded-full">{totalNormal} selecionado{totalNormal !== 1 ? 's' : ''}</span>
                </div>
                <p className="text-gray-500 text-xs mb-5">Clique no card para ver detalhes. Clique no círculo para selecionar.</p>
                <div className="flex flex-col gap-6">
                  {normalByType('unlimited').length > 0 && (
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <h3 className="text-sm font-bold text-green-400">⚡ Sem Limite</h3>
                        {counterBadge(selectedSkills.unlimited.length, 2)}
                      </div>
                      <div className="flex flex-col gap-2">
                        {normalByType('unlimited').map(s => <SkillCard key={s.ID} skill={s} selectable selected={isSelected(s)} disabled={!canSelect('unlimited', s) && !isSelected(s)} onToggle={toggleSkill} defaultExpanded />)}
                      </div>
                    </div>
                  )}
                  {normalByType('encounter').length > 0 && (
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <h3 className="text-sm font-bold text-yellow-400">⚔️ Por Encontro</h3>
                        {counterBadge(selectedSkills.encounter.length, 1)}
                      </div>
                      <div className="flex flex-col gap-2">
                        {normalByType('encounter').map(s => <SkillCard key={s.ID} skill={s} selectable selected={isSelected(s)} disabled={!canSelect('encounter', s) && !isSelected(s)} onToggle={toggleSkill} defaultExpanded />)}
                      </div>
                    </div>
                  )}
                  {normalByType('daily').length > 0 && (
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <h3 className="text-sm font-bold text-red-400">📅 Diário</h3>
                        {counterBadge(selectedSkills.daily.length, 1)}
                      </div>
                      <div className="flex flex-col gap-2">
                        {normalByType('daily').map(s => <SkillCard key={s.ID} skill={s} selectable selected={isSelected(s)} disabled={!canSelect('daily', s) && !isSelected(s)} onToggle={toggleSkill} defaultExpanded />)}
                      </div>
                    </div>
                  )}
                  {utilitySkills.length > 0 && (
                    <div className="opacity-40">
                      <h3 className="text-sm font-bold text-blue-400 mb-1">🔧 Utilitário</h3>
                      <p className="text-gray-500 text-xs">Disponível a partir do nível 2.</p>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* PASSO 10 — Armadura */}
            <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
              {sectionHeader(`Passo ${S()} — Armadura`)}
              <select {...register('armor_id')} className="rpg-select"
                onChange={e => { const v = Number(e.target.value); setValue('armor_id', v); setSelectedArmorData(armors?.find((a: Armor) => a.ID === v) ?? null) }}>
                <option value="">Sem armadura</option>
                {armors?.filter((a: Armor) => a.armor_type !== 'shield').map((a: Armor) => (
                  <option key={a.ID} value={a.ID}>{a.name} (CA base {a.base_ac})</option>
                ))}
              </select>
              {selectedArmorData && <p className="text-gray-500 text-xs mt-2">{selectedArmorData.description}</p>}
            </div>

            {/* PASSO 11 — Atributos */}
            <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
              {sectionHeader(`Passo ${S()} — Atributos`)}
              <p className="text-gray-500 text-xs mb-4">
                Alterar a <span style={{ color: '#c9a84c' }}>Constituição</span> recalcula o HP automaticamente.
                {is5e && ' Array padrão sugerido: 15, 14, 13, 12, 10, 8.'}
              </p>
              <div className="grid grid-cols-2 gap-3">
                {attrFields.map(attr => (
                  <div key={attr.key}>
                    <label className={`text-xs mb-1.5 block uppercase tracking-wider ${attr.key === 'constitution' ? 'font-semibold' : 'text-gray-500'}`}
                      style={attr.key === 'constitution' ? { color: '#c9a84c' } : {}}>
                      {attr.label}
                    </label>
                    <input type="number" {...register(attr.key)} min={1} max={20} className="rpg-input"
                      style={attr.key === 'constitution' ? { borderColor: 'rgba(201,168,76,0.3)' } : {}} />
                  </div>
                ))}
              </div>
            </div>

            {/* PASSO 12 — Hit Points */}
            <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
              <div className="flex justify-between items-center mb-3">
                {sectionHeader(`Passo ${S()} — Hit Points`)}
                {selectedClassData && <span className="text-xs text-gray-500 bg-gray-700/60 px-3 py-1 rounded-full">{hpLabel()}</span>}
              </div>
              <input type="number" {...register('hit_points')} min={1} className="rpg-input" />
            </div>

            {/* Resumo */}
            {(is4e || is5e) && (totalSelected > 0 || selectedPericias.length > 0 || selectedTalentos.length > 0 || selectedBackground || selectedAlignment) && (
              <div className="rounded-xl p-4" style={{ background: 'rgba(201,168,76,0.06)', border: '1px solid rgba(201,168,76,0.2)' }}>
                <p className="text-xs font-semibold mb-2 uppercase tracking-widest" style={{ color: '#c9a84c' }}>Resumo da Criação</p>
                <div className="flex flex-wrap gap-2">
                  {selectedBackground && <span className="text-xs px-2 py-1 rounded-full bg-violet-900/60 text-violet-300">📜 {selectedBackground.name}</span>}
                  {selectedAlignment  && <span className="text-xs px-2 py-1 rounded-full bg-sky-900/60 text-sky-300">⚖️ {selectedAlignment}</span>}
                  {Object.entries(choiceSelections).filter(([, s]) => s.is_class_feature).map(([, s]) => (
                    <span key={s.ID} className="text-xs px-2 py-1 rounded-full bg-purple-900/60 text-purple-300">{s.name}</span>
                  ))}
                  {Object.entries(choiceSelections).filter(([, s]) => s.is_race_feature).map(([, s]) => (
                    <span key={s.ID} className="text-xs px-2 py-1 rounded-full bg-amber-900/60 text-amber-300">{s.name}</span>
                  ))}
                  {Object.values(selectedSkills).flat().map(s => (
                    <span key={s.ID} className={`text-xs px-2 py-1 rounded-full ${powerConfig[(s.power_type as PowerType) ?? 'unlimited'].badge}`}>{s.name}</span>
                  ))}
                  {selectedPericias.map(name => (
                    <span key={name} className={`text-xs px-2 py-1 rounded-full ${backgroundSkills.includes(name) ? 'bg-indigo-900/60 text-indigo-300' : 'bg-teal-900/60 text-teal-300'}`}>📚 {name}</span>
                  ))}
                  {selectedTalentos.map(t => <span key={t.ID} className="text-xs px-2 py-1 rounded-full bg-orange-900/60 text-orange-300">🏆 {t.name}</span>)}
                </div>
              </div>
            )}

            {pendingItems.length > 0 && (
              <div className="rounded-xl p-4" style={{ background: 'rgba(220,38,38,0.08)', border: '1px solid rgba(220,38,38,0.25)' }}>
                <p className="text-red-400 text-xs font-semibold mb-2 uppercase tracking-wider">Complete os itens abaixo para criar o personagem:</p>
                <ul className="space-y-1">
                  {pendingItems.map(item => <li key={item} className="text-red-400/70 text-xs">• {item}</li>)}
                </ul>
              </div>
            )}

            <button type="submit" disabled={!canSubmit}
              className="w-full rounded-lg py-2 text-sm font-semibold transition mb-4"
              style={canSubmit
                ? { background: '#c9a84c', color: '#0a0a0a', border: '1px solid #c9a84c', cursor: 'pointer', fontFamily: 'Georgia, serif', letterSpacing: '0.04em' }
                : { background: 'transparent', color: '#52525b', border: '1px solid #3f3f46', cursor: 'not-allowed' }
              }
            >
              {createMutation.isPending ? 'Criando...' : canSubmit ? '✦ Criar Personagem' : 'Complete todos os passos'}
            </button>

          </form>
        )}
      </div>
    </div>
  )
}