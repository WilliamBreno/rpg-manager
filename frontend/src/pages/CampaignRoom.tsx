import { useEffect, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { campaignService } from '../services/campaignService'
import { sessionService } from '../services/sessionService'
import { sceneService } from '../services/sceneService'
import { chatService, type ChatMessage } from '../services/chatService'
import { membershipService } from '../services/membershipService'
import { useCampaignSocket } from '../lib/useCampaignSocket'
import { useAuthStore } from '../store/authStore'
import MasterSceneCanvas from '../components/MasterSceneCanvas'

// Sala ao vivo — Sistema do Mestre, Etapa 5/6. Usada tanto pelo mestre
// quanto pelo jogador: mostra o cenário que a sessão ativa aponta como
// ativo (read-only aqui — mover token continua um recurso só do mestre, na
// tela de gerenciamento), o chat em tempo real, e o elenco da campanha
// (resolve a ambiguidade "adicionar jogadores" do documento original como
// "ver quem mais está na campanha", não um sistema de contatos/amigos — ver
// CLAUDE.md pra a justificativa completa).
export default function CampaignRoom() {
  const { id } = useParams()
  const campaignId = Number(id)
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const user = useAuthStore(s => s.user)

  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [draft, setDraft] = useState('')
  const [nowPlaying, setNowPlaying] = useState('')
  const chatEndRef = useRef<HTMLDivElement>(null)
  const musicRef = useRef<HTMLAudioElement>(null)
  const sfxRef = useRef<HTMLAudioElement>(null)

  const { data: campaign } = useQuery({ queryKey: ['campaigns', campaignId], queryFn: () => campaignService.getByID(campaignId) })
  const { data: sessions } = useQuery({ queryKey: ['sessions', campaignId], queryFn: () => sessionService.getByCampaign(campaignId) })
  const { data: members } = useQuery({ queryKey: ['members', campaignId], queryFn: () => membershipService.getMembers(campaignId) })
  const { data: history } = useQuery({ queryKey: ['chat', campaignId], queryFn: () => chatService.getHistory(campaignId) })

  const activeSession = (sessions ?? []).find(s => !s.ended_at)
  const { data: activeScene } = useQuery({
    queryKey: ['scene', activeSession?.active_scene_id],
    queryFn: () => sceneService.getByID(activeSession!.active_scene_id as number),
    enabled: !!activeSession?.active_scene_id,
  })

  useEffect(() => { setMessages(history ?? []) }, [history])
  useEffect(() => { chatEndRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  useCampaignSocket(campaignId, event => {
    if (event.type === 'chat_message') {
      setMessages(prev => [...prev, event.data as ChatMessage])
    }
    if (event.type === 'scene_changed') {
      queryClient.invalidateQueries({ queryKey: ['sessions', campaignId] })
    }
    if (event.type === 'play_audio') {
      const data = event.data as { kind: string; url: string; playing?: boolean; text?: string; enemy_name?: string }
      if (data.kind === 'music') {
        if (musicRef.current) {
          if (data.playing) {
            musicRef.current.src = data.url
            musicRef.current.loop = true
            musicRef.current.play().catch(() => {})
          } else {
            musicRef.current.pause()
          }
        }
      } else {
        if (sfxRef.current) {
          sfxRef.current.src = data.url
          sfxRef.current.loop = false
          sfxRef.current.play().catch(() => {})
        }
        setNowPlaying(data.kind === 'line' ? `${data.enemy_name}: "${data.text}"` : `Som de ${data.enemy_name}`)
      }
    }
  })

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!draft.trim()) return
    const text = draft
    setDraft('')
    await chatService.send(campaignId, text, activeSession?.ID)
  }

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-6xl mx-auto">
        <button onClick={() => navigate(-1)} className="transition mb-4 block text-sm" style={{ color: 'rgba(201,168,76,0.5)' }}>← Voltar</button>
        <h1 className="font-rpg text-2xl font-bold mb-1" style={{ color: '#c9a84c' }}>{campaign?.name ?? 'Sala'}</h1>
        <p className="text-gray-500 text-sm mb-1">
          {activeSession ? 'Sessão em andamento' : 'Nenhuma sessão em andamento agora'}
        </p>
        {nowPlaying && <p className="text-rpg-gold text-xs mb-6">🔊 {nowPlaying}</p>}
        {!nowPlaying && <div className="mb-6" />}
        <audio ref={musicRef} />
        <audio ref={sfxRef} />

        <div className="flex flex-col lg:flex-row gap-6">
          <div className="flex-1 min-w-0">
            <h2 className="text-white font-semibold text-sm mb-2 uppercase tracking-wider">Cenário Ativo</h2>
            {activeScene ? (
              <MasterSceneCanvas scene={activeScene} readOnly onMoveToken={() => {}} onDeleteToken={() => {}} />
            ) : (
              <div className="rpg-card p-8 text-center text-gray-500 text-sm">
                {activeSession ? 'O mestre ainda não escolheu um cenário ativo.' : 'Aguardando o mestre abrir a sessão.'}
              </div>
            )}

            <h2 className="text-white font-semibold text-sm mt-6 mb-2 uppercase tracking-wider">Elenco da Campanha</h2>
            <div className="flex flex-wrap gap-2">
              <span className="text-xs px-3 py-1 rounded-full bg-rpg-gold-muted text-rpg-gold">
                {campaign?.master_id === user?.id ? 'Você (Mestre)' : 'Mestre'}
              </span>
              {(members ?? []).filter(m => m.status === 'accepted').map(m => (
                <span key={m.ID} className="text-xs px-3 py-1 rounded-full bg-gray-700/60 text-gray-300">
                  {m.user?.name ?? `Jogador #${m.user_id}`}
                </span>
              ))}
            </div>
          </div>

          <div className="lg:w-80 flex-shrink-0 flex flex-col rpg-card p-3 h-[500px]">
            <h2 className="text-white font-semibold text-sm mb-2 uppercase tracking-wider">Chat</h2>
            <div className="flex-1 overflow-y-auto flex flex-col gap-2 mb-2 pr-1">
              {messages.map(m => (
                <div key={m.ID} className="text-sm">
                  <span className="text-rpg-gold font-semibold">{m.sender?.name ?? '...'}: </span>
                  <span className="text-gray-300">{m.text}</span>
                </div>
              ))}
              <div ref={chatEndRef} />
            </div>
            <form onSubmit={sendMessage} className="flex gap-2">
              <input value={draft} onChange={e => setDraft(e.target.value)} placeholder="Mensagem..." className="rpg-input text-sm flex-1" />
              <button type="submit" className="btn-rpg-primary px-3 py-1.5 text-sm">Enviar</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
