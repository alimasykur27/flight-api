package database_test

import (
	"flight-api/config"
	"flight-api/pkg/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Ping()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}
