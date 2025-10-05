package migrations

import (
	"flight-api/config"
	"flight-api/pkg/database"
	"flight-api/pkg/logger"

	"github.com/sirupsen/logrus"
)

var (
	GetMigrationStatus = database.GetMigrationStatus
	runMigrations      = database.RunMigrations
	rollbackMigrations = database.RollbackMigrations
)

type Migrations struct {
	cfg    config.Config
	logger *logger.Logger
}

func NewMigrations(cfg config.Config, logger *logger.Logger) *Migrations {
	return &Migrations{
		cfg:    cfg,
		logger: logger,
	}
}

func (m *Migrations) GetLogger() *logger.Logger {
	return m.logger
}

func (m *Migrations) HandleMigrations(down bool, status bool) {
	m.logger.Info("Starting migration process...")

	if status {
		m.logger.Info("Checking migration status")
		version, dirty, err := GetMigrationStatus(m.cfg.DatabaseURL, "migrations")

		if err != nil {
			m.logger.Fatalw(logrus.Fields{"error": err}, "Failed to get migration status")
		}

		if version == 0 {
			m.logger.Info("No migrations have been applied yet")
		} else {
			m.logger.Infow(logrus.Fields{
				"version": version,
				"dirty":   dirty,
				"state":   m.getDirtyStateMessage(dirty),
			}, "Current migration status")
		}

		return
	}

	if down {
		m.logger.Info("Rolling back migrations")
		if err := rollbackMigrations(m.cfg.DatabaseURL, "migrations"); err != nil {
			m.logger.Fatalw(logrus.Fields{
				"error": err,
			}, "Failed to rollback migration")
		}
		m.logger.Info("Successfully rolled back migrations")
	} else {
		m.logger.Info("Running migrations")
		if err := runMigrations(m.cfg.DatabaseURL, "migrations"); err != nil {
			m.logger.Fatalw(logrus.Fields{
				"error": err,
			}, "Failed to run migrations")
		}
		m.logger.Info("Successfully ran migrations")
	}
}

func (m *Migrations) getDirtyStateMessage(dirty bool) string {
	if dirty {
		return "Migration is in a dirty state and may need manual intervention"
	}
	return "Clean - migrations have been applied successfully"
}

// Run executes database migrations
func (m *Migrations) Run() {
	m.HandleMigrations(false, false)
}

// Rollback reverts database migrations
func (m *Migrations) Rollback() {
	m.HandleMigrations(true, false)
}

// Status shows the current migration status
func (m *Migrations) Status() {
	m.HandleMigrations(false, true)
}
