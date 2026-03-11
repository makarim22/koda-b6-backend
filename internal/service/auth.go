package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/lib"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password using Argon2
	hashedPassword, err := lib.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Return user without password hash
	return &models.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
