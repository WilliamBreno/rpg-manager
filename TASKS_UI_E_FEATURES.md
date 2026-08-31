# RPG Manager — UI, boas-vindas, export de PDF 5e e base de conhecimento

> Gerado no claude.ai a partir de uma sessão de planejamento, e corrigido depois de ler o
> `CLAUDE.md` real do projeto. Leia `CLAUDE.md` primeiro — ele é a fonte da verdade sobre
> arquitetura; este arquivo só descreve as tarefas desta rodada.

## Documentos relacionados (todos na raiz do projeto)

- `TASKS_UI_E_FEATURES.md` (este arquivo) — fases de UI, boas-vindas, export de PDF e tema.
- `DND5E_CRIACAO_PERSONAGEM_CHECKLIST.md` — auditoria de mecânicas na criação de personagem.
- `AUDITORIA_COMPLETUDE_LISTAS_FECHADAS.md` — auditoria de listas fechadas (subclasses, raças,
  antecedentes etc.) incompletas no seed.
- `PROGRESSAO_DE_NIVEL.md` — o que deve desbloquear automaticamente em cada nível, e escopo da
  auditoria mestra de completude 5e.
- `SISTEMA_MESTRE.md` — mesa virtual leve pro mestre (campanhas, NPCs/inimigos/boss/vilão,
  sessões, cenários, sala ao vivo via WebSocket, chat, dados do mestre, recompensas, áudio).

Este arquivo também tem uma seção própria (perto do final) sobre o alerta de "servidor
acordando", que precisa ser revisado à luz da migração do banco pro Neon.

Os quatro se complementam: a Fase 2 deste arquivo (export de PDF) só sai correta se os outros
três estiverem resolvidos primeiro — um personagem exportado é só tão completo quanto os dados
que existem sobre ele.

## Contexto rápido (resumo do CLAUDE.md, não substitui a leitura completa)

- Três apps no monorepo: `backend/` (Go + Gin + GORM, :8080), `frontend/` (React 19 + TS + Vite,
  :5173), `ai-service/` (Python FastAPI, :8000, RAG via ChromaDB + Ollama).
- `Character.Edition` é `"4e"` ou `"5e"` — é o fato arquitetural central, tudo é ramificado por
  esse campo. **A opção de exportar PDF só existe para `Edition == "5e"`.**
- Backend em camadas estritas `handler → service → repository → domain`, sem DI framework,
  fiação manual em `cmd/api/main.go`. Regras 4e vs 5e concentradas em
  `internal/service/character_service.go`.
- `Pericia` (perícias 5e, ex: Furtividade/Percepção) e `Talento` (feats 5e) são modelos
  **diferentes** do `Skill` mais antigo (poderes de classe 4e: at-will/encounter/daily/utility).
  Não confundir os dois ao montar os dados pra ficha PDF.
- Frontend sem feature folders: `pages/`, `components/`, `services/` (um módulo Axios por
  recurso, todos usando `services/api.ts`), `store/authStore.ts` (Zustand,
  persiste token/user em localStorage), `types/index.ts`. Rotas em `App.tsx`, tudo exceto
  `/login` e `/register` dentro de `PrivateRoute`.
- Comentários, mensagens de erro e commits em português — manter consistência.
- Ambiente Windows — comandos de CLI compatíveis com PowerShell/cmd.
- Estilo de trabalho: explicar a abordagem **antes** de escrever código, pra eu validar cada
  passo.

## Localização dos arquivos-fonte

Livros e ficha 5e (o usuário já confirmou que ficam aqui):

```
C:\Users\PC\OneDrive\Área de Trabalho\RPG Manager\ai-service\books
```

Essa é literalmente a pasta que `ai-service/ingest.py` já lê para indexar no `chroma_db/` — não
é uma pasta nova a criar.

