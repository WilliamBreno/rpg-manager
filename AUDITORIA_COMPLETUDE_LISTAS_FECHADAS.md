# Auditoria e correção: listas fechadas incompletas no seed de dados

> Gatilho: um usuário relatou que, ao criar um Mago, só apareciam 4 das 8 Escolas de Magia
> (faltavam Necromancia, Encantamento e outras). Escola de Magia é a subclasse do Mago (Tradição
> Arcana) — ou seja, isso não é um problema isolado das escolas, é uma instância de um padrão
> maior: **toda classe de 5e tem uma lista fechada de subclasses, e se uma lista fechada ficou
> parcialmente populada silenciosamente, é provável que outras também estejam.** Este arquivo
> existe pra achar e corrigir todas as instâncias desse padrão, não só a que foi relatada.

## Metodologia (aplicar a cada lista abaixo)

1. Definir a fonte canônica exata (livro + edição) — para 5e, isso é limitado ao que já existe
   em `ai-service/books`: Livro do Jogador (2014 e/ou 2024, conforme a decisão já pendente sobre
   qual regra o projeto segue) + Guia de Xanathar para Todas as Coisas. Não importar conteúdo de
   nenhum outro livro que não esteja fisicamente na pasta.
2. Contar quantos itens da lista existem hoje em `internal/seed`/`internal/domain`.
3. Contar quantos itens a fonte canônica realmente tem.
4. Se houver divergência, **listar por nome quem falta** — não é suficiente reportar "faltam
   alguns", porque isso é exatamente como esse bug passou despercebido até um usuário reportar.
5. Popular os que faltam **sem remover ou alterar os que já existem** — só corrigir se algo que
   já está lá estiver com nome/descrição errada, e nesse caso avisar antes de sobrescrever.
6. Pra cada item novo, aplicar o princípio já registrado no
   `DND5E_CRIACAO_PERSONAGEM_CHECKLIST.md`: é algo que o jogador **escolhe** (ex: escola de magia
   é escolhida no nível 2 do Mago) ou algo **informativo/automático** (ex: um traço racial fixo)?
   Confirmar que a UI trata isso do jeito certo pro tipo de dado, não só que o dado existe no
   banco.

## Correção confirmada: Escolas de Magia (Tradição Arcana do Mago)

Lista completa das 8 escolas, todas do Livro do Jogador básico: Abjuração, Convocação,
Adivinhação, Encantamento, Evocação, Ilusão, Necromancia, Transmutação. Popular as que
faltarem em `internal/seed` sem duplicar as 4 que já existem, e conferir que a escolha de
escola está de fato ligada ao personagem no nível certo (nível 2 do Mago), não só existindo
solta numa tabela sem uso real na criação/progressão.

## Generalização: subclasses de TODAS as classes (verificar cada uma com a mesma atenção)

Cada classe de 5e tem sua própria lista fechada de subclasses, com nome de categoria próprio
(a "escola de magia" do Mago é só o nome que a subclasse de Mago recebe). Referência de partida
abaixo, montada de memória — **tratar como ponto de partida a conferir contra o PDF real do
Livro do Jogador e do Guia de Xanathar, não como fonte definitiva**, exatamente pelo mesmo
motivo que gerou o bug relatado (um humano ou uma IA pode errar a contagem de cabeça):

| Classe | Nome da categoria de subclasse | Nº esperado no PHB básico | Nº esperado somando Xanathar's |
|---|---|---|---|
| Bárbaro | Trilha Primitiva | 2 | 5 |
| Bardo | Colégio de Bardo | 2 | 5 |
| Clérigo | Domínio Divino | 7 | 9 |
| Druida | Círculo Druídico | 2 | 4 |
| Guerreiro | Arquétipo Marcial | 3 | 6 |
| Ladino | Arquétipo de Ladino | 3 | 7 |
| Mago | Tradição Arcana (as "escolas") | 8 | 8 (sem adição no Xanathar's) |
| Monge | Tradição Monástica | 3 | 6 |
| Paladino | Juramento Sagrado | 3 | 5 |
| Patrulheiro (Ranger) | Arquétipo de Patrulheiro | 2 | 5 |
| Bruxo (Warlock) | Pacto Extraplanar | 3 | 5 |
| Feiticeiro | Origem Feiticeira | 2 | 5 |

Rodar a mesma auditoria de contagem-por-nome pra cada linha desta tabela. Se a contagem atual do
seed for menor que "Nº esperado no PHB básico", já é uma lacuna confirmada, independente do
Xanathar's.

## Outras listas fechadas a auditar com a mesma lente (não assumir que estão completas)

- **Raças e sub-raças** (5e): lista do PHB (Humano, Elfo com sub-raças, Anão com sub-raças,
  Halfling com sub-raças, Draconato, Gnomo com sub-raças, Meio-Elfo, Meio-Orc, Tiefling) —
  conferir se todas as sub-raças estão presentes, não só a raça-mãe.
- **Antecedentes**: lista do PHB (Acólito, Charlatão, Criminoso, Artista, Herói do Povo, Artesão
  de Guilda, Eremita, Nobre, Forasteiro, Sábio, Marinheiro, Soldado, Órfão de Rua) — conferir
  contagem exata contra o livro.
- **Talentos**: lista completa do PHB (é uma lista longa — contar exatamente quantos o livro
  tem e comparar).
- **Estilos de Luta** (Guerreiro, Patrulheiro, Paladino): Arquearia, Defesa, Duelismo, Combate
  com Arma Grande, Proteção, Combate com Duas Armas — conferir se as 6 do PHB básico estão
  todas presentes.
- **Perícias**: já coberto no checklist anterior (18 no total), mas vale reconferir a contagem
  aqui também, já que é o mesmo tipo de bug.
- **Idiomas**: comuns (8) + exóticos (8) do PHB.
- **Alinhamentos**: os 9 combinações de eixo lei/caos e bem/mal.
- **Condições** (cego, enfeitiçado, paralisado etc.): usadas como referência, mesmo que não
  afetem diretamente a ficha exportável, valem conferência se a Biblioteca do app expõe isso.
- **Tabela de progressão de magia** (espaços de magia por nível de personagem, por tipo de
  conjurador — completo/meio/terço): isso é uma lista fechada numérica, não de nomes, mas o
  mesmo risco de "só preenchi até certo nível e esqueci o resto" se aplica.

## E o lado 4e?

Diferente do 5e, **não há nenhum livro de 4e em `ai-service/books`** — então não dá pra aplicar
a mesma auditoria "contagem contra o PDF" sem antes descobrir de onde o conteúdo 4e do seed
veio originalmente. Antes de tentar auditar completude de 4e (classes, poderes, raças), o
Claude Code precisa responder: qual foi a fonte usada quando o seed de 4e foi escrito? Isso é
provavelmente a mesma causa raiz do bug já conhecido #1 (perícias do livro "Poder Arcano" não
aparecem nas queries de IA) — os dois sintomas apontam pra a mesma pergunta: o conteúdo de 4e
tem uma fonte de referência completa em algum lugar, ou foi digitado manualmente e por isso é
inerentemente parcial? Essa resposta decide se completude de 4e é uma tarefa de "auditar contra
um livro" ou de "primeiro adquirir/organizar a fonte de referência".

## Critério de aceite

Uma tabela final, por lista fechada, com: contagem atual → contagem oficial → itens faltantes
nomeados → status (corrigido / aguardando confirmação do usuário / bloqueado por falta de fonte
de referência, no caso do 4e). Nenhuma lista marcada como "completa" sem a contagem ter sido
conferida de verdade contra a fonte, não de memória.
