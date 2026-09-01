import { useState } from 'react'

interface MasterWelcomeModalProps {
  onClose: () => void
  onCreateCampaign: () => void
}

const POWERS = [
  { icon: '🗺️', label: 'Campanhas', text: 'Erga mundos com história, NPCs, inimigos e boss sob medida.' },
  { icon: '🕯️', label: 'Sala ao vivo', text: 'Cenários, chat e dados em tempo real com seus jogadores.' },
  { icon: '💰', label: 'Recompensas', text: 'Distribua moeda e itens mágicos com um gesto.' },
]

const EMBER_COUNT = 16

function generateEmbers() {
  return Array.from({ length: EMBER_COUNT }, (_, i) => {
    const angle = (i / EMBER_COUNT) * Math.PI * 2 + Math.random() * 0.4
    const dist = 100 + Math.random() * 120
    return {
      dx: Math.cos(angle) * dist,
      dy: Math.sin(angle) * dist,
      delay: Math.random() * 0.25,
      size: 3 + Math.random() * 4,
    }
  })
}

// Tela de boas-vindas exclusiva de Mestre — mostrada uma única vez na
// primeira visita a "Minhas Campanhas" (User.MasterWelcomeSeen), separada da
// boas-vindas geral de jogador (WelcomeModal). Reaproveita o vocabulário
// visual já estabelecido (anéis dourados, partículas, animações rpg-*) num
// tema de "coroação" — o Mestre senta ao trono da própria mesa.
export default function MasterWelcomeModal({ onClose, onCreateCampaign }: MasterWelcomeModalProps) {
  const [embers] = useState(generateEmbers)

  const handleCreate = () => {
    onClose()
    onCreateCampaign()
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 overflow-y-auto animate-rpg-fade-in"
      style={{ background: 'rgba(5,3,0,0.92)', backdropFilter: 'blur(4px)' }}
      role="dialog"
      aria-modal="true"
      aria-labelledby="master-welcome-title"
    >
      <button
        onClick={onClose}
        className="fixed top-5 right-5 text-gray-500 hover:text-gray-300 transition text-sm z-10"
        aria-label="Fechar"
      >
        ✕ pular
      </button>

      <div className="w-full max-w-lg text-center py-8">

        {/* ── Selo do Mestre: anéis + brasas + coroa ── */}
        <div className="relative flex items-center justify-center mb-8 mx-auto" style={{ width: 140, height: 140 }}>
          <div className="absolute inset-0 rounded-full border-2 border-rpg-gold animate-rpg-ring-expand" />
          <div className="absolute inset-0 rounded-full border-2 border-rpg-gold-light animate-rpg-ring-expand" style={{ animationDelay: '0.25s' }} />
          <div className="absolute inset-0 rounded-full border border-rpg-gold animate-rpg-ring-expand" style={{ animationDelay: '0.5s' }} />

          {embers.map((p, i) => (
            <div
              key={i}
              className="absolute rounded-full animate-rpg-burst-particle"
              style={{
                width: p.size,
                height: p.size,
                background: '#e8c46a',
                '--dx': `${p.dx}px`,
                '--dy': `${p.dy}px`,
                animationDelay: `${p.delay}s`,
                boxShadow: '0 0 8px rgba(232,196,106,0.9)',
              } as React.CSSProperties}
            />
          ))}

          <div
            className="relative flex items-center justify-center w-24 h-24 rounded-full text-6xl animate-rpg-dice-entrance"
            style={{
              background: 'radial-gradient(circle, rgba(201,168,76,0.2) 0%, rgba(0,0,0,0) 72%)',
              border: '2px solid rgba(201,168,76,0.6)',
              boxShadow: '0 0 70px rgba(201,168,76,0.3)',
            }}
          >
            👑
          </div>
        </div>

        <p
          className="text-xs uppercase tracking-[0.35em] mb-3 animate-rpg-rise"
          style={{ color: 'rgba(201,168,76,0.7)', animationDelay: '0.15s' }}
        >
          Sistema do Mestre
        </p>

        <h1
          id="master-welcome-title"
          className="font-rpg text-3xl sm:text-4xl font-bold mb-4 animate-rpg-rise"
          style={{ color: '#e8c46a', letterSpacing: '0.02em', animationDelay: '0.3s' }}
        >
          Sente-se ao trono, Mestre.
        </h1>
        <p
          className="font-rpg text-base sm:text-lg text-gray-300 mb-8 leading-relaxed max-w-md mx-auto animate-rpg-rise"
          style={{ animationDelay: '0.42s' }}
        >
          Daqui você tece campanhas, dá vida a vilões e comanda a mesa inteira —
          seus jogadores só verão a história que você decidir contar.
        </p>

        {/* ── Pergaminho de poderes ── */}
        <div
          className="rounded-xl p-5 mb-9 text-left animate-rpg-rise"
          style={{
            background: 'linear-gradient(180deg, rgba(201,168,76,0.06), rgba(0,0,0,0) 60%)',
            border: '1px solid rgba(201,168,76,0.22)',
            animationDelay: '0.54s',
          }}
        >
          <div className="flex flex-col gap-4">
            {POWERS.map(p => (
              <div key={p.label} className="flex items-start gap-3">
                <span className="text-xl flex-shrink-0" style={{ filter: 'drop-shadow(0 0 6px rgba(201,168,76,0.4))' }}>{p.icon}</span>
                <div>
                  <p className="text-sm font-semibold" style={{ color: '#e8c46a' }}>{p.label}</p>
                  <p className="text-gray-400 text-sm leading-snug">{p.text}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <button
          onClick={handleCreate}
          className="btn-rpg-primary px-8 py-3 text-base animate-rpg-rise"
          style={{ animationDelay: '0.68s' }}
        >
          ✦ Fundar minha primeira campanha
        </button>
      </div>
    </div>
  )
}
