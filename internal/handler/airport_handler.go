package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IAirportHandler interface {
	RegisterRouter(r chi.Router)
	Create(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetWeatherCondition(w http.ResponseWriter, r *http.Request)
}
