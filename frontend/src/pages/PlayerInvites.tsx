import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { membershipService } from '../services/membershipService'

// Área "adicionar mestre(s)" do lado do jogador (Etapa 6 do
// SISTEMA_MESTRE.md) — aceitar/recusar convite de campanha. Distinta da
// área "adicionar jogadores" (ambígua no documento original, resolvida como
// "ver o elenco da campanha" — ver o painel de elenco em CampaignRoom.tsx e
// a nota em CLAUDE.md).
export default function PlayerInvites() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: pending, isLoading } = useQuery({ queryKey: ['my-invites'], queryFn: membershipService.getMyPending })
  const { data: myCampaigns } = useQuery({ queryKey: ['my-campaigns'], queryFn: membershipService.getMyCampaigns })

  const respond = useMutation({
    mutationFn: ({ id, accept }: { id: number; accept: boolean }) => membershipService.respond(id, accept),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-invites'] })
      queryClient.invalidateQueries({ queryKey: ['my-campaigns'] })
    },
  })

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-2xl mx-auto">
        <button onClick={() => navigate('/characters')} className="transition mb-6 block text-sm" style={{ color: 'rgba(201,168,76,0.5)' }}>← Voltar</button>
        <h1 className="font-rpg text-2xl font-bold mb-6" style={{ color: '#c9a84c' }}>Convites de Campanha</h1>

        {isLoading && <p className="text-gray-500 text-sm">Carregando...</p>}
        {!isLoading && (pending ?? []).length === 0 && (
          <p className="text-gray-500 text-sm">Nenhum convite pendente no momento.</p>
        )}

        <div className="flex flex-col gap-3">
          {(pending ?? []).map(inv => (
            <div key={inv.ID} className="rpg-card p-4 flex flex-wrap items-center justify-between gap-3">
              <div className="min-w-0 flex-1">
                <p className="text-white font-semibold truncate">{inv.campaign?.name ?? `Campanha #${inv.campaign_id}`}</p>
                <p className="text-gray-500 text-xs">Convite de mesa recebido</p>
              </div>
              <div className="flex gap-2 flex-shrink-0">
                <button onClick={() => respond.mutate({ id: inv.ID, accept: true })} disabled={respond.isPending}
                  className="btn-rpg-primary px-3 py-1.5 text-xs">Aceitar</button>
                <button onClick={() => respond.mutate({ id: inv.ID, accept: false })} disabled={respond.isPending}
                  className="text-xs px-3 py-1.5 rounded-lg border border-red-800/50 text-red-400 hover:bg-red-900/20 transition">Recusar</button>
              </div>
            </div>
          ))}
        </div>

        {(myCampaigns ?? []).length > 0 && (
          <>
            <h2 className="font-rpg text-lg font-bold mt-8 mb-3" style={{ color: '#c9a84c' }}>Minhas Mesas</h2>
            <div className="flex flex-col gap-3">
              {(myCampaigns ?? []).map(m => (
                <div key={m.ID} className="rpg-card p-4 flex flex-wrap items-center justify-between gap-3">
                  <p className="text-white font-semibold min-w-0 flex-1 truncate">{m.campaign?.name ?? `Campanha #${m.campaign_id}`}</p>
                  <button onClick={() => navigate(`/campaigns/${m.campaign_id}/room`)} className="btn-rpg-primary px-3 py-1.5 text-xs flex-shrink-0">
                    🔴 Entrar na Sala
                  </button>
                </div>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  )
}
