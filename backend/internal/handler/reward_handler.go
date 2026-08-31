package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
)

type RewardHandler struct {
	svc           *service.RewardService
	itemSvc       *service.MagicItemService
	campaignSvc   *service.CampaignService
	membershipSvc *service.CampaignMembershipService
}

func NewRewardHandler(svc *service.RewardService, itemSvc *service.MagicItemService, campaignSvc *service.CampaignService, membershipSvc *service.CampaignMembershipService) *RewardHandler {
	return &RewardHandler{svc: svc, itemSvc: itemSvc, campaignSvc: campaignSvc, membershipSvc: membershipSvc}
}

type magicItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Effect      string `json:"effect"`
}

// POST /campaigns/:id/magic-items
func (h *RewardHandler) CreateMagicItem(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req magicItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := domain.MagicItem{CampaignID: campaignID, Name: req.Name, Description: req.Description, Effect: req.Effect}
	if err := h.itemSvc.Create(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// GET /campaigns/:id/magic-items
func (h *RewardHandler) GetMagicItems(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	items, err := h.itemSvc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// DELETE /magic-items/:id
func (h *RewardHandler) DeleteMagicItem(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	item, err := h.itemSvc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, item.CampaignID) {
		return
	}
	if err := h.itemSvc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item removido"})
}

// GET /campaigns/:id/rewards — histórico de entregas da campanha.
func (h *RewardHandler) GetHistory(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	rewards, err := h.svc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rewards)
}

// recipientCharacterIDs resolve pra quem a entrega vai: um personagem
// específico, ou (All: true) todo personagem vinculado a um membro aceito da
// campanha — "dar a todos" gera uma linha de Reward por personagem (ver
// domain.Reward), não uma entrega multi-destinatário.
func (h *RewardHandler) recipientCharacterIDs(campaignID uint, characterID *uint, all bool) ([]uint, error) {
	if !all {
		if characterID == nil {
			return nil, errBadRecipient
		}
		return []uint{*characterID}, nil
	}
	members, err := h.membershipSvc.GetByCampaign(campaignID)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, m := range members {
		if m.Status == domain.MembershipAccepted && m.CharacterID != nil {
			ids = append(ids, *m.CharacterID)
		}
	}
	return ids, nil
}

var errBadRecipient = errors.New("informe character_id ou all=true")

type grantCurrencyRequest struct {
	CharacterID    *uint  `json:"character_id"`
	All            bool   `json:"all"`
	CopperPieces   int    `json:"copper_pieces"`
	SilverPieces   int    `json:"silver_pieces"`
	ElectrumPieces int    `json:"electrum_pieces"`
	GoldPieces     int    `json:"gold_pieces"`
	PlatinumPieces int    `json:"platinum_pieces"`
	Note           string `json:"note"`
}

// POST /campaigns/:id/rewards/currency
func (h *RewardHandler) GrantCurrency(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req grantCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipients, err := h.recipientCharacterIDs(campaignID, req.CharacterID, req.All)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grantedBy := c.GetUint("userID")
	var rewards []domain.Reward
	for _, characterID := range recipients {
		reward, err := h.svc.GrantCurrency(campaignID, characterID, grantedBy,
			req.CopperPieces, req.SilverPieces, req.ElectrumPieces, req.GoldPieces, req.PlatinumPieces, req.Note)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rewards = append(rewards, reward)
	}
	c.JSON(http.StatusCreated, rewards)
}

type grantItemRequest struct {
	CharacterID *uint  `json:"character_id"`
	All         bool   `json:"all"`
	MagicItemID uint   `json:"magic_item_id"`
	Note        string `json:"note"`
}

// POST /campaigns/:id/rewards/item
func (h *RewardHandler) GrantItem(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req grantItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipients, err := h.recipientCharacterIDs(campaignID, req.CharacterID, req.All)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grantedBy := c.GetUint("userID")
	var rewards []domain.Reward
	for _, characterID := range recipients {
		reward, err := h.svc.GrantItem(campaignID, characterID, grantedBy, req.MagicItemID, req.Note)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rewards = append(rewards, reward)
	}
	c.JSON(http.StatusCreated, rewards)
}
