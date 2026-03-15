package handlers

import (
	"net/http"
	"strconv"

	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	var req models.AddToCartRequest

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

	cart, err := h.cartService.AddToCart(ctx, customerID.(int), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Product added to cart",
		Data:    cart,
	})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()

	cart, err := h.cartService.GetCart(ctx, customerID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Cart retrieved successfully",
		Data:    cart,
	})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	cartItemID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid cart item ID",
		})
		return
	}

	var req models.UpdateCartItemRequest

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

	cart, err := h.cartService.UpdateCartItemQuantity(ctx, customerID.(int), cartItemID, req)
	if err != nil {
		if err.Error() == "unauthorized access to cart item" {
			c.JSON(http.StatusForbidden, models.ApiResponse{
				Success: false,
				Message: "You don't have access to this cart item",
			})
			return
		}

		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Cart item updated",
		Data:    cart,
	})
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	cartItemID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid cart item ID",
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

	cart, err := h.cartService.RemoveFromCart(ctx, customerID.(int), cartItemID)
	if err != nil {
		if err.Error() == "unauthorized access to cart item" {
			c.JSON(http.StatusForbidden, models.ApiResponse{
				Success: false,
				Message: "You don't have access to this cart item",
			})
			return
		}

		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Item removed from cart",
		Data:    cart,
	})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()

	if err := h.cartService.ClearCart(ctx, customerID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Cart cleared successfully",
	})
}

//func (h *CartHandler) GetCartSummary(c *gin.Context) {
//	// Get customer ID from JWT context
//	customerID, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, models.ApiResponse{
//			Success: false,
//			Message: "Unauthorized",
//		})
//		return
//	}
//
//	ctx := c.Request.Context()
//
//	summary, err := h.cartService.GetCartSummary(ctx, customerID.(int))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, models.ApiResponse{
//			Success: false,
//			Message: err.Error(),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, models.ApiResponse{
//		Success: true,
//		Message: "Cart summary retrieved",
//		Data:    summary,
//	})
//}
