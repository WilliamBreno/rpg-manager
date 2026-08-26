package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpg-manager/internal/service"
)

type ItemHandler struct{ svc *service.ItemService }

func NewItemHandler(svc *service.ItemService) *ItemHandler { return &ItemHandler{svc: svc} }

// GET /items?edition=5e&category=arma
func (h *ItemHandler) GetAll(c *gin.Context) {
	items, err := h.svc.GetAll(c.Query("edition"), c.Query("category"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
