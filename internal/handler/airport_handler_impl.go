package handler

import (
	airport_dto "flight-api/internal/dto/airport"
	queryparams "flight-api/internal/dto/query_params"
	response_dto "flight-api/internal/dto/response"
	service_airport "flight-api/internal/service/airport"
	"flight-api/pkg/logger"
	"flight-api/util"
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
		r.Get("/weathers", h.GetWeatherCondition)
	}

	// Airports Endpoints
	r.Route("/v1/airports", routes)
}

// Create Airport Data
func (h *AirportHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	airportReq := airport_dto.AirportRequestDto{}
	err := util.ReadFromRequestBody(r, &airportReq)

	if err != nil {
		h.logger.Errorf("[Create] Failed to read request body: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "invalid request body: "+err.Error()))
		return
	}

	// Call service
	airportResponse, err := h.airportService.Create(r.Context(), airportReq)
	if err != nil {
		h.logger.Errorf("[Create] Failed to create airport: %v", err.Error())
		util.ErrorHandler(w, util.NewErrorException(err, "failed to create airport: "+err.Error()))
		return
	}

	// Response (201 Created)
	response := response_dto.ResponseDto{
		Code:    http.StatusCreated,
		Status:  "Created",
		Data:    airportResponse,
		Message: "airport created successfully",
	}

	util.WriteToResponseBody(w, http.StatusCreated, response)
}

// Find All data
func (h *AirportHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	query := queryparams.GetQueryParams(r)

	airportResponses, err := h.airportService.FindAll(r.Context(), query)
	if err != nil {
		h.logger.Errorf("[FindAll] Failed to fetch airports: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "failed to fetch airports"))
		return
	}

	// Response (200 OK)
	response := response_dto.ResponseDto{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    airportResponses,
		Message: "airports fetched successfully",
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}

// Find By ID
func (h *AirportHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	airportResponse, err := h.airportService.FindByID(r.Context(), id)

	if err != nil {
		h.logger.Errorf("[FindByID] Failed to fetch airport by ID: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "failed to fetch airport by ID: "+err.Error()))
		return
	}

	// Response (200 OK)
	response := response_dto.ResponseDto{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    airportResponse,
		Message: "airport fetched successfully",
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
		h.logger.Errorf("[Update] Failed to update airport: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "failed to update airport"))
		return
	}

	h.logger.Debugf("Airport with ID %s updated successfully", id)

	// Response (200 OK)
	response := response_dto.ResponseDto{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    airportResponse,
		Message: "airport updated successfully",
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}

// Delete
func (h *AirportHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.airportService.Delete(r.Context(), id)

	if err != nil {
		h.logger.Errorf("[Delete] Failed to delete airport: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "failed to delete airport: "+err.Error()))
		return
	}

	h.logger.Debugf("Airport with ID %s deleted successfully", id)
	response := response_dto.ResponseDto{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "airport deleted successfully",
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}

func (h *AirportHandler) GetWeatherCondition(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter optional
	var code, name string
	var response response_dto.ResponseDto

	query := queryparams.GetQueryParams(r)
	code = r.URL.Query().Get("code")
	name = r.URL.Query().Get("name")

	if code == "" && name == "" {
		// Should has value
		response = response_dto.ResponseDto{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Data:    nil,
			Message: "neither 'code' or 'name' query parameter is required",
		}

		util.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	} else if code != "" && name != "" {
		// can't has value at the same time
		response = response_dto.ResponseDto{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Data:    nil,
			Message: "'code' and 'name' query parameter can't be used at the same time",
		}

		util.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	}

	// Call service
	data, err := h.airportService.GetWeatherCondition(r.Context(), code, name, query)

	if err != nil {
		h.logger.Errorf("[GetWeatherCondition] Failed to get weather condition: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "failed to get weather condition: "+err.Error()))
		return
	}
	h.logger.Debugf("Get weather condition success")

	// Response (200 OK)
	response = response_dto.ResponseDto{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    data,
		Message: "weather condition fetched successfully",
	}
	util.WriteToResponseBody(w, http.StatusOK, response)
}
