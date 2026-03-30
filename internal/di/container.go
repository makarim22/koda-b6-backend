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

	//cart
	cartRepo    *repository.CartRepository
	cartService *service.CartService
	cartHandler *handlers.CartHandler

	//reviews
	reviewsRepo    *repository.ReviewsRepository
	reviewsService *service.ReviewsService
	reviewsHandler *handlers.ReviewsHandler

	//product_category
	productCategoryRepo    *repository.ProductCategoryRepository
	productCategoryService *service.ProductCategoryService
	productCategoryHandler *handlers.ProductCategoryHandler

	//order_detail
	orderDetailRepo    *repository.OrderDetailRepository
	orderDetailService *service.OrderDetailService
	orderDetailHandler *handlers.OrderDetailHandler

	//payment
	paymentRepo    *repository.PaymentRepository
	paymentService *service.PaymentService
	paymentHandler *handlers.PaymentHandler

	//variant
	variantRepo    *repository.VariantRepository
	variantService *service.VariantService
	variantHandler *handlers.VariantHandler
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

	//cart
	c.cartRepo = repository.NewCartRepository(c.db)
	c.cartService = service.NewCartService(c.cartRepo, c.productRepo)
	c.cartHandler = handlers.NewCartHandler(c.cartService)

	//reviews
	c.reviewsRepo = repository.NewReviewsRepository(c.db)
	c.reviewsService = service.NewReviewsService(c.reviewsRepo, c.orderRepo, c.productRepo)
	c.reviewsHandler = handlers.NewReviewsHandler(c.reviewsService)

	//product_category
	c.productCategoryRepo = repository.NewProductCategoryRepository(c.db)
	c.productCategoryService = service.NewProductCategoryService(c.productCategoryRepo)
	c.productCategoryHandler = handlers.NewProductCategoryHandler(c.productCategoryService)

	//order_detail
	c.orderDetailRepo = repository.NewOrderDetailRepository(c.db)
	c.orderDetailService = service.NewOrderDetailService(c.orderDetailRepo, c.orderRepo, c.productRepo)
	c.orderDetailHandler = handlers.NewOrderDetailHandler(c.orderDetailService)

	//payment
	c.paymentRepo = repository.NewPaymentRepository(c.db)
	c.paymentService = service.NewPaymentService(c.paymentRepo, c.orderRepo)
	c.paymentHandler = handlers.NewPaymentHandler(c.paymentService)

	//variant
	c.variantRepo = repository.NewVariantRepository(c.db)
	c.variantService = service.NewVariantService(c.variantRepo)
	c.variantHandler = handlers.NewVariantHandler(c.variantService)

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

func (c *Container) CartHandler() *handlers.CartHandler {
	return c.cartHandler
}

func (c *Container) ReviewsHandler() *handlers.ReviewsHandler {
	return c.reviewsHandler
}

func (c *Container) ProductCategoryHandler() *handlers.ProductCategoryHandler {
	return c.productCategoryHandler
}

func (c *Container) OrderDetailHandler() *handlers.OrderDetailHandler { return c.orderDetailHandler }

func (c *Container) PaymentHandler() *handlers.PaymentHandler { return c.paymentHandler }

func (c *Container) VariantHandler() *handlers.VariantHandler { return c.variantHandler }
