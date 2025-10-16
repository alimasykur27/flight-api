package middleware

// Custom middleware to recover from panics and return a 500 error
// and the message to the client.

import (
	response_dto "flight-api/internal/dto/response"
	"flight-api/util"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				response := response_dto.ResponseDto{
					Code:    http.StatusInternalServerError,
					Status:  "Internal Server Error",
					Data:    nil,
					Message: "A server error occurred",
					Errors:  []interface{}{"internal server error", rvr},
				}
				util.WriteToResponseBody(w, http.StatusInternalServerError, response)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
