package handler

import (
	airport_dto "flight-api/internal/dto/airport"
	queryparams "flight-api/internal/dto/query_params"
	response_dto "flight-api/internal/dto/response"
	service_airport "flight-api/internal/service/airport"
	"flight-api/pkg/logger"
	"flight-api/util"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AirportHandler struct {
	airportService service_airport.IAirportService
	logger         *logger.Logger
}

// NewAirportHandler
func NewAirportHandler(service service_airport.IAirportService, logger *logger.Logger) IAirportHandler {
	return &AirportHandler{
		airportService: service,
		logger:         logger,
	}
}

// RegisterRoutes
func (h *AirportHandler) RegisterRouter(r chi.Router) {
	routes := func(r chi.Router) {
		// Create Airport Data
		r.Post("/", h.Create)
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	}

	// Airports Endpoints
	r.Route("/v1/airports", routes)
}

// Create Airport Data
func (h *AirportHandler) Create(w http.ResponseWriter, r *http.Request) {
	airportReq := airport_dto.AirportRequestDto{}
	util.ReadFromRequestBody(r, &airportReq)

	airportResponse := h.airportService.Create(r.Context(), airportReq)
	response := response_dto.ResponseDto{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   airportResponse,
	}

	util.WriteToResponseBody(w, http.StatusCreated, response)
}

// Find All data
func (h *AirportHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	query := queryparams.GetQueryParams(r)

	airportResponses := h.airportService.FindAll(r.Context(), query)
	response := response_dto.ResponseDto{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   airportResponses,
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}

// Find By ID
func (h *AirportHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	airportResponse, err := h.airportService.FindByID(r.Context(), id)

	if err != nil {
		if err == util.ErrNotFound {
			response := response_dto.ResponseDto{
				Code:   http.StatusNotFound,
				Status: "Not Found",
				Data:   nil,
			}
			util.WriteToResponseBody(w, http.StatusNotFound, response)

			return
		}

		response := response_dto.ResponseDto{
			Code:   http.StatusInternalServerError,
			Status: "Internal Server Error",
			Data:   nil,
		}
		util.WriteToResponseBody(w, http.StatusInternalServerError, response)

		return
	}

	response := response_dto.ResponseDto{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   airportResponse,
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}

// Update
func (h *AirportHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	airportUpdate := airport_dto.AirportUpdateDto{}
	util.ReadFromRequestBody(r, &airportUpdate)

	airportResponse, err := h.airportService.Update(r.Context(), id, airportUpdate)

	if err != nil {
		switch {
		case err == util.ErrBadRequest:
			response := response_dto.ResponseDto{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   nil,
			}
			util.WriteToResponseBody(w, http.StatusBadRequest, response)

			return
		case err == util.ErrNotFound:
			response := response_dto.ResponseDto{
				Code:   http.StatusNotFound,
				Status: "Not Found",
				Data:   nil,
			}
			util.WriteToResponseBody(w, http.StatusNotFound, response)

			return
		default:
			response := response_dto.ResponseDto{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Data:   nil,
			}
			util.WriteToResponseBody(w, http.StatusInternalServerError, response)

			return
		}
	}

	response := response_dto.ResponseDto{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   airportResponse,
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}

// Delete
func (h *AirportHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.airportService.Delete(r.Context(), id)

	if err != nil {
		switch {
		case err == util.ErrNotFound:
			response := response_dto.ResponseDto{
				Code:    http.StatusNotFound,
				Status:  "Not Found",
				Data:    nil,
				Message: fmt.Sprintf("Airport with ID %s not found", id),
			}
			util.WriteToResponseBody(w, http.StatusNotFound, response)

			return
		default:
			response := response_dto.ResponseDto{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Data:    nil,
				Message: err.Error(),
			}
			util.WriteToResponseBody(w, http.StatusInternalServerError, response)

			return
		}
	}

	h.logger.Debugf("Airport with ID %s deleted successfully", id)

	response := response_dto.ResponseDto{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: fmt.Sprintf("Airport with ID %s deleted successfully", id),
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}
