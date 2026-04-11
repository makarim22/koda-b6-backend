package service

import (
	"context"
	"errors"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type ProductImageService struct {
	repo      *repository.ProductImageRepository
	uploadDir string
}

func NewProductImageService(repo *repository.ProductImageRepository) *ProductImageService {
	return &ProductImageService{
		repo:      repo,
		uploadDir: "upload/products",
	}
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func (s *ProductImageService) GetImagesByProductID(ctx context.Context, productId int) ([]models.ProductImage, error) {
	if productId <= 0 {
		return nil, errors.New("productId must be positive")
	}
	images, err := s.repo.GetByProductImageID(ctx, productId)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (s *ProductImageService) saveFile(header *multipart.FileHeader, productID, index int) (string, error) {
	ext := filepath.Ext(header.Filename)
	if !allowedExtensions[ext] {
		return "", errors.New("only jpg, jpeg, png, webp are allowed")
	}

	if err := os.MkdirAll(s.uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filename := fmt.Sprintf("%d_%d_%d%s", productID, time.Now().UnixNano(), index, ext)
	savePath := filepath.Join(s.uploadDir, filename)

	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	buf := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err := dst.Write(buf[:n])
			if err != nil {
				return "", err
			}
		}
		if err != nil {
			break
		}
	}

	return "/" + savePath, nil
}

func (s *ProductImageService) UploadImage(ctx context.Context, input models.UploadImageInput) (*models.ProductImage, error) {
	savePath, err := s.saveFile(input.Header, input.ProductID, 0)
	if err != nil {
		return nil, err
	}

	// Auto-primary if this product has no images yet
	count, err := s.repo.CountByProductID(ctx, input.ProductID)
	if err != nil {
		err := os.Remove("." + savePath)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	if count == 0 {
		input.IsPrimary = true
	}

	// Unset existing primary before assigning new one
	if input.IsPrimary {
		if err := s.repo.UnsetPrimary(ctx, input.ProductID); err != nil {
			err := os.Remove("." + savePath)
			if err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	image := &models.ProductImage{
		ProductID: input.ProductID,
		Path:      savePath,
		IsPrimary: input.IsPrimary,
	}

	if err := s.repo.Save(ctx, image); err != nil {
		err := os.Remove("." + savePath)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return image, nil
}

func (s *ProductImageService) UploadMultipleImages(ctx context.Context, input models.UploadMultipleInput) ([]models.ProductImage, error) {
	if len(input.Files) == 0 {
		return nil, errors.New("no images provided")
	}

	count, err := s.repo.CountByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}

	var saved []models.ProductImage

	for i, fileHeader := range input.Files {
		savePath, err := s.saveFile(fileHeader, input.ProductID, i)
		if err != nil {
			continue // skip invalid files, don't abort the whole batch
		}

		isPrimary := i == 0 && count == 0

		image := models.ProductImage{
			ProductID: input.ProductID,
			Path:      savePath,
			IsPrimary: isPrimary,
		}

		if err := s.repo.Save(ctx, &image); err != nil {
			err := os.Remove("." + savePath)
			if err != nil {
				return nil, err
			}
			continue
		}

		saved = append(saved, image)
	}

	if len(saved) == 0 {
		return nil, errors.New("no valid images were uploaded")
	}

	return saved, nil
}

func (s *ProductImageService) SetPrimaryImage(ctx context.Context, imageID, productID int) error {
	image, err := s.repo.FindByID(ctx, imageID, productID)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image with ID %d not found", imageID)
	}

	if err := s.repo.UnsetPrimary(ctx, productID); err != nil {
		return err
	}

	return s.repo.SetPrimary(ctx, imageID, productID)
}

func (s *ProductImageService) DeleteImage(ctx context.Context, imageID, productID int) error {
	image, err := s.repo.FindByID(ctx, imageID, productID)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image with ID %d not found", imageID)
	}

	err = os.Remove("." + image.Path)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, imageID, productID); err != nil {
		return err
	}

	// Promote next image if deleted one was primary
	if image.IsPrimary {
		err := s.repo.PromoteNextPrimary(ctx, productID)
		if err != nil {
			return err
		}
	}

	return nil
}
