package routes

import (
	// "fmt"
	"koda-b6-backend/internal/di"

	"github.com/gin-gonic/gin"
	// "github.com/jackc/pgx/v5"
)

func SetupRoutes(router *gin.Engine, container *di.Container) {
	userHandler := container.UserHandler()
	productHandler := container.ProductHandler()
	forgotPasswordHandler := container.ForgotPasswordHandler()

	api := router.Group("/admin")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
	{
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts)
			products.GET("/:id", productHandler.GetById)
			products.GET("/recommended-products", productHandler.MostReviewedProduct)
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
		}
	}
	{
		forgotPasswords := api.Group("/auth")
		{
			forgotPasswords.POST("/forgot-password", forgotPasswordHandler.ResetPassword)

		}
	}
}
