package util

import (
	"errors"
	"flight-api/pkg/logger"

	"github.com/sirupsen/logrus"
)

var ErrBadRequest = errors.New("bad request")                 // 400
var ErrUnauthorized = errors.New("unauthorized")              // 401
var ErrPaymentRequired = errors.New("payment required")       // 402
var ErrForbidden = errors.New("forbidden")                    // 403
var ErrNotFound = errors.New("record not found")              // 404
var ErrConflict = errors.New("data conflict")                 // 409
var ErrInternalServer = errors.New("internal server error")   // 500
var ErrNotImplemented = errors.New("not implemented")         // 501
var ErrBadGateway = errors.New("bad gateway")                 // 502
var ErrServiceUnavailable = errors.New("service unavailable") // 503
var ErrGatewayTimeout = errors.New("gateway timeout")         // 504

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
