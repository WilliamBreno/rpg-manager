package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/service"
)

type CharacterHandler struct {
    Service *service.CharacterService
    ArmorService *service.ArmorService
}

func NewCharacterHandler(service *service.CharacterService, armorService *service.ArmorService) *CharacterHandler {
    return &CharacterHandler{Service: service, ArmorService: armorService}
}

func (h *CharacterHandler) GetAll(c *gin.Context) {
    userID := c.GetUint("userID")
    characters, err := h.Service.GetAll(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, characters)
}

func (h *CharacterHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    character, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }
    c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) Create(c *gin.Context) {
    var character domain.Character
    if err := c.ShouldBindJSON(&character); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    character.UserID = c.GetUint("userID")

    if err := h.Service.Create(&character); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, character)
}

func (h *CharacterHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var character domain.Character
    if err := c.ShouldBindJSON(&character); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    character.ID = uint(id)
    character.UserID = c.GetUint("userID")
    
    if err := h.Service.Update(&character); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    if err := h.Service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Personagem deletado com sucesso"})
}

func (h *CharacterHandler) LevelUp(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    character, err := h.Service.LevelUp(uint(id))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) AddSkill(c *gin.Context) {
    characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    skillID, err := strconv.ParseUint(c.Param("skill_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "skill_id inválido"})
        return
    }

    if err := h.Service.AddSkill(uint(characterID), uint(skillID)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Habilidade adicionada com sucesso"})
}

func (h *CharacterHandler) RemoveSkill(c *gin.Context) {
    characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    skillID, err := strconv.ParseUint(c.Param("skill_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "skill_id inválido"})
        return
    }

    if err := h.Service.RemoveSkill(uint(characterID), uint(skillID)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Habilidade removida com sucesso"})
}
func (h *CharacterHandler) GetAC(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    character, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }

    ac := h.ArmorService.CalculateAC(character)
    c.JSON(http.StatusOK, gin.H{"ac": ac})
}
func (h *CharacterHandler) TakeDamage(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var body struct {
        Damage int `json:"damage"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    character, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }

    // Desconta do TempHP primeiro
    if character.TempHP > 0 {
        if body.Damage <= character.TempHP {
            character.TempHP -= body.Damage
            body.Damage = 0
        } else {
            body.Damage -= character.TempHP
            character.TempHP = 0
        }
    }

    // Desconta do HP real
    character.HitPoints -= body.Damage
    if character.HitPoints < 0 {
        character.HitPoints = 0
    }

    if err := h.Service.Update(&character); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) Heal(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var body struct {
        Amount int `json:"amount"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    character, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }

    character.HitPoints += body.Amount
    if character.HitPoints > character.MaxHP {
        character.HitPoints = character.MaxHP
    }

    if err := h.Service.Update(&character); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, character)
}

func (h *CharacterHandler) AddTempHP(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var body struct {
        Amount int `json:"amount"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    character, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }

    // Temp HP não acumula — fica com o maior valor
    if body.Amount > character.TempHP {
        character.TempHP = body.Amount
    }

    if err := h.Service.Update(&character); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, character)
}