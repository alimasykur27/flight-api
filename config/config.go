package config

import (
	"flight-api/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName     string        `mapstructure:"SERVICE_NAME"`
	HTTPPort        string        `mapstructure:"HTTP_PORT"`
	LogLevel        string        `mapstructure:"LOG_LEVEL"`
	AppEnv          string        `mapstructure:"APP_ENV"`
	DatabaseURL     string        `mapstructure:"DATABASE_URL"`
	RedisURL        string        `mapstructure:"REDIS_URL"`
	RedisEnable     bool          `mapstructure:"REDIS_ENABLE"`
	AviationURL     string        `mapstructure:"AVIATION_API_URL"`
	WeatherURL      string        `mapstructure:"WEATHER_API_URL"`
	WeatherAPIKey   string        `mapstructure:"WEATHER_API_KEY"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

func Load() (config Config, err error) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	if path := os.Getenv("CONFIG_FILE"); path != "" {
		viper.SetConfigFile(path)
	} else {
		// Setup viper to read from .env file
		viper.SetConfigName(".env")
		viper.SetConfigType("env")

		// repo root & current dir
		_, thisFile, _, _ := runtime.Caller(0) // .../config/config.go
		configDir := filepath.Dir(thisFile)    // .../config
		repoRoot := filepath.Clean(filepath.Join(configDir, ".."))
		viper.AddConfigPath(repoRoot)

		// Current dir
		viper.AddConfigPath(".")
	}

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
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("REDIS_URL", "")
	viper.SetDefault("REDIS_ENABLE", false)
	viper.SetDefault("AVIATION_API_URL", "")
	viper.SetDefault("WEATHER_API_URL", "")
	viper.SetDefault("WEATHER_API_KEY", "")
	viper.SetDefault("SHUTDOWN_TIMEOUT", 5*time.Second)

	err = viper.Unmarshal(&config)
	return config, err
}