Arquivos esperados:
- `dampd-5e---livro-do-jogador-2024.pdf`
- `dd-5e-ficha-de-personagem-completavel-biblioteca-elfica.pdf` ← **template usado na Fase 2**
- `dd-5e-guia-de-xanathar-para-todas-as-coisas-fundo-colorido-biblioteca-elfica.pdf`
- `dd-5e-livro-do-jogador-fundo-branco-biblioteca-elfica.pdf`
- `dd-5e-regras-basicas-dos-jogadores-biblioteca-elfica.pdf`
- `DnD 5ed - Guia do Mestre.pdf`
- `DnD 5ed - Manual de Monstros.pdf`

Antes da Fase 2, mover para `ai-service/pdf_export/reference/` (criar a pasta) os três arquivos
gerados numa sessão anterior do claude.ai:
- `dnd5e_pdf_field_map.json` — mapa categorizado dos 334 campos AcroForm da ficha oficial
- `fill_dnd5e_sheet.py` — função Python testada que preenche a ficha a partir de um dict de
  valores de campo já resolvidos (**não** assume o schema real do projeto — a função
  `character_to_field_values()` nesse arquivo é só um exemplo com nomes de campo inventados;
  precisa ser reescrita usando os nomes reais de `internal/domain` e `Pericia`)
- `ficha_5e_demo_preenchida.pdf` — exemplo de referência já validado visualmente

---

## Fase 0 — Preparação

1. Ler `CLAUDE.md` e este arquivo por completo.
2. Ler também `AUDITORIA_COMPLETUDE_LISTAS_FECHADAS.md` — trata de um bug real já confirmado
   (Escolas de Magia do Mago incompletas) que provavelmente se repete em outras listas fechadas
   do seed (subclasses de outras classes, raças/sub-raças, antecedentes, talentos). Isso afeta
   diretamente a qualidade dos dados que alimentam a Fase 2 (export PDF) e a Fase 1 (IA), então
   vale corrigir cedo.
3. Ler `internal/domain/character.go` (ou onde estiver o struct `Character`) e o model
   `Pericia`/`Talento` antes de planejar a Fase 2 — não assumir nomes de campo.
4. Ler `ai-service/main.py` e `ai-service/ingest.py` por completo antes de planejar as Fases 1
   e 2 — em especial, confirmar se `ai-service` tem qualquer acesso a Postgres/GORM ou se é
   puramente RAG (isso decide a arquitetura da Fase 2, ver nota abaixo).
5. Rodar `go vet ./...` e `go build ./...` no backend como baseline — **não existe nenhum teste
   de backend hoje** (`go test ./...` roda mas não encontra arquivos de teste), então não há
   suíte prévia pra comparar. Os testes da Fase 2 serão, na prática, os primeiros testes de
   backend do projeto.
6. Rodar `npm run lint` e `npm run build` no frontend como baseline.

## Fase 1 — Ingestão dos livros 5e na base de conhecimento

**Objetivo:** os 6 livros de regras 5e (não a ficha) ficam disponíveis para o endpoint
`/skills` do `ai-service`, que auto-gera poderes/perícias de classe.

Passos:
1. Confirmar que os PDFs já estão em `ai-service/books/` (mesma pasta usada hoje para os
   livros de 4e).
2. Ler `ingest.py` pra entender exatamente como ele processa cada PDF hoje (OCR com Tesseract +
   PyMuPDF, chunking com `langchain-text-splitters`) antes de rodar — em especial, checar se ele
   já grava metadata por livro/edição no ChromaDB, ou se todos os livros caem numa collection
   única sem diferenciação. Se não houver diferenciação por edição, avaliar (comigo) se vale
   adicionar `{"livro": "<nome>", "edicao": "4e"|"5e"}` no metadata de cada chunk antes de
   indexar os novos livros — isso importa pro upgrade de IA já planejado (Claude API Haiku +
   cache por classe/edição/nível), que depende de conseguir filtrar por edição.
3. Rodar `python ingest.py` para indexar os novos PDFs.
4. Testar uma query real em `/skills` para uma classe 5e (ex: Mago) e confirmar que o resultado
   vem de um dos livros novos, não é alucinado.

