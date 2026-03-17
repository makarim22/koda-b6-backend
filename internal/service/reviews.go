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

func (s *ReviewsService) GetAllReviews(ctx context.Context) ([]models.Reviews, error) {
	reviews, err := s.reviewsRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve reviews")
	}

	if len(reviews) == 0 {
		return []models.Reviews{}, nil
	}

	return reviews, nil
}
