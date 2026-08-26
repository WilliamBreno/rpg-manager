package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/domain"
	"rpg-manager/internal/service"
)

type InventoryHandler struct {
	Service          *service.InventoryService
	CharacterService *service.CharacterService
}

func NewInventoryHandler(svc *service.InventoryService, characterService *service.CharacterService) *InventoryHandler {
	return &InventoryHandler{Service: svc, CharacterService: characterService}
}

// isOwner verifica que o personagem existe e pertence ao usuário autenticado.
// Retorna false (e já escreve a resposta HTTP) se a checagem falhar.
func (h *InventoryHandler) isOwner(c *gin.Context, characterID uint) bool {
	character, err := h.CharacterService.GetByID(characterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return false
	}
	if character.UserID != c.GetUint("userID") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
		return false
	}
	return true
}

// GET /characters/:id/inventory
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if !h.isOwner(c, uint(id)) {
		return
	}
	items, armors, err := h.Service.GetInventory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "armors": armors})
}

type purchaseRequest struct {
	Quantity int `json:"quantity"`
}

// POST /characters/:id/shop/items/:item_id
func (h *InventoryHandler) PurchaseItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if !h.isOwner(c, uint(id)) {
		return
	}
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item_id inválido"})
		return
	}
	var req purchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Quantity <= 0 {
		req.Quantity = 1
	}
	character, err := h.Service.PurchaseItem(uint(id), uint(itemID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

// POST /characters/:id/shop/armors/:armor_id
func (h *InventoryHandler) PurchaseArmor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if !h.isOwner(c, uint(id)) {
		return
	}
	armorID, err := strconv.ParseUint(c.Param("armor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "armor_id inválido"})
		return
	}
	var req purchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Quantity <= 0 {
		req.Quantity = 1
	}
	character, err := h.Service.PurchaseArmor(uint(id), uint(armorID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}

type setCurrencyRequest struct {
	CopperPieces   int `json:"copper_pieces"`
	SilverPieces   int `json:"silver_pieces"`
	ElectrumPieces int `json:"electrum_pieces"`
	GoldPieces     int `json:"gold_pieces"`
	PlatinumPieces int `json:"platinum_pieces"`
}

// PATCH /characters/:id/currency — só o Mestre pode conceder/ajustar moedas,
// de qualquer personagem (não há conceito de campanha nesse sistema hoje).
func (h *InventoryHandler) SetCurrency(c *gin.Context) {
	role, _ := c.Get("userRole")
	if role != domain.RoleMaster {
		c.JSON(http.StatusForbidden, gin.H{"error": "Só o Mestre pode conceder moedas"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var req setCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	character, err := h.Service.SetCurrency(uint(id), req.CopperPieces, req.SilverPieces, req.ElectrumPieces, req.GoldPieces, req.PlatinumPieces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, character)
}
