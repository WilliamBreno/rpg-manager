import { useState, useRef, useCallback, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { characterService } from '../services/characterService'
import type { Character } from '../types'
import './HPManager.css'

const DICE_TYPES = [4, 6, 8, 10, 12, 20] as const

// Posição final que mostra a face frontal claramente
const SETTLE = { x: -15, y: 15, z: 0 }

interface Props { character: Character }

export default function HPManager({ character }: Props) {
  const queryClient = useQueryClient()

  // ── HP state ────────────────────────────────────────────────────────────────
  const [damageInput, setDamageInput] = useState('')
  const [healInput, setHealInput]     = useState('')
  const [tempInput, setTempInput]     = useState('')
  const [mode, setMode] = useState<'damage' | 'heal' | 'temp' | null>(null)

  // ── Dice state ──────────────────────────────────────────────────────────────
  const [showDice, setShowDice]       = useState(false)
  const [numDice, setNumDice]         = useState(1)
  const [selectedDie, setSelectedDie] = useState(20)
  const [isRolling, setIsRolling]     = useState(false)
  const [rollResult, setRollResult]   = useState<{ rolls: number[]; total: number } | null>(null)
  const [rollHistory, setRollHistory] = useState<string[]>([])

  // ── Refs ────────────────────────────────────────────────────────────────────
  const diceRef      = useRef<HTMLDivElement>(null)
  const wrapperRef   = useRef<HTMLDivElement>(null)
  const diceAngleRef = useRef({ ...SETTLE })
  const animFrameRef = useRef<number | null>(null)
  const tickerRef    = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    return () => {
      if (animFrameRef.current) cancelAnimationFrame(animFrameRef.current)
      if (tickerRef.current) clearInterval(tickerRef.current)
    }
  }, [])

  // ── Mutations ───────────────────────────────────────────────────────────────
  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: ['character', String(character.ID)] })

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

  // ── HP helpers ──────────────────────────────────────────────────────────────
  const maxHP      = character.max_hp || character.hit_points
  const currentHP  = character.hit_points
  const tempHP     = character.temp_hp || 0
  const hpPercent  = Math.max(0, Math.min(100, (currentHP / maxHP) * 100))
  const isCritical    = hpPercent <= 25 && currentHP > 0
  const isUnconscious = currentHP === 0

  const getHPColor = () => {
    if (isUnconscious)   return 'bg-red-600'
    if (hpPercent <= 25) return 'bg-orange-500'
    if (hpPercent <= 50) return 'bg-yellow-500'
    return 'bg-green-500'
  }
  const getHPStatus = () => {
    if (isUnconscious)   return { icon: '💀', label: 'Inconsciente', color: 'text-red-400',    pulse: false }
    if (hpPercent <= 25) return { icon: '🩸', label: 'Crítico',      color: 'text-orange-400', pulse: false }
    if (hpPercent <= 50) return { icon: '⚠️', label: 'Ferido',       color: 'text-yellow-400', pulse: true  }
    return                      { icon: '✅', label: 'Saudável',     color: 'text-green-400',  pulse: false }
  }
  const status = getHPStatus()

  // ── Partículas douradas ─────────────────────────────────────────────────────
  const spawnParticles = useCallback((die: number) => {
    if (!wrapperRef.current) return
    const w = wrapperRef.current
    const count = die === 20 ? 16 : 8
    for (let i = 0; i < count; i++) {
      const p = document.createElement('div')
      p.className = 'rpg-particle'
      const a = (i / count) * Math.PI * 2 + Math.random() * .4
      const dist = 48 + Math.random() * 30
      w.appendChild(p)
      requestAnimationFrame(() => requestAnimationFrame(() => {
        p.style.transform = `translate(calc(-50% + ${Math.cos(a) * dist}px), calc(-50% + ${Math.sin(a) * dist}px))`
        p.style.opacity = '0'
      }))
      setTimeout(() => p.remove(), 900)
    }
  }, [])

  // ── Roll animation ──────────────────────────────────────────────────────────
  const rollDice = useCallback(() => {
    if (isRolling || !diceRef.current) return
    setIsRolling(true)
    setRollResult(null)

    // Pré-calcula resultado antes de começar a animação
    const rolls = Array.from({ length: numDice }, () =>
      Math.floor(Math.random() * selectedDie) + 1
    )
    const total = rolls.reduce((a, b) => a + b, 0)

    const dice  = diceRef.current
    dice.style.transition = ''
    const faces = dice.querySelectorAll<HTMLDivElement>('.rpg-face')

    // Destino = posição frontal + N rotações completas em direção aleatória
    // O dado naturalmente termina mostrando a face frontal — sem transição extra
    const spinsX = 4 + Math.floor(Math.random() * 4)
    const spinsY = 4 + Math.floor(Math.random() * 4)
    const spinsZ = 2 + Math.floor(Math.random() * 3)
    const dirX   = Math.random() < .5 ? 1 : -1
    const dirY   = Math.random() < .5 ? 1 : -1
    const dirZ   = Math.random() < .5 ? 1 : -1

    const target = {
      x: SETTLE.x + dirX * spinsX * 360,
      y: SETTLE.y + dirY * spinsY * 360,
      z: SETTLE.z + dirZ * spinsZ * 360,
    }

    const start = { ...diceAngleRef.current }
    const dur   = 1200 + Math.floor(Math.random() * 400) // 1200-1600ms
    const t0    = performance.now()
    let resultPlaced = false

    // Faces piscam números aleatórios enquanto gira rápido
    tickerRef.current = setInterval(() => {
      faces.forEach(f => { f.textContent = String(Math.floor(Math.random() * selectedDie) + 1) })
    }, 65)

    // easeOutCubic: desaceleração física natural, sem bounce
    const easeOut = (t: number) => 1 - Math.pow(1 - t, 3)

    function frame(now: number) {
      const t = Math.min((now - t0) / dur, 1)
      const e = easeOut(t)

      if (diceRef.current) {
        diceRef.current.style.transform =
          `rotateX(${start.x + (target.x - start.x) * e}deg) ` +
          `rotateY(${start.y + (target.y - start.y) * e}deg) ` +
          `rotateZ(${start.z + (target.z - start.z) * e}deg)`
      }

      // A 75% do tempo o dado está quase parado (97% da rotação completa)
      // Coloca o resultado na face frontal — o dado ainda se move levemente,
      // então a troca de número não é percebida
      if (t >= 0.75 && !resultPlaced) {
        resultPlaced = true
        if (tickerRef.current) { clearInterval(tickerRef.current); tickerRef.current = null }
        const front = diceRef.current?.querySelector<HTMLDivElement>('.rpg-face-front')
        if (front) front.textContent = numDice === 1 ? String(total) : '∑'
      }

      if (t < 1) {
        animFrameRef.current = requestAnimationFrame(frame)
      } else {
        // Garante posição final exata
        if (diceRef.current) {
          diceRef.current.style.transform =
            `rotateX(${SETTLE.x}deg) rotateY(${SETTLE.y}deg) rotateZ(${SETTLE.z}deg)`
        }
        diceAngleRef.current = { ...SETTLE }

        // Garante resultado na face frontal
        const front = diceRef.current?.querySelector<HTMLDivElement>('.rpg-face-front')
        if (front) front.textContent = numDice === 1 ? String(total) : '∑'

        setRollResult({ rolls, total })
        const entry = numDice === 1
          ? `${total} (d${selectedDie})`
          : `${total} (${numDice}d${selectedDie})`
        setRollHistory(prev => [entry, ...prev].slice(0, 6))

        const isCrit = numDice === 1 && selectedDie === 20 && total === 20
        const isMax  = total === selectedDie * numDice
        if (isCrit || isMax) spawnParticles(selectedDie)

        setIsRolling(false)
      }
    }
    animFrameRef.current = requestAnimationFrame(frame)
  }, [isRolling, numDice, selectedDie, spawnParticles])

  // ── Input Row helper ────────────────────────────────────────────────────────
  const InputRow = ({
    value, onChange, placeholder, onApply, isPending, color,
  }: {
    value: string
    onChange: (v: string) => void
    placeholder: string
    onApply: () => void
    isPending: boolean
    color: 'red' | 'green' | 'blue'
  }) => {
    const map = {
      red:   { ring: 'focus:ring-red-500',   btn: 'bg-red-600   hover:bg-red-700'   },
      green: { ring: 'focus:ring-green-500', btn: 'bg-green-600 hover:bg-green-700' },
      blue:  { ring: 'focus:ring-blue-500',  btn: 'bg-blue-600  hover:bg-blue-700'  },
    }
    return (
      <div className="flex flex-col sm:flex-row gap-2">
        <input
          type="number" value={value} min={1}
          onChange={e => onChange(e.target.value)}
          placeholder={placeholder}
          className={`flex-1 bg-gray-700 text-white rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 ${map[color].ring}`}
        />
        <button
          onClick={onApply}
          disabled={!value || isPending}
          className={`${map[color].btn} disabled:opacity-50 text-white font-semibold px-5 py-2 rounded-lg transition text-sm whitespace-nowrap`}
        >
          {isPending ? '...' : 'Aplicar'}
        </button>
      </div>
    )
  }

  // ── Render ──────────────────────────────────────────────────────────────────
  return (
    <div className={`bg-gray-800 rounded-lg p-4 sm:p-6 mb-4 border transition-all duration-300 ${
      isUnconscious ? 'border-red-700' : isCritical ? 'border-orange-700' : 'border-gray-700'
    }`}>

      <h2 className="text-base sm:text-lg font-semibold text-white mb-3">Pontos de Vida</h2>

      <div className="flex justify-between items-center mb-2">
        <span className={`text-sm font-semibold flex items-center gap-1 ${status.color}`}>
          <span className={status.pulse ? 'animate-pulse' : ''}>{status.icon}</span>
          {status.label}
        </span>
        <div className="flex gap-2 text-sm">
          <span className="text-white font-bold">{currentHP} / {maxHP} HP</span>
          {tempHP > 0 && <span className="text-blue-400 font-bold">+{tempHP}</span>}
        </div>
      </div>

      <div className="w-full bg-gray-700 rounded-full h-3 sm:h-4 mb-2 overflow-hidden">
        <div className={`h-full rounded-full transition-all duration-500 ${getHPColor()}`} style={{ width: `${hpPercent}%` }} />
      </div>

      {tempHP > 0 && (
        <div className="w-full bg-gray-700 rounded-full h-2 mb-3 overflow-hidden">
          <div className="h-2 rounded-full bg-blue-500 transition-all duration-500" style={{ width: `${Math.min(100, (tempHP / maxHP) * 100)}%` }} />
        </div>
      )}

      {isCritical && <p className="text-orange-400 text-xs text-center mb-3 font-semibold">⚠️ HP crítico — procure cura imediatamente!</p>}
      {isUnconscious && <p className="text-red-400 text-xs text-center mb-3 font-semibold animate-pulse">💀 Inconsciente — salvaguardas de morte!</p>}

      <div className="grid grid-cols-3 gap-2 mb-3">
        <button onClick={() => setMode(mode === 'damage' ? null : 'damage')}
          className={`py-2 rounded-lg font-semibold transition flex flex-col sm:flex-row items-center justify-center gap-1 ${mode === 'damage' ? 'bg-red-600 text-white' : 'bg-gray-700 text-red-400 hover:bg-gray-600'}`}>
          <span>⚔️</span>
          <span className="text-xs sm:text-sm leading-tight"><span className="hidden sm:inline">Receber </span>Dano</span>
        </button>
        <button onClick={() => setMode(mode === 'heal' ? null : 'heal')}
          className={`py-2 rounded-lg font-semibold transition flex flex-col sm:flex-row items-center justify-center gap-1 ${mode === 'heal' ? 'bg-green-600 text-white' : 'bg-gray-700 text-green-400 hover:bg-gray-600'}`}>
          <span>💚</span>
          <span className="text-xs sm:text-sm">Curar</span>
        </button>
        <button onClick={() => setMode(mode === 'temp' ? null : 'temp')}
          className={`py-2 rounded-lg font-semibold transition flex flex-col sm:flex-row items-center justify-center gap-1 ${mode === 'temp' ? 'bg-blue-600 text-white' : 'bg-gray-700 text-blue-400 hover:bg-gray-600'}`}>
          <span>🛡️</span>
          <span className="text-xs sm:text-sm"><span className="hidden sm:inline">HP </span>Temp</span>
        </button>
      </div>

      {mode === 'damage' && (
        <InputRow value={damageInput} onChange={setDamageInput} placeholder="Quantidade de dano"
          onApply={() => damageMutation.mutate(Number(damageInput))} isPending={damageMutation.isPending} color="red" />
      )}
      {mode === 'heal' && (
        <InputRow value={healInput} onChange={setHealInput} placeholder="Quantidade de cura"
          onApply={() => healMutation.mutate(Number(healInput))} isPending={healMutation.isPending} color="green" />
      )}
      {mode === 'temp' && (
        <div className="flex flex-col gap-1">
          <InputRow value={tempInput} onChange={setTempInput} placeholder="HP temporário"
            onApply={() => tempMutation.mutate(Number(tempInput))} isPending={tempMutation.isPending} color="blue" />
          <p className="text-gray-500 text-xs">HP temporário não acumula — fica com o maior valor.</p>
        </div>
      )}

      {/* ── Simulador de Dados ─────────────────────────────────────────────── */}
      <div className="border-t border-gray-700 mt-4 pt-4">

        <button type="button" onClick={() => setShowDice(v => !v)} className={`rpg-toggle-btn ${showDice ? 'open' : ''}`}>
          <span style={{ fontSize: 16 }}>✦</span>
          <span>Simulador de Dados</span>
          <span className="ml-auto text-xs" style={{ color: 'rgba(201,168,76,.35)' }}>{showDice ? '▲' : '▼'}</span>
        </button>

        {showDice && (
          <div className="mt-5 flex flex-col items-center gap-4">

            {/* Cena 3D */}
            <div ref={wrapperRef} style={{ width: 148, height: 148, position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <div style={{ position: 'absolute', inset: 0, borderRadius: '50%', background: '#060606', border: '1.5px solid rgba(201,168,76,.55)', pointerEvents: 'none' }} />
              <div style={{ position: 'absolute', inset: 10, borderRadius: '50%', border: '.5px solid rgba(201,168,76,.18)', pointerEvents: 'none' }} />
              <div style={{ position: 'absolute', inset: 20, borderRadius: '50%', border: '.5px dashed rgba(201,168,76,.1)', pointerEvents: 'none' }} />
              <div className="rpg-star rpg-star-top" />
              <div className="rpg-star rpg-star-bottom" />
              <div className="rpg-star rpg-star-left" />
              <div className="rpg-star rpg-star-right" />
              <div style={{ width: 96, height: 96, perspective: '440px', cursor: 'pointer', position: 'relative', zIndex: 2 }} onClick={rollDice}>
                <div ref={diceRef} className="rpg-dice">
                  <div className="rpg-face rpg-face-front">20</div>
                  <div className="rpg-face rpg-face-back">1</div>
                  <div className="rpg-face rpg-face-right">6</div>
                  <div className="rpg-face rpg-face-left">14</div>
                  <div className="rpg-face rpg-face-top">3</div>
                  <div className="rpg-face rpg-face-bottom">17</div>
                </div>
              </div>
            </div>

            {/* Controles */}
            <div className="w-full flex flex-col gap-3">
              <div className="flex items-center justify-center gap-3">
                <span style={{ color: 'rgba(201,168,76,.6)', fontSize: 13 }}>Qtd:</span>
                <button type="button" className="rpg-qty-btn" onClick={() => setNumDice(v => Math.max(1, v - 1))}>−</button>
                <span style={{ color: '#c9a84c' }} className="font-bold w-6 text-center text-lg tabular-nums">{numDice}</span>
                <button type="button" className="rpg-qty-btn" onClick={() => setNumDice(v => Math.min(10, v + 1))}>+</button>
              </div>
              <div className="grid grid-cols-3 gap-1.5">
                {DICE_TYPES.map(d => (
                  <button key={d} type="button" onClick={() => setSelectedDie(d)}
                    className={`rpg-die-btn ${selectedDie === d ? 'active' : ''}`}>
                    d{d}
                  </button>
                ))}
              </div>
              <button type="button" className="rpg-roll-btn" onClick={rollDice} disabled={isRolling}>
                {isRolling ? '✦ Rolando...' : `✦ Rolar ${numDice > 1 ? numDice : ''}d${selectedDie}`}
              </button>
            </div>

            {/* Resultado */}
            {(isRolling || rollResult) && (
              <div className="rpg-result-box">
                {isRolling ? (
                  <p className="rpg-result-number animate-pulse" style={{ fontSize: 40 }}>✦</p>
                ) : rollResult && (
                  <>
                    <p className="rpg-result-number">{rollResult.total}</p>
                    {numDice > 1 && <p className="rpg-result-detail">[ {rollResult.rolls.join(' + ')} ]</p>}
                    <p className="rpg-result-detail flex items-center gap-2">
                      <span>{numDice}d{selectedDie}</span>
                      {numDice === 1 && selectedDie === 20 && rollResult.total === 20 && <span className="rpg-badge rpg-badge-crit">⚡ Crítico natural!</span>}
                      {numDice === 1 && selectedDie === 20 && rollResult.total === 1  && <span className="rpg-badge rpg-badge-fail">💀 Falha crítica!</span>}
                      {rollResult.total === selectedDie * numDice && numDice > 1        && <span className="rpg-badge rpg-badge-max">✦ Máximo!</span>}
                    </p>
                  </>
                )}
              </div>
            )}

            {/* Histórico */}
            {rollHistory.length > 0 && (
              <div className="w-full flex flex-col gap-2">
                <div className="flex justify-between items-center">
                  <span style={{ color: 'rgba(201,168,76,.45)', letterSpacing: '.08em' }} className="text-xs font-medium uppercase">Histórico</span>
                  <button type="button" onClick={() => setRollHistory([])}
                    style={{ color: 'rgba(201,168,76,.35)' }} className="text-xs bg-transparent border-none cursor-pointer hover:opacity-70">limpar</button>
                </div>
                <div className="flex flex-wrap gap-1.5">
                  {rollHistory.map((entry, i) => (
                    <span key={i} className={`rpg-chip ${i === 0 ? 'latest' : ''}`}>{entry}</span>
                  ))}
                </div>
              </div>
            )}

          </div>
        )}
      </div>
    </div>
  )
}