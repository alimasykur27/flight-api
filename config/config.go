package config

import (
	"flight-api/pkg/logger"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName     string        `mapstructure:"SERVICE_NAME"`
	HTTPPort        string        `mapstructure:"HTTP_PORT"`
	LogLevel        string        `mapstructure:"LOG_LEVEL"`
	DatabaseURL     string        `mapstructure:"DATABASE_URL"`
	AviationURL     string        `mapstructure:"AVIATION_API_URL"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

func Load() (config Config, err error) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	// Setup viper to read from .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Handle for repo root
	_, thisFile, _, _ := runtime.Caller(0) // .../config/config.go
	configDir := filepath.Dir(thisFile)    // .../config
	repoRoot := filepath.Clean(filepath.Join(configDir, ".."))
	viper.AddConfigPath(repoRoot)

	// Current dir
	viper.AddConfigPath(".")

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("error reading config file: %w", err)
		}

		logger.Warnf("No .env file found, using environment variables only")
	} else {
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Also read environment variables
	viper.AutomaticEnv()

	// set default values
	viper.SetDefault("SERVICE_NAME", "flight-api")
	viper.SetDefault("HTTP_PORT", "3000")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("AVIATION_API_URL", "")
	viper.SetDefault("SHUTDOWN_TIMEOUT", 5*time.Second)

	err = viper.Unmarshal(&config)
	return config, err
}
