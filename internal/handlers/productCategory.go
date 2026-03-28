package handlers

import (
	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductCategoryHandler struct {
	service *service.ProductCategoryService
}

func NewProductCategoryHandler(service *service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		service: service,
	}
}

func (h *ProductCategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	category := &models.ProductCategory{Name: req.Name}
	err := h.service.Create(ctx, category)
	if err != nil {
		switch {
		case errors.IsValidationError(err):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.IsConflictError(err):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *ProductCategoryHandler) GetCategoryByID(c *gin.Context) {
	categoryID := c.Param("id")

	ctx := c.Request.Context()

	categoryIDInt, err := strconv.Atoi(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID format - must be a number",
		})
		return
	}

	category, err := h.service.GetByID(ctx, categoryIDInt)
	if err != nil {
		switch {
		case errors.IsNotFoundError(err):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *ProductCategoryHandler) GetAllCategories(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *ProductCategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")

	categoryIDInt, err := strconv.Atoi(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID format - must be a number",
		})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	category := &models.ProductCategory{
		ID:   categoryIDInt,
		Name: req.Name,
	}
	err = h.service.Update(ctx, category)
	if err != nil {
		switch {
		case errors.IsValidationError(err):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.IsNotFoundError(err):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *ProductCategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	categoryIDInt, err := strconv.Atoi(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID format - must be a number",
		})
		return
	}

	ctx := c.Request.Context()

	err = h.service.Delete(ctx, categoryIDInt)
	if err != nil {
		switch {
		case errors.IsNotFoundError(err):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
