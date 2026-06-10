package service

import (
	"context"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

type PointService interface {
	AddPoints(ctx context.Context, tx pgx.Tx, userID int, orderID *int, points int, description string) error
	DeductPoints(ctx context.Context, tx pgx.Tx, userID int, orderID *int, points int, description string) error
	GetPointHistory(ctx context.Context, userID int, limit, offset int) ([]models.PointLedger, error)
	GetUserBalance(ctx context.Context, userID int) (int, error)
}

type pointService struct {
	pointRepo repository.PointRepository
}

func NewPointService(pointRepo repository.PointRepository) PointService {
	return &pointService{pointRepo: pointRepo}
}

func (s *pointService) AddPoints(ctx context.Context, tx pgx.Tx, userID int, orderID *int, points int, description string) error {
	if points <= 0 {
		return nil
	}
	return s.pointRepo.AddPoints(ctx, tx, userID, orderID, points, description)
}

func (s *pointService) DeductPoints(ctx context.Context, tx pgx.Tx, userID int, orderID *int, points int, description string) error {
	if points <= 0 {
		return nil
	}
	// Deduct is just adding negative points
	return s.pointRepo.AddPoints(ctx, tx, userID, orderID, -points, description)
}

func (s *pointService) GetPointHistory(ctx context.Context, userID int, limit, offset int) ([]models.PointLedger, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.pointRepo.GetPointHistory(ctx, userID, limit, offset)
}

func (s *pointService) GetUserBalance(ctx context.Context, userID int) (int, error) {
	return s.pointRepo.GetUserBalance(ctx, userID)
}
