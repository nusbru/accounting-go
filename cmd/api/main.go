package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"accounting/internal/handler/http/health"
	"accounting/internal/handler/http/router"
	"accounting/internal/middleware"
	"accounting/internal/pkg/database"
	"accounting/internal/pkg/logger"
	"accounting/internal/repository/postgres"
	"accounting/internal/service"

	_ "accounting/docs" // Import generated docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Accounting API
// @version 1.0
// @description A REST API for personal accounting management with support for users, accounts, and transactions
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@accounting.app

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	// Initialize logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "json"
	}
	log := logger.New(logFormat, logLevel)

	// Initialize database
	dbCfg := database.LoadConfigFromEnv()
	db, err := database.InitDB(dbCfg, log)
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo, userRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)

	// Create router with all handlers
	r := router.NewRouter(userService, accountService, transactionService)

	// Setup HTTP server
	mux := http.NewServeMux()

	// API routes
	mux.Handle("/api/", r)

	// Swagger documentation
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Health check endpoint with database connectivity check
	healthHandler := health.NewHandler(db)
	mux.HandleFunc("/health", healthHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Build middleware chain: RequestID -> Logging -> Recovery -> Handler
	handler := middleware.RequestID(
		middleware.Logging(log)(
			middleware.Recovery(log)(mux),
		),
	)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting API server", "port", port)
		log.Info("Swagger UI available at", "url", "http://localhost:"+port+"/swagger/")
		log.Info("API documentation at", "url", "http://localhost:"+port+"/swagger/doc.json")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
