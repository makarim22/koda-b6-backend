package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"strings"
	"regexp"
	"golang.org/x/crypto/bcrypt" 
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Phone string
	Gender string
	Age int
	Address string
}

type Product struct {
	ID int
	Name string
	Description string
	Stock int
	VariantId int `json:"variant_id"`
	SizeId int `json:"size_id"`
}

type LoginPayload struct {
	Email string
	Password string
}

// type RegisterPayload struct {
// 	Fullname string
// 	Email string
// 	Password string
// 	Phone string
// 	Gender string
// 	Age int
// 	Address string
// }

var users = map[int]User{
	1: {ID: 1, Name: "Budi", Email: "budi@email.com", Password: "hashed123"},
	2: {ID: 2, Name: "Siti", Email: "siti@email.com", Password: "hashed456"},
}

var products = map[int]Product{}

var nextID = 3

var userEmails = map[string]int{}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

var nextProductId = 1


func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain", []byte("Hello!"))
	})

	r.GET("/user", func(ctx *gin.Context) {
       
		if len (users) == 0 {
			ctx.JSON(404, gin.H{
				"error": "No users found",
			})
			return
		}

		userList := make([]User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}

		ctx.JSON(200, gin.H{
			"data":  userList,
			"count": len(userList),
		})
	})

	r.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		user , exists := users[userId]
		if !exists {
          ctx.JSON(400, gin.H{
			"error" : "tidak dapat menemukan user",
		  })
		  return
		}

		
	    ctx.JSON(200, user)

	})

	r.POST("/users", func(ctx *gin.Context) {
		var newUser User
		fmt.Println("newUser", newUser)

		err := ctx.ShouldBindJSON(&newUser)

		if err != nil {
			ctx.JSON(400, gin.H{
				"Success": false,
				"Message": "Gagal membuat user",
			})
			return
		} 

		name := strings.TrimSpace(newUser.Name)
		email := strings.ToLower(strings.TrimSpace(newUser.Email))
		password := strings.TrimSpace(newUser.Password)

        if name == "" || email == "" || password == "" {
            ctx.JSON(400, gin.H{
				"success": false,
				"message": "Nama, email, dan password tidak boleh kosong",
				"error":   "validation_error",
			})
			return
		}

		for _ , existingUser := range users {
                if existingUser.Email == email {
				ctx.JSON(409, gin.H{ 
					"success": false,
					"message": "Email sudah terdaftar",
					"error":   "duplicate_email",
				})
				return
			}
		}

		newUser.ID = nextID
		nextID++
		users[newUser.ID]= newUser

		ctx.JSON(201, gin.H{ 
			"success": true,
			"message": "Berhasil membuat user",
			"data":    newUser, 
		})

	})

	r.PUT("/user/:id", func(ctx *gin.Context){ 
        id := ctx.Param("id")

		userID, err := strconv.Atoi(id)
	    fmt.Println("usernya", userID)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		user, exists := users[userID]
		if !exists {
			ctx.JSON(404, gin.H{
				"error": "User not found",
			})
			return
		}

		var updateData User
		if err := ctx.ShouldBindJSON(&updateData); err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		name := strings.TrimSpace(updateData.Name)
		email := strings.ToLower(strings.TrimSpace(updateData.Email))

		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)

		if err != nil {
			return
		}
		
		password := string(HashedPassword)


		fmt.Println("name", name)
		fmt.Println("email", email)
		fmt.Println("password", password)

		user.Name = name
		user.Email = email
		user.Password = password

		users[userID] = user

		ctx.JSON(200, gin.H{
			"message": "berhasil mengupdate user",
			"data":    user,
		})
	})

	/// endpoint auth 

	r.POST("/register", func(ctx *gin.Context) {
		var newUser User

		if err := ctx.ShouldBindJSON(&newUser); err != nil {
			ctx.JSON(400, gin.H{
				"success": false, 
				"message": "Invalid request body", 
				"error": err.Error()})
			return
		}

		Name := strings.TrimSpace(newUser.Name)
		Email := strings.ToLower(strings.TrimSpace(newUser.Email))

		if Name == "" || Email == "" || newUser.Password == "" {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Nama, email, dan password tidak boleh kosong",
				"error": "validation_error"})
			return
		}

		if !emailRegex.MatchString(Email) {
			ctx.JSON(400, gin.H{
				"success": false, 
				"message": "Format email tidak valid", 
				"error": "invalid_email_format"})
			return
		}

		if _, emailExists := userEmails[newUser.Email]; emailExists {
			ctx.JSON(400, gin.H{
				"success": false,
				 "message": "Email sudah terdaftar", 
				 "error": "duplicate_email"})
			return
		}

		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

		if err != nil {
			return
		}
		newUser.Password = string(HashedPassword)

		newUser.ID = nextID
		nextID++
		users[newUser.ID] = newUser
		userEmails[newUser.Email] = newUser.ID

		responseUser := newUser
		responseUser.Password = "" // hilangkan password
		ctx.JSON(200, gin.H{
			"success": true, 
			"message": "Registrasi berhasil", 
			"data": responseUser})
	})

	r.POST("/login", func(ctx *gin.Context) {
		var payload LoginPayload

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{
				"success": false, 
				"message": "Invalid request body", 
				"error": err.Error()})
			return
		}

		Email := strings.ToLower(strings.TrimSpace(payload.Email))

		userID, emailExists := userEmails[Email]
		if !emailExists {
			ctx.JSON(403, gin.H{
				"success": false,
				 "message": "Email atau password salah", 
				 "error": "invalid_credentials"})
			return
		}

		user := users[userID]
		fmt.Println("user", user)

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
		if err != nil {
			ctx.JSON(403, gin.H{
				"success": false, 
				"message": "Email atau password salah", 
				"error": "invalid_credentials"})
			return
		}
		responseUser := user
		responseUser.Password = ""
		ctx.JSON(200, gin.H{
			"success": true, 
			"message": "Login berhasil", 
			"data": responseUser})
	})

	// endpoint product
	r.POST("/products", func(ctx *gin.Context) {
		var newProduct Product

		if err := ctx.ShouldBindJSON(&newProduct); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Invalid request body",
				"error":   err.Error()})
			return
		}

		newProduct.ID = nextProductId
		products[newProduct.ID] = newProduct 

		nextProductId++ 

		fmt.Printf("Added product: %+v\n", newProduct)
		fmt.Println("Current products map:", products) 

		ctx.JSON(201, gin.H{ 
			"success": true,
			"message": "Product created successfully",
			"data":    newProduct,
		})
	})

	r.GET("/products", func(ctx *gin.Context) {
		if len(products) == 0 {
			ctx.JSON(404, gin.H{
				"error": "No products found",
			})
			return
		}

		productList := make([]Product, 0, len(products))
		for _, product := range products {
			productList = append(productList, product)
		}

		ctx.JSON(200, gin.H{
			"data":  productList,
			"count": len(productList),
		})
	})



	r.Run("localhost:8888")
}