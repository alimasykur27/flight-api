package httpserver

import (
	"net/http"
	"time"

	"flight-api/pkg/logger"
	mid "flight-api/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	RegisterRouter(r chi.Router)
}

func NewRouter(handlers ...Handler) *chi.Mux {
	logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	r := chi.NewRouter()

	// Middleware
	logger.Info("Setup Middleware ...")
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mid.HTTPLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Credentials", "true"))

	// Handle OPTIONS requests
	r.Options("/v1/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Health check endpoint
	r.Get("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logger.Errorw(logrus.Fields{"error": err}, "[NewRouter] Error write")
		}
	})

	// Register all handlers
	for _, h := range handlers {
		h.RegisterRouter(r)
	}

	return r
}
