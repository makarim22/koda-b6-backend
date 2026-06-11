package handlers

import (
	"fmt"
	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	payment := &models.Payment{
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Method:  req.Method,
		Status:  "pending",
	}

	err := h.service.Create(ctx, payment)
	if err != nil {
		if ve, ok := err.(*errors.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	ctx := c.Request.Context()
	payment, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	ctx := c.Request.Context()
	payment, err := h.service.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	var req models.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		if ve, ok := err.(*errors.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Error()})
			return
		}
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment status updated"})
}

func (h *PaymentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.Delete(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete payment"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PaymentHandler) ProcessMidtransCallback(c *gin.Context) {
	var req models.MidtransCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.service.ProcessMidtransCallback(ctx, req)
	if err != nil {
		fmt.Printf("Midtrans Callback Error: %v\n", err)
		if ve, ok := err.(*errors.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process midtrans callback", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
