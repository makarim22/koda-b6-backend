package repository

import (
	"context"
	"errors"
	"fmt"
	"koda-b6-backend/internal/models"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

type ForgotPasswordRepository struct {
	db *pgx.Conn
}

func NewForgotPasswordRepository(db *pgx.Conn) *ForgotPasswordRepository {
	return &ForgotPasswordRepository{db: db}
}

func (r *ForgotPasswordRepository) CreateForgotPassword(ctx context.Context, email string) (models.ForgotPassword, error) {

	codeOTP := rand.Intn(900000) + 100000

	now := time.Now()
	expiredAt := now.Add(15 * time.Minute)

	const query = "INSERT INTO forgot_password (code_otp, email, expired_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, code_otp, email, expired_at, created_at, updated_at"

	var result models.ForgotPassword

	err := r.db.QueryRow(ctx, query, codeOTP, email, expiredAt, now, now).Scan(
		&result.ID,
		&result.CodeOTP,
		&result.Email,
		&result.ExpiredAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		log.Printf("CreateForgotPassword error: %v", err)
		return models.ForgotPassword{}, fmt.Errorf("failed to create forgot password record: %w", err)
	}

	log.Printf("Forgot password record created for email: %s", email)
	return result, nil

}

func (r *ForgotPasswordRepository) GetDataByEmail(ctx context.Context, email string) (models.ForgotPassword, error) {
	var result models.ForgotPassword
	err := r.db.QueryRow(ctx,
		"SELECT id, code_otp, email, expired_at, created_at, updated_at FROM forgot_password WHERE email = $1",
		email,
	).Scan(&result.ID, &result.CodeOTP, &result.Email, &result.ExpiredAt, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ForgotPassword{}, errors.New("OTP record not found for email: " + email)
		}
		log.Printf("GetDataByEmail error: %v", err)
		return models.ForgotPassword{}, err
	}

	return result, nil

}

func (r *ForgotPasswordRepository) DeleteDataByCode(ctx context.Context, codeOTP string) error {
	const query = "DELETE FROM forgot_password WHERE code_otp = $1"

	result, err := r.db.Exec(ctx, query, codeOTP)
	if err != nil {
		log.Printf("DeleteDataByCode Exec error: %v", err)
		return fmt.Errorf("failed to delete OTP record: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("OTP record not found")
	}

	log.Printf("Successfully deleted OTP record. Rows affected: %d", result.RowsAffected())
	return nil

}
