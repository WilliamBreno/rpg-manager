import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { characterService } from '../services/characterService'
import type { Character } from '../types'

interface Props {
  character: Character
}

export default function HPManager({ character }: Props) {
  const queryClient = useQueryClient()
  const [damageInput, setDamageInput] = useState('')
  const [healInput, setHealInput] = useState('')
  const [tempInput, setTempInput] = useState('')
  const [mode, setMode] = useState<'damage' | 'heal' | 'temp' | null>(null)

  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['character', String(character.ID)] })
  }

  const damageMutation = useMutation({
    mutationFn: (dmg: number) => characterService.takeDamage(character.ID, dmg),
    onSuccess: () => { invalidate(); setDamageInput(''); setMode(null) },
  })

  const healMutation = useMutation({
    mutationFn: (amount: number) => characterService.heal(character.ID, amount),
    onSuccess: () => { invalidate(); setHealInput(''); setMode(null) },
  })

  const tempMutation = useMutation({
    mutationFn: (amount: number) => characterService.addTempHP(character.ID, amount),
    onSuccess: () => { invalidate(); setTempInput(''); setMode(null) },
  })

  const maxHP = character.max_hp || character.hit_points
  const currentHP = character.hit_points
  const tempHP = character.temp_hp || 0
  const hpPercent = Math.max(0, Math.min(100, (currentHP / maxHP) * 100))

  const isCritical = hpPercent <= 25 && currentHP > 0
  const isUnconscious = currentHP === 0

  const getHPColor = () => {
    if (isUnconscious) return 'bg-red-600'
    if (hpPercent <= 25) return 'bg-orange-500'
    if (hpPercent <= 50) return 'bg-yellow-500'
    return 'bg-green-500'
  }

  const getHPStatus = () => {
    if (isUnconscious) return { icon: '💀', label: 'Inconsciente', color: 'text-red-400', pulse: false }
    if (hpPercent <= 25) return { icon: '🩸', label: 'Crítico', color: 'text-orange-400', pulse: false }
    if (hpPercent <= 50) return { icon: '⚠️', label: 'Ferido', color: 'text-yellow-400', pulse: true }
    return { icon: '✅', label: 'Saudável', color: 'text-green-400', pulse: false }
  }

  const status = getHPStatus()

  return (
    <div className={`bg-gray-800 rounded-lg p-6 mb-4 border transition-all duration-300 ${
      isUnconscious
        ? 'border-red-700 animate-pulse-red'
        : isCritical
        ? 'border-orange-700 animate-shake'
        : 'border-gray-700'
    }`}>
      <h2 className="text-lg font-semibold text-white mb-4">Pontos de Vida</h2>

      {/* Status */}
      <div className="flex justify-between items-center mb-2">
        <span className={`text-sm font-semibold flex items-center gap-1 ${status.color}`}>
          <span className={status.pulse ? 'animate-pulse' : ''}>{status.icon}</span>
          {status.label}
        </span>
        <div className="flex gap-3 text-sm">
          <span className="text-white font-bold">{currentHP} / {maxHP} HP</span>
          {tempHP > 0 && (
            <span className="text-blue-400 font-bold">+{tempHP} temp</span>
          )}
        </div>
      </div>

      {/* Barra de HP */}
      <div className="w-full bg-gray-700 rounded-full h-4 mb-2 overflow-hidden">
        <div
          className={`h-4 rounded-full transition-all duration-500 ${getHPColor()}`}
          style={{ width: `${hpPercent}%` }}
        />
      </div>

      {/* Barra de Temp HP */}
      {tempHP > 0 && (
        <div className="w-full bg-gray-700 rounded-full h-2 mb-4 overflow-hidden">
          <div
            className="h-2 rounded-full bg-blue-500 transition-all duration-500"
            style={{ width: `${Math.min(100, (tempHP / maxHP) * 100)}%` }}
          />
        </div>
      )}

      {/* Aviso crítico */}
      {isCritical && (
        <p className="text-orange-400 text-xs text-center mb-3 font-semibold">
          ⚠️ PERIGO! HP crítico — procure cura imediatamente!
        </p>
      )}

      {/* Aviso inconsciente */}
      {isUnconscious && (
        <p className="text-red-400 text-xs text-center mb-3 font-semibold animate-pulse">
          💀 Personagem inconsciente — salvaguardas de morte!
        </p>
      )}

      {/* Botões de ação */}
      <div className="flex gap-2 mb-4">
        <button
          onClick={() => setMode(mode === 'damage' ? null : 'damage')}
          className={`flex-1 py-2 rounded-lg text-sm font-semibold transition ${
            mode === 'damage'
              ? 'bg-red-600 text-white'
              : 'bg-gray-700 text-red-400 hover:bg-gray-600'
          }`}
        >
          ⚔️ Receber Dano
        </button>
        <button
          onClick={() => setMode(mode === 'heal' ? null : 'heal')}
          className={`flex-1 py-2 rounded-lg text-sm font-semibold transition ${
            mode === 'heal'
              ? 'bg-green-600 text-white'
              : 'bg-gray-700 text-green-400 hover:bg-gray-600'
          }`}
        >
          💚 Curar
        </button>
        <button
          onClick={() => setMode(mode === 'temp' ? null : 'temp')}
          className={`flex-1 py-2 rounded-lg text-sm font-semibold transition ${
            mode === 'temp'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-700 text-blue-400 hover:bg-gray-600'
          }`}
        >
          🛡 HP Temp
        </button>
      </div>

      {/* Input de dano */}
      {mode === 'damage' && (
        <div className="flex gap-2">
          <input
            type="number"
            value={damageInput}
            onChange={e => setDamageInput(e.target.value)}
            placeholder="Quantidade de dano"
            min={1}
            className="flex-1 bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-red-500"
          />
          <button
            onClick={() => damageMutation.mutate(Number(damageInput))}
            disabled={!damageInput || damageMutation.isPending}
            className="bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-lg transition"
          >
            {damageMutation.isPending ? '...' : 'Aplicar'}
          </button>
        </div>
      )}

      {/* Input de cura */}
      {mode === 'heal' && (
        <div className="flex gap-2">
          <input
            type="number"
            value={healInput}
            onChange={e => setHealInput(e.target.value)}
            placeholder="Quantidade de cura"
            min={1}
            className="flex-1 bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
          />
          <button
            onClick={() => healMutation.mutate(Number(healInput))}
            disabled={!healInput || healMutation.isPending}
            className="bg-green-600 hover:bg-green-700 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-lg transition"
          >
            {healMutation.isPending ? '...' : 'Aplicar'}
          </button>
        </div>
      )}

      {/* Input de HP temporário */}
      {mode === 'temp' && (
        <div>
          <div className="flex gap-2">
            <input
              type="number"
              value={tempInput}
              onChange={e => setTempInput(e.target.value)}
              placeholder="HP temporário"
              min={1}
              className="flex-1 bg-gray-700 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={() => tempMutation.mutate(Number(tempInput))}
              disabled={!tempInput || tempMutation.isPending}
              className="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold px-4 py-2 rounded-lg transition"
            >
              {tempMutation.isPending ? '...' : 'Aplicar'}
            </button>
          </div>
          <p className="text-gray-400 text-xs mt-1">HP temporário não acumula — fica com o maior valor.</p>
        </div>
      )}
    </div>
  )
}