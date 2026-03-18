package service

import (
	"context"
	"errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type ReviewsService struct {
	reviewsRepo *repository.ReviewsRepository
}

func NewReviewsService(reviewsRepo *repository.ReviewsRepository) *ReviewsService {
	return &ReviewsService{
		reviewsRepo: reviewsRepo,
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
