package di

import (
	"koda-b6-backend/internal/handlers"
	"koda-b6-backend/internal/repository"
	"koda-b6-backend/internal/service"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Container struct {
	db *pgx.Conn

	// user
	userRepo *repository.UserRepository
	userService *service.UserService
	userHandler *handlers.UserHandler
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

func  (c *Container) initDependencies(){
//Users
	c.userRepo = repository.NewUserRepository(c.db)
	c.userService = service.NewUserService(c.userRepo)
	c.userHandler = handlers.NewUserHandler(c.userService)

}

func (c *Container) UserHandler() *handlers.UserHandler {
	return c.userHandler
}