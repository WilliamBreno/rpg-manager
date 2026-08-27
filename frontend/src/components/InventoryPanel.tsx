import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { characterService } from '../services/characterService'
import { parseWeaponDescription, damageAbilityMod, formatSigned, acFor } from '../lib/weaponParser'
import type { Currency } from '../types'

const COIN_LABELS: { key: keyof Currency; label: string }[] = [
  { key: 'platinum_pieces', label: 'PL' },
  { key: 'gold_pieces',     label: 'PO' },
  { key: 'electrum_pieces', label: 'PE' },
  { key: 'silver_pieces',   label: 'PP' },
  { key: 'copper_pieces',   label: 'PC' },
]

interface Props {
  characterId: number
  currency: Currency
  strengthMod: number
  dexterityMod: number
}

export default function InventoryPanel({ characterId, currency, strengthMod, dexterityMod }: Props) {
  const navigate = useNavigate()

  const { data: inventory } = useQuery({
    queryKey: ['inventory', characterId],
    queryFn: () => characterService.getInventory(characterId),
  })

  const hasItems = (inventory?.items?.length ?? 0) > 0 || (inventory?.armors?.length ?? 0) > 0

  return (
    <div className="bg-gray-800 rounded-lg p-4 sm:p-6 mb-4 border border-gray-700">
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-base sm:text-lg font-semibold text-white">Inventário</h2>
        <button onClick={() => navigate(`/characters/${characterId}/shop`)} className="btn-rpg-outline text-xs px-3 py-1.5">
          🛒 Ir para a loja
        </button>
      </div>

      <div className="flex gap-2 flex-wrap mb-4">
        {COIN_LABELS.map(({ key, label }) => (
          <div key={key} className="bg-gray-900 rounded-lg px-3 py-1.5 text-center" style={{ minWidth: 52 }}>
            <div className="text-rpg-gold font-bold text-sm">{currency[key]}</div>
            <div className="text-gray-500 text-[10px] uppercase">{label}</div>
          </div>
        ))}
      </div>

      {!hasItems && <p className="text-gray-500 text-sm">Nenhum item comprado ainda.</p>}

      {hasItems && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
          {inventory?.armors?.map(a => {
            const ca = acFor(a.armor.base_ac, a.armor.max_dex_bonus, dexterityMod)
            return (
              <div key={`armor-${a.armor_id}`} className="bg-gray-900 rounded-lg px-3 py-2 flex justify-between text-sm">
                <span className="text-gray-200">🛡️ {a.armor.name} <span className="text-gray-500">— CA {ca}</span></span>
                <span className="text-gray-500">x{a.quantity}</span>
              </div>
            )
          })}
          {inventory?.items?.map(i => {
            const weapon = i.item.category === 'arma' ? parseWeaponDescription(i.item.description) : null
            return (
              <div key={`item-${i.item_id}`} className="bg-gray-900 rounded-lg px-3 py-2 flex justify-between text-sm">
                <span className="text-gray-200">
                  {i.item.category === 'item_magico' ? '✨' : i.item.category === 'arma' ? '⚔️' : '🎒'} {i.item.name}
                  {weapon && (
                    <span className="text-gray-500"> — {weapon.dice} {weapon.damageType} ({formatSigned(damageAbilityMod(weapon, strengthMod, dexterityMod))})</span>
                  )}
                </span>
                <span className="text-gray-500">x{i.quantity}</span>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
