package handlers

import (
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductImageHandler struct {
	service *service.ProductImageService
}

func NewProductImageHandler(service *service.ProductImageService) *ProductImageHandler {
	return &ProductImageHandler{
		service: service,
	}
}

func (h *ProductImageHandler) GetImageByID(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}
	if productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product ID must be greater than zero",
		})
		return
	}

	ctx := c.Request.Context()
	fmt.Println("productID:", productID)
	images, err := h.service.GetImagesByProductID(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"images": images,
	})

}

func (h *ProductImageHandler) parseProductID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return 0, false
	}
	return id, true
}

func (h *ProductImageHandler) parseImageID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("imageId"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id"})
		return 0, false
	}
	return id, true
}

// GET /api/products/:id/images
func (h *ProductImageHandler) GetImages(c *gin.Context) {
	productID, ok := h.parseProductID(c)
	if !ok {
		return
	}

	images, err := h.service.GetImagesByProductID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "images fetched successfully",
		"data":    images,
	})
}

// POST /api/products/:id/images
func (h *ProductImageHandler) UploadImage(c *gin.Context) {
	productID, ok := h.parseProductID(c)
	if !ok {
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large or invalid form"})
		return
	}

	_, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	isPrimary := c.PostForm("is_primary") == "true"

	image, err := h.service.UploadImage(c.Request.Context(), models.UploadImageInput{
		ProductID: productID,
		Header:    header,
		IsPrimary: isPrimary,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "image uploaded successfully",
		"data":    image,
	})
}

// POST /api/products/:id/images/multiple
func (h *ProductImageHandler) UploadMultipleImages(c *gin.Context) {
	productID, ok := h.parseProductID(c)
	if !ok {
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no images provided"})
		return
	}

	saved, err := h.service.UploadMultipleImages(c.Request.Context(), models.UploadMultipleInput{
		ProductID: productID,
		Files:     files,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "images uploaded successfully",
		"data":    saved,
	})
}

// PATCH /api/products/:id/images/:imageId/set-primary
func (h *ProductImageHandler) SetPrimaryImage(c *gin.Context) {
	productID, ok := h.parseProductID(c)
	if !ok {
		return
	}
	imageID, ok := h.parseImageID(c)
	if !ok {
		return
	}

	if err := h.service.SetPrimaryImage(c.Request.Context(), imageID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "primary image updated successfully"})
}

// DELETE /api/products/:id/images/:imageId
func (h *ProductImageHandler) DeleteImage(c *gin.Context) {
	productID, ok := h.parseProductID(c)
	if !ok {
		return
	}
	imageID, ok := h.parseImageID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteImage(c.Request.Context(), imageID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image deleted successfully"})
}