**Nota:** o bug conhecido "perícias do livro Poder Arcano não aparecem" é sobre conteúdo de
**4e** — pipeline separado, esta fase não deve mexer na ingestão de 4e.

**Critério de aceite:** uma query em `/skills` para uma classe/nível 5e retorna algo
rastreável a um dos livros recém-indexados.

## Fase 2 — Exportar personagem 5e para a ficha PDF preenchível

**Antes de qualquer coisa nesta fase:** rodar a auditoria descrita em
`DND5E_CRIACAO_PERSONAGEM_CHECKLIST.md` (arquivo separado, gerado numa sessão de planejamento
específica sobre isso). Ela existe porque exportar a ficha só sai correta se o personagem tiver,
em algum lugar do sistema, todo o dado que a ficha pede — e hoje não está claro quais mecânicas
de criação de personagem 5e já existem de verdade e quais são só lacunas. **Preencher a tabela
de status daquele arquivo primeiro, me mostrar o resultado, e só então seguir pros passos abaixo.**

**Objetivo:** um jogador com personagem `Edition == "5e"` pode baixar a ficha oficial
preenchida com os dados do personagem. O botão só aparece para personagens 5e.

**Decisão de arquitetura a confirmar comigo antes de implementar** (depende do que a Fase 0.3
descobrir sobre o `ai-service`):
- Se `ai-service` **não** tem acesso a Postgres/GORM (provável, já que ele é descrito só como
  RAG/ChromaDB/Ollama): o backend Go monta o dicionário de valores já resolvidos (usando
  `internal/service`, respeitando `Pericia`/`Talento`/o bug de PV) e faz um `POST` HTTP interno
  para um novo endpoint no `ai-service` (ex: `POST /export/pdf/5e`, recebendo o JSON de valores
  + `character_id` só pra nome do arquivo), que só roda o preenchimento do PDF com `pypdf` e
  devolve os bytes. O Go backend repassa esses bytes pro frontend. Vantagem: reaproveita
  `fill_dnd5e_sheet.py` quase direto, sem reescrever preenchimento de PDF em Go.
- Se `ai-service` já acessa o banco por algum outro motivo: reavaliar, pode fazer mais sentido
  ele buscar o personagem direto.
- Alternativa (evitar chamada entre serviços): usar uma lib Go de preenchimento de PDF (ex:
  `pdfcpu`) e fazer tudo dentro do backend Go, sem depender do `ai-service`. Mais simples
  operacionalmente, mas exige reescrever a lógica de `fill_dnd5e_sheet.py` em Go.

Passos sugeridos (assumindo a primeira opção):
1. Em `internal/service`, nova função que monta o dict de campos a partir de um `Character` 5e
   real — atributos, `ProficiencyBonus`, PV atual/máximo (**respeitar o bug #2: nunca
   sobrescrever o PV digitado pelo jogador**), `Pericia`s do personagem, `Talento`s, etc. Usar
   `dnd5e_pdf_field_map.json` como referência dos nomes de campo válidos do PDF.
2. Novo endpoint no `ai-service` (seguir o padrão de estrutura que `main.py` já usa para
   `/skills`) que recebe esse JSON e devolve o PDF preenchido, adaptando
   `fill_dnd5e_sheet.py`.
3. Novo handler Go (`internal/handler`) que valida `Edition == "5e"`, valida que o personagem
   pertence ao usuário autenticado (`AuthMiddleware`/`userID` do contexto Gin), monta os dados,
   chama o `ai-service`, e repassa o PDF (`Content-Type: application/pdf`). Registrar a rota em
   `main.go`.
4. No frontend, botão "Exportar ficha PDF" em `pages/` (tela do personagem), visível só quando
   `character.edition === "5e"`, chamando via um novo módulo em `services/`.
5. **Antes de produção:** as ~68 caixas de seleção do PDF (proficiências) têm IDs sem nome
   semântico (ex: `Check Box 23`). Rodar `convert_pdf_to_images.py` sobre a ficha em
   `ai-service/books` e casar visualmente cada ID com a perícia/salvaguarda correspondente
   antes de preenchê-las — não adivinhar a partir da ordem dos campos.
