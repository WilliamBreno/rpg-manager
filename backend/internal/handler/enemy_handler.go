package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
	"rpg-manager/internal/ws"
)

type EnemyHandler struct {
	svc         *service.EnemyService
	campaignSvc *service.CampaignService
	wsManager   *ws.Manager
}

func NewEnemyHandler(svc *service.EnemyService, campaignSvc *service.CampaignService, wsManager *ws.Manager) *EnemyHandler {
	return &EnemyHandler{svc: svc, campaignSvc: campaignSvc, wsManager: wsManager}
}

type enemyAbilityInput struct {
	Name        string `json:"name"`
	Damage      string `json:"damage"`
	Description string `json:"description"`
}

type enemyLineInput struct {
	Text     string                 `json:"text"`
	AudioURL string                 `json:"audio_url"`
	Source   domain.EnemyLineSource `json:"source"`
}

type enemyRequest struct {
	Kind            domain.EnemyKind    `json:"kind"`
	Name            string              `json:"name"`
	HP              int                 `json:"hp"`
	ChallengeRating string              `json:"challenge_rating"`
	Race            string              `json:"race"`
	PhotoURL        string              `json:"photo_url"`
	SoundURL        string              `json:"sound_url"`
	Class           string              `json:"class"`
	Armor           int                 `json:"armor"`
	History         string              `json:"history"`
	Bonds           string              `json:"bonds"`
	Notes           string              `json:"notes"`
	Abilities       []enemyAbilityInput `json:"abilities"`
	Lines           []enemyLineInput    `json:"lines"`
}

func buildEnemy(campaignID uint, req enemyRequest) domain.Enemy {
	e := domain.Enemy{
		CampaignID: campaignID, Kind: req.Kind, Name: req.Name, HP: req.HP,
		ChallengeRating: req.ChallengeRating, Race: req.Race, PhotoURL: req.PhotoURL,
		SoundURL: req.SoundURL, Class: req.Class, Armor: req.Armor,
		History: req.History, Bonds: req.Bonds, Notes: req.Notes,
	}
	for _, a := range req.Abilities {
		e.Abilities = append(e.Abilities, domain.EnemyAbility{Name: a.Name, Damage: a.Damage, Description: a.Description})
	}
	for _, l := range req.Lines {
		e.Lines = append(e.Lines, domain.EnemyLine{Text: l.Text, AudioURL: l.AudioURL, Source: l.Source})
	}
	return e
}

// POST /campaigns/:id/enemies
func (h *EnemyHandler) Create(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	var req enemyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	enemy := buildEnemy(campaignID, req)
	warnings, err := h.svc.Create(&enemy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"enemy": enemy, "warnings": warnings})
}

// GET /campaigns/:id/enemies
func (h *EnemyHandler) GetByCampaign(c *gin.Context) {
	campaignID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, campaignID) {
		return
	}
	enemies, err := h.svc.GetByCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enemies)
}

func (h *EnemyHandler) getOwned(c *gin.Context) (domain.Enemy, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.Enemy{}, false
	}
	enemy, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inimigo não encontrado"})
		return domain.Enemy{}, false
	}
	if !requireCampaignOwner(c, h.campaignSvc, enemy.CampaignID) {
		return domain.Enemy{}, false
	}
	return enemy, true
}

// GET /enemies/:id
func (h *EnemyHandler) GetByID(c *gin.Context) {
	enemy, ok := h.getOwned(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, enemy)
}

// PUT /enemies/:id
func (h *EnemyHandler) Update(c *gin.Context) {
	enemy, ok := h.getOwned(c)
	if !ok {
		return
	}
	var req enemyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	enemy.Kind, enemy.Name, enemy.HP = req.Kind, req.Name, req.HP
	enemy.ChallengeRating, enemy.Race = req.ChallengeRating, req.Race
	enemy.PhotoURL, enemy.SoundURL, enemy.Class, enemy.Armor = req.PhotoURL, req.SoundURL, req.Class, req.Armor
	enemy.History, enemy.Bonds, enemy.Notes = req.History, req.Bonds, req.Notes
	if err := h.svc.Update(&enemy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enemy)
}

// DELETE /enemies/:id
func (h *EnemyHandler) Delete(c *gin.Context) {
	enemy, ok := h.getOwned(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(enemy.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inimigo removido"})
}

// POST /enemies/:id/abilities
func (h *EnemyHandler) CreateAbility(c *gin.Context) {
	enemyID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	enemy, err := h.svc.GetByID(enemyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inimigo não encontrado"})
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, enemy.CampaignID) {
		return
	}
	var req enemyAbilityInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ability := domain.EnemyAbility{EnemyID: enemyID, Name: req.Name, Damage: req.Damage, Description: req.Description}
	warning, err := h.svc.CreateAbility(&ability)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ability": ability, "warning": warning})
}

func (h *EnemyHandler) getOwnedAbility(c *gin.Context) (domain.EnemyAbility, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.EnemyAbility{}, false
	}
	ability, err := h.svc.GetAbilityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habilidade não encontrada"})
		return domain.EnemyAbility{}, false
	}
	enemy, err := h.svc.GetByID(ability.EnemyID)
	if err != nil || !requireCampaignOwner(c, h.campaignSvc, enemy.CampaignID) {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Habilidade não encontrada"})
		}
		return domain.EnemyAbility{}, false
	}
	return ability, true
}

