package seed

import (
	"gorm.io/gorm"
	"log"
	"rpg-manager/internal/domain"
)

// seedTalentos5e semeia os talentos (feats) do PHB 2024, extraídos diretamente
// do PDF (capítulo 5, "D&D 5.0 - livro-do-jogador-2024.pdf") via PyMuPDF -- não
// inventados. Cobre as 4 categorias do livro: Origem, Geral, Estilo de Luta e
// Dádiva Épica (75 talentos ao todo, conferido 1 a 1 contra a tabela oficial do
// capítulo). O pré-requisito de nível/atributo vai no campo Prerequisite como
// texto (não há enforcement automático hoje, mesmo padrão dos talentos 4e).
func seedTalentos5e(db *gorm.DB) {
	talentos := []domain.Talento{
		{
			Name: "Alerta", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Proficiência em Iniciativa: quando você joga Iniciativa, pode adicionar seu Bônus de Proficiência à jogada. Troca de Iniciativa: imediatamente após jogar Iniciativa, você pode trocar sua Iniciativa com a de um aliado voluntário no mesmo combate (não é possível se você ou o aliado estiver Incapacitado).",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Artifista", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Proficiência com Ferramentas: adquire proficiência com três Ferramentas de Artesão diferentes à sua escolha. Desconto: sempre que compra um item não mágico, recebe 20% de desconto. Fabricação Rápida: ao completar um Descanso Longo, pode fabricar uma peça de equipamento da tabela Fabricação Rápida se tiver as ferramentas associadas; o item se desfaz no próximo Descanso Longo.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Atacante Selvagem", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você treinou para causar ataques particularmente nocivos. Uma vez por turno, quando atinge um alvo com uma arma, pode jogar os dados de dano da arma duas vezes e usar qualquer uma das jogadas contra o alvo.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Curandeiro", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Médico de Combate: com um Kit de Curandeiro, pode gastar um uso e cuidar de uma criatura a até 1,5 metro como uma ação Usar Objeto — ela gasta um Dado de Vida e você joga o dado; ela recupera PV iguais à jogada mais seu Bônus de Proficiência. Cura Garantida: sempre que jogar um dado para determinar PV recuperados por magia ou pelo benefício Médico de Combate, pode jogar novamente se tirar 1, usando a nova jogada.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Habilidoso", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire proficiência em qualquer combinação de três perícias ou ferramentas à sua escolha. Repetível: você pode adquirir este talento mais de uma vez.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Iniciado em Magia", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Dois Truques: aprende dois truques à escolha da lista de Clérigo, Druida ou Mago (Inteligência, Sabedoria ou Carisma é seu atributo de conjuração, escolhido ao selecionar o talento). Magia de 1º Círculo: escolhe uma magia de 1º círculo da mesma lista, sempre preparada, conjurável uma vez sem espaço de magia (restaura em Descanso Longo) ou com qualquer espaço de magia que tiver. Substituição de Magia: a cada novo nível pode trocar uma das magias escolhidas por outra do mesmo círculo da mesma lista. Repetível: pode adquirir de novo, mas deve escolher uma lista de magias diferente a cada vez.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Músico", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Treinamento em Instrumentos: proficiência com três Instrumentos Musicais à escolha. Canção Encorajadora: ao completar um Descanso Curto ou Longo, pode tocar uma música e conceder Inspiração Heroica a um número de aliados que ouvem a música igual ao seu Bônus de Proficiência.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Sortudo", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Pontos de Sorte: você tem Pontos de Sorte iguais ao seu Bônus de Proficiência (restauram em Descanso Longo). Vantagem: gaste 1 Ponto de Sorte para ter Vantagem numa jogada de d20 que você fizer. Desvantagem: gaste 1 Ponto de Sorte para impor Desvantagem numa jogada de ataque contra você.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Valentão de Taverna", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Você adquire os seguintes benefícios. Ataque Desarmado Aprimorado: seu Ataque Desarmado causa 1d4 + mod. Força de dano Contundente. Dano Garantido: pode jogar novamente o dado de dano do Ataque Desarmado se tirar 1, usando a nova jogada. Armamento Improvisado: proficiência com armas improvisadas. Empurrar: ao atingir com Ataque Desarmado como parte da ação Atacar, pode empurrar o alvo 1,5 metro (1x por turno).",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Vigoroso", Edition: "5e", Category: "Origem", Prerequisite: "",
			Description: "Seus Pontos de Vida máximos aumentam em um valor igual ao dobro do seu nível de personagem quando você obtém este talento. Sempre que alcançar um novo nível depois disso, seus PV máximos aumentam em mais 2.",
			Tooltip:     "Talento de Origem — normalmente concedido pelo antecedente escolhido na criação.",
		},

		{
			Name: "Adepto Elemental", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Característica Conjuração ou Magia de Pacto — repetível (escolhendo um tipo de dano diferente a cada vez)",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Domínio Elemental: escolha Ácido, Elétrico, Gélido, Ígneo ou Trovejante — suas magias ignoram Resistência a esse dano, e ao jogar dano desse tipo pode tratar qualquer 1 como um 2.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Característica Conjuração ou Magia de Pacto — repetível (escolhendo um tipo de dano diferente a cada vez)).",
		},

		{
			Name: "Agressor", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força ou Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Corrida Aprimorada: ao Correr, Deslocamento +3m nessa ação. Ataque em Investida: movendo-se ao menos 3m em linha reta antes de atingir um alvo corpo a corpo, escolha +1d8 de dano ou empurrar o alvo até 3m (1x por turno).",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força ou Destreza 13 ou superior).",
		},

		{
			Name: "Analítico", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Inteligência ou Sabedoria 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência ou Sabedoria +1 (máx. 20). Observador Atento: escolha Intuição, Investigação ou Percepção — ganha proficiência, ou Especialização se já for proficiente. Pesquisa Rápida: pode executar a ação Procurar como Ação Bônus.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Inteligência ou Sabedoria 13 ou superior).",
		},

		{
			Name: "Atirador Arcano", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Característica Conjuração ou Magia de Pacto",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Ignorar Cobertura: ataques de magia ignoram Cobertura Parcial e de Três Quartos. Conjuração à Queima-Roupa: inimigo a 1,5m não impõe Desvantagem em ataques de magia. Alcance Aumentado: magias de ataque com alcance de ao menos 3m ganham +18m de alcance.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Característica Conjuração ou Magia de Pacto).",
		},

		{
			Name: "Atleta", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força ou Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Deslocamento de Escalada igual ao seu Deslocamento. Levantar-se do Caído com apenas 1,5m de movimento. Salto em Distância ou Altura correndo após mover-se só 1,5m.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força ou Destreza 13 ou superior).",
		},

		{
			Name: "Ator", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Carisma 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Carisma +1 (máx. 20). Personificação: Vantagem em Carisma (Atuação ou Enganação) ao se passar por alguém real ou fictício. Mimetismo: imita sons e falas de outras criaturas (quem ouve precisa de teste de Sabedoria (Intuição) CD 8 + mod. Carisma + Bônus de Proficiência pra perceber o embuste).",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Carisma 13 ou superior).",
		},

		{
			Name: "Aumento no Valor de Atributo", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior — repetível",
			Description: "Aumente um valor de atributo à sua escolha em 2, ou dois valores de atributo à sua escolha em 1 cada. Não pode aumentar um atributo acima de 20.",
			Tooltip:     "Talento Geral (Nível 4 ou superior — repetível).",
		},

		{
			Name: "Chef", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Constituição ou Sabedoria +1 (máx. 20). Utensílios de Cozinheiro: proficiência, se ainda não tiver. Refeição Satisfatória: em Descanso Curto, cozinha comida pra 4 + Bônus de Proficiência criaturas; quem come e gasta Dados de Vida recupera 1d8 PV adicionais. Guloseimas Revigorantes: com 1h de trabalho ou em Descanso Longo, prepara guloseimas (iguais ao Bônus de Proficiência) que duram 8h — comer uma como Ação Bônus dá PV temporários iguais ao seu Bônus de Proficiência.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Combatente Montado", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força, Destreza ou Sabedoria +1 (máx. 20). Golpe Montado: Vantagem em ataques contra criaturas desmontadas a até 1,5m da sua montaria, pelo menos um tamanho menor que ela. Pulo Lateral: sua montaria evita totalmente dano de salvaguarda de Destreza bem-sucedida (metade se falhar), desde que você a esteja montando e nenhum dos dois esteja Incapacitado. Redirecionar Ataque: enquanto montado, força um ataque que atingiria sua montaria a te atingir, se não estiver Incapacitado.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Conjurador Bélico", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Característica Conjuração ou Magia de Pacto",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Concentração: Vantagem em salvaguardas de Constituição pra manter Concentração. Magia Reativa: ao provocar Ataque de Oportunidade por sair do alcance, pode usar Reação pra conjurar uma magia (tempo de conjuração de 1 ação, alvo único) em vez de sofrer o ataque. Componentes Somáticos: pode realizá-los mesmo com armas/Escudo nas mãos.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Característica Conjuração ou Magia de Pacto).",
		},

		{
			Name: "Conjurador Ritualista", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Inteligência, Sabedoria ou Carisma 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Magias Rituais: escolhe magias de 1º círculo com o marcador Ritual, em número igual ao seu Bônus de Proficiência, sempre preparadas e conjuráveis com qualquer espaço de magia; ganha mais uma sempre que o Bônus de Proficiência aumentar. Ritual Rápido: conjura uma dessas magias com o tempo normal em vez do tempo de Ritual, sem gastar espaço de magia, uma vez por Descanso Longo.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Inteligência, Sabedoria ou Carisma 13 ou superior).",
		},

		{
			Name: "Duelista Defensivo", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Destreza +1 (máx. 20). Aparar: segurando uma arma de Acuidade, ao ser atingido por ataque corpo a corpo pode usar Reação para somar seu Bônus de Proficiência à CA contra esse ataque, possivelmente fazendo-o errar (bônus dura até o início do seu próximo turno).",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Destreza 13 ou superior).",
		},

		{
			Name: "Envenenador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Destreza ou Inteligência +1 (máx. 20). Veneno Potente: ignora Resistência a dano Venenoso nas suas jogadas de dano. Preparar Veneno: proficiência com Kit de Veneno; com 1h e 50 PO em materiais fabrica doses iguais ao Bônus de Proficiência; aplicar é Ação Bônus, dura 1 minuto ou até causar dano; alvo faz salvaguarda de Constituição (CD 8 + mod. do atributo aumentado + Bônus de Proficiência) ou sofre 2d8 de dano Venenoso e fica Envenenado até o fim do próximo turno.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Esmagador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Constituição +1 (máx. 20). Empurrar: 1x por turno, ao causar dano Contundente pode mover o alvo 1,5m para um espaço desocupado (se não for maior que você). Crítico Melhorado: acerto crítico com dano Contundente dá Vantagem contra esse alvo até seu próximo turno.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Especialista Ambidestro", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força ou Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Combate com Duas Armas Aprimorado: ao Atacar com arma Leve, faz um ataque adicional como Ação Bônus com outra arma corpo a corpo sem a propriedade Duas Mãos (sem somar modificador de atributo ao dano, salvo se negativo). Saque Rápido: desembainha/embainha duas armas sem Duas Mãos quando normalmente faria só uma.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força ou Destreza 13 ou superior).",
		},

		{
			Name: "Especialista em Armaduras Leves", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Treinamento com Armadura Leve e Escudos.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Especialista em Armaduras Médias", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Treinamento com Armadura Leve",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Treinamento com Armadura Média.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Treinamento com Armadura Leve).",
		},

		{
			Name: "Especialista em Armaduras Pesadas", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Treinamento com Armadura Média",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Constituição ou Força +1 (máx. 20). Treinamento com Armadura Pesada.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Treinamento com Armadura Média).",
		},

		{
			Name: "Especialista em Besta", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Destreza +1 (máx. 20). Ignorar Recarga das bestas de Mão/Leve/Pesada (pode recarregar sem mão livre). Disparo à Queima-Roupa: inimigo a 1,5m não impõe Desvantagem com bestas. Combate com Duas Armas: soma modificador de atributo ao ataque adicional da propriedade Leve se for com besta Leve.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Destreza 13 ou superior).",
		},

		{
			Name: "Especialista em Perícia", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 20). Proficiência em Perícia: uma perícia à escolha. Especialização: escolha uma perícia em que já é proficiente (mas sem Especialização) e ganhe Especialização nela.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Exterminador de Conjuradores", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Quebrador de Concentração: dano causado por você impõe Desvantagem na salvaguarda de Concentração do alvo. Resguardo Mental: ao falhar salvaguarda de Inteligência, Sabedoria ou Carisma, pode escolher ser bem-sucedido (1x por Descanso Curto ou Longo).",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Imobilizador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força ou Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Socar e Imobilizar: 1x por turno ao acertar Ataque Desarmado, pode causar dano e Imobilizar. Vantagem no Ataque contra criatura Imobilizada por você. Imobilizador Veloz: não gasta movimento extra movendo-se com uma criatura Imobilizada do seu tamanho ou menor.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força ou Destreza 13 ou superior).",
		},

		{
			Name: "Líder Inspirador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Sabedoria ou Carisma 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Sabedoria ou Carisma +1 (máx. 20). Atuação Encorajadora: em Descanso Curto ou Longo, faz uma atuação e concede a até 6 aliados a 9m PV temporários iguais ao seu nível + mod. do atributo aumentado.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Sabedoria ou Carisma 13 ou superior).",
		},

		{
			Name: "Mente Aguçada", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Inteligência 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência +1 (máx. 20). Conhecimento Vasto: escolha Arcanismo, História, Investigação, Natureza ou Religião — ganha proficiência, ou Especialização se já proficiente. Análise Rápida: pode executar a ação Analisar como Ação Bônus.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Inteligência 13 ou superior).",
		},

		{
			Name: "Mestre das Armas", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Propriedade de Maestria: usa a maestria de um tipo de arma Simples ou Marcial em que seja proficiente à escolha, podendo trocar em Descanso Longo.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Mestre em Armaduras Médias", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Treinamento com Armadura Média",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Portador Ágil: com armadura Média e Destreza 16+, soma +3 (em vez de +2) à CA pela Destreza.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Treinamento com Armadura Média).",
		},

		{
			Name: "Mestre em Armaduras Pesadas", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Treinamento com Armadura Pesada",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Constituição +1 (máx. 20). Redução de Dano: dano Contundente, Cortante e Perfurante sofrido com armadura Pesada é reduzido em pontos iguais ao seu Bônus de Proficiência.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Treinamento com Armadura Pesada).",
		},

		{
			Name: "Mestre em Armas de Haste", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força ou Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Golpe de Haste: com Cajado, Lança ou arma Extensão+Pesado, ataca com Ação Bônus usando a extremidade oposta (dano Contundente 1d4). Golpe Reativo: com essas armas, usa Reação para atacar criatura que entra no seu alcance.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força ou Destreza 13 ou superior).",
		},

		{
			Name: "Mestre em Armas Grandes", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força +1 (máx. 20). Maestria em Armas Pesadas: dano adicional igual ao Bônus de Proficiência com armas de propriedade Pesada. Cortar: após Acerto Crítico ou reduzir criatura a 0 PV com arma corpo a corpo, ataca de novo com a mesma arma como Ação Bônus.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força 13 ou superior).",
		},

		{
			Name: "Mestre em Escudos", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Treinamento com Escudo",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força +1 (máx. 20). Golpe de Escudo: ao acertar com arma corpo a corpo, ataca com o Escudo forçando salvaguarda de Força (CD 8 + mod. Força + Bônus de Proficiência) — se falhar, empurra 1,5m ou derruba (1x por turno). Interpor Escudo: com Reação e Escudo equipado, evita totalmente dano de salvaguarda de Destreza bem-sucedida.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Treinamento com Escudo).",
		},

		{
			Name: "Mestre-Atirador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Destreza +1 (máx. 20). Ignorar Cobertura Parcial/Três Quartos em ataques à distância. Disparo à Queima-Roupa: inimigo a 1,5m não impõe Desvantagem. Tiro Longo: alcance máximo não impõe Desvantagem.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Destreza 13 ou superior).",
		},

		{
			Name: "Perfurador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Punção: 1x por turno, ao causar dano Perfurante pode jogar novamente um dos dados de dano, usando a nova jogada. Crítico Melhorado: em Acerto Crítico com dano Perfurante, joga um dado de dano adicional.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Resiliente", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: escolha um atributo sem proficiência em salvaguarda e aumente-o em 1 (máx. 20). Proficiência em Salvaguarda com o atributo escolhido.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Resistente", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Constituição +1 (máx. 20). Desafie a Morte: Vantagem em Salvaguardas Contra Morte. Recuperação Rápida: como Ação Bônus, gasta um Dado de Vida e recupera PV iguais ao resultado.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Sentinela", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Força ou Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Diligente: Ataque de Oportunidade contra criatura a 1,5m que Desengaje ou ataque outro alvo. Deter: ao acertar Ataque de Oportunidade, Deslocamento da criatura vira 0 pelo resto do turno.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Força ou Destreza 13 ou superior).",
		},

		{
			Name: "Sorrateiro", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Destreza 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Destreza +1 (máx. 20). Visão às Cegas com alcance de 3m. Névoa de Guerra: Vantagem em Destreza (Furtividade) ao Esconder-se durante combate. Atirador: errar um ataque enquanto escondido não revela sua localização.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Destreza 13 ou superior).",
		},

		{
			Name: "Talhador", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Debilitar: 1x por turno, ao causar dano Cortante reduz Deslocamento do alvo em 3m até seu próximo turno. Crítico Melhorado: Acerto Crítico com dano Cortante impõe Desvantagem nos ataques do alvo até seu próximo turno.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Telecinético", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Telecinese Menor: aprende Mãos Mágicas (sem componentes V/S, mão Invisível, alcance +9m). Empurrão Telecinético: Ação Bônus empurra criatura a até 9m (salvaguarda de Força CD 8 + mod. + Bônus de Proficiência ou move 1,5m).",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Telepático", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Enunciado Telepático: fala telepaticamente até 18m em idioma que conhece. Detectar Pensamentos sempre preparada, conjurável 1x sem espaço de magia por Descanso Longo (ou com espaço de magia depois).",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Tocado Pelas Sombras", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Sua exposição à magia do Sombral lhe concede os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Magia Sombria: escolhe uma magia de 1º círculo de Ilusão ou Necromancia — ela e Invisibilidade ficam sempre preparadas, conjuráveis sem espaço de magia 1x por Descanso Longo cada (ou com espaço de magia depois).",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Tocado Por Fadas", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Sua exposição à magia Feérica lhe concede os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 20). Magia Feérica: escolhe uma magia de 1º círculo de Adivinhação ou Encantamento — ela e Passo Nebuloso ficam sempre preparadas, conjuráveis sem espaço de magia 1x por Descanso Longo cada (ou com espaço de magia depois).",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Treinamento com Armas Marciais", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 20). Proficiência com armas Marciais.",
			Tooltip:     "Talento Geral (Nível 4 ou superior).",
		},

		{
			Name: "Velocista", Edition: "5e", Category: "Geral", Prerequisite: "Nível 4 ou superior, Destreza ou Constituição 13 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Destreza ou Constituição +1 (máx. 20). Aumento de Deslocamento: +3m. Correr em Terreno Difícil sem custo adicional de movimento. Movimentação Ágil: Ataques de Oportunidade contra você têm Desvantagem.",
			Tooltip:     "Talento Geral (Nível 4 ou superior, Destreza ou Constituição 13 ou superior).",
		},

		{
			Name: "Arquearia", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Você recebe um bônus de +2 nas jogadas de ataque com armas à Distância.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Combate com Armas de Arremesso", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Ao acertar um ataque à distância com arma de propriedade Arremesso, ganha +2 na jogada de dano.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Combate com Armas Grandes", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Ao jogar dano com arma corpo a corpo empunhada com as duas mãos (propriedade Duas Mãos ou Versátil), pode tratar 1 ou 2 no dado como 3.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Combate com Duas Armas", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "No ataque adicional de uma arma com propriedade Leve, pode somar o modificador de atributo ao dano, se ainda não estiver somando.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Combate Desarmado", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Seu Ataque Desarmado causa 1d6 (1d8 se não estiver segurando arma ou Escudo) + mod. Força de dano Contundente. No início de cada turno, pode causar 1d4 de dano Contundente a uma criatura Imobilizada por você.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Defensivo", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Usando armadura Leve, Média ou Pesada, ganha +1 na Classe de Armadura.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Duelismo", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Empunhando uma arma corpo a corpo em uma mão e nenhuma outra arma, ganha +2 nas jogadas de dano dessa arma.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Interceptação", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Quando uma criatura à sua vista atinge outra a até 1,5m de você, pode usar Reação para reduzir o dano em 1d10 + Bônus de Proficiência (precisa estar com Escudo ou arma Simples/Marcial).",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Luta às Cegas", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Você tem Visão às Cegas com alcance de 3 metros.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Protetivo", Edition: "5e", Category: "Estilo de Luta", Prerequisite: "Característica de Estilo de Luta",
			Description: "Quando uma criatura à sua vista ataca um alvo (que não você) a até 1,5m de distância, pode usar Reação pra interpor o Escudo, impondo Desvantagem no ataque e em ataques contra o alvo até seu próximo turno, enquanto estiver a até 1,5m dele.",
			Tooltip:     "Talento de Estilo de Luta — requer a característica Estilo de Luta da classe.",
		},

		{
			Name: "Dádiva da Fortitude", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Saúde Fortalecida: PV máximos +40; ao recuperar PV, recupera PV adicionais iguais ao mod. Constituição (1x até seu próximo turno).",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Proeza em Combate", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Pontaria Inigualável: ao errar uma jogada de ataque, pode transformá-la em acerto (1x até seu próximo turno).",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Proficiência em Perícia", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Assecla Completo: proficiência em todas as perícias. Especialização em uma perícia na qual já seja proficiente.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Recordação de Magia", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Inteligência, Sabedoria ou Carisma +1 (máx. 30). Conjuração Livre: ao conjurar magia com espaço de 1º a 4º círculo, jogue 1d4 — se o resultado for igual ao círculo do espaço, ele não é gasto.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Recuperação", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Até a Morte: ao ser reduzido a 0 PV, pode ficar com 1 PV e recuperar metade dos PV máximos (1x por Descanso Longo). Recuperar Vitalidade: reserva de dez d10 — como Ação Bônus, gasta dados da reserva e recupera PV iguais ao total (restaura em Descanso Longo).",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Resistência à Energia", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Resistências à Energia: Resistência a dois tipos de dano à escolha entre Ácido, Elétrico, Gélido, Ígneo, Necrótico, Psíquico, Radiante ou Venenoso (troca em Descanso Longo). Redirecionamento de Energia: com Reação, redireciona dano de um desses tipos pra outra criatura à vista a até 18m (salvaguarda de Destreza CD 8 + mod. Constituição + Bônus de Proficiência ou sofre 2d12 + mod. Constituição).",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Velocidade", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Artista de Fuga: Ação Bônus para Desengajar, encerrando também a condição Imobilizado. Agilidade: Deslocamento +9m.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Viagem Dimensional", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Passos Fugazes: imediatamente após Atacar ou Usar Magia, teleporta-se até 9m para um espaço desocupado à vista.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva da Visão Verdadeira", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Visão Verdadeira com alcance de 18 metros.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva do Ataque Irresistível", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: Força ou Destreza +1 (máx. 30). Superar Defesas: dano Contundente, Cortante e Perfurante que você causa ignora Resistência. Golpe Devastador: ao tirar 20 natural num ataque, causa dano adicional igual ao valor do atributo aumentado por este talento.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva do Destino", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Aprimorar Destino: quando você ou outra criatura a até 18m for bem-sucedida ou falhar num Teste de D20, pode jogar 2d4 e aplicar como bônus ou penalidade (1x até jogar Iniciativa ou Descanso Curto/Longo).",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},

		{
			Name: "Dádiva do Espírito da Noite", Edition: "5e", Category: "Dádiva Épica", Prerequisite: "Nível 19 ou superior",
			Description: "Você adquire os seguintes benefícios. Aumento no Valor de Atributo: um atributo à escolha +1 (máx. 30). Fundir-se com Sombras: em Meia-luz ou Escuridão, concede-se Invisível como Ação Bônus (encerra ao agir). Forma Sombria: em Meia-luz ou Escuridão, Resistência a todo dano exceto Psíquico e Radiante.",
			Tooltip:     "Dádiva Épica — só disponível a partir do nível 19.",
		},
	}

	for _, t := range talentos {
		var existing domain.Talento
		if db.Where("name = ? AND edition = ?", t.Name, t.Edition).First(&existing).Error != nil {
			db.Create(&t)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"description":  t.Description,
				"prerequisite": t.Prerequisite,
				"category":     t.Category,
				"tooltip":      t.Tooltip,
			})
		}
	}
	log.Println("  ✓ Talentos 5e (PHB 2024) seedados")
}
