package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

type RegisterRequest struct {
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     domain.UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome, email e senha são obrigatórios"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Senha deve ter no mínimo 6 caracteres"})
		return
	}

	// Se não informar role, padrão é player
	if req.Role == "" {
		req.Role = domain.RolePlayer
	}

	user, err := h.Service.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário criado com sucesso",
		"user": gin.H{
			"id":                  user.ID,
			"name":                user.Name,
			"email":               user.Email,
			"role":                user.Role,
			"welcome_seen":        user.WelcomeSeen,
			"master_welcome_seen": user.MasterWelcomeSeen,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":                  user.ID,
			"name":                user.Name,
			"email":               user.Email,
			"role":                user.Role,
			"welcome_seen":        user.WelcomeSeen,
			"master_welcome_seen": user.MasterWelcomeSeen,
		},
	})
}

// MarkWelcomeSeen — PATCH /users/me/welcome-seen
// Marca que o usuário autenticado já viu o modal de boas-vindas, pra não
// mostrar de novo em outro login (inclusive de outro navegador/dispositivo).
func (h *AuthHandler) MarkWelcomeSeen(c *gin.Context) {
	userID := c.GetUint("userID")
	if err := h.Service.MarkWelcomeSeen(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"welcome_seen": true})
}

// MarkMasterWelcomeSeen — PATCH /users/me/master-welcome-seen
// Mesma ideia de MarkWelcomeSeen, mas pra tela de boas-vindas específica de
// Mestre (Sistema do Mestre) — as duas são independentes.
func (h *AuthHandler) MarkMasterWelcomeSeen(c *gin.Context) {
	userID := c.GetUint("userID")
	if err := h.Service.MarkMasterWelcomeSeen(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"master_welcome_seen": true})
}
