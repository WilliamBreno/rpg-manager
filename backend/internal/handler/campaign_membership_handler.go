package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type CampaignMembershipHandler struct {
	svc         *service.CampaignMembershipService
	campaignSvc *service.CampaignService
}

func NewCampaignMembershipHandler(svc *service.CampaignMembershipService, campaignSvc *service.CampaignService) *CampaignMembershipHandler {
	return &CampaignMembershipHandler{svc: svc, campaignSvc: campaignSvc}
}

type inviteRequest struct {
	Email string `json:"email"`
}

// POST /campaigns/:id/invites — o mestre convida um jogador pelo e-mail.
func (h *CampaignMembershipHandler) Invite(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req inviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	membership, err := h.svc.Invite(campaignID, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, membership)
}

// GET /campaigns/:id/members — elenco de membros aceitos (e convites
// pendentes) da campanha. É a área "ver quem mais está na campanha" do lado
// do jogador citada na Etapa 6 (ver CLAUDE.md pra a decisão tomada sobre a
// ambiguidade "adicionar jogadores" do documento original) — qualquer membro
// aceito também pode ver o elenco, não só o mestre.
func (h *CampaignMembershipHandler) GetByCampaign(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return
	}
	if campaign.MasterID != userID && !h.svc.IsAcceptedMember(campaignID, userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return
	}
	members, err := h.svc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

// GET /users/me/campaign-invites — convites pendentes do jogador logado.
func (h *CampaignMembershipHandler) GetMyPending(c *gin.Context) {
	userID := c.GetUint("userID")
	pending, err := h.svc.GetPendingForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pending)
}

// GET /users/me/campaigns — campanhas que o jogador já participa (convite
// aceito), pra navegar até a Sala ao vivo de cada uma.
func (h *CampaignMembershipHandler) GetMyCampaigns(c *gin.Context) {
	userID := c.GetUint("userID")
	accepted, err := h.svc.GetAcceptedForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accepted)
}

type respondRequest struct {
	Accept      bool  `json:"accept"`
	CharacterID *uint `json:"character_id"`
}

// PATCH /campaign-memberships/:id/respond — o jogador aceita ou recusa.
func (h *CampaignMembershipHandler) Respond(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	membership, err := h.svc.GetByID(id)
	if err != nil || membership.UserID != c.GetUint("userID") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Convite não encontrado"})
		return
	}
	var req respondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Respond(&membership, req.Accept, req.CharacterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, membership)
}
