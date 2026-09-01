import { useRef, useState } from 'react'
import { uploadService, type UploadKind } from '../services/uploadService'

const ACCEPT: Record<UploadKind, string> = {
  image: 'image/jpeg,image/png,image/webp,image/gif',
  audio: 'audio/mpeg,audio/wav,audio/ogg,audio/mp4,audio/webm,.mp3,.wav,.ogg,.m4a,.webm',
}

// Campo de upload de arquivo reaproveitável — substitui os antigos inputs de
// URL simples pra foto/som de inimigo, imagem de cenário, áudio de fala e
// música de sessão (Sistema do Mestre). Sobe o arquivo pra
// POST /uploads/file, que devolve um data URI já pronto pra ir no payload
// normal do recurso — não fala com nenhum outro endpoint.
export default function FileUploadField({
  label,
  kind,
  value,
  onChange,
  hint,
}: {
  label: string
  kind: UploadKind
  value: string
  onChange: (dataUrl: string) => void
  hint?: string
}) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState('')

  const handleFile = async (file: File) => {
    setError('')
    setUploading(true)
    try {
      const dataUrl = await uploadService.uploadFile(file, kind)
      onChange(dataUrl)
    } catch (err) {
      const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
        ?? 'Falha ao enviar arquivo.'
      setError(message)
    } finally {
      setUploading(false)
      if (inputRef.current) inputRef.current.value = ''
    }
  }

  return (
    <div>
      <label className="text-gray-500 text-xs mb-1 block uppercase tracking-wider">{label}</label>

      <input
        ref={inputRef}
        type="file"
        accept={ACCEPT[kind]}
        className="hidden"
        onChange={e => {
          const file = e.target.files?.[0]
          if (file) handleFile(file)
        }}
      />

      {!value && (
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          disabled={uploading}
          className="btn-rpg-outline w-full py-2 text-sm disabled:opacity-50"
        >
          {uploading ? 'Enviando...' : kind === 'image' ? '📁 Escolher imagem' : '📁 Escolher áudio'}
        </button>
      )}

      {value && kind === 'image' && (
        <div className="flex items-center gap-3 bg-gray-800/40 rounded-lg p-2">
          <img src={value} alt="" className="w-14 h-14 rounded object-cover flex-shrink-0 border border-gray-700" />
          <div className="flex-1 min-w-0 flex gap-2">
            <button type="button" onClick={() => inputRef.current?.click()} disabled={uploading}
              className="btn-rpg-outline px-3 py-1 text-xs">{uploading ? 'Enviando...' : 'Trocar'}</button>
            <button type="button" onClick={() => onChange('')} className="text-xs text-red-400 hover:underline">Remover</button>
          </div>
        </div>
      )}

      {value && kind === 'audio' && (
        <div className="flex flex-col gap-2 bg-gray-800/40 rounded-lg p-2">
          <audio src={value} controls className="w-full h-9" />
          <div className="flex gap-2">
            <button type="button" onClick={() => inputRef.current?.click()} disabled={uploading}
              className="btn-rpg-outline px-3 py-1 text-xs">{uploading ? 'Enviando...' : 'Trocar'}</button>
            <button type="button" onClick={() => onChange('')} className="text-xs text-red-400 hover:underline">Remover</button>
          </div>
        </div>
      )}

      {error && <p className="text-red-400 text-xs mt-1">{error}</p>}
      {hint && !error && <p className="text-gray-600 text-xs mt-1">{hint}</p>}
    </div>
  )
}
