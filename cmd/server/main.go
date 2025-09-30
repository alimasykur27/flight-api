package main

import (
	"context"
	"flight-api/config"
	"flight-api/internal/handler"
	"flight-api/internal/repository"
	"flight-api/internal/service"
	"flight-api/pkg/database"
	"flight-api/pkg/httpserver"
	"flight-api/pkg/logger"
	"flight-api/util"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	// Create a context that will be canceled on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create a logger
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	logger.Info("Connecting to database ...")
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalw(logrus.Fields{
			"error": err,
		}, "Failed to connect to database")
	} else {
		logger.Info("Succesfully connected to database!")
	}
	defer db.Close()

	// Initialize validator
	validate := util.NewValidator()

	// Initialize repository
	airportRepository := repository.NewAirportRepository(logger)

	// Initialize service
	airportService := service.NewAirportService(airportRepository, db, validate, logger)

	// Initialize Handlers
	airportHandler := handler.NewAirportHandler(airportService, logger)
	syncHandler := handler.NewSyncHandler(logger)

	// Setup router
	logger.Info("Setup Router ...")
	router := httpserver.NewRouter(
		airportHandler,
		syncHandler,
	)

	// Start HTTP server
	server := httpserver.NewServer(router, cfg.HTTPPort)

	serverErrCh := make(chan error, 1)
	go func() {
		logger.Infow(logrus.Fields{"port": cfg.HTTPPort}, "Starting HTTP Server")
		serverErrCh <- server.Start()
	}()

	select {
	case err := <-serverErrCh:
		logger.Errorw(logrus.Fields{"error": err}, "Server error")
	case <-ctx.Done():
		logger.Info("Received termination signal")
	}

	// Graceful shutdown
	logger.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorw(logrus.Fields{"error": err}, "Error during server shutdown")
	}

	logger.Info("Server stopped gracefully!")
}
