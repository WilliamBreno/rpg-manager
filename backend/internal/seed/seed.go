package seed

import (
	"log"
	"rpg-manager/internal/domain"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	log.Println("🌱 Rodando seed...")
	seedClasses(db)
	seedSkills(db)
	log.Println("✅ Seed concluído!")
}

func seedClasses(db *gorm.DB) {
	classes := []domain.Class{
		// ── D&D 4e ──────────────────────────────────────────────────
		{Name: "Bardo", Edition: "4e", Description: "Líder arcano que usa música e conhecimento como armas.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Bárbaro", Edition: "4e", Description: "Combatente primitivo de força e fúria selvagem.", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 8, FortBonus: 2, IsDefault: true},
		{Name: "Clérigo", Edition: "4e", Description: "Líder divino que canaliza o poder de sua divindade.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 1, IsDefault: true},
		{Name: "Druida", Edition: "4e", Description: "Controlador primitivo com maestria da natureza.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Feiticeiro", Edition: "4e", Description: "Arcano que canaliza poder dracônico ou do caos.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Guardião", Edition: "4e", Description: "Defensor primitivo e protetor do mundo natural.", BaseHP: 17, HPPerLevel: 7, SurgesPerDay: 9, FortBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Guerreiro", Edition: "4e", Description: "Defensor marcial especialista em combate corpo a corpo.", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 9, FortBonus: 2, IsDefault: true},
		{Name: "Invocador", Edition: "4e", Description: "Controlador divino que invoca o poder dos deuses.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Ladino", Edition: "4e", Description: "Agressor furtivo especializado em ataques precisos.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, ReflBonus: 2, IsDefault: true},
		{Name: "Mago", Edition: "4e", Description: "Controlador arcano de grande poder mágico.", BaseHP: 10, HPPerLevel: 4, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Monge", Edition: "4e", Description: "Agressor psiônico de disciplina e artes marciais.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Paladino", Edition: "4e", Description: "Defensor divino e campeão sagrado.", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 10, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Patrulheiro", Edition: "4e", Description: "Agressor marcial e batedor das regiões ermas.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6, FortBonus: 1, ReflBonus: 1, IsDefault: true},
		{Name: "Psionista", Edition: "4e", Description: "Controlador psiônico de grande poder mental.", BaseHP: 12, HPPerLevel: 4, SurgesPerDay: 6, WillBonus: 2, IsDefault: true},
		{Name: "Rastreador", Edition: "4e", Description: "Controlador primitivo especialista em emboscadas.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Vingador", Edition: "4e", Description: "Agressor divino executor da vontade dos deuses.", BaseHP: 14, HPPerLevel: 6, SurgesPerDay: 7, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
		{Name: "Xamã", Edition: "4e", Description: "Líder primitivo que canaliza espíritos do mundo natural.", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7, FortBonus: 1, WillBonus: 1, IsDefault: true},
		// ── D&D 5e ──────────────────────────────────────────────────
		{Name: "Bárbaro", Edition: "5e", Description: "Guerreiro furioso de força bruta.", HitDie: 12, IsDefault: true},
		{Name: "Bardo", Edition: "5e", Description: "Conjurador versátil com magia e música.", HitDie: 8, IsDefault: true},
		{Name: "Bruxo", Edition: "5e", Description: "Conjurador com poder concedido por um patrono.", HitDie: 8, IsDefault: true},
		{Name: "Clérigo", Edition: "5e", Description: "Servidor divino com poderes de cura e combate.", HitDie: 8, IsDefault: true},
		{Name: "Druida", Edition: "5e", Description: "Guardião da natureza com magia elemental.", HitDie: 8, IsDefault: true},
		{Name: "Feiticeiro", Edition: "5e", Description: "Conjurador de magia inata e poderosa.", HitDie: 6, IsDefault: true},
		{Name: "Guerreiro", Edition: "5e", Description: "Mestre do combate e das armas.", HitDie: 10, IsDefault: true},
		{Name: "Ladino", Edition: "5e", Description: "Especialista furtivo e hábil.", HitDie: 8, IsDefault: true},
		{Name: "Mago", Edition: "5e", Description: "Estudioso da magia arcana.", HitDie: 6, IsDefault: true},
		{Name: "Monge", Edition: "5e", Description: "Combatente de artes marciais e energia mística.", HitDie: 8, IsDefault: true},
		{Name: "Paladino", Edition: "5e", Description: "Guerreiro sagrado jurado a um ideal.", HitDie: 10, IsDefault: true},
		{Name: "Patrulheiro", Edition: "5e", Description: "Caçador e rastreador das regiões selvagens.", HitDie: 10, IsDefault: true},
		{Name: "Xamã", Edition: "5e", Description: "Líder espiritual conectado ao mundo natural.", HitDie: 8, IsDefault: true},
	}

	for _, c := range classes {
		var existing domain.Class
		result := db.Where("name = ? AND edition = ?", c.Name, c.Edition).First(&existing)
		if result.Error != nil {
			// Não existe — cria
			if err := db.Create(&c).Error; err != nil {
				log.Printf("Erro ao criar classe %s (%s): %v", c.Name, c.Edition, err)
			}
		} else {
			// Já existe — atualiza os campos numéricos
			db.Model(&existing).Updates(map[string]interface{}{
				"base_hp":       c.BaseHP,
				"hp_per_level":  c.HPPerLevel,
				"surges_per_day": c.SurgesPerDay,
				"fort_bonus":    c.FortBonus,
				"refl_bonus":    c.ReflBonus,
				"will_bonus":    c.WillBonus,
				"hit_die":       c.HitDie,
				"description":   c.Description,
			})
		}
	}
	log.Println("  ✓ Classes seedadas")
}

func seedSkills(db *gorm.DB) {
	// Busca o ID do Bardo 4e
	var bardoClass domain.Class
	if err := db.Where("name = ? AND edition = ?", "Bardo", "4e").First(&bardoClass).Error; err != nil {
		log.Println("  ✗ Bardo 4e não encontrado, pulando skills")
		return
	}
	bardoID := bardoClass.ID

	skills := []domain.Skill{
		// ── Nível 1 — Sem Limite ────────────────────────────────────
		{Name: "Golpe Condutor", Description: "O golpe de sua arma guia seus aliados, mostrando a eles onde devem focalizar seus ataques.", PowerType: domain.PowerUnlimited, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Golpe da Canção de Guerra", Description: "Com sua canção de vitória e de guerra, seus aliados se sentem revigorados a cada ataque.", PowerType: domain.PowerUnlimited, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Marca Indireta", Description: "Ao ocultar seu ataque arcano, você engana o adversário, fazendo-o crer que o ataque veio de um aliado.", PowerType: domain.PowerUnlimited, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Zombaria Malévola", Description: "Emitindo impropérios contra seu adversário, você o envolve na magia dos bardos e o coloca num estado de fúria cega.", PowerType: domain.PowerUnlimited, Level: 1, Edition: "4e", ClassID: &bardoID},
		// ── Nível 1 — Por Encontro ──────────────────────────────────
		{Name: "Grito do Triunfo", Description: "Você emite um poderoso chamado ao combate, espalhando os inimigos e instigando seus aliados avante.", PowerType: domain.PowerEncounter, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Refrão Inspirador", Description: "Sua arma entoa uma canção arcana que guia seus aliados para a vitória.", PowerType: domain.PowerEncounter, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Amigos Rápidos", Description: "Entoando uma melodia de falsa amizade, você leva o alvo a um estado de puro devaneio.", PowerType: domain.PowerEncounter, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Gafe", Description: "Suas palavras certeiras confundem e perturbam o inimigo, deixando-o vulnerável.", PowerType: domain.PowerEncounter, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Ecos do Guardião", Description: "Sons arcanos envolvem o alvo, dificultando seus movimentos e ataques.", PowerType: domain.PowerEncounter, Level: 1, Edition: "4e", ClassID: &bardoID},
		// ── Nível 1 — Diário ────────────────────────────────────────
		{Name: "Grito Inspirador", Description: "Seu brado furioso apunhala a mente do inimigo. Sempre que ele for atingido por um aliado, o vigor lhe será negado.", PowerType: domain.PowerDaily, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Canção do Matador", Description: "Sua canção de batalha enfraquece as defesas do inimigo a cada golpe desferido.", PowerType: domain.PowerDaily, Level: 1, Edition: "4e", ClassID: &bardoID},
		{Name: "Verso do Triunfo", Description: "Com palavras inspiradoras, você incentiva seus aliados ao ataque com vigor renovado.", PowerType: domain.PowerDaily, Level: 1, Edition: "4e", ClassID: &bardoID},
		// ── Nível 2 — Utilitário ────────────────────────────────────
		{Name: "Melodia do Caçador", Description: "Ao moldar uma corrente sonora, você elimina todos os ruídos, criando uma área de silêncio absoluto.", PowerType: domain.PowerUtility, Level: 2, Edition: "4e", ClassID: &bardoID},
		{Name: "Canção da Coragem", Description: "Uma canção de encorajamento fortalece seus aliados para enfrentar os desafios à frente.", PowerType: domain.PowerUtility, Level: 2, Edition: "4e", ClassID: &bardoID},
		{Name: "Canção da Defesa", Description: "Uma melodia defensiva protege você e seus aliados de ataques iminentes.", PowerType: domain.PowerUtility, Level: 2, Edition: "4e", ClassID: &bardoID},
		{Name: "Inspirar Competência", Description: "Sua música inspira um aliado a realizar tarefas além de suas capacidades normais.", PowerType: domain.PowerUtility, Level: 2, Edition: "4e", ClassID: &bardoID},
		// ── Nível 3 — Sem Limite ────────────────────────────────────
		{Name: "Energia Impulsora", Description: "A magia cria um ataque de energia que afasta um inimigo próximo de um aliado.", PowerType: domain.PowerUnlimited, Level: 3, Edition: "4e", ClassID: &bardoID},
		{Name: "Estrofe Dissonante", Description: "Sons discordantes perturbam o inimigo, reduzindo sua eficácia em combate.", PowerType: domain.PowerUnlimited, Level: 3, Edition: "4e", ClassID: &bardoID},
		{Name: "Ferocidade Astuta", Description: "Combinando astúcia e ferocidade, você lança um ataque desconcertante contra o inimigo.", PowerType: domain.PowerUnlimited, Level: 3, Edition: "4e", ClassID: &bardoID},
		// ── Nível 3 — Diário ────────────────────────────────────────
		{Name: "Chamado aos Cavalos de Batalha", Description: "Com palavras inspiradoras, você incentiva seus aliados ao ataque com renovado vigor.", PowerType: domain.PowerDaily, Level: 3, Edition: "4e", ClassID: &bardoID},
		// ── Nível 5 — Sem Limite ────────────────────────────────────
		{Name: "Palavra de Proteção Mística", Description: "Uma palavra de poder cria uma barreira mística protetora ao redor de um aliado.", PowerType: domain.PowerUnlimited, Level: 5, Edition: "4e", ClassID: &bardoID},
		{Name: "Sátira da Bravura", Description: "Sua sátira desmoraliza os inimigos enquanto inspira bravura em seus aliados.", PowerType: domain.PowerUnlimited, Level: 5, Edition: "4e", ClassID: &bardoID},
		// ── Nível 5 — Por Encontro ──────────────────────────────────
		{Name: "Melodia do Gelo e do Vento", Description: "Uma melodia gélida lança ventos cortantes contra seus inimigos, congelando-os no lugar.", PowerType: domain.PowerEncounter, Level: 5, Edition: "4e", ClassID: &bardoID},
		// ── Nível 5 — Diário ────────────────────────────────────────
		{Name: "Canção da Discórdia", Description: "Enchendo um inimigo de desconfiança, você o obriga a atacar um aliado.", PowerType: domain.PowerDaily, Level: 5, Edition: "4e", ClassID: &bardoID},
		// ── Nível 6 — Utilitário ────────────────────────────────────
		{Name: "Alegro", Description: "Você cria um ritmo apressado que concede velocidade a você e seus aliados.", PowerType: domain.PowerUtility, Level: 6, Edition: "4e", ClassID: &bardoID},
		{Name: "Canção da Conquista", Description: "Uma canção triunfante motiva seus aliados a avançar com determinação renovada.", PowerType: domain.PowerUtility, Level: 6, Edition: "4e", ClassID: &bardoID},
		{Name: "Cura do Trapaceiro", Description: "Uma melodia enganosa restaura a vitalidade de um aliado ferido.", PowerType: domain.PowerUtility, Level: 6, Edition: "4e", ClassID: &bardoID},
		{Name: "Ode ao Sacrifício", Description: "Uma ode ao sacrifício heroico inspira um aliado a agir com abnegação.", PowerType: domain.PowerUtility, Level: 6, Edition: "4e", ClassID: &bardoID},
		// ── Nível 7 — Sem Limite ────────────────────────────────────
		{Name: "Desviar a Atenção", Description: "Suas palavras desviam a atenção do inimigo, permitindo que aliados se movam livremente.", PowerType: domain.PowerUnlimited, Level: 7, Edition: "4e", ClassID: &bardoID},
		{Name: "Golpe da Garra do Escorpião", Description: "Sua distração permite que um aliado se desloque ao redor do inimigo para uma posição vantajosa.", PowerType: domain.PowerUnlimited, Level: 7, Edition: "4e", ClassID: &bardoID},
		{Name: "Lâmina do Trovão", Description: "Sua arma retumba como um trovão, atingindo um inimigo e permitindo movimentação aliada.", PowerType: domain.PowerUnlimited, Level: 7, Edition: "4e", ClassID: &bardoID},
		// ── Nível 7 — Diário ────────────────────────────────────────
		{Name: "Grito de Distração", Description: "Um grito ensurdecedor distrai os inimigos, deixando-os vulneráveis aos ataques aliados.", PowerType: domain.PowerDaily, Level: 7, Edition: "4e", ClassID: &bardoID},
		{Name: "Má-Sorte", Description: "O que parecia ser uma ode ao destino se torna uma maldição que persegue o inimigo.", PowerType: domain.PowerDaily, Level: 7, Edition: "4e", ClassID: &bardoID},
		// ── Nível 9 — Por Encontro ──────────────────────────────────
		{Name: "Hino do Resgate Audacioso", Description: "Seu ataque ressoa uma canção arcana que permite a um aliado se deslocar para posição segura.", PowerType: domain.PowerEncounter, Level: 9, Edition: "4e", ClassID: &bardoID},
		{Name: "Riso Repugnante", Description: "Um terrível ataque de riso convulsivo assola o alvo, incapacitando-o temporariamente.", PowerType: domain.PowerEncounter, Level: 9, Edition: "4e", ClassID: &bardoID},
		{Name: "Condutor Poderoso", Description: "Como um maestro poderoso, você dirige seus aliados em um ataque coordenado devastador.", PowerType: domain.PowerEncounter, Level: 9, Edition: "4e", ClassID: &bardoID},
		// ── Nível 10 — Utilitário ───────────────────────────────────
		{Name: "Canção da Recuperação", Description: "Com uma canção inspiradora, seus aliados ficam repletos de perseverança e recuperam vigor.", PowerType: domain.PowerUtility, Level: 10, Edition: "4e", ClassID: &bardoID},
		{Name: "Palavra da Vida", Description: "Uma palavra de poder restaura a vitalidade de um aliado gravemente ferido.", PowerType: domain.PowerUtility, Level: 10, Edition: "4e", ClassID: &bardoID},
		{Name: "Rasura Ilusória", Description: "Uma ilusão sonora confunde os inimigos sobre a posição real de seus aliados.", PowerType: domain.PowerUtility, Level: 10, Edition: "4e", ClassID: &bardoID},
		// ── Nível 10 — Diário ───────────────────────────────────────
		{Name: "Véu", Description: "Um véu de ilusão sonora oculta você e seus aliados dos sentidos inimigos.", PowerType: domain.PowerDaily, Level: 10, Edition: "4e", ClassID: &bardoID},
	}

	for _, skill := range skills {
		var existing domain.Skill
		result := db.Where("name = ? AND edition = ? AND class_id = ?", skill.Name, skill.Edition, bardoID).First(&existing)
		if result.Error != nil {
			if err := db.Create(&skill).Error; err != nil {
				log.Printf("  Erro ao criar skill %s: %v", skill.Name, err)
			}
		}
	}
	log.Printf("  ✓ %d habilidades do Bardo 4e seedadas", len(skills))
}