import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { characterService } from '../services/characterService'
import type { Character } from '../types'

interface Props {
  characterId: number
  character: Character
}

// Só renderiza quando HP = 0 em edição 5e
export default function DeathSaves({ characterId, character }: Props) {
  const queryClient = useQueryClient()
  const [feedback, setFeedback] = useState<{ message: string; dead: boolean } | null>(null)

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ['character', String(characterId)] })

  const saveMutation = useMutation({
    mutationFn: (body: { success: boolean; critical: boolean }) =>
      characterService.deathSave(characterId, body),
    onSuccess: (data) => {
      invalidate()
      if (data.stabilized && data.character?.hit_points > 0) {
        setFeedback({ message: '⭐ 20 Natural! Recuperou 1 HP e acordou.', dead: false })
        setTimeout(() => setFeedback(null), 6000)
      } else if (data.stabilized) {
        setFeedback({ message: '✦ Estabilizado — 3 sucessos. Continua inconsciente mas não morrerá.', dead: false })
        setTimeout(() => setFeedback(null), 8000)
      } else if (data.dead) {
        setFeedback({ message: '💀 3 falhas — o personagem morreu.', dead: true })
      }
    },
  })

  const resetMutation = useMutation({
    mutationFn: () => characterService.resetDeathSaves(characterId),
    onSuccess: () => { invalidate(); setFeedback(null) },
  })

  // Não mostra enquanto o personagem estiver de pé
  if (character.hit_points > 0) return null

  const successes = character.death_save_successes ?? 0
  const failures  = character.death_save_failures  ?? 0
  const isDead    = failures >= 3
  const isPending = saveMutation.isPending || resetMutation.isPending

  return (
    <div className="bg-gray-800 rounded-xl p-4 sm:p-6 mb-4 border"
      style={{ borderColor: isDead ? 'rgba(220,38,38,0.5)' : 'rgba(234,179,8,0.4)' }}>

      {/* Título */}
      <h2 className="text-sm font-semibold uppercase tracking-widest mb-3"
        style={{ color: isDead ? '#f87171' : '#facc15' }}>
        {isDead ? '💀 Personagem Morto' : '⚠️ Testes de Morte'}
      </h2>

      {!isDead && (
        <>
          <p className="text-gray-500 text-xs mb-4 leading-relaxed">
            No início do seu turno, role 1d20.{' '}
            <span className="text-green-400 font-medium">10+ = sucesso</span> ·{' '}
            <span className="text-red-400 font-medium">9 ou menos = falha</span> ·{' '}
            3 sucessos = estabilizado · 3 falhas = morto.
          </p>

          {/* Circulos de sucesso / falha */}
          <div className="flex flex-col gap-3 mb-5">
            {/* Sucessos */}
            <div className="flex items-center gap-4">
              <span className="text-xs text-gray-500 w-16 text-right uppercase tracking-wider shrink-0">Sucessos</span>
              <div className="flex gap-2.5">
                {[0, 1, 2].map(i => (
                  <div key={i}
                    className="w-8 h-8 rounded-full border-2 flex items-center justify-center transition"
                    style={i < successes
                      ? { background: '#16a34a', borderColor: '#16a34a' }
                      : { background: '#18181b', borderColor: '#3f3f46' }
                    }
                  >
                    {i < successes && <span className="text-white text-xs font-bold">✓</span>}
                  </div>
                ))}
              </div>
              <span className="text-green-600 text-xs font-mono">{successes}/3</span>
            </div>

            {/* Falhas */}
            <div className="flex items-center gap-4">
              <span className="text-xs text-gray-500 w-16 text-right uppercase tracking-wider shrink-0">Falhas</span>
              <div className="flex gap-2.5">
                {[0, 1, 2].map(i => (
                  <div key={i}
                    className="w-8 h-8 rounded-full border-2 flex items-center justify-center transition"
                    style={i < failures
                      ? { background: '#dc2626', borderColor: '#dc2626' }
                      : { background: '#18181b', borderColor: '#3f3f46' }
                    }
                  >
                    {i < failures && <span className="text-white text-xs font-bold">✗</span>}
                  </div>
                ))}
              </div>
              <span className="text-red-600 text-xs font-mono">{failures}/3</span>
            </div>
          </div>

          {/* Botões de resultado */}
          <div className="grid grid-cols-2 gap-2 mb-3">
            <button
              disabled={isPending || successes >= 3}
              onClick={() => saveMutation.mutate({ success: true, critical: false })}
              className="py-2.5 rounded-lg text-sm font-semibold transition border"
              style={{ background: 'rgba(22,163,74,0.12)', borderColor: 'rgba(22,163,74,0.35)', color: '#4ade80' }}
            >
              ✓ Sucesso
            </button>
            <button
              disabled={isPending || failures >= 3}
              onClick={() => saveMutation.mutate({ success: false, critical: false })}
              className="py-2.5 rounded-lg text-sm font-semibold transition border"
              style={{ background: 'rgba(220,38,38,0.12)', borderColor: 'rgba(220,38,38,0.35)', color: '#f87171' }}
            >
              ✗ Falha
            </button>
            <button
              disabled={isPending}
              onClick={() => saveMutation.mutate({ success: true, critical: true })}
              className="py-2 rounded-lg text-xs font-semibold transition border"
              style={{ background: 'rgba(234,179,8,0.1)', borderColor: 'rgba(234,179,8,0.3)', color: '#facc15' }}
            >
              ⭐ 20 Natural — Acorda com 1 HP
            </button>
            <button
              disabled={isPending || failures >= 3}
              onClick={() => saveMutation.mutate({ success: false, critical: true })}
              className="py-2 rounded-lg text-xs font-semibold transition border"
              style={{ background: 'rgba(220,38,38,0.07)', borderColor: 'rgba(220,38,38,0.2)', color: '#fca5a5' }}
            >
              💀 1 Natural — +2 falhas
            </button>
          </div>

          {/* Resetar */}
          <button
            disabled={isPending || (successes === 0 && failures === 0)}
            onClick={() => resetMutation.mutate()}
            className="w-full py-1.5 rounded-lg text-xs transition"
            style={{
              background: 'transparent',
              border: '1px solid #27272a',
              color: successes === 0 && failures === 0 ? '#3f3f46' : '#71717a',
            }}
          >
            Resetar testes de morte
          </button>
        </>
      )}

      {/* Se morto — só botão de reset (ressurreição) */}
      {isDead && (
        <div className="flex flex-col gap-3">
          <p className="text-red-400/70 text-xs">
            O personagem acumulou 3 falhas nos testes de morte.
            Use "Resetar" se uma magia de ressurreição for aplicada pelo mestre.
          </p>
          <button
            disabled={isPending}
            onClick={() => resetMutation.mutate()}
            className="w-full py-2 rounded-lg text-xs font-semibold transition border"
            style={{ background: 'rgba(99,102,241,0.1)', borderColor: 'rgba(99,102,241,0.3)', color: '#a5b4fc' }}
          >
            ✦ Ressurreição — Resetar testes
          </button>
        </div>
      )}

      {/* Feedback */}
      {feedback && (
        <div className="mt-3 text-xs text-center font-semibold rounded-lg p-2.5"
          style={feedback.dead
            ? { background: 'rgba(220,38,38,0.1)', border: '1px solid rgba(220,38,38,0.2)', color: '#f87171' }
            : { background: 'rgba(201,168,76,0.1)', border: '1px solid rgba(201,168,76,0.2)', color: '#c9a84c' }
          }
        >
          {feedback.message}
        </div>
      )}
    </div>
  )
}