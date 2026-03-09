package models

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Phone    string `json:"phone" db:"phone"`
}

type LoginPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type RegisterPayload struct {
	Fullname string `json:"full_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Phone string `json:"phone"`
	Gender string `json:"gender"`
	Age int `json:"age"`
	Address string `json:"address"`
}

// var users = map[int]User{
// 	1: {ID: 1, Name: "Budi", Email: "budi@email.com", Password: "hashed123"},
// 	2: {ID: 2, Name: "Siti", Email: "siti@email.com", Password: "hashed456"},
// }

// var nextID = 3

// var userEmails = map[string]int{}

// var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)