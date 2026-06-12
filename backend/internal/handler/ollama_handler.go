package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OllamaHandler struct{}

func NewOllamaHandler() *OllamaHandler {
	return &OllamaHandler{}
}

type SkillRequest struct {
	ClassName string `json:"class_name"`
	Edition   string `json:"edition"`
	Level     int    `json:"level"`
}

func (h *OllamaHandler) GetSkills(c *gin.Context) {
	var req SkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao montar requisição"})
		return
	}

	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	resp, err := client.Post("http://localhost:8001/skills", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Serviço de IA não está rodando."})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler resposta da IA"})
		return
	}

	// Repassa a resposta do Python direto para o front
	c.Data(resp.StatusCode, "application/json", respBody)
}