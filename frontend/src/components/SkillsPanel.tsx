import { useEffect, useState } from 'react'
import { skillService } from '../services/skillService'
import { SkillCard, powerConfig } from '../components/SkillCard'
import type { Skill, PowerType } from '../types'

interface Props {
  characterID: number
  classID: number
  raceID: number
  className: string
  edition: string
  level: number
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

  // Separa características de classe das habilidades normais
  const classFeatures      = skills.filter(s => s.is_class_feature && !s.requires_choice)
  const choiceFeatures     = skills.filter(s => s.is_class_feature && s.requires_choice)
  const normalSkills       = skills.filter(s => !s.is_class_feature)

  // Agrupa escolhas por grupo
  const choiceGroups = choiceFeatures.reduce<Record<string, Skill[]>>((acc, s) => {
    const g = s.choice_group ?? 'default'
    acc[g] = [...(acc[g] ?? []), s]
    return acc
  }, {})

  // Agrupa habilidades normais por tipo
  const byType = (type: PowerType) => normalSkills.filter(s => s.power_type === type)

  if (loading) return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <p className="text-gray-400 text-sm animate-pulse">Carregando habilidades...</p>
    </div>
  )

  if (error) return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <p className="text-red-400 text-sm">Erro ao carregar habilidades.</p>
    </div>
  )

  if (skills.length === 0) return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <p className="text-gray-400 text-sm text-center py-4">Nenhuma habilidade cadastrada para essa classe ainda.</p>
    </div>
  )

  return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <h2 className="text-lg font-semibold text-white mb-5">Habilidades da Classe</h2>

      <div className="flex flex-col gap-6">

        {/* Características automáticas */}
        {classFeatures.length > 0 && (
          <div>
            <h3 className="text-sm font-bold text-indigo-400 mb-2">📖 Características de Classe</h3>
            <p className="text-gray-500 text-xs mb-3">Concedidas automaticamente — não requerem seleção.</p>
            <div className="flex flex-col gap-2">
              {classFeatures.map(s => <SkillCard key={s.ID} skill={s} informative />)}
            </div>
          </div>
        )}

        {/* Características com escolha */}
        {Object.entries(choiceGroups).map(([group, options]) => (
          <div key={group}>
            <h3 className="text-sm font-bold text-purple-400 mb-1">🎯 Escolha de Característica</h3>
            <p className="text-gray-500 text-xs mb-3">Você escolheu uma das opções abaixo durante a criação.</p>
            <div className="flex flex-col gap-2">
              {options.map(s => <SkillCard key={s.ID} skill={s} informative />)}
            </div>
          </div>
        ))}

        {/* Habilidades normais agrupadas por tipo (4e) */}
        {is4e ? (
          <>
            {(['unlimited', 'encounter', 'daily', 'utility'] as PowerType[]).map(type => {
              const list = byType(type)
              if (list.length === 0) return null
              const cfg = powerConfig[type]
              return (
                <div key={type}>
                  <h3 className={`text-sm font-bold mb-2 ${cfg.color}`}>
                    {cfg.label} ({list.length})
                  </h3>
                  <div className="flex flex-col gap-2">
                    {list.map(s => <SkillCard key={s.ID} skill={s} />)}
                  </div>
                </div>
              )
            })}
          </>
        ) : (
          <div className="flex flex-col gap-2">
            {normalSkills.map(s => <SkillCard key={s.ID} skill={s} />)}
          </div>
        )}

      </div>
    </div>
  )
}