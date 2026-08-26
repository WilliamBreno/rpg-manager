import { useNavigate } from 'react-router-dom'

interface WelcomeModalProps {
  onClose: () => void
}

export default function WelcomeModal({ onClose }: WelcomeModalProps) {
  const navigate = useNavigate()

  const handleCreateCharacter = () => {
    onClose()
    navigate('/characters/new')
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 animate-rpg-fade-in"
      style={{ background: 'rgba(0,0,0,0.88)', backdropFilter: 'blur(3px)' }}
      role="dialog"
      aria-modal="true"
      aria-labelledby="welcome-modal-title"
    >
      <button
        onClick={onClose}
        className="absolute top-5 right-5 text-gray-500 hover:text-gray-300 transition text-sm"
        aria-label="Fechar"
      >
        ✕ pular
      </button>

      <div className="w-full max-w-md text-center">
        <div
          className="inline-flex items-center justify-center w-24 h-24 rounded-full mb-8 text-6xl animate-rpg-dice-entrance"
          style={{
            background: 'radial-gradient(circle, rgba(201,168,76,0.15) 0%, rgba(0,0,0,0) 70%)',
            border: '2px solid rgba(201,168,76,0.5)',
            boxShadow: '0 0 60px rgba(201,168,76,0.25)',
          }}
        >
          🎲
        </div>

        <h1
          id="welcome-modal-title"
          className="font-rpg text-3xl sm:text-4xl font-bold mb-4 animate-rpg-rise"
          style={{ color: '#e8c46a', letterSpacing: '0.02em', animationDelay: '0.3s' }}
        >
          O dado foi lançado.
        </h1>
        <p
          className="font-rpg text-lg sm:text-xl text-gray-300 mb-10 leading-relaxed animate-rpg-rise"
          style={{ animationDelay: '0.45s' }}
        >
          Sua jornada de uma lenda começa aqui.
        </p>

        <button
          onClick={handleCreateCharacter}
          className="btn-rpg-primary px-8 py-3 text-base animate-rpg-rise"
          style={{ animationDelay: '0.6s' }}
        >
          ✦ Criar meu primeiro personagem
        </button>
      </div>
    </div>
  )
}
