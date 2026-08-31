package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
)

type NPCHandler struct {
	svc         *service.NPCService
	campaignSvc *service.CampaignService
}

func NewNPCHandler(svc *service.NPCService, campaignSvc *service.CampaignService) *NPCHandler {
	return &NPCHandler{svc: svc, campaignSvc: campaignSvc}
}

type npcRequest struct {
	Name        string `json:"name"`
	HP          int    `json:"hp"`
	History     string `json:"history"`
	Bonds       string `json:"bonds"`
	Alignment   string `json:"alignment"`
	Personality string `json:"personality"`
	Notes       string `json:"notes"`
}

// POST /campaigns/:id/npcs
func (h *NPCHandler) Create(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req npcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	npc := domain.NPC{
		CampaignID: campaignID, Name: req.Name, HP: req.HP, History: req.History,
		Bonds: req.Bonds, Alignment: req.Alignment, Personality: req.Personality, Notes: req.Notes,
	}
	if err := h.svc.Create(&npc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, npc)
}

// GET /campaigns/:id/npcs
func (h *NPCHandler) GetByCampaign(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	npcs, err := h.svc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, npcs)
}

func (h *NPCHandler) getOwned(c *gin.Context) (domain.NPC, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.NPC{}, false
	}
	npc, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NPC não encontrado"})
		return domain.NPC{}, false
	}
	if !requireCampaignOwner(c, h.campaignSvc, npc.CampaignID) {
		return domain.NPC{}, false
	}
	return npc, true
}

// PUT /npcs/:id
func (h *NPCHandler) Update(c *gin.Context) {
	npc, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req npcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	npc.Name, npc.HP, npc.History = req.Name, req.HP, req.History
	npc.Bonds, npc.Alignment, npc.Personality, npc.Notes = req.Bonds, req.Alignment, req.Personality, req.Notes
	if err := h.svc.Update(&npc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, npc)
}

// DELETE /npcs/:id
func (h *NPCHandler) Delete(c *gin.Context) {
	npc, ok := h.getOwned(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(npc.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "NPC removido"})
}
