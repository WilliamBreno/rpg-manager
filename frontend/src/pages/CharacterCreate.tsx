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
import type { Class, Armor } from '../types'

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

// Apenas 4e e 5e
const editions = ['4e', '5e']

function getConMod(con: number) {
  return Math.floor((con - 10) / 2)
}

export default function CharacterCreate() {
  const navigate = useNavigate()
  const [selectedEdition, setSelectedEdition] = useState<string | null>(null)
  const [selectedClass, setSelectedClass] = useState<number | null>(null)
  const [selectedRace, setSelectedRace] = useState<number | null>(null)
  const [selectedClassData, setSelectedClassData] = useState<Class | null>(null)
  const [selectedArmorData, setSelectedArmorData] = useState<Armor | null>(null)

  const { register, handleSubmit, setValue, reset, control, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema) as any,
    defaultValues: {
      strength: 10, dexterity: 10, constitution: 10,
      intelligence: 10, wisdom: 10, charisma: 10,
      hit_points: 10,
    }
  })

  const constitution = useWatch({ control, name: 'constitution' })

  // Recalcula HP automaticamente quando classe ou CON mudar
  useEffect(() => {
    if (!selectedClassData || !selectedEdition) return

    const conValue = Number(constitution)
    const conMod = getConMod(conValue)

    let newHP: number

    if (selectedEdition === '4e') {
      // 4e: BaseHP da classe + valor de CON (não modificador!)
      newHP = (selectedClassData.base_hp ?? 0) + conValue
    } else {
      // 5e: HitDie máximo + modificador de CON
      newHP = (selectedClassData.hit_die ?? 0) + conMod
    }

    setValue('hit_points', Math.max(1, newHP))
  }, [selectedClassData, constitution, selectedEdition, setValue])

  const { data: classes } = useQuery({
    queryKey: ['classes', selectedEdition],
    queryFn: () => classService.getAll(selectedEdition!),
    enabled: !!selectedEdition,
  })

  const { data: armors } = useQuery({
    queryKey: ['armors', selectedEdition],
    queryFn: () => armorService.getByEdition(selectedEdition!),
    enabled: !!selectedEdition,
  })

  const { data: races } = useQuery({
    queryKey: ['races', selectedEdition],
    queryFn: () => raceService.getAll(selectedEdition!),
    enabled: !!selectedEdition,
  })

  const { data: skills } = useQuery({
    queryKey: ['skills', selectedClass, selectedRace],
    queryFn: () => skillService.getByFilter(selectedClass!, selectedRace!),
    enabled: !!selectedClass && !!selectedRace,
  })

  const createMutation = useMutation({
    mutationFn: characterService.create,
    onSuccess: () => navigate('/characters'),
  })

  const handleEditionChange = (edition: string) => {
    setSelectedEdition(edition)
    setSelectedClass(null)
    setSelectedRace(null)
    setSelectedClassData(null)
    reset({
      name: '',
      edition,
      class_id: 0,
      race_id: 0,
      hit_points: 10,
      strength: 10,
      dexterity: 10,
      constitution: 10,
      intelligence: 10,
      wisdom: 10,
      charisma: 10,
    })
    setValue('edition', edition)
  }

  const handleClassChange = (classId: number) => {
    setSelectedClass(classId)
    setValue('class_id', classId)
    const classData = classes?.find(c => c.ID === classId) ?? null
    setSelectedClassData(classData)
  }

  const attributes = [
    { label: 'Força', key: 'strength' },
    { label: 'Destreza', key: 'dexterity' },
    { label: 'Constituição', key: 'constitution' },
    { label: 'Inteligência', key: 'intelligence' },
    { label: 'Sabedoria', key: 'wisdom' },
    { label: 'Carisma', key: 'charisma' },
  ] as const

  const powerTypeLabel: Record<string, string> = {
    utility: 'Utilitário',
    unlimited: 'Sem Limite',
    encounter: 'Por Encontro',
    daily: 'Diário',
  }

  const powerTypeColor: Record<string, string> = {
    utility: 'bg-blue-900 text-blue-300',
    unlimited: 'bg-green-900 text-green-300',
    encounter: 'bg-yellow-900 text-yellow-300',
    daily: 'bg-red-900 text-red-300',
  }

  // Info da classe conforme edição
  const classInfo = () => {
    if (!selectedClassData) return null
    if (selectedEdition === '4e') {
      return `PV base: ${selectedClassData.base_hp ?? '?'} + CON — ${selectedClassData.description}`
    }
    return `Hit Die: d${selectedClassData.hit_die} — ${selectedClassData.description}`
  }

  // Label do HP conforme edição
  const hpLabel = () => {
    if (!selectedClassData) return null
    const conValue = Number(constitution)
    const conMod = getConMod(conValue)

    if (selectedEdition === '4e') {
      return `${selectedClassData.base_hp ?? '?'} (base) + ${conValue} (CON) = ${Math.max(1, (selectedClassData.base_hp ?? 0) + conValue)} PV`
    }
    return `d${selectedClassData.hit_die} + mod CON (${conMod >= 0 ? '+' : ''}${conMod})`
  }

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-2xl mx-auto">
        <button
          onClick={() => navigate('/characters')}
          className="text-gray-400 hover:text-white transition mb-6 block"
        >
          ← Voltar
        </button>

        <h1 className="text-3xl font-bold text-white mb-8">Criar Personagem</h1>

        {/* Passo 1 — Edição */}
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 mb-6">
          <h2 className="text-lg font-semibold text-white mb-2">Passo 1 — Escolha a Edição</h2>
          <p className="text-gray-400 text-sm mb-4">Selecione qual edição do D&D você vai jogar.</p>
          <div className="grid grid-cols-2 gap-3">
            {editions.map(edition => (
              <button
                key={edition}
                type="button"
                onClick={() => handleEditionChange(edition)}
                className={`py-4 rounded-lg font-semibold text-base transition border ${
                  selectedEdition === edition
                    ? 'bg-indigo-600 border-indigo-500 text-white'
                    : 'bg-gray-700 border-gray-600 text-gray-300 hover:bg-gray-600'
                }`}
              >
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
              <input
                {...register('name')}
                className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="Ex: Aragorn"
              />
              {errors.name && <p className="text-red-400 text-xs mt-1">{errors.name.message}</p>}
            </div>

            {/* Passo 3 — Classe e Raça */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 flex flex-col gap-4">
              <h2 className="text-lg font-semibold text-white">Passo 3 — Classe e Raça</h2>

              <div>
                <label className="text-gray-400 text-sm mb-1 block">Classe</label>
                <select
                  {...register('class_id', { valueAsNumber: true })}
                  onChange={e => handleClassChange(Number(e.target.value))}
                  className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="">Selecione a classe</option>
                  {classes?.map(c => (
                    <option key={c.ID} value={c.ID}>{c.name}</option>
                  ))}
                </select>
                {selectedClassData && (
                  <p className="text-gray-400 text-xs mt-1">{classInfo()}</p>
                )}
                {errors.class_id && <p className="text-red-400 text-xs mt-1">Selecione uma classe</p>}
              </div>

              <div>
                <label className="text-gray-400 text-sm mb-1 block">Raça</label>
                <select
                  {...register('race_id', { valueAsNumber: true })}
                  onChange={e => {
                    const val = Number(e.target.value)
                    setValue('race_id', val)
                    setSelectedRace(val)
                  }}
                  className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="">Selecione a raça</option>
                  {races?.map(r => (
                    <option key={r.ID} value={r.ID}>{r.name}</option>
                  ))}
                </select>
                {errors.race_id && <p className="text-red-400 text-xs mt-1">Selecione uma raça</p>}
              </div>

              {skills && skills.length > 0 && (
                <div>
                  <p className="text-gray-400 text-sm mb-2">Habilidades disponíveis:</p>
                  <div className="flex flex-col gap-2">
                    {skills.map(skill => (
                      <div key={skill.ID} className="bg-gray-700 rounded-lg p-3">
                        <div className="flex justify-between items-center">
                          <p className="text-white text-sm font-semibold">{skill.name}</p>
                          <span className={`text-xs px-2 py-1 rounded-full ${powerTypeColor[skill.power_type ?? '']}`}>
                            {powerTypeLabel[skill.power_type ?? '']}
                          </span>
                        </div>
                        <p className="text-gray-400 text-xs mt-1">{skill.description}</p>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Armadura */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <h2 className="text-lg font-semibold text-white mb-4">Passo 3.5 — Armadura</h2>
              <select
                {...register('armor_id')}
                className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                onChange={e => {
                  const val = Number(e.target.value)
                  setValue('armor_id', val)
                  const armor = armors?.find((a: Armor) => a.ID === val) ?? null
                  setSelectedArmorData(armor)
                }}
              >
                <option value="">Sem armadura</option>
                {armors?.filter((a: Armor) => a.armor_type !== 'shield').map((a: Armor) => (
                  <option key={a.ID} value={a.ID}>
                    {a.name} (CA base {a.base_ac})
                  </option>
                ))}
              </select>
              {selectedArmorData && (
                <p className="text-gray-400 text-xs mt-2">{selectedArmorData.description}</p>
              )}
            </div>

            {/* Passo 4 — Atributos */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <h2 className="text-lg font-semibold text-white mb-1">Passo 4 — Atributos</h2>
              <p className="text-gray-400 text-sm mb-4">
                Alterar a <span className="text-indigo-400">Constituição</span> recalcula o HP automaticamente.
              </p>
              <div className="grid grid-cols-2 gap-4">
                {attributes.map(attr => (
                  <div key={attr.key}>
                    <label className={`text-sm mb-1 block ${attr.key === 'constitution' ? 'text-indigo-400 font-semibold' : 'text-gray-400'}`}>
                      {attr.label}
                    </label>
                    <input
                      type="number"
                      {...register(attr.key)}
                      min={1}
                      max={20}
                      className={`w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 ${
                        attr.key === 'constitution' ? 'focus:ring-indigo-400 ring-1 ring-indigo-800' : 'focus:ring-indigo-500'
                      }`}
                    />
                  </div>
                ))}
              </div>
            </div>

            {/* Passo 5 — HP calculado */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
              <div className="flex justify-between items-center mb-2">
                <h2 className="text-lg font-semibold text-white">Passo 5 — Hit Points</h2>
                {selectedClassData && (
                  <span className="text-xs text-gray-400 bg-gray-700 px-3 py-1 rounded-full">
                    {hpLabel()}
                  </span>
                )}
              </div>
              {selectedClassData && (
                <p className="text-gray-400 text-xs mb-3">
                  Calculado automaticamente com base na sua classe e Constituição. Você pode ajustar manualmente.
                </p>
              )}
              <input
                type="number"
                {...register('hit_points')}
                min={1}
                className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
              {errors.hit_points && <p className="text-red-400 text-xs mt-1">{errors.hit_points.message}</p>}
            </div>

            <button
              type="submit"
              disabled={createMutation.isPending}
              className="bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white font-bold py-3 rounded-lg transition"
            >
              {createMutation.isPending ? 'Criando...' : '✨ Criar Personagem'}
            </button>

          </form>
        )}
      </div>
    </div>
  )
}