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
import { SkillCard, powerConfig } from '../components/SkillCard'
import type { Class, Armor, Skill, PowerType } from '../types'

const schema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  edition: z.string().min(1, 'Edição é obrigatória'),
  class_id: z.coerce.number().min(1, 'Selecione uma classe'),
  race_id: z.coerce.number().min(1, 'Selecione uma raça'),
  hit_points: z.coerce.number().min(1),
  strength: z.coerce.number().min(1).max(20),
  dexterity: z.coerce.number().min(1).max(20),
  constitution: z.coerce.number().min(1).max(20),
  intelligence: z.coerce.number().min(1).max(20),
  wisdom: z.coerce.number().min(1).max(20),
  charisma: z.coerce.number().min(1).max(20),
  armor_id: z.coerce.number().optional(),
})
type FormData = z.infer<typeof schema>

// Limites de seleção por tipo no nível 1
const SKILL_LIMITS: Record<PowerType, number> = { unlimited: 2, encounter: 1, daily: 1, utility: 1 }

const editions = ['4e', '5e']
function getConMod(con: number) { return Math.floor((con - 10) / 2) }

export default function CharacterCreate() {
  const navigate = useNavigate()
  const [selectedEdition, setSelectedEdition] = useState<string | null>(null)
  const [selectedClass, setSelectedClass] = useState<number | null>(null)
  const [selectedRace, setSelectedRace] = useState<number | null>(null)
  const [selectedClassData, setSelectedClassData] = useState<Class | null>(null)
  const [selectedArmorData, setSelectedArmorData] = useState<Armor | null>(null)

  // Poderes selecionados por tipo
  const [selectedSkills, setSelectedSkills] = useState<Record<PowerType, Skill[]>>({
    unlimited: [], encounter: [], daily: [], utility: [],
  })
  // Escolha de características (grupos mutuamente exclusivos)
  const [choiceSelections, setChoiceSelections] = useState<Record<string, Skill>>({})

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

  const { data: classes }   = useQuery({ queryKey: ['classes', selectedEdition],  queryFn: () => classService.getAll(selectedEdition!), enabled: !!selectedEdition })
  const { data: armors }    = useQuery({ queryKey: ['armors', selectedEdition],   queryFn: () => armorService.getByEdition(selectedEdition!), enabled: !!selectedEdition })
  const { data: races }     = useQuery({ queryKey: ['races', selectedEdition],    queryFn: () => raceService.getAll(selectedEdition!), enabled: !!selectedEdition })
  const { data: allSkills } = useQuery({
    queryKey: ['skills-filter', selectedClass, selectedRace],
    queryFn: () => skillService.getByFilter(selectedClass!, selectedRace ?? undefined),
    enabled: !!selectedClass,
  })

  // Separa skills por categoria
  const classFeatures   = (allSkills ?? []).filter(s => s.is_class_feature && !s.requires_choice)
  const choiceFeatures  = (allSkills ?? []).filter(s => s.is_class_feature && s.requires_choice)
  const choiceGroups    = choiceFeatures.reduce<Record<string, Skill[]>>((acc, s) => {
    const g = s.choice_group ?? 'default'
    acc[g] = [...(acc[g] ?? []), s]; return acc
  }, {})
  const normalByType = (type: PowerType): Skill[] =>
    (allSkills ?? []).filter(s => !s.is_class_feature && s.power_type === type && (!s.level || s.level <= 1))
  const utilitySkills = (allSkills ?? []).filter(s => !s.is_class_feature && s.power_type === 'utility')

  const createMutation = useMutation({
    mutationFn: async (data: FormData) => {
      const character = await characterService.create(data)
      // Poderes selecionados
      for (const skill of Object.values(selectedSkills).flat()) {
        await characterService.addSkill(character.ID, skill.ID)
      }
      // Escolhas de características
      for (const skill of Object.values(choiceSelections)) {
        await characterService.addSkill(character.ID, skill.ID)
      }
      // Características automáticas
      for (const skill of classFeatures) {
        await characterService.addSkill(character.ID, skill.ID)
      }
      return character
    },
    onSuccess: () => navigate('/characters'),
  })

  const handleEditionChange = (edition: string) => {
    setSelectedEdition(edition); setSelectedClass(null); setSelectedRace(null); setSelectedClassData(null)
    setSelectedSkills({ unlimited: [], encounter: [], daily: [], utility: [] })
    setChoiceSelections({})
    reset({ name: '', edition, class_id: 0, race_id: 0, hit_points: 10, strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10 })
    setValue('edition', edition)
  }

  const handleClassChange = (classId: number) => {
    setSelectedClass(classId); setValue('class_id', classId)
    setSelectedClassData(classes?.find(c => c.ID === classId) ?? null)
    setSelectedSkills({ unlimited: [], encounter: [], daily: [], utility: [] })
    setChoiceSelections({})
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

  const selectChoice = (skill: Skill) => {
    setChoiceSelections(prev => ({ ...prev, [skill.choice_group ?? 'default']: skill }))
  }
  const isChoiceSelected = (skill: Skill) =>
    choiceSelections[skill.choice_group ?? 'default']?.ID === skill.ID

  const totalNormal   = Object.values(selectedSkills).flat().length
  const totalChoices  = Object.keys(choiceSelections).length
  const totalSelected = totalNormal + totalChoices

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
  const hasSkills = allSkills && allSkills.length > 0

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-2xl mx-auto">
        <button onClick={() => navigate('/characters')} className="text-gray-400 hover:text-white transition mb-6 block">← Voltar</button>
        <h1 className="text-3xl font-bold text-white mb-8">Criar Personagem</h1>

        {/* Passo 1 — Edição */}
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 mb-6">
          <h2 className="text-lg font-semibold text-white mb-2">Passo 1 — Escolha a Edição</h2>
          <div className="grid grid-cols-2 gap-3">
            {editions.map(e => (
              <button key={e} type="button" onClick={() => handleEditionChange(e)}
                className={`py-4 rounded-lg font-semibold text-base transition border ${selectedEdition === e ? 'bg-indigo-600 border-indigo-500 text-white' : 'bg-gray-700 border-gray-600 text-gray-300 hover:bg-gray-600'}`}>
                D&D {e}
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

            {/* Passo 4 — Habilidades (4e) */}
            {is4e && selectedClass && hasSkills && (
              <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
                <div className="flex justify-between items-center mb-1">
                  <h2 className="text-lg font-semibold text-white">Passo 4 — Habilidades</h2>
                  <span className="text-xs text-gray-400 bg-gray-700 px-3 py-1 rounded-full">{totalSelected} selecionado{totalSelected !== 1 ? 's' : ''}</span>
                </div>
                <p className="text-gray-400 text-sm mb-5">Clique em qualquer habilidade para ver os detalhes completos.</p>

                <div className="flex flex-col gap-6">

                  {/* Características automáticas */}
                  {classFeatures.length > 0 && (
                    <div>
                      <h3 className="text-sm font-bold text-indigo-400 mb-1">📖 Características de Classe</h3>
                      <p className="text-gray-500 text-xs mb-3">Concedidas automaticamente ao criar o personagem.</p>
                      <div className="flex flex-col gap-2">
                        {classFeatures.map(s => <SkillCard key={s.ID} skill={s} informative defaultExpanded />)}
                      </div>
                    </div>
                  )}

                  {/* Grupos de escolha de características */}
                  {Object.entries(choiceGroups).map(([group, options]) => (
                    <div key={group}>
                      <h3 className="text-sm font-bold text-purple-400 mb-1">🎯 Escolha uma Característica</h3>
                      <p className="text-gray-500 text-xs mb-3">Selecione uma das opções abaixo — esta será sua característica permanente.</p>
                      <div className="flex flex-col gap-2">
                        {options.map(s => (
                          <SkillCard
                            key={s.ID} skill={s}
                            selectable selected={isChoiceSelected(s)} disabled={false}
                            onToggle={selectChoice} defaultExpanded
                          />
                        ))}
                      </div>
                    </div>
                  ))}

                  {/* Sem Limite */}
                  {normalByType('unlimited').length > 0 && (
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <h3 className="text-sm font-bold text-green-400">⚡ Sem Limite</h3>
                        <span className={`text-xs px-2 py-1 rounded-full ${selectedSkills.unlimited.length >= 2 ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-400'}`}>
                          {selectedSkills.unlimited.length} / 2
                        </span>
                      </div>
                      <div className="flex flex-col gap-2">
                        {normalByType('unlimited').map(s => (
                          <SkillCard key={s.ID} skill={s} selectable selected={isSelected(s)} disabled={!canSelect('unlimited', s) && !isSelected(s)} onToggle={toggleSkill} defaultExpanded />
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Por Encontro */}
                  {normalByType('encounter').length > 0 && (
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <h3 className="text-sm font-bold text-yellow-400">⚔️ Por Encontro</h3>
                        <span className={`text-xs px-2 py-1 rounded-full ${selectedSkills.encounter.length >= 1 ? 'bg-yellow-900 text-yellow-300' : 'bg-gray-700 text-gray-400'}`}>
                          {selectedSkills.encounter.length} / 1
                        </span>
                      </div>
                      <div className="flex flex-col gap-2">
                        {normalByType('encounter').map(s => (
                          <SkillCard key={s.ID} skill={s} selectable selected={isSelected(s)} disabled={!canSelect('encounter', s) && !isSelected(s)} onToggle={toggleSkill} defaultExpanded />
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Diário */}
                  {normalByType('daily').length > 0 && (
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <h3 className="text-sm font-bold text-red-400">📅 Diário</h3>
                        <span className={`text-xs px-2 py-1 rounded-full ${selectedSkills.daily.length >= 1 ? 'bg-red-900 text-red-300' : 'bg-gray-700 text-gray-400'}`}>
                          {selectedSkills.daily.length} / 1
                        </span>
                      </div>
                      <div className="flex flex-col gap-2">
                        {normalByType('daily').map(s => (
                          <SkillCard key={s.ID} skill={s} selectable selected={isSelected(s)} disabled={!canSelect('daily', s) && !isSelected(s)} onToggle={toggleSkill} defaultExpanded />
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Utilitário — bloqueado no nível 1 */}
                  {utilitySkills.length > 0 && (
                    <div className="opacity-40">
                      <h3 className="text-sm font-bold text-blue-400 mb-1">🔧 Utilitário</h3>
                      <p className="text-gray-500 text-xs">Disponível a partir do nível 2.</p>
                    </div>
                  )}

                </div>
              </div>
            )}

            {/* Armadura */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <h2 className="text-lg font-semibold text-white mb-4">{is4e ? 'Passo 5' : 'Passo 4'} — Armadura</h2>
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
              <h2 className="text-lg font-semibold text-white mb-1">{is4e ? 'Passo 6' : 'Passo 5'} — Atributos</h2>
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
                <h2 className="text-lg font-semibold text-white">{is4e ? 'Passo 7' : 'Passo 6'} — Hit Points</h2>
                {selectedClassData && <span className="text-xs text-gray-400 bg-gray-700 px-3 py-1 rounded-full">{hpLabel()}</span>}
              </div>
              <input type="number" {...register('hit_points')} min={1} className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500" />
            </div>

            {/* Resumo de seleção */}
            {is4e && totalSelected > 0 && (
              <div className="bg-gray-800 rounded-lg p-4 border border-indigo-700">
                <p className="text-indigo-300 text-sm font-semibold mb-2">✨ Selecionados ({totalSelected})</p>
                <div className="flex flex-wrap gap-2">
                  {Object.entries(choiceSelections).map(([, s]) => (
                    <span key={s.ID} className="text-xs px-2 py-1 rounded-full bg-purple-900 text-purple-300">{s.name}</span>
                  ))}
                  {Object.values(selectedSkills).flat().map(s => (
                    <span key={s.ID} className={`text-xs px-2 py-1 rounded-full ${powerConfig[(s.power_type as PowerType) ?? 'unlimited'].badge}`}>{s.name}</span>
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