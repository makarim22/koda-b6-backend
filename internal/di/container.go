package di

import (
	"fmt"
	"koda-b6-backend/internal/handlers"
	"koda-b6-backend/internal/repository"
	"koda-b6-backend/internal/service"

	"github.com/jackc/pgx/v5"
)

type Container struct {
	db *pgx.Conn

	// user
	userRepo    *repository.UserRepository
	userService *service.UserService
	userHandler *handlers.UserHandler
	authHandler *handlers.AuthHandler
	authService *service.AuthService

	// product
	productRepo    *repository.ProductRepository
	productService *service.ProductService
	productHandler *handlers.ProductHandler

	//forgotPassword
	forgotPasswordRepo    *repository.ForgotPasswordRepository
	forgotPasswordService *service.ForgotPasswordService
	forgotPasswordHandler *handlers.ForgotPasswordHandler

	// order
	orderRepo    *repository.OrderRepository
	orderService *service.OrderService
	orderHandler *handlers.OrderHandler
}

func NewContainer(db *pgx.Conn) (*Container, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}

	container := &Container{
		db: db,
	}

	container.initDependencies()

	return container, nil
}

func (c *Container) initDependencies() {
	//Users
	c.userRepo = repository.NewUserRepository(c.db)
	c.userService = service.NewUserService(c.userRepo)
	c.userHandler = handlers.NewUserHandler(c.userService)
	c.authService = service.NewAuthService(c.userRepo)
	c.authHandler = handlers.NewAuthHandler(c.userService, c.authService)

	//Products
	c.productRepo = repository.NewProductRepository(c.db)
	c.productService = service.NewProductService(c.productRepo)
	c.productHandler = handlers.NewProductHandler(c.productService)

	//forgotPassword
	c.forgotPasswordRepo = repository.NewForgotPasswordRepository(c.db)
	c.forgotPasswordService = service.NewForgotPasswordService(c.userRepo, c.forgotPasswordRepo)
	c.forgotPasswordHandler = handlers.NewForgotPasswordHandler(c.forgotPasswordService)

	//order
	c.orderRepo = repository.NewOrderRepository(c.db)
	c.orderService = service.NewOrderService(c.orderRepo, c.productRepo)
	c.orderHandler = handlers.NewOrderHandler(c.orderService)

}

func (c *Container) UserHandler() *handlers.UserHandler {
	return c.userHandler
}

func (c *Container) ProductHandler() *handlers.ProductHandler {
	return c.productHandler
}

func (c *Container) ForgotPasswordHandler() *handlers.ForgotPasswordHandler {
	return c.forgotPasswordHandler
}

func (c *Container) AuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

func (c *Container) OrderHandler() *handlers.OrderHandler {
	return c.orderHandler
}
