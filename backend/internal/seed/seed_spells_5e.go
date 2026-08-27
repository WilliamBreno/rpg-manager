package seed

import (
	"log"

	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

// seedSpells5e povoa o catálogo de magias do PHB 2024 (capítulo "Magias"),
// extraído das listas de magias de cada classe conjuradora (capítulo 3,
// seções "Lista de Magias de X") via PyMuPDF — não inventado. Cobre as 397
// magias únicas conhecidas por Bardo, Bruxo, Clérigo, Druida, Feiticeiro,
// Guardião, Mago e Paladino (Guerreiro e Ladino usam a lista de Mago via suas
// subclasses conjuradoras — Cavaleiro Místico e Trapaceiro Arcano — e não têm
// lista própria; Monge não conjura). Descrição mecânica completa (dano,
// alcance detalhado etc.) não foi seedada nesta primeira passada — ver nota
// em CLAUDE.md — só nome, círculo, escola, ritual/concentração e quais
// classes a conhecem em que círculo, o suficiente para a tela de escolha na
// criação de personagem e para exibir a lista de magias conhecidas.
func seedSpells5e(db *gorm.DB) {
	spells := []domain.Spell{
		{
			Name: "Acalmar Emoções", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Clérigo": 2}`,
		},
		{
			Name: "Acudir os Moribundos", Edition: "5e", Level: 0, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 0, "Druida": 0}`,
		},
		{
			Name: "Alarme", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: true, Concentration: false,
			Classes: `{"Guardião": 1, "Mago": 1}`,
		},
		{
			Name: "Aliado Extraplanar", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 6}`,
		},
		{
			Name: "Aljava Veloz", Edition: "5e", Level: 5, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Guardião": 5}`,
		},
		{
			Name: "Alterar-se", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Amigos", Edition: "5e", Level: 0, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Amizade Animal", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Druida": 1, "Guardião": 1}`,
		},
		{
			Name: "Animar Mortos", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 3, "Mago": 3}`,
		},
		{
			Name: "Animar Objetos", Edition: "5e", Level: 5, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Antipatia/Simpatia", Edition: "5e", Level: 8, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 8, "Druida": 8, "Mago": 8}`,
		},
		{
			Name: "Aprimorar Atributo", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Druida": 2, "Feiticeiro": 2, "Guardião": 2, "Mago": 2}`,
		},
		{
			Name: "Aprisionamento", Edition: "5e", Level: 9, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 9, "Mago": 9}`,
		},
		{
			Name: "Arca Secreta de Leomund", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Arma Elemental", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 3, "Guardião": 3, "Paladino": 3}`,
		},
		{
			Name: "Arma Espiritual", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 2}`,
		},
		{
			Name: "Arma Mágica", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 2, "Guardião": 2, "Mago": 2, "Paladino": 2}`,
		},
		{
			Name: "Armadura Arcana", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Armadura de Agathys", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 1}`,
		},
		{
			Name: "Arrombar", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Arte Druídica", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 0}`,
		},
		{
			Name: "Assassino Fantasmagórico", Edition: "5e", Level: 4, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4, "Mago": 4}`,
		},
		{
			Name: "Augúrio", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 2, "Druida": 2, "Mago": 2}`,
		},
		{
			Name: "Aumentar/Reduzir", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Druida": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Aura Mágica de Nystul", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 2}`,
		},
		{
			Name: "Aura Sagrada", Edition: "5e", Level: 8, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 8}`,
		},
		{
			Name: "Aura de Pureza", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 4, "Paladino": 4}`,
		},
		{
			Name: "Aura de Vida", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 4, "Paladino": 4}`,
		},
		{
			Name: "Aura de Vitalidade", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 3, "Druida": 3, "Paladino": 3}`,
		},
		{
			Name: "Auxílio", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Druida": 2, "Guardião": 2, "Paladino": 2}`,
		},
		{
			Name: "Badalar Fúnebre", Edition: "5e", Level: 0, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 0, "Clérigo": 0, "Mago": 0}`,
		},
		{
			Name: "Banimento", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 4, "Clérigo": 4, "Feiticeiro": 4, "Mago": 4, "Paladino": 4}`,
		},
		{
			Name: "Banquete de Heróis", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 6, "Clérigo": 6, "Druida": 6}`,
		},
		{
			Name: "Barreira de Lâminas", Edition: "5e", Level: 6, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 6}`,
		},
		{
			Name: "Benção", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 1, "Paladino": 1}`,
		},
		{
			Name: "Boca Encantada", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 2, "Mago": 2}`,
		},
		{
			Name: "Bola de Fogo", Edition: "5e", Level: 3, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Bola de Fogo Adiável", Edition: "5e", Level: 7, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Bolha Ácida", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Bom Fruto", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 1, "Guardião": 1}`,
		},
		{
			Name: "Bordão Místico", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 0}`,
		},
		{
			Name: "Braços de Hadar", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 1}`,
		},
		{
			Name: "Caldeirão Borbulhante de Tasha", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 6}`,
		},
		{
			Name: "Caldeirão Borbulhante de Tasha Invocação M Círculo da Morte", Edition: "5e", Level: 6, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 6}`,
		},
		{
			Name: "Caminhar Sobre as Águas", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 3, "Druida": 3, "Feiticeiro": 3, "Guardião": 3}`,
		},
		{
			Name: "Caminhar no Vento", Edition: "5e", Level: 6, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 6}`,
		},
		{
			Name: "Campo Antimagia", Edition: "5e", Level: 8, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 8, "Mago": 8}`,
		},
		{
			Name: "Cativar", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Bruxo": 2}`,
		},
		{
			Name: "Cegueira/Surdez", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Celeridade", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Chama Contínua", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 2, "Druida": 2, "Mago": 2}`,
		},
		{
			Name: "Chama Sagrada", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 0}`,
		},
		{
			Name: "Chicote de Espinhos", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 0}`,
		},
		{
			Name: "Chuva de Meteoros", Edition: "5e", Level: 9, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 9, "Mago": 9}`,
		},
		{
			Name: "Clarividência", Edition: "5e", Level: 3, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Clérigo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Clone", Edition: "5e", Level: 8, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 8}`,
		},
		{
			Name: "Coluna de Chamas", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 5}`,
		},
		{
			Name: "Comando", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Clérigo": 1, "Paladino": 1}`,
		},
		{
			Name: "Compreender Idiomas", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Compulsão", Edition: "5e", Level: 4, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4}`,
		},
		{
			Name: "Comunhão", Edition: "5e", Level: 5, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 5}`,
		},
		{
			Name: "Comunhão com a Natureza", Edition: "5e", Level: 5, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Druida": 5, "Guardião": 5}`,
		},
		{
			Name: "Cone de Frio", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Confusão", Edition: "5e", Level: 4, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4, "Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Consagrar", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 5}`,
		},
		{
			Name: "Contato Extraplanar", Edition: "5e", Level: 5, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Bruxo": 5, "Mago": 5}`,
		},
		{
			Name: "Contingência", Edition: "5e", Level: 6, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 6}`,
		},
		{
			Name: "Contramagia", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Controlar o Clima", Edition: "5e", Level: 8, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 8, "Druida": 8, "Mago": 8}`,
		},
		{
			Name: "Controlar Água", Edition: "5e", Level: 4, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 4, "Druida": 4, "Mago": 4}`,
		},
		{
			Name: "Contágio", Edition: "5e", Level: 5, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 5, "Druida": 5}`,
		},
		{
			Name: "Convocar Celestial", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 5, "Paladino": 5}`,
		},
		{
			Name: "Convocar Elemental", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Guardião": 4, "Mago": 4}`,
		},
		{
			Name: "Convocar Familiar", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: true, Concentration: false,
			Classes: `{"Mago": 1}`,
		},
		{
			Name: "Convocar Feérico", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 3, "Druida": 3, "Guardião": 3, "Mago": 3}`,
		},
		{
			Name: "Convocar Montaria", Edition: "5e", Level: 2, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 2}`,
		},
		{
			Name: "Convocar Relâmpagos", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 3}`,
		},
		{
			Name: "Corda Extradimensional", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 2}`,
		},
		{
			Name: "Cordão de Flechas", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Guardião": 2}`,
		},
		{
			Name: "Coroa da Loucura", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Corrente de Relâmpagos", Edition: "5e", Level: 6, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Crescer Espinhos", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Crescimento de Plantas", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Druida": 3, "Guardião": 3}`,
		},
		{
			Name: "Criar Chamas", Edition: "5e", Level: 0, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 0}`,
		},
		{
			Name: "Criar Comida e Água", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 3, "Paladino": 3}`,
		},
		{
			Name: "Criar Mortos-Vivos", Edition: "5e", Level: 6, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 6, "Clérigo": 6, "Mago": 6}`,
		},
		{
			Name: "Criar Passagem", Edition: "5e", Level: 5, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 5}`,
		},
		{
			Name: "Criar ou Destruir Água", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 1, "Druida": 1}`,
		},
		{
			Name: "Criação", Edition: "5e", Level: 5, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Cura Completa", Edition: "5e", Level: 6, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 6, "Druida": 6}`,
		},
		{
			Name: "Cura Completa em Massa", Edition: "5e", Level: 9, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 9}`,
		},
		{
			Name: "Curar Ferimentos", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Clérigo": 1, "Druida": 1, "Guardião": 1, "Paladino": 1}`,
		},
		{
			Name: "Curar Ferimentos em Massa", Edition: "5e", Level: 5, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Clérigo": 5, "Druida": 5}`,
		},
		{
			Name: "Cárcere de Energia", Edition: "5e", Level: 7, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 7, "Bruxo": 7, "Mago": 7}`,
		},
		{
			Name: "Cão Fiel de Mordenkainen", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Círculo Mágico", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 3, "Clérigo": 3, "Mago": 3, "Paladino": 3}`,
		},
		{
			Name: "Círculo da Morte", Edition: "5e", Level: 6, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 6, "Feiticeiro": 6}`,
		},
		{
			Name: "Círculo de Poder", Edition: "5e", Level: 5, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 5, "Mago": 5, "Paladino": 5}`,
		},
		{
			Name: "Círculo de Teleporte", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Cúpula Antivida", Edition: "5e", Level: 5, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 5}`,
		},
		{
			Name: "Danação", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 1}`,
		},
		{
			Name: "Dança Irresistível de Otto", Edition: "5e", Level: 6, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 6, "Mago": 6}`,
		},
		{
			Name: "De Carne para Pedra", Edition: "5e", Level: 6, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Dedo da Morte", Edition: "5e", Level: 7, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 7, "Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Defensor da Fé", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 4}`,
		},
		{
			Name: "Desejo", Edition: "5e", Level: 9, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 9, "Mago": 9}`,
		},
		{
			Name: "Desintegrar", Edition: "5e", Level: 6, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Despedaçar", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Despertar", Edition: "5e", Level: 5, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Druida": 5}`,
		},
		{
			Name: "Despistar", Edition: "5e", Level: 5, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Mago": 5}`,
		},
		{
			Name: "Destruição Atordoante", Edition: "5e", Level: 4, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 4}`,
		},
		{
			Name: "Destruição Banidora", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Paladino": 5}`,
		},
		{
			Name: "Destruição Cauterizante", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 1}`,
		},
		{
			Name: "Destruição Cegante", Edition: "5e", Level: 3, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 3}`,
		},
		{
			Name: "Destruição Colérica", Edition: "5e", Level: 1, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 1}`,
		},
		{
			Name: "Destruição Divina", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 1}`,
		},
		{
			Name: "Destruição Estrondosa", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 1}`,
		},
		{
			Name: "Destruição Radiante", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Paladino": 2}`,
		},
		{
			Name: "Detectar Magia", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: true, Concentration: true,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Clérigo": 1, "Druida": 1, "Feiticeiro": 1, "Guardião": 1, "Mago": 1, "Paladino": 1}`,
		},
		{
			Name: "Detectar Pensamentos", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Detectar Veneno e Doença", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: true, Concentration: true,
			Classes: `{"Clérigo": 1, "Druida": 1, "Guardião": 1, "Paladino": 1}`,
		},
		{
			Name: "Detectar o Bem e o Mal", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 1, "Paladino": 1}`,
		},
		{
			Name: "Disco Flutuante de Tenser", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: true, Concentration: false,
			Classes: `{"Mago": 1}`,
		},
		{
			Name: "Disfarçar-se", Edition: "5e", Level: 1, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Dissipar Magia", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Bruxo": 3, "Clérigo": 3, "Druida": 3, "Feiticeiro": 3, "Guardião": 3, "Mago": 3, "Paladino": 3}`,
		},
		{
			Name: "Dissipar o Bem e o Mal", Edition: "5e", Level: 5, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 5, "Paladino": 5}`,
		},
		{
			Name: "Dominar Fera", Edition: "5e", Level: 4, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Feiticeiro": 4, "Guardião": 4}`,
		},
		{
			Name: "Dominar Monstro", Edition: "5e", Level: 8, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 8, "Bruxo": 8, "Feiticeiro": 8, "Mago": 8}`,
		},
		{
			Name: "Dominar Pessoa", Edition: "5e", Level: 5, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Duelo Compelido", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Paladino": 1}`,
		},
		{
			Name: "Elementalismo", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Emaranhar", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 1, "Guardião": 1}`,
		},
		{
			Name: "Encarnação Fantasmagórica", Edition: "5e", Level: 9, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 9, "Mago": 9}`,
		},
		{
			Name: "Encontrar Armadilhas", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 2, "Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Encontrar o Caminho", Edition: "5e", Level: 6, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 6, "Clérigo": 6, "Druida": 6}`,
		},
		{
			Name: "Enfeitiçar Monstro", Edition: "5e", Level: 4, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 4, "Bruxo": 4, "Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Enfeitiçar Pessoa", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Druida": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Escalada de Aranha", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Escrita Ilusória", Edition: "5e", Level: 1, School: "Ilusão",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Mago": 1}`,
		},
		{
			Name: "Escudo Arcano", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Escudo Ardente", Edition: "5e", Level: 4, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Escudo da Fé", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 1, "Paladino": 1}`,
		},
		{
			Name: "Escuridão", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Esfera Congelante de Otiluke", Edition: "5e", Level: 6, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Esfera Flamejante", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Esfera Resiliente de Otiluke", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Esfera Vitriólica", Edition: "5e", Level: 4, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 4}`,
		},
		{
			Name: "Espada de Mordenkainen", Edition: "5e", Level: 7, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 7, "Mago": 7}`,
		},
		{
			Name: "Espinho Mental", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Esquentar Metal", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Druida": 2}`,
		},
		{
			Name: "Estática Sináptica", Edition: "5e", Level: 5, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Explosão Elemental", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 0}`,
		},
		{
			Name: "Explosão Solar", Edition: "5e", Level: 8, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 8, "Druida": 8, "Feiticeiro": 8, "Mago": 8}`,
		},
		{
			Name: "Fabricar", Edition: "5e", Level: 4, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Faca de Gelo", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Fagulha Estelar", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Druida": 0}`,
		},
		{
			Name: "Falar com Animais", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Druida": 1, "Guardião": 1}`,
		},
		{
			Name: "Falar com Mortos", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Clérigo": 3, "Mago": 3}`,
		},
		{
			Name: "Falar com Plantas", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Druida": 3, "Guardião": 3}`,
		},
		{
			Name: "Favor Divino", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 1}`,
		},
		{
			Name: "Flecha Relâmpago", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Guardião": 3}`,
		},
		{
			Name: "Flecha Ácida de Melf", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 2}`,
		},
		{
			Name: "Fogo das Fadas", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 1, "Druida": 1}`,
		},
		{
			Name: "Fome de Hadar Conjuração C Forma Gasosa", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 3}`,
		},
		{
			Name: "Fonte do Luar", Edition: "5e", Level: 4, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4, "Druida": 4}`,
		},
		{
			Name: "Forma Etérea", Edition: "5e", Level: 7, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Bruxo": 7, "Clérigo": 7, "Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Forma Gasosa", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Formas Animais", Edition: "5e", Level: 8, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 8}`,
		},
		{
			Name: "Força Espectral", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Gargalhada Nefasta de Tasha", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Mago": 1}`,
		},
		{
			Name: "Glifo de Proteção", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Clérigo": 3, "Mago": 3}`,
		},
		{
			Name: "Globo de Invulnerabilidade", Edition: "5e", Level: 6, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Golpe Certeiro", Edition: "5e", Level: 0, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Golpe Constritor", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Guardião": 1}`,
		},
		{
			Name: "Golpe de Arço", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Guardião": 5, "Mago": 5}`,
		},
		{
			Name: "Graxa", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Guardiões Espirituais", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 3}`,
		},
		{
			Name: "Heroísmo", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 1, "Paladino": 1}`,
		},
		{
			Name: "Identificar", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 1, "Mago": 1}`,
		},
		{
			Name: "Ilusão Menor", Edition: "5e", Level: 0, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Ilusão Programada", Edition: "5e", Level: 6, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 6, "Mago": 6}`,
		},
		{
			Name: "Imagem Maior", Edition: "5e", Level: 3, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Bruxo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Imagem Silenciosa", Edition: "5e", Level: 1, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Indetectável", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Guardião": 3, "Mago": 3}`,
		},
		{
			Name: "Infligir Ferimentos", Edition: "5e", Level: 1, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 1}`,
		},
		{
			Name: "Inseto Gigante", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4}`,
		},
		{
			Name: "Inverter a Gravidade", Edition: "5e", Level: 7, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 7, "Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Invisibilidade", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Invisibilidade Maior", Edition: "5e", Level: 4, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Invocar Aberração", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 4, "Mago": 4}`,
		},
		{
			Name: "Invocar Animais", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 3, "Guardião": 3}`,
		},
		{
			Name: "Invocar Barragem", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Guardião": 3}`,
		},
		{
			Name: "Invocar Celestial", Edition: "5e", Level: 7, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 7}`,
		},
		{
			Name: "Invocar Constructo", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Invocar Dragão", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 5}`,
		},
		{
			Name: "Invocar Elementais Menores", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Mago": 4}`,
		},
		{
			Name: "Invocar Elemental", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 5, "Mago": 5}`,
		},
		{
			Name: "Invocar Fera", Edition: "5e", Level: 2, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Invocar Feérico", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 6}`,
		},
		{
			Name: "Invocar Morto-vivo", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 3, "Mago": 3}`,
		},
		{
			Name: "Invocar Saraivada", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Guardião": 5}`,
		},
		{
			Name: "Invocar Seres da Floresta", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Guardião": 4}`,
		},
		{
			Name: "Invocar Ínfero", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 6, "Mago": 6}`,
		},
		{
			Name: "Invocação Instantânea de Drawmij", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: true, Concentration: false,
			Classes: `{"Mago": 6}`,
		},
		{
			Name: "JANE KATSUBO Esfera Vitriólica", Edition: "5e", Level: 4, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Labirinto", Edition: "5e", Level: 8, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 8}`,
		},
		{
			Name: "Lendas e Histórias", Edition: "5e", Level: 5, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Clérigo": 5, "Mago": 5}`,
		},
		{
			Name: "Lentidão", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Leque Cromático", Edition: "5e", Level: 1, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Levitação", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Ligação Telepática de Rary", Edition: "5e", Level: 5, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 5, "Mago": 5}`,
		},
		{
			Name: "Limpar a Mente", Edition: "5e", Level: 8, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 8, "Mago": 8}`,
		},
		{
			Name: "Localizar Animais ou Plantas", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 2, "Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Localizar Criatura", Edition: "5e", Level: 4, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4, "Clérigo": 4, "Druida": 4, "Guardião": 4, "Mago": 4, "Paladino": 4}`,
		},
		{
			Name: "Localizar Objeto", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Druida": 2, "Guardião": 2, "Mago": 2, "Paladino": 2}`,
		},
		{
			Name: "Loquacidade", Edition: "5e", Level: 8, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 8, "Bruxo": 8}`,
		},
		{
			Name: "Lufada de Vento", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2, "Feiticeiro": 2, "Guardião": 2, "Mago": 2}`,
		},
		{
			Name: "Luz", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Clérigo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Luz do Dia", Edition: "5e", Level: 3, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 3, "Druida": 3, "Feiticeiro": 3, "Guardião": 3, "Paladino": 3}`,
		},
		{
			Name: "Luzes Dançantes", Edition: "5e", Level: 0, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Lâmina Flamejante", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2, "Feiticeiro": 2}`,
		},
		{
			Name: "Línguas", Edition: "5e", Level: 3, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Bruxo": 3, "Clérigo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Malogro", Edition: "5e", Level: 4, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 4, "Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Mansão Magnífica de Mordenkainen", Edition: "5e", Level: 7, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Mago": 7}`,
		},
		{
			Name: "Manto do Cruzado", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Paladino": 3}`,
		},
		{
			Name: "Marca do Predador", Edition: "5e", Level: 1, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Guardião": 1}`,
		},
		{
			Name: "Mau Olhado", Edition: "5e", Level: 6, School: "Necromancia",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 6, "Bruxo": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Medo", Edition: "5e", Level: 3, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Bruxo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Mensageiro Animal", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 2, "Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Mensagem", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Druida": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Mesclar-se às Rochas", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 3, "Druida": 3, "Guardião": 3}`,
		},
		{
			Name: "Metamorfose", Edition: "5e", Level: 9, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 9, "Mago": 9}`,
		},
		{
			Name: "Miragem Arcana", Edition: "5e", Level: 7, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Druida": 7, "Mago": 7}`,
		},
		{
			Name: "Missão", Edition: "5e", Level: 5, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Clérigo": 5, "Druida": 5, "Mago": 5, "Paladino": 5}`,
		},
		{
			Name: "Modificar Memória", Edition: "5e", Level: 5, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Mago": 5}`,
		},
		{
			Name: "Moldar Rochas", Edition: "5e", Level: 4, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 4, "Druida": 4, "Mago": 4}`,
		},
		{
			Name: "Moléstia", Edition: "5e", Level: 6, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 6}`,
		},
		{
			Name: "Montaria Fantasmagórica", Edition: "5e", Level: 3, School: "Ilusão",
			Ritual: true, Concentration: false,
			Classes: `{"Mago": 3}`,
		},
		{
			Name: "Mover Terra", Edition: "5e", Level: 6, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Movimentação Livre", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 4, "Clérigo": 4, "Druida": 4, "Guardião": 4}`,
		},
		{
			Name: "Muralha Prismática", Edition: "5e", Level: 9, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 9, "Mago": 9}`,
		},
		{
			Name: "Muralha de Energia", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 5}`,
		},
		{
			Name: "Muralha de Espinhos", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 6}`,
		},
		{
			Name: "Muralha de Fogo", Edition: "5e", Level: 4, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Muralha de Gelo", Edition: "5e", Level: 6, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 6}`,
		},
		{
			Name: "Muralha de Pedra", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Muralha de Vento", Edition: "5e", Level: 3, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 3, "Guardião": 3}`,
		},
		{
			Name: "Mão de Bigby", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Mãos Flamejantes", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Mãos Mágicas", Edition: "5e", Level: 0, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Mísseis Mágicos", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Nevasca", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Nuvem Fétida", Edition: "5e", Level: 3, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Nuvem Incendiária", Edition: "5e", Level: 8, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 8, "Feiticeiro": 8, "Mago": 8}`,
		},
		{
			Name: "Nuvem de Adagas", Edition: "5e", Level: 2, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Névoa Mortal", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Névoa Obscurecente", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 1, "Feiticeiro": 1, "Guardião": 1, "Mago": 1}`,
		},
		{
			Name: "Olho Arcano", Edition: "5e", Level: 4, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Onda Destrutiva", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Paladino": 5}`,
		},
		{
			Name: "Onda Trovejante", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Druida": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Oração de Cura", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 2, "Paladino": 2}`,
		},
		{
			Name: "Orbe Cromático", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Orientação", Edition: "5e", Level: 0, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 0, "Druida": 0}`,
		},
		{
			Name: "Padrão Hipnótico", Edition: "5e", Level: 3, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Bruxo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Palavra Curativa", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Clérigo": 1, "Druida": 1}`,
		},
		{
			Name: "Palavra Curativa em Massa", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Clérigo": 3}`,
		},
		{
			Name: "Palavra Sagrada", Edition: "5e", Level: 7, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 7}`,
		},
		{
			Name: "Palavra de Poder: Atordoar", Edition: "5e", Level: 8, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 8, "Bruxo": 8, "Feiticeiro": 8, "Mago": 8}`,
		},
		{
			Name: "Palavra de Poder: Fortificar", Edition: "5e", Level: 7, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Clérigo": 7}`,
		},
		{
			Name: "Palavra de Poder: Matar", Edition: "5e", Level: 9, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 9, "Bruxo": 9, "Feiticeiro": 9, "Mago": 9}`,
		},
		{
			Name: "Palavra de Poder: Salvar", Edition: "5e", Level: 9, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 9, "Clérigo": 9}`,
		},
		{
			Name: "Palavra de Radiância", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 0}`,
		},
		{
			Name: "Palavra de Regresso", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 6}`,
		},
		{
			Name: "Paralisar Monstro", Edition: "5e", Level: 5, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Paralisar Pessoa", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Bruxo": 2, "Clérigo": 2, "Druida": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Parar o Tempo", Edition: "5e", Level: 9, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 9, "Mago": 9}`,
		},
		{
			Name: "Passo Arbóreo", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 5, "Guardião": 5}`,
		},
		{
			Name: "Passo Nebuloso", Edition: "5e", Level: 2, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Passo Sem Rastro", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Passos Largos", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Druida": 1, "Guardião": 1, "Mago": 1}`,
		},
		{
			Name: "Pele-Rocha", Edition: "5e", Level: 4, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Feiticeiro": 4, "Guardião": 4, "Mago": 4}`,
		},
		{
			Name: "Pele-casca", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Pequeno Refúgio de Leomund", Edition: "5e", Level: 3, School: "Evocação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 3, "Mago": 3}`,
		},
		{
			Name: "Perdição", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Clérigo": 1}`,
		},
		{
			Name: "Piscar", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Polimorfia", Edition: "5e", Level: 4, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 4, "Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Polimorfia Total", Edition: "5e", Level: 9, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 9, "Bruxo": 9, "Mago": 9}`,
		},
		{
			Name: "Porta Dimensional", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 4, "Bruxo": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Portais Arcanos", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Portal", Edition: "5e", Level: 9, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 9, "Clérigo": 9, "Feiticeiro": 9, "Mago": 9}`,
		},
		{
			Name: "Praga de Insetos", Edition: "5e", Level: 5, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 5, "Druida": 5, "Feiticeiro": 5}`,
		},
		{
			Name: "Presença Régia de Yolande", Edition: "5e", Level: 5, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Mago": 5}`,
		},
		{
			Name: "Presságio", Edition: "5e", Level: 4, School: "Adivinhação",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 4, "Druida": 4, "Mago": 4}`,
		},
		{
			Name: "Prestidigitação Arcana", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Proibição", Edition: "5e", Level: 6, School: "Abjuração",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 6}`,
		},
		{
			Name: "Projetar Imagem", Edition: "5e", Level: 7, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 7, "Mago": 7}`,
		},
		{
			Name: "Projeção Astral", Edition: "5e", Level: 9, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 9, "Clérigo": 9, "Mago": 9}`,
		},
		{
			Name: "Proteger Fortaleza", Edition: "5e", Level: 6, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 6, "Mago": 6}`,
		},
		{
			Name: "Proteção Contra Lâminas", Edition: "5e", Level: 0, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Proteção Contra Veneno", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 2, "Druida": 2, "Guardião": 2, "Paladino": 2}`,
		},
		{
			Name: "Proteção Contra a Morte", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 4, "Paladino": 4}`,
		},
		{
			Name: "Proteção Contra o Bem e o Mal", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 1, "Clérigo": 1, "Paladino": 1}`,
		},
		{
			Name: "Proteção Contra o Bem e o Mal Abjuração C, M Purificar Alimentos e Bebidas", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: true, Concentration: false,
			Classes: `{"Druida": 1}`,
		},
		{
			Name: "Proteção Contra o Bem e o Mal Abjuração C, M Queda Suave", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 1}`,
		},
		{
			Name: "Proteção contra Energia", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 3, "Druida": 3, "Feiticeiro": 3, "Guardião": 3, "Mago": 3}`,
		},
		{
			Name: "Purificar Alimentos e Bebidas", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 1, "Paladino": 1}`,
		},
		{
			Name: "Queda Suave", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1, "Feiticeiro": 1}`,
		},
		{
			Name: "Raio Ardente", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Raio Guia", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 1}`,
		},
		{
			Name: "Raio Lunar", Edition: "5e", Level: 2, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 2}`,
		},
		{
			Name: "Raio Místico", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 0}`,
		},
		{
			Name: "Raio Nauseante", Edition: "5e", Level: 1, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Raio Solar", Edition: "5e", Level: 6, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 6, "Druida": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Raio de Bruxa", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Raio de Fogo", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Raio de Gelo", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Raio do Enfraquecimento", Edition: "5e", Level: 2, School: "Necromancia",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 2, "Mago": 2}`,
		},
		{
			Name: "Rajada Prismática", Edition: "5e", Level: 7, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Rajada de Veneno", Edition: "5e", Level: 0, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Receptáculo Arcano", Edition: "5e", Level: 6, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 6}`,
		},
		{
			Name: "Reencarnar", Edition: "5e", Level: 5, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 5}`,
		},
		{
			Name: "Reflexos", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Refugiar", Edition: "5e", Level: 7, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 7}`,
		},
		{
			Name: "Regeneração", Edition: "5e", Level: 7, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Clérigo": 7, "Druida": 7}`,
		},
		{
			Name: "Relâmpago", Edition: "5e", Level: 3, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Remeter", Edition: "5e", Level: 3, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 3, "Clérigo": 3, "Mago": 3}`,
		},
		{
			Name: "Remover Maldição", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 3, "Clérigo": 3, "Mago": 3, "Paladino": 3}`,
		},
		{
			Name: "Reparar", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Clérigo": 0, "Druida": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Repouso Tranquilo", Edition: "5e", Level: 2, School: "Necromancia",
			Ritual: true, Concentration: false,
			Classes: `{"Clérigo": 2, "Mago": 2, "Paladino": 2}`,
		},
		{
			Name: "Repreensão Diabólica", Edition: "5e", Level: 1, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 1}`,
		},
		{
			Name: "Resistência", Edition: "5e", Level: 0, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 0, "Druida": 0}`,
		},
		{
			Name: "Respirar na Água", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: true, Concentration: false,
			Classes: `{"Druida": 3, "Feiticeiro": 3, "Guardião": 3, "Mago": 3}`,
		},
		{
			Name: "Ressurreição", Edition: "5e", Level: 7, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Clérigo": 7}`,
		},
		{
			Name: "Ressurreição Verdadeira", Edition: "5e", Level: 9, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 9, "Druida": 9}`,
		},
		{
			Name: "Restauração Maior", Edition: "5e", Level: 5, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Clérigo": 5, "Druida": 5, "Guardião": 5, "Paladino": 5}`,
		},
		{
			Name: "Restauração Menor", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Druida": 2, "Guardião": 2, "Paladino": 2}`,
		},
		{
			Name: "Retirada Acelerada", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Reviver os Mortos", Edition: "5e", Level: 5, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Clérigo": 5, "Paladino": 5}`,
		},
		{
			Name: "Revivificar", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 3, "Druida": 3, "Guardião": 3, "Paladino": 3}`,
		},
		{
			Name: "Rogar Maldição", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 3, "Clérigo": 3, "Mago": 3}`,
		},
		{
			Name: "Salto", Edition: "5e", Level: 1, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 1, "Feiticeiro": 1, "Guardião": 1, "Mago": 1}`,
		},
		{
			Name: "Santuário", Edition: "5e", Level: 1, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 1}`,
		},
		{
			Name: "Santuário Particular de Mordenkainen", Edition: "5e", Level: 4, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Saraivada de Espinhos", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Guardião": 1}`,
		},
		{
			Name: "Semiplano", Edition: "5e", Level: 8, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 8, "Feiticeiro": 8, "Mago": 8}`,
		},
		{
			Name: "Sentido Feral", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: true, Concentration: true,
			Classes: `{"Druida": 2, "Guardião": 2}`,
		},
		{
			Name: "Servo Invisível", Edition: "5e", Level: 1, School: "Invocação",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 1, "Bruxo": 1, "Mago": 1}`,
		},
		{
			Name: "Sexto Sentido", Edition: "5e", Level: 9, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 9, "Bruxo": 9, "Druida": 9, "Mago": 9}`,
		},
		{
			Name: "Silêncio", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: true, Concentration: true,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Guardião": 2}`,
		},
		{
			Name: "Similaridade", Edition: "5e", Level: 5, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Simulacro", Edition: "5e", Level: 7, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 7}`,
		},
		{
			Name: "Simular Morte", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: true, Concentration: false,
			Classes: `{"Bardo": 3, "Clérigo": 3, "Druida": 3, "Mago": 3}`,
		},
		{
			Name: "Sinal de Esperança", Edition: "5e", Level: 3, School: "Abjuração",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 3}`,
		},
		{
			Name: "Sonho", Edition: "5e", Level: 5, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Mago": 5}`,
		},
		{
			Name: "Sono", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 1, "Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Sopro de Dragão", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Sugestão", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 2, "Bruxo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Sugestão em Massa", Edition: "5e", Level: 6, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Suplício", Edition: "5e", Level: 8, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 8, "Bruxo": 8, "Druida": 8, "Mago": 8}`,
		},
		{
			Name: "Sussurros Dissonantes", Edition: "5e", Level: 1, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 1}`,
		},
		{
			Name: "Símbolo", Edition: "5e", Level: 7, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Clérigo": 7, "Druida": 7, "Mago": 7}`,
		},
		{
			Name: "Talho Mental", Edition: "5e", Level: 0, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Taumaturgia", Edition: "5e", Level: 0, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 0}`,
		},
		{
			Name: "Teia", Edition: "5e", Level: 2, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Telecinese", Edition: "5e", Level: 5, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 5, "Mago": 5}`,
		},
		{
			Name: "Telepatia", Edition: "5e", Level: 8, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 8}`,
		},
		{
			Name: "Teleporte", Edition: "5e", Level: 7, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 7, "Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Tempestade Glacial", Edition: "5e", Level: 4, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 4, "Feiticeiro": 4, "Mago": 4}`,
		},
		{
			Name: "Tempestade Radiante de Jallarzi", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 5}`,
		},
		{
			Name: "Tempestade Radiante de Jallazar", Edition: "5e", Level: 5, School: "Evocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 5}`,
		},
		{
			Name: "Tempestade da Vingança", Edition: "5e", Level: 9, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 9}`,
		},
		{
			Name: "Tempestade de Fogo", Edition: "5e", Level: 7, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 7, "Druida": 7, "Feiticeiro": 7}`,
		},
		{
			Name: "Tentáculos Negros de Evard", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Mago": 4}`,
		},
		{
			Name: "Terremoto", Edition: "5e", Level: 8, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Clérigo": 8, "Druida": 8, "Feiticeiro": 8}`,
		},
		{
			Name: "Terreno Alucinatório", Edition: "5e", Level: 4, School: "Ilusão",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 4, "Bruxo": 4, "Druida": 4, "Mago": 4}`,
		},
		{
			Name: "Toque Chocante", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Toque Necrótico", Edition: "5e", Level: 0, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Toque Vampírico", Edition: "5e", Level: 3, School: "Necromancia",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Tranca Arcana", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Mago": 2}`,
		},
		{
			Name: "Transição Planar", Edition: "5e", Level: 7, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 7, "Clérigo": 7, "Druida": 7, "Feiticeiro": 7, "Mago": 7}`,
		},
		{
			Name: "Transporte via Plantas", Edition: "5e", Level: 6, School: "Invocação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 6}`,
		},
		{
			Name: "Trovão", Edition: "5e", Level: 0, School: "Evocação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0, "Bruxo": 0, "Druida": 0, "Feiticeiro": 0, "Mago": 0}`,
		},
		{
			Name: "Tsunami", Edition: "5e", Level: 8, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 8}`,
		},
		{
			Name: "Turvar", Edition: "5e", Level: 2, School: "Ilusão",
			Ritual: false, Concentration: true,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Ver o Invisível", Edition: "5e", Level: 2, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Vidência", Edition: "5e", Level: 5, School: "Adivinhação",
			Ritual: false, Concentration: true,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Clérigo": 5, "Druida": 5, "Mago": 5}`,
		},
		{
			Name: "Vigor Arcano", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 2, "Mago": 2}`,
		},
		{
			Name: "Vinha Agarradora", Edition: "5e", Level: 4, School: "Invocação",
			Ritual: false, Concentration: true,
			Classes: `{"Druida": 4, "Guardião": 4}`,
		},
		{
			Name: "Visão da Verdade", Edition: "5e", Level: 6, School: "Adivinhação",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 6, "Bruxo": 6, "Clérigo": 6, "Feiticeiro": 6, "Mago": 6}`,
		},
		{
			Name: "Visão no Escuro", Edition: "5e", Level: 2, School: "Transmutação",
			Ritual: false, Concentration: false,
			Classes: `{"Druida": 2, "Feiticeiro": 2, "Guardião": 2, "Mago": 2}`,
		},
		{
			Name: "Vitalidade Vazia", Edition: "5e", Level: 1, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Feiticeiro": 1, "Mago": 1}`,
		},
		{
			Name: "Voo", Edition: "5e", Level: 3, School: "Transmutação",
			Ritual: false, Concentration: true,
			Classes: `{"Bruxo": 3, "Feiticeiro": 3, "Mago": 3}`,
		},
		{
			Name: "Vínculo de Proteção", Edition: "5e", Level: 2, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Clérigo": 2, "Paladino": 2}`,
		},
		{
			Name: "WAYNE ENGLAND Rajada de Veneno", Edition: "5e", Level: 0, School: "Necromancia",
			Ritual: false, Concentration: false,
			Classes: `{"Bruxo": 0}`,
		},
		{
			Name: "Zombaria Perversa", Edition: "5e", Level: 0, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 0}`,
		},
		{
			Name: "Zona da Verdade", Edition: "5e", Level: 2, School: "Encantamento",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 2, "Clérigo": 2, "Paladino": 2}`,
		},
		{
			Name: "Âncora Planar", Edition: "5e", Level: 5, School: "Abjuração",
			Ritual: false, Concentration: false,
			Classes: `{"Bardo": 5, "Bruxo": 5, "Clérigo": 5, "Druida": 5, "Mago": 5}`,
		},
	}
	for _, s := range spells {
		var existing domain.Spell
		if db.Where("name = ? AND edition = ?", s.Name, "5e").First(&existing).Error != nil {
			db.Create(&s)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"level": s.Level, "school": s.School, "ritual": s.Ritual,
				"concentration": s.Concentration, "classes": s.Classes,
			})
		}
	}
	log.Println("  ✓ Magias 5e seedadas:", len(spells))
}
