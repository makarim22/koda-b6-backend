package service

import (
	"context"
	"errors"
	"fmt"
	"koda-b6-backend/internal/lib"

	//"github.com/matthewhartstonge/argon2"
	"koda-b6-backend/internal/models"
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

func (s *ForgotPasswordService) ForgotPassword(ctx context.Context, email string) (models.ForgotPassword, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("ForgotPassword: tidak ada user dengan email tersebut %s", email)
		return models.ForgotPassword{}, errors.New("apabila email terdaftar di sistem kami, kode OTP akan dikirim via email")
	}

	forgotPassword, err := s.forgotPasswordRepo.CreateForgotPassword(ctx, email)
	if err != nil {
		log.Printf("ForgotPassword: gagal membuat kode OTP untuk user %d: %v", user.ID, err)
		return models.ForgotPassword{}, fmt.Errorf("gagal membuat kode OTP: %w", err)
	}

	fmt.Println("OTP codenya", forgotPassword.CodeOTP)

	return forgotPassword, nil

}

func (s *ForgotPasswordService) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {

	//if req.NewPassword != req.ConfirmPassword {
	//	return errors.New("password tidak sesuai")
	//}

	otp, err := s.forgotPasswordRepo.GetDataByEmail(ctx, req.Email)
	fmt.Println("otp", otp)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("ResetPassword: gagal menemukan user dengan email - email: %s", req.Email)
		return errors.New("user tidak ditemukan")
	}

	//todos
	//hashedPassword, err := argon2.Hash(req.NewPassword, argon2.DefaultConfig())
	//if err != nil {
	//	log.Printf("ResetPassword: failed to hash password - userID: %d", user.ID)
	//	return errors.New("failed to reset password")
	//}

	// Hash password
	hashedPassword, err := lib.HashPassword(req.NewPassword)
	fmt.Println("hashedPassword", hashedPassword)
	if err != nil {
		fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, hashedPassword, user.ID); err != nil {
		log.Printf("ResetPassword: gagal mengupdate password - userID: %d, error: %v", user.ID, err)
		return fmt.Errorf("gagal mengupdate password: %w", err)
	}

	_ = s.forgotPasswordRepo.DeleteDataByCode(ctx, req.CodeOTP)

	log.Printf("ResetPassword: success - userID: %d, email: %s", user.ID, req.Email)
	return nil

}
