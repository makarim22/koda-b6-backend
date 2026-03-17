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
	authHandler := container.AuthHandler()
	orderHandler := container.OrderHandler()
	cartHandler := container.CartHandler()

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
		auth := api.Group("/auth")
		{
			auth.POST("/forgot-password", forgotPasswordHandler.ForgotPassword)
			auth.POST("/reset-password", forgotPasswordHandler.ResetPassword)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

		}
	}
	{
		orders := api.Group("/orders") // Add orders routes
		{
			orders.GET("", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.POST("", orderHandler.CreateOrder)
			orders.PUT("/:id", orderHandler.UpdateOrderStatus)
			orders.DELETE("/:id", orderHandler.DeleteOrder)
		}
	}
	cartGroup := router.Group("/api/cart")
	{
		cartGroup.GET("", cartHandler.GetCart)
		//cartGroup.GET("/summary", cartHandler.GetCartSummary)
		cartGroup.POST("", cartHandler.AddToCart)
		cartGroup.PUT("/:cart_id", cartHandler.UpdateCartItem)
		cartGroup.DELETE("/:cart_id", cartHandler.RemoveFromCart)
		cartGroup.DELETE("", cartHandler.ClearCart)
	}

}
