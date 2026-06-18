package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
)

type CharacterHandler struct {
	Service      *service.CharacterService
	ArmorService *service.ArmorService
}

func NewCharacterHandler(service *service.CharacterService, armorService *service.ArmorService) *CharacterHandler {
	return &CharacterHandler{Service: service, ArmorService: armorService}
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

// ── AddXP ─────────────────────────────────────────────────────────────────────
// PATCH /characters/:id/add-xp
// Body: { "xp": 250 }
// Response: { character, leveled_up, needs_asi, new_level }
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

// ── ApplyASI ──────────────────────────────────────────────────────────────────
// PATCH /characters/:id/apply-asi
// Body: { "strength": 0, "dexterity": 0, "constitution": 1, "intelligence": 0, "wisdom": 1, "charisma": 0 }
// Total deve ser 1 ou 2. Cada campo indica quantos pontos adicionar naquele atributo.
// Response: { character, leveled_up, needs_asi, new_level }
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

// ── LevelUp manual (mantido para compatibilidade) ─────────────────────────────
// PATCH /characters/:id/level-up
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

	character.HitPoints += body.Amount
	if character.HitPoints > character.MaxHP {
		character.HitPoints = character.MaxHP
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

	// Temp HP não acumula — fica com o maior valor
	if body.Amount > character.TempHP {
		character.TempHP = body.Amount
	}

	if err := h.Service.Update(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}