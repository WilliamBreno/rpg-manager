import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { enemyService, type EnemyAbilityInput, type EnemyInput, type EnemyLineInput } from '../services/enemyService'
import FileUploadField from './FileUploadField'
import type { EnemyKind } from '../types'

const KIND_LABEL: Record<EnemyKind, string> = { enemy: 'Inimigo', boss: 'Boss', villain: 'Vilão' }

const DICE_RE = /^(\d+)d(\d+)([+-]\d+)?$/i

// Só uma dica visual client-side (a validação de verdade é sempre no
// back-end, ver dice_notation.go) — mostra a média e compara com a faixa
// "Dano/Rodada" sugerida pro ND escolhido, mesma margem larga do back-end
// (metade do mínimo a dobro do máximo), já que uma única habilidade
// normalmente é só parte do dano total por rodada de um inimigo.
function damageHint(damage: string, cr: string, table: { cr: string; damage_per_round_min: number; damage_per_round_max: number }[] | undefined) {
  const m = DICE_RE.exec(damage.trim().toLowerCase())
  if (!m) return { text: damage ? 'Use notação de dado real, ex: 2d6+3' : '', bad: !!damage }
  const count = parseInt(m[1], 10)
  const sides = parseInt(m[2], 10)
  const mod = m[3] ? parseInt(m[3], 10) : 0
  const avg = Math.floor((count * (sides + 1)) / 2) + mod
  if (!cr || !table) return { text: `Média: ${avg}`, bad: false }
  const row = table.find(r => r.cr === cr)
  if (!row) return { text: `Média: ${avg}`, bad: false }
  const bad = avg < row.damage_per_round_min / 2 || avg > row.damage_per_round_max * 2
  return {
    text: `Média: ${avg} (ND ${cr} sugere ${row.damage_per_round_min}–${row.damage_per_round_max}/rodada)`,
    bad,
  }
}

