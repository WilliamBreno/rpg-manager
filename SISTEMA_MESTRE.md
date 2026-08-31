# Sistema do Mestre — mesa virtual leve integrada ao RPG Manager

> Complementa `TASKS_UI_E_FEATURES.md`, `DND5E_CRIACAO_PERSONAGEM_CHECKLIST.md`,
> `AUDITORIA_COMPLETUDE_LISTAS_FECHADAS.md` e `PROGRESSAO_DE_NIVEL.md`. Este documento nomeia
> algo importante: o que está descrito aqui é, na prática, uma **mesa virtual (VTT) leve**,
> integrada ao sistema de personagens que já existe — no mesmo território de Roll20/Foundry
> VTT/Owlbear Rodeo, só que sob medida e mais enxuta. É um projeto grande. Vai ser dividido em
> etapas, e cada etapa deve ser reportada antes de seguir pra próxima — não tentar tudo de uma
> vez.

## Decisões de arquitetura já tomadas (confirmadas com o usuário)

1. **Sincronização em tempo real via WebSockets desde já** — quando o mestre troca de cenário
   ou envia uma mensagem no chat, os jogadores veem instantaneamente, sem precisar atualizar a
   página.
2. **Falas de boss/vilão: as duas opções disponíveis** — o mestre pode gravar/fazer upload do
   próprio áudio, ou gerar via texto-pra-voz (TTS). Ambas devem existir desde o início, não só
   uma.
3. **Mapa interativo: Konva.js / react-konva** (MIT, grátis, sem tier pago, binding oficial pra
   React 19) — biblioteca de canvas pra desenhar o mapa e tornar tokens arrastáveis. É a única
   biblioteca nova necessária além das 5 já adotadas (shadcn/ui, Unlumen UI, Magic UI, ReactBits,
   Uiverse), porque nenhuma delas cobre canvas interativo — são bibliotecas de UI, não de
   desenho/drag-and-drop de objetos num mapa. Prova de conceito de referência (não é biblioteca
   pra importar, é inspiração de UX): o Owlbear Rodeo é elogiado justamente por deixar trocar de
   cenário (mapa aberto → masmorra e vice-versa) de forma muito rápida e sem fricção — é
   exatamente esse padrão de "trocar de cenário de forma prática" que vale reproduzir aqui.

## Decisões que ainda precisam de confirmação antes de implementar (não decidir sozinho)

- **Armazenamento de mídia** (imagens de inimigo/NPC/cenário, áudio de música/som/falas): hoje
  não existe nenhum serviço de storage no projeto (nem Render, nem Vercel, nem Neon servem esse
  propósito — Neon é só banco relacional). Alguém precisa escolher um provedor (ex: Cloudinary
  ou Supabase Storage, ambos com tier gratuito generoso e API simples para imagem+áudio). Propor
  opções pro usuário antes de integrar qualquer um.
- **Provedor de TTS**: opções variam muito em custo e qualidade — desde a Web Speech API do
  navegador (grátis, qualidade robótica, zero configuração) até serviços pagos de voz mais
  realista. Como isso tem custo recorrente dependendo da escolha, propor opções com o
  usuário antes de integrar, não escolher um serviço pago sozinho.
- **Servidor WebSocket no Go**: Gin não tem suporte nativo a WebSocket — vai precisar de uma
  lib como `gorilla/websocket` ou `nhooyr.io/websocket`. Antes de implementar, desenhar o
  padrão de "hub/sala" (uma sala por sessão de campanha ativa, broadcast de eventos: troca de
  cenário, mensagem de chat, e futuramente movimento de token) e mostrar esse desenho pro
  usuário antes de codar — é a peça de infraestrutura mais nova de todo o projeto até agora.

## Escopo confirmado como fora desta rodada (o próprio usuário já adiou)

- **Token 2D jogável com movimento por turno e sistema de ações**: o usuário já sinalizou que
  isso será organizado no futuro. Nesta rodada, só preparar o terreno pra não bloquear depois —
  por exemplo, já incluir campos de posição (x/y no grid) no modelo de token, mesmo sem
  implementar a lógica de turno/ação ainda. Não construir o motor de combate agora.
