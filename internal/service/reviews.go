package service

import (
	"context"
	"errors"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type ReviewsService struct {
	reviewsRepo *repository.ReviewsRepository
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewReviewsService(reviewsRepo *repository.ReviewsRepository, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *ReviewsService {
	return &ReviewsService{
		reviewsRepo: reviewsRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *ReviewsService) GetAllReviews(ctx context.Context) ([]models.ReviewsResponse, error) {
	reviews, err := s.reviewsRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve reviews")
	}

	if len(reviews) == 0 {
		return []models.ReviewsResponse{}, nil
	}

	return reviews, nil
}

func (s *ReviewsService) GetById(ctx context.Context, id int) (models.ReviewsResponse, error) {
	review, err := s.reviewsRepo.GetById(ctx, id)
	if err != nil {
		return models.ReviewsResponse{}, err
	}
	return *review, nil
}

func (s *ReviewsService) CreateReview(ctx context.Context, review *models.ReviewsRequest) error {
	if review.Message == "" || review.Rating <= 0 {
		return errors.New("review message cannot be empty and rating must be greater than zero")
	}

	existingOrder, _ := s.orderRepo.GetOrderByID(ctx, review.OrderId)
	if existingOrder == nil {
		return errors.New("order with this id IS NOT exists")
	}
	existingProduct, _ := s.productRepo.GetByID(ctx, review.ProductId)
	if existingProduct == nil {
		return errors.New("product with this id IS NOT exists")
	}

	existingReview, err := s.reviewsRepo.GetByUserProductOrder(ctx, review.UserId, review.ProductId, review.OrderId)
	if err != nil {
		return fmt.Errorf("failed to check existing review: %w", err)
	}
	if existingReview != nil {
		return errors.New("review already exists for this user, product, and order combination")
	}

	err = s.reviewsRepo.CreateReview(ctx, review)
	if err != nil {
		return errors.New("failed to create review")
	}
	return nil
}
