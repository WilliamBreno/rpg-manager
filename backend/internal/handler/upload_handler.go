package handler

import (
    "fmt"
    "net/http"
    "path/filepath"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/repository"
)

type UploadHandler struct {
    CharacterRepo *repository.CharacterRepository
}

func NewUploadHandler(characterRepo *repository.CharacterRepository) *UploadHandler {
    return &UploadHandler{CharacterRepo: characterRepo}
}

func (h *UploadHandler) UploadAvatar(c *gin.Context) {
    characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    // Busca o personagem
    character, err := h.CharacterRepo.FindByID(uint(characterID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }

    // Recebe o arquivo
    file, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não encontrado"})
        return
    }

    // Valida o tipo do arquivo
    ext := filepath.Ext(file.Filename)
    allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
    if !allowed[ext] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido. Use jpg, jpeg, png ou webp"})
        return
    }

    // Gera um nome único para o arquivo
    filename := fmt.Sprintf("character_%d_%d%s", characterID, time.Now().Unix(), ext)
    savePath := filepath.Join("uploads", filename)

    // Salva o arquivo
    if err := c.SaveUploadedFile(file, savePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
        return
    }

    // Atualiza o personagem com a URL do avatar
    character.AvatarURL = fmt.Sprintf("/uploads/%s", filename)
    if err := h.CharacterRepo.Update(&character); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar personagem"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "avatar_url": character.AvatarURL,
        "message":    "Avatar atualizado com sucesso",
    })
}