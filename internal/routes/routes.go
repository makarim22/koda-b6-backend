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

	api := router.Group("/api/v1")
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
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
		}
	}
}