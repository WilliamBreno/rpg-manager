package handler

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
	"rpg-manager/internal/ws"
)

type WSHandler struct {
	manager       *ws.Manager
	authSvc       *service.AuthService
	campaignSvc   *service.CampaignService
	membershipSvc *service.CampaignMembershipService
}

func NewWSHandler(manager *ws.Manager, authSvc *service.AuthService, campaignSvc *service.CampaignService, membershipSvc *service.CampaignMembershipService) *WSHandler {
	return &WSHandler{manager: manager, authSvc: authSvc, campaignSvc: campaignSvc, membershipSvc: membershipSvc}
}

// GET /ws/campaign/:id?token=<jwt> — o token vem via query param, não header
// Authorization, porque o navegador não permite header customizado no
// handshake de WebSocket (decisão de arquitetura já aprovada na Etapa 0).
// Entra quem for o mestre dono da campanha OU tiver CampaignMembership com
// Status accepted — mesmo critério em espírito do 404-não-403 usado no resto
// da API (aqui vira só um 401/404 antes mesmo do upgrade, a conexão nunca
// chega a abrir pra quem não pode entrar).
func (h *WSHandler) JoinCampaign(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	token := c.Query("token")
	claims, err := h.authSvc.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return
	}
	allowed := campaign.MasterID == claims.UserID || h.membershipSvc.IsAcceptedMember(campaignID, claims.UserID)
	if !allowed {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return
	}

	// OriginPatterns precisa espelhar o allowlist de CORS do main.go — o
	// coder/websocket rejeita o handshake por padrão se Origin não bater com
	// o Host (frontend em :5173/Vercel, backend em :8080/Render, origens
	// diferentes).
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173", "rpg-manager-smoky.vercel.app"},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	h.manager.Join(c.Request.Context(), campaignID, conn)
}
