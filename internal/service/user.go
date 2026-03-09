package service

import (
	"context"
	"errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
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

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	if user.ID < 0  {
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

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}