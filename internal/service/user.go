package service

import (
	"context"
	"errors"
	"fmt"
	"koda-b6-backend/internal/lib"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
	"strconv"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve users")
	}

	if len(users) == 0 {
		return []models.User{}, nil
	}

	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	if id < 0 {
		return nil, errors.New("Invalid user id")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("empty email")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	if user.Email == "" || user.Name == "" {
		return errors.New("email and full name are required")
	}

	existingUser, _ := s.userRepo.GetByID(ctx, user.ID)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *UserService) CreateUserNew(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Validate input
	if req.Email == "" || req.Name == "" {
		return nil, errors.New("email and name are required")
	}

	if !lib.IsValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	//if req.Password != req.ConfirmPassword {
	//	return nil, errors.New("make sure the confirm password is correct")
	//}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := lib.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save to repository
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return user without password hash
	return &models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	fmt.Println("user", user)
	if user.ID == 0 {
		return errors.New("invalid User Id")
	}

	if user.Email == "" || user.Name == "" {
		return errors.New("email and full name are required")
	}

	err := s.userRepo.Update(ctx, user)
	if err != nil {
		return errors.New("failed to update user")
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID cannot be empty")
	}

	idInt, err := strconv.Atoi(id)

	fmt.Println(idInt)

	err = s.userRepo.Delete(ctx, idInt)
	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}
