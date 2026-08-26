import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { characterService } from '../services/characterService'
import { itemService } from '../services/itemService'
import { armorService } from '../services/armorService'
import { useAuthStore } from '../store/authStore'
import type { Item, ItemCategory, Armor, Currency } from '../types'

const COIN_LABELS: { key: keyof Currency; label: string; color: string }[] = [
  { key: 'platinum_pieces', label: 'PL', color: '#cbd5e1' },
  { key: 'gold_pieces',     label: 'PO', color: '#c9a84c' },
  { key: 'electrum_pieces', label: 'PE', color: '#e5e7eb' },
  { key: 'silver_pieces',   label: 'PP', color: '#9ca3af' },
  { key: 'copper_pieces',   label: 'PC', color: '#b45309' },
]

const CATEGORY_TABS: { key: ItemCategory | 'armadura'; label: string }[] = [
  { key: 'arma', label: '⚔️ Armas' },
  { key: 'armadura', label: '🛡️ Armaduras' },
  { key: 'equipamento', label: '🎒 Equipamento' },
  { key: 'item_magico', label: '✨ Itens Mágicos' },
]

function totalCopper(c: Currency): number {
  return c.copper_pieces + c.silver_pieces * 10 + c.electrum_pieces * 50 + c.gold_pieces * 100 + c.platinum_pieces * 1000
}

function formatCost(costCopper: number): string {
  if (costCopper % 100 === 0) return `${costCopper / 100} PO`
  if (costCopper >= 100) return `${(costCopper / 100).toFixed(2)} PO`
  if (costCopper % 10 === 0) return `${costCopper / 10} PP`
  return `${costCopper} PC`
}

const RARITY_LABELS: Record<string, string> = {
  comum: 'Comum', incomum: 'Incomum', raro: 'Raro', muito_raro: 'Muito Raro', lendario: 'Lendário',
}

