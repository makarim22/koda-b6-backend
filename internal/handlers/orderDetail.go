package handlers

import (
	"koda-b6-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderDetailHandler struct {
	service *service.OrderDetailService
}

func NewOrderDetailHandler(service *services.OrderDetailService) *OrderDetailHandler {
	return &OrderDetailHandler{service: service}
}

type CreateOrderDetailRequest struct {
	OrderID   string  `json:"order_id" binding:"required"`
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

func (h *OrderDetailHandler) CreateOrderDetail(c *gin.Context) {
	var req CreateOrderDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}

	ctx := c.Request.Context()
	detail, err := h.service.CreateOrderDetail(ctx, req.OrderID, req.ProductID, req.Quantity, req.Price)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, detail)
}

func (h *OrderDetailHandler) GetOrderDetailByID(c *gin.Context) {
	detailID := c.Param("id")
	ctx := c.Request.Context()

	detail, err := h.service.GetOrderDetailByID(ctx, detailID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (h *OrderDetailHandler) GetOrderDetails(c *gin.Context) {
	orderID := c.Param("order_id")
	ctx := c.Request.Context()

	details, err := h.service.GetOrderDetails(ctx, orderID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *OrderDetailHandler) DeleteOrderDetail(c *gin.Context) {
	detailID := c.Param("id")
	ctx := c.Request.Context()

	err := h.service.DeleteOrderDetail(ctx, detailID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
