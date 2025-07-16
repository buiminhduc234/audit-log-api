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
	"github.com/buiminhduc234/audit-log-api/internal/repository/composite"
	"github.com/buiminhduc234/audit-log-api/internal/service"
	"github.com/buiminhduc234/audit-log-api/internal/service/queue"
	"github.com/buiminhduc234/audit-log-api/pkg/logger"
)

// @title           Audit log Swagger API
// @version         1.0
// @description     This is a Audit log swagger server.

// @host      localhost:10000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize logger
	appLogger := logger.NewLogger(os.Getenv("APP_ENV"))

	// Initialize config
	cfg := &config.Config{
		ServerPort:         10000,
		JWTSecretKey:       os.Getenv("JWT_SECRET_KEY"),
		JWTExpirationHours: 24,
	}

	// Initialize database
	db, err := config.NewDatabase()
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}

	// Initialize OpenSearch
	osConfig := config.DefaultOpenSearchConfig()
	osClient, err := osConfig.GetClient()
	if err != nil {
		appLogger.Fatal("Failed to connect to OpenSearch", err)
	}

	// Initialize SQS
	sqsConfig := config.DefaultSQSConfig()
	sqsClient, err := sqsConfig.GetClient()
	if err != nil {
		appLogger.Fatal("Failed to connect to SQS", err)
	}
	sqsService := queue.NewSQSService(sqsClient, sqsConfig)

	// Initialize repositories
	repo := composite.NewCompositeRepository(db, osClient, osConfig)

	// Initialize services
	tenantService := service.NewTenantService(repo)
	auditLogService := service.NewAuditLogService(repo, sqsService)

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
			appLogger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	// Shutdown the HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", err)
	}

	appLogger.Info("Server exiting")
	appLogger.Sync()
}