export default function Shop() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const characterId = Number(id)

  const [tab, setTab] = useState<ItemCategory | 'armadura'>('arma')
  const [feedback, setFeedback] = useState<string | null>(null)
  const [grantForm, setGrantForm] = useState<Currency>({
    copper_pieces: 0, silver_pieces: 0, electrum_pieces: 0, gold_pieces: 0, platinum_pieces: 0,
  })

  const { data: character } = useQuery({
    queryKey: ['character', id],
    queryFn: () => characterService.getByID(characterId),
  })

  const { data: items } = useQuery({
    queryKey: ['items', '5e', tab],
    queryFn: () => itemService.getAll('5e', tab as ItemCategory),
    enabled: tab !== 'armadura',
  })

  const { data: armors } = useQuery({
    queryKey: ['armors', '5e'],
    queryFn: () => armorService.getByEdition('5e'),
    enabled: tab === 'armadura',
  })

  const invalidateCharacter = () => {
    queryClient.invalidateQueries({ queryKey: ['character', id] })
    queryClient.invalidateQueries({ queryKey: ['inventory', id] })
  }

  const buyItemMutation = useMutation({
    mutationFn: (item: Item) => characterService.buyItem(characterId, item.ID),
    onSuccess: (_, item) => { invalidateCharacter(); setFeedback(`Comprado: ${item.name}`); setTimeout(() => setFeedback(null), 2500) },
    onError: () => { setFeedback('Ouro insuficiente para essa compra.'); setTimeout(() => setFeedback(null), 2500) },
  })

  const buyArmorMutation = useMutation({
    mutationFn: (armor: Armor) => characterService.buyArmor(characterId, armor.ID),
    onSuccess: (_, armor) => { invalidateCharacter(); setFeedback(`Comprado: ${armor.name}`); setTimeout(() => setFeedback(null), 2500) },
    onError: () => { setFeedback('Ouro insuficiente para essa compra.'); setTimeout(() => setFeedback(null), 2500) },
  })

  const grantMutation = useMutation({
    mutationFn: () => characterService.setCurrency(characterId, grantForm),
    onSuccess: () => { invalidateCharacter(); setFeedback('Moedas concedidas.'); setTimeout(() => setFeedback(null), 2500) },
    onError: () => { setFeedback('Não foi possível conceder moedas.'); setTimeout(() => setFeedback(null), 2500) },
  })

  if (!character) return <div className="min-h-screen bg-gray-900 p-8 text-gray-400">Carregando...</div>

  const currency: Currency = {
    copper_pieces: character.copper_pieces ?? 0,
    silver_pieces: character.silver_pieces ?? 0,
    electrum_pieces: character.electrum_pieces ?? 0,
    gold_pieces: character.gold_pieces ?? 0,
    platinum_pieces: character.platinum_pieces ?? 0,
  }
  const available = totalCopper(currency)
  const isMaster = user?.role === 'master'

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-4xl mx-auto">
        <button onClick={() => navigate(`/characters/${id}`)} className="text-sm mb-4 block transition"
          style={{ color: 'rgba(201,168,76,0.5)' }}
        >← Voltar pro personagem</button>

        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mb-6">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-white">Loja</h1>
            <p className="text-gray-500 text-sm mt-1">{character.name} — compre equipamento com as moedas disponíveis.</p>
          </div>
          <div className="flex gap-2 flex-wrap">
            {COIN_LABELS.map(({ key, label, color }) => (
              <div key={key} className="rpg-card px-3 py-1.5 text-sm" style={{ minWidth: 56, textAlign: 'center' }}>
                <div className="font-bold" style={{ color }}>{currency[key]}</div>
                <div className="text-gray-500 text-[10px] uppercase">{label}</div>
              </div>
            ))}
          </div>
        </div>

        {feedback && (
          <div className="rpg-card-gold px-4 py-2 mb-4 text-sm text-center" style={{ color: '#e8c46a' }}>{feedback}</div>
        )}

        {isMaster && (
          <div className="rpg-card p-4 mb-6">
            <p className="text-rpg-gold-light text-sm font-semibold mb-3">👑 Conceder moedas (Mestre)</p>
            <div className="grid grid-cols-2 sm:grid-cols-5 gap-2 mb-3">
              {COIN_LABELS.map(({ key, label }) => (
                <div key={key}>
                  <label className="text-gray-500 text-[10px] uppercase block mb-1">{label}</label>
                  <input
                    type="number" min={0} value={grantForm[key]}
                    onChange={e => setGrantForm(prev => ({ ...prev, [key]: Math.max(0, Number(e.target.value)) }))}
                    className="rpg-input text-sm py-1.5"
                  />
                </div>
              ))}
            </div>
            <button onClick={() => grantMutation.mutate()} disabled={grantMutation.isPending} className="btn-rpg-primary text-sm px-4 py-1.5">
              {grantMutation.isPending ? 'Concedendo...' : 'Definir moedas do personagem'}
            </button>
            <p className="text-gray-600 text-xs mt-2">Isso define o total absoluto de cada moeda (não soma ao que já existe).</p>
          </div>
        )}

        <div className="flex gap-2 mb-5 flex-wrap">
          {CATEGORY_TABS.map(t => (
            <button
              key={t.key}
              onClick={() => setTab(t.key)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition ${
                tab === t.key ? 'bg-rpg-gold-muted text-rpg-gold' : 'bg-gray-800 text-gray-400 hover:text-gray-200'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          {tab !== 'armadura' && items?.map(item => {
            const canAfford = available >= item.cost_copper
            return (
              <div key={item.ID} className="rpg-card p-4 flex flex-col gap-2">
                <div className="flex justify-between items-start gap-2">
                  <div>
                    <p className="text-white font-semibold text-sm">{item.name}</p>
                    {item.rarity && (
                      <span className="rpg-badge">{RARITY_LABELS[item.rarity] ?? item.rarity}</span>
                    )}
                  </div>
                  <span className="text-rpg-gold text-sm font-semibold whitespace-nowrap">{formatCost(item.cost_copper)}</span>
                </div>
                {item.description && <p className="text-gray-500 text-xs">{item.description}</p>}
                {item.weight && item.weight !== '—' && <p className="text-gray-600 text-xs">Peso: {item.weight}</p>}
                <button
                  onClick={() => buyItemMutation.mutate(item)}
                  disabled={!canAfford || buyItemMutation.isPending}
                  className="btn-rpg-outline text-xs px-3 py-1.5 self-start disabled:opacity-40 disabled:cursor-not-allowed"
                >
                  {canAfford ? 'Comprar' : 'Ouro insuficiente'}
                </button>
              </div>
            )
          })}

          {tab === 'armadura' && armors?.filter(a => a.armor_type !== 'none').map(armor => {
            const cost = armor.cost_copper ?? 0
            const canAfford = available >= cost
            return (
              <div key={armor.ID} className="rpg-card p-4 flex flex-col gap-2">
                <div className="flex justify-between items-start gap-2">
                  <p className="text-white font-semibold text-sm">{armor.name}</p>
                  <span className="text-rpg-gold text-sm font-semibold whitespace-nowrap">{formatCost(cost)}</span>
                </div>
                <p className="text-gray-500 text-xs">
                  CA base {armor.base_ac}{armor.max_dex_bonus !== 0 && ' + mod. Destreza'}
                  {armor.max_dex_bonus > 0 && ` (máx. ${armor.max_dex_bonus})`}
                </p>
                {armor.weight && <p className="text-gray-600 text-xs">Peso: {armor.weight}</p>}
                <button
                  onClick={() => buyArmorMutation.mutate(armor)}
                  disabled={!canAfford || buyArmorMutation.isPending}
                  className="btn-rpg-outline text-xs px-3 py-1.5 self-start disabled:opacity-40 disabled:cursor-not-allowed"
                >
                  {canAfford ? 'Comprar' : 'Ouro insuficiente'}
                </button>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
