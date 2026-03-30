package handlers

import (
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductDiscountHandler struct {
	service *service.ProductDiscountService
}

func NewProductDiscountHandler(service *service.ProductDiscountService) *ProductDiscountHandler {
	return &ProductDiscountHandler{
		service: service,
	}
}

func (h *ProductDiscountHandler) GetDiscountsByProductID(c *gin.Context) {
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
	discounts, err := h.service.GetDiscountsByProductID(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch variants",
		})
		return
	}

	if discounts == nil {
		discounts = []models.ProductDiscount{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": discounts,
	})
}
