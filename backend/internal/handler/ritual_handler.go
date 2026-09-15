package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type RitualHandler struct{ svc *service.RitualService }

func NewRitualHandler(svc *service.RitualService) *RitualHandler {
	return &RitualHandler{svc: svc}
}

// GET /rituals?edition=4e
func (h *RitualHandler) GetAll(c *gin.Context) {
	rituals, err := h.svc.GetAll(c.Query("edition"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rituals)
}

// GET /characters/:id/rituals
func (h *RitualHandler) GetByCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	rituals, err := h.svc.GetByCharacter(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rituals)
}

// GET /characters/:id/rituals/access
func (h *RitualHandler) GetAccess(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	access, err := h.svc.GetAccess(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, access)
}

// POST /characters/:id/rituals/:ritual_id
func (h *RitualHandler) Add(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	ritualID, err := strconv.Atoi(c.Param("ritual_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ritual ID inválido"})
		return
	}
	if err := h.svc.Add(uint(charID), uint(ritualID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ritual adicionado"})
}

// DELETE /characters/:id/rituals/:ritual_id
func (h *RitualHandler) Remove(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	ritualID, err := strconv.Atoi(c.Param("ritual_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ritual ID inválido"})
		return
	}
	if err := h.svc.Remove(uint(charID), uint(ritualID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ritual removido"})
}
