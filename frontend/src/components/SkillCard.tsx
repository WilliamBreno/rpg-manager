import { useState } from 'react'
import type { Skill, PowerType } from '../types'

export const powerConfig: Record<string, { label: string; color: string; border: string; bg: string; badge: string }> = {
  unlimited: { label: '⚡ Sem Limite',    color: 'text-green-400',  border: 'border-green-800',  bg: 'bg-green-950',  badge: 'bg-green-900 text-green-300'  },
  encounter: { label: '⚔️ Por Encontro', color: 'text-yellow-400', border: 'border-yellow-800', bg: 'bg-yellow-950', badge: 'bg-yellow-900 text-yellow-300' },
  daily:     { label: '📅 Diário',        color: 'text-red-400',    border: 'border-red-800',    bg: 'bg-red-950',    badge: 'bg-red-900 text-red-300'      },
  utility:   { label: '🔧 Utilitário',    color: 'text-blue-400',   border: 'border-blue-800',   bg: 'bg-blue-950',   badge: 'bg-blue-900 text-blue-300'    },
}

interface SkillCardProps {
  skill: Skill
  // Modo seleção
  selectable?: boolean
  selected?: boolean
  disabled?: boolean
  onToggle?: (skill: Skill) => void
  // Modo informativo (característica de classe automática)
  informative?: boolean
}

export function SkillCard({ skill, selectable, selected, disabled, onToggle, informative }: SkillCardProps) {
  const [expanded, setExpanded] = useState(false)
  const cfg = powerConfig[skill.power_type ?? 'unlimited']

  const hasDetails = skill.action_type || skill.range || skill.target || skill.attack ||
    skill.hit || skill.miss || skill.effect || skill.special || skill.level_scaling

  const handleClick = () => {
    if (selectable && onToggle && !disabled) {
      onToggle(skill)
    } else {
      setExpanded(e => !e)
    }
  }

  return (
    <div
      onClick={handleClick}
      className={`rounded-lg border p-4 transition-all ${
        selectable
          ? selected
            ? `${cfg.border} ${cfg.bg} ring-2 ring-offset-1 ring-offset-gray-800 cursor-pointer`
            : disabled
            ? 'border-gray-700 bg-gray-800 opacity-40 cursor-not-allowed'
            : 'border-gray-600 bg-gray-700 hover:border-gray-500 hover:bg-gray-600 cursor-pointer'
          : informative
          ? 'border-indigo-900 bg-indigo-950 cursor-pointer'
          : `${cfg.border} ${cfg.bg} cursor-pointer`
      }`}
    >
      {/* Header */}
      <div className="flex justify-between items-start gap-2">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-white font-semibold text-sm">{skill.name}</span>
            {skill.keywords && (
              <span className="text-gray-400 text-xs">• {skill.keywords}</span>
            )}
            {informative && (
              <span className="text-xs px-2 py-0.5 rounded-full bg-indigo-900 text-indigo-300">
                Automático
              </span>
            )}
          </div>
          {skill.description && (
            <p className="text-gray-400 text-xs mt-1 italic line-clamp-2">{skill.description}</p>
          )}
        </div>
        <div className="flex items-center gap-2 shrink-0">
          {selectable && (
            <span className={`w-5 h-5 rounded-full border flex items-center justify-center text-xs ${
              selected ? `${cfg.badge} border-transparent` : 'border-gray-500 text-gray-500'
            }`}>
              {selected ? '✓' : ''}
            </span>
          )}
          {hasDetails && (
            <span className="text-gray-500 text-xs">{expanded ? '▲' : '▼'}</span>
          )}
        </div>
      </div>

      {/* Detalhes expandidos */}
      {expanded && hasDetails && (
        <div className="mt-3 pt-3 border-t border-gray-700 flex flex-col gap-1.5">
          {skill.action_type && <DetailRow label="Ação" value={skill.action_type} />}
          {skill.range      && <DetailRow label="Alcance" value={skill.range} />}
          {skill.target     && <DetailRow label="Alvo" value={skill.target} />}
          {skill.attack     && <DetailRow label="Ataque" value={skill.attack} color="text-indigo-300" />}
          {skill.hit        && <DetailRow label="Sucesso" value={skill.hit} color="text-green-300" />}
          {skill.miss       && <DetailRow label="Fracasso" value={skill.miss} color="text-red-300" />}
          {skill.effect     && <DetailRow label="Efeito" value={skill.effect} color="text-yellow-200" />}
          {skill.special    && <DetailRow label="Especial" value={skill.special} color="text-purple-300" />}
          {skill.level_scaling && <DetailRow label="Escalonamento" value={skill.level_scaling} color="text-gray-300" />}
        </div>
      )}
    </div>
  )
}

function DetailRow({ label, value, color = 'text-white' }: { label: string; value: string; color?: string }) {
  return (
    <div className="flex gap-2 text-xs">
      <span className="text-gray-400 w-24 shrink-0">{label}:</span>
      <span className={color}>{value}</span>
    </div>
  )
}