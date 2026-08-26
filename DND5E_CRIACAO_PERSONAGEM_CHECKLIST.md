# Auditoria: mecânicas de criação de personagem 5e no RPG Manager

> Objetivo deste arquivo: antes de qualquer trabalho novo, o Claude Code audita o que já existe
> hoje no fluxo de criação de personagem 5e (frontend de criação + `internal/domain` +
> `internal/service` + seed) contra a lista abaixo, e preenche a coluna Status. **Não é pra
> implementar nada nesta etapa** — é pra descobrir com precisão o que falta antes de tocar em
> código, porque a exportação pra ficha PDF (Fase 2 do `TASKS_UI_E_FEATURES.md`) só fica correta
> se o personagem tiver todo esse dado pra exportar.
>
> Status possíveis: `✅ existe` / `⚠️ parcial` (explicar o que falta) / `❌ não existe`.
> Pra cada item, citar o arquivo/linha onde encontrou (ou confirmar a ausência).

> **Nota de pesquisa (adicionada depois da primeira versão deste checklist):** o usuário
> forneceu a transcrição de dois vídeos introdutórios sobre criação de personagem 5e. Ambos
> cobrem só o fluxo básico (raça/classe → atributos → antecedente → personalidade →
> equipamento/magias; e nome/classe/nível/raça/antecedente → atributos e como modificam
> perícias → proficiências → CA/PV/velocidade → ataques e magias). Isso corresponde às seções
> 0, 1, 2, 3, 4, 6 e 8 abaixo — nenhum item novo surgiu da revisão. Nenhum dos vídeos menciona
> Expertise, Jack of All Trades, tabelas de espaço de magia, subclasses, talentos ou a diferença
> 2014/2024 — ou seja, as lacunas mais prováveis de causar exportação errada continuam sendo as
> das seções 3 e 5, não cobertas por esses vídeos.

## 0. Pergunta que decide o resto do checklist

**Qual conjunto de regras 5e o projeto segue: 2014, 2024, ou os dois?** Isso não é
cosmético — muda mecânica real:

| | 2014 | 2024 |
|---|---|---|
| Origem do bônus de atributo | Raça | Antecedente (Background) |
| Feat no nível 1 | Não (variante opcional) | Sim, todo personagem recebe um "Origin Feat" |
| Nome do conceito "raça" | Race | Species (algumas mudam traços) |
| Perícias do antecedente | Fixas por antecedente | Jogador escolhe 2 de uma lista |

O `ai-service/books` tem tanto o Livro do Jogador 2024 quanto versões que parecem ser de 2014.
**Parar aqui e confirmar com o usuário qual regra vale** antes de marcar qualquer item abaixo
como "falta implementar" — pode ser que о projeto já tenha escolhido uma das duas de propósito.

## 0.5 Princípio de automação (importante pra interpretar o resto deste checklist)

A ideia por trás desta auditoria não é só "existe o campo no banco?" — é "o jogador vê isso
automaticamente ao escolher raça/classe, sem precisar abrir um livro ou pedir pra IA?". Isso
separa cada item em duas categorias, e a categoria certa muda a solução técnica certa:

- **Universal** (existe pra qualquer personagem, independente de raça/classe): os 6 atributos,
  a lista de 18 perícias, as opções de alinhamento, os antecedentes disponíveis, a lista de
  equipamentos pra comprar, os idiomas comuns. Esses **sempre aparecem** na tela de criação,
  não dependem de nenhuma escolha prévia.
- **Condicional a raça/classe** (só existe, ou muda, de acordo com a escolha): traços raciais
  (visão no escuro, resistências), perícias de classe pra escolher, salvaguardas proficientes,
  dado de vida, atributo de conjuração, magias conhecidas/preparadas, talento de origem. Esses
  devem **aparecer automaticamente assim que o jogador escolhe a raça/classe**, sem ele
  precisar ir atrás em lugar nenhum.

**Atenção arquitetural, e isso vale destacar pro Claude Code:** "automático" não significa
"gerado pela IA em tempo real". A maior parte do que é condicional acima (Seções 3, 4, 5 e 7)
tem valor fixo e conhecido pelas regras: visão no escuro do Elfo é sempre 18m, salvaguarda
proficiente do Mago é sempre Int/Sab, dado de vida do Bárbaro é sempre d12. Isso deve vir de
**seed estático** (`internal/seed`), consultado instantaneamente — não de uma chamada ao
`ai-service`/RAG a cada vez, porque uma resposta gerada por LLM pode variar entre chamadas ou
alucinar um valor, e isso contaminaria exatamente o dado que vai pro PDF exportado. Reservar o
`/skills` (RAG) do `ai-service` pro que o CLAUDE.md já descreve como seu propósito real: gerar
as descrições dos poderes de classe estilo 4e (`power_type`: at-will/encounter/daily/utility),
que são numerosos e mais fáceis de indexar de um livro do que tabular à mão — não pra decidir
bônus, proficiência ou espaço de magia, que têm resposta certa e única.

Na prática: quando a auditoria abaixo encontrar uma lacuna nas Seções 3 a 7, a pergunta certa
não é "falta plugar isso no `/skills`", é "falta essa tabela em `internal/seed`". Marcar isso
explicitamente na coluna de observação de cada item da tabela final.

---

## 1. Identidade básica

