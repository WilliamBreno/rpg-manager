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
	userRepo       := repository.NewUserRepository(config.DB)
	classRepo      := repository.NewClassRepository(config.DB)
	raceRepo       := repository.NewRaceRepository(config.DB)
	skillRepo      := repository.NewSkillRepository(config.DB)
	characterRepo  := repository.NewCharacterRepository(config.DB)
	backgroundRepo := repository.NewBackgroundRepository(config.DB)
	armorRepo      := repository.NewArmorRepository(config.DB)
	periciaRepo    := repository.NewPericiaRepository(config.DB)
	talentoRepo    := repository.NewTalentoRepository(config.DB)
	spellRepo      := repository.NewSpellRepository(config.DB)
	itemRepo       := repository.NewItemRepository(config.DB)

	_ = backgroundRepo // usado indiretamente via backgroundService

	// Services
	authService       := service.NewAuthService(userRepo)
	classService      := service.NewClassService(classRepo)
	raceService       := service.NewRaceService(raceRepo)
	skillService      := service.NewSkillService(skillRepo)
	characterService  := service.NewCharacterService(characterRepo, skillRepo, classRepo, raceRepo, talentoRepo, config.DB)
	backgroundService := service.NewBackgroundService(config.DB)
	antecedentSvc     := service.NewAntecedentService(config.DB)
	armorService      := service.NewArmorService(armorRepo)
	periciaService    := service.NewPericiaService(periciaRepo)
	talentoService    := service.NewTalentoService(talentoRepo)
	spellService      := service.NewSpellService(spellRepo)
	itemService       := service.NewItemService(itemRepo)
	inventoryService  := service.NewInventoryService(config.DB)

	// Handlers
	antecedentHandler := handler.NewAntecedentHandler(antecedentSvc)
	authHandler       := handler.NewAuthHandler(authService)
	classHandler      := handler.NewClassHandler(classService)
	raceHandler       := handler.NewRaceHandler(raceService)
	skillHandler      := handler.NewSkillHandler(skillService)
	characterHandler  := handler.NewCharacterHandler(characterService, armorService, periciaService, userRepo)
	backgroundHandler := handler.NewBackgroundHandler(backgroundService)
	uploadHandler     := handler.NewUploadHandler(characterRepo)
	armorHandler      := handler.NewArmorHandler(armorService)
	ollamaHandler     := handler.NewOllamaHandler()
	periciaHandler    := handler.NewPericiaHandler(periciaService)
	talentoHandler    := handler.NewTalentoHandler(talentoService)
	spellHandler      := handler.NewSpellHandler(spellService)
	itemHandler       := handler.NewItemHandler(itemService)
	inventoryHandler  := handler.NewInventoryHandler(inventoryService, characterService)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := map[string]bool{
			"http://localhost:5173":                true,
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

		// Rotas públicas
		api.GET("/armors", armorHandler.GetByEdition)
		api.GET("/items", itemHandler.GetAll)
		api.GET("/pericias", periciaHandler.GetAll)
		api.GET("/talentos", talentoHandler.GetAll)
		api.GET("/spells", spellHandler.GetAll)
		api.GET("/antecedentes", antecedentHandler.GetAll)
		api.GET("/antecedentes/:id", antecedentHandler.GetByID)
		api.POST("/ai/skills", ollamaHandler.GetSkills)
		api.PATCH("/characters/:id/add-xp",    characterHandler.AddXP)
		api.PATCH("/characters/:id/apply-asi", characterHandler.ApplyASI)
		api.PATCH("/characters/:id/death-save",        characterHandler.DeathSave)
		api.PATCH("/characters/:id/reset-death-saves", characterHandler.ResetDeathSaves)
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

		// Rotas protegidas (JWT)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.PATCH("/users/me/welcome-seen", authHandler.MarkWelcomeSeen)

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
				characters.GET("/:id/export/pdf", characterHandler.ExportPDF5e)
				characters.GET("/:id/pericias", periciaHandler.GetByCharacter)
				characters.POST("/:id/pericias", periciaHandler.Save)
				characters.GET("/:id/talentos", talentoHandler.GetByCharacter)
				characters.POST("/:id/talentos/:talento_id", talentoHandler.Add)
				characters.DELETE("/:id/talentos/:talento_id", talentoHandler.Remove)
				characters.GET("/:id/spells", spellHandler.GetByCharacter)
				characters.POST("/:id/spells/:spell_id", spellHandler.Add)
				characters.DELETE("/:id/spells/:spell_id", spellHandler.Remove)
				characters.GET("/:id/inventory", inventoryHandler.GetInventory)
				characters.POST("/:id/shop/items/:item_id", inventoryHandler.PurchaseItem)
				characters.POST("/:id/shop/armors/:armor_id", inventoryHandler.PurchaseArmor)
				characters.PATCH("/:id/currency", inventoryHandler.SetCurrency)
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