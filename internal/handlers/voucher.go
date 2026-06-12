package handlers

import (
	"net/http"
	"strconv"

	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	voucherService *service.VoucherService
}

func NewVoucherHandler(voucherService *service.VoucherService) *VoucherHandler {
	return &VoucherHandler{
		voucherService: voucherService,
	}
}

func (h *VoucherHandler) ValidateVoucher(c *gin.Context) {
	var req models.ValidateVoucherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidateVoucherResponse{
			Valid:   false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ValidateVoucherResponse{
			Valid:   false,
			Message: "Unauthorized",
		})
		return
	}

	_, discountAmount, err := h.voucherService.CalculateDiscount(ctx, req.Code, req.Subtotal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ValidateVoucherResponse{
			Valid:   false,
			Message: err.Error(),
		})
		return
	}

	finalTotal := req.Subtotal - discountAmount

	c.JSON(http.StatusOK, models.ValidateVoucherResponse{
		Valid:          true,
		Code:           req.Code,
		DiscountAmount: discountAmount,
		FinalTotal:     finalTotal,
		Message:        "Voucher applied successfully",
	})
}

func (h *VoucherHandler) GetAllVouchers(c *gin.Context) {
	vouchers, err := h.voucherService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vouchers", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vouchers)
}

func (h *VoucherHandler) CreateVoucher(c *gin.Context) {
	var v models.Voucher
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := h.voucherService.Create(c.Request.Context(), &v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create voucher", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *VoucherHandler) UpdateVoucher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var v models.Voucher
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	if err := h.voucherService.Update(c.Request.Context(), id, &v); err != nil {
		if err.Error() == "voucher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voucher", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voucher updated successfully"})
}

func (h *VoucherHandler) DeleteVoucher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.voucherService.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "voucher not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete voucher", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voucher deleted successfully"})
}
