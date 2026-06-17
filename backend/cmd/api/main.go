package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"rpg-manager/internal/handler"
	"rpg-manager/internal/repository"
	"rpg-manager/internal/service"
	"rpg-manager/pkg/config"
	"rpg-manager/pkg/middleware"
	"rpg-manager/internal/seed"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variáveis do sistema")
	}

	config.ConnectDatabase()
	seed.Run(config.DB)

	// Repositories
	userRepo      := repository.NewUserRepository(config.DB)
	classRepo     := repository.NewClassRepository(config.DB)
	raceRepo      := repository.NewRaceRepository(config.DB)
	skillRepo     := repository.NewSkillRepository(config.DB)
	characterRepo := repository.NewCharacterRepository(config.DB)
	backgroundRepo := repository.NewBackgroundRepository(config.DB)
	armorRepo     := repository.NewArmorRepository(config.DB)
	periciaRepo   := repository.NewPericiaRepository(config.DB)   // NOVO
	talentoRepo   := repository.NewTalentoRepository(config.DB)   // NOVO

	// Services
	authService      := service.NewAuthService(userRepo)
	classService     := service.NewClassService(classRepo)
	raceService      := service.NewRaceService(raceRepo)
	skillService     := service.NewSkillService(skillRepo)
	characterService := service.NewCharacterService(characterRepo, skillRepo)
	backgroundService := service.NewBackgroundService(backgroundRepo)
	armorService     := service.NewArmorService(armorRepo)
	periciaService   := service.NewPericiaService(periciaRepo)    // NOVO
	talentoService   := service.NewTalentoService(talentoRepo)    // NOVO

	// Handlers
	authHandler      := handler.NewAuthHandler(authService)
	classHandler     := handler.NewClassHandler(classService)
	raceHandler      := handler.NewRaceHandler(raceService)
	skillHandler     := handler.NewSkillHandler(skillService)
	characterHandler := handler.NewCharacterHandler(characterService, armorService)
	backgroundHandler := handler.NewBackgroundHandler(backgroundService)
	uploadHandler    := handler.NewUploadHandler(characterRepo)
	armorHandler     := handler.NewArmorHandler(armorService)
	ollamaHandler    := handler.NewOllamaHandler()
	periciaHandler   := handler.NewPericiaHandler(periciaService) // NOVO
	talentoHandler   := handler.NewTalentoHandler(talentoService) // NOVO

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := map[string]bool{
			"http://localhost:5173":               true,
			"https://rpg-manager-smoky.vercel.app": true,
		}
		if allowed[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "RPG Manager API rodando!"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		api.GET("/armors", armorHandler.GetByEdition)
		api.POST("/ai/skills", ollamaHandler.GetSkills)

		// Perícias e Talentos — públicos (para leitura no CharacterCreate)
		api.GET("/pericias", periciaHandler.GetAll)   // NOVO
		api.GET("/talentos", talentoHandler.GetAll)   // NOVO

		classes := api.Group("/classes")
		{
			classes.GET("", classHandler.GetAll)
			classes.GET("/:id", classHandler.GetByID)
		}
		races := api.Group("/races")
		{
			races.GET("", raceHandler.GetAll)
			races.GET("/:id", raceHandler.GetByID)
		}
		skills := api.Group("/skills")
		{
			skills.GET("", skillHandler.GetAll)
			skills.GET("/filter", skillHandler.GetByClassAndRace)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.POST("/classes", classHandler.Create)
			protected.PUT("/classes/:id", classHandler.Update)
			protected.DELETE("/classes/:id", classHandler.Delete)

			protected.POST("/races", raceHandler.Create)
			protected.PUT("/races/:id", raceHandler.Update)
			protected.DELETE("/races/:id", raceHandler.Delete)

			protected.POST("/skills", skillHandler.Create)
			protected.PUT("/skills/:id", skillHandler.Update)
			protected.DELETE("/skills/:id", skillHandler.Delete)

			characters := protected.Group("/characters")
			{
				characters.GET("", characterHandler.GetAll)
				characters.GET("/:id", characterHandler.GetByID)
				characters.POST("", characterHandler.Create)
				characters.PUT("/:id", characterHandler.Update)
				characters.DELETE("/:id", characterHandler.Delete)
				characters.PATCH("/:id/level-up", characterHandler.LevelUp)
				characters.POST("/:id/skills/:skill_id", characterHandler.AddSkill)
				characters.DELETE("/:id/skills/:skill_id", characterHandler.RemoveSkill)
				characters.GET("/:id/background", backgroundHandler.Get)
				characters.POST("/:id/background", backgroundHandler.Save)
				characters.POST("/:id/avatar", uploadHandler.UploadAvatar)
				characters.PATCH("/:id/take-damage", characterHandler.TakeDamage)
				characters.PATCH("/:id/heal", characterHandler.Heal)
				characters.PATCH("/:id/temp-hp", characterHandler.AddTempHP)
				characters.GET("/:id/ac", characterHandler.GetAC)

				// NOVO — Perícias e Talentos do personagem
				characters.GET("/:id/pericias", periciaHandler.GetByCharacter)
				characters.POST("/:id/pericias", periciaHandler.Save)
				characters.GET("/:id/talentos", talentoHandler.GetByCharacter)
				characters.POST("/:id/talentos/:talento_id", talentoHandler.Add)
				characters.DELETE("/:id/talentos/:talento_id", talentoHandler.Remove)
			}
		}
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	fmt.Println("🎲 Servidor rodando na porta 8080")
	srv.ListenAndServe()
}