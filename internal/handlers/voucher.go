package handlers

import (
	"net/http"

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

	// Need to get user_id to ensure user is logged in (optional based on your biz logic)
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
