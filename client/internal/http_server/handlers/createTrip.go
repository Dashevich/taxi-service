package handlers

import (
	"client/internal/kafka_client"
	"client/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"time"
)

func (h *Handlers) CreateTrip(w http.ResponseWriter, r *http.Request) {
	store := h.controller
	var model models.Trip

	err := json.NewDecoder(r.Body).Decode(&model)
	fmt.Println(r)
	fmt.Println(r.Header.Get("user_id"))
	fmt.Println(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := http.Get("http://offering:8080/offers/" + model.OfferID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var order models.Order
	err = json.Unmarshal(bytes, &order)
	fmt.Println(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	// Validate and insert trip into MongoDB
	offer := models.Trip{
		OfferID: model.OfferID,
		ID:      r.Header.Get("user_id"),
		From: models.LatLngLiteral{
			Lat: order.From.Lat,
			Lng: order.From.Lng,
		},
		To: models.LatLngLiteral{
			Lat: order.To.Lat,
			Lng: order.To.Lng,
		},
		Price: models.Money{
			Amount:   order.Price.Amount,
			Currency: order.Price.Currency,
		},
		Status: "DRIVER_SEARCH",
	}
	result, err := store.CreateTrip(&offer)
	offerData := models.CommandCreate{
		OfferId: model.OfferID,
	}
	jsonData, err := json.Marshal(offerData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req := models.Request{
		Id:              result.InsertedID.(primitive.ObjectID).Hex(),
		Source:          "/client",
		Type:            "trip.command.create",
		DataContentType: "application/jsonapplication/json",
		Time:            time.Now(),
		Data:            jsonData,
	}

	toTrip, err := kafka_client.ConnectKafka(context.Background(), "kafka:9092", "driver-client-trip-topic", 0)
	jsonRequest, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = kafka_client.SendToTopic(toTrip, jsonRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
