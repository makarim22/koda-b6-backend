package service

import (
	"context"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type WishlistService struct {
	wishlistRepo *repository.WishlistRepository
}

func NewWishlistService(wishlistRepo *repository.WishlistRepository) *WishlistService {
	return &WishlistService{wishlistRepo: wishlistRepo}
}

func (s *WishlistService) AddToWishlist(ctx context.Context, customerID int, productID int) error {
	return s.wishlistRepo.Add(ctx, customerID, productID)
}

func (s *WishlistService) RemoveFromWishlist(ctx context.Context, customerID int, productID int) error {
	return s.wishlistRepo.Remove(ctx, customerID, productID)
}

func (s *WishlistService) GetUserWishlist(ctx context.Context, customerID int) ([]models.WishlistItemResponse, error) {
	return s.wishlistRepo.GetUserWishlist(ctx, customerID)
}

func (s *WishlistService) CheckStatus(ctx context.Context, customerID int, productID int) (bool, error) {
	return s.wishlistRepo.CheckStatus(ctx, customerID, productID)
}
