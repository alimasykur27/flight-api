package handler

import (
	response_dto "flight-api/internal/dto/response"
	"flight-api/pkg/logger"
	"flight-api/util"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SyncHandler struct {
	logger *logger.Logger
}

func NewSyncHandler(logger *logger.Logger) ISyncHandler {
	return &SyncHandler{
		logger: logger,
	}
}

// RegisterRoutes
func (h *SyncHandler) RegisterRouter(r chi.Router) {
	routes := func(r chi.Router) {
		// Sync Airport Data
		r.Post("/airports", h.SyncAirport)
	}

	// Sync Endpoints
	r.Route("/v1/sync", routes)
}

func (h *SyncHandler) SyncAirport(w http.ResponseWriter, r *http.Request) {
	response := response_dto.ResponseDto{
		Code:   200,
		Status: "OK",
		Data:   "adalah",
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}
