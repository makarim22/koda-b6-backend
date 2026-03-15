package handlers

import (
	"net/http"
	"strconv"

	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest

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

	order, err := h.orderService.CreateOrder(ctx, customerID.(int), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid order ID",
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

	order, err := h.orderService.GetOrder(ctx, orderID, customerID.(int))
	if err != nil {
		if err.Error() == "unauthorized access to order" {
			c.JSON(http.StatusForbidden, models.ApiResponse{
				Success: false,
				Message: "You don't have access to this order",
			})
			return
		}

		c.JSON(http.StatusNotFound, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()

	orders, err := h.orderService.GetUserOrders(ctx, customerID.(int), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid order ID",
		})
		return
	}

	var req models.UpdateOrderStatusRequest

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

	if err := h.orderService.UpdateOrderStatus(ctx, orderID, customerID.(int), req.Status); err != nil {
		if err.Error() == "unauthorized access to order" {
			c.JSON(http.StatusForbidden, models.ApiResponse{
				Success: false,
				Message: "You don't have access to this order",
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
		Message: "Order status updated successfully",
	})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid order ID",
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

	if err := h.orderService.DeleteOrder(ctx, orderID, customerID.(int)); err != nil {
		if err.Error() == "unauthorized access to order" {
			c.JSON(http.StatusForbidden, models.ApiResponse{
				Success: false,
				Message: "You don't have access to this order",
			})
			return
		}

		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, models.ApiResponse{
				Success: false,
				Message: "Order not found",
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
		Message: "Order deleted successfully",
	})
}
