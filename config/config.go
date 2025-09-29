package config

import (
	"flight-api/pkg/logger"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func Load() (config Config, err error) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	// Setup viper to read from .env file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("error reading config file: %w", err)
		}

		logger.Info("No .env file found, using environment variables only")
	} else {
		logger.Info("Using config file: %s", viper.ConfigFileUsed())
	}

	// Also read environment variables
	viper.AutomaticEnv()

	// set default values
	viper.SetDefault("SERVICE_NAME", "flight-api")
	viper.SetDefault("HTTP_PORT", "3000")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DATABASE_URL", "")

	err = viper.Unmarshal(&config)

	return config, err
}
