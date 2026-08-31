package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

// requireCampaignOwner é reaproveitado por todo handler de recurso que
// pertence a uma Campaign (NPC, Enemy, Scene, Session, Reward, mensagens de
// chat enviadas pelo mestre) — confirma que campaignID pertence ao usuário
// autenticado antes de deixar prosseguir. 404 (não 403) pra não revelar a
// existência de campanha de outro mestre, mesmo padrão já usado em
// CampaignHandler.getOwned/ExportPDF5e.
func requireCampaignOwner(c *gin.Context, campaignSvc *service.CampaignService, campaignID uint) bool {
	campaign, err := campaignSvc.GetByID(campaignID)
	if err != nil || campaign.MasterID != c.GetUint("userID") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campanha não encontrada"})
		return false
	}
	return true
}

// parseUintParam é um pequeno atalho pra strconv.ParseUint + resposta 400,
// repetido em praticamente todo handler desta feature.
func parseUintParam(c *gin.Context, name string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return 0, false
	}
	return uint(v), true
}
