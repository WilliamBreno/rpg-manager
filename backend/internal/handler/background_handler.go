package handler

import (
	"net/http"
	"strconv"

	"rpg-manager/internal/service"

	"github.com/gin-gonic/gin"
)

// BackgroundHandler gerencia a biografia/notas do personagem
// (rota: GET/POST /characters/:id/background)
type BackgroundHandler struct {
	service *service.BackgroundService
}

func NewBackgroundHandler(s *service.BackgroundService) *BackgroundHandler {
	return &BackgroundHandler{service: s}
}

// GET /api/characters/:id/background
func (h *BackgroundHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	bg, err := h.service.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bg)
}

// POST /api/characters/:id/background
func (h *BackgroundHandler) Save(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Save(uint(id), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Background salvo"})
}