export default function MasterEnemyForm({
  onSubmit,
  onCancel,
  submitting,
}: {
  onSubmit: (payload: EnemyInput) => void
  onCancel: () => void
  submitting: boolean
}) {
  const [kind, setKind] = useState<EnemyKind>('enemy')
  const [name, setName] = useState('')
  const [hp, setHp] = useState('')
  const [cr, setCr] = useState('')
  const [race, setRace] = useState('')
  const [klass, setKlass] = useState('')
  const [armor, setArmor] = useState('')
  const [photoUrl, setPhotoUrl] = useState('')
  const [soundUrl, setSoundUrl] = useState('')
  const [history, setHistory] = useState('')
  const [bonds, setBonds] = useState('')
  const [notes, setNotes] = useState('')
  const [abilities, setAbilities] = useState<EnemyAbilityInput[]>([{ name: '', damage: '', description: '' }])
  const [lines, setLines] = useState<EnemyLineInput[]>([])

  const { data: crTable } = useQuery({ queryKey: ['cr-damage-table'], queryFn: enemyService.getCRDamageTable })

  const isBossOrVillain = kind !== 'enemy'

  const updateAbility = (i: number, field: keyof EnemyAbilityInput, value: string) => {
    setAbilities(prev => prev.map((a, idx) => (idx === i ? { ...a, [field]: value } : a)))
  }
  const addAbility = () => setAbilities(prev => [...prev, { name: '', damage: '', description: '' }])
  const removeAbility = (i: number) => setAbilities(prev => prev.filter((_, idx) => idx !== i))

  const updateLine = (i: number, field: keyof EnemyLineInput, value: string) => {
    setLines(prev => prev.map((l, idx) => (idx === i ? { ...l, [field]: value } : l)))
  }
  const addLine = () => setLines(prev => [...prev, { text: '', audio_url: '', source: 'upload' }])
  const removeLine = (i: number) => setLines(prev => prev.filter((_, idx) => idx !== i))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({
      kind, name, hp: Number(hp) || 0, challenge_rating: cr, race,
      photo_url: photoUrl, sound_url: soundUrl, class: klass, armor: Number(armor) || 0,
      history: isBossOrVillain ? history : '', bonds: isBossOrVillain ? bonds : '', notes: isBossOrVillain ? notes : '',
      abilities: abilities.filter(a => a.name.trim() !== ''),
      lines: isBossOrVillain ? lines.filter(l => l.text.trim() !== '') : [],
    })
  }

  return (
    <form onSubmit={handleSubmit} className="rpg-card p-5 mb-6 flex flex-col gap-4">
      <div className="flex gap-2">
        {(['enemy', 'boss', 'villain'] as EnemyKind[]).map(k => (
          <button key={k} type="button" onClick={() => setKind(k)}
            className={`px-3 py-1.5 rounded-lg text-xs font-medium transition ${kind === k ? 'bg-rpg-gold-muted text-rpg-gold' : 'bg-gray-800 text-gray-400 hover:text-gray-200'}`}>
            {KIND_LABEL[k]}
          </button>
        ))}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Nome</label>
          <input value={name} onChange={e => setName(e.target.value)} required className="rpg-input" />
        </div>
        <div>
          <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Raça</label>
          <input value={race} onChange={e => setRace(e.target.value)} className="rpg-input" />
        </div>
        <div>
          <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Pontos de Vida</label>
          <input type="number" value={hp} onChange={e => setHp(e.target.value)} className="rpg-input" />
        </div>
        <div>
          <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">CA (armadura)</label>
          <input type="number" value={armor} onChange={e => setArmor(e.target.value)} className="rpg-input" />
        </div>
        <div>
          <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Classe (opcional)</label>
          <input value={klass} onChange={e => setKlass(e.target.value)} className="rpg-input" />
        </div>
        <div>
          <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Nível de Desafio (ND)</label>
          <input list="cr-options" value={cr} onChange={e => setCr(e.target.value)} className="rpg-input" placeholder="ex: 1/4, 3, 12" />
          <datalist id="cr-options">
            {(crTable ?? []).map(r => <option key={r.cr} value={r.cr} />)}
          </datalist>
          <p className="text-gray-600 text-xs mt-1">Opcional — usado só pra sugerir faixa de dano.</p>
        </div>
        <FileUploadField label="Foto" kind="image" value={photoUrl} onChange={setPhotoUrl} />
        <FileUploadField label="Som" kind="audio" value={soundUrl} onChange={setSoundUrl} />
      </div>

      {isBossOrVillain && (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div>
            <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">História</label>
            <textarea value={history} onChange={e => setHistory(e.target.value)} rows={2} className="rpg-input resize-none" />
          </div>
          <div>
            <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Vínculos</label>
            <textarea value={bonds} onChange={e => setBonds(e.target.value)} rows={2} className="rpg-input resize-none" />
          </div>
          <div>
            <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">Observações</label>
            <textarea value={notes} onChange={e => setNotes(e.target.value)} rows={2} className="rpg-input resize-none" />
          </div>
        </div>
      )}

      <div>
        <div className="flex items-center justify-between mb-2">
          <label className="text-gray-500 text-xs uppercase tracking-wider">Habilidades customizadas</label>
          <button type="button" onClick={addAbility} className="text-xs text-rpg-gold hover:underline">+ adicionar</button>
        </div>
        <div className="flex flex-col gap-2">
          {abilities.map((a, i) => {
            const hint = damageHint(a.damage, cr, crTable)
            return (
              <div key={i} className="flex flex-col sm:flex-row gap-2 bg-gray-800/40 p-2 rounded-lg">
                <input value={a.name} onChange={e => updateAbility(i, 'name', e.target.value)}
                  placeholder="Nome da habilidade" className="rpg-input flex-1" />
                <div className="flex-1">
                  <input value={a.damage} onChange={e => updateAbility(i, 'damage', e.target.value)}
                    placeholder="Dano (ex: 2d6+3)" className="rpg-input w-full" />
                  {hint.text && <p className={`text-xs mt-1 ${hint.bad ? 'text-amber-400' : 'text-gray-500'}`}>{hint.text}</p>}
                </div>
                <input value={a.description} onChange={e => updateAbility(i, 'description', e.target.value)}
                  placeholder="Descrição (opcional)" className="rpg-input flex-1" />
                <button type="button" onClick={() => removeAbility(i)} className="text-red-400 text-xs px-2 self-start">Remover</button>
              </div>
            )
          })}
        </div>
      </div>

      {isBossOrVillain && (
        <div>
          <div className="flex items-center justify-between mb-2">
            <label className="text-gray-500 text-xs uppercase tracking-wider">Falas</label>
            <button type="button" onClick={addLine} className="text-xs text-rpg-gold hover:underline">+ adicionar</button>
          </div>
          <div className="flex flex-col gap-2">
            {lines.map((l, i) => (
              <div key={i} className="flex flex-col gap-2 bg-gray-800/40 p-3 rounded-lg">
                <div className="flex flex-col sm:flex-row gap-2">
                  <input value={l.text} onChange={e => updateLine(i, 'text', e.target.value)}
                    placeholder="Texto da fala" className="rpg-input flex-1" />
                  <select value={l.source} onChange={e => updateLine(i, 'source', e.target.value)} className="rpg-input sm:w-44">
                    <option value="upload">Áudio enviado</option>
                    <option value="tts">Gerado por TTS</option>
                  </select>
                  <button type="button" onClick={() => removeLine(i)} className="text-red-400 text-xs px-2 flex-shrink-0 self-center">Remover</button>
                </div>
                <FileUploadField label="Áudio da fala" kind="audio" value={l.audio_url} onChange={v => updateLine(i, 'audio_url', v)} />
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="flex gap-2">
        <button type="submit" disabled={submitting} className="btn-rpg-primary px-4 py-2 text-sm">
          {submitting ? 'Salvando...' : `Criar ${KIND_LABEL[kind]}`}
        </button>
        <button type="button" onClick={onCancel} className="btn-rpg-outline px-4 py-2 text-sm">Cancelar</button>
      </div>
    </form>
  )
}
