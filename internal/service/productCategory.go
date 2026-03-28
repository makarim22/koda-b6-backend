package service

import (
	"context"
	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type ProductCategoryService struct {
	categoryRepo *repository.ProductCategoryRepository
}

func NewProductCategoryService(categoryRepo *repository.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{categoryRepo: categoryRepo}
}

func (s *ProductCategoryService) Create(ctx context.Context, category *models.ProductCategory) error {
	if category.Name == "" {
		return errors.NewValidationError("category name", "cannot be empty")
	}

	return s.categoryRepo.Create(ctx, category)
}

func (s *ProductCategoryService) GetByID(ctx context.Context, id int) (*models.ProductCategory, error) {
	if id <= 0 {
		return nil, errors.NewValidationError("category ID", "must be positive")
	}

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, err
		}
		return nil, errors.NewServiceError("get category", err)
	}
	return category, nil
}

func (s *ProductCategoryService) GetAll(ctx context.Context) ([]models.ProductCategory, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.NewServiceError("get all categories", err)
	}
	return categories, nil
}

func (s *ProductCategoryService) Update(ctx context.Context, category *models.ProductCategory) error {
	if category.ID <= 0 {
		return errors.NewValidationError("category ID", "must be positive")
	}
	if category.Name == "" {
		return errors.NewValidationError("category name", "cannot be empty")
	}

	err := s.categoryRepo.Update(ctx, category)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return err
		}
		return errors.NewServiceError("update category", err)
	}
	return nil
}

func (s *ProductCategoryService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.NewValidationError("category ID", "must be positive")
	}

	err := s.categoryRepo.Delete(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return err
		}
		return errors.NewServiceError("delete category", err)
	}
	return nil
}
