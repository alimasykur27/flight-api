package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var sqlOpen = sql.Open

// Connect establishes a connection to the database
func Connect(connectionString string) (*sql.DB, error) {
	db, err := sqlOpen("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db, nil
}

// Run database migrations
func RunMigrations(connectionString string, migrationsPath string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create migrations directory if doesn't exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsPath, 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %w", err)
		}
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver %w", err)
	}

	// Get absolute path for migrations
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute migrations path %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations %w", err)
	}

	return nil
}

// Rollback database migrations
func RollbackMigrations(connectionString string, migrationsPath string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get absolute path for migrations
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

// Get current migrations version and dirty state
func GetMigrationStatus(connectionString string, migrationsPath string) (uint, bool, error) {
	db, err := sqlOpen("postgres", connectionString)
	if err != nil {
		return 0, false, err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get absolute path for migrations
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return 0, false, fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			// No migrations have been applied yet
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}
