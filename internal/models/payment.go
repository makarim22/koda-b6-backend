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

type CreatePaymentRequest struct {
	OrderID int     `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,min=0"`
	Method  string  `json:"method" binding:"required"`
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type MidtransCallbackRequest struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	ApprovalCode      string `json:"approval_code"`
}
