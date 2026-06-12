package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/service"
)

type BackgroundHandler struct {
    Service *service.BackgroundService
}

func NewBackgroundHandler(service *service.BackgroundService) *BackgroundHandler {
    return &BackgroundHandler{Service: service}
}

func (h *BackgroundHandler) Get(c *gin.Context) {
    characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    background, err := h.Service.GetByCharacterID(uint(characterID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Background não encontrado"})
        return
    }
    c.JSON(http.StatusOK, background)
}

func (h *BackgroundHandler) Save(c *gin.Context) {
    characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var background domain.Background
    if err := c.ShouldBindJSON(&background); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    background.CharacterID = uint(characterID)

    if err := h.Service.Save(&background); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, background)
}