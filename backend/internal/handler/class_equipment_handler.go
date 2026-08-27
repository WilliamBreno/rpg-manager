package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type ClassEquipmentHandler struct{ Service *service.ClassEquipmentService }

func NewClassEquipmentHandler(s *service.ClassEquipmentService) *ClassEquipmentHandler {
	return &ClassEquipmentHandler{Service: s}
}

// GetByClass — GET /classes/:id/equipment-options
func (h *ClassEquipmentHandler) GetByClass(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de classe inválido"})
		return
	}
	options, err := h.Service.GetByClass(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, options)
}