| Item | Onde verificar |
|---|---|
| Nome do personagem | `internal/domain` (`Character.Name` ou equivalente) |
| Nome do jogador (dono da conta, não do personagem) | relação `Character` → `User` |
| Raça/Espécie (+ sub-raça/linhagem, se existir no ruleset escolhido) | `internal/domain`, `internal/seed` |
| Classe (+ subclasse — em qual nível a subclasse é escolhida? varia por classe: Clérigo/Feiticeiro no 1, Guerreiro/Mago no 2/3) | `internal/domain`, `internal/service` |
| Nível | `Character.Level` |
| Antecedente (Background) | `internal/domain` — CLAUDE.md cita tanto `Antecedent` quanto `Background` como models separados, **confirmar se isso é intencional (dois conceitos diferentes) ou duplicação/confusão a resolver** |
| Alinhamento | `Character.Alignment` (CLAUDE.md já cita como campo 5e-only) |
| Pontos de experiência (XP) | `Character.XP`, tabela `xpTable5e` |

## 2. Atributos e modificadores

| Item | Onde verificar |
|---|---|
| 6 atributos (For/Des/Con/Int/Sab/Car) armazenados | `internal/domain` |
| Método de geração do atributo — **qual(is) o sistema oferece?** Array padrão (15,14,13,12,10,8), point buy, 4d6 descarta o menor, ou fixo. Se for mais de um, o jogador escolhe? | fluxo de criação no frontend |
| Bônus de atributo aplicado corretamente conforme a resposta do item 0 (raça em 2014, antecedente em 2024) | `internal/service` |
| Cálculo do modificador `(score-10)/2` arredondado pra baixo | `internal/service` |
| ASI (Ability Score Improvement) nos níveis corretos pra 5e (4, 8, 12, 16, 19) — já documentado no CLAUDE.md como existente | confirmar que `ApplyASI` cobre também a opção de pegar um **feat** em vez de +2/+1/+1+1, não só pontos de atributo |

## 3. Perícias, proficiências e talentos

| Item | Onde verificar |
|---|---|
| Lista das 18 perícias (`Pericia`) com atributo associado correto | `internal/seed` |
| Proficiência em perícia vinda da classe (jogador escolhe N de uma lista, não todas) | fluxo de criação |
| Proficiência em perícia vinda do antecedente | fluxo de criação |
| Proficiência em perícia vinda de raça/feat (alguns concedem perícia fixa, ex: Elfo → Percepção) | `internal/seed` |
| **Expertise** (Ladino nível 1 e 6, Bardo nível 2 — dobra o bônus de proficiência em perícias escolhidas) — mecânica frequentemente esquecida | procurar por algo tipo `Expertise`/`Pericia.Expertise` no domain; se não existir, é uma lacuna real, não cosmética, porque muda o número final que vai pro PDF |
| **Jack of All Trades** (Bardo — some do modificador de proficiência em testes de perícia sem proficiência) | idem, procurar tratamento especial pra Bardo |
| Talentos (`Talento`) — de onde vêm: feat de antecedente (2024), ASI trocado por feat (ambos), ou variante humano (2014)? | `internal/domain`, `internal/service` |
| Bônus de Sabedoria Passiva (Percepção) = 10 + bônus total de Percepção | verificar se é calculado ou só um campo livre |
| Proficiências em armas/armaduras/ferramentas (da classe + antecedente + raça) | `internal/domain` — existe algum campo pra isso, ou só fica em "ProficienciesLang" como texto solto? Se for só texto livre, a ficha PDF até aceita (é um campo de texto), mas a lógica de cálculo de ataque (item 4) não pode depender de texto livre |
| Idiomas conhecidos (comum + raça + antecedente + feats) | `internal/domain` |

## 4. Combate

| Item | Onde verificar |
|---|---|
| Classe de Armadura — cálculo correto por tipo de armadura (leve/média/pesada, limite de mod de Destreza pra média, sem limite pra leve, sem bônus de Destreza pra pesada), + escudo +2 | `internal/service` — **atenção**: o bug conhecido #3 (CA errada de Bardo 4e) mostra que esse cálculo já teve erro; conferir se a mesma classe de bug não existe também do lado 5e |
| Defesa Desarmada (Unarmored Defense) — Bárbaro usa Con, Monge usa Sab, em vez do cálculo padrão | procurar tratamento de classe especial |
| Iniciativa = mod de Destreza (+ feats como Alerta) | `internal/service` |
| Velocidade de deslocamento vinda da raça (Fase B do roadmap já cobre isso — confirmar que ficou pronta antes de assumir que falta) | `internal/domain`/`internal/seed` |
| Dado de vida por classe (`d6`/`d8`/`d10`/`d12`) e cálculo de PV máximo (dado + mod de Con no nível 1; média do dado arredondada pra cima + mod de Con por nível seguinte, ou role manual) | `internal/service` — **cruzar com o bug #2 já conhecido**: o PV deve poder ser digitado pelo jogador e não sobrescrito |
| Bônus de proficiência por nível (tabela 5e: +2 até nv4, +3 até nv8, +4 até nv12, +5 até nv16, +6 até nv20) | `internal/service` |
| Testes de resistência (salvaguardas) — quais 2 atributos cada classe é proficiente | `internal/seed` |
| Testes contra a morte (sucessos/fracassos) — campo já existe no domain segundo o CLAUDE.md (`DeathSaveSuccesses/Failures`) | confirmar se o frontend expõe isso em algum lugar, ou só existe no banco sem UI |
| Ataques: bônus de ataque = proficiência (se proficiente na arma) + mod de atributo (Força ou Destreza pra armas "finesse"/à distância) + bônus mágico | `internal/service` — checar se existe cálculo de ataque em algum lugar ou se hoje é só texto livre digitado pelo jogador |
| Dano da arma: dado + mod de atributo | idem |
| Ataque Extra (Extra Attack) — Guerreiro nv5, outras classes marciais — múltiplos ataques por turno | procurar feature de classe correspondente |

## 5. Conjuração (se a classe conjura)

