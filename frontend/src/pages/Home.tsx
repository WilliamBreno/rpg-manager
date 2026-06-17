import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

export default function Home() {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-rpg-dark flex flex-col items-center justify-center p-4 relative overflow-hidden">

      {/* ── Anéis decorativos (inspirados na logo) ── */}
      {[500, 700, 900].map((size, i) => (
        <div key={i} style={{
          position: 'absolute',
          width: size,
          height: size,
          borderRadius: '50%',
          border: `1px solid rgba(201,168,76,${0.06 - i * 0.015})`,
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          pointerEvents: 'none',
        }} />
      ))}

      {/* ── Conteúdo principal ── */}
      <div className="relative z-10 flex flex-col items-center gap-8 text-center">

        {/* Logo circular com estrelas */}
        <div style={{ position: 'relative', display: 'inline-block' }}>
          <div style={{
            width: 110, height: 110,
            borderRadius: '50%',
            background: '#0e0e0e',
            border: '1.5px solid rgba(201,168,76,0.5)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: 50,
            boxShadow: '0 0 50px rgba(201,168,76,0.12)',
          }}>
            🎲
          </div>
          {/* Estrelas de 4 pontas nos cantos */}
          {[
            { top: '-4px', left: '50%', transform: 'translateX(-50%)' },
            { bottom: '-4px', left: '50%', transform: 'translateX(-50%)' },
            { left: '-4px', top: '50%', transform: 'translateY(-50%) scale(.65)' },
            { right: '-4px', top: '50%', transform: 'translateY(-50%) scale(.65)' },
          ].map((style, i) => (
            <div key={i} style={{
              position: 'absolute',
              width: 10, height: 10,
              background: '#c9a84c',
              clipPath: 'polygon(50% 0%,57% 43%,100% 50%,57% 57%,50% 100%,43% 57%,0% 50%,43% 43%)',
              ...style,
            }} />
          ))}
        </div>

        {/* Título */}
        <div>
          <h1 className="font-rpg font-bold text-4xl sm:text-5xl" style={{ color: '#c9a84c', letterSpacing: '0.06em' }}>
            RPG Manager
          </h1>
          <div className="flex items-center justify-center gap-3 mt-3">
            <div style={{ height: 1, width: 48, background: 'rgba(201,168,76,0.3)' }} />
            <span style={{ color: 'rgba(201,168,76,0.5)', fontSize: 10 }}>✦</span>
            <div style={{ height: 1, width: 48, background: 'rgba(201,168,76,0.3)' }} />
          </div>
        </div>

        {/* Saudação */}
        <div className="flex flex-col items-center gap-2">
          <p className="text-gray-300">
            Bem-vindo,{' '}
            <span className="font-rpg font-semibold" style={{ color: '#e8c46a' }}>
              {user?.name}
            </span>
          </p>
          <span style={{
            background: 'rgba(201,168,76,0.1)',
            border: '1px solid rgba(201,168,76,0.3)',
            color: '#c9a84c',
            fontSize: 12,
            padding: '3px 14px',
            borderRadius: 20,
            fontFamily: 'Georgia, serif',
            letterSpacing: '0.04em',
          }}>
            {user?.role === 'master' ? '👑 Mestre' : '🧙 Jogador'}
          </span>
        </div>

        {/* Botões de ação */}
        <div className="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
          <button
            onClick={() => navigate('/characters')}
            className="btn-rpg-outline px-6 py-2"
          >
            📜 Ver Personagens
          </button>
          <button
            onClick={() => navigate('/characters/new')}
            className="btn-rpg-primary px-6 py-2"
          >
            ✦ Criar Personagem
          </button>
        </div>

        {/* Logout */}
        <button
          onClick={handleLogout}
          className="text-gray-600 hover:text-red-400 transition text-sm mt-2"
        >
          Sair da conta
        </button>

      </div>
    </div>
  )
}