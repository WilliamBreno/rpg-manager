package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/service"
)

type RaceHandler struct {
    Service *service.RaceService
}

func NewRaceHandler(service *service.RaceService) *RaceHandler {
    return &RaceHandler{Service: service}
}

func (h *RaceHandler) GetAll(c *gin.Context) {
    edition := c.Query("edition")

    var races []domain.Race
    var err error

    if edition != "" {
        races, err = h.Service.GetByEdition(edition)
    } else {
        races, err = h.Service.GetAll()
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, races)
}

func (h *RaceHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    race, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Raça não encontrada"})
        return
    }
    c.JSON(http.StatusOK, race)
}

func (h *RaceHandler) Create(c *gin.Context) {
    var race domain.Race
    if err := c.ShouldBindJSON(&race); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Service.Create(&race); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, race)
}

func (h *RaceHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var race domain.Race
    if err := c.ShouldBindJSON(&race); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    race.ID = uint(id)
    if err := h.Service.Update(&race); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, race)
}

func (h *RaceHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    if err := h.Service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Raça deletada com sucesso"})
}