| Item | Onde verificar |
|---|---|
| Atributo de conjuração por classe (Int para Mago, Sab para Clérigo/Druida/Patrulheiro, Car para Bardo/Feiticeiro/Bruxo/Paladino) | `internal/seed`/`internal/service` |
| CD de resistência de magia = 8 + proficiência + mod do atributo de conjuração | `internal/service` |
| Bônus de ataque de magia = proficiência + mod do atributo | `internal/service` |
| Truques conhecidos (cantrips) — não gastam espaço de magia | `internal/domain` |
| Magias conhecidas vs. preparadas — **isso varia por classe** (Feiticeiro/Bardo/Patrulheiro/Ladino-Arcano "conhecem" um número fixo; Clérigo/Druida/Mago/Paladino "preparam" a partir de uma lista maior todo dia). O sistema já distingue isso ou trata todo mundo igual? | `internal/domain`/`internal/service` — se tratar igual, é uma lacuna real |
| Tabela de espaços de magia por nível/classe (conjurador completo, meio-conjurador tipo Paladino/Patrulheiro, terço-conjurador tipo Cavaleiro Arcano) | `internal/seed` |
| Espaços de magia usados/disponíveis (tracker, não só o total) | `internal/domain` |
| Magia ritual (algumas classes podem conjurar certas magias sem gastar espaço, fora de combate) | provavelmente não crítico pro PDF, mas vale marcar como não coberto |
| Foco de conjuração / bolsa de componentes (afeta se o personagem consegue pagar componentes materiais) | equipamento |

## 6. Antecedente e traços de personagem

| Item | Onde verificar |
|---|---|
| Traços de personalidade, ideais, vínculos, defeitos (campos de texto livre do antecedente) | `internal/domain` — já mapeados no `dnd5e_pdf_field_map.json` |
| Perícias/ferramentas/idiomas concedidos pelo antecedente | `internal/seed` |
| Equipamento inicial concedido pelo antecedente | `internal/seed` — cruza com a Fase D do roadmap (inventário/equipamento), que ainda está só planejada, não implementada — **não assumir que isso já existe** |
| Traço de antecedente (feature nomeada, ex: "Refúgio dos Fiéis") | `internal/domain` |

## 7. Traços raciais

| Item | Onde verificar |
|---|---|
| Visão no escuro (Darkvision) e alcance | `internal/seed` — Fase F do roadmap ("traços raciais como entidade separada") ainda está só planejada; se isso não existe ainda, é esperado, não é bug |
| Resistências/imunidades de raça (ex: Anão resiste a veneno) | idem |
| Tamanho (Pequeno/Médio) — afeta capacidade de carga e algumas armas | `internal/domain` |
| Magias raciais inatas (ex: Tiefling conjura Taumaturgia) | Fase G do roadmap ("rituais e itens mágicos") — ainda planejada |
| Idiomas raciais | Fase E do roadmap ("idiomas raciais com seed data") — ainda planejada, confirmar se já saiu do papel |

## 8. Equipamento e recursos

| Item | Onde verificar |
|---|---|
| Inventário de equipamento (Fase D do roadmap) | confirmar status real — se ainda não implementado, a ficha PDF vai exportar o campo "Equipment" como texto livre digitado, não uma lista estruturada; **decidir se isso é aceitável pra v1 do export ou se bloqueia a Fase 2** |
| Moedas (PC/PP/PE/PO/PL) | `internal/domain` |
| Capacidade de carga (Força × 15) | provavelmente não crítico, marcar como não coberto se não existir |
| Armas equipadas (pra puxar automaticamente no bloco "Ataques" da ficha) | depende do item de inventário acima |
| Armadura equipada (`ArmorID`, já corrigido pra nullable segundo o histórico do projeto) | `internal/domain` |

## 9. Descrição física e narrativa

| Item | Onde verificar |
|---|---|
| Idade, altura, peso, olhos, pele, cabelo | `internal/domain` — Fase C do roadmap ("detalhes pessoais") ainda planejada, confirmar status |
| Símbolo de fé/divindade | idem, ligado à Fase C |
| História do personagem (texto livre) | `internal/domain` |
| Aliados e organizações, tesouro (texto livre) | `internal/domain` |

---

## Cruzamento com os vídeos-tutorial consultados

O usuário pediu análise de três vídeos como insumo pra este checklist. Dois foram transcritos
manualmente (o terceiro, um vídeo de playlist mais longo, não foi). Registro aqui o que eles
confirmam, pra deixar claro o que já estava coberto e por quê:

- **Vídeo 1** (visão geral em 5 passos): raça+classe → atributos (4d6 descarta o menor, ou
  arranjo padrão) → antecedente → traços de personalidade/ideais/vínculos/defeitos →
  equipamento e magias. Cobre só as seções 1, 2 (parcialmente) e 6 deste checklist.
- **Vídeo 2** (visão geral com foco nos números derivados): nome/classe/nível/raça/antecedente
  como base; atributos gerados por dado, sistema de pontos ou arranjo padrão; perícias vindas de
  classe+antecedente; CA = armadura + Destreza; PV = classe + Constituição; velocidade vinda da
  raça; bloco de ataques/magias. Cobre partes das seções 1, 2, 3 e 4.

**Conclusão honesta:** os dois vídeos são conteúdo introdutório de divulgação, não referência de
regras. Eles confirmam a camada mais superficial do processo (a mesma que qualquer resumo de uma
página do PHB cobre) e **não mencionam nenhum item das seções que este checklist já sinalizava
como "frequentemente esquecido"**: Expertise, Jack of All Trades, a distinção entre conjurador
que "conhece" magias fixas e o que "prepara" a partir de uma lista maior, Defesa Desarmada
(Bárbaro/Monge), testes contra a morte, ou a divergência de regras 2014 vs 2024 da Seção 0. Isso
não invalida os vídeos — só confirma que uma auditoria baseada neles teria ficado rasa
exatamente nos pontos que mais custam caro se faltarem na hora de exportar a ficha. Nenhum item
do checklist foi alterado por causa dessa checagem; ela serve só como registro de que a fonte
foi considerada.

