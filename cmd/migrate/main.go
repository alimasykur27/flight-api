package main

import (
	"flag"
	"flight-api/config"
	"flight-api/pkg/logger"
	migrations "flight-api/pkg/migration"
)

func main() {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	// Define command-line flags for migrations
	upFlag := flag.Bool("up", false, "Run migrations up")
	downFlag := flag.Bool("down", false, "Roll migrations down")
	statusFlag := flag.Bool("status", false, "Show migrations status")
	helpFlag := flag.Bool("help", false, "Show help information")
	flag.Parse()

	// Show help information if requested
	if *helpFlag {
		logger.Info("Please specify a migration command: --up, --down, or --status")
		return
	}

	// Load application connfiguration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration %w", err)
	}

	// Migrations handler
	migration := migrations.NewMigrations(cfg, logger)

	// Execute migration command
	if *downFlag {
		migration.Rollback()
	} else if *statusFlag {
		migration.Status()
	} else if *upFlag {
		migration.Run()
	} else {
		// Default behavior - show help
		logger.Info("Please specify a migration command: --up, --down, or --status")
		logger.Info("Use --help for more information")
	}
}
