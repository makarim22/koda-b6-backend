package handlers

import (
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VariantHandler struct {
	service *service.VariantService
}

func NewVariantHandler(service *service.VariantService) *VariantHandler {
	return &VariantHandler{
		service: service,
	}
}

func (h *VariantHandler) GetVariantsByProductID(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}

	if productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product ID must be greater than zero",
		})
		return
	}

	ctx := c.Request.Context()
	fmt.Println("productID:", productID)
	variants, err := h.service.GetVariantsByProductID(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch variants",
		})
		return
	}

	// Return empty array if no variants found instead of null
	if variants == nil {
		variants = []models.Variant{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": variants,
	})
}