---

## Saída esperada desta auditoria

Uma tabela markdown (pode ser neste mesmo arquivo, numa seção nova ao final) com uma linha por
item acima, preenchida com Status + arquivo/linha + observação. **Não implementar nada ainda.**
Depois de eu revisar essa tabela, decidimos juntos:

1. Quais lacunas bloqueiam a Fase 2 (export PDF) — ou seja, campos que a ficha oficial espera e
   que hoje não existem em lugar nenhum do sistema, nem como texto livre.
2. Quais lacunas são aceitáveis por ora (ex: o jogador preenche manualmente aquele campo, e o
   PDF exporta o texto livre) — não é preciso ter a mecânica 100% automatizada pra exportar
   corretamente, só é preciso que o dado exista em algum lugar.
3. Quais lacunas são trabalho de fases futuras já planejadas no roadmap (C, D, E, F, G) e não
   devem ser antecipadas sem eu pedir.

---

# Resultado da auditoria

## Resposta à Seção 0

**Nem 2014 nem 2024 de forma automática, pra bônus de atributo — nenhum dos dois está
implementado.** `domain.Race` (`internal/domain/race.go`) só tem `Speed`,
`BonusTrainedSkills`/`BonusSkillValues`/`BonusTalentos` (todos explicitamente comentados como
"(4e)" no código) — nenhum campo de bônus de atributo. `domain.Antecedent`
(`internal/domain/antecedent.go`) também não tem nenhum campo de bônus de atributo. A criação de
personagem (`CharacterCreate.tsx`, schema Zod linhas 25-30) usa entrada livre `1-20` pros 6
atributos, sem nenhum bônus automático aplicado por cima — já documentado no `CLAUDE.md` como
simplificação deliberada. Então isso não é uma ambiguidade 2014-vs-2024 que precise de decisão
seu — é um gap único que serve pros dois modelos igualmente.

Os outros três pontos da tabela da Seção 0 também dão pra responder por evidência direta, sem
precisar perguntar:
- **Feat no nível 1:** o campo `Class.TalentosCount` só é populado para classes 4e
  (`fixClassTalentosCount4e(db)`, único call site de `TalentosCount` em todo `seed.go`) — pra
  toda classe 5e ele fica 0/não setado, e o frontend cai no fallback `?? 1`
  (`CharacterCreate.tsx:172`), então **todo personagem 5e criado hoje é obrigado a escolher
  exatamente 1 talento**, sem restrição de lista por classe/antecedente. Funciona parecido com o
  "Origin Feat" de 2024 na prática (1 feat garantido no nível 1), mas não é modelado como tal:
  não está ligado ao antecedente escolhido, e `Talento` (`internal/domain/talento.go`) não tem
  nenhum campo que marque "é um feat de origem" vs "é um feat qualquer".
- **Nome do conceito "raça":** o código usa `Race`/`RaceID`/"raça" em todo lugar (domain, seed,
  frontend) — nunca "Species"/"espécie". Ruleset 2014 nesse ponto específico.
- **Perícias do antecedente:** `Antecedent.SkillProficiencies` (`antecedent.go:10`) é uma lista
  JSON fixa por antecedente (ex. `["Intuição","Religião"]`), não uma escolha de 2 de uma lista
  maior — mas isso não decide 2014 vs 2024 sozinho, porque no PHB 2024 os antecedentes também têm
  2 perícias fixas por antecedente (a diferença real de 2024 nesse ponto é o *Origin Feat* ligado
  ao antecedente, coberto acima).

**Conclusão prática:** o projeto não seguiu conscientemente nem 2014 nem 2024 pra esses pontos —
implementou uma mistura pragmática (raça com nome antigo, sem bônus de atributo de fonte nenhuma,
antecedente com perícias fixas, 1 feat garantido no nível 1 meio-2024-mas-desacoplado do
antecedente). Ver a Seção 2 abaixo pra decisão sobre se vale a pena formalizar isso.

## 1. Identidade básica

