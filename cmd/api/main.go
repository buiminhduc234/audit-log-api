package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/buiminhduc234/audit-log-api/internal/api"
	"github.com/buiminhduc234/audit-log-api/internal/config"
	"github.com/buiminhduc234/audit-log-api/internal/middleware"
	"github.com/buiminhduc234/audit-log-api/internal/repository/postgres"
	"github.com/buiminhduc234/audit-log-api/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize config
	cfg := &config.Config{
		ServerPort:         10000,
		JWTSecretKey:       os.Getenv("JWT_SECRET_KEY"),
		JWTExpirationHours: 24,
	}

	// Initialize database
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	repo := postgres.NewPostgresRepository(db)

	// Initialize services
	tenantService := service.NewTenantService(repo)
	auditLogService := service.NewAuditLogService(repo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware()

	// Initialize server
	server := api.NewServer(
		tenantService,
		auditLogService,
		authMiddleware,
		rateLimitMiddleware,
	)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup API routes
	api := router.Group("/api/v1")
	server.SetupRoutes(api)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
