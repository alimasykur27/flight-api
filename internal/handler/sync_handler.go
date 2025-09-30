package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ISyncHandler interface {
	RegisterRouter(r chi.Router)
	SyncAirport(w http.ResponseWriter, r *http.Request)
}
