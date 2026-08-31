import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { campaignService } from '../services/campaignService'
import { useAuthStore } from '../store/authStore'
import type { Campaign } from '../types'

// Sistema do Mestre — Etapa 1 (CRUD de campanha). Edição fixa em 5e por
// enquanto (campo já existe no backend, só não expor a opção 4e na UI ainda,
// conforme SISTEMA_MESTRE.md).
export default function MasterCampaigns() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const user = useAuthStore(s => s.user)

  const [showCreate, setShowCreate] = useState(false)
  const [name, setName] = useState('')
  const [mainStory, setMainStory] = useState('')
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState('')
  const [editStory, setEditStory] = useState('')

  const { data: campaigns, isLoading } = useQuery({
    queryKey: ['campaigns'],
    queryFn: campaignService.getAll,
    enabled: user?.role === 'master',
  })

  const createMutation = useMutation({
    mutationFn: () => campaignService.create({ name, main_story: mainStory }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      setName(''); setMainStory(''); setShowCreate(false)
    },
  })

  const updateMutation = useMutation({
    mutationFn: (id: number) => campaignService.update(id, { name: editName, main_story: editStory }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
      setEditingId(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => campaignService.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['campaigns'] }),
  })

  const startEdit = (c: Campaign) => {
    setEditingId(c.ID); setEditName(c.name); setEditStory(c.main_story)
  }

  if (user && user.role !== 'master') {
    return (
      <div className="min-h-screen bg-gray-900 px-4 py-8 flex items-center justify-center">
        <div className="rpg-card p-6 text-center max-w-md">
          <p className="text-gray-300">Essa área é só para contas de Mestre.</p>
          <button onClick={() => navigate('/characters')} className="btn-rpg-outline mt-4 px-4 py-2">
            Voltar
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-900 px-4 py-6 sm:px-8 sm:py-8">
      <div className="max-w-3xl mx-auto">
        <button onClick={() => navigate('/characters')} className="transition mb-6 block text-sm"
          style={{ color: 'rgba(201,168,76,0.5)' }}>← Voltar</button>

        <div className="flex items-center justify-between mb-6">
          <h1 className="font-rpg text-2xl sm:text-3xl font-bold" style={{ color: '#c9a84c' }}>
            Minhas Campanhas
          </h1>
          <button onClick={() => setShowCreate(v => !v)} className="btn-rpg-primary px-4 py-2 text-sm">
            {showCreate ? 'Cancelar' : '+ Nova Campanha'}
          </button>
        </div>

        {showCreate && (
          <form onSubmit={e => { e.preventDefault(); createMutation.mutate() }}
            className="rpg-card p-5 mb-6 flex flex-col gap-3">
            <div>
              <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">Nome da Campanha</label>
              <input value={name} onChange={e => setName(e.target.value)} required
                className="rpg-input" placeholder="Ex: A Sombra sobre Baldur" />
            </div>
            <div>
              <label className="text-gray-500 text-xs mb-1.5 block uppercase tracking-wider">História Principal</label>
              <textarea value={mainStory} onChange={e => setMainStory(e.target.value)} rows={4}
                className="rpg-input resize-none" placeholder="O enredo geral da campanha..." />
            </div>
            <p className="text-gray-600 text-xs">Edição: D&D 5e</p>
            <button type="submit" disabled={createMutation.isPending} className="btn-rpg-primary py-2 text-sm">
              {createMutation.isPending ? 'Criando...' : 'Criar Campanha'}
            </button>
          </form>
        )}

        {isLoading && <p className="text-gray-500 text-sm">Carregando...</p>}
        {!isLoading && (campaigns ?? []).length === 0 && !showCreate && (
          <p className="text-gray-500 text-sm">Nenhuma campanha ainda — crie a primeira acima.</p>
        )}

        <div className="flex flex-col gap-3">
          {(campaigns ?? []).map(c => (
            <div key={c.ID} className="rpg-card p-4">
              {editingId === c.ID ? (
                <div className="flex flex-col gap-3">
                  <input value={editName} onChange={e => setEditName(e.target.value)} className="rpg-input" />
                  <textarea value={editStory} onChange={e => setEditStory(e.target.value)} rows={3} className="rpg-input resize-none" />
                  <div className="flex gap-2">
                    <button onClick={() => updateMutation.mutate(c.ID)} disabled={updateMutation.isPending}
                      className="btn-rpg-primary px-3 py-1.5 text-xs">Salvar</button>
                    <button onClick={() => setEditingId(null)} className="btn-rpg-outline px-3 py-1.5 text-xs">Cancelar</button>
                  </div>
                </div>
              ) : (
                <div className="flex items-start justify-between gap-3">
                  <div className="flex-1 min-w-0">
                    <p className="text-white font-semibold">{c.name}</p>
                    <span className="text-xs text-gray-500 bg-gray-700/60 px-2 py-0.5 rounded-full inline-block mb-2">D&D {c.edition}</span>
                    <p className="text-gray-400 text-sm whitespace-pre-wrap">{c.main_story || 'Sem história registrada ainda.'}</p>
                  </div>
                  <div className="flex flex-col gap-1.5 flex-shrink-0">
                    <button onClick={() => navigate(`/campaigns/${c.ID}`)} className="btn-rpg-primary px-3 py-1 text-xs">Gerenciar</button>
                    <button onClick={() => startEdit(c)} className="btn-rpg-outline px-3 py-1 text-xs">Editar</button>
                    <button
                      onClick={() => { if (confirm(`Excluir a campanha "${c.name}"? Isso não pode ser desfeito.`)) deleteMutation.mutate(c.ID) }}
                      className="text-xs px-3 py-1 rounded-lg border border-red-800/50 text-red-400 hover:bg-red-900/20 transition"
                    >Excluir</button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
