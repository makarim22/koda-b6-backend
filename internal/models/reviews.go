package models

import "time"

type Reviews struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	ProductId int       `json:"product_id" db:"product_id"`
	OrderId   int       `json:"order_id" db:"order_id"`
	Message   string    `json:"message" db:"message"`
	Rating    int       `json:"rating" db:"rating"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
type ReviewsResponse struct {
	Id          int    `json:"id" db:"id"`
	UserId      int    `json:"user_id" db:"user_id"`
	UserName    string `json:"user_name" db:"full_name"`
	Email       string `json:"email" db:"email"`
	ProductId   int    `json:"product_id" db:"product_id"`
	ProductName string `json:"product_name" db:"product_name"`
	OrderId     int    `json:"order_id" db:"order_id"`
	Message     string `json:"message" db:"message"`
	Rating      int    `json:"rating" db:"rating"`
}

type ReviewsRequest struct {
	Id        int    `json:"id" db:"id"`
	UserId    int    `json:"user_id" db:"user_id"`
	ProductId int    `json:"product_id" db:"product_id"`
	OrderId   int    `json:"order_id" db:"order_id"`
	Message   string `json:"message" db:"message"`
	Rating    int    `json:"rating" db:"rating"`
}

type UpdateReviewsRequest struct {
	Message string `json:"message" db:"message"`
	Rating  int    `json:"rating" db:"rating"`
}
