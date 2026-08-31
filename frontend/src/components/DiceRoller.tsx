import { useState } from 'react'

const DICE_TYPES = [4, 6, 8, 10, 12, 20, 100] as const

interface RollEntry {
  id: number
  label: string
  rolls: number[]
  modifier: number
  total: number
}

// Roller de dados genérico — Sistema do Mestre, Etapa 7 ("mesmo tema
// visual/biblioteca usada pro dado do jogador... não criar um componente de
// dado paralelo com visual diferente"). Não existe ainda nenhum roller do
// lado do jogador no código (Unlumen UI/Magic UI mencionados no
// TASKS_UI_E_FEATURES.md são só um objetivo de Fase 5, sem componente
// construído) — pra não criar dois componentes de dado com visual diferente
// mais tarde, este é escrito como o único componente de dado do projeto,
// reaproveitável por qualquer tela (mestre agora, jogador quando existir).
export default function DiceRoller({ onShare }: { onShare?: (text: string) => void }) {
  const [quantity, setQuantity] = useState(1)
  const [modifier, setModifier] = useState(0)
  const [rolling, setRolling] = useState<number | null>(null)
  const [history, setHistory] = useState<RollEntry[]>([])

  const roll = (sides: number) => {
    setRolling(sides)
    setTimeout(() => {
      const rolls = Array.from({ length: quantity }, () => 1 + Math.floor(Math.random() * sides))
      const total = rolls.reduce((a, b) => a + b, 0) + modifier
      const entry: RollEntry = {
        id: Date.now(),
        label: `${quantity}d${sides}${modifier !== 0 ? (modifier > 0 ? `+${modifier}` : modifier) : ''}`,
        rolls, modifier, total,
      }
      setHistory(prev => [entry, ...prev].slice(0, 20))
      setRolling(null)
    }, 350)
  }

  return (
    <div className="rpg-card p-4">
      <div className="flex flex-wrap items-center gap-3 mb-4">
        <div className="flex items-center gap-1">
          <label className="text-gray-500 text-xs uppercase tracking-wider">Qtd</label>
          <input type="number" min={1} max={20} value={quantity} onChange={e => setQuantity(Math.max(1, Number(e.target.value) || 1))}
            className="rpg-input w-16 text-sm py-1" />
        </div>
        <div className="flex items-center gap-1">
          <label className="text-gray-500 text-xs uppercase tracking-wider">Mod</label>
          <input type="number" value={modifier} onChange={e => setModifier(Number(e.target.value) || 0)}
            className="rpg-input w-16 text-sm py-1" />
        </div>
      </div>

      <div className="flex flex-wrap gap-2 mb-4">
        {DICE_TYPES.map(sides => (
          <button
            key={sides}
            onClick={() => roll(sides)}
            className={`w-16 h-16 rounded-xl border-2 flex items-center justify-center font-rpg font-bold text-lg transition
              ${rolling === sides ? 'animate-spin border-rpg-gold text-rpg-gold' : 'border-gray-700 text-gray-300 hover:border-rpg-gold hover:text-rpg-gold'}`}
          >
            d{sides}
          </button>
        ))}
      </div>

      <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
        {history.length === 0 && <p className="text-gray-600 text-sm">Nenhuma rolagem ainda — clique num dado acima.</p>}
        {history.map(h => (
          <div key={h.id} className="flex items-center justify-between bg-gray-800/40 rounded-lg px-3 py-2">
            <div>
              <span className="text-gray-400 text-xs">{h.label}</span>
              <span className="text-gray-500 text-xs ml-2">[{h.rolls.join(', ')}]</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-rpg-gold font-bold text-lg">{h.total}</span>
              {onShare && (
                <button onClick={() => onShare(`🎲 rolou ${h.label} = ${h.total} [${h.rolls.join(', ')}]`)}
                  className="text-xs text-gray-500 hover:text-rpg-gold transition">enviar ao chat</button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
