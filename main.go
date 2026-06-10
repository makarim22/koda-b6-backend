package main

import (
	"context"
	"errors"
	"koda-b6-backend/internal/di"
	"koda-b6-backend/internal/lib"
	"koda-b6-backend/internal/routes"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}

// slogMiddleware integrates slog into Gin
func slogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.Info("HTTP Request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.Duration("duration", duration),
		)
	}
}

var db *pgx.Conn

func main() {
	// 1. Initialize slog (JSON handler)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if err := lib.InitConfig(); err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}
	slog.Info("Loaded database URL", "url_length", len(databaseURL))

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
		slog.Error("Failed to parse database config", "error", err)
		os.Exit(1)
	}

	// Configure pool for concurrent requests
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Could not ping database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully connected to the database!")

	container, err := di.NewContainer(pool)
	if err != nil {
		slog.Error("Failed to initialize DI container", "error", err)
		os.Exit(1)
	}

	// Disable Gin's default console color output
	gin.SetMode(gin.ReleaseMode)
	router := gin.New() // Create engine without default middlewares

	// Use custom middlewares
	router.Use(gin.Recovery())
	router.Use(slogMiddleware(logger))
	router.Use(corsMiddleware())

	router.Static("/uploads", "./uploads")

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	routes.SetupRoutes(router, container)

	// Configure robust HTTP Server
	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Server is starting...", "port", serverPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown signal received, shutting down gracefully...")

	// The context is used to inform the server it has 5 seconds to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting gracefully")
}
