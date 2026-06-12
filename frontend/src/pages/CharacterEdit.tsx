import { useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { characterService } from '../services/characterService'

const schema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  edition: z.string(),
  class_id: z.coerce.number(),
  race_id: z.coerce.number(),
  hit_points: z.coerce.number().min(1, 'HP deve ser maior que zero'),
  strength: z.coerce.number().min(1).max(20),
  dexterity: z.coerce.number().min(1).max(20),
  constitution: z.coerce.number().min(1).max(20),
  intelligence: z.coerce.number().min(1).max(20),
  wisdom: z.coerce.number().min(1).max(20),
  charisma: z.coerce.number().min(1).max(20),
})

type FormData = z.infer<typeof schema>

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
        name: character.name,
        edition: character.edition,
        class_id: character.class_id,
        race_id: character.race_id,
        hit_points: character.hit_points,
        strength: character.strength,
        dexterity: character.dexterity,
        constitution: character.constitution,
        intelligence: character.intelligence,
        wisdom: character.wisdom,
        charisma: character.charisma,
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

  const attributes = [
    { label: 'Força', key: 'strength' },
    { label: 'Destreza', key: 'dexterity' },
    { label: 'Constituição', key: 'constitution' },
    { label: 'Inteligência', key: 'intelligence' },
    { label: 'Sabedoria', key: 'wisdom' },
    { label: 'Carisma', key: 'charisma' },
  ] as const

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <p className="text-gray-400">Carregando...</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-2xl mx-auto">
        <button
          onClick={() => navigate(`/characters/${id}`)}
          className="text-gray-400 hover:text-white transition mb-6 block"
        >
          ← Voltar
        </button>

        <h1 className="text-3xl font-bold text-white mb-8">Editar Personagem</h1>

        <form onSubmit={handleSubmit(data => updateMutation.mutate(data))} className="flex flex-col gap-6">

          {/* Informações travadas */}
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700 flex flex-col gap-4">
            <div className="flex items-center justify-between mb-2">
              <h2 className="text-lg font-semibold text-white">Informações do Personagem</h2>
              <span className="text-xs text-gray-500 bg-gray-700 px-3 py-1 rounded-full">
                🔒 Edição, classe e raça não podem ser alteradas
              </span>
            </div>

            {/* Edição travada */}
            <div>
              <label className="text-gray-400 text-sm mb-1 block">Edição</label>
              <div className="w-full bg-gray-900 text-gray-400 rounded-lg px-4 py-2 border border-gray-700 cursor-not-allowed">
                {character?.edition}
              </div>
              <input type="hidden" {...register('edition')} />
            </div>

            {/* Classe travada */}
            <div>
              <label className="text-gray-400 text-sm mb-1 block">Classe</label>
              <div className="w-full bg-gray-900 text-gray-400 rounded-lg px-4 py-2 border border-gray-700 cursor-not-allowed">
                {character?.class?.name}
              </div>
              <input type="hidden" {...register('class_id')} />
            </div>

            {/* Raça travada */}
            <div>
              <label className="text-gray-400 text-sm mb-1 block">Raça</label>
              <div className="w-full bg-gray-900 text-gray-400 rounded-lg px-4 py-2 border border-gray-700 cursor-not-allowed">
                {character?.race?.name}
              </div>
              <input type="hidden" {...register('race_id')} />
            </div>
          </div>

          {/* Nome editável */}
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h2 className="text-lg font-semibold text-white mb-4">Nome</h2>
            <input
              {...register('name')}
              className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            {errors.name && <p className="text-red-400 text-xs mt-1">{errors.name.message}</p>}
          </div>

          {/* Atributos editáveis */}
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h2 className="text-lg font-semibold text-white mb-4">Atributos</h2>
            <div className="grid grid-cols-2 gap-4">
              {attributes.map(attr => (
                <div key={attr.key}>
                  <label className="text-gray-400 text-sm mb-1 block">{attr.label}</label>
                  <input
                    type="number"
                    {...register(attr.key)}
                    min={1}
                    max={20}
                    className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
              ))}
            </div>
          </div>

          {/* HP editável */}
          <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <h2 className="text-lg font-semibold text-white mb-4">Hit Points</h2>
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
            disabled={updateMutation.isPending}
            className="bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-bold py-3 rounded-lg transition"
          >
            {updateMutation.isPending ? 'Salvando...' : 'Salvar Alterações'}
          </button>

        </form>
      </div>
    </div>
  )
}