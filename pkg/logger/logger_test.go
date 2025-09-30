package logger_test

import (
	"flight-api/pkg/logger"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	if logger == nil {
		t.Error("Expected logger to be initialized, got nil")
	}

	if logger.GetLogger().Level != logrus.InfoLevel {
		t.Errorf("Expected log level to be InfoLevel, got %v", logger.GetLogger().Level)
	}

	logger.Info("This is an info message")
	logger.Debug("This is a debug message - should not appear in Info level")
}
