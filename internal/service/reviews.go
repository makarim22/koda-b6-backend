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
	if review.Message == "" || review.Rating <= 0 || review.Rating > 5 {
		return errors.New("review message cannot be empty and rating must be between 1 and 5")
	}

	existingProduct, _ := s.productRepo.GetByID(ctx, review.ProductId)
	if existingProduct == nil {
		return errors.New("product with this id IS NOT exists")
	}

	// Find an eligible order for this user and product
	orderID, err := s.reviewsRepo.GetEligibleOrderForReview(ctx, review.UserId, review.ProductId)
	if err != nil {
		return fmt.Errorf("error checking eligibility: %w", err)
	}
	if orderID == 0 {
		return errors.New("user is not eligible to review this product")
	}
	review.OrderId = orderID

	err = s.reviewsRepo.CreateReview(ctx, review)
	if err != nil {
		return errors.New("failed to create review")
	}
	return nil
}

func (s *ReviewsService) UpdateReview(ctx context.Context, review *models.ReviewsRequest) error {
	if review.Message == "" || review.Rating <= 0 {
		return errors.New("review message cannot be empty and rating must be greater than zero")
	}
	existingReview, _ := s.reviewsRepo.GetById(ctx, review.Id)
	if existingReview == nil {
		return errors.New("review with this id IS NOT exists")
	}
	err := s.reviewsRepo.UpdateReview(ctx, review)
	if err != nil {
		return errors.New("failed to update review")
	}
	return nil
}

func (s *ReviewsService) DeleteReview(ctx context.Context, id int) error {
	err := s.reviewsRepo.DeleteReview(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ReviewsService) GetByProductId(ctx context.Context, productID int, limit, offset int) ([]models.ReviewsResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	reviews, err := s.reviewsRepo.GetByProductId(ctx, productID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve reviews by product id")
	}
	if len(reviews) == 0 {
		return []models.ReviewsResponse{}, nil
	}
	return reviews, nil
}

func (s *ReviewsService) GetRatingSummary(ctx context.Context, productID int) (*models.RatingSummary, error) {
	summary, err := s.reviewsRepo.GetRatingSummary(ctx, productID)
	if err != nil {
		return nil, errors.New("failed to retrieve rating summary")
	}
	return summary, nil
}

func (s *ReviewsService) CheckEligible(ctx context.Context, userID, productID int) (bool, error) {
	orderID, err := s.reviewsRepo.GetEligibleOrderForReview(ctx, userID, productID)
	if err != nil {
		return false, err
	}
	return orderID != 0, nil
}
