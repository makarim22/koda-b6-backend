package models

import "time"

type Cart struct {
	ID         int       `db:"id" json:"id"`
	CustomerID int       `db:"customer_id" json:"customer_id"`
	ProductID  int       `db:"product_id" json:"product_id"`
	SizeID     *int      `db:"size_id" json:"size_id,omitempty"`
	VariantID  *int      `db:"variant_id" json:"variant_id,omitempty"`
	Quantity   int       `db:"quantity" json:"quantity"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type CartItem struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Subtotal    float64   `json:"subtotal"`
	SizeID      *int      `json:"size_id,omitempty"`
	SizeName    string    `json:"size_name,omitempty"`
	VariantID   *int      `json:"variant_id,omitempty"`
	VariantName string    `json:"variant_name,omitempty"`
	Image       *string   `json:"image,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	Subtotal   float64    `json:"subtotal"`
}

type AddToCartRequest struct {
	ProductID int  `json:"product_id" binding:"required,gt=0"`
	SizeID    *int `json:"size_id"`
	VariantID *int `json:"variant_id"`
	Quantity  int  `json:"quantity" binding:"required,gt=0,lte=100"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0,lte=100"`
}

type RemoveCartItemRequest struct {
	CartItemID int `json:"cart_item_id" binding:"required,gt=0"`
}
