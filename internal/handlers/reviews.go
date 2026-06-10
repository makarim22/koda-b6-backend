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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	switch v := userID.(type) {
	case int:
		review.UserId = v
	case float64:
		review.UserId = int(v)
	default:
		// Attempt to parse string or fallback
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id type invalid"})
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

func (h *ReviewsHandler) UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	id, err := strconv.Atoi(reviewID)
	if err != nil {
		return
	}
	var review models.ReviewsRequest
	if err := c.BindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
	}

	review.Id = id

	err = h.reviewsService.UpdateReview(c.Request.Context(), &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengubah data review",
		"data":    review,
	})
}

func (h *ReviewsHandler) GetByProductId(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, err := h.reviewsService.GetByProductId(c.Request.Context(), productID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reviews,
	})
}

func (h *ReviewsHandler) GetRatingSummary(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	summary, err := h.reviewsService.GetRatingSummary(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

func (h *ReviewsHandler) CheckEligible(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var uid int
	switch v := userID.(type) {
	case int:
		uid = v
	case float64:
		uid = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id type invalid"})
		return
	}

	eligible, err := h.reviewsService.CheckEligible(c.Request.Context(), uid, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"eligible": eligible,
	})
}
