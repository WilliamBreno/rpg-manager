package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
	"rpg-manager/internal/ws"
)

type SessionHandler struct {
	svc         *service.SessionService
	campaignSvc *service.CampaignService
	wsManager   *ws.Manager
}

func NewSessionHandler(svc *service.SessionService, campaignSvc *service.CampaignService, wsManager *ws.Manager) *SessionHandler {
	return &SessionHandler{svc: svc, campaignSvc: campaignSvc, wsManager: wsManager}
}

// POST /campaigns/:id/sessions — abre uma nova sessão.
func (h *SessionHandler) Start(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	session, err := h.svc.Start(campaignID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

// GET /campaigns/:id/sessions — histórico de sessões da campanha.
func (h *SessionHandler) GetByCampaign(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	sessions, err := h.svc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) getOwned(c *gin.Context) (domain.Session, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.Session{}, false
	}
	session, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sessão não encontrada"})
		return domain.Session{}, false
	}
	if !requireCampaignOwner(c, h.campaignSvc, session.CampaignID) {
		return domain.Session{}, false
	}
	return session, true
}

type sessionSummaryRequest struct {
	Summary string `json:"summary"`
}

// PATCH /sessions/:id/end — encerra, opcionalmente com o resumo final.
func (h *SessionHandler) End(c *gin.Context) {
	session, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req sessionSummaryRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.End(&session, req.Summary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

// PATCH /sessions/:id/summary — atualiza o diário durante a sessão.
func (h *SessionHandler) UpdateSummary(c *gin.Context) {
	session, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req sessionSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateSummary(&session, req.Summary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

type sessionMusicRequest struct {
	MusicURL string `json:"music_url"`
	Playing  bool   `json:"playing"`
}

// PATCH /sessions/:id/music — toca/para a música de fundo da sessão pros
// jogadores conectados (Etapa 9). URL simples por enquanto (mesmo padrão já
// usado em Enemy.PhotoURL/SoundURL) — upload de arquivo de verdade depende
// da integração com Cloudinary, ainda não escrita.
func (h *SessionHandler) SetMusic(c *gin.Context) {
	session, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req sessionMusicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.SetMusic(&session, req.MusicURL, req.Playing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.wsManager.Broadcast(session.CampaignID, ws.Event{
		Type: "play_audio",
		Data: gin.H{"kind": "music", "url": req.MusicURL, "playing": req.Playing},
	})
	c.JSON(http.StatusOK, session)
}