| Item | Status | Onde | Observação |
|---|---|---|---|
| Nome do personagem | ✅ existe | `Character.Name` (`character.go:7`) | — |
| Nome do jogador (dono da conta) | ✅ existe | `Character.UserID` → `User` (`character.go:28`) | relação existe; não há um campo de "display name do jogador" separado, mas dá pra puxar de `User.Name` |
| Raça/Espécie (+ sub-raça/linhagem) | ⚠️ parcial | `Character.RaceID` → `Race` (`race.go`) | raça existe; **sub-raça/linhagem não existe como campo separado** — `Race` é uma entidade "achatada" (ex: "Elfo" é uma raça, não "Elfo" + sub-raça "Alto Elfo"); linhagens tipo Tiferino/Aasimar aparecem só como `RequiresChoice`/`ChoiceGroup` em `Skill` (`legado_tiferino`, `revelacao_aasimar` em `seed_races_5e.go`), não como um campo estruturado de personagem |
| Classe (+ subclasse) | ⚠️ parcial | `Character.ClassID` → `Class`; subclasse via `Skill.ChoiceGroup` | **implementação inconsistente entre classes, e nenhuma respeita o nível real de escolha do RAW.** Clérigo (`dominio_clerigo`), Bruxo (`patrono_bruxo`), Feiticeiro (`origem_feiticeiro`), Guardião/Ladino (`especialista_guardiao`/`especialista_ladino`) têm uma "escolha de subclasse" seedada, mas **todas em `Level: 1`**, mesmo quando o RAW manda nível 2 (Mago) ou 3 (Guerreiro). E pior: **Mago não tem nenhuma Escola de Magia/Tradição Arcana seedada** (busquei "Escola"/"Tradição Arcana" em `seed.go`, zero resultados) e **Guerreiro não tem Arquétipo Marcial** (só tem "Estilo de Luta", que é uma feature diferente e RAW-correta de estar no nível 1, mas não substitui o Arquétipo Marcial do nível 3) — essas duas classes simplesmente não têm mecânica de subclasse nenhuma hoje |
| Nível | ✅ existe | `Character.Level` (`character.go:9`) | — |
| Antecedente (Background) | ⚠️ confirmar como intencional | dois models: `domain.Antecedent` (`antecedent.go`) e `domain.Background` (`background.go`) | **são de fato dois conceitos diferentes, mas o segundo tem um nome enganoso.** `Antecedent` é o catálogo real de antecedentes 5e (Acólito, Criminoso etc, com perícias/ferramentas/idiomas/equipamento/feature) e é o que `Character.AntecedentID` referencia de verdade. `Background` (`background.go`) parece um catálogo idêntico (mesmos campos: `SkillProficiencies`/`ToolProficiencies`/`Languages`/etc) mas **não é usado como catálogo em lugar nenhum** — `BackgroundService.Get/Save` (`background_service.go`) na verdade lê/escreve direto 4 campos soltos em `Character` (`PersonalityTraits`/`Ideals`/`Bonds`/`Flaws`, agora +`Rumors`) via um `map[string]interface{}`, sem nunca tocar na tabela `backgrounds`. Ou seja: `domain.Background` é uma tabela morta/vestigial que só coincide de nome com a rota `/characters/:id/background`, que na real é só um editor de campos de biografia do `Character`. Vale confirmar com você se dá pra remover `domain.Background` (parece não ter nenhum dado real usando a tabela) ou se ela é usada em algum lugar que não encontrei |
| Alinhamento | ✅ existe | `Character.Alignment` (`character.go:43`) | — |
| Pontos de experiência (XP) | ✅ existe | `Character.ExperiencePoints` (`character.go:12`), `xpTable5e` (`character_service.go:32`) | — |

## 2. Atributos e modificadores

| Item | Status | Onde | Observação |
|---|---|---|---|
| 6 atributos armazenados | ✅ existe | `character.go:33-38` | — |
| Método de geração do atributo | ❌ não existe | `CharacterCreate.tsx` schema (linhas 25-30) | não existe NENHUM método assistido (array padrão/point buy/4d6) — é digitação livre `1-20` validada só por `min(1).max(20)` no Zod. Já documentado no CLAUDE.md como simplificação deliberada, então isto é uma confirmação, não uma descoberta nova |
| Bônus de atributo por raça/antecedente | ❌ não existe | — | ver resposta da Seção 0 acima — nenhum campo de bônus em `Race` nem `Antecedent`, então isso nunca é aplicado automaticamente, pra nenhuma das duas fontes possíveis |
| Modificador `floor((score-10)/2)` | ✅ existe (corrigido) | `mod()` em `character_service.go` | já documentado no CLAUDE.md — era `(attr-10)/2` sem floor, corrigido nesta sessão |
| ASI nos níveis corretos (4,8,12,16,19) | ✅ existe | `isASILevel()` (`character_service.go`) | confirmado já existente e correto |
| ASI com opção de trocar por um feat | ❌ não existe | `ApplyASI`/`ASIChoice` (`character_service.go:227-234`) | `ASIChoice` só tem os 6 campos de atributo (`Strength`...`Charisma`) — não existe um caminho pra "pegar um talento em vez de +2/+1+1" nessa struct nem no handler. Hoje, se o jogador quer um talento num nível de ASI, o único jeito é usar `POST /characters/:id/talentos/:talento_id` por fora, sem nenhuma verificação de que isso substituiu (ou não) o ASI daquele nível — ou seja, tecnicamente dá pra pegar os dois ao mesmo tempo sem o sistema perceber |

## 3. Perícias, proficiências e talentos

| Item | Status | Onde | Observação |
|---|---|---|---|
| 18 perícias com atributo correto | ✅ existe | `domain.Pericia` (`pericia.go`), seed correspondente | não conferi as 18 uma por uma contra o atributo oficial, mas a estrutura existe e é usada de ponta a ponta (`pdf_export_service.go:83-91`) |
| Proficiência de perícia vinda da classe (escolhe N) | ✅ existe | `Class.AvailableSkills` + `Class.TrainedSkillsCount`, fluxo em `CharacterCreate.tsx` | — |
| Proficiência de perícia vinda do antecedente | ✅ existe | `Antecedent.SkillProficiencies`, aplicado em `CharacterCreate.tsx` (`backgroundSkills`) | fixo por antecedente, não escolhido (ver Seção 0) |
| Proficiência de perícia vinda de raça/feat fixa | ❌ não existe pra 5e | `Race.BonusSkillValues` existe mas é comentado "(4e)" (`race.go:14`) | nenhuma raça 5e (ex: Elfo → Percepção) concede perícia fixa automaticamente hoje |
| **Expertise** | ❌ não existe | busquei `Expertise`/"dobra" em todo o backend, zero resultados | confirmado gap real, não cosmético — `pdf_export_service.go:86-91` faz só `mod + profBonus` linear pra toda perícia proficiente, nunca dobra pra Ladino nv1/6 ou Bardo nv2. Isso silenciosamente exporta o número errado na ficha PDF pra qualquer Ladino/Bardo com Expertise |
| **Jack of All Trades** | ❌ não existe | mesmo arquivo/linhas acima | perícia sem proficiência sempre soma só o mod puro, nunca metade do bônus de proficiência pro Bardo |
| Talentos — de onde vêm | ⚠️ parcial | ver Seção 0 acima | 1 talento obrigatório e livre no nível 1 pra todo personagem 5e (fallback de `TalentosCount`), escolhido de toda a lista de `Talento` sem filtro por classe/antecedente/pré-requisito real (pré-requisito é só texto exibido, nunca validado — `CharacterCreate.tsx:760-762`) |
| Sabedoria Passiva (Percepção) | ✅ existe, calculado | `pdf_export_service.go:93-96` | `10 + valor de Percepção` — só é calculado no momento do export do PDF, não é um campo que aparece na tela de detalhe do personagem hoje |
| Proficiências em armas/armaduras/ferramentas | ⚠️ parcial | `Antecedent.ToolProficiencies` (texto livre); nada estruturado pra armas/armaduras | não existe nenhum campo em `Class` pra proficiência de arma/armadura (compare com `Class.SavingThrows`, que É estruturado) — hoje não tem como o sistema saber programaticamente "esse personagem é proficiente com espada longa?" |
| Idiomas conhecidos | ❌ não existe estruturado | `Antecedent.Languages` é só texto livre (ex: "Dois idiomas à sua escolha") | não existe uma lista real de idiomas conhecidos por personagem — nem um catálogo de idiomas, nem uma tabela de associação |

