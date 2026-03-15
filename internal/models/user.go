package models

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Phone    string `json:"phone" db:"phone"`
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
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
