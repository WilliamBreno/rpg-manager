import { useState } from 'react'
import type { Skill } from '../types'

export const powerConfig: Record<string, {
  label: string; color: string; border: string; bg: string; badge: string
}> = {
  unlimited: { label: '⚡ Sem Limite',    color: 'text-green-400',  border: 'border-green-700',  bg: 'bg-green-950',  badge: 'bg-green-900 text-green-300'  },
  encounter: { label: '⚔️ Por Encontro', color: 'text-yellow-400', border: 'border-yellow-700', bg: 'bg-yellow-950', badge: 'bg-yellow-900 text-yellow-300' },
  daily:     { label: '📅 Diário',        color: 'text-red-400',    border: 'border-red-700',    bg: 'bg-red-950',    badge: 'bg-red-900 text-red-300'      },
  utility:   { label: '🔧 Utilitário',    color: 'text-blue-400',   border: 'border-blue-700',   bg: 'bg-blue-950',   badge: 'bg-blue-900 text-blue-300'    },
}

interface SkillCardProps {
  skill: Skill
  // Modo seleção
  selectable?: boolean
  selected?: boolean
  disabled?: boolean
  onToggle?: (skill: Skill) => void
  // Modo informativo (característica automática)
  informative?: boolean
  // Sempre expandido (usado na criação para ver todos os detalhes)
  defaultExpanded?: boolean
}

export function SkillCard({
  skill, selectable, selected, disabled, onToggle, informative, defaultExpanded = false
}: SkillCardProps) {
  const [expanded, setExpanded] = useState(defaultExpanded)
  const cfg = powerConfig[skill.power_type ?? 'unlimited']

  const hasDetails = !!(skill.action_type || skill.range || skill.target || skill.attack ||
    skill.hit || skill.miss || skill.effect || skill.special || skill.level_scaling)

  const handleExpandClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    setExpanded(v => !v)
  }

  const handleCheckboxClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    if (selectable && onToggle && !disabled) onToggle(skill)
  }

  return (
    <div
      className={`rounded-xl border transition-all ${
        selectable
          ? selected
            ? `${cfg.border} ${cfg.bg} shadow-lg shadow-black/30`
            : disabled
            ? 'border-gray-700 bg-gray-800/50 opacity-40'
            : 'border-gray-600 bg-gray-700/80 hover:border-gray-500'
          : informative
          ? 'border-indigo-800 bg-indigo-950/80'
          : `${cfg.border} ${cfg.bg}`
      }`}
    >
      {/* ── CABEÇALHO ─────────────────────────────────────────── */}
      <div className="p-3 sm:p-4">
        <div className="flex items-start justify-between gap-2">
          <div className="flex-1 min-w-0">

            {/* Linha 1: nome + badge de tipo */}
            <div className="flex flex-wrap items-center gap-2 mb-1">
              <span className="text-white font-bold text-sm">{skill.name}</span>
              <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${cfg.badge}`}>
                {cfg.label}
              </span>
              {informative && (
                <span className="text-xs px-2 py-0.5 rounded-full bg-indigo-900 text-indigo-300">
                  Automático
                </span>
              )}
            </div>

            {/* Linha 2: palavras-chave */}
            {skill.keywords && (
              <p className="text-gray-400 text-xs mb-1">
                <span className={`font-semibold ${cfg.color}`}>
                  {powerConfig[skill.power_type ?? 'unlimited']?.label.replace(/^[^\s]+ /, '')}
                </span>
                {skill.keywords && ` • ${skill.keywords}`}
              </p>
            )}

            {/* Linha 3: descrição */}
            {skill.description && (
              <p className="text-gray-300 text-xs italic mt-1">{skill.description}</p>
            )}

          </div>

          {/* Botões direita */}
          <div className="flex items-center gap-1.5 shrink-0 mt-0.5">
            {selectable && (
              <button
                type="button"
                onClick={handleCheckboxClick}
                disabled={disabled && !selected}
                className={`w-6 h-6 rounded-full border-2 flex items-center justify-center text-xs font-bold transition-all ${
                  selected
                    ? `${cfg.badge} border-transparent scale-110`
                    : 'border-gray-500 text-transparent hover:border-gray-300'
                }`}
              >
                ✓
              </button>
            )}
            {hasDetails && (
              <button
                type="button"
                onClick={handleExpandClick}
                className="text-gray-500 hover:text-gray-300 text-xs px-1 transition-colors"
              >
                {expanded ? '▲' : '▼'}
              </button>
            )}
          </div>
        </div>

        {/* ── DETALHES EXPANDIDOS ──────────────────────────────── */}
        {expanded && hasDetails && (
          <div className="mt-3 pt-3 border-t border-white/10 grid grid-cols-1 gap-1.5">

            {/* Ação + Alcance na mesma linha se couber */}
            {(skill.action_type || skill.range) && (
              <div className="flex flex-wrap gap-x-4 gap-y-1">
                {skill.action_type && (
                  <div className="flex gap-2 text-xs">
                    <span className="text-gray-500 shrink-0">Ação:</span>
                    <span className="text-gray-200">{skill.action_type}</span>
                  </div>
                )}
                {skill.range && (
                  <div className="flex gap-2 text-xs">
                    <span className="text-gray-500 shrink-0">Alcance:</span>
                    <span className="text-gray-200">{skill.range}</span>
                  </div>
                )}
              </div>
            )}

            {skill.target  && <DetailRow label="Alvo"          value={skill.target} />}
            {skill.attack  && <DetailRow label="Ataque"        value={skill.attack}        color="text-indigo-300 font-semibold" />}
            {skill.hit     && <DetailRow label="Sucesso"       value={skill.hit}           color="text-green-300" />}
            {skill.miss    && <DetailRow label="Fracasso"      value={skill.miss}          color="text-red-300" />}
            {skill.effect  && <DetailRow label="Efeito"        value={skill.effect}        color="text-yellow-200" />}
            {skill.special && <DetailRow label="Especial"      value={skill.special}       color="text-purple-300" />}
            {skill.level_scaling && (
              <div className="mt-1 pt-1 border-t border-white/5">
                <DetailRow label="Escalonamento" value={skill.level_scaling} color="text-gray-400 italic" />
              </div>
            )}

          </div>
        )}
      </div>
    </div>
  )
}

function DetailRow({ label, value, color = 'text-gray-200' }: {
  label: string; value: string; color?: string
}) {
  return (
    <div className="flex gap-2 text-xs">
      <span className="text-gray-500 w-20 sm:w-24 shrink-0">{label}:</span>
      <span className={color}>{value}</span>
    </div>
  )
}