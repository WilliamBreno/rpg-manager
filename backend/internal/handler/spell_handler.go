package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type SpellHandler struct{ svc *service.SpellService }

func NewSpellHandler(svc *service.SpellService) *SpellHandler {
	return &SpellHandler{svc: svc}
}

// GET /spells?edition=5e
func (h *SpellHandler) GetAll(c *gin.Context) {
	spells, err := h.svc.GetAll(c.Query("edition"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spells)
}

// GET /characters/:id/spells
func (h *SpellHandler) GetByCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	spells, err := h.svc.GetByCharacter(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spells)
}

// POST /characters/:id/spells/:spell_id
func (h *SpellHandler) Add(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	spellID, err := strconv.Atoi(c.Param("spell_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spell ID inválido"})
		return
	}
	if err := h.svc.Add(uint(charID), uint(spellID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Magia adicionada"})
}

// DELETE /characters/:id/spells/:spell_id
func (h *SpellHandler) Remove(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	spellID, err := strconv.Atoi(c.Param("spell_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spell ID inválido"})
		return
	}
	if err := h.svc.Remove(uint(charID), uint(spellID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Magia removida"})
}