6. Escrever o primeiro teste de backend do projeto para este fluxo (gerar personagem 5e de
   teste, montar os valores, verificar que os campos essenciais batem) e um teste do lado
   Python que verifica que o PDF retornado tem os valores esperados (pode usar `pypdf` pra ler
   de volta e comparar).

**Critério de aceite:** baixar a ficha de um personagem 5e real do sistema mostra os dados
corretos nos campos certos ao abrir no leitor de PDF.

## Fase 3 — Tela de boas-vindas (Fase A do roadmap)

**Objetivo:** modal único, exibido apenas no primeiro registro do jogador, antes da tela de
personagens.

- Frase: *"O dado foi lançado. Sua jornada de uma lenda começa aqui."*
- Direção visual (protótipo já aprovado numa sessão anterior): fundo escuro, dado com animação
  de entrada (rotação + fade), frase principal em fonte mais editorial, botão levando pra
  criação de personagem. O protótipo era só referência de direção — na implementação real, use
  Motion (via Unlumen UI, Fase 5) livremente pra ir além do que o protótipo mostrou.
- Decidir onde guardar a flag "já viu o welcome": local mais confiável é um campo no `User`
  (`internal/domain`) e não só `localStorage` via `authStore` — `localStorage` some ao trocar de
  navegador/dispositivo e mostraria o modal de novo. Confirmar comigo antes de migrar o schema
  do `User`.
- Fluxo provável: `App.tsx`/`PrivateRoute` checa a flag depois do login e redireciona pro modal
  antes da tela de personagens.

**Critério de aceite:** criar conta nova mostra o modal uma única vez; logins seguintes (mesmo
em outro navegador) não mostram mais.

## Fase 4 — Tela "Meus Personagens" (dashboard de personagens)

**Objetivo:** redesenhar a tela que lista os personagens do jogador, usando o mockup
"ArcanaRPG" (imagem anexada pelo usuário numa sessão de planejamento no claude.ai) como
referência visual — não pixel-perfect, é direção de estilo, não um design final fechado.

Elementos do mockup a reproduzir:
- Grid de cards de personagem, cada um com: retrato do personagem, nome, raça • classe, nível,
  e uma barra de progresso de XP (atual/necessário pro próximo nível) com uma cor de destaque
  própria por card. Menu de três pontos no canto do card para ações (usar as ações que já
  existem hoje no fluxo do projeto — editar, excluir etc.).
- Cabeçalho da página: título "Meus Personagens", subtítulo curto, botão de destaque "Criar
  Personagem" levando ao fluxo de criação já existente.
- Painel lateral com resumo em números e uma lista de atividade recente.

**Pedido específico do usuário: moldura temática de RPG em cada card.** Não é uma borda de
card genérica — pensar em algo como ornamentos nos cantos, textura de couro/pergaminho/runa, ou
contorno com brilho sutil ao passar o mouse (hover). Verificar primeiro se Unlumen UI ou algum
snippet do Uiverse (Fase 5) já cobre algo parecido antes de construir do zero.

**Atenção ao painel de resumo — não inventar dado que não existe:** "Personagens" é real (dá
pra contar). "Campanhas" e "Conquistas" não aparecem em nenhum model do `internal/domain`
descrito no `CLAUDE.md` hoje — antes de replicar esses blocos, checar se essas entidades
existem no banco. Se não existirem, não simular números falsos; ou omitir o bloco por enquanto,
ou confirmar comigo se vale criar um placeholder "em breve" no lugar. O mesmo vale pra
"Atividade Recente": só popular com eventos que o sistema realmente registra hoje (ex: subida
de nível), não inventar tipos de evento que não são rastreados em lugar nenhum do backend.

**Menu lateral de navegação (Início, Campanhas, Inventário, Biblioteca, Mapas,
Configurações): NÃO implementar essas páginas ainda — é trabalho futuro, fora do escopo desta
rodada.** Pode reproduzir a estrutura visual do menu, mas só deixando clicáveis os itens que já
têm uma página real e funcional hoje. Antes de decidir como tratar os itens que ainda não
existem (escondê-los vs. mostrar desabilitado/"em breve"), parar e confirmar comigo — mostrar
um item de menu morto pro jogador clicar é pior do que não mostrar.

