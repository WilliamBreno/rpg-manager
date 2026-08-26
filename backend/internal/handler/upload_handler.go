package handler

import (
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "path/filepath"
    "strconv"

    "github.com/gin-gonic/gin"
    "rpg-manager/internal/repository"
)

type UploadHandler struct {
    CharacterRepo *repository.CharacterRepository
}

func NewUploadHandler(characterRepo *repository.CharacterRepository) *UploadHandler {
    return &UploadHandler{CharacterRepo: characterRepo}
}

// maxAvatarBytes limita o tamanho do arquivo original (antes da inflação de
// ~33% do base64) para não sobrecarregar a coluna de texto no Postgres nem o
// payload de cada GET /characters.
const maxAvatarBytes = 2 * 1024 * 1024 // 2MB

var avatarMimeByExt = map[string]string{
    ".jpg":  "image/jpeg",
    ".jpeg": "image/jpeg",
    ".png":  "image/png",
    ".webp": "image/webp",
}

// UploadAvatar recebe a imagem e a guarda como data URI (base64) diretamente
// na coluna avatar_url do personagem, em vez de salvar em disco.
//
// Por quê: o disco local do backend (./uploads) não é persistente entre
// reinícios/redeploys do serviço (hospedagem free tier), então avatares
// enviados anteriormente desapareciam silenciosamente enquanto o banco ainda
// apontava para eles, quebrando a imagem na UI. O Postgres, por outro lado,
// já é a fonte de persistência real do resto dos dados do personagem — guardar
// a imagem ali resolve o problema sem depender de um serviço de storage
// externo (S3/Cloudinary/etc.) que exigiria credenciais novas.
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
    characterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    character, err := h.CharacterRepo.FindByID(uint(characterID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
        return
    }

    fileHeader, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não encontrado"})
        return
    }

    ext := filepath.Ext(fileHeader.Filename)
    mimeType, ok := avatarMimeByExt[ext]
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido. Use jpg, jpeg, png ou webp"})
        return
    }

    if fileHeader.Size > maxAvatarBytes {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Imagem muito grande. Máximo de 2MB"})
        return
    }

    file, err := fileHeader.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir arquivo"})
        return
    }
    defer file.Close()

    bytes, err := io.ReadAll(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler arquivo"})
        return
    }

    dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(bytes))

    character.AvatarURL = dataURI
    if err := h.CharacterRepo.Update(&character); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar personagem"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "avatar_url": character.AvatarURL,
        "message":    "Avatar atualizado com sucesso",
    })
}
