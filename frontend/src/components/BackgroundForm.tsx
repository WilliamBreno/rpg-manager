import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { Background } from '../types'
import { backgroundService } from '../services/backgroundService'

interface Props {
  characterID: number
  background?: Background
}

export default function BackgroundForm({ characterID, background }: Props) {
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(!background)

  const [form, setForm] = useState<Background>({
    history: background?.history ?? '',
    personality_traits: background?.personality_traits ?? '',
    ideals: background?.ideals ?? '',
    bonds: background?.bonds ?? '',
    flaws: background?.flaws ?? '',
  })

  const saveMutation = useMutation({
    mutationFn: () => backgroundService.save(characterID, form),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character', String(characterID)] })
      setIsEditing(false)
    },
  })

  const fields = [
    { key: 'history', label: 'História', placeholder: 'Conte a origem do seu personagem...' },
    { key: 'personality_traits', label: 'Traços de Personalidade', placeholder: 'Como seu personagem se comporta?' },
    { key: 'ideals', label: 'Ideais', placeholder: 'No que seu personagem acredita?' },
    { key: 'bonds', label: 'Vínculos', placeholder: 'O que seu personagem se importa?' },
    { key: 'flaws', label: 'Defeitos', placeholder: 'Quais são as fraquezas do seu personagem?' },
  ] as const

  return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold text-white">Background</h2>
        {!isEditing && (
          <button
            onClick={() => setIsEditing(true)}
            className="text-indigo-400 hover:text-indigo-300 text-sm transition"
          >
            ✏️ Editar
          </button>
        )}
      </div>

      {isEditing ? (
        <div className="flex flex-col gap-4">
          {fields.map(field => (
            <div key={field.key}>
              <label className="text-gray-400 text-sm mb-1 block">{field.label}</label>
              <textarea
                value={form[field.key]}
                onChange={e => setForm(prev => ({ ...prev, [field.key]: e.target.value }))}
                placeholder={field.placeholder}
                rows={3}
                className="w-full bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500 resize-none"
              />
            </div>
          ))}

          <div className="flex gap-3">
            <button
              onClick={() => saveMutation.mutate()}
              disabled={saveMutation.isPending}
              className="bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-semibold px-6 py-2 rounded-lg transition"
            >
              {saveMutation.isPending ? 'Salvando...' : '💾 Salvar'}
            </button>
            {background && (
              <button
                onClick={() => setIsEditing(false)}
                className="text-gray-400 hover:text-white transition px-4 py-2"
              >
                Cancelar
              </button>
            )}
          </div>
        </div>
      ) : (
        <div className="flex flex-col gap-4">
          {fields.map(field => (
            <div key={field.key}>
              <p className="text-gray-400 text-sm mb-1">{field.label}</p>
              <p className="text-white text-sm bg-gray-700 rounded-lg px-4 py-3 max-h-24 overflow-y-auto">
                {form[field.key] || <span className="text-gray-500 italic">Não preenchido</span>}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}