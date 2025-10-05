// config/config_test.go
package config_test

import (
	"flight-api/config"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad_WithEnvFile(t *testing.T) {
	viper.Reset()

	dir := t.TempDir()
	dotenv := filepath.Join(dir, ".env")
	err := os.WriteFile(dotenv, []byte(
		"SERVICE_NAME=from-file\nHTTP_PORT=7777\nSHUTDOWN_TIMEOUT=11s\n",
	), 0o644)
	assert.NoError(t, err)

	// arahkan ke file yang barusan dibuat
	t.Setenv("CONFIG_FILE", dotenv)

	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "from-file", cfg.ServiceName)
	assert.Equal(t, "7777", cfg.HTTPPort)
	assert.Equal(t, 11*time.Second, cfg.ShutdownTimeout)

}

func TestLoad_NoEnvFile(t *testing.T) {
	viper.Reset()
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.Equal(t, "flight-api", cfg.ServiceName) // default
	assert.Equal(t, "3000", cfg.HTTPPort)          // default
}

func TestLoad_ValidConfig(t *testing.T) {
	t.Helper()
	viper.Reset()
	cfg, err := config.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Service Name
	assert.Equal(t, "flight-api", cfg.ServiceName)

	// Is Port int
	port, err := strconv.Atoi(cfg.HTTPPort)
	assert.NoError(t, err)
	assert.IsType(t, 0, port)

	// LogLevel
	levels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	assert.Contains(t, levels, cfg.LogLevel)

	// App Env
	envs := []string{"development", "production", "test"}
	assert.Contains(t, envs, cfg.AppEnv)

	// DatabaseUrl
	assert.NotEmpty(t, cfg.DatabaseURL)
	assert.IsType(t, "", cfg.DatabaseURL)
	assert.Contains(t, cfg.DatabaseURL, "postgres://")
	assert.Contains(t, cfg.DatabaseURL, "@")
	assert.Contains(t, cfg.DatabaseURL, ":")
	assert.Contains(t, cfg.DatabaseURL, "/")
	assert.Contains(t, cfg.DatabaseURL, "?")

	// AviationUrl
	assert.NotEmpty(t, cfg.AviationURL)
	assert.IsType(t, "", cfg.AviationURL)
	assert.Contains(t, cfg.AviationURL, "https://")

	// WeatherUrl
	assert.NotEmpty(t, cfg.WeatherURL)
	assert.IsType(t, "", cfg.WeatherAPIKey)
	assert.Contains(t, cfg.WeatherURL, "https://")

	// WeatherAPIKey
	assert.NotEmpty(t, cfg.WeatherAPIKey)
	assert.IsType(t, "", cfg.WeatherAPIKey)

	// ShutdownTimeout
	assert.NotEmpty(t, cfg.ShutdownTimeout)
	assert.IsType(t, time.Duration(0), cfg.ShutdownTimeout)

}
