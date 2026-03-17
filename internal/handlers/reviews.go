package handlers

import (
	"koda-b6-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewsHandler struct {
	reviewsService *service.ReviewsService
}

func NewReviewsHandler(reviewsService *service.ReviewsService) *ReviewsHandler {
	return &ReviewsHandler{
		reviewsService: reviewsService,
	}
}

func (h *ReviewsHandler) GetAllReviews(c *gin.Context) {
	reviews, err := h.reviewsService.GetAllReviews(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengambil data reviews",
		"data":    reviews,
	})
}
