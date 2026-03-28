package models

import "time"

type Payment struct {
	ID            int        `json:"id" db:"id"`
	OrderID       int        `json:"order_id" db:"order_id"`
	Method        string     `json:"method" db:"method"`
	Amount        float64    `json:"amount" db:"amount"`
	Status        string     `json:"status" db:"status"`
	TransactionID *string    `json:"transaction_id,omitempty" db:"transaction_id"`
	PaymentDate   *time.Time `json:"payment_date,omitempty" db:"payment_date"`
}

type PaymentRequest struct {
	OrderID int    `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

type PaymentResponse struct {
	ID            int        `json:"id"`
	OrderID       int        `json:"order_id"`
	Method        string     `json:"method"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	TransactionID *string    `json:"transaction_id,omitempty"`
	PaymentDate   *time.Time `json:"payment_date,omitempty"`
}
