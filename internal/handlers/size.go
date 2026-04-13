package handlers

import (
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SizeHandler struct {
	service *service.SizeService
}

func NewSizeHandler(service *service.SizeService) *SizeHandler {
	return &SizeHandler{
		service: service,
	}
}

func (h *SizeHandler) GetSizeByProductID(c *gin.Context) {
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
	sizes, err := h.service.GetSizeByProductID(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch sizes options",
		})
		return
	}

	if sizes == nil {
		sizes = []models.Size{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": sizes,
	})
}


func (h *SizeHandler) CreateSize(c *gin.Context){
	var size models.Size

	if err := c.ShouldBindJSON(&size); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

    err := h.service.CreateSize(c.Request.Context(), &size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    size,
	})
}

func (h *SizeHandler) GetAllSizes(c *gin.Context){

	ctx := c.Request.Context()
	sizes, err := h.service.GetAllSizes(ctx)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success retreiving sizes",
		"data":    sizes,
	})
}