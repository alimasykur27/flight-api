package handler

import (
	response_dto "flight-api/internal/dto/response"
	sync_dto "flight-api/internal/dto/sync"
	service_sync "flight-api/internal/service/sync"
	"flight-api/pkg/logger"
	"flight-api/util"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SyncHandler struct {
	service service_sync.ISyncService
	logger  *logger.Logger
}

func NewSyncHandler(service service_sync.ISyncService, logger *logger.Logger) ISyncHandler {
	return &SyncHandler{
		service: service,
		logger:  logger,
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
	// Parse body
	var req sync_dto.SyncAirportRequest
	util.ReadFromRequestBody(r, &req)

	// Call service
	data, err := h.service.SyncAirports(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to sync airports: ", err)
		response := response_dto.ResponseDto{
			Code:   http.StatusInternalServerError,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		util.WriteToResponseBody(w, http.StatusInternalServerError, response)
		return
	}

	h.logger.Debug("Successfully synced airports")

	// Return response
	response := response_dto.ResponseDto{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}
