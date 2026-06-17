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
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
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
        &domain.Pericia{},
        &domain.CharacterPericia{},
        &domain.Talento{},
    )
    if err != nil {
        log.Fatal("Erro ao executar AutoMigrate: ", err)
    }

    // Popula o banco com dados padrão
    SeedDatabase(db)

    log.Println("Banco de dados conectado com sucesso!")
    DB = db
}