package util

import (
	"errors"
	response_dto "flight-api/internal/dto/response"
	"flight-api/pkg/logger"
	"fmt"
	"net/http"

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

func ErrorHandler(w http.ResponseWriter, err error) {
	switch err {
	case ErrBadRequest:
		response := response_dto.ResponseDto{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Data:    nil,
			Message: fmt.Sprintf("Invalid request data: %v", err),
		}
		WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	case ErrUnauthorized:
		response := response_dto.ResponseDto{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Data:    nil,
			Message: "Authentication is required and has failed or has not yet been provided",
		}
		WriteToResponseBody(w, http.StatusUnauthorized, response)
		return
	case ErrPaymentRequired:
		response := response_dto.ResponseDto{
			Code:    http.StatusPaymentRequired,
			Status:  "Payment Required",
			Data:    nil,
			Message: "Payment is required to access the requested resource",
		}
		WriteToResponseBody(w, http.StatusPaymentRequired, response)
		return
	case ErrForbidden:
		response := response_dto.ResponseDto{
			Code:    http.StatusForbidden,
			Status:  "Forbidden",
			Data:    nil,
			Message: "You do not have permission to access the requested resource",
		}
		WriteToResponseBody(w, http.StatusForbidden, response)
		return
	case ErrNotFound:
		response := response_dto.ResponseDto{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Data:    nil,
			Message: "The requested resource could not be found",
		}
		WriteToResponseBody(w, http.StatusNotFound, response)
		return
	case ErrConflict:
		response := response_dto.ResponseDto{
			Code:    http.StatusConflict,
			Status:  "Conflict",
			Data:    nil,
			Message: "Airport with the same ICAO ID already exists",
		}
		WriteToResponseBody(w, http.StatusConflict, response)
		return
	case ErrInternalServer:
		response := response_dto.ResponseDto{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Data:    nil,
			Message: "An unexpected error occurred on the server",
		}
		WriteToResponseBody(w, http.StatusInternalServerError, response)
		return
	case ErrNotImplemented:
		response := response_dto.ResponseDto{
			Code:    http.StatusNotImplemented,
			Status:  "Not Implemented",
			Data:    nil,
			Message: "The requested functionality is not implemented",
		}
		WriteToResponseBody(w, http.StatusNotImplemented, response)
		return
	case ErrBadGateway:
		response := response_dto.ResponseDto{
			Code:    http.StatusBadGateway,
			Status:  "Bad Gateway",
			Data:    nil,
			Message: "Received an invalid response from the upstream server",
		}
		WriteToResponseBody(w, http.StatusBadGateway, response)
		return
	case ErrServiceUnavailable:
		response := response_dto.ResponseDto{
			Code:    http.StatusServiceUnavailable,
			Status:  "Service Unavailable",
			Data:    nil,
			Message: "The server is currently unable to handle the request",
		}
		WriteToResponseBody(w, http.StatusServiceUnavailable, response)
		return
	case ErrGatewayTimeout:
		response := response_dto.ResponseDto{
			Code:    http.StatusGatewayTimeout,
			Status:  "Gateway Timeout",
			Data:    nil,
			Message: "The server did not receive a timely response from the upstream server",
		}
		WriteToResponseBody(w, http.StatusGatewayTimeout, response)
		return
	default:
		response := response_dto.ResponseDto{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Data:    nil,
			Message: fmt.Sprintf("An unexpected error occurred: %v", err),
		}
		WriteToResponseBody(w, http.StatusInternalServerError, response)
		return
	}
}
