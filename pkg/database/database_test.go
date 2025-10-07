package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestConnect_Success_UsesPoolSettings(t *testing.T) {
	// Arrange: create a fake *sql.DB with sqlmock
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orig := sqlOpen
	defer func() { sqlOpen = orig }()

	capturedDriver := ""
	capturedDSN := ""
	sqlOpen = func(driverName, dsn string) (*sql.DB, error) {
		capturedDriver = driverName
		capturedDSN = dsn
		// Return our sqlmock DB regardless of args
		return db, nil
	}

	// Act
	got, err := Connect("host=localhost user=me password=secret dbname=test sslmode=disable")

	// Assert
	assert.NoError(t, err)
	assert.Same(t, db, got, "should return the DB from the opener")
	assert.Equal(t, "postgres", capturedDriver)
	assert.Contains(t, capturedDSN, "host=localhost")

	// Can only reliably assert MaxOpenConnections via Stats
	stats := got.Stats()
	assert.Equal(t, 25, stats.MaxOpenConnections, "MaxOpenConns should be set to 25")

	// (We won't assert other pool values since database/sql doesn't expose getters)
}

func TestConnect_OpenError_Wrapped(t *testing.T) {
	orig := sqlOpen
	defer func() { sqlOpen = orig }()

	sqlOpen = func(driverName, dsn string) (*sql.DB, error) {
		return nil, errors.New("boom")
	}

	got, err := Connect("any")
	assert.Nil(t, got)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "failed to connect to database")
		assert.Contains(t, err.Error(), "boom") // wrapped original error
	}
}

// --- Test Get Migration Status ---
