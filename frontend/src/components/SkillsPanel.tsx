import type { Skill, PowerType } from '../types'
import { SkillCard, powerConfig } from './SkillCard'

interface Props {
  skills: Skill[]  // character.skills — já filtrados/escolhidos
  edition: string
}

export default function SkillsPanel({ skills, edition }: Props) {
  if (!skills || skills.length === 0) {
    return (
      <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
        <h2 className="text-lg font-semibold text-white mb-2">Habilidades da Classe</h2>
        <p className="text-gray-400 text-sm text-center py-4">
          Nenhuma habilidade atribuída ainda.
        </p>
      </div>
    )
  }

  const is4e = edition === '4e'

  // Separa características de classe das habilidades normais
  const classFeatures  = skills.filter(s => s.is_class_feature && !s.requires_choice)
  const chosenFeatures = skills.filter(s => s.is_class_feature && s.requires_choice)
  const normalSkills   = skills.filter(s => !s.is_class_feature)

  const byType = (type: PowerType) => normalSkills.filter(s => s.power_type === type)

  return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <h2 className="text-lg font-semibold text-white mb-5">Habilidades da Classe</h2>

      <div className="flex flex-col gap-6">

        {/* Características automáticas */}
        {classFeatures.length > 0 && (
          <div>
            <h3 className="text-sm font-bold text-indigo-400 mb-1">📖 Características de Classe</h3>
            <p className="text-gray-500 text-xs mb-3">Concedidas automaticamente pela sua classe.</p>
            <div className="flex flex-col gap-2">
              {classFeatures.map(s => <SkillCard key={s.ID} skill={s} informative />)}
            </div>
          </div>
        )}

        {/* Características escolhidas */}
        {chosenFeatures.length > 0 && (
          <div>
            <h3 className="text-sm font-bold text-purple-400 mb-1">🎯 Características Escolhidas</h3>
            <div className="flex flex-col gap-2">
              {chosenFeatures.map(s => <SkillCard key={s.ID} skill={s} informative />)}
            </div>
          </div>
        )}

        {/* Poderes normais (4e agrupados por tipo, 5e lista simples) */}
        {normalSkills.length > 0 && (
          is4e ? (
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
          )
        )}

      </div>
    </div>
  )
}