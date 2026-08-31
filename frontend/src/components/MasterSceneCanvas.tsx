import { useEffect, useState } from 'react'
import { Stage, Layer, Image as KonvaImage, Text, Rect, Group } from 'react-konva'
import useImage from 'use-image'
import type { Scene, SceneToken } from '../types'

const CANVAS_WIDTH = 800
const CANVAS_HEIGHT = 500

function TokenNode({ token, onMove, onDelete, readOnly }: { token: SceneToken; onMove: (x: number, y: number) => void; onDelete: () => void; readOnly?: boolean }) {
  const [image] = useImage(token.image_url || '')
  const size = 48

  return (
    <Group
      x={token.x}
      y={token.y}
      draggable={!readOnly}
      onDragEnd={e => onMove(e.target.x(), e.target.y())}
      onDblClick={readOnly ? undefined : onDelete}
      onDblTap={readOnly ? undefined : onDelete}
    >
      {image ? (
        <KonvaImage image={image} width={size} height={size} cornerRadius={size / 2} />
      ) : (
        <Rect width={size} height={size} fill="#c9a84c" cornerRadius={size / 2} stroke="#1a1a1a" strokeWidth={2} />
      )}
      <Text text={token.label} y={size + 2} width={size * 2} x={-size / 2} align="center" fontSize={12} fill="#e5e5e5" />
    </Group>
  )
}

// Renderiza o cenário ativo (imagem de fundo) com tokens arrastáveis por
// cima — Sistema do Mestre, Etapa 4. Só posição/arrastar, sem lógica de
// movimento por turno (limite de escopo já combinado). Dê um duplo-clique
// num token pra removê-lo.
export default function MasterSceneCanvas({
  scene,
  onMoveToken,
  onDeleteToken,
  readOnly,
}: {
  scene: Scene
  onMoveToken: (tokenId: number, x: number, y: number) => void
  onDeleteToken: (tokenId: number) => void
  readOnly?: boolean
}) {
  const [bgImage] = useImage(scene.image_url || '')
  const [tokens, setTokens] = useState<SceneToken[]>(scene.tokens ?? [])

  useEffect(() => { setTokens(scene.tokens ?? []) }, [scene.tokens, scene.ID])

  const handleMove = (tokenId: number, x: number, y: number) => {
    setTokens(prev => prev.map(t => (t.ID === tokenId ? { ...t, x, y } : t)))
    onMoveToken(tokenId, x, y)
  }

  return (
    <div className="rounded-xl overflow-hidden border border-gray-800 inline-block bg-gray-950">
      <Stage width={CANVAS_WIDTH} height={CANVAS_HEIGHT}>
        <Layer>
          {bgImage ? (
            <KonvaImage image={bgImage} width={CANVAS_WIDTH} height={CANVAS_HEIGHT} />
          ) : (
            <Rect width={CANVAS_WIDTH} height={CANVAS_HEIGHT} fill="#111827" />
          )}
          {!scene.image_url && (
            <Text text="Sem imagem de fundo" x={0} y={CANVAS_HEIGHT / 2 - 8} width={CANVAS_WIDTH} align="center" fill="#6b7280" fontSize={14} />
          )}
        </Layer>
        <Layer>
          {tokens.map(t => (
            <TokenNode key={t.ID} token={t} readOnly={readOnly} onMove={(x, y) => handleMove(t.ID, x, y)} onDelete={() => onDeleteToken(t.ID)} />
          ))}
        </Layer>
      </Stage>
    </div>
  )
}