// PUT /enemy-abilities/:id
func (h *EnemyHandler) UpdateAbility(c *gin.Context) {
	ability, ok := h.getOwnedAbility(c)
	if !ok {
		return
	}
	var req enemyAbilityInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ability.Name, ability.Damage, ability.Description = req.Name, req.Damage, req.Description
	warning, err := h.svc.UpdateAbility(&ability)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ability": ability, "warning": warning})
}

// DELETE /enemy-abilities/:id
func (h *EnemyHandler) DeleteAbility(c *gin.Context) {
	ability, ok := h.getOwnedAbility(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteAbility(ability.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Habilidade removida"})
}

// POST /enemies/:id/lines
func (h *EnemyHandler) CreateLine(c *gin.Context) {
	enemyID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	enemy, err := h.svc.GetByID(enemyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inimigo não encontrado"})
		return
	}
	if !requireCampaignOwner(c, h.campaignSvc, enemy.CampaignID) {
		return
	}
	var req enemyLineInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	line := domain.EnemyLine{EnemyID: enemyID, Text: req.Text, AudioURL: req.AudioURL, Source: req.Source}
	if err := h.svc.CreateLine(&line); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, line)
}

func (h *EnemyHandler) getOwnedLine(c *gin.Context) (domain.EnemyLine, bool) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return domain.EnemyLine{}, false
	}
	line, err := h.svc.GetLineByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fala não encontrada"})
		return domain.EnemyLine{}, false
	}
	enemy, err := h.svc.GetByID(line.EnemyID)
	if err != nil || !requireCampaignOwner(c, h.campaignSvc, enemy.CampaignID) {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Fala não encontrada"})
		}
		return domain.EnemyLine{}, false
	}
	return line, true
}

// PUT /enemy-lines/:id
func (h *EnemyHandler) UpdateLine(c *gin.Context) {
	line, ok := h.getOwnedLine(c)
	if !ok {
		return
	}
	var req enemyLineInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	line.Text, line.AudioURL, line.Source = req.Text, req.AudioURL, req.Source
	if err := h.svc.UpdateLine(&line); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, line)
}

// DELETE /enemy-lines/:id
func (h *EnemyHandler) DeleteLine(c *gin.Context) {
	line, ok := h.getOwnedLine(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteLine(line.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fala removida"})
}

// POST /enemies/:id/play-sound — o mestre traz o inimigo à cena e aciona o
// som que ele faz (Etapa 9). Só dispara o broadcast, não altera nada — dá
// 400 se o inimigo não tiver SoundURL preenchida (nada pra tocar).
func (h *EnemyHandler) PlaySound(c *gin.Context) {
	enemy, ok := h.getOwned(c)
	if !ok {
		return
	}
	if enemy.SoundURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Esse inimigo não tem som cadastrado"})
		return
	}
	h.wsManager.Broadcast(enemy.CampaignID, ws.Event{
		Type: "play_audio",
		Data: gin.H{"kind": "enemy_sound", "url": enemy.SoundURL, "enemy_name": enemy.Name},
	})
	c.JSON(http.StatusOK, gin.H{"message": "Som disparado"})
}

// POST /enemy-lines/:id/play — o mestre aciona uma fala de Boss/Vilão
// (Etapa 9), pros jogadores conectados ouvirem.
func (h *EnemyHandler) PlayLine(c *gin.Context) {
	line, ok := h.getOwnedLine(c)
	if !ok {
		return
	}
	if line.AudioURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Essa fala não tem áudio cadastrado"})
		return
	}
	enemy, err := h.svc.GetByID(line.EnemyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inimigo não encontrado"})
		return
	}
	h.wsManager.Broadcast(enemy.CampaignID, ws.Event{
		Type: "play_audio",
		Data: gin.H{"kind": "line", "url": line.AudioURL, "text": line.Text, "enemy_name": enemy.Name},
	})
	c.JSON(http.StatusOK, gin.H{"message": "Fala disparada"})
}

// GET /dnd/cr-damage-table — tabela pública "Estatísticas de Monstro por
// Nível de Desafio" (Guia do Mestre p.275), pra UI mostrar a faixa sugerida
// de dano por rodada conforme o mestre digita o ND do inimigo.
func CRDamageTableHandler(c *gin.Context) {
	c.JSON(http.StatusOK, service.CRDamageTable)
}
