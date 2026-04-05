package service

import (
	"context"
	"errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type ProductImageService struct {
	repo *repository.ProductImageRepository
}

func NewProductImageService(repo *repository.ProductImageRepository) *ProductImageService {
	return &ProductImageService{
		repo: repo,
	}
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
