import { useRef, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { characterService } from '../services/characterService'

// ✅ Usa a variável de ambiente — funciona local E em produção
const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

interface Props {
  characterID: number
  avatarURL?: string
  characterName: string
}

export default function AvatarUpload({ characterID, avatarURL, characterName }: Props) {
  const queryClient = useQueryClient()
  const fileInputRef = useRef<HTMLInputElement>(null)

  // ✅ Monta a URL correta com a base do ambiente
  const buildUrl = (path?: string) => path ? `${API_URL}${path}` : null

  const [preview, setPreview] = useState<string | null>(buildUrl(avatarURL))

  const uploadMutation = useMutation({
    mutationFn: (file: File) => characterService.uploadAvatar(characterID, file),
    onSuccess: (data) => {
      // ✅ Usa a mesma função para montar a URL
      setPreview(buildUrl(data.avatar_url))
      queryClient.invalidateQueries({ queryKey: ['character', String(characterID)] })
      queryClient.invalidateQueries({ queryKey: ['characters'] })
    },
  })

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onloadend = () => setPreview(reader.result as string)
    reader.readAsDataURL(file)
    uploadMutation.mutate(file)
  }

  return (
    <div className="flex flex-col items-center gap-3">
      <div
        onClick={() => fileInputRef.current?.click()}
        className="w-28 h-28 rounded-full overflow-hidden bg-gray-700 border-2 border-gray-600 hover:border-indigo-500 cursor-pointer transition flex items-center justify-center"
      >
        {preview ? (
          <img src={preview} alt={characterName} className="w-full h-full object-cover" />
        ) : (
          <span className="text-4xl">🧙</span>
        )}
      </div>

      <button
        onClick={() => fileInputRef.current?.click()}
        disabled={uploadMutation.isPending}
        className="text-indigo-400 hover:text-indigo-300 text-xs transition"
      >
        {uploadMutation.isPending ? 'Enviando...' : 'Alterar foto'}
      </button>

      <input
        ref={fileInputRef}
        type="file"
        accept=".jpg,.jpeg,.png,.webp"
        onChange={handleFileChange}
        className="hidden"
      />
    </div>
  )
}