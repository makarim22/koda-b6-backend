package models

import "time"

type Voucher struct {
	ID            int       `db:"id" json:"id"`
	Code          string    `db:"code" json:"code"`
	DiscountType  string    `db:"discount_type" json:"discount_type"` // PERCENTAGE, FIXED
	DiscountValue float64   `db:"discount_value" json:"discount_value"`
	MinPurchase   float64   `db:"min_purchase" json:"min_purchase"`
	MaxDiscount   *float64  `db:"max_discount" json:"max_discount,omitempty"`
	ValidFrom     time.Time `db:"valid_from" json:"valid_from"`
	ValidUntil    time.Time `db:"valid_until" json:"valid_until"`
	UsageLimit    int       `db:"usage_limit" json:"usage_limit"`
	UsedCount     int       `db:"used_count" json:"used_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type ValidateVoucherRequest struct {
	Code     string  `json:"code" binding:"required"`
	Subtotal float64 `json:"subtotal" binding:"required,min=0"`
}

type ValidateVoucherResponse struct {
	Valid          bool    `json:"valid"`
	Code           string  `json:"code,omitempty"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalTotal     float64 `json:"final_total"`
	Message        string  `json:"message"`
}
