# Progressão de nível automática e auditoria mestra de completude (5e)

> Complementa `DND5E_CRIACAO_PERSONAGEM_CHECKLIST.md` e
> `AUDITORIA_COMPLETUDE_LISTAS_FECHADAS.md`. Aqueles cobrem "o personagem nasce completo"; este
> arquivo cobre "o personagem evolui completo, sem o jogador precisar consultar o livro pra
> saber o que ganhou".

## Antes de tudo: definindo o que "completo" significa aqui

Os livros em `ai-service/books` incluem Manual dos Monstros e Guia do Mestre — conteúdo pra
mestre, não pra ficha de jogador. Auditar isso com o mesmo rigor das regras de personagem é um
projeto à parte (provavelmente ligado à Fase G do roadmap — itens mágicos — ou a uma futura
funcionalidade de bestiário pro mestre). **Escopo desta auditoria: tudo que afeta a ficha de um
personagem jogador** — classes, subclasses, raças, antecedentes, talentos, magias que PJs podem
usar, equipamento, progressão de nível. Se o Claude Code achar algo relevante de mestre no
caminho, anotar à parte, não tentar resolver junto.

Segundo ponto: uma auditoria genuinamente completa de centenas de magias e dezenas de talentos
não sai de uma vez. A prioridade certa é **o que afeta 100% dos personagens primeiro** (regras
de progressão, subclasses, listas fechadas já mapeadas), **depois o que afeta uma fração**
(listas completas de magias por classe, talentos individuais). Pedir pro Claude Code trabalhar
nessa ordem e reportar progresso por etapa, não tentar tudo de uma vez e arriscar qualidade.

## Recomendação técnica: extração direta do PDF, não RAG, para listas/tabelas

O `ingest.py` existente indexa os livros no ChromaDB pra busca semântica (`/skills`) — ótimo pra
"ache um trecho relevante sobre X", ruim pra "extraia a lista completa e exata de Y". Uma busca
semântica pode não retornar todos os chunks de uma tabela longa. Pra qualquer tarefa desta
auditoria que seja "extrair uma lista/tabela completa" (progressão de classe por nível, lista de
magias, lista de talentos), ler o texto do PDF diretamente (extração de texto simples do
capítulo relevante) em vez de depender de queries no ChromaDB.

## Sistema de progressão de nível — o que deve acontecer automaticamente a cada level up

`checkAndApplyLevelUps` já existe e já pausa em níveis de ASI aguardando `ApplyASI` (conforme o
CLAUDE.md). Verificar se ele cobre tudo abaixo — e não só o ganho de nível em si:

### Estrutura geral (vale pra qualquer classe)
- Bônus de proficiência aumenta em breakpoints fixos (2/3/4/5/6) — já implementado segundo o
  CLAUDE.md, só confirmar que o valor recalculado se propaga pra CD de magia, bônus de ataque e
  perícias, não só fica um número solto.
- PV máximo aumenta (dado de vida ou média + mod de Constituição) — já é um bug conhecido (#2)
  que isso não deve sobrescrever PV atual já ajustado manualmente pelo jogador; confirmar que a
  correção anterior também vale em level up, não só na criação.
- Escolha de subclasse acontece **uma vez, no nível certo pra cada classe** — e esse nível
  **não é o mesmo pra todas as classes**: Clérigo, Feiticeiro e Bruxo escolhem no nível 1;
  Mago no nível 2; a maioria das outras (Bárbaro, Bardo, Guerreiro, Ladino, Monge, Druida,
  Patrulheiro) no nível 3; Paladino no nível 3. Conferir se o sistema dispara essa escolha no
  nível certo por classe, e não um nível fixo genérico pra todas.
- ASI/talento nos níveis 4, 8, 12, 16, 19 — **mas atenção**: Guerreiro ganha ASI extra nos
  níveis 6 e 14 (além dos padrões), e Ladino ganha ASI extra no nível 10. Se o sistema trata ASI
  como "mesmos níveis pra qualquer classe", isso é uma lacuna real e silenciosa — exatamente o
  mesmo padrão do bug das escolas de magia, só que na progressão em vez de na criação.
- Ataque Extra: a maioria das classes marciais (Guerreiro, Bárbaro, Paladino, Patrulheiro,
  Monge) ganha no nível 5. Guerreiro é a exceção que ganha mais vezes: um segundo ataque extra
  no nível 11 (3 ataques) e um terceiro no nível 20 (4 ataques). Se isso afeta o cálculo de
  "Ataques e Magias" da ficha, precisa refletir o número certo de ataques por turno no nível
  certo, não só ligar/desligar uma flag booleana de "tem ataque extra".
- Para conjuradores: espaços de magia por nível seguem uma de três tabelas (conjurador completo,
  meio-conjurador como Paladino/Patrulheiro começando no nível 2, terço-conjurador como
  Cavaleiro Arcano/Ladino Arcano começando no nível 3) — confirmar que o sistema aplica a tabela
  certa por classe/subclasse, não uma tabela genérica.
- Truques (cantrips) conhecidos aumentam em níveis específicos (tipicamente 4 e 10 pra a maioria
  dos conjuradores de truques) — conferir contra o livro pra cada classe conjuradora, não
  assumir que é igual pra todas.
- Toda feature de classe ganha em cada nível (não só as "grandes" como subclasse/ASI/ataque
  extra) deveria gerar uma notificação visível pro jogador — isso conecta direto com o painel
  "Atividade Recente" do mockup da tela de personagens (Fase 4 do `TASKS_UI_E_FEATURES.md`):
  “[Personagem] subiu para o nível X e desbloqueou [feature]” em vez do jogador precisar abrir o
  livro pra saber o que mudou. Isso é literalmente o "automatizar pra facilitar o jogador" que
  foi pedido.

### O que fazer com o que não dá pra automatizar direito
Algumas escolhas exigem julgamento do jogador que o sistema não deve decidir sozinho (ex: qual
magia aprender entre as conhecidas por um Feiticeiro, qual talento escolher em vez de ASI). Pra
essas, o objetivo não é a automação escolher por ele — é **apresentar as opções válidas
automaticamente**, sem o jogador precisar ir atrás no livro pra saber quais são as opções
elegíveis naquele nível/classe/subclasse.

## O que peço pra não fazer sem antes confirmar comigo

- Não assumir que a lista de features/níveis que estou descrevendo acima está 100% certa em
  cada detalhe fino — é o framework estrutural, confiável no nível geral, mas o Claude Code tem
  acesso ao livro real e deve conferir contra ele, principalmente em classes onde eu não citei
  números específicos.
- Não expandir esta auditoria pra conteúdo de mestre (monstros, itens mágicos, encontros) sem
  perguntar antes — está fora do escopo de "ficha de personagem completa".
- Não tentar "completar tudo" numa única resposta gigante — trabalhar por prioridade (regras de
  progressão → subclasses → raças/antecedentes → talentos → magias) e reportar o que foi
  auditado/corrigido em cada etapa antes de seguir pra próxima.

## Critério de aceite

Um personagem de teste subindo de nível 1 a 20 (em pelo menos duas classes diferentes: uma
conjuradora completa e uma marcial) mostra, em cada nível, exatamente as features/escolhas que o
livro prevê pra aquele nível — nem a mais, nem a menos — sem o jogador precisar consultar nada
fora do próprio sistema.
