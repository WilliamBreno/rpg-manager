package config

import (
    "log"
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
    seedClasses(db)
    seedRaces(db)
    seedArmors(db)
    log.Println("Seed concluído com sucesso!")
}

func seedClasses(db *gorm.DB) {
    classes := []domain.Class{
        // D&D 1e (AD&D)
        {Name: "Guerreiro", Edition: "1e", Description: "Especialista em combate corpo a corpo.", HitDie: 10, IsDefault: true},
        {Name: "Mago", Edition: "1e", Description: "Usuário de magia arcana poderosa.", HitDie: 4, IsDefault: true},
        {Name: "Clérigo", Edition: "1e", Description: "Servo divino com poderes de cura.", HitDie: 8, IsDefault: true},
        {Name: "Ladrão", Edition: "1e", Description: "Especialista em furtividade e trapaças.", HitDie: 6, IsDefault: true},
        {Name: "Paladino", Edition: "1e", Description: "Guerreiro sagrado a serviço do bem.", HitDie: 10, IsDefault: true},
        {Name: "Ranger", Edition: "1e", Description: "Guerreiro das terras selvagens.", HitDie: 8, IsDefault: true},
        {Name: "Druida", Edition: "1e", Description: "Guardião da natureza com magias naturais.", HitDie: 8, IsDefault: true},
        {Name: "Ilusionista", Edition: "1e", Description: "Mago especializado em ilusões.", HitDie: 4, IsDefault: true},
        {Name: "Assassino", Edition: "1e", Description: "Especialista em eliminar alvos silenciosamente.", HitDie: 6, IsDefault: true},
        {Name: "Monge", Edition: "1e", Description: "Combatente de artes marciais com disciplina espiritual.", HitDie: 6, IsDefault: true},
        {Name: "Bardo", Edition: "1e", Description: "Artista versátil com habilidades mágicas.", HitDie: 6, IsDefault: true},

        // D&D 2e (AD&D 2nd Edition)
        {Name: "Guerreiro", Edition: "2e", Description: "Especialista em combate corpo a corpo.", HitDie: 10, IsDefault: true},
        {Name: "Paladino", Edition: "2e", Description: "Guerreiro sagrado a serviço do bem.", HitDie: 10, IsDefault: true},
        {Name: "Ranger", Edition: "2e", Description: "Guerreiro das terras selvagens.", HitDie: 8, IsDefault: true},
        {Name: "Mago", Edition: "2e", Description: "Usuário de magia arcana poderosa.", HitDie: 4, IsDefault: true},
        {Name: "Ilusionista", Edition: "2e", Description: "Mago especializado em ilusões.", HitDie: 4, IsDefault: true},
        {Name: "Clérigo", Edition: "2e", Description: "Servo divino com poderes de cura.", HitDie: 8, IsDefault: true},
        {Name: "Druida", Edition: "2e", Description: "Guardião da natureza com magias naturais.", HitDie: 8, IsDefault: true},
        {Name: "Xamã", Edition: "2e", Description: "Líder espiritual com conexão com os espíritos.", HitDie: 6, IsDefault: true},
        {Name: "Ladrão", Edition: "2e", Description: "Especialista em furtividade e trapaças.", HitDie: 6, IsDefault: true},
        {Name: "Bardo", Edition: "2e", Description: "Artista versátil com habilidades mágicas.", HitDie: 6, IsDefault: true},
        {Name: "Monge", Edition: "2e", Description: "Combatente de artes marciais com disciplina espiritual.", HitDie: 6, IsDefault: true},

        // D&D 3e
        {Name: "Bárbaro", Edition: "3e", Description: "Guerreiro selvagem movido pela fúria.", HitDie: 12, IsDefault: true},
        {Name: "Bardo", Edition: "3e", Description: "Artista versátil com habilidades mágicas.", HitDie: 6, IsDefault: true},
        {Name: "Clérigo", Edition: "3e", Description: "Servo divino com poderes de cura.", HitDie: 8, IsDefault: true},
        {Name: "Druida", Edition: "3e", Description: "Guardião da natureza com magias naturais.", HitDie: 8, IsDefault: true},
        {Name: "Guerreiro", Edition: "3e", Description: "Especialista em combate corpo a corpo.", HitDie: 10, IsDefault: true},
        {Name: "Monge", Edition: "3e", Description: "Combatente de artes marciais com disciplina espiritual.", HitDie: 8, IsDefault: true},
        {Name: "Paladino", Edition: "3e", Description: "Guerreiro sagrado a serviço do bem.", HitDie: 10, IsDefault: true},
        {Name: "Ranger", Edition: "3e", Description: "Guerreiro das terras selvagens.", HitDie: 8, IsDefault: true},
        {Name: "Ladino", Edition: "3e", Description: "Especialista em furtividade e trapaças.", HitDie: 6, IsDefault: true},
        {Name: "Feiticeiro", Edition: "3e", Description: "Lançador de magias com poder inato.", HitDie: 4, IsDefault: true},
        {Name: "Mago", Edition: "3e", Description: "Usuário de magia arcana poderosa.", HitDie: 4, IsDefault: true},

        // D&D 3.5e
        {Name: "Bárbaro", Edition: "3.5e", Description: "Guerreiro selvagem movido pela fúria.", HitDie: 12, IsDefault: true},
        {Name: "Bardo", Edition: "3.5e", Description: "Artista versátil com habilidades mágicas.", HitDie: 6, IsDefault: true},
        {Name: "Clérigo", Edition: "3.5e", Description: "Servo divino com poderes de cura.", HitDie: 8, IsDefault: true},
        {Name: "Druida", Edition: "3.5e", Description: "Guardião da natureza com magias naturais.", HitDie: 8, IsDefault: true},
        {Name: "Guerreiro", Edition: "3.5e", Description: "Especialista em combate corpo a corpo.", HitDie: 10, IsDefault: true},
        {Name: "Monge", Edition: "3.5e", Description: "Combatente de artes marciais com disciplina espiritual.", HitDie: 8, IsDefault: true},
        {Name: "Paladino", Edition: "3.5e", Description: "Guerreiro sagrado a serviço do bem.", HitDie: 10, IsDefault: true},
        {Name: "Ranger", Edition: "3.5e", Description: "Guerreiro das terras selvagens.", HitDie: 8, IsDefault: true},
        {Name: "Ladino", Edition: "3.5e", Description: "Especialista em furtividade e trapaças.", HitDie: 6, IsDefault: true},
        {Name: "Feiticeiro", Edition: "3.5e", Description: "Lançador de magias com poder inato.", HitDie: 4, IsDefault: true},
        {Name: "Mago", Edition: "3.5e", Description: "Usuário de magia arcana poderosa.", HitDie: 4, IsDefault: true},

        // D&D 4e
        {Name: "Clérigo", Edition: "4e", Description: "Líder divino com poderes de cura e suporte.", HitDie: 8, IsDefault: true},
        {Name: "Guerreiro", Edition: "4e", Description: "Tanque resistente especialista em combate.", HitDie: 10, IsDefault: true},
        {Name: "Paladino", Edition: "4e", Description: "Guerreiro sagrado defensor dos aliados.", HitDie: 10, IsDefault: true},
        {Name: "Ranger", Edition: "4e", Description: "Atacante ágil especialista em dois alvos.", HitDie: 8, IsDefault: true},
        {Name: "Ladino", Edition: "4e", Description: "Atacante furtivo com vantagem tática.", HitDie: 6, IsDefault: true},
        {Name: "Mago", Edition: "4e", Description: "Controlador arcano com magias devastadoras.", HitDie: 4, IsDefault: true},
        {Name: "Bárbaro", Edition: "4e", Description: "Atacante selvagem com fúria incontrolável.", HitDie: 12, IsDefault: true},
        {Name: "Bardo", Edition: "4e", Description: "Líder versátil com habilidades arcanas.", HitDie: 6, IsDefault: true},
        {Name: "Druida", Edition: "4e", Description: "Controlador da natureza com transformações.", HitDie: 8, IsDefault: true},
        {Name: "Invoker", Edition: "4e", Description: "Canal divino de poder destrutivo.", HitDie: 6, IsDefault: true},
        {Name: "Shaman", Edition: "4e", Description: "Líder espiritual com companheiro espírito.", HitDie: 6, IsDefault: true},
        {Name: "Warlock", Edition: "4e", Description: "Usuário de pactos com poderes sombrios.", HitDie: 6, IsDefault: true},
        {Name: "Warlord", Edition: "4e", Description: "Líder tático que potencializa aliados.", HitDie: 8, IsDefault: true},

        // D&D 5e
        {Name: "Bárbaro", Edition: "5e", Description: "Guerreiro selvagem movido pela fúria.", HitDie: 12, IsDefault: true},
        {Name: "Bardo", Edition: "5e", Description: "Artista versátil com habilidades mágicas.", HitDie: 8, IsDefault: true},
        {Name: "Clérigo", Edition: "5e", Description: "Servo divino com poderes de cura.", HitDie: 8, IsDefault: true},
        {Name: "Druida", Edition: "5e", Description: "Guardião da natureza com magias naturais.", HitDie: 8, IsDefault: true},
        {Name: "Guerreiro", Edition: "5e", Description: "Especialista em combate corpo a corpo.", HitDie: 10, IsDefault: true},
        {Name: "Monge", Edition: "5e", Description: "Combatente de artes marciais com disciplina espiritual.", HitDie: 8, IsDefault: true},
        {Name: "Paladino", Edition: "5e", Description: "Guerreiro sagrado a serviço do bem.", HitDie: 10, IsDefault: true},
        {Name: "Ranger", Edition: "5e", Description: "Guerreiro das terras selvagens.", HitDie: 8, IsDefault: true},
        {Name: "Ladino", Edition: "5e", Description: "Especialista em furtividade e trapaças.", HitDie: 8, IsDefault: true},
        {Name: "Feiticeiro", Edition: "5e", Description: "Lançador de magias com poder inato.", HitDie: 6, IsDefault: true},
        {Name: "Bruxo", Edition: "5e", Description: "Usuário de pactos com entidades poderosas.", HitDie: 8, IsDefault: true},
        {Name: "Mago", Edition: "5e", Description: "Usuário de magia arcana poderosa.", HitDie: 6, IsDefault: true},
    }

    for _, class := range classes {
        var existing domain.Class
        result := db.Where("name = ? AND edition = ? AND is_default = ?", class.Name, class.Edition, true).First(&existing)
        if result.Error != nil {
            db.Create(&class)
        }
    }

    log.Println("Classes populadas com sucesso!")
}

func seedRaces(db *gorm.DB) {
    races := []domain.Race{
        // D&D 1e
        {Name: "Humano", Edition: "1e", Description: "A raça mais versátil e ambiciosa.", Speed: 30, IsDefault: true},
        {Name: "Elfo", Edition: "1e", Description: "Seres graciosos com afinidade pela magia.", Speed: 30, IsDefault: true},
        {Name: "Anão", Edition: "1e", Description: "Seres robustos com resistência natural.", Speed: 25, IsDefault: true},
        {Name: "Halfling", Edition: "1e", Description: "Pequenos seres com grande sorte.", Speed: 25, IsDefault: true},
        {Name: "Gnomo", Edition: "1e", Description: "Inventores curiosos com afinidade arcana.", Speed: 25, IsDefault: true},
        {Name: "Meio-Elfo", Edition: "1e", Description: "Mistura de humano e elfo.", Speed: 30, IsDefault: true},
        {Name: "Meio-Orc", Edition: "1e", Description: "Mistura de humano e orc com força bruta.", Speed: 30, IsDefault: true},

        // D&D 2e
        {Name: "Humano", Edition: "2e", Description: "A raça mais versátil e ambiciosa.", Speed: 30, IsDefault: true},
        {Name: "Elfo", Edition: "2e", Description: "Seres graciosos com afinidade pela magia.", Speed: 30, IsDefault: true},
        {Name: "Elfo da Floresta", Edition: "2e", Description: "Elfos selvagens guardiões da floresta.", Speed: 30, IsDefault: true},
        {Name: "Elfo do Mar", Edition: "2e", Description: "Elfos com afinidade pelo oceano.", Speed: 30, IsDefault: true},
        {Name: "Anão das Montanhas", Edition: "2e", Description: "Anões fortes adaptados às montanhas.", Speed: 25, IsDefault: true},
        {Name: "Anão das Colinas", Edition: "2e", Description: "Anões resistentes das colinas.", Speed: 25, IsDefault: true},
        {Name: "Halfling", Edition: "2e", Description: "Pequenos seres com grande sorte.", Speed: 25, IsDefault: true},
        {Name: "Gnomo", Edition: "2e", Description: "Inventores curiosos com afinidade arcana.", Speed: 25, IsDefault: true},
        {Name: "Meio-Elfo", Edition: "2e", Description: "Mistura de humano e elfo.", Speed: 30, IsDefault: true},
        {Name: "Meio-Orc", Edition: "2e", Description: "Mistura de humano e orc com força bruta.", Speed: 30, IsDefault: true},

        // D&D 3e
        {Name: "Humano", Edition: "3e", Description: "A raça mais versátil e ambiciosa.", Speed: 30, IsDefault: true},
        {Name: "Elfo", Edition: "3e", Description: "Seres graciosos com afinidade pela magia.", Speed: 30, IsDefault: true},
        {Name: "Anão", Edition: "3e", Description: "Seres robustos com resistência natural.", Speed: 20, IsDefault: true},
        {Name: "Halfling", Edition: "3e", Description: "Pequenos seres com grande sorte.", Speed: 20, IsDefault: true},
        {Name: "Gnomo", Edition: "3e", Description: "Inventores curiosos com afinidade arcana.", Speed: 20, IsDefault: true},
        {Name: "Meio-Elfo", Edition: "3e", Description: "Mistura de humano e elfo.", Speed: 30, IsDefault: true},
        {Name: "Meio-Orc", Edition: "3e", Description: "Mistura de humano e orc com força bruta.", Speed: 30, IsDefault: true},

        // D&D 3.5e
        {Name: "Humano", Edition: "3.5e", Description: "A raça mais versátil e ambiciosa.", Speed: 30, IsDefault: true},
        {Name: "Elfo", Edition: "3.5e", Description: "Seres graciosos com afinidade pela magia.", Speed: 30, IsDefault: true},
        {Name: "Anão", Edition: "3.5e", Description: "Seres robustos com resistência natural.", Speed: 20, IsDefault: true},
        {Name: "Halfling", Edition: "3.5e", Description: "Pequenos seres com grande sorte.", Speed: 20, IsDefault: true},
        {Name: "Gnomo", Edition: "3.5e", Description: "Inventores curiosos com afinidade arcana.", Speed: 20, IsDefault: true},
        {Name: "Meio-Elfo", Edition: "3.5e", Description: "Mistura de humano e elfo.", Speed: 30, IsDefault: true},
        {Name: "Meio-Orc", Edition: "3.5e", Description: "Mistura de humano e orc com força bruta.", Speed: 30, IsDefault: true},
        {Name: "Tiefling", Edition: "3.5e", Description: "Descendente de demônios com traços infernais.", Speed: 30, IsDefault: true},

        // D&D 4e
        {Name: "Humano", Edition: "4e", Description: "A raça mais versátil e ambiciosa.", Speed: 30, IsDefault: true},
        {Name: "Elfo", Edition: "4e", Description: "Seres graciosos com afinidade pela magia.", Speed: 35, IsDefault: true},
        {Name: "Eladrin", Edition: "4e", Description: "Elfos feéricos com teleporte natural.", Speed: 30, IsDefault: true},
        {Name: "Anão", Edition: "4e", Description: "Seres robustos com resistência natural.", Speed: 25, IsDefault: true},
        {Name: "Halfling", Edition: "4e", Description: "Pequenos seres com grande sorte.", Speed: 30, IsDefault: true},
        {Name: "Meio-Elfo", Edition: "4e", Description: "Mistura de humano e elfo.", Speed: 30, IsDefault: true},
        {Name: "Tiefling", Edition: "4e", Description: "Descendente de demônios com traços infernais.", Speed: 30, IsDefault: true},
        {Name: "Draconato", Edition: "4e", Description: "Humanoides com herança dracônica.", Speed: 30, IsDefault: true},

        // D&D 5e
        {Name: "Humano", Edition: "5e", Description: "A raça mais versátil e ambiciosa.", Speed: 30, IsDefault: true},
        {Name: "Elfo", Edition: "5e", Description: "Seres graciosos com afinidade pela magia.", Speed: 30, IsDefault: true},
        {Name: "Elfo da Floresta", Edition: "5e", Description: "Elfos ágeis com visão apurada.", Speed: 35, IsDefault: true},
        {Name: "Elfo Negro (Drow)", Edition: "5e", Description: "Elfos das profundezas com magia das trevas.", Speed: 30, IsDefault: true},
        {Name: "Alto Elfo", Edition: "5e", Description: "Elfos com talento arcano natural.", Speed: 30, IsDefault: true},
        {Name: "Anão das Montanhas", Edition: "5e", Description: "Anões fortes adaptados às montanhas.", Speed: 25, IsDefault: true},
        {Name: "Anão das Colinas", Edition: "5e", Description: "Anões resistentes com grande vitalidade.", Speed: 25, IsDefault: true},
        {Name: "Halfling Pés-Leves", Edition: "5e", Description: "Halflings ágeis e furtivos.", Speed: 25, IsDefault: true},
        {Name: "Halfling Robusto", Edition: "5e", Description: "Halflings resistentes a venenos.", Speed: 25, IsDefault: true},
        {Name: "Humano Variante", Edition: "5e", Description: "Humanos com talento especial.", Speed: 30, IsDefault: true},
        {Name: "Draconato", Edition: "5e", Description: "Humanoides com herança dracônica.", Speed: 30, IsDefault: true},
        {Name: "Gnomo da Floresta", Edition: "5e", Description: "Gnomos com ilusões naturais.", Speed: 25, IsDefault: true},
        {Name: "Gnomo das Rochas", Edition: "5e", Description: "Gnomos inventores com resistência mágica.", Speed: 25, IsDefault: true},
        {Name: "Meio-Elfo", Edition: "5e", Description: "Mistura de humano e elfo.", Speed: 30, IsDefault: true},
        {Name: "Meio-Orc", Edition: "5e", Description: "Mistura de humano e orc com força bruta.", Speed: 30, IsDefault: true},
        {Name: "Tiefling", Edition: "5e", Description: "Descendente de demônios com traços infernais.", Speed: 30, IsDefault: true},
        {Name: "Aasimar", Edition: "5e", Description: "Descendente de celestiais com luz divina.", Speed: 30, IsDefault: true},
        {Name: "Firbolg", Edition: "5e", Description: "Gigantes gentis guardiões da natureza.", Speed: 30, IsDefault: true},
        {Name: "Goblin", Edition: "5e", Description: "Pequenos seres ágeis e astutos.", Speed: 30, IsDefault: true},
        {Name: "Hobgoblin", Edition: "5e", Description: "Guerreiros disciplinados de origem goblinoide.", Speed: 30, IsDefault: true},
        {Name: "Kenku", Edition: "5e", Description: "Humanoides aviários sem asas.", Speed: 30, IsDefault: true},
        {Name: "Kobold", Edition: "5e", Description: "Pequenos draconatos de pack tactics.", Speed: 30, IsDefault: true},
        {Name: "Lizardfolk", Edition: "5e", Description: "Humanoides répteis pragmáticos.", Speed: 30, IsDefault: true},
        {Name: "Tabaxi", Edition: "5e", Description: "Humanoides felinos curiosos e ágeis.", Speed: 30, IsDefault: true},
        {Name: "Triton", Edition: "5e", Description: "Guardiões do oceano profundo.", Speed: 30, IsDefault: true},
        {Name: "Yuan-ti Pureblood", Edition: "5e", Description: "Humanoides serpentinos com resistência mágica.", Speed: 30, IsDefault: true},
    }

    for _, race := range races {
        var existing domain.Race
        result := db.Where("name = ? AND edition = ? AND is_default = ?", race.Name, race.Edition, true).First(&existing)
        if result.Error != nil {
            db.Create(&race)
        }
    }

    log.Println("Raças populadas com sucesso!")

    
}

