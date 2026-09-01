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

// maxImageBytes/maxAudioBytes — mesmo raciocínio de maxAvatarBytes acima:
// tamanho do arquivo original antes da inflação de ~33% do base64. Áudio
// aceita mais porque uma música de fundo ou fala de vilão curta ainda assim
// pesa mais que uma foto.
const (
	maxImageBytes = 3 * 1024 * 1024 // 3MB
	maxAudioBytes = 8 * 1024 * 1024 // 8MB
)

var imageMimeByExt = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
	".gif":  "image/gif",
}

var audioMimeByExt = map[string]string{
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".ogg":  "audio/ogg",
	".m4a":  "audio/mp4",
	".webm": "audio/webm",
}

// UploadFile é um upload genérico, desacoplado de qualquer entidade — recebe
// um arquivo (campo "file") e devolve o data URI base64 correspondente, sem
// persistir nada. Pensado pro Sistema do Mestre (foto/som de Enemy, imagem
// de Scene, áudio de EnemyLine, música de Session): esses recursos não têm
// (e não deveriam precisar de) uma rota de upload própria cada um — o
// front-end sobe o arquivo aqui e manda o data URI resultante junto do
// payload normal de criação/edição do recurso, exatamente como já
// funcionava seria feito colando uma URL antes. Mesma justificativa de
// armazenar como data URI em vez de disco/S3 já documentada em UploadAvatar
// acima (sem custo de storage externo, sem credencial nova).
func (h *UploadHandler) UploadFile(c *gin.Context) {
	kind := c.Query("kind")
	var mimeByExt map[string]string
	var maxBytes int64
	var maxLabel string
	switch kind {
	case "image":
		mimeByExt, maxBytes, maxLabel = imageMimeByExt, maxImageBytes, "3MB"
	case "audio":
		mimeByExt, maxBytes, maxLabel = audioMimeByExt, maxAudioBytes, "8MB"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro kind deve ser 'image' ou 'audio'"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	mimeType, ok := mimeByExt[ext]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Formato inválido para %s", kind)})
		return
	}

	if fileHeader.Size > maxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Arquivo muito grande. Máximo de %s", maxLabel)})
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
	c.JSON(http.StatusOK, gin.H{"data_url": dataURI})
}