## 4. Combate

| Item | Status | Onde | Observação |
|---|---|---|---|
| Classe de Armadura (leve/média/pesada + escudo) | ✅ existe, correto | `ArmorService.CalculateAC` (`armor_service.go:21-77`) | limite de DEX por `Armor.MaxDexBonus` (-1 = sem limite, 0 = sem bônus, N = capado), escudo tratado à parte. **Não achei sinal do bug de CA do Bardo 4e mencionado no checklist original nem de um equivalente 5e** — o cálculo 5e (`default` no switch de edição, linha 58-77) está correto pra armadura leve/média/pesada |
| Defesa Desarmada (Bárbaro Con, Monge Sab) | ✅ existe, correto | `specialClassAC()` (`armor_service.go:81-102`) | `10 + dexMod + wisMod` pro Monge, `10 + dexMod + conMod` pro Bárbaro — bate com o RAW |
| Iniciativa = mod Destreza | ✅ existe | `pdf_export_service.go:132` | só calculado no momento do export, não exposto na tela de detalhe hoje |
| Velocidade de deslocamento da raça | ❌ **bug real, não "falta implementar"** | `Race.Speed` existe (`race.go:10`) mas **nunca é copiado pra `Character.Speed`** | busquei toda atribuição a `character.Speed`/`Speed:` em `character_service.go` — nenhuma. `CharacterCreate.tsx` também nunca envia `speed` no payload de criação. Resultado: `Character.Speed` fica sempre `0` (zero-value do Go) pra todo personagem 5e criado, e `pdf_export_service.go:133` exporta esse `0` direto pro campo "Deslocamento" da ficha PDF. **Isso é um bug ativo que quebra a exportação hoje, não uma lacuna de fase futura** — merece prioridade separada das lacunas "aceitáveis por ora" |
| Dado de vida por classe + cálculo de PV | ✅ existe | `calcHP5e` (`character_service.go`), já documentado no CLAUDE.md | PV manual do jogador é respeitado na criação (bug #2 do histórico, já corrigido) |
| Bônus de proficiência por nível | ✅ existe, correto | `proficiencyBonus5e()` (`character_service.go`) | tabela +2/+3/+4/+5/+6 confirmada batendo com o RAW |
| Salvaguardas proficientes por classe | ✅ existe | `Class.SavingThrows` (seedado pras 12 classes 5e, `seed.go:287-377`), consumido em `pdf_export_service.go:52-57,73-81` | cálculo correto (`mod + profBonus` se proficiente); só é exposto no export do PDF, não na tela de detalhe |
| Testes contra a morte | ✅ existe, com UI | `DeathSaves.tsx`, renderizado em `CharacterDetail.tsx:471` (`{is5e && <DeathSaves .../>}`) | **não é só banco sem UI como o checklist original suspeitava** — está exposto e interativo |
| Bônus de ataque / dano de arma | ❌ não existe | busquei `AttackBonus`/`WeaponDamage`/cálculo equivalente em `internal/service`, zero resultados | confirmado: hoje não existe nenhum cálculo de ataque/dano — não há sequer um campo de "arma equipada" ligado a um cálculo; isso teria que vir do sistema de inventário (`CharacterItem`, já existe desde esta sessão) cruzado com um novo cálculo de ataque, que não existe |
| Ataque Extra (Extra Attack) | ❌ não existe | busquei "Ataque Extra"/"Extra Attack" em todo `internal/seed`, zero resultados | nenhuma classe marcial tem essa feature seedada, nem como `Skill` descritiva |

## 5. Conjuração

| Item | Status | Onde | Observação |
|---|---|---|---|
| Atributo de conjuração por classe | ⚠️ só como texto | `Skill.Description` (ex: "Atributo de conjuração: Sabedoria" na descrição de "Conjuração de Clérigo") | existe só como prosa dentro da descrição da `Skill`, não como um campo estruturado em `Class` (compare com `SavingThrows`, que é estruturado) — nada no código lê esse valor pra calcular DC/ataque de magia |
| CD de resistência de magia | ❌ não existe | busquei em `pdf_export_service.go` e `character_service.go`, não calculado em nenhum lugar | confirmado pelo próprio comentário do field map da IA (ver abaixo) |
| Bônus de ataque de magia | ❌ não existe | idem | idem |
| Truques conhecidos (cantrips) | ❌ não existe | nenhum campo em `domain` | — |
| Magias conhecidas vs. preparadas | ❌ não existe distinção | nenhum campo em `domain` | trata todo mundo igual porque **não trata magia nenhuma** — não há nenhuma lista de magias por personagem hoje |
| Tabela de espaços de magia | ❌ não existe | nenhum arquivo de seed relacionado | — |
| Espaços de magia usados/disponíveis | ❌ não existe | — | — |
| Magia ritual | ❌ não existe | — | não crítico, como o checklist original já antecipava |
| Foco de conjuração / bolsa de componentes | ⚠️ parcial via loja genérica | `Item`/`CharacterItem` (sistema de loja desta sessão) pode incluir um foco arcano como item comprável, mas não há nenhuma lógica que verifique "esse personagem tem um foco equipado" pra permitir conjurar | — |
| **Confirmação direto da fonte:** `ai-service/pdf_export/reference/dnd5e_pdf_field_map.json:251` já documenta isso: *"página 2 inteira (aparência/biografia) e página 3 inteira (conjuração de magias) não são preenchidas — o domain layer não modela spellcasting nem esses dados de biografia livre"* — quem construiu o export do PDF já sabia e registrou esse gap inteiro. Não é uma descoberta nova desta auditoria, é uma confirmação por uma segunda fonte independente. |

## 6. Antecedente e traços de personagem

| Item | Status | Onde | Observação |
|---|---|---|---|
| Traços/ideais/vínculos/defeitos | ✅ existe (bug de carregamento corrigido nesta sessão) | `Character.PersonalityTraits/Ideals/Bonds/Flaws`, mapeado em `dnd5e_pdf_field_map.json` | ver commit mais recente — `BackgroundForm.tsx` não recebia os valores salvos, corrigido |
| Perícias/ferramentas/idiomas do antecedente | ✅ existe | `Antecedent.SkillProficiencies/ToolProficiencies/Languages` | perícias estruturadas (lista), ferramentas/idiomas só texto livre |
| Equipamento inicial do antecedente | ⚠️ existe só como texto, não aplicado | `Antecedent.Equipment` (string livre) | o texto existe e aparece na tela (`CharacterDetail.tsx:712-716`), mas **nada credita esse equipamento automaticamente no inventário do personagem** — o sistema de loja/inventário (`CharacterItem`) desta sessão é 100% baseado em compra com moedas, não em "receber de graça pelo antecedente". O jogador teria que comprar manualmente os mesmos itens que o texto do antecedente já diz que ele deveria ter de graça |
| Traço de antecedente nomeado (feature) | ✅ existe | `Antecedent.Feature`/`FeatureDescription`, exibido em `CharacterDetail.tsx:697-698` | — |

## 7. Traços raciais

| Item | Status | Onde | Observação |
|---|---|---|---|
| Visão no escuro (Darkvision) | ❌ não existe | `Race` (`race.go`) não tem nenhum campo de visão | confirmado ainda não implementado, como o checklist original já esperava (roadmap Fase F) |
| Resistências/imunidades | ❌ não existe | idem | idem |
| Tamanho | ❌ não existe | `Race` não tem campo `Size` | — |
| Magias raciais inatas | ❌ não existe | idem | roadmap Fase G, como já esperado |
| Idiomas raciais | ❌ não existe | idem | roadmap Fase E, como já esperado |

## 8. Equipamento e recursos

| Item | Status | Onde | Observação |
|---|---|---|---|
| Inventário de equipamento | ✅ **existe agora** (o checklist original está desatualizado neste ponto) | `Item`/`CharacterItem`/`CharacterArmorOwned`/`InventoryService` — construído numa sessão posterior à criação deste checklist | o checklist original (Seção 8) foi escrito antes do sistema de loja existir e assumia que isso ainda era só planejado ("Fase D"). Hoje existe de verdade: catálogo de ~38 armas + ~92 itens + 13 armaduras + ~22 itens mágicos, compra com as 5 moedas oficiais. Ver seção "Shop / equipment / currency (5e)" do `CLAUDE.md` |
| Moedas (PC/PP/PE/PO/PL) | ✅ existe | `Character.CopperPieces...PlatinumPieces` (`character.go:67-71`) | — |
| Capacidade de carga (Força × 15) | ❌ não existe | nenhum cálculo encontrado | não crítico, como o checklist original já marcava |
| Armas equipadas → bloco de Ataques | ❌ não existe | depende do cálculo de ataque da Seção 4, que não existe | o inventário sabe "quais itens o personagem tem", mas nada os liga a um cálculo de ataque (que também não existe) |
| Armadura equipada (`ArmorID`) | ✅ existe, nullable | `Character.ArmorID *uint` (`character.go:27`) | — |

## 9. Descrição física e narrativa

| Item | Status | Onde | Observação |
|---|---|---|---|
| Idade, altura, peso, olhos, pele, cabelo | ❌ não existe | nenhum campo em `Character` | confirma o comentário do próprio field map da IA citado na Seção 5 — página de aparência do PDF não é preenchida por não ter de onde puxar o dado |
| Símbolo de fé/divindade | ❌ não existe | idem | — |
| História do personagem | ⚠️ existe só decorativamente | campo "História" em `BackgroundForm.tsx` | **nunca foi persistido** — `BackgroundService.Save`'s allowlist (`background_service.go`) nunca incluiu `"history"`; o campo aparece no formulário, mas escrever nele e salvar não guarda nada. Bug pré-existente, não introduzido nesta sessão, e ainda não corrigido (fora do escopo do que foi pedido até agora) |
| Aliados/organizações, tesouro (texto livre) | ❌ não existe | nenhum campo em `Character` | — |

---

## Classificação final (pra decisão em conjunto)

**Bloqueia a Fase 2 (export PDF) por já causar dado ERRADO exportado, não só ausente:**
- ~~`Character.Speed` nunca é preenchido na criação → PDF exporta deslocamento "0" pra todo personagem 5e (Seção 4).~~ **Corrigido.** `Create()` agora carrega `domain.Race` (mesmo padrão do fix de `Class` já existente) e copia `race.Speed` pra `character.Speed` em personagens 5e. Verificado via API real: `Humano` → 30, `Elfo da Floresta` → 35, valor persiste após reload. Ver `character_service.go` e `CLAUDE.md`.
- Expertise/Jack of All Trades ausentes (Seção 3) fazem a ficha exportar o bônus de perícia errado pra qualquer Ladino/Bardo — não trava a exportação, mas o número que sai está incorreto pras regras do personagem.

**Aceitável por ora (campo ausente = PDF exporta em branco/zero, sem mentir um valor errado):**
- Conjuração inteira (Seção 5) — já assumido e documentado no próprio field map da IA.
- Bônus de ataque/dano de arma, Ataque Extra (Seção 4).
- Idiomas estruturados, proficiências de arma/armadura estruturadas (Seção 3).
- Capacidade de carga (Seção 8).
- Descrição física/narrativa (Seção 9), exceto o bug do campo "História" que nunca salva (esse é um bug de UX enganoso — o jogador acha que salvou e não salvou — mas não afeta o PDF, que também nunca teve esse campo mapeado).

**Trabalho de fases futuras já no roadmap, não antecipar:**
- Traços raciais inteiros — Darkvision, resistências, tamanho, magias raciais, idiomas raciais (Seção 7, Fases E/F/G).
- Equipamento inicial do antecedente creditado automaticamente (Seção 6) — depende de decisão de produto sobre como isso interage com o sistema de compra já existente.

**Decisões de produto, não bugs — perguntar antes de mexer:**
- ~~Bônus de atributo de raça/antecedente ausente (Seção 0/2)~~ **Respondido e implementado.** Bônus de atributo vem sempre do Antecedente (nunca da raça/espécie, seja ela 2014 ou 2024) — `Antecedent.AbilityBonusOptions` + `OriginFeatName`, escolha do jogador em `Character.AbilityBonusChoice`, aplicado e validado em `CharacterService.Create`/`applyAntecedentAbilityBonus`, com UI dedicada em `CharacterCreate.tsx`. Testado via API (+2/+1, +1/+1/+1, rejeição de escolha inválida, rejeição de ausência de escolha) e via UI (Playwright). Detalhes completos na seção "Mixed 2014/2024 ruleset" do `CLAUDE.md`. Pendência isolada: o antecedente "Herói do Povo" é conteúdo só de 2014 sem equivalente direto no PHB 2024, então ficou sem bônus automático — decidir se remapeia pro antecedente 2024 mais parecido (Fazendeiro) ou mantém assim.
- ~~ASI vs. Feat na mesma escolha (Seção 2)~~ **Respondido: impedir. Implementado e testado.** `ApplyASI` agora rejeita (`400`) se atributo e `talento_id` vierem juntos na mesma chamada; escolher `talento_id` concede o talento e não mexe em nenhum atributo. UI: aba "Atributo"/"Talento" no modal de melhoria (só aparece pra 5e). **Gap descoberto ao testar isso:** não existe nenhum `Talento` com `Edition: "5e"` no banco — o seed de talentos é 100% 4e. A troca por talento funciona mecanicamente (testado com um talento 4e só pra provar o mecanismo), mas não há nenhum talento 5e de verdade pra um jogador escolher hoje. Isso também significa que o painel de Talentos na ficha estava escondido pra 5e por engano (`is4e && hasTalentos` → corrigido pra `hasTalentos`).
- ~~Subclasses inconsistentes entre classes~~ **Respondido: corrigir níveis e adicionar as que faltavam. Implementado e testado.** As 9 classes sem subclasse (Bárbaro, Bardo, Druida, Guardião, Guerreiro, Ladino, Mago, Monge, Paladino) agora têm as 4 opções reais do PHB 2024 cada, todas no **nível 3** (nome e nível confirmados direto no PDF via `rag_5e.query_5e()`, não por memória — a 2024 unificou quase toda escolha de subclasse pro nível 3, inclusive Mago e Druida que eram nível 2 no 2014). As 3 que já existiam (Clérigo/Bruxo/Feiticeiro, nível 1) conferem com o PHB 2024 e não precisaram de correção. **Bug de dados encontrado nesse processo, não corrigido ainda:** existem 2 linhas "Guerreiro" 5e duplicadas no banco (uma completa, uma incompleta sem saving throws/perícias) — não mexi nisso sem confirmar com você, já que apagar/mesclar uma linha de classe é uma ação destrutiva se algum personagem já referenciar a errada.
- `domain.Background` parece uma tabela morta (Seção 1) — **respondido: deixar por ora**, revisitar se precisar mais adiante.

## Pergunta em aberto (bloqueando a implementação de 2014+2024)

Você pediu pra formalizar **tanto** as regras de 2014 quanto as de 2024 pro bônus de atributo — mas isso são dois mecanismos diferentes e incompatíveis pro mesmo personagem (2014: raça dá o bônus; 2024: antecedente dá o bônus + 1 talento de origem). Antes de implementar isso, preciso saber: **como o sistema decide qual das duas regras vale pra um personagem específico?** Por exemplo:
- Um campo novo no personagem (`rules_version`: "2014" ou "2024") escolhido na criação, e a UI muda o que aparece (bônus vem da raça OU do antecedente) conforme essa escolha?
- Ou é uma configuração global do sistema/tabela (todo mundo usa a mesma regra)?
- Ou outra coisa que você tinha em mente?
