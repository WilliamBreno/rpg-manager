import { useState } from 'react'
import { aiService } from '../services/aiService'
import type { AISkill } from '../services/aiService'

interface Props {
  characterID: number
  className: string
  edition: string
  level: number
}

const powerConfig = {
  'at-will': {
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

function SkillCard({ skill }: { skill: AISkill }) {
  const [expanded, setExpanded] = useState(false)
  const config = powerConfig[skill.power_type] ?? powerConfig['at-will']

  return (
    <div
      className={`rounded-lg border ${config.border} ${config.bg} p-4 cursor-pointer transition-all duration-200`}
      onClick={() => setExpanded(!expanded)}
    >
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-2">
          <h3 className="text-white font-semibold text-sm">{skill.name}</h3>
          {skill.min_level > 1 && (
            <span className="text-gray-500 text-xs">Nível {skill.min_level}</span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <span className={`text-xs px-2 py-1 rounded-full ${config.badge}`}>
            {config.label}
          </span>
          <span className="text-gray-400 text-xs">{expanded ? '▲' : '▼'}</span>
        </div>
      </div>

      {expanded && (
        <div className="mt-3 flex flex-col gap-2 border-t border-gray-700 pt-3">
          {skill.action_type && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Ação:</span>
              <span className="text-white">{skill.action_type}</span>
            </div>
          )}
          {skill.range && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Alcance:</span>
              <span className="text-white">{skill.range}</span>
            </div>
          )}
          {skill.target && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Alvo:</span>
              <span className="text-white">{skill.target}</span>
            </div>
          )}
          {skill.attack && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Ataque:</span>
              <span className="text-indigo-300 font-semibold">{skill.attack}</span>
            </div>
          )}
          {skill.hit && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Sucesso:</span>
              <span className="text-green-300">{skill.hit}</span>
            </div>
          )}
          {skill.miss && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Fracasso:</span>
              <span className="text-red-300">{skill.miss}</span>
            </div>
          )}
          {skill.effect && (
            <div className="flex gap-2 text-xs">
              <span className="text-gray-400 w-24 shrink-0">Efeito:</span>
              <span className="text-yellow-200">{skill.effect}</span>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default function SkillsPanel({ className, edition, level }: Props) {
  const [skills, setSkills] = useState<AISkill[] | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(false)

  const handleLoad = async () => {
    setLoading(true)
    setError(false)
    try {
      const result = await aiService.getSkills(className, edition, level)
      setSkills(result)
    } catch {
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  const grouped = {
    'at-will': skills?.filter(s => s.power_type === 'at-will') ?? [],
    encounter: skills?.filter(s => s.power_type === 'encounter') ?? [],
    daily: skills?.filter(s => s.power_type === 'daily') ?? [],
    utility: skills?.filter(s => s.power_type === 'utility') ?? [],
  }

  return (
    <div className="bg-gray-800 rounded-lg p-6 mb-4 border border-gray-700">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold text-white">Habilidades da Classe</h2>
        <button
          onClick={handleLoad}
          disabled={loading}
          className="bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white text-sm font-semibold px-4 py-2 rounded-lg transition flex items-center gap-2"
        >
          {loading ? (
            <>
              <span className="animate-spin inline-block">⟳</span>
              Consultando IA...
            </>
          ) : (
            <>🎲 {skills ? 'Recarregar' : 'Carregar'} Habilidades</>
          )}
        </button>
      </div>

      {!skills && !loading && !error && (
        <p className="text-gray-400 text-sm text-center py-4">
          Clique em "Carregar Habilidades" para consultar os livros de D&D.
        </p>
      )}

      {error && (
        <p className="text-red-400 text-sm text-center py-4">
          Erro ao consultar a IA. Verifique se o Ollama está rodando.
        </p>
      )}

      {loading && (
        <div className="text-center py-8">
          <p className="text-indigo-400 text-sm animate-pulse">
            🧠 Consultando os livros de D&D...
          </p>
          <p className="text-gray-500 text-xs mt-2">
            Isso pode levar alguns segundos
          </p>
        </div>
      )}

      {skills && skills.length > 0 && (
        <div className="flex flex-col gap-6">
          {(Object.keys(grouped) as Array<keyof typeof grouped>).map(type => (
            grouped[type].length > 0 && (
              <div key={type}>
                <h3 className={`text-sm font-bold mb-3 ${powerConfig[type].color}`}>
                  {powerConfig[type].label} ({grouped[type].length})
                </h3>
                <div className="flex flex-col gap-2">
                  {grouped[type].map((skill, i) => (
                    <SkillCard key={i} skill={skill} />
                  ))}
                </div>
              </div>
            )
          ))}
        </div>
      )}

      {skills && skills.length === 0 && (
        <p className="text-gray-400 text-sm text-center py-4">
          Nenhuma habilidade encontrada.
        </p>
      )}
    </div>
  )
}