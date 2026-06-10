package models

import "time"

type Wishlist struct {
	ID         int       `db:"id" json:"id"`
	CustomerID int       `db:"customer_id" json:"customer_id"`
	ProductID  int       `db:"product_id" json:"product_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type AddWishlistRequest struct {
	ProductID int `json:"product_id" binding:"required"`
}

type WishlistItemResponse struct {
	ProductID   int      `json:"product_id"`
	ProductName string   `json:"product_name"`
	BasePrice   float64  `json:"base_price"`
	Image       *string  `json:"image,omitempty"`
	Description string   `json:"description,omitempty"`
	AddedAt     time.Time `json:"added_at"`
}

type WishlistStatusResponse struct {
	IsFavorite bool `json:"is_favorite"`
}
