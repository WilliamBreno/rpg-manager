import { useEffect, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { campaignService } from '../services/campaignService'
import { npcService } from '../services/npcService'
import { enemyService } from '../services/enemyService'
import { sessionService } from '../services/sessionService'
import { sceneService } from '../services/sceneService'
import { membershipService } from '../services/membershipService'
import { chatService } from '../services/chatService'
import { rewardService } from '../services/rewardService'
import MasterEnemyForm from '../components/MasterEnemyForm'
import MasterSceneCanvas from '../components/MasterSceneCanvas'
import DiceRoller from '../components/DiceRoller'
import FileUploadField from '../components/FileUploadField'
import type { EnemyInput } from '../services/enemyService'
import type { EnemyKind } from '../types'

const KIND_LABEL: Record<EnemyKind, string> = { enemy: 'Inimigo', boss: 'Boss', villain: 'Vilão' }
const KIND_COLOR: Record<EnemyKind, string> = {
  enemy: 'bg-gray-700/60 text-gray-300',
  boss: 'bg-amber-900/40 text-amber-400',
  villain: 'bg-red-900/40 text-red-400',
}

// Sistema do Mestre — Etapa 2: NPCs, Inimigos, Boss e Vilão de uma campanha.
export default function MasterCampaignDetail() {
  const { id } = useParams()
  const campaignId = Number(id)
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [tab, setTab] = useState<'npcs' | 'enemies' | 'sessions' | 'scenes' | 'dice' | 'rewards'>('npcs')
  const tabBarRef = useRef<HTMLDivElement>(null)

  // Garante que a aba escolhida fique visível mesmo quando a barra de abas
  // rola horizontalmente no mobile (sem isso, tocar numa aba fora da área
  // visível troca o conteúdo mas o indicador ativo continua fora de vista).
  useEffect(() => {
    tabBarRef.current?.querySelector(`[data-tab="${tab}"]`)
      ?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'nearest' })
  }, [tab])
  const [editingSummaryId, setEditingSummaryId] = useState<number | null>(null)
  const [summaryDraft, setSummaryDraft] = useState('')
  const [showSceneForm, setShowSceneForm] = useState(false)
  const [sceneName, setSceneName] = useState('')
  const [sceneImageUrl, setSceneImageUrl] = useState('')
  const [openSceneId, setOpenSceneId] = useState<number | null>(null)
  const [tokenLabel, setTokenLabel] = useState('')
  const [inviteEmail, setInviteEmail] = useState('')
  const [inviteMsg, setInviteMsg] = useState('')
  const [musicUrl, setMusicUrl] = useState('')
  const [showItemForm, setShowItemForm] = useState(false)
  const [itemName, setItemName] = useState('')
  const [itemDescription, setItemDescription] = useState('')
  const [itemEffect, setItemEffect] = useState('')
  const [grantTarget, setGrantTarget] = useState<'all' | number>('all')
  const [grantGold, setGrantGold] = useState('')
  const [grantNote, setGrantNote] = useState('')
  const [grantItemId, setGrantItemId] = useState<number | ''>('')
  const [rewardMsg, setRewardMsg] = useState('')
  const [showNpcForm, setShowNpcForm] = useState(false)
  const [showEnemyForm, setShowEnemyForm] = useState(false)
  const [enemyWarnings, setEnemyWarnings] = useState<string[]>([])
  const [npcName, setNpcName] = useState('')
  const [npcHp, setNpcHp] = useState('')
  const [npcAlignment, setNpcAlignment] = useState('')
  const [npcPersonality, setNpcPersonality] = useState('')
  const [npcHistory, setNpcHistory] = useState('')
  const [npcBonds, setNpcBonds] = useState('')
  const [npcNotes, setNpcNotes] = useState('')

  const { data: campaign } = useQuery({ queryKey: ['campaigns', campaignId], queryFn: () => campaignService.getByID(campaignId) })
  const { data: npcs, isLoading: loadingNpcs } = useQuery({ queryKey: ['npcs', campaignId], queryFn: () => npcService.getByCampaign(campaignId) })
  const { data: enemies, isLoading: loadingEnemies } = useQuery({ queryKey: ['enemies', campaignId], queryFn: () => enemyService.getByCampaign(campaignId) })
  const { data: sessions, isLoading: loadingSessions } = useQuery({ queryKey: ['sessions', campaignId], queryFn: () => sessionService.getByCampaign(campaignId) })
  const { data: scenes, isLoading: loadingScenes } = useQuery({ queryKey: ['scenes', campaignId], queryFn: () => sceneService.getByCampaign(campaignId) })
  const { data: openScene } = useQuery({
    queryKey: ['scene', openSceneId],
    queryFn: () => sceneService.getByID(openSceneId as number),
    enabled: openSceneId !== null,
  })

  const activeSession = (sessions ?? []).find(s => !s.ended_at)

  const createScene = useMutation({
    mutationFn: () => sceneService.create(campaignId, { name: sceneName, image_url: sceneImageUrl }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['scenes', campaignId] })
      setSceneName(''); setSceneImageUrl(''); setShowSceneForm(false)
    },
  })

  const deleteScene = useMutation({
    mutationFn: (sceneId: number) => sceneService.delete(sceneId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['scenes', campaignId] })
      setOpenSceneId(null)
    },
  })

  const addToken = useMutation({
    mutationFn: () => sceneService.addToken(openSceneId as number, { label: tokenLabel, image_url: '', x: 50, y: 50 }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['scene', openSceneId] })
      setTokenLabel('')
    },
  })

  const moveToken = useMutation({
    mutationFn: ({ tokenId, x, y }: { tokenId: number; x: number; y: number }) => sceneService.moveToken(tokenId, x, y),
  })

  const deleteToken = useMutation({
    mutationFn: (tokenId: number) => sceneService.deleteToken(tokenId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['scene', openSceneId] }),
  })

  const setActiveScene = useMutation({
    mutationFn: (sceneId: number) => sceneService.setActiveScene((activeSession as { ID: number }).ID, sceneId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['sessions', campaignId] }),
  })

  const { data: members } = useQuery({ queryKey: ['members', campaignId], queryFn: () => membershipService.getMembers(campaignId) })
  const { data: magicItems } = useQuery({ queryKey: ['magic-items', campaignId], queryFn: () => rewardService.getMagicItems(campaignId) })
  const { data: rewardHistory } = useQuery({ queryKey: ['rewards', campaignId], queryFn: () => rewardService.getHistory(campaignId) })

  const createMagicItem = useMutation({
    mutationFn: () => rewardService.createMagicItem(campaignId, { name: itemName, description: itemDescription, effect: itemEffect }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['magic-items', campaignId] })
      setItemName(''); setItemDescription(''); setItemEffect(''); setShowItemForm(false)
    },
  })

  const deleteMagicItem = useMutation({
    mutationFn: (id: number) => rewardService.deleteMagicItem(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['magic-items', campaignId] }),
  })

  const grantCurrency = useMutation({
    mutationFn: () => rewardService.grantCurrency(campaignId, {
      ...(grantTarget === 'all' ? { all: true } : { character_id: grantTarget }),
      gold_pieces: Number(grantGold) || 0, note: grantNote,
    }),
    onSuccess: (rewards) => {
      queryClient.invalidateQueries({ queryKey: ['rewards', campaignId] })
      setGrantGold(''); setGrantNote('')
      setRewardMsg(`${rewards.length} personagem(ns) recebeu(ram) a recompensa em moeda.`)
    },
    onError: (err: unknown) => {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? 'Erro ao conceder recompensa.'
      setRewardMsg(message)
    },
  })

  const grantItem = useMutation({
    mutationFn: () => rewardService.grantItem(campaignId, {
      ...(grantTarget === 'all' ? { all: true } : { character_id: grantTarget }),
      magic_item_id: Number(grantItemId), note: grantNote,
    }),
    onSuccess: (rewards) => {
      queryClient.invalidateQueries({ queryKey: ['rewards', campaignId] })
      setGrantItemId(''); setGrantNote('')
      setRewardMsg(`${rewards.length} personagem(ns) recebeu(ram) o item.`)
    },
    onError: (err: unknown) => {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? 'Erro ao conceder item.'
      setRewardMsg(message)
    },
  })

  const invitePlayer = useMutation({
    mutationFn: () => membershipService.invite(campaignId, inviteEmail),
    onSuccess: () => { setInviteEmail(''); setInviteMsg('Convite enviado!') },
    onError: (err: unknown) => {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? 'Erro ao convidar.'
      setInviteMsg(message)
    },
  })

  const startSession = useMutation({
    mutationFn: () => sessionService.start(campaignId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['sessions', campaignId] }),
  })

  const saveSummary = useMutation({
    mutationFn: (sessionId: number) => sessionService.updateSummary(sessionId, summaryDraft),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions', campaignId] })
      setEditingSummaryId(null)
    },
  })

  const endSession = useMutation({
    mutationFn: (sessionId: number) => sessionService.end(sessionId, summaryDraft),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions', campaignId] })
      setEditingSummaryId(null)
    },
  })

  const createNpc = useMutation({
    mutationFn: () => npcService.create(campaignId, {
      name: npcName, hp: Number(npcHp) || 0, alignment: npcAlignment, personality: npcPersonality,
      history: npcHistory, bonds: npcBonds, notes: npcNotes,
    }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npcs', campaignId] })
      setNpcName(''); setNpcHp(''); setNpcAlignment(''); setNpcPersonality(''); setNpcHistory(''); setNpcBonds(''); setNpcNotes('')
      setShowNpcForm(false)
    },
  })

  const deleteNpc = useMutation({
    mutationFn: (npcId: number) => npcService.delete(npcId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['npcs', campaignId] }),
  })

  const createEnemy = useMutation({
    mutationFn: (payload: EnemyInput) => enemyService.create(campaignId, payload),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ['enemies', campaignId] })
      setEnemyWarnings(result.warnings ?? [])
      setShowEnemyForm(false)
    },
  })

  const deleteEnemy = useMutation({
    mutationFn: (enemyId: number) => enemyService.delete(enemyId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['enemies', campaignId] }),
  })

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-4xl mx-auto">
        <button onClick={() => navigate('/campaigns')} className="transition mb-4 block text-sm" style={{ color: 'rgba(201,168,76,0.5)' }}>← Voltar às campanhas</button>

        <div className="flex flex-wrap items-start justify-between gap-3 mb-1">
          <h1 className="font-rpg text-2xl sm:text-3xl font-bold" style={{ color: '#c9a84c' }}>{campaign?.name ?? 'Campanha'}</h1>
          <button onClick={() => navigate(`/campaigns/${campaignId}/room`)} className="btn-rpg-primary px-4 py-2 text-sm whitespace-nowrap flex-shrink-0">
            🔴 Entrar na Sala
          </button>
        </div>
        <p className="text-gray-500 text-sm mb-3">{campaign?.main_story}</p>

        <form onSubmit={e => { e.preventDefault(); invitePlayer.mutate() }} className="flex flex-col sm:flex-row gap-2 mb-6 max-w-md">
          <input type="email" value={inviteEmail} onChange={e => setInviteEmail(e.target.value)} required
            placeholder="E-mail do jogador a convidar" className="rpg-input text-sm flex-1" />
          <button type="submit" disabled={invitePlayer.isPending} className="btn-rpg-outline px-3 py-1.5 text-xs whitespace-nowrap flex-shrink-0">
            {invitePlayer.isPending ? 'Enviando...' : 'Convidar'}
          </button>
        </form>
        {inviteMsg && <p className="text-xs text-gray-400 -mt-4 mb-6">{inviteMsg}</p>}

        <div className="relative mb-6">
          <div ref={tabBarRef} className="flex gap-2 border-b border-gray-800 overflow-x-auto no-scrollbar">
            <button data-tab="npcs" onClick={() => setTab('npcs')} className={`flex-shrink-0 px-4 py-2 text-sm font-medium border-b-2 transition whitespace-nowrap ${tab === 'npcs' ? 'border-rpg-gold text-rpg-gold' : 'border-transparent text-gray-500 hover:text-gray-300'}`}>
              NPCs
            </button>
            <button data-tab="enemies" onClick={() => setTab('enemies')} className={`flex-shrink-0 px-4 py-2 text-sm font-medium border-b-2 transition whitespace-nowrap ${tab === 'enemies' ? 'border-rpg-gold text-rpg-gold' : 'border-transparent text-gray-500 hover:text-gray-300'}`}>
              Inimigos / Boss / Vilão
            </button>
            <button data-tab="sessions" onClick={() => setTab('sessions')} className={`flex-shrink-0 px-4 py-2 text-sm font-medium border-b-2 transition whitespace-nowrap ${tab === 'sessions' ? 'border-rpg-gold text-rpg-gold' : 'border-transparent text-gray-500 hover:text-gray-300'}`}>
              Sessões
            </button>
            <button data-tab="scenes" onClick={() => setTab('scenes')} className={`flex-shrink-0 px-4 py-2 text-sm font-medium border-b-2 transition whitespace-nowrap ${tab === 'scenes' ? 'border-rpg-gold text-rpg-gold' : 'border-transparent text-gray-500 hover:text-gray-300'}`}>
              Cenários
            </button>
            <button data-tab="dice" onClick={() => setTab('dice')} className={`flex-shrink-0 px-4 py-2 text-sm font-medium border-b-2 transition whitespace-nowrap ${tab === 'dice' ? 'border-rpg-gold text-rpg-gold' : 'border-transparent text-gray-500 hover:text-gray-300'}`}>
              Dados
            </button>
            <button data-tab="rewards" onClick={() => setTab('rewards')} className={`flex-shrink-0 px-4 py-2 text-sm font-medium border-b-2 transition whitespace-nowrap ${tab === 'rewards' ? 'border-rpg-gold text-rpg-gold' : 'border-transparent text-gray-500 hover:text-gray-300'}`}>
              Recompensas
            </button>
          </div>
          {/* Dica visual de que a barra rola horizontalmente no mobile — sem
              isso o corte abrupto no fim das abas parece um problema de
              layout em vez de "arraste pra ver mais". */}
          <div className="md:hidden pointer-events-none absolute right-0 top-0 bottom-[2px] w-8 bg-gradient-to-l from-gray-900 to-transparent" />
        </div>

        {tab === 'npcs' && (
          <div>
            <div className="flex justify-end mb-4">
              <button onClick={() => setShowNpcForm(v => !v)} className="btn-rpg-primary px-4 py-2 text-sm">
                {showNpcForm ? 'Cancelar' : '+ Novo NPC'}
              </button>
            </div>

            {showNpcForm && (
              <form onSubmit={e => { e.preventDefault(); createNpc.mutate() }} className="rpg-card p-5 mb-6 flex flex-col gap-3">
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <input value={npcName} onChange={e => setNpcName(e.target.value)} required placeholder="Nome" className="rpg-input" />
                  <input type="number" value={npcHp} onChange={e => setNpcHp(e.target.value)} placeholder="Pontos de Vida" className="rpg-input" />
                  <input value={npcAlignment} onChange={e => setNpcAlignment(e.target.value)} placeholder="Tendência" className="rpg-input" />
                </div>
                <input value={npcPersonality} onChange={e => setNpcPersonality(e.target.value)} placeholder="Personalidade" className="rpg-input" />
                <textarea value={npcHistory} onChange={e => setNpcHistory(e.target.value)} rows={2} placeholder="História" className="rpg-input resize-none" />
                <textarea value={npcBonds} onChange={e => setNpcBonds(e.target.value)} rows={2} placeholder="Vínculos" className="rpg-input resize-none" />
                <textarea value={npcNotes} onChange={e => setNpcNotes(e.target.value)} rows={2} placeholder="Observações" className="rpg-input resize-none" />
                <button type="submit" disabled={createNpc.isPending} className="btn-rpg-primary py-2 text-sm">
                  {createNpc.isPending ? 'Criando...' : 'Criar NPC'}
                </button>
              </form>
            )}

            {loadingNpcs && <p className="text-gray-500 text-sm">Carregando...</p>}
            {!loadingNpcs && (npcs ?? []).length === 0 && <p className="text-gray-500 text-sm">Nenhum NPC ainda.</p>}
            <div className="flex flex-col gap-3">
              {(npcs ?? []).map(n => (
                <div key={n.ID} className="rpg-card p-4 flex items-start justify-between gap-3">
                  <div className="flex-1 min-w-0">
                    <p className="text-white font-semibold">{n.name} <span className="text-gray-500 text-xs">PV {n.hp}</span></p>
                    {n.alignment && <p className="text-gray-500 text-xs">{n.alignment}</p>}
                    {n.personality && <p className="text-gray-400 text-sm mt-1">{n.personality}</p>}
                    {n.history && <p className="text-gray-500 text-xs mt-1 whitespace-pre-wrap">{n.history}</p>}
                  </div>
                  <button onClick={() => { if (confirm(`Excluir "${n.name}"?`)) deleteNpc.mutate(n.ID) }}
                    className="text-xs px-3 py-1 rounded-lg border border-red-800/50 text-red-400 hover:bg-red-900/20 transition flex-shrink-0">
                    Excluir
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {tab === 'enemies' && (
          <div>
            {enemyWarnings.length > 0 && (
              <div className="rpg-card p-3 mb-4 border border-amber-800/50 bg-amber-900/10">
                {enemyWarnings.map((w, i) => <p key={i} className="text-amber-400 text-xs">{w}</p>)}
              </div>
            )}

            <div className="flex justify-end mb-4">
              <button onClick={() => setShowEnemyForm(v => !v)} className="btn-rpg-primary px-4 py-2 text-sm">
                {showEnemyForm ? 'Cancelar' : '+ Novo Inimigo / Boss / Vilão'}
              </button>
            </div>

            {showEnemyForm && (
              <MasterEnemyForm
                submitting={createEnemy.isPending}
                onCancel={() => setShowEnemyForm(false)}
                onSubmit={payload => createEnemy.mutate(payload)}
              />
            )}

            {loadingEnemies && <p className="text-gray-500 text-sm">Carregando...</p>}
            {!loadingEnemies && (enemies ?? []).length === 0 && <p className="text-gray-500 text-sm">Nenhum inimigo ainda.</p>}
            <div className="flex flex-col gap-3">
              {(enemies ?? []).map(en => (
                <div key={en.ID} className="rpg-card p-4">
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <p className="text-white font-semibold">{en.name}</p>
                        <span className={`text-xs px-2 py-0.5 rounded-full ${KIND_COLOR[en.kind]}`}>{KIND_LABEL[en.kind]}</span>
                        {en.challenge_rating && <span className="text-xs text-gray-500">ND {en.challenge_rating}</span>}
                      </div>
                      <p className="text-gray-500 text-xs">PV {en.hp} · CA {en.armor} {en.race && `· ${en.race}`} {en.class && `· ${en.class}`}</p>
                      {en.abilities?.length > 0 && (
                        <ul className="text-gray-400 text-xs mt-2 flex flex-col gap-0.5">
                          {en.abilities.map(a => <li key={a.ID}>⚔ {a.name} — {a.damage}</li>)}
                        </ul>
                      )}
                      {en.lines && en.lines.length > 0 && (
                        <ul className="text-gray-500 text-xs mt-2 flex flex-col gap-1 italic">
                          {en.lines.map(l => (
                            <li key={l.ID} className="flex items-center gap-2 not-italic">
                              <span className="italic">"{l.text}"</span>
                              {l.audio_url && (
                                <button onClick={() => enemyService.playLine(l.ID)} className="text-rpg-gold text-xs hover:underline">▶ tocar</button>
                              )}
                            </li>
                          ))}
                        </ul>
                      )}
                      {en.sound_url && (
                        <button onClick={() => enemyService.playSound(en.ID)} className="text-rpg-gold text-xs mt-2 hover:underline">
                          🔊 Tocar som do inimigo
                        </button>
                      )}
                    </div>
                    <button onClick={() => { if (confirm(`Excluir "${en.name}"?`)) deleteEnemy.mutate(en.ID) }}
                      className="text-xs px-3 py-1 rounded-lg border border-red-800/50 text-red-400 hover:bg-red-900/20 transition flex-shrink-0">
                      Excluir
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {tab === 'sessions' && (
          <div>
            <div className="flex justify-end mb-4">
              <button
                onClick={() => startSession.mutate()}
                disabled={!!activeSession || startSession.isPending}
                className="btn-rpg-primary px-4 py-2 text-sm disabled:opacity-40"
              >
                {activeSession ? 'Já existe sessão em andamento' : (startSession.isPending ? 'Abrindo...' : '+ Abrir Sessão')}
              </button>
            </div>

            {activeSession && (
              <div className="rpg-card p-4 mb-4 flex flex-col gap-3">
                <FileUploadField label="Música de fundo" kind="audio" value={musicUrl} onChange={setMusicUrl} />
                <div className="flex gap-2">
                  <button type="button" disabled={!musicUrl} onClick={() => sessionService.setMusic(activeSession.ID, musicUrl, true)}
                    className="btn-rpg-primary px-3 py-1.5 text-xs disabled:opacity-40">▶ Tocar pros jogadores</button>
                  <button type="button" onClick={() => sessionService.setMusic(activeSession.ID, '', false)}
                    className="text-xs px-3 py-1.5 rounded-lg border border-gray-700 text-gray-400 hover:bg-gray-800 transition">⏹ Parar</button>
                </div>
              </div>
            )}

            {loadingSessions && <p className="text-gray-500 text-sm">Carregando...</p>}
            {!loadingSessions && (sessions ?? []).length === 0 && <p className="text-gray-500 text-sm">Nenhuma sessão registrada ainda.</p>}

            <div className="flex flex-col gap-3">
              {(sessions ?? []).map(s => {
                const editing = editingSummaryId === s.ID
                return (
                  <div key={s.ID} className="rpg-card p-4">
                    <div className="flex items-center justify-between gap-3 mb-2">
                      <p className="text-white font-semibold text-sm">
                        {new Date(s.started_at).toLocaleString('pt-BR')}
                        {!s.ended_at && <span className="ml-2 text-xs px-2 py-0.5 rounded-full bg-green-900/40 text-green-400">em andamento</span>}
                        {s.ended_at && <span className="ml-2 text-xs text-gray-500">até {new Date(s.ended_at).toLocaleString('pt-BR')}</span>}
                      </p>
                    </div>

                    {editing ? (
                      <div className="flex flex-col gap-2">
                        <textarea value={summaryDraft} onChange={e => setSummaryDraft(e.target.value)} rows={4}
                          className="rpg-input resize-none" placeholder="O que aconteceu nessa sessão..." />
                        <div className="flex gap-2">
                          <button onClick={() => saveSummary.mutate(s.ID)} disabled={saveSummary.isPending}
                            className="btn-rpg-primary px-3 py-1.5 text-xs">Salvar diário</button>
                          {!s.ended_at && (
                            <button onClick={() => endSession.mutate(s.ID)} disabled={endSession.isPending}
                              className="text-xs px-3 py-1.5 rounded-lg border border-red-800/50 text-red-400 hover:bg-red-900/20 transition">
                              Encerrar sessão
                            </button>
                          )}
                          <button onClick={() => setEditingSummaryId(null)} className="btn-rpg-outline px-3 py-1.5 text-xs">Cancelar</button>
                        </div>
                      </div>
                    ) : (
                      <div>
                        <p className="text-gray-400 text-sm whitespace-pre-wrap mb-2">{s.summary || 'Nenhum registro ainda.'}</p>
                        <button onClick={() => { setEditingSummaryId(s.ID); setSummaryDraft(s.summary) }}
                          className="btn-rpg-outline px-3 py-1 text-xs">
                          {s.ended_at ? 'Editar diário' : 'Registrar / encerrar'}
                        </button>
                      </div>
                    )}
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {tab === 'scenes' && (
          <div>
            <div className="flex justify-end mb-4">
              <button onClick={() => setShowSceneForm(v => !v)} className="btn-rpg-primary px-4 py-2 text-sm">
                {showSceneForm ? 'Cancelar' : '+ Novo Cenário'}
              </button>
            </div>

            {showSceneForm && (
              <form onSubmit={e => { e.preventDefault(); createScene.mutate() }} className="rpg-card p-5 mb-6 flex flex-col gap-3">
                <input value={sceneName} onChange={e => setSceneName(e.target.value)} required placeholder="Nome do cenário" className="rpg-input" />
                <FileUploadField label="Imagem de fundo" kind="image" value={sceneImageUrl} onChange={setSceneImageUrl} />
                <button type="submit" disabled={createScene.isPending} className="btn-rpg-primary py-2 text-sm">
                  {createScene.isPending ? 'Criando...' : 'Criar Cenário'}
                </button>
              </form>
            )}

            {loadingScenes && <p className="text-gray-500 text-sm">Carregando...</p>}
            {!loadingScenes && (scenes ?? []).length === 0 && <p className="text-gray-500 text-sm">Nenhum cenário na biblioteca ainda.</p>}

            <div className="flex flex-wrap gap-3 mb-6">
              {(scenes ?? []).map(sc => (
                <div key={sc.ID} className={`rpg-card p-3 w-48 cursor-pointer transition ${openSceneId === sc.ID ? 'ring-2 ring-rpg-gold' : ''}`}
                  onClick={() => setOpenSceneId(sc.ID)}>
                  <p className="text-white text-sm font-semibold truncate">{sc.name}</p>
                  {activeSession?.active_scene_id === sc.ID && (
                    <span className="text-xs text-green-400">● ativo na sessão</span>
                  )}
                  <div className="flex gap-2 mt-2">
                    {activeSession && (
                      <button onClick={e => { e.stopPropagation(); setActiveScene.mutate(sc.ID) }}
                        className="text-xs text-rpg-gold hover:underline">Ativar</button>
                    )}
                    <button onClick={e => { e.stopPropagation(); if (confirm(`Excluir "${sc.name}"?`)) deleteScene.mutate(sc.ID) }}
                      className="text-xs text-red-400 hover:underline">Excluir</button>
                  </div>
                </div>
              ))}
            </div>

            {openScene && (
              <div>
                <div className="flex flex-col sm:flex-row sm:items-center gap-2 mb-3">
                  <h3 className="text-white font-semibold">{openScene.name}</h3>
                  <form onSubmit={e => { e.preventDefault(); addToken.mutate() }} className="flex gap-2 sm:ml-auto">
                    <input value={tokenLabel} onChange={e => setTokenLabel(e.target.value)} placeholder="Nome do token" className="rpg-input text-sm py-1 flex-1 sm:flex-none" />
                    <button type="submit" disabled={!tokenLabel || addToken.isPending} className="btn-rpg-primary px-3 py-1 text-xs flex-shrink-0">+ Token</button>
                  </form>
                </div>
                <MasterSceneCanvas
                  scene={openScene}
                  onMoveToken={(tokenId, x, y) => moveToken.mutate({ tokenId, x, y })}
                  onDeleteToken={tokenId => deleteToken.mutate(tokenId)}
                />
                <p className="text-gray-600 text-xs mt-2">Arraste um token pra mover. Duplo-clique num token pra removê-lo.</p>
              </div>
            )}
          </div>
        )}

        {tab === 'dice' && (
          <DiceRoller onShare={text => { chatService.send(campaignId, text, activeSession?.ID) }} />
        )}

        {tab === 'rewards' && (
          <div>
            <div className="flex items-center gap-2 mb-2">
              <label className="text-gray-500 text-xs uppercase tracking-wider">Destinatário</label>
            </div>
            <select
              value={grantTarget}
              onChange={e => setGrantTarget(e.target.value === 'all' ? 'all' : Number(e.target.value))}
              className="rpg-input text-sm mb-4 max-w-xs"
            >
              <option value="all">Todos os jogadores da campanha</option>
              {(members ?? []).filter(m => m.status === 'accepted' && m.character).map(m => (
                <option key={m.ID} value={m.character!.ID}>{m.character!.name} ({m.user?.name})</option>
              ))}
            </select>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 mb-6">
              <form onSubmit={e => { e.preventDefault(); grantCurrency.mutate() }} className="rpg-card p-4 flex flex-col gap-2">
                <h3 className="text-white font-semibold text-sm mb-1">Conceder Moeda</h3>
                <input type="number" value={grantGold} onChange={e => setGrantGold(e.target.value)} placeholder="Peças de ouro (PO)" className="rpg-input text-sm" />
                <input value={grantNote} onChange={e => setGrantNote(e.target.value)} placeholder="Nota (opcional)" className="rpg-input text-sm" />
                <button type="submit" disabled={grantCurrency.isPending || !grantGold} className="btn-rpg-primary py-1.5 text-sm">
                  {grantCurrency.isPending ? 'Concedendo...' : 'Conceder Moeda'}
                </button>
              </form>

              <form onSubmit={e => { e.preventDefault(); grantItem.mutate() }} className="rpg-card p-4 flex flex-col gap-2">
                <h3 className="text-white font-semibold text-sm mb-1">Conceder Item Mágico</h3>
                <select value={grantItemId} onChange={e => setGrantItemId(e.target.value ? Number(e.target.value) : '')} className="rpg-input text-sm">
                  <option value="">Selecione um item...</option>
                  {(magicItems ?? []).map(mi => <option key={mi.ID} value={mi.ID}>{mi.name}</option>)}
                </select>
                <button type="submit" disabled={grantItem.isPending || !grantItemId} className="btn-rpg-primary py-1.5 text-sm">
                  {grantItem.isPending ? 'Concedendo...' : 'Conceder Item'}
                </button>
              </form>
            </div>
            {rewardMsg && <p className="text-xs text-gray-400 mb-6">{rewardMsg}</p>}

            <div className="flex flex-wrap items-center justify-between gap-2 mb-3">
              <h3 className="text-white font-semibold text-sm uppercase tracking-wider">Catálogo de Itens Mágicos</h3>
              <button onClick={() => setShowItemForm(v => !v)} className="btn-rpg-outline px-3 py-1 text-xs">
                {showItemForm ? 'Cancelar' : '+ Novo Item'}
              </button>
            </div>
            {showItemForm && (
              <form onSubmit={e => { e.preventDefault(); createMagicItem.mutate() }} className="rpg-card p-4 mb-4 flex flex-col gap-2">
                <input value={itemName} onChange={e => setItemName(e.target.value)} required placeholder="Nome do item" className="rpg-input text-sm" />
                <textarea value={itemDescription} onChange={e => setItemDescription(e.target.value)} rows={2} placeholder="Descrição" className="rpg-input text-sm resize-none" />
                <input value={itemEffect} onChange={e => setItemEffect(e.target.value)} placeholder="Efeito (ex: +1 no ataque)" className="rpg-input text-sm" />
                <button type="submit" disabled={createMagicItem.isPending} className="btn-rpg-primary py-1.5 text-sm">Criar Item</button>
              </form>
            )}
            <div className="flex flex-col gap-2 mb-8">
              {(magicItems ?? []).map(mi => (
                <div key={mi.ID} className="rpg-card p-3 flex flex-wrap items-center justify-between gap-3">
                  <div className="min-w-0 flex-1">
                    <p className="text-white text-sm font-semibold truncate">{mi.name}</p>
                    <p className="text-gray-500 text-xs">{mi.effect}</p>
                  </div>
                  <button onClick={() => deleteMagicItem.mutate(mi.ID)} className="text-xs text-red-400 hover:underline flex-shrink-0">Excluir</button>
                </div>
              ))}
            </div>

            <h3 className="text-white font-semibold text-sm uppercase tracking-wider mb-3">Histórico de Recompensas</h3>
            <div className="flex flex-col gap-2">
              {(rewardHistory ?? []).length === 0 && <p className="text-gray-500 text-sm">Nenhuma recompensa concedida ainda.</p>}
              {(rewardHistory ?? []).map(r => (
                <div key={r.ID} className="rpg-card p-3 text-sm">
                  <span className="text-gray-300">{r.character?.name ?? `Personagem #${r.character_id}`}</span>
                  <span className="text-gray-500"> recebeu </span>
                  <span className="text-rpg-gold">
                    {r.kind === 'currency' ? `${r.gold_pieces}po/${r.platinum_pieces}pp/${r.electrum_pieces}pe/${r.silver_pieces}ps/${r.copper_pieces}pc` : r.magic_item?.name}
                  </span>
                  {r.note && <span className="text-gray-600 text-xs"> — {r.note}</span>}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
