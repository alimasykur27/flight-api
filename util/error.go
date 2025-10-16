package util

import (
	"errors"
	response_dto "flight-api/internal/dto/response"
	"net/http"
)

// Standard Error Variables
var ErrValidation = errors.New("validation error")
var ErrNilPointer = errors.New("nil pointer dereference")
var ErrInvalidInput = errors.New("invalid input")
var ErrDatabase = errors.New("database error")
var ErrTimeout = errors.New("operation timed out")
var ErrUnknown = errors.New("unknown error")

// HTTP Standard Error Variables
var ErrBadRequest = errors.New("bad request")                   // HTTP Status: 400
var ErrUnauthorized = errors.New("unauthorized")                // HTTP Status: 401
var ErrPaymentRequired = errors.New("payment required")         // HTTP Status: 402
var ErrForbidden = errors.New("forbidden")                      // HTTP Status: 403
var ErrNotFound = errors.New("record not found")                // HTTP Status: 404
var ErrConflict = errors.New("data conflict")                   // HTTP Status: 409
var ErrUnprocessableEntity = errors.New("unprocessable entity") // HTTP Status: 422
var ErrInternalServer = errors.New("internal server error")     // HTTP Status: 500
var ErrNotImplemented = errors.New("not implemented")           // HTTP Status: 501
var ErrBadGateway = errors.New("bad gateway")                   // HTTP Status: 502
var ErrServiceUnavailable = errors.New("service unavailable")   // HTTP Status: 503
var ErrGatewayTimeout = errors.New("gateway timeout")           // HTTP Status: 504

type ErrorException struct {
	Cause   error
	Message string
	Meta    *map[string]any
}

func (e *ErrorException) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause != nil && e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func (e *ErrorException) Unwrap() error { return e.Cause }

func NewErrorException(cause error, msg string) error {
	return &ErrorException{Cause: cause, Message: msg, Meta: nil}
}

func mapErrorStatus(err error) (int, string) {
	switch {
	case errors.Is(err, ErrValidation),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, "Bad Request"
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "Unauthorized"
	case errors.Is(err, ErrPaymentRequired):
		return http.StatusPaymentRequired, "Payment Required"
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "Forbidden"
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "Not Found"
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, "Conflict"
	case errors.Is(err, ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity, "Unprocessable Entity"
	case errors.Is(err, ErrTimeout), errors.Is(err, ErrGatewayTimeout):
		return http.StatusGatewayTimeout, "Gateway Timeout"
	case errors.Is(err, ErrServiceUnavailable), errors.Is(err, ErrDatabase):
		return http.StatusServiceUnavailable, "Service Unavailable"
	case errors.Is(err, ErrNotImplemented):
		return http.StatusNotImplemented, "Not Implemented"
	case errors.Is(err, ErrBadGateway):
		return http.StatusBadGateway, "Bad Gateway"
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

func firstNonEmpty(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func collectErrors(err error) []interface{} {
	out := []interface{}{}

	// Collect all error messages in the chain
	for e := err; e != nil; e = errors.Unwrap(e) {
		out = append(out, e.Error())
	}

	// Reverse the slice to maintain the original error order
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return out
}

func ErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	code, status := mapErrorStatus(err)

	var ee *ErrorException
	if errors.As(err, &ee) {
		response := response_dto.ResponseDto{
			Code:    code,
			Status:  status,
			Data:    nil,
			Message: firstNonEmpty(ee.Message, status),
			Errors:  collectErrors(err),
		}
		WriteToResponseBody(w, code, response)
		return
	}

	// Fallback for generic errors
	response := response_dto.ResponseDto{
		Code:    code,
		Status:  status,
		Data:    nil,
		Message: err.Error(),
		Errors:  []any{err.Error()},
	}
	WriteToResponseBody(w, code, response)
}
