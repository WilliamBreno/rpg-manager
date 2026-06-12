package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/service"
)

type ArmorHandler struct {
    Service *service.ArmorService
}

func NewArmorHandler(service *service.ArmorService) *ArmorHandler {
    return &ArmorHandler{Service: service}
}

func (h *ArmorHandler) GetByEdition(c *gin.Context) {
    edition := c.Query("edition")
    if edition == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Edição é obrigatória"})
        return
    }

    armors, err := h.Service.GetByEdition(edition)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, armors)
}