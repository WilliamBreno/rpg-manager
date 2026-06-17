package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/service"
)

type SkillHandler struct {
    Service *service.SkillService
}

func NewSkillHandler(service *service.SkillService) *SkillHandler {
    return &SkillHandler{Service: service}
}

func (h *SkillHandler) GetAll(c *gin.Context) {
    skills, err := h.Service.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, skills)
}

func (h *SkillHandler) GetByClassAndRace(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Query("class_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id inválido"})
		return
	}

	// ── race_id é opcional ──────────────────────────────────────────
	var raceIDPtr *uint
	if raceStr := c.Query("race_id"); raceStr != "" {
		raceID, err := strconv.ParseUint(raceStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "race_id inválido"})
			return
		}
		id := uint(raceID)
		raceIDPtr = &id
	}
	// ───────────────────────────────────────────────────────────────

	skills, err := h.Service.GetByClassAndRace(uint(classID), raceIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, skills)
}

func (h *SkillHandler) Create(c *gin.Context) {
    var skill domain.Skill
    if err := c.ShouldBindJSON(&skill); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Service.Create(&skill); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, skill)
}

func (h *SkillHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var skill domain.Skill
    if err := c.ShouldBindJSON(&skill); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    skill.ID = uint(id)
    if err := h.Service.Update(&skill); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, skill)
}

func (h *SkillHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    if err := h.Service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Habilidade deletada com sucesso"})
}