package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handlers) ListTrips(w http.ResponseWriter, r *http.Request) {
	store := h.controller
	user_id := r.Header.Get("user_id")
	trips, err := store.ListTrips(user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	answer, err := json.Marshal(trips)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(answer)
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handlers) GetTripByID(w http.ResponseWriter, r *http.Request) {
	store := h.controller
	tripID := chi.URLParam(r, "trip_id")

	trip, err := store.GetTripByID(tripID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	answer, err := json.Marshal(trip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(answer)
	w.WriteHeader(http.StatusOK)
	return
}
