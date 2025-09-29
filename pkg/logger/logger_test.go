package logger_test

import (
	"flight-api/pkg/logger"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	logger.Infow(logrus.Fields{
		"id": "123",
	}, "halo %s", "sjdfhj")
}
