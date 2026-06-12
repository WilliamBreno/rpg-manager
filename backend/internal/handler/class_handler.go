package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/service"
)

type ClassHandler struct {
    Service *service.ClassService
}

func NewClassHandler(service *service.ClassService) *ClassHandler {
    return &ClassHandler{Service: service}
}

func (h *ClassHandler) GetAll(c *gin.Context) {
    edition := c.Query("edition")

    var classes []domain.Class
    var err error

    if edition != "" {
        classes, err = h.Service.GetByEdition(edition)
    } else {
        classes, err = h.Service.GetAll()
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, classes)
}

func (h *ClassHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    class, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Classe não encontrada"})
        return
    }
    c.JSON(http.StatusOK, class)
}

func (h *ClassHandler) Create(c *gin.Context) {
    var class domain.Class
    if err := c.ShouldBindJSON(&class); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Service.Create(&class); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, class)
}

func (h *ClassHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var class domain.Class
    if err := c.ShouldBindJSON(&class); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    class.ID = uint(id)
    if err := h.Service.Update(&class); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, class)
}

func (h *ClassHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    if err := h.Service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Classe deletada com sucesso"})
}