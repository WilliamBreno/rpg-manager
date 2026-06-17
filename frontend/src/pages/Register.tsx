import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { authService } from '../services/authService'

export default function Register() {
  const navigate = useNavigate()
  const [form, setForm] = useState({
    name: '',
    email: '',
    password: '',
    role: 'player' as 'player' | 'master',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await authService.register(form)
      navigate('/login')
    } catch (err: any) {
      setError(err.response?.data?.error ?? 'Erro ao cadastrar')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-rpg-dark flex items-center justify-center p-4">
      <div className="w-full max-w-sm">

        {/* ── Logo ── */}
        <div className="text-center mb-8">
          <div
            className="inline-flex items-center justify-center w-16 h-16 rounded-full mb-4 text-4xl"
            style={{ background: '#0e0e0e', border: '1.5px solid rgba(201,168,76,0.45)', boxShadow: '0 0 30px rgba(201,168,76,0.08)' }}
          >
            🎲
          </div>
          <h1 className="font-rpg text-3xl font-bold" style={{ color: '#c9a84c', letterSpacing: '0.05em' }}>
            RPG Manager
          </h1>
          <div className="flex items-center justify-center gap-3 mt-2">
            <div style={{ height: 1, width: 32, background: 'rgba(201,168,76,0.25)' }} />
            <span style={{ color: 'rgba(201,168,76,0.45)', fontSize: 9 }}>✦</span>
            <p className="text-gray-500 text-sm">Crie sua conta</p>
            <span style={{ color: 'rgba(201,168,76,0.45)', fontSize: 9 }}>✦</span>
            <div style={{ height: 1, width: 32, background: 'rgba(201,168,76,0.25)' }} />
          </div>
        </div>

        {/* ── Card ── */}
        <div
          className="rounded-xl p-7 flex flex-col gap-5"
          style={{ background: '#161616', border: '1px solid rgba(201,168,76,0.2)', boxShadow: '0 0 50px rgba(201,168,76,0.05)' }}
        >
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">

            <div>
              <label className="text-gray-500 text-xs font-medium mb-1.5 block uppercase tracking-widest">Nome</label>
              <input
                type="text"
                value={form.name}
                onChange={e => setForm(prev => ({ ...prev, name: e.target.value }))}
                className="rpg-input"
                placeholder="Seu nome"
              />
            </div>

            <div>
              <label className="text-gray-500 text-xs font-medium mb-1.5 block uppercase tracking-widest">Email</label>
              <input
                type="email"
                value={form.email}
                onChange={e => setForm(prev => ({ ...prev, email: e.target.value }))}
                className="rpg-input"
                placeholder="seu@email.com"
              />
            </div>

            <div>
              <label className="text-gray-500 text-xs font-medium mb-1.5 block uppercase tracking-widest">Senha</label>
              <input
                type="password"
                value={form.password}
                onChange={e => setForm(prev => ({ ...prev, password: e.target.value }))}
                className="rpg-input"
                placeholder="Mínimo 6 caracteres"
              />
            </div>

            <div>
              <label className="text-gray-500 text-xs font-medium mb-1.5 block uppercase tracking-widest">Tipo de conta</label>
              <select
                value={form.role}
                onChange={e => setForm(prev => ({ ...prev, role: e.target.value as 'player' | 'master' }))}
                className="rpg-select"
              >
                <option value="player">🧙 Jogador</option>
                <option value="master">👑 Mestre</option>
              </select>
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
              {loading ? '⚔️ Cadastrando...' : '✦ Criar Conta'}
            </button>

          </form>

          <div className="pt-4 text-center" style={{ borderTop: '1px solid rgba(201,168,76,0.1)' }}>
            <p className="text-gray-600 text-sm">
              Já tem conta?{' '}
              <Link to="/login" className="transition" style={{ color: '#c9a84c' }}>
                Entrar
              </Link>
            </p>
          </div>
        </div>

      </div>
    </div>
  )
}