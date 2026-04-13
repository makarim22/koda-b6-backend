package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type SizeService struct {
	repo *repository.SizeRepository
}

func NewSizeService(repo *repository.SizeRepository) *SizeService {
	return &SizeService{
		repo: repo,
	}
}

func (s *SizeService) GetSizeByProductID(ctx context.Context, productID int) ([]models.Size, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("invalid product ID: %d", productID)
	}
	sizes, err := s.repo.GetSizesByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sizes for product: %w", err)
	}
	return sizes, nil
}

func (s *SizeService) CreateSize (ctx context.Context, size *models.Size) error {
    err := s.repo.CreateSize(ctx, size)
	if err != nil {
		return fmt.Errorf("cannot create size: %w", err)
	}
	return nil
}