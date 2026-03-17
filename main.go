package main

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/di"
	"koda-b6-backend/internal/lib"
	"koda-b6-backend/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func corsMiddleware() gin.HandlerFunc {
	godotenv.Load()
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.GetHeader("Content-Type")
		if ctx.Request.Method == http.MethodOptions {
			ctx.Data(http.StatusOK, "", []byte(""))
		} else {
			ctx.Next()
		}
	}
}

var db *pgx.Conn

func main() {

	if err := lib.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	fmt.Println("📌 Database URL:", databaseURL)

	if databaseURL == "" {
		log.Fatal("❌ DATABASE_URL environment variable is not set")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8888"
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		log.Fatalf("❌ Could not ping database: %v", err)
	}

	log.Println("✅ Successfully connected to the database!")
	defer conn.Close(ctx)

	container, err := di.NewContainer(conn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize container: %v", err)
	}

	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	routes.SetupRoutes(router, container)

	log.Printf("🚀 Server started on http://localhost%s\n", serverPort)
	if err := router.Run(serverPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