**Critério de aceite:** a tela de personagens mostra os personagens reais do jogador logado em
cards com a moldura temática, sem nenhum número ou evento inventado nos painéis auxiliares, e
sem links para páginas que não existem.

## Fase 5 — Bibliotecas de UI e tema visual geral

| Biblioteca | Uso |
|---|---|
| **shadcn/ui** | Base: forms, dialogs, cards, inputs — prioridade #1 do projeto |
| **Unlumen UI** (tem MCP próprio) | Botões com glow, cards com tilt, scramble/shimmer text — boas-vindas, level up |
| **Magic MCP (21st.dev)** | Celebrações maiores: subir de nível, crítico natural 20 |
| **ReactBits** | Só fundo decorativo/ambiente (partículas, texto animado) — com moderação |
| **Uiverse** | Snippets avulsos de CSS/HTML copiados manualmente (não é lib instalável) — adaptar cores/fontes aos tokens do projeto |

Se os MCPs (Shadcn MCP, Magic MCP, Unlumen UI MCP) ainda não estiverem configurados neste
ambiente de Claude Code, avisar antes de tentar usá-los.

Onde não houver biblioteca específica cobrindo o componente: visual com "cara de sistema de
RPG", efeitos e animações onde fizer sentido (transição ao subir de nível, hover em cards de
personagem, feedback ao rolar dados), sem prejudicar legibilidade/performance.

## Fase 6 — Testes finais

- `go vet ./...`, `go build ./...` e `go test ./...` no backend (incluindo os testes novos da
  Fase 2 — os primeiros do projeto).
- `npm run lint` e `npm run build` no frontend.
- Testes do `ai-service` para o novo endpoint de export e para a ingestão da Fase 1.
- Checklist manual: criar conta nova (vê welcome uma vez), criar personagem 5e, exportar PDF e
  conferir visualmente, fazer uma query em `/skills` pra uma classe 5e recém-indexada.
- Reportar qualquer teste que falhar antes de considerar a fase concluída — não marcar como
  "100% funcionando" sem ter rodado de fato.

---

## Verificação extra: alerta de "servidor acordando" (Render → Neon)

Contexto: o toast "servidor acordando" em `services/api.ts` foi criado quando o cold start
relevante era o do backend Go no plano gratuito do Render (arrancar o container inteiro do
zero — tipicamente dezenas de segundos). Desde então, o banco migrou pro Neon, que tem seu
próprio comportamento de "scale to zero": a instância de computação do Postgres suspende após
um período de inatividade e "acorda" na próxima query. **São dois cold starts possivelmente
diferentes, não o mesmo fenômeno com nome trocado** — antes de decidir o que fazer com o toast,
confirmar:

1. O backend Go ainda está hospedado no plano gratuito do Render (com spin-down completo do
   container), ou migrou pra outro lugar? Se ainda estiver no Render free tier, esse cold start
   — o mais lento dos dois — continua sendo o principal suspeito.
2. Qual é o tempo de "autosuspend" configurado no projeto Neon (painel do Neon → Settings →
   Compute)? E qual a latência real observada numa query feita depois desse tempo de
   inatividade — **medir de verdade, não assumir**. O resume do Neon costuma ser bem mais
   rápido que um cold start de container inteiro (frequentemente abaixo de 1-2 segundos nas
   versões recentes), mas isso varia por plano/região.
