package handler

import (
	"net/http"
	"strconv"

	"rpg-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type AntecedentHandler struct {
	service *service.AntecedentService
}

func NewAntecedentHandler(s *service.AntecedentService) *AntecedentHandler {
	return &AntecedentHandler{service: s}
}

// GET /api/antecedentes?edition=5e
func (h *AntecedentHandler) GetAll(c *gin.Context) {
	edition := c.DefaultQuery("edition", "5e")
	backgrounds, err := h.service.GetAll(edition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, backgrounds)
}

// GET /api/antecedentes/:id
func (h *AntecedentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	bg, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Antecedente não encontrado"})
		return
	}
	c.JSON(http.StatusOK, bg)
}