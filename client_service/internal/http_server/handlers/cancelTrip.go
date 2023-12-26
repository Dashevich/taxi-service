package handlers

import (
	"client_service/internal/kafka_client"
	"client_service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func (h *Handlers) CancelTrip(w http.ResponseWriter, r *http.Request) {
	store := h.controller
	tripID := chi.URLParam(r, "trip_id")

	offer := models.CommandCancel{
		TripId: tripID,
	}
	data, err := json.Marshal(offer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req := models.Request{
		Id:              tripID,
		Source:          "/client",
		Type:            "trip.command.cancel",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data:            data,
	}
	toTrip, err := kafka_client.ConnectKafka(context.Background(), "kafka:9092", "driver-client-trip-topic", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonRequest, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = kafka_client.SendToTopic(toTrip, jsonRequest)
	if err != nil {
		fmt.Println(err)
	}

	err = store.CancelTrip(tripID)
	if err != nil {
		http.Error(w, "Error canceling", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