- **Importar inimigos/NPCs pro mapa como tokens quando os jogadores entram na sessão**: também
  adiado pelo próprio usuário pra depois que a base (campanha, sessão, sala, cenário) estiver
  pronta. Só desenhar o modelo de dados de forma que isso encaixe depois sem retrabalho (ex: um
  inimigo/NPC já ter os campos de imagem/HP que um token vai precisar exibir).

---

## Etapa 0 — Modelagem de dados

Antes de qualquer UI, desenhar os novos modelos em `internal/domain` (seguindo o padrão já
existente de GORM models) e me mostrar o desenho antes de migrar o banco:

- `Campaign` — nome, edição (`5e` por enquanto, campo já pronto pra aceitar `4e` no futuro),
  história principal, mestre (dono), lista de jogadores membros.
- `CampaignInvite`/`CampaignMembership` — relação mestre↔jogador por campanha (quem pode
  entrar na sala de quem).
- `NPC` — vida, nome, história, vínculos, tendência (alinhamento), personalidade, observações/
  notas. Pertence a uma `Campaign`.
- `Enemy` — pontos de vida, raça, foto (URL de mídia), som que o inimigo faz (URL de áudio),
  classe (opcional), pontos de armadura (CA), lista de habilidades customizadas (nome + dano).
  Pertence a uma `Campaign`.
- `Boss` — tudo que `Enemy` tem (pode ser composição/embed do mesmo struct, ou uma flag
  `IsBoss` no próprio `Enemy` — decidir o padrão junto com quem for implementar, mas sem
  duplicar campos à toa) + nome de destaque + lista de falas (texto + referência ao áudio,
  gerado por TTS ou upload).
- `Villain` — tudo que `Enemy` tem + nome + história + vínculos + observações + lista de falas
  (mesmo esquema do Boss).
- `MagicItem`/`Object` — nome, descrição, efeito, pertence a uma `Campaign`, pode ser dado como
  recompensa.
- `Session` — pertence a uma `Campaign`, tem início, fim, e um registro/resumo do que aconteceu
  (texto, editável pelo mestre — o "diário de bordo" da sessão).
- `Scene`/`Battleground` — pertence a uma `Campaign`, imagem de fundo do mapa, nome do cenário.
  Uma campanha tem uma biblioteca de cenários; uma `Session` referencia qual cenário está ativo
  no momento.
- `ChatMessage` — remetente, campanha/sessão, texto, timestamp. Serve tanto pro chat
  mestre↔jogadores quanto jogador↔jogador dentro da mesma campanha.
- `Reward`/`Transaction` — registro de moeda ou item dado pelo mestre a um jogador/personagem
  (não só alterar o saldo direto — manter um histórico de "o que foi dado, quando, por quem").

**Habilidade de inimigo com dano real de D&D 5e:** ao criar uma habilidade customizada de
inimigo/boss/vilão, o campo de dano não deveria ser um número solto — a UI deveria orientar o
mestre a usar notação de dado real (`XdY+Z`, ex: `2d6+3`), e idealmente sugerir faixas de dano
coerentes com o nível/CR do grupo (isso pode se apoiar no que já foi levantado em
`PROGRESSAO_DE_NIVEL.md` sobre progressão) em vez de deixar o mestre digitar qualquer número.
Não impedir o mestre de fugir da faixa sugerida (é uma sugestão, não uma regra travada), mas dar
o ponto de partida certo pra ele não ter que ir calcular isso no livro.

## Etapa 1 — Campanha e história principal

CRUD de campanha: criar, editar a história principal, listar campanhas do mestre. Edição fixa
em `5e` por enquanto (campo já existe no modelo, só não expor a opção `4e` na UI ainda).

## Etapa 2 — Criação de NPCs, Inimigos, Boss e Vilão

Um formulário por tipo, exatamente com os campos listados na Etapa 0. Pontos de atenção:
- Boss e Vilão reaproveitam o formulário de Inimigo por trás (evitar duplicar a UI do zero).
- O upload de foto/som já depende da decisão de storage pendente acima — não implementar upload
  de arquivo antes dessa decisão estar tomada.
- As falas de Boss/Vilão: campo de texto + botão "gerar áudio (TTS)" e/ou "enviar áudio
  gravado" — as duas vias devem coexistir na mesma tela, sem uma substituir a outra.

## Etapa 3 — Sessões

Criar sessão (associada a uma campanha), encerrar sessão, e durante/depois dela permitir que o
mestre registre o que aconteceu (texto livre, tipo um diário). Listar histórico de sessões
passadas de uma campanha, cada uma com seu registro.

