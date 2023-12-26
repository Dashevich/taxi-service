package models

import (
	"encoding/json"
	"google.golang.org/genproto/googleapis/type/decimal"
	"time"
)

type Trip struct {
	ID      string        `bson:"id"`
	OfferID string        `bson:"offer_id"`
	From    LatLngLiteral `bson:"from"`
	To      LatLngLiteral `bson:"to"`
	Price   Money         `bson:"price"`
	Status  string        `bson:"status"`
}

type LatLngLiteral struct {
	Lat decimal.Decimal `bson:"lat"`
	Lng decimal.Decimal `bson:"lng"`
}

type Money struct {
	Amount   decimal.Decimal `bson:"amount"`
	Currency string          `bson:"currency"`
}

type Request struct {
	Id              string          `json:"id"`
	Source          string          `json:"source"`
	Type            string          `json:"type"`
	DataContentType string          `json:"datacontenttype"`
	Time            time.Time       `json:"time"`
	Data            json.RawMessage `json:"data"`
}

type TripRequest struct {
	OfferID string `json:"offer_id"`
}

type Order struct {
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientID string        `json:"client_id"`
	Price    Money         `json:"price"`
}
