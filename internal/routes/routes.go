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
	reviewsHandler := container.ReviewsHandler()
	productCategoryHandler := container.ProductCategoryHandler()
	orderDetailHandler := container.OrderDetailHandler()
	paymentHandler := container.PaymentHandler()

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
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PUT("/:id", orderHandler.UpdateOrderStatus)
			orders.DELETE("/:id", orderHandler.DeleteOrder)

			orderDetails := orders.Group("/:id/details")
			{
				orderDetails.GET("", orderDetailHandler.GetByOrderID)
				orderDetails.POST("", orderDetailHandler.Create)
				orderDetails.PUT("/:detail_id", orderDetailHandler.Update)
				orderDetails.DELETE("/:detail_id", orderDetailHandler.Delete)
			}
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
	reviewsGroup := api.Group("/reviews")
	{
		reviewsGroup.GET("", reviewsHandler.GetAllReviews)
		reviewsGroup.GET("/:id", reviewsHandler.GetReview)
		reviewsGroup.POST("", reviewsHandler.CreateReview)
		reviewsGroup.PUT("/:id", reviewsHandler.UpdateReview)
	}
	productCategories := api.Group("/product-categories")
	{
		productCategories.GET("", productCategoryHandler.GetAllCategories)
		productCategories.GET("/:id", productCategoryHandler.GetCategoryByID)
		productCategories.POST("", productCategoryHandler.CreateCategory)
		productCategories.PUT("/:id", productCategoryHandler.UpdateCategory)
		productCategories.DELETE("/:id", productCategoryHandler.DeleteCategory)
	}
	{
		payments := api.Group("/payments")
		{
			payments.POST("", paymentHandler.Create)
			payments.GET("/:id", paymentHandler.GetByID)
			payments.PUT("/:id", paymentHandler.UpdateStatus)
			payments.DELETE("/:id", paymentHandler.Delete)
		}
	}
}
