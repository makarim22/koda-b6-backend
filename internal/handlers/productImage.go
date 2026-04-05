package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"
)

type ProductImageHandler struct {
	service *service.ProductImageService
}

func NewProductImageHandler(service *service.ProductImageService) *ProductImageHandler {
	return &ProductImageHandler{
		service: service,
	}
}

func (h *ProductImageHandler) GetImageByID(c *gin.Context) {
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
	images, err := h.service.GetImagesByProductID(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"images": images,
	})

}
