package models

import "time"

type ForgotPassword struct {
	ID        int       `db:"id" json:"id"`
	CodeOTP   string    `db:"code_otp" json:"code_otp"`
	Email     string    `db:"email" json:"email"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