## Etapa 4 — Cenários / Battlegrounds

Upload de imagem de mapa, nome do cenário, biblioteca de cenários por campanha (pra já ter todos
prontos antes da sessão começar, como o usuário pediu). Tela de troca rápida de cenário ativo
durante uma sessão — essa é a peça de UX inspirada no Owlbear Rodeo citado acima. Usar
Konva/react-konva pra renderizar o cenário ativo com suporte a tokens arrastáveis (mesmo que,
por enquanto, os únicos "tokens" sejam decorativos/manuais, já que importar inimigo/NPC como
token automaticamente foi adiado).

## Etapa 5 — Sala ao vivo (WebSockets)

O mestre "abre a sala" de uma campanha. Jogadores com convite aceito (`CampaignMembership`)
podem entrar. Ao entrar, o jogador é levado pro cenário que o mestre definiu como ativo. Se o
mestre trocar de cenário durante a sessão, todos os jogadores conectados veem a troca em tempo
real, via WebSocket — não precisam atualizar a página. Desenhar o padrão de sala/hub no Go antes
de implementar (um hub por campanha com sessão ativa; eventos: `scene_changed`,
`chat_message`, e mais eventos que forem necessários depois).

## Etapa 6 — Chat

Um `ChatMessage` por campanha, visível tanto na área do mestre quanto na área do jogador,
entregue em tempo real pelo mesmo canal WebSocket da Etapa 5. Do lado do jogador, precisa
existir uma área pra "adicionar mestre(s)" (aceitar convite de campanha) separada de uma área
pra "adicionar jogadores" (não é o jogador quem convida outros jogadores pra campanha do mestre
— checar com o usuário se essa segunda área é sobre contatos/amigos em geral ou sobre visualizar
quem mais está na mesma campanha, porque o pedido original é ambíguo nesse ponto específico e
vale confirmar antes de implementar errado).

## Etapa 7 — Dados do mestre

Um "roller" de dados completo (todos os dados de D&D: d4, d6, d8, d10, d12, d20, d100) do lado
do mestre, com o mesmo tema visual/biblioteca usada pro dado do jogador (Unlumen UI/Magic UI,
conforme já decidido na Fase 5 do `TASKS_UI_E_FEATURES.md`) — não criar um componente de dado
paralelo com visual diferente.

## Etapa 8 — Recompensas

O mestre pode dar moeda (PO/PP/PE/PC/PL, mesmo esquema já usado na ficha 5e) ou um `MagicItem`/
item de campanha a um jogador específico ou a todos. Cada entrega gera um `Reward`/`Transaction`
no histórico, e atualiza o saldo/inventário do personagem do jogador.

## Etapa 9 — Áudio de sessão

Música/som de fundo por sessão (upload, tocando em loop pros jogadores conectados), som que
cada inimigo faz (tocado quando o mestre aciona, ex: ao trazer o inimigo à cena), e as falas de
Boss/Vilão da Etapa 2 (tocadas quando o mestre aciona, pros jogadores ouvirem). Tudo depende da
decisão de storage pendente.

---

## Ordem recomendada geral

Etapa 0 (modelagem) → Etapa 1 (campanha) → Etapa 2 (NPCs/inimigos/boss/vilão) → Etapa 3
(sessões) → Etapa 4 (cenários) → Etapa 5 (sala ao vivo/WebSocket) → Etapa 6 (chat) → Etapa 7
(dados do mestre) → Etapa 8 (recompensas) → Etapa 9 (áudio). Essa ordem existe porque cada
etapa depende de dado que a anterior cria — não pular pra frente sem a anterior estar de pé.

## Critério de aceite geral

Um mestre consegue: criar uma campanha 5e com história principal, criar um NPC/inimigo/boss/
vilão com todos os campos pedidos, montar uma biblioteca de cenários, abrir uma sessão, ver um
jogador convidado entrar na sala e cair no cenário certo, trocar de cenário e ver isso refletir
pro jogador sem ele atualizar a página, conversar por chat com o jogador em tempo real, rolar um
dado com o mesmo visual do jogador, dar uma recompensa em moeda pro jogador, tocar uma música de
fundo e a fala gravada/gerada de um vilão — e encerrar a sessão com um registro do que aconteceu
salvo pra consulta depois.
