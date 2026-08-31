package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
	"rpg-manager/internal/ws"
)

type SceneHandler struct {
	svc         *service.SceneService
	campaignSvc *service.CampaignService
	sessionSvc  *service.SessionService
	wsManager   *ws.Manager
}

func NewSceneHandler(svc *service.SceneService, campaignSvc *service.CampaignService, sessionSvc *service.SessionService, wsManager *ws.Manager) *SceneHandler {
	return &SceneHandler{svc: svc, campaignSvc: campaignSvc, sessionSvc: sessionSvc, wsManager: wsManager}
}

type sceneRequest struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// POST /campaigns/:id/scenes
func (h *SceneHandler) Create(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req sceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scene := domain.Scene{CampaignID: campaignID, Name: req.Name, ImageURL: req.ImageURL}
	if err := h.svc.Create(&scene); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, scene)
}

// GET /campaigns/:id/scenes — biblioteca de cenários da campanha.
func (h *SceneHandler) GetByCampaign(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	scenes, err := h.svc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scenes)
}

func (h *SceneHandler) getOwned(c *gin.Context) (domain.Scene, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.Scene{}, false
	}
	scene, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cenário não encontrado"})
		return domain.Scene{}, false
	}
	if !requireCampaignOwner(c, h.campaignSvc, scene.CampaignID) {
		return domain.Scene{}, false
	}
	return scene, true
}

// GET /scenes/:id
func (h *SceneHandler) GetByID(c *gin.Context) {
	scene, ok := h.getOwned(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, scene)
}

// PUT /scenes/:id
func (h *SceneHandler) Update(c *gin.Context) {
	scene, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req sceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scene.Name, scene.ImageURL = req.Name, req.ImageURL
	if err := h.svc.Update(&scene); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scene)
}

// DELETE /scenes/:id
func (h *SceneHandler) Delete(c *gin.Context) {
	scene, ok := h.getOwned(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(scene.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cenário removido"})
}

type tokenRequest struct {
	Label    string  `json:"label"`
	ImageURL string  `json:"image_url"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	EnemyID  *uint   `json:"enemy_id"`
	NPCID    *uint   `json:"npc_id"`
}

// POST /scenes/:id/tokens
func (h *SceneHandler) CreateToken(c *gin.Context) {
	sceneID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	scene, err := h.svc.GetByID(sceneID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cenário não encontrado"})
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, scene.CampaignID) {
		return
	}
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := domain.Token{
		SceneID: sceneID, Label: req.Label, ImageURL: req.ImageURL,
		X: req.X, Y: req.Y, EnemyID: req.EnemyID, NPCID: req.NPCID,
	}
	if err := h.svc.CreateToken(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, token)
}

func (h *SceneHandler) getOwnedToken(c *gin.Context) (domain.Token, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.Token{}, false
	}
	token, err := h.svc.GetTokenByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token não encontrado"})
		return domain.Token{}, false
	}
	scene, err := h.svc.GetByID(token.SceneID)
	if err != nil || !requireCampaignOwner(c, h.campaignSvc, scene.CampaignID) {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Token não encontrado"})
		}
		return domain.Token{}, false
	}
	return token, true
}

type tokenMoveRequest struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// PATCH /tokens/:id/move — arrastar-e-soltar do Konva chama isso.
func (h *SceneHandler) MoveToken(c *gin.Context) {
	token, ok := h.getOwnedToken(c)
	if !ok {
		return
	}
	var req tokenMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.MoveToken(&token, req.X, req.Y); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

// DELETE /tokens/:id
func (h *SceneHandler) DeleteToken(c *gin.Context) {
	token, ok := h.getOwnedToken(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteToken(token.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token removido"})
}

type activeSceneRequest struct {
	SceneID uint `json:"scene_id"`
}

// PATCH /sessions/:id/active-scene — troca o cenário ativo da sessão.
// Valida que o cenário pertence à mesma campanha da sessão (não só que o
// mestre é dono de ambos separadamente) pra não deixar apontar pra um
// cenário de outra campanha por engano.
func (h *SceneHandler) SetActiveScene(c *gin.Context) {
	sessionID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	session, err := h.sessionSvc.GetByID(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sessão não encontrada"})
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, session.CampaignID) {
		return
	}
	var req activeSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scene, err := h.svc.GetByID(req.SceneID)
	if err != nil || scene.CampaignID != session.CampaignID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cenário não pertence a essa campanha"})
		return
	}
	if err := h.sessionSvc.SetActiveScene(&session, req.SceneID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.wsManager.Broadcast(session.CampaignID, ws.Event{
		Type: "scene_changed",
		Data: gin.H{"session_id": session.ID, "scene_id": req.SceneID},
	})
	c.JSON(http.StatusOK, session)
}
