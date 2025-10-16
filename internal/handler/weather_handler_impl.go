package handler

import (
	response_dto "flight-api/internal/dto/response"
	service_weather "flight-api/internal/service/weather"
	"flight-api/pkg/logger"
	"flight-api/util"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	service service_weather.IWeatherService
	logger  *logger.Logger
}

func NewWeatherHandler(service service_weather.IWeatherService, logger *logger.Logger) IWeatherHandler {
	return &WeatherHandler{
		service: service,
		logger:  logger,
	}
}

func (h *WeatherHandler) RegisterRouter(r chi.Router) {
	routes := func(r chi.Router) {
		r.Get("/", h.GetWeatherCondition)
	}

	// weather Endpoint
	r.Route("/v1/weathers", routes)
}

func (h *WeatherHandler) GetWeatherCondition(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter
	var loc string
	var response response_dto.ResponseDto
	loc = r.URL.Query().Get("loc")

	if loc == "" {
		response = response_dto.ResponseDto{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Data:    nil,
			Message: "'loc' query parameter is required",
		}
		util.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	}

	// Call service
	data, err := h.service.GetWeatherCondition(r.Context(), &loc)
	if err != nil {
		h.logger.Errorf("[GetWeatherCondition] Failed to get weather condition: %v", err)
		util.ErrorHandler(w, util.NewErrorException(err, "failed to get weather condition: "+err.Error()))
		return
	}

	// Return response
	response = response_dto.ResponseDto{
		Code:    200,
		Status:  "OK",
		Data:    data,
		Message: "Success",
	}

	util.WriteToResponseBody(w, http.StatusOK, response)
}
