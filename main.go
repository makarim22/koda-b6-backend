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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func corsMiddleware() gin.HandlerFunc {
	godotenv.Load()
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
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

	// serverPort := os.Getenv("SERVER_PORT")
	// if serverPort == "" {
	// 	serverPort = ":3002"
	// }

	serverPort := os.Getenv("SERVER_PORT")
    if serverPort == "" {
    serverPort = ":3002"
    } else if !strings.HasPrefix(serverPort, ":") {
    serverPort = ":" + serverPort
    }

	ctx := context.Background()

	// Use pgxpool instead of single connection
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to parse database config: %v", err)
	}

	// Configure pool for concurrent requests
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("❌ Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Could not ping database: %v", err)
	}

	log.Println("✅ Successfully connected to the database!")

	container, err := di.NewContainer(pool)
	if err != nil {
		log.Fatalf("❌ Failed to initialize container: %v", err)
	}

	router := gin.Default()

	router.Use(corsMiddleware())

	router.Static("/uploads", "./uploads")

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
