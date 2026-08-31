package config

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rpg-manager/internal/domain"
)

var DB *gorm.DB

func ConnectDatabase() {
    ssl := os.Getenv("DB_SSL")
    if ssl == "" {
        ssl = "disable" // local: disable | Render: configure DB_SSL=require
    }

    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        ssl,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Erro ao conectar ao banco de dados: ", err)
    }

    err = db.AutoMigrate(
        &domain.User{},
        &domain.Class{},
        &domain.Race{},
        &domain.Armor{},
        &domain.Skill{},
        &domain.Character{},
        &domain.Background{},
        &domain.Antecedent{}, // ← NOVO: cria tabela antecedents
        &domain.Pericia{},
        &domain.CharacterPericia{},
        &domain.Talento{},
        &domain.Spell{},
        &domain.Item{},
        &domain.CharacterItem{},
        &domain.CharacterArmorOwned{},
        &domain.ClassEquipmentOption{},
        &domain.ClassEquipmentComponent{},
        &domain.Language{},
        // ── Sistema do Mestre (mesa virtual) ──────────────────────────────
        &domain.Campaign{},
        &domain.CampaignMembership{},
        &domain.NPC{},
        &domain.Enemy{},
        &domain.EnemyAbility{},
        &domain.EnemyLine{},
        &domain.MagicItem{},
        &domain.Scene{},
        &domain.Token{},
        &domain.Session{},
        &domain.ChatMessage{},
        &domain.Reward{},
    )
    if err != nil {
        log.Fatal("Erro ao executar AutoMigrate: ", err)
    }

    log.Println("Banco de dados conectado com sucesso!")
    DB = db
}