3. Com os dois números reais em mãos, escolher entre estas opções — meu palpite é começar pela
   (A) e só cair pras outras se ela não for suficiente:
   - **(A) Aquecimento silencioso em segundo plano:** disparar uma requisição leve (ex: um
     health-check) assim que o app carrega, antes do jogador precisar dos dados de verdade —
     enquanto ele ainda está na tela de login/splash. Se terminar antes de ele navegar pra uma
     tela com dados, ele nunca vê aviso nenhum. Manter o toast atual só como *fallback*, caso
     mesmo assim uma requisição real ultrapasse o limiar de tempo.
   - **(B) Trocar o toast por um estado de carregamento comum** (skeleton/spinner na própria
     tela, sem o texto de "servidor dormindo") — faz sentido se a medição do item 2 mostrar que
     o delay do Neon é curto o bastante pra não assustar ninguém, e se o Render não for mais o
     gargalo principal.
   - **(C) Ping periódico externo pra manter tudo acordado** (cron/GitHub Action batendo num
     health-check a cada poucos minutos) — funciona, mas tem custo real: em planos gratuitos
     (Render e Neon) isso consome a cota de horas gratuitas mais rápido, então não é "grátis" —
     é trocar UX por cota de uso. Só faz sentido se o projeto puder absorver isso ou já estiver
     num plano pago.
   - Se algum texto de aviso for mantido, atualizar a mensagem pra refletir a causa real (banco
     de dados, especificamente, não "servidor" genérico) — hoje o texto pode estar descrevendo
     uma causa que já não é exatamente a que está acontecendo.

**Critério de aceite:** o jogador não vê um aviso alarmante de "servidor dormindo" pra um atraso
de 1-2 segundos; se o atraso real ainda for de dezenas de segundos (Render), ele continua sendo
avisado, só que com um texto que descreve a causa certa.

## Atualizações que este trabalho vai deixar pendentes no CLAUDE.md

O `CLAUDE.md` já tem uma regra de manutenção própria (atualizar a seção correspondente sempre
que algo mudar, sem changelog separado). Pra não depender de lembrar isso de forma genérica,
aqui estão os pontos específicos que cada fase vai deixar desatualizados — resolver junto com o
critério de aceite de cada fase, não deixar acumular pro final:

- **Fase 1:** seção "AI service" — mencionar que `ai-service/books/` agora inclui livros 5e e,
  se `ingest.py` passou a gravar metadata de edição por chunk, documentar isso.
- **Fase 2:** seção "Architecture → Backend layering" e "AI service" — descrever o novo endpoint
  de export, onde ele mora (ai-service ou backend, conforme a decisão tomada na Fase 2) e a
  chamada entre serviços envolvida. Também remover o "(no test files currently exist)" do
  comando `go test ./...` assim que os primeiros testes existirem.
- **Fase 3:** se a flag "já viu o welcome" virar um campo novo no `User` (domain), atualizar a
  lista de campos do model `User` na seção Architecture.
- **Fase 4:** se o painel de resumo revelar que "Campanhas"/"Conquistas" não existem no domain
  hoje, e a decisão for criar esses models, documentar isso na seção Architecture. Também vale
  registrar en passant que a tela de personagens ganhou uma versão nova (útil pra próximas
  sessões que forem mexer nela).
- **Fase 5:** vale acrescentar uma linha na seção "Frontend structure" citando as bibliotecas de
  UI adotadas (shadcn/ui, Unlumen UI, etc.), pra próximas sessões saberem que existem sem
  precisar redescobrir olhando o `package.json`.
- **Verificação do alerta de servidor acordando:** a frase atual do CLAUDE.md ("shows a 'waking
  up server' toast if a request takes >4s, the backend is hosted on a free tier with cold
  starts") precisa refletir a decisão tomada — se a causa raiz virou o Neon, se o texto/UX
  mudou, ou se o comportamento no fundo continua o mesmo e só a explicação estava desatualizada.

---

## Bugs conhecidos (fora do escopo desta rodada, mas não piorar)

1. Perícias do livro "Poder Arcano" (4e) não aparecem em queries de IA — pipeline separado do
   da Fase 1.
2. PV informado na criação é sobrescrito pelo valor calculado em `character_service.go` —
   relevante pra Fase 2, não repetir esse erro ao montar os dados do export.
3. Classe de Armadura incorreta para Bardo em 4e usando Inteligência como maior modificador de
   defesa com armadura de couro.
