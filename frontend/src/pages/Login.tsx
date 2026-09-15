import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { authService } from '../services/authService'
import { useAuthStore } from '../store/authStore'

export default function Login() {
  const navigate = useNavigate()
  const { setAuth } = useAuthStore()
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const response = await authService.login(form)
      setAuth(response.token, response.user)
      navigate('/characters')
    } catch (err: any) {
      setError(err.response?.data?.error ?? '⚠ Não foi possível entrar')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-rpg-dark flex items-center justify-center p-4 relative overflow-hidden">

      {/* Vinheta radial suave atrás de tudo — dá profundidade sem competir com o selo */}
      <div
        className="pointer-events-none absolute inset-0"
        style={{ background: 'radial-gradient(ellipse at 50% 32%, rgba(201,168,76,0.07) 0%, rgba(0,0,0,0) 55%)' }}
      />

      <div className="w-full max-w-sm relative">

        {/* ── Selo ── */}
        <div className="text-center mb-9">
          <div
            className="relative inline-flex items-center justify-center w-28 h-28 mb-5 animate-rpg-dice-entrance"
          >
            <div
              className="absolute inset-0 rounded-full"
              style={{ background: 'radial-gradient(circle, rgba(201,168,76,0.16) 0%, rgba(0,0,0,0) 68%)' }}
            />
            <img src="/logo.png" alt="RPG Manager" className="relative w-24 h-24 object-contain drop-shadow-[0_0_18px_rgba(201,168,76,0.25)]" />
          </div>
          <h1
            className="font-rpg text-3xl font-bold"
            style={{ color: '#c9a84c', letterSpacing: '0.08em' }}
          >
            RPG MANAGER
          </h1>
          <p className="font-rpg text-sm italic mt-2" style={{ color: 'rgba(201,168,76,0.55)' }}>
            Sua ficha de aventureiro te aguarda.
          </p>
        </div>

        {/* ── Card com cantoneiras douradas ── */}
        <div className="relative rounded-xl p-7" style={{ background: '#161616', border: '1px solid rgba(201,168,76,0.2)', boxShadow: '0 0 50px rgba(201,168,76,0.05)' }}>
          {[
            { style: { top: -1, left: -1 }, sides: ['top', 'left'] },
            { style: { top: -1, right: -1 }, sides: ['top', 'right'] },
            { style: { bottom: -1, left: -1 }, sides: ['bottom', 'left'] },
            { style: { bottom: -1, right: -1 }, sides: ['bottom', 'right'] },
          ].map((corner, i) => (
            <div
              key={i}
              className="pointer-events-none absolute w-5 h-5"
              style={{
                ...corner.style,
                borderColor: 'rgba(201,168,76,0.55)',
                borderTopWidth: corner.sides.includes('top') ? 2 : 0,
                borderLeftWidth: corner.sides.includes('left') ? 2 : 0,
                borderRightWidth: corner.sides.includes('right') ? 2 : 0,
                borderBottomWidth: corner.sides.includes('bottom') ? 2 : 0,
                borderStyle: 'solid',
              }}
            />
          ))}

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">

            <div>
              <label className="text-gray-500 text-xs font-medium mb-1.5 block uppercase tracking-widest">
                Email
              </label>
              <input
                type="email"
                value={form.email}
                onChange={e => setForm(prev => ({ ...prev, email: e.target.value }))}
                className="rpg-input"
                placeholder="seu@email.com"
              />
            </div>

            <div>
              <label className="text-gray-500 text-xs font-medium mb-1.5 block uppercase tracking-widest">
                Senha
              </label>
              <input
                type="password"
                value={form.password}
                onChange={e => setForm(prev => ({ ...prev, password: e.target.value }))}
                className="rpg-input"
                placeholder="••••••"
              />
            </div>

            {error && (
              <p
                className="text-red-400 text-xs text-center py-2 px-3 rounded-lg"
                style={{ background: 'rgba(220,38,38,0.08)', border: '1px solid rgba(220,38,38,0.25)' }}
              >
                {error}
              </p>
            )}

            <button type="submit" disabled={loading} className="btn-rpg-primary w-full mt-1">
              {loading ? '⚔️ Entrando...' : '✦ Entrar'}
            </button>

          </form>

          <div className="pt-4 mt-4 text-center" style={{ borderTop: '1px solid rgba(201,168,76,0.1)' }}>
            <p className="text-gray-600 text-sm">
              Não tem conta?{' '}
              <Link to="/register" className="transition" style={{ color: '#c9a84c' }}>
                Cadastre-se
              </Link>
            </p>
          </div>
        </div>

      </div>
    </div>
  )
}
