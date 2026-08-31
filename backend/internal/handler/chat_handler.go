package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
	"rpg-manager/internal/ws"
)

type ChatHandler struct {
	svc           *service.ChatService
	campaignSvc   *service.CampaignService
	membershipSvc *service.CampaignMembershipService
	wsManager     *ws.Manager
}

func NewChatHandler(svc *service.ChatService, campaignSvc *service.CampaignService, membershipSvc *service.CampaignMembershipService, wsManager *ws.Manager) *ChatHandler {
	return &ChatHandler{svc: svc, campaignSvc: campaignSvc, membershipSvc: membershipSvc, wsManager: wsManager}
}

// requireCampaignAccess é mais permissivo que requireCampaignOwner: aceita o
// mestre OU qualquer membro com convite aceito — o chat (e o resto da sala
// ao vivo) precisa ser visível pros dois lados, diferente dos recursos de
// gerenciamento (NPC/Enemy/Scene/Session) que são só do mestre.
func (h *ChatHandler) requireCampaignAccess(c *gin.Context, campaignID uint) bool {
	userID := c.GetUint("userID")
	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return false
	}
	if campaign.MasterID == userID || h.membershipSvc.IsAcceptedMember(campaignID, userID) {
		return true
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
	return false
}

// GET /campaigns/:id/chat — histórico (últimas mensagens).
func (h *ChatHandler) GetHistory(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !h.requireCampaignAccess(c, campaignID) {
		return
	}
	messages, err := h.svc.History(campaignID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

type chatMessageRequest struct {
	Text      string `json:"text"`
	SessionID *uint  `json:"session_id"`
}

// POST /campaigns/:id/chat — envia e transmite via WebSocket pra quem está
// conectado na sala (ver internal/ws — REST persiste, WS só notifica).
func (h *ChatHandler) Send(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !h.requireCampaignAccess(c, campaignID) {
		return
	}
	var req chatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg := domain.ChatMessage{
		CampaignID: campaignID, SessionID: req.SessionID,
		SenderUserID: c.GetUint("userID"), Text: req.Text,
	}
	if err := h.svc.Send(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.wsManager.Broadcast(campaignID, ws.Event{Type: "chat_message", Data: msg})
	c.JSON(http.StatusCreated, msg)
}
