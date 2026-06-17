package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type TalentoHandler struct{ svc *service.TalentoService }

func NewTalentoHandler(svc *service.TalentoService) *TalentoHandler {
	return &TalentoHandler{svc: svc}
}

// GET /talentos?edition=4e
func (h *TalentoHandler) GetAll(c *gin.Context) {
	talentos, err := h.svc.GetAll(c.Query("edition"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, talentos)
}

// GET /characters/:id/talentos
func (h *TalentoHandler) GetByCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	talentos, err := h.svc.GetByCharacter(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, talentos)
}

// POST /characters/:id/talentos/:talento_id
func (h *TalentoHandler) Add(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	talentoID, err := strconv.Atoi(c.Param("talento_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talento ID inválido"})
		return
	}
	if err := h.svc.Add(uint(charID), uint(talentoID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Talento adicionado"})
}

// DELETE /characters/:id/talentos/:talento_id
func (h *TalentoHandler) Remove(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	talentoID, err := strconv.Atoi(c.Param("talento_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talento ID inválido"})
		return
	}
	if err := h.svc.Remove(uint(charID), uint(talentoID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Talento removido"})
}