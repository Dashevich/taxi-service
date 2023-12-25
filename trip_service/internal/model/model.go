package model

import (
	"encoding/json"
)

type DataBase struct {
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Event struct {
	Id              string `json:"id"`
	Source          string `json:"source"`
	Type            string `json:"type"`
	DataContentType string `json:"datacontenttype"`
	Time            string `json:"time"`
	Data            Data   `json:"data"`
}

type Command struct {
	UserId          string          `json:"user"`
	Source          string          `json:"source"`
	Type            string          `json:"type"`
	DataContentType string          `json:"datacontenttype"`
	Time            string          `json:"time"`
	Data            json.RawMessage `json:"data"`
}

type Data interface{}

type EventDataStartEnd struct {
	TripId string `json:"trip_id"`
}
type EventDataAccept struct {
	TripId string `json:"trip_id"`
}
type EventDataCancel struct {
	TripId string `json:"trip_id"`
}
type EventDataCreate struct {
	OfferId string `json:"offer_id"`
	TripId  string `json:"trip_id"`
}

type CommandDataStartEnd struct {
	TripeId string `json:"tripe_id"`
}
type CommandDataAccept struct {
	TripeId  string `json:"tripe_id"`
	DriverId string `json:"driver_id"`
}
type CommandDataCancel struct {
	TripeId string `json:"tripe_id"`
	Reason  string `json:"reason"`
}
type CommandDataCreate struct {
	OfferId string `json:"offer_id"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Price struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}
type Offer struct {
	OfferId  string        `json:"offer_id"`
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
}
