package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/pdfexport"
	"rpg-manager/internal/repository"
	"rpg-manager/internal/service"
)

type CharacterHandler struct {
	Service        *service.CharacterService
	ArmorService   *service.ArmorService
	PericiaService *service.PericiaService
	UserRepo       *repository.UserRepository
}

func NewCharacterHandler(service *service.CharacterService, armorService *service.ArmorService, periciaService *service.PericiaService, userRepo *repository.UserRepository) *CharacterHandler {
	return &CharacterHandler{Service: service, ArmorService: armorService, PericiaService: periciaService, UserRepo: userRepo}
}

func (h *CharacterHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("userID")
	characters, err := h.Service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, characters)
}

func (h *CharacterHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}
	c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) Create(c *gin.Context) {
	var character domain.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character.UserID = c.GetUint("userID")
	if err := h.Service.Create(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, character)
}

func (h *CharacterHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var character domain.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character.ID = uint(id)
	character.UserID = c.GetUint("userID")
	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if err := h.Service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Personagem deletado com sucesso"})
}

// ── XP / Level Up ─────────────────────────────────────────────────────────────

func (h *CharacterHandler) AddXP(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var body struct {
		XP int `json:"xp"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, result, err := h.Service.AddXP(uint(id), body.XP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"character":  character,
		"leveled_up": result.LeveledUp,
		"needs_asi":  result.NeedsASI,
		"new_level":  result.NewLevel,
	})
}

func (h *CharacterHandler) ApplyASI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var choice service.ASIChoice
	if err := c.ShouldBindJSON(&choice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, result, err := h.Service.ApplyASI(uint(id), choice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"character":  character,
		"leveled_up": result.LeveledUp,
		"needs_asi":  result.NeedsASI,
		"new_level":  result.NewLevel,
	})
}

func (h *CharacterHandler) LevelUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	character, err := h.Service.LevelUp(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

// ── Skills ────────────────────────────────────────────────────────────────────

func (h *CharacterHandler) AddSkill(c *gin.Context) {
	characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	skillID, err := strconv.ParseUint(c.Param("skill_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "skill_id inválido"})
		return
	}
	if err := h.Service.AddSkill(uint(characterID), uint(skillID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Habilidade adicionada com sucesso"})
}

func (h *CharacterHandler) RemoveSkill(c *gin.Context) {
	characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	skillID, err := strconv.ParseUint(c.Param("skill_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "skill_id inválido"})
		return
	}
	if err := h.Service.RemoveSkill(uint(characterID), uint(skillID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Habilidade removida com sucesso"})
}

// ── HP Management ─────────────────────────────────────────────────────────────

func (h *CharacterHandler) GetAC(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}
	ac := h.ArmorService.CalculateAC(character)
	c.JSON(http.StatusOK, gin.H{"ac": ac})
}

func (h *CharacterHandler) TakeDamage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var body struct {
		Damage int `json:"damage"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}
	// Desconta TempHP primeiro
	if character.TempHP > 0 {
		if body.Damage <= character.TempHP {
			character.TempHP -= body.Damage
			body.Damage = 0
		} else {
			body.Damage -= character.TempHP
			character.TempHP = 0
		}
	}
	character.HitPoints -= body.Damage
	if character.HitPoints < 0 {
		character.HitPoints = 0
	}
	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

// Heal — reseta testes de morte automaticamente quando HP sobe acima de 0 (5e)
func (h *CharacterHandler) Heal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var body struct {
		Amount int `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}
	wasAtZero := character.HitPoints == 0
	character.HitPoints += body.Amount
	if character.HitPoints > character.MaxHP {
		character.HitPoints = character.MaxHP
	}
	// 5e: se o personagem estava em 0 HP e foi curado, reseta os testes de morte
	if wasAtZero && character.HitPoints > 0 && character.Edition == "5e" {
		character.DeathSaveSuccesses = 0
		character.DeathSaveFailures = 0
	}
	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) AddTempHP(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var body struct {
		Amount int `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}
	// Temp HP não acumula — fica com o maior
	if body.Amount > character.TempHP {
		character.TempHP = body.Amount
	}
	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

// ── Death Saving Throws (5e) ──────────────────────────────────────────────────

// DeathSave — registra o resultado de um teste de morte
// PATCH /characters/:id/death-save
// Body: { "success": bool, "critical": bool }
//   success=true,  critical=false → +1 sucesso (3 = estabilizado)
//   success=true,  critical=true  → 20 natural: recupera 1 HP e acorda
//   success=false, critical=false → +1 falha
//   success=false, critical=true  → 1 natural: +2 falhas
func (h *CharacterHandler) DeathSave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var body struct {
		Success  bool `json:"success"`
		Critical bool `json:"critical"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}

	stabilized := false
	dead := false

	if body.Success {
		if body.Critical {
			// 20 natural: recupera 1 HP e acorda — reseta testes
			character.HitPoints = 1
			character.DeathSaveSuccesses = 0
			character.DeathSaveFailures = 0
			stabilized = true
		} else {
			character.DeathSaveSuccesses++
			if character.DeathSaveSuccesses >= 3 {
				// Estabilizado: reseta contadores mas continua em 0 HP
				character.DeathSaveSuccesses = 0
				character.DeathSaveFailures = 0
				stabilized = true
			}
		}
	} else {
		failures := 1
		if body.Critical {
			failures = 2 // 1 natural = 2 falhas
		}
		character.DeathSaveFailures += failures
		if character.DeathSaveFailures >= 3 {
			dead = true
		}
	}

	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"character":  character,
		"stabilized": stabilized,
		"dead":       dead,
	})
}

// ResetDeathSaves — reseta os testes de morte (usado em ressurreições)
// PATCH /characters/:id/reset-death-saves
func (h *CharacterHandler) ResetDeathSaves(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}
	character.DeathSaveSuccesses = 0
	character.DeathSaveFailures = 0
	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

// ── Export de PDF (5e) ──────────────────────────────────────────────────────────

// ExportPDF5e — GET /characters/:id/export/pdf
// Só disponível para personagens de edição 5e. Todo o cálculo de regras
// (BuildPDF5eExportPayload) E o preenchimento do AcroForm (pdfexport,
// pacote Go puro embutido no binário) acontecem aqui — não depende mais do
// ai-service Python, que nunca foi hospedado em produção e sempre fazia
// essa rota falhar com 503 fora do ambiente local. Ver internal/pdfexport
// para o porquê da migração.
func (h *CharacterHandler) ExportPDF5e(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	character, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}

	if character.UserID != c.GetUint("userID") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return
	}

	if character.Edition != "5e" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exportação de ficha em PDF só está disponível para personagens de 5ª edição"})
		return
	}

	allPericias, err := h.PericiaService.GetAll("5e")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar catálogo de perícias"})
		return
	}
	characterPericias, err := h.PericiaService.GetByCharacter(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar perícias do personagem"})
		return
	}

	playerName := ""
	if user, err := h.UserRepo.FindByID(character.UserID); err == nil {
		playerName = user.Name
	}

	payload := service.BuildPDF5eExportPayload(character, allPericias, characterPericias, h.ArmorService, playerName)

	pdfBytes, err := pdfexport.FillCharacterSheet5e(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao preencher a ficha: %v", err)})
		return
	}

	filename := fmt.Sprintf("ficha_%s.pdf", character.Name)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}