func seedArmors(db *gorm.DB) {
    armors := []domain.Armor{

        // ========== D&D 1e e 2e ==========
        // Nas edições 1e/2e a CA é decrescente (menor = melhor)
        // Guardamos o valor base e calculamos inversamente
        {Name: "Sem Armadura", Edition: "1e", ArmorType: domain.ArmorNone, BaseAC: 10, MaxDexBonus: -1, IsDefault: true, Description: "Sem proteção. CA base 10 (sistema decrescente)."},
        {Name: "Gambeson", Edition: "1e", ArmorType: domain.ArmorLight, BaseAC: 8, MaxDexBonus: -1, IsDefault: true, Description: "Armadura acolchoada leve."},
        {Name: "Couro", Edition: "1e", ArmorType: domain.ArmorLight, BaseAC: 7, MaxDexBonus: -1, IsDefault: true, Description: "Armadura de couro endurecido."},
        {Name: "Couro Batido", Edition: "1e", ArmorType: domain.ArmorMedium, BaseAC: 6, MaxDexBonus: -1, IsDefault: true, Description: "Couro reforçado com placas metálicas."},
        {Name: "Cota de Malha", Edition: "1e", ArmorType: domain.ArmorMedium, BaseAC: 5, MaxDexBonus: -1, IsDefault: true, Description: "Anéis de metal entrelaçados."},
        {Name: "Brunea", Edition: "1e", ArmorType: domain.ArmorMedium, BaseAC: 4, MaxDexBonus: -1, IsDefault: true, Description: "Armadura de escamas metálicas."},
        {Name: "Lamelar", Edition: "1e", ArmorType: domain.ArmorHeavy, BaseAC: 3, MaxDexBonus: -1, IsDefault: true, Description: "Placas sobrepostas de metal."},
        {Name: "Plate Mail", Edition: "1e", ArmorType: domain.ArmorHeavy, BaseAC: 2, MaxDexBonus: -1, IsDefault: true, Description: "Armadura completa de placas de metal."},
        {Name: "Escudo", Edition: "1e", ArmorType: domain.ArmorShield, BaseAC: -1, MaxDexBonus: -1, IsDefault: true, Description: "Reduz a CA em 1 (sistema decrescente)."},

        {Name: "Sem Armadura", Edition: "2e", ArmorType: domain.ArmorNone, BaseAC: 10, MaxDexBonus: -1, IsDefault: true, Description: "Sem proteção. CA base 10 (sistema decrescente)."},
        {Name: "Gambeson", Edition: "2e", ArmorType: domain.ArmorLight, BaseAC: 8, MaxDexBonus: -1, IsDefault: true, Description: "Armadura acolchoada leve."},
        {Name: "Couro", Edition: "2e", ArmorType: domain.ArmorLight, BaseAC: 7, MaxDexBonus: -1, IsDefault: true, Description: "Armadura de couro endurecido."},
        {Name: "Couro Batido", Edition: "2e", ArmorType: domain.ArmorMedium, BaseAC: 6, MaxDexBonus: -1, IsDefault: true, Description: "Couro reforçado com placas metálicas."},
        {Name: "Cota de Malha", Edition: "2e", ArmorType: domain.ArmorMedium, BaseAC: 5, MaxDexBonus: -1, IsDefault: true, Description: "Anéis de metal entrelaçados."},
        {Name: "Brunea", Edition: "2e", ArmorType: domain.ArmorMedium, BaseAC: 4, MaxDexBonus: -1, IsDefault: true, Description: "Armadura de escamas metálicas."},
        {Name: "Lamelar", Edition: "2e", ArmorType: domain.ArmorHeavy, BaseAC: 3, MaxDexBonus: -1, IsDefault: true, Description: "Placas sobrepostas de metal."},
        {Name: "Plate Mail", Edition: "2e", ArmorType: domain.ArmorHeavy, BaseAC: 2, MaxDexBonus: -1, IsDefault: true, Description: "Armadura completa de placas de metal."},
        {Name: "Escudo", Edition: "2e", ArmorType: domain.ArmorShield, BaseAC: -1, MaxDexBonus: -1, IsDefault: true, Description: "Reduz a CA em 1 (sistema decrescente)."},

        // ========== D&D 3e e 3.5e ==========
        {Name: "Sem Armadura", Edition: "3e", ArmorType: domain.ArmorNone, BaseAC: 10, MaxDexBonus: -1, IsDefault: true, Description: "CA = 10 + mod DEX."},
        {Name: "Gambeson", Edition: "3e", ArmorType: domain.ArmorLight, BaseAC: 11, MaxDexBonus: -1, IsDefault: true, Description: "Armadura acolchoada. CA = 11 + mod DEX."},
        {Name: "Couro", Edition: "3e", ArmorType: domain.ArmorLight, BaseAC: 12, MaxDexBonus: -1, IsDefault: true, Description: "Couro endurecido. CA = 12 + mod DEX."},
        {Name: "Couro Batido", Edition: "3e", ArmorType: domain.ArmorLight, BaseAC: 13, MaxDexBonus: 5, IsDefault: true, Description: "Couro reforçado. CA = 13 + mod DEX (máx +5)."},
        {Name: "Cota de Malha", Edition: "3e", ArmorType: domain.ArmorMedium, BaseAC: 13, MaxDexBonus: 2, IsDefault: true, Description: "Anéis de metal. CA = 13 + mod DEX (máx +2)."},
        {Name: "Escamas", Edition: "3e", ArmorType: domain.ArmorMedium, BaseAC: 14, MaxDexBonus: 3, IsDefault: true, Description: "Escamas metálicas sobrepostas."},
        {Name: "Brunea", Edition: "3e", ArmorType: domain.ArmorMedium, BaseAC: 15, MaxDexBonus: 2, IsDefault: true, Description: "Armadura de placas metálicas médias."},
        {Name: "Meia Armadura", Edition: "3e", ArmorType: domain.ArmorHeavy, BaseAC: 15, MaxDexBonus: 0, IsDefault: true, Description: "Armadura pesada parcial."},
        {Name: "Cota de Placas", Edition: "3e", ArmorType: domain.ArmorHeavy, BaseAC: 17, MaxDexBonus: 0, IsDefault: true, Description: "Armadura de placas completa."},
        {Name: "Armadura Completa", Edition: "3e", ArmorType: domain.ArmorHeavy, BaseAC: 18, MaxDexBonus: 0, IsDefault: true, Description: "Proteção máxima em placas de metal."},
        {Name: "Escudo", Edition: "3e", ArmorType: domain.ArmorShield, BaseAC: 2, MaxDexBonus: -1, IsDefault: true, Description: "Adiciona +2 à CA."},

        {Name: "Sem Armadura", Edition: "3.5e", ArmorType: domain.ArmorNone, BaseAC: 10, MaxDexBonus: -1, IsDefault: true, Description: "CA = 10 + mod DEX."},
        {Name: "Gambeson", Edition: "3.5e", ArmorType: domain.ArmorLight, BaseAC: 11, MaxDexBonus: -1, IsDefault: true, Description: "Armadura acolchoada. CA = 11 + mod DEX."},
        {Name: "Couro", Edition: "3.5e", ArmorType: domain.ArmorLight, BaseAC: 12, MaxDexBonus: -1, IsDefault: true, Description: "Couro endurecido. CA = 12 + mod DEX."},
        {Name: "Couro Batido", Edition: "3.5e", ArmorType: domain.ArmorLight, BaseAC: 13, MaxDexBonus: 5, IsDefault: true, Description: "Couro reforçado. CA = 13 + mod DEX (máx +5)."},
        {Name: "Cota de Malha", Edition: "3.5e", ArmorType: domain.ArmorMedium, BaseAC: 13, MaxDexBonus: 2, IsDefault: true, Description: "Anéis de metal. CA = 13 + mod DEX (máx +2)."},
        {Name: "Escamas", Edition: "3.5e", ArmorType: domain.ArmorMedium, BaseAC: 14, MaxDexBonus: 3, IsDefault: true, Description: "Escamas metálicas sobrepostas."},
        {Name: "Brunea", Edition: "3.5e", ArmorType: domain.ArmorMedium, BaseAC: 15, MaxDexBonus: 2, IsDefault: true, Description: "Armadura de placas metálicas médias."},
        {Name: "Meia Armadura", Edition: "3.5e", ArmorType: domain.ArmorHeavy, BaseAC: 15, MaxDexBonus: 0, IsDefault: true, Description: "Armadura pesada parcial."},
        {Name: "Cota de Placas", Edition: "3.5e", ArmorType: domain.ArmorHeavy, BaseAC: 17, MaxDexBonus: 0, IsDefault: true, Description: "Armadura de placas completa."},
        {Name: "Armadura Completa", Edition: "3.5e", ArmorType: domain.ArmorHeavy, BaseAC: 18, MaxDexBonus: 0, IsDefault: true, Description: "Proteção máxima em placas de metal."},
        {Name: "Escudo", Edition: "3.5e", ArmorType: domain.ArmorShield, BaseAC: 2, MaxDexBonus: -1, IsDefault: true, Description: "Adiciona +2 à CA."},

        // ========== D&D 4e ==========
        {Name: "Sem Armadura", Edition: "4e", ArmorType: domain.ArmorNone, BaseAC: 10, MaxDexBonus: -1, IsDefault: true, Description: "CA = 10 + metade do nível + mod DEX ou INT."},
        {Name: "Couro", Edition: "4e", ArmorType: domain.ArmorLight, BaseAC: 12, MaxDexBonus: -1, IsDefault: true, Description: "CA = 12 + metade do nível + mod DEX ou INT."},
        {Name: "Couro Escondido", Edition: "4e", ArmorType: domain.ArmorLight, BaseAC: 13, MaxDexBonus: -1, IsDefault: true, Description: "Couro reforçado discreto."},
        {Name: "Cota de Malha", Edition: "4e", ArmorType: domain.ArmorMedium, BaseAC: 14, MaxDexBonus: 0, IsDefault: true, Description: "CA = 14 + metade do nível. Sem bônus DEX."},
        {Name: "Escamas", Edition: "4e", ArmorType: domain.ArmorMedium, BaseAC: 15, MaxDexBonus: 0, IsDefault: true, Description: "CA = 15 + metade do nível. Sem bônus DEX."},
        {Name: "Armadura de Placas", Edition: "4e", ArmorType: domain.ArmorHeavy, BaseAC: 17, MaxDexBonus: 0, IsDefault: true, Description: "CA = 17 + metade do nível. Sem bônus DEX."},
        {Name: "Armadura Completa", Edition: "4e", ArmorType: domain.ArmorHeavy, BaseAC: 18, MaxDexBonus: 0, IsDefault: true, Description: "CA = 18 + metade do nível. Sem bônus DEX."},
        {Name: "Escudo Leve", Edition: "4e", ArmorType: domain.ArmorShield, BaseAC: 1, MaxDexBonus: -1, IsDefault: true, Description: "Adiciona +1 à CA."},
        {Name: "Escudo Pesado", Edition: "4e", ArmorType: domain.ArmorShield, BaseAC: 2, MaxDexBonus: -1, IsDefault: true, Description: "Adiciona +2 à CA."},

        // ========== D&D 5e ==========
        {Name: "Sem Armadura", Edition: "5e", ArmorType: domain.ArmorNone, BaseAC: 10, MaxDexBonus: -1, IsDefault: true, Description: "CA = 10 + mod DEX."},
        {Name: "Gambeson", Edition: "5e", ArmorType: domain.ArmorLight, BaseAC: 11, MaxDexBonus: -1, IsDefault: true, Description: "CA = 11 + mod DEX."},
        {Name: "Couro", Edition: "5e", ArmorType: domain.ArmorLight, BaseAC: 11, MaxDexBonus: -1, IsDefault: true, Description: "CA = 11 + mod DEX."},
        {Name: "Couro Batido", Edition: "5e", ArmorType: domain.ArmorLight, BaseAC: 12, MaxDexBonus: -1, IsDefault: true, Description: "CA = 12 + mod DEX."},
        {Name: "Cota de Malha", Edition: "5e", ArmorType: domain.ArmorMedium, BaseAC: 13, MaxDexBonus: 2, IsDefault: true, Description: "CA = 13 + mod DEX (máx +2)."},
        {Name: "Couro Escondido", Edition: "5e", ArmorType: domain.ArmorMedium, BaseAC: 12, MaxDexBonus: 2, IsDefault: true, Description: "CA = 12 + mod DEX (máx +2)."},
        {Name: "Escamas", Edition: "5e", ArmorType: domain.ArmorMedium, BaseAC: 14, MaxDexBonus: 2, IsDefault: true, Description: "CA = 14 + mod DEX (máx +2)."},
        {Name: "Meia Armadura", Edition: "5e", ArmorType: domain.ArmorMedium, BaseAC: 15, MaxDexBonus: 2, IsDefault: true, Description: "CA = 15 + mod DEX (máx +2)."},
        {Name: "Brunea", Edition: "5e", ArmorType: domain.ArmorMedium, BaseAC: 14, MaxDexBonus: 2, IsDefault: true, Description: "CA = 14 + mod DEX (máx +2)."},
        {Name: "Cota de Anéis", Edition: "5e", ArmorType: domain.ArmorHeavy, BaseAC: 14, MaxDexBonus: 0, IsDefault: true, Description: "CA = 14. Sem bônus DEX."},
        {Name: "Cota de Placas", Edition: "5e", ArmorType: domain.ArmorHeavy, BaseAC: 16, MaxDexBonus: 0, IsDefault: true, Description: "CA = 16. Sem bônus DEX."},
        {Name: "Meia Placa", Edition: "5e", ArmorType: domain.ArmorHeavy, BaseAC: 15, MaxDexBonus: 0, IsDefault: true, Description: "CA = 15. Sem bônus DEX."},
        {Name: "Armadura Completa", Edition: "5e", ArmorType: domain.ArmorHeavy, BaseAC: 18, MaxDexBonus: 0, IsDefault: true, Description: "CA = 18. Sem bônus DEX."},
        {Name: "Escudo", Edition: "5e", ArmorType: domain.ArmorShield, BaseAC: 2, MaxDexBonus: -1, IsDefault: true, Description: "Adiciona +2 à CA."},
    }

    for _, armor := range armors {
        var existing domain.Armor
        result := db.Where("name = ? AND edition = ? AND is_default = ?", armor.Name, armor.Edition, true).First(&existing)
        if result.Error != nil {
            db.Create(&armor)
        }
    }

    log.Println("Armaduras populadas com sucesso!")
}