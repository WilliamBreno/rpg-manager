import { useEffect, useState } from 'react'

interface LevelUpCelebrationProps {
  level: number
  onDone: () => void
}

const PARTICLE_COUNT = 20

function generateParticles() {
  return Array.from({ length: PARTICLE_COUNT }, (_, i) => {
    const angle = (i / PARTICLE_COUNT) * Math.PI * 2 + Math.random() * 0.3
    const dist = 90 + Math.random() * 90
    return {
      dx: Math.cos(angle) * dist,
      dy: Math.sin(angle) * dist,
      delay: Math.random() * 0.15,
      size: 4 + Math.random() * 5,
    }
  })
}

export default function LevelUpCelebration({ level, onDone }: LevelUpCelebrationProps) {
  const [closing, setClosing] = useState(false)
  // Lazy initializer: roda uma única vez na montagem, não a cada render —
  // por isso pode chamar Math.random() sem violar a regra de pureza de render.
  const [particles] = useState(generateParticles)

  useEffect(() => {
    const closeTimer = setTimeout(() => setClosing(true), 2600)
    const doneTimer = setTimeout(onDone, 3000)
    return () => { clearTimeout(closeTimer); clearTimeout(doneTimer) }
  }, [onDone])

  return (
    <div
      className={`fixed inset-0 z-50 flex items-center justify-center pointer-events-none transition-opacity duration-300 ${closing ? 'opacity-0' : 'opacity-100 animate-rpg-fade-in'}`}
      style={{ background: 'rgba(0,0,0,0.6)' }}
      role="status"
      aria-live="polite"
    >
      <div className="relative flex flex-col items-center">
        {/* Anéis de expansão */}
        <div className="absolute w-40 h-40 rounded-full border-2 border-rpg-gold animate-rpg-ring-expand" />
        <div className="absolute w-40 h-40 rounded-full border-2 border-rpg-gold-light animate-rpg-ring-expand" style={{ animationDelay: '0.2s' }} />

        {/* Partículas douradas */}
        {particles.map((p, i) => (
          <div
            key={i}
            className="absolute rounded-full bg-rpg-gold-light animate-rpg-burst-particle"
            style={{
              width: p.size,
              height: p.size,
              '--dx': `${p.dx}px`,
              '--dy': `${p.dy}px`,
              animationDelay: `${p.delay}s`,
              boxShadow: '0 0 6px rgba(232,196,106,0.8)',
            } as React.CSSProperties}
          />
        ))}

        <div className="relative flex flex-col items-center animate-rpg-level-pop">
          <p className="font-rpg text-sm sm:text-base tracking-[0.3em] uppercase mb-1" style={{ color: '#e8c46a' }}>
            Nível alcançado
          </p>
          <p className="font-rpg font-bold leading-none" style={{ color: '#c9a84c', fontSize: '5rem', textShadow: '0 0 40px rgba(201,168,76,0.5)' }}>
            {level}
          </p>
          <p className="text-gray-300 text-sm sm:text-base mt-2">Você ficou mais forte.</p>
        </div>
      </div>
    </div>
  )
}
