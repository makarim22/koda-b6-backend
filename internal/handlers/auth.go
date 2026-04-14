package handlers

import (
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewAuthHandler(userService *service.UserService, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}


	if req.Role == "" {
		req.Role = "user" 
	}
	if req.Role != "user" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid role. Must be 'user' or 'admin'",
		})
		return
	}

	ctx := c.Request.Context()

	user, err := h.authService.RegisterWithRole(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusConflict, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "User registered successfully",
		Data: models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Call service to login user
	user, token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data: models.LoginResponse{
			ID:    user.ID,
			Email: user.Email,
			Role: user.Role,
			Token: token,
		},
	})
}

