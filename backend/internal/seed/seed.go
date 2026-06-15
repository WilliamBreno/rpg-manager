package seed

import (
    "log"
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

func Run(db *gorm.DB) {
    log.Println("🌱 Rodando seed...")

    classes := []domain.Class{
        // 4e
        {Name: "Bardo",       Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7,  ReflBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Bárbaro",     Edition: "4e", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 8,  FortBonus: 2,               IsDefault: true},
        {Name: "Clérigo",     Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6,                WillBonus: 1, IsDefault: true},
        {Name: "Druida",      Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7,  ReflBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Feiticeiro",  Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6,                WillBonus: 2, IsDefault: true},
        {Name: "Guardião",    Edition: "4e", BaseHP: 17, HPPerLevel: 7, SurgesPerDay: 9,  FortBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Guerreiro",   Edition: "4e", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 9,  FortBonus: 2,               IsDefault: true},
        {Name: "Invocador",   Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6,                WillBonus: 2, IsDefault: true},
        {Name: "Ladino",      Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6,  ReflBonus: 2,               IsDefault: true},
        {Name: "Mago",        Edition: "4e", BaseHP: 10, HPPerLevel: 4, SurgesPerDay: 6,                WillBonus: 2, IsDefault: true},
        {Name: "Monge",       Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7,  FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Paladino",    Edition: "4e", BaseHP: 15, HPPerLevel: 6, SurgesPerDay: 10, FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Patrulheiro", Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 6,  FortBonus: 1, ReflBonus: 1,               IsDefault: true},
        {Name: "Psionista",   Edition: "4e", BaseHP: 12, HPPerLevel: 4, SurgesPerDay: 6,                WillBonus: 2, IsDefault: true},
        {Name: "Rastreador",  Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7,  ReflBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Vingador",    Edition: "4e", BaseHP: 14, HPPerLevel: 6, SurgesPerDay: 7,  FortBonus: 1, ReflBonus: 1, WillBonus: 1, IsDefault: true},
        {Name: "Xamã",        Edition: "4e", BaseHP: 12, HPPerLevel: 5, SurgesPerDay: 7,  FortBonus: 1, WillBonus: 1, IsDefault: true},

        // 5e
        {Name: "Bárbaro",     Edition: "5e", HitDie: 12, IsDefault: true},
        {Name: "Bardo",       Edition: "5e", HitDie: 8,  IsDefault: true},
        {Name: "Bruxo",       Edition: "5e", HitDie: 8,  IsDefault: true},
        {Name: "Clérigo",     Edition: "5e", HitDie: 8,  IsDefault: true},
        {Name: "Druida",      Edition: "5e", HitDie: 8,  IsDefault: true},
        {Name: "Feiticeiro",  Edition: "5e", HitDie: 6,  IsDefault: true},
        {Name: "Guerreiro",   Edition: "5e", HitDie: 10, IsDefault: true},
        {Name: "Ladino",      Edition: "5e", HitDie: 8,  IsDefault: true},
        {Name: "Mago",        Edition: "5e", HitDie: 6,  IsDefault: true},
        {Name: "Monge",       Edition: "5e", HitDie: 8,  IsDefault: true},
        {Name: "Paladino",    Edition: "5e", HitDie: 10, IsDefault: true},
        {Name: "Patrulheiro", Edition: "5e", HitDie: 10, IsDefault: true},
        {Name: "Xamã",        Edition: "5e", HitDie: 8,  IsDefault: true},
    }

    for _, class := range classes {
        result := db.Where(domain.Class{Name: class.Name, Edition: class.Edition}).FirstOrCreate(&class)
        if result.Error != nil {
            log.Printf("Erro ao seed classe %s (%s): %v", class.Name, class.Edition, result.Error)
        }
    }

    log.Println("✅ Seed concluído!")
}