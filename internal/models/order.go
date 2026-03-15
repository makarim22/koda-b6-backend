package models

import "time"

type Order struct {
	ID          int       `db:"id" json:"id"`
	CustomerID  int       `db:"customer_id" json:"customer_id"`
	OrderDate   time.Time `db:"order_date" json:"order_date"`
	Subtotal    float64   `db:"subtotal" json:"subtotal"`
	Tax         float64   `db:"tax" json:"tax"`
	DeliveryFee float64   `db:"delivery_fee" json:"delivery_fee"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type OrderDetail struct {
	ID               int     `db:"id" json:"id"`
	OrderID          int     `db:"order_id" json:"order_id"`
	ProductID        int     `db:"product_id" json:"product_id"`
	SizeID           *int    `db:"size_id" json:"size_id,omitempty"`
	TemperatureID    *int    `db:"temperature_id" json:"temperature_id,omitempty"`
	Quantity         int     `db:"quantity" json:"quantity"`
	Price            float64 `db:"unit_price" json:"price"`
	ProductName      string  `json:"product_name,omitempty"`
	SizeName         string  `json:"size_name,omitempty"`
	TemperatureLabel string  `json:"temperature_label,omitempty"`
}

type CreateOrderRequest struct {
	DeliveryFee float64 `json:"delivery_fee" binding:"gte=0"`
	Tax         float64 `json:"tax" binding:"gte=0"`
}

type OrderResponse struct {
	ID          int                   `json:"id"`
	CustomerID  int                   `json:"customer_id"`
	OrderDate   time.Time             `json:"order_date"`
	Subtotal    float64               `json:"subtotal"`
	Tax         float64               `json:"tax"`
	DeliveryFee float64               `json:"delivery_fee"`
	Total       float64               `json:"total"`
	Status      string                `json:"status"`
	Items       []OrderDetailResponse `json:"items"`
	CreatedAt   time.Time             `json:"created_at"`
}

type OrderDetailResponse struct {
	ID               int     `json:"id"`
	ProductID        int     `json:"product_id"`
	ProductName      string  `json:"product_name"`
	Quantity         int     `json:"quantity"`
	Price            float64 `json:"price"`
	SizeID           *int    `json:"size_id,omitempty"`
	SizeName         string  `json:"size_name,omitempty"`
	TemperatureID    *int    `json:"temperature_id,omitempty"`
	TemperatureLabel string  `json:"temperature_label,omitempty"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing shipped delivered cancelled"`
}
