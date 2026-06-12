package service

import (
	"context"
	"database/sql"
	"errors"
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
		Role:     models.RoleUser,
		
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Return user without password hash
	return &models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthService) RegisterWithRole(ctx context.Context, email, password, role string) (*models.User, error) {

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := lib.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", fmt.Errorf("invalid email or password")
		}
		return nil, "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user == nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Verify password using Argon2
	valid, err := lib.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to verify password: %w", err)
	}

	if !valid {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := lib.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Return user without password
	return &models.User{
		ID:    user.ID,
		Email: user.Email,
		Role: user.Role,
	}, token, nil
}

func (s *AuthService) OAuthLogin(ctx context.Context, email, name, profileImage string) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// Just to be safe, if GetByEmail returns a different error than NoRows, but actually we use pgxpool which returns pgx.ErrNoRows or standard err.
		// Wait, user.go repository does not return sql.ErrNoRows because it's pgx. 
		// Actually, GetByEmail in pgx returns pgx.ErrNoRows. If user is nil, it wasn't found.
	}

	if user == nil || err != nil {
		// User does not exist, create them
		dummyPassword, _ := lib.HashPassword("") // Set a dummy password
		newUser := &models.User{
			Name:         name,
			Email:        email,
			Password:     dummyPassword,
			Role:         models.RoleUser,
		}
		if profileImage != "" {
			newUser.ProfileImage = &profileImage
		}
		
		err = s.userRepo.Create(ctx, newUser)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create oauth user: %w", err)
		}
		user = newUser
	}

	// Generate JWT token
	token, err := lib.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, token, nil
}
