package util

import (
	"errors"
	"flight-api/pkg/logger"

	"github.com/sirupsen/logrus"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrNotFound = errors.New("record not found")
var ErrBadRequest = errors.New("bad request")
var ErrGatewayTimeout = errors.New("gateway timeout")
var ErrInternalServer = errors.New("internal server error")

// LogPanicError will log the error and panic if the error is not nil
func LogPanicError(err error) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	logger.Errorw(logrus.Fields{"error": err}, "[FLIGHT API ERROR]")
}

// PanicIfError will panic if the error is not nil
func PanicIfError(err error) {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

	if err != nil {
		logger.Errorw(logrus.Fields{"error": err}, "[FLIGHT API ERROR]")
		panic(err)
	}
}

func RecoverPanic(err *error) {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case string:
			*err = errors.New(x)
		case error:
			*err = x
		default:
			*err = errors.New("unknown panic")
		}
	}
}
