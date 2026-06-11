package models


import "time"

type User struct {
	ID            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"full_name"`
	Email         string  `json:"email" db:"email"`
	Password      string  `json:"password" db:"password"`
	Role          string  `db:"role"` 
	Phone         *string `json:"phone" db:"phone"`
	Address       *string `json:"address" db:"address"`
	ProfileImage  *string `json:"profile_image" db:"profile_image"`
	PointsBalance int     `json:"points_balance" db:"points_balance"`
}

type PointLedger struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	OrderID     *int      `json:"order_id" db:"order_id"`
	Points      int       `json:"points" db:"points"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}


type UserRole struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterPayload struct {
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
	Role string  `json:"role"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type UserResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Phone         *string `json:"phone"`
	Address       *string `json:"address"`
	ProfileImage  *string `json:"profile_image"`
	PointsBalance int     `json:"points_balance"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)