package handlers

import (
	defaulterr "errors"
	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderDetailHandler struct {
	service *service.OrderDetailService
}

func NewOrderDetailHandler(service *service.OrderDetailService) *OrderDetailHandler {
	return &OrderDetailHandler{service: service}
}

func (h *OrderDetailHandler) Create(c *gin.Context) {
	var req models.CreateOrderDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	detail := &models.OrderDetail{
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		SizeID:    &req.SizeID,
		VariantID: &req.VariantID,
		Quantity:  req.Quantity,
		Price:     req.UnitPrice,
	}

	err := h.service.Create(ctx, detail)
	if err != nil {
		var validationErr *errors.ValidationError
		var notFoundErr *errors.NotFoundError

		if defaulterr.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		if defaulterr.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order detail"})
		return
	}

	c.JSON(http.StatusCreated, detail)
}

//func (h *OrderDetailHandler) GetByID(c *gin.Context) {
//	idParam := c.Param("id")
//	id, err := strconv.Atoi(idParam)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order detail id"})
//		return
//	}
//
//	ctx := c.Request.Context()
//
//	detail, err := h.service.GetByID(ctx, id)
//	if err != nil {
//		var notFoundErr *errors.NotFoundError
//
//		if errors.As(err, &notFoundErr) {
//			c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
//			return
//		}
//
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order detail"})
//		return
//	}
//
//	c.JSON(http.StatusOK, detail)
//}

func (h *OrderDetailHandler) GetByOrderID(c *gin.Context) {
	orderIDParam := c.Param("order_id")
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	ctx := c.Request.Context()

	details, err := h.service.GetByOrderID(ctx, orderID)
	if err != nil {
		var notFoundErr *errors.NotFoundError

		if defaulterr.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order details"})
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *OrderDetailHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order detail id"})
		return
	}

	var req models.UpdateOrderDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	detail := &models.OrderDetail{
		ID:       id,
		Quantity: req.Quantity,
		Price:    req.UnitPrice,
	}

	err = h.service.Update(ctx, detail)
	if err != nil {
		var validationErr *errors.ValidationError
		var notFoundErr *errors.NotFoundError

		if defaulterr.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		if defaulterr.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order detail"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (h *OrderDetailHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order detail id"})
		return
	}

	ctx := c.Request.Context()

	err = h.service.Delete(ctx, id)
	if err != nil {
		var notFoundErr *errors.NotFoundError

		if defaulterr.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete order detail"})
		return
	}

	c.Status(http.StatusNoContent)
}
