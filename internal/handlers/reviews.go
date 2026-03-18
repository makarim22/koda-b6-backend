package handlers

import (
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

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

func (h *ReviewsHandler) GetReview(c *gin.Context) {
	reviewID := c.Param("id")
	id, err := strconv.Atoi(reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	review, err := h.reviewsService.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengambil data review",
		"data":    review,
	})
}

func (h *ReviewsHandler) CreateReview(c *gin.Context) {
	var review models.ReviewsRequest
	if err := c.BindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.reviewsService.CreateReview(c.Request.Context(), &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil membuat data review",
		"data":    review,
	})
}
