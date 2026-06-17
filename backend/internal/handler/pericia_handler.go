package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type PericiaHandler struct{ svc *service.PericiaService }

func NewPericiaHandler(svc *service.PericiaService) *PericiaHandler {
	return &PericiaHandler{svc: svc}
}

// GET /pericias?edition=4e
func (h *PericiaHandler) GetAll(c *gin.Context) {
	pericias, err := h.svc.GetAll(c.Query("edition"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pericias)
}

// GET /characters/:id/pericias
func (h *PericiaHandler) GetByCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	pericias, err := h.svc.GetByCharacter(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pericias)
}

// POST /characters/:id/pericias
// Body: {"pericias": ["Atletismo", "Percepção"]}
func (h *PericiaHandler) Save(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var body struct {
		Pericias []string `json:"pericias"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Save(uint(id), body.Pericias); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Perícias salvas"})
}