import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { membershipService } from '../services/membershipService'
import { characterService } from '../services/characterService'

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
  // Personagens do próprio jogador — a escolha de qual representa ele nesta
  // campanha agora é obrigatória (o back-end recusa aceitar sem character_id).
  const { data: myCharacters } = useQuery({ queryKey: ['characters'], queryFn: characterService.getAll })

  // Personagem escolhido por convite, antes de confirmar "Aceitar".
  const [chosenCharacter, setChosenCharacter] = useState<Record<number, number>>({})

  const respond = useMutation({
    mutationFn: ({ id, accept, characterId }: { id: number; accept: boolean; characterId?: number }) =>
      membershipService.respond(id, accept, characterId),
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
          {(pending ?? []).map(inv => {
            const selectedId = chosenCharacter[inv.ID]
            return (
              <div key={inv.ID} className="rpg-card p-4 flex flex-col gap-3">
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div className="min-w-0 flex-1">
                    <p className="text-white font-semibold truncate">{inv.campaign?.name ?? `Campanha #${inv.campaign_id}`}</p>
                    <p className="text-gray-500 text-xs">Convite de mesa recebido</p>
                  </div>
                  <button onClick={() => respond.mutate({ id: inv.ID, accept: false })} disabled={respond.isPending}
                    className="text-xs px-3 py-1.5 rounded-lg border border-red-800/50 text-red-400 hover:bg-red-900/20 transition flex-shrink-0">Recusar</button>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  <select
                    value={selectedId ?? ''}
                    onChange={e => setChosenCharacter(prev => ({ ...prev, [inv.ID]: Number(e.target.value) }))}
                    className="rpg-input text-sm flex-1 min-w-[180px]"
                  >
                    <option value="">Escolha o personagem para esta campanha...</option>
                    {(myCharacters ?? []).map(c => (
                      <option key={c.ID} value={c.ID}>{c.name} — {c.class?.name} (Nv. {c.level}, {c.edition})</option>
                    ))}
                  </select>
                  <button
                    onClick={() => respond.mutate({ id: inv.ID, accept: true, characterId: selectedId })}
                    disabled={respond.isPending || !selectedId}
                    className="btn-rpg-primary px-3 py-1.5 text-xs flex-shrink-0"
                  >Aceitar</button>
                </div>
                {(myCharacters ?? []).length === 0 && (
                  <p className="text-amber-400/80 text-xs">Você ainda não tem nenhum personagem — crie um antes de aceitar este convite.</p>
                )}
              </div>
            )
          })}
        </div>

        {(myCampaigns ?? []).length > 0 && (
          <>
            <h2 className="font-rpg text-lg font-bold mt-8 mb-3" style={{ color: '#c9a84c' }}>Minhas Mesas</h2>
            <div className="flex flex-col gap-3">
              {(myCampaigns ?? []).map(m => (
                <div key={m.ID} className="rpg-card p-4 flex flex-wrap items-center justify-between gap-3">
                  <div className="min-w-0 flex-1">
                    <p className="text-white font-semibold truncate">{m.campaign?.name ?? `Campanha #${m.campaign_id}`}</p>
                    {m.character && <p className="text-gray-500 text-xs truncate">Jogando com {m.character.name}</p>}
                  </div>
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
