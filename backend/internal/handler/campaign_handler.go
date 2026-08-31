package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
)

type CampaignHandler struct{ svc *service.CampaignService }

func NewCampaignHandler(svc *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{svc: svc}
}

type createCampaignRequest struct {
	Name      string `json:"name"`
	MainStory string `json:"main_story"`
}

// POST /campaigns — só um usuário com Role == RoleMaster pode criar campanha
// (mesmo padrão de gate já usado em InventoryHandler.SetCurrency).
func (h *CampaignHandler) Create(c *gin.Context) {
	role, _ := c.Get("userRole")
	if role != domain.RoleMaster {
		c.JSON(http.StatusForbidden, gin.H{"error": "Só o Mestre pode criar campanhas"})
		return
	}
	var req createCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campaign := domain.Campaign{
		Name:      req.Name,
		MainStory: req.MainStory,
		Edition:   "5e",
		MasterID:  c.GetUint("userID"),
	}
	if err := h.svc.Create(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, campaign)
}

// GET /campaigns — campanhas do mestre autenticado (não lista campanhas de
// outros mestres, e não lista campanhas em que o usuário é só jogador — isso
// é uma tela futura, ver Etapa 6).
func (h *CampaignHandler) GetAll(c *gin.Context) {
	campaigns, err := h.svc.GetByMaster(c.GetUint("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campaigns)
}

// getOwned resolve a campanha do :id e confirma que pertence ao usuário
// autenticado — 404 (não 403) se não for dono, mesmo padrão já usado em
// ExportPDF5e pra não revelar a existência de campanhas de outros mestres.
func (h *CampaignHandler) getOwned(c *gin.Context) (domain.Campaign, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return domain.Campaign{}, false
	}
	campaign, err := h.svc.GetByID(uint(id))
	if err != nil || campaign.MasterID != c.GetUint("userID") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return domain.Campaign{}, false
	}
	return campaign, true
}

// GET /campaigns/:id
func (h *CampaignHandler) GetByID(c *gin.Context) {
	campaign, ok := h.getOwned(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, campaign)
}

type updateCampaignRequest struct {
	Name      string `json:"name"`
	MainStory string `json:"main_story"`
}

// PUT /campaigns/:id
func (h *CampaignHandler) Update(c *gin.Context) {
	campaign, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req updateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.svc.Update(campaign.ID, req.Name, req.MainStory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DELETE /campaigns/:id
func (h *CampaignHandler) Delete(c *gin.Context) {
	campaign, ok := h.getOwned(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(campaign.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Campanha removida"})
}
