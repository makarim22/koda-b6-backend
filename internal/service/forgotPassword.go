package service

import (
	"context"
	"errors"
	"fmt"
	"koda-b6-backend/internal/repository"
	"log"
)

type ForgotPasswordService struct {
	userRepo           *repository.UserRepository
	forgotPasswordRepo *repository.ForgotPasswordRepository
}

func NewForgotPasswordService(userRepository *repository.UserRepository, passwordRepository *repository.ForgotPasswordRepository) *ForgotPasswordService {
	return &ForgotPasswordService{
		userRepo:           userRepository,
		forgotPasswordRepo: passwordRepository,
	}
}

func (s *ForgotPasswordService) ForgotPassword(ctx context.Context, email string) error {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("ForgotPassword: tidak ada user dengan email tersebut %s", email)
		return errors.New("apabila email terdaftar di sistem kami, kode OTP akan dikirim via email")
	}

	forgotPassword, err := s.forgotPasswordRepo.CreateForgotPassword(ctx, email)
	if err != nil {
		log.Printf("ForgotPassword: gagal membuat kode OTP untuk user %d: %v", user.ID, err)
		return fmt.Errorf("gagal membuat kode OTP: %w", err)
	}

	fmt.Println("OTP codenya", forgotPassword.CodeOTP)

	return nil

}
