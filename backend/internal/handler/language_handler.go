package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type LanguageHandler struct{ svc *service.LanguageService }

func NewLanguageHandler(svc *service.LanguageService) *LanguageHandler {
	return &LanguageHandler{svc: svc}
}

// GET /languages?edition=5e
func (h *LanguageHandler) GetAll(c *gin.Context) {
	languages, err := h.svc.GetAll(c.Query("edition"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}

// GET /characters/:id/languages
func (h *LanguageHandler) GetByCharacter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	languages, err := h.svc.GetByCharacter(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}

// POST /characters/:id/languages/:language_id
func (h *LanguageHandler) Add(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	languageID, err := strconv.Atoi(c.Param("language_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language ID inválido"})
		return
	}
	if err := h.svc.Add(uint(charID), uint(languageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Idioma adicionado"})
}

// DELETE /characters/:id/languages/:language_id
func (h *LanguageHandler) Remove(c *gin.Context) {
	charID, _ := strconv.Atoi(c.Param("id"))
	languageID, err := strconv.Atoi(c.Param("language_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language ID inválido"})
		return
	}
	if err := h.svc.Remove(uint(charID), uint(languageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Idioma removido"})
}
