import { useState, useRef, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { characterService } from '../services/characterService'
import { xpProgressFor } from '../lib/xpTables'
import type { Character } from '../types'

const API_BASE = (import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api').replace(/\/api$/, '')

// Avatares novos vêm como data: URI (base64, guardado direto no banco — ver
// upload_handler.go). Avatares antigos ainda podem ter um caminho relativo
// tipo "/uploads/arquivo.png" (arquivo que não existe mais em disco, cai no
// fallback via onError abaixo); só prefixa com a API_BASE nesse caso legado.
function resolveAvatarSrc(avatarURL: string): string {
  return avatarURL.startsWith('data:') ? avatarURL : `${API_BASE}${avatarURL}`
}

// Cor de destaque por card — determinística por ID (não sorteia de novo a cada render).
const ACCENTS = [
  { bar: '#c9a84c', ring: 'rgba(201,168,76,0.5)' },   // dourado (tema base)
  { bar: '#34d399', ring: 'rgba(52,211,153,0.5)' },   // esmeralda
  { bar: '#38bdf8', ring: 'rgba(56,189,248,0.5)' },   // céu
  { bar: '#f87171', ring: 'rgba(248,113,113,0.5)' },  // rubro
  { bar: '#a78bfa', ring: 'rgba(167,139,250,0.5)' },  // violeta
]

interface Props { character: Character }

export default function CharacterCard({ character }: Props) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [menuOpen, setMenuOpen] = useState(false)
  const [avatarFailed, setAvatarFailed] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const deleteMutation = useMutation({
    mutationFn: () => characterService.delete(character.ID),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['characters'] }),
  })

  useEffect(() => {
    if (!menuOpen) return
    const onClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) setMenuOpen(false)
    }
    document.addEventListener('mousedown', onClickOutside)
    return () => document.removeEventListener('mousedown', onClickOutside)
  }, [menuOpen])

  const accent = ACCENTS[character.ID % ACCENTS.length]
  const { currentXP, currentLevelXP, nextLevelXP, progressPercent, isMaxLevel } =
    xpProgressFor(character.edition, character.level, character.experience_points ?? 0)

  const handleDelete = () => {
    setMenuOpen(false)
    if (window.confirm(`Excluir ${character.name}? Essa ação não pode ser desfeita.`)) {
      deleteMutation.mutate()
    }
  }

  return (
    <div
      onClick={() => navigate(`/characters/${character.ID}`)}
      className="group relative bg-gray-800 rounded-xl cursor-pointer overflow-hidden border border-gray-700 transition-all duration-300 hover:-translate-y-1"
      style={{ ['--accent-ring' as string]: accent.ring }}
    >
      {/* Moldura temática: cantoneiras douradas + glow sutil no hover */}
      <div
        className="pointer-events-none absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"
        style={{ boxShadow: `0 0 0 1px ${accent.ring}, 0 0 24px ${accent.ring}` }}
      />
      {[
        { position: { top: 6, left: 6 }, sides: ['top', 'left'] },
        { position: { top: 6, right: 6 }, sides: ['top', 'right'] },
        { position: { bottom: 6, left: 6 }, sides: ['bottom', 'left'] },
        { position: { bottom: 6, right: 6 }, sides: ['bottom', 'right'] },
      ].map((corner, i) => (
        <div
          key={i}
          className="pointer-events-none absolute w-4 h-4 opacity-60 group-hover:opacity-100 transition-opacity duration-300"
          style={{
            ...corner.position,
            borderColor: accent.ring,
            borderTopWidth: corner.sides.includes('top') ? 2 : 0,
            borderLeftWidth: corner.sides.includes('left') ? 2 : 0,
            borderRightWidth: corner.sides.includes('right') ? 2 : 0,
            borderBottomWidth: corner.sides.includes('bottom') ? 2 : 0,
            borderStyle: 'solid',
          }}
        />
      ))}

      {/* Menu de três pontos */}
      <div ref={menuRef} className="absolute top-2 right-2 z-10">
        <button
          onClick={(e) => { e.stopPropagation(); setMenuOpen(v => !v) }}
          className="w-7 h-7 rounded-full bg-black/50 hover:bg-black/70 text-white flex items-center justify-center transition"
          aria-label="Ações do personagem"
        >
          ⋮
        </button>
        {menuOpen && (
          <div
            onClick={(e) => e.stopPropagation()}
            className="absolute right-0 mt-1 w-36 bg-gray-900 border border-gray-700 rounded-lg shadow-lg overflow-hidden text-sm"
          >
            <button
              onClick={() => navigate(`/characters/${character.ID}/edit`)}
              className="w-full text-left px-3 py-2 text-gray-200 hover:bg-gray-800 transition"
            >
              ✎ Editar
            </button>
            <button
              onClick={handleDelete}
              className="w-full text-left px-3 py-2 text-red-400 hover:bg-gray-800 transition"
            >
              🗑 Excluir
            </button>
          </div>
        )}
      </div>

      {/* Retrato */}
      <div className="h-36 bg-gray-700 flex items-center justify-center overflow-hidden relative">
        <span
          className="absolute top-2 left-2 z-10 text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full"
          style={{ background: 'rgba(0,0,0,0.6)', color: accent.bar, border: `1px solid ${accent.ring}` }}
        >
          D&D {character.edition}
        </span>
        {character.avatar_url && !avatarFailed ? (
          <img
            src={resolveAvatarSrc(character.avatar_url)}
            alt={character.name}
            className="w-full h-full object-cover"
            onError={() => setAvatarFailed(true)}
          />
        ) : (
          <span className="text-5xl">🧙</span>
        )}
      </div>

      <div className="p-4">
        <h2 className="text-base font-bold text-white truncate">{character.name}</h2>
        <p className="text-gray-400 text-xs mt-0.5">
          {character.race?.name ?? 'Sem raça'} • {character.class?.name ?? 'Sem classe'}
        </p>
        <p className="text-gray-500 text-xs mt-1">Nível {character.level}</p>

        <div className="mt-2">
          <div className="flex justify-between text-[11px] text-gray-500 mb-1">
            <span>{currentXP.toLocaleString('pt-BR')} / {isMaxLevel ? currentLevelXP.toLocaleString('pt-BR') : nextLevelXP.toLocaleString('pt-BR')} XP</span>
          </div>
          <div className="h-1.5 bg-gray-700 rounded-full overflow-hidden">
            <div
              className="h-full rounded-full transition-all duration-500"
              style={{ width: `${progressPercent}%`, background: accent.bar }}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
