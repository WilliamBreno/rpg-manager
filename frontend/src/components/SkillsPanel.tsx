import { useEffect, useState } from 'react'
import type { Skill } from '../types'
import { skillService } from '../services/skillService'

interface Props {
  characterID: number
  classID: number
  raceID: number
  className: string
  edition: string
  level: number
}

const powerConfig = {
  unlimited: {
    label: '⚡ Sem Limite',
    color: 'text-green-400',
    border: 'border-green-800',
    bg: 'bg-green-950',
    badge: 'bg-green-900 text-green-300',
  },
  encounter: {
    label: '⚔️ Por Encontro',
    color: 'text-yellow-400',
    border: 'border-yellow-800',
    bg: 'bg-yellow-950',
    badge: 'bg-yellow-900 text-yellow-300',
  },
  daily: {
    label: '📅 Diário',
    color: 'text-red-400',
    border: 'border-red-800',
    bg: 'bg-red-950',
    badge: 'bg-red-900 text-red-300',
  },
  utility: {
    label: '🔧 Utilitário',
    color: 'text-blue-400',
    border: 'border-blue-800',
    bg: 'bg-blue-950',
    badge: 'bg-blue-900 text-blue-300',
  },
}

function SkillCard({ skill }: { skill: Skill }) {
  const [expanded, setExpanded] = useState(false)
  const config = skill.power_type
    ? powerConfig[skill.power_type]
    : powerConfig['utility']

  return (
    <div
      className={`rounded-lg border ${config.border} ${config.bg} p-4 cursor-pointer transition-all duration-200`}
      onClick={() => setExpanded(!expanded)}
    >
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-2">
          <h3 className="text-white font-semibold text-sm">{skill.name}</h3>
          {skill.level && skill.level > 1 && (
            <span className="text-gray-500 text-xs">Nível {skill.level}</span>
          )}
        </div>
        <div className="flex items-center gap-2">
          {skill.power_type && (
            <span className={`text-xs px-2 py-1 rounded-full ${config.badge}`}>
              {config.label}
            </span>
          )}
          {skill.description && (
            <span className="text-gray-400 text-xs">{expanded ? '▲' : '▼'}</span>
          )}
        </div>
      </div>

      {expanded && skill.description && (
        <p className="mt-3 text-gray-300 text-xs border-t border-gray-700 pt-3">
          {skill.description}
        </p>
      )}
    </div>
  )
}

export default function SkillsPanel({ classID, raceID, edition, level }: Props) {
  const [skills, setSkills] = useState<Skill[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  useEffect(() => {
    const load = async () => {
      setLoading(true)
      setError(false)
      try {
        const all = await skillService.getByFilter(classID, raceID)
        // Filtra apenas habilidades disponíveis para o nível atual
        const available = all.filter(s => !s.level || s.level <= level)
        setSkills(available)
      } catch {
        setError(true)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [classID, raceID, level])

  const is4e = edition === '4e'

  const grouped = {
    unlimited: skills.filter(s => s.power_type === 'unlimited'),
    encounter: skills.filter(s => s.power_type === 'encounter'),
    daily: skills.filter(s => s.power_type === 'daily'),
    utility: skills.filter(s => s.power_type === 'utility'),
    other: skills.filter(s => !s.power_type),
  }

  return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <h2 className="text-lg font-semibold text-white mb-4">
        Habilidades da Classe
      </h2>

      {loading && (
        <p className="text-gray-400 text-sm text-center py-4 animate-pulse">
          Carregando habilidades...
        </p>
      )}

      {error && (
        <p className="text-red-400 text-sm text-center py-4">
          Erro ao carregar habilidades.
        </p>
      )}

      {!loading && !error && skills.length === 0 && (
        <p className="text-gray-400 text-sm text-center py-4">
          Nenhuma habilidade cadastrada para essa classe ainda.
        </p>
      )}

      {!loading && !error && skills.length > 0 && (
        <>
          {is4e ? (
            // 4e — agrupado por tipo de poder
            <div className="flex flex-col gap-6">
              {(Object.keys(grouped) as Array<keyof typeof grouped>).map(type => {
                if (type === 'other' || grouped[type].length === 0) return null
                const cfg = powerConfig[type as keyof typeof powerConfig]
                return (
                  <div key={type}>
                    <h3 className={`text-sm font-bold mb-3 ${cfg.color}`}>
                      {cfg.label} ({grouped[type].length})
                    </h3>
                    <div className="flex flex-col gap-2">
                      {grouped[type].map(skill => (
                        <SkillCard key={skill.ID} skill={skill} />
                      ))}
                    </div>
                  </div>
                )
              })}
              {grouped.other.length > 0 && (
                <div>
                  <h3 className="text-sm font-bold mb-3 text-gray-400">
                    Outras ({grouped.other.length})
                  </h3>
                  <div className="flex flex-col gap-2">
                    {grouped.other.map(skill => (
                      <SkillCard key={skill.ID} skill={skill} />
                    ))}
                  </div>
                </div>
              )}
            </div>
          ) : (
            // 5e — lista simples sem agrupamento
            <div className="flex flex-col gap-2">
              {skills.map(skill => (
                <SkillCard key={skill.ID} skill={skill} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}