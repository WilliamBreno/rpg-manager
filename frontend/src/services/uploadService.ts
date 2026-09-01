import api from './api'

export type UploadKind = 'image' | 'audio'

export const uploadService = {
  // Sobe um arquivo genérico (imagem ou áudio) e devolve o data URI base64
  // pronto pra usar em qualquer campo que hoje aceita uma URL (foto/som de
  // inimigo, imagem de cenário, áudio de fala, música de sessão) — não
  // persiste nada por si só, quem chama decide onde colocar o resultado.
  uploadFile: async (file: File, kind: UploadKind): Promise<string> => {
    const formData = new FormData()
    formData.append('file', file)
    const { data } = await api.post(`/uploads/file?kind=${kind}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return data.data_url as string
  },
}
