import { useEffect, useRef } from 'react'
import { useAuthStore } from '../store/authStore'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export interface CampaignSocketEvent {
  type: string
  data: unknown
}

// Conecta na sala ao vivo de uma campanha (Sistema do Mestre, Etapa 5) — o
// WebSocket aqui só recebe eventos que o servidor empurra (scene_changed,
// chat_message); toda escrita continua indo pelas rotas REST normais, que
// disparam o broadcast no back-end (ver internal/ws/hub.go). Reconecta
// automaticamente se a conexão cair (rede instável, cold start do Render).
export function useCampaignSocket(campaignId: number | null, onEvent: (event: CampaignSocketEvent) => void) {
  const token = useAuthStore(s => s.token)
  const onEventRef = useRef(onEvent)
  onEventRef.current = onEvent

  useEffect(() => {
    if (!campaignId || !token) return
    let socket: WebSocket | null = null
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null
    let closedByEffect = false

    const wsBase = API_URL.replace(/^http/, 'ws')

    const connect = () => {
      socket = new WebSocket(`${wsBase}/ws/campaign/${campaignId}?token=${encodeURIComponent(token)}`)
      socket.onmessage = e => {
        try {
          onEventRef.current(JSON.parse(e.data))
        } catch {
          // evento não-JSON (ping/controle) — ignora
        }
      }
      socket.onclose = () => {
        if (!closedByEffect) {
          reconnectTimer = setTimeout(connect, 2000)
        }
      }
    }
    connect()

    return () => {
      closedByEffect = true
      if (reconnectTimer) clearTimeout(reconnectTimer)
      socket?.close()
    }
  }, [campaignId, token])
}
