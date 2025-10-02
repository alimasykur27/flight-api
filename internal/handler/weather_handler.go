package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IWeatherHandler interface {
	RegisterRouter(r chi.Router)
	GetWeatherCondition(w http.ResponseWriter, r *http.Request)
}
