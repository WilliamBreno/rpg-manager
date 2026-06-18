import { useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { characterService } from '../services/characterService'

const schema = z.object({
  name:         z.string().min(1, 'Nome é obrigatório'),
  edition:      z.string(),
  class_id:     z.coerce.number(),
  race_id:      z.coerce.number(),
  hit_points:   z.coerce.number().min(1, 'HP deve ser maior que zero'),
  strength:     z.coerce.number().min(1).max(20),
  dexterity:    z.coerce.number().min(1).max(20),
  constitution: z.coerce.number().min(1).max(20),
  intelligence: z.coerce.number().min(1).max(20),
  wisdom:       z.coerce.number().min(1).max(20),
  charisma:     z.coerce.number().min(1).max(20),
})

type FormData = z.infer<typeof schema>

const attributes = [
  { label: 'Força',        key: 'strength'     },
  { label: 'Destreza',     key: 'dexterity'    },
  { label: 'Constituição', key: 'constitution' },
  { label: 'Inteligência', key: 'intelligence' },
  { label: 'Sabedoria',    key: 'wisdom'       },
  { label: 'Carisma',      key: 'charisma'     },
] as const

export default function CharacterEdit() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: character, isLoading } = useQuery({
    queryKey: ['character', id],
    queryFn: () => characterService.getByID(Number(id)),
  })

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema) as any,
  })

  useEffect(() => {
    if (character) {
      reset({
        name:         character.name,
        edition:      character.edition,
        class_id:     character.class_id,
        race_id:      character.race_id,
        hit_points:   character.hit_points,
        strength:     character.strength,
        dexterity:    character.dexterity,
        constitution: character.constitution,
        intelligence: character.intelligence,
        wisdom:       character.wisdom,
        charisma:     character.charisma,
      })
    }
  }, [character, reset])

  const updateMutation = useMutation({
    mutationFn: (data: FormData) => characterService.update(Number(id), data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character', id] })
      queryClient.invalidateQueries({ queryKey: ['characters'] })
      navigate(`/characters/${id}`)
    },
  })

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <p className="text-gray-500">Carregando...</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-2xl mx-auto">

        {/* Voltar */}
        <button
          onClick={() => navigate(`/characters/${id}`)}
          className="transition mb-6 block text-sm"
          style={{ color: 'rgba(201,168,76,0.5)' }}
          onMouseEnter={e => (e.currentTarget.style.color = '#c9a84c')}
          onMouseLeave={e => (e.currentTarget.style.color = 'rgba(201,168,76,0.5)')}
        >
          ← Voltar
        </button>

        <h1 className="font-rpg text-2xl sm:text-3xl font-bold mb-6 sm:mb-8" style={{ color: '#c9a84c' }}>
          Editar Personagem
        </h1>

        <form onSubmit={handleSubmit(data => updateMutation.mutate(data))} className="flex flex-col gap-4">

          {/* ── Informações travadas ─────────────────────────────────────────── */}
          <div className="bg-gray-800 rounded-xl p-5 border border-gray-700 flex flex-col gap-4">
            <div className="flex items-center justify-between">
              <h2 className="text-sm font-semibold uppercase tracking-widest" style={{ color: 'rgba(201,168,76,0.7)' }}>
                Informações do Personagem
              </h2>
              <span
                className="text-xs px-3 py-1 rounded-full"
                style={{ background: 'rgba(201,168,76,0.08)', border: '1px solid rgba(201,168,76,0.2)', color: 'rgba(201,168,76,0.5)' }}
              >
                Não editável
              </span>
            </div>

            {[
              { label: 'Edição',  value: character?.edition,     field: 'edition'  },
              { label: 'Classe',  value: character?.class?.name, field: 'class_id' },
              { label: 'Raça',    value: character?.race?.name,  field: 'race_id'  },
            ].map(item => (
              <div key={item.field}>
                <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">
                  {item.label}
                </label>
                <div
                  className="w-full rounded-lg px-4 py-2 text-sm cursor-not-allowed"
                  style={{ background: '#0f0f0f', border: '1px solid #2a2a2a', color: '#52525b' }}
                >
                  {item.value}
                </div>
                <input type="hidden" {...register(item.field as any)} />
              </div>
            ))}
          </div>

          {/* ── Nome ────────────────────────────────────────────────────────── */}
          <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Nome
            </h2>
            <input
              {...register('name')}
              className="rpg-input"
              placeholder="Nome do personagem"
            />
            {errors.name && <p className="text-red-400 text-xs mt-1">{errors.name.message}</p>}
          </div>

          {/* ── Atributos ───────────────────────────────────────────────────── */}
          <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-4" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Atributos
            </h2>
            <div className="grid grid-cols-2 gap-3">
              {attributes.map(attr => (
                <div key={attr.key}>
                  <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">
                    {attr.label}
                  </label>
                  <input
                    type="number"
                    {...register(attr.key)}
                    min={1}
                    max={20}
                    className="rpg-input"
                  />
                </div>
              ))}
            </div>
          </div>

          {/* ── Hit Points ──────────────────────────────────────────────────── */}
          <div className="bg-gray-800 rounded-xl p-5 border border-gray-700">
            <h2 className="text-sm font-semibold uppercase tracking-widest mb-3" style={{ color: 'rgba(201,168,76,0.7)' }}>
              Hit Points
            </h2>
            <input
              type="number"
              {...register('hit_points')}
              min={1}
              className="rpg-input"
            />
            {errors.hit_points && (
              <p className="text-red-400 text-xs mt-1">{errors.hit_points.message}</p>
            )}
          </div>

          {/* ── Botões ──────────────────────────────────────────────────────── */}
          <div className="flex gap-2 pt-1 pb-6">
            <button
              type="button"
              onClick={() => navigate(`/characters/${id}`)}
              className="btn-rpg-outline flex-1 sm:flex-none"
            >
              Cancelar
            </button>
            <button
              type="submit"
              disabled={updateMutation.isPending}
              className="btn-rpg-primary flex-1"
            >
              {updateMutation.isPending ? 'Salvando...' : 'Salvar'}
            </button>
          </div>

        </form>
      </div>
    </div>
  )
}