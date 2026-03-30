package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type VariantService struct {
	variantRepo *repository.VariantRepository
}

func NewVariantService(variantRepo *repository.VariantRepository) *VariantService {
	return &VariantService{
		variantRepo: variantRepo,
	}
}

func (s *VariantService) GetVariantsByProductID(ctx context.Context, productID int) ([]models.Variant, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("invalid product ID: %d", productID)
	}
	variants, err := s.variantRepo.GetVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variants for product: %w", err)
	}
	return variants, nil
}
