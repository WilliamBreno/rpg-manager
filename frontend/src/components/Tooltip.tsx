import { useState } from 'react'

interface TooltipProps {
  content: string
}

export function Tooltip({ content }: TooltipProps) {
  const [open, setOpen] = useState(false)

  return (
    <div className="relative flex-shrink-0">
      <button
        type="button"
        className="w-5 h-5 rounded-full bg-gray-600 hover:bg-indigo-600 text-gray-300 text-xs font-bold flex items-center justify-center transition"
        onMouseEnter={() => setOpen(true)}
        onMouseLeave={() => setOpen(false)}
        onClick={e => { e.stopPropagation(); setOpen(v => !v) }}
        aria-label="Ver detalhes"
      >
        ?
      </button>

      {open && (
        <div className="absolute right-0 bottom-7 w-72 bg-gray-950 text-gray-200 text-xs rounded-lg px-3 py-2.5 border border-gray-600 shadow-2xl z-50 leading-relaxed">
          {content}
          {/* seta apontando para baixo */}
          <div className="absolute top-full right-2 border-4 border-transparent border-t-gray-600" />
        </div>
      )}
    </div>
  )
}