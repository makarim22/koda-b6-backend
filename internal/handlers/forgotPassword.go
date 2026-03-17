package handlers

import (
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ForgotPasswordHandler struct {
	forgotPasswordService *service.ForgotPasswordService
}

func NewForgotPasswordHandler(forgotPassService *service.ForgotPasswordService) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		forgotPasswordService: forgotPassService,
	}
}

func (h *ForgotPasswordHandler) ResetPassword(c *gin.Context) {
	var reqPassword models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&reqPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := h.forgotPasswordService.ResetPassword(c.Request.Context(), reqPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "password berhasil diubah",
		"data":    reqPassword,
	})

}

func (h *ForgotPasswordHandler) ForgotPassword(c *gin.Context) {
	var reqEmail models.ForgotPasswordRequest

	fmt.Println("email", reqEmail.Email)

	if err := c.ShouldBindJSON(&reqEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := h.forgotPasswordService.ForgotPassword(c.Request.Context(), reqEmail.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "request berhasil dikirim",
		"data":    reqEmail,
	})

}
