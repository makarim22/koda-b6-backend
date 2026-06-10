package handlers

import (
	"net/http"
	"strconv"

	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	wishlistService *service.WishlistService
}

func NewWishlistHandler(wishlistService *service.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	var req models.AddWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	err := h.wishlistService.AddToWishlist(ctx, customerID.(int), req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "Product added to wishlist",
	})
}

func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid product ID",
		})
		return
	}

	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	err = h.wishlistService.RemoveFromWishlist(ctx, customerID.(int), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Product removed from wishlist",
	})
}

func (h *WishlistHandler) GetUserWishlist(c *gin.Context) {
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	items, err := h.wishlistService.GetUserWishlist(ctx, customerID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Wishlist retrieved successfully",
		Data:    items,
	})
}

func (h *WishlistHandler) CheckStatus(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid product ID",
		})
		return
	}

	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	isFav, err := h.wishlistService.CheckStatus(ctx, customerID.(int), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Status checked",
		Data:    models.WishlistStatusResponse{IsFavorite: isFav},
	})
}
