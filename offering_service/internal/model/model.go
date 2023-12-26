package model

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Request struct {
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientId string        `json:"client_id"`
}

type Response struct {
	OfferId  string        `json:"offer_id"`
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
}

type Offer struct {
	ClientId string        `json:"client_id"`
	Price    Price         `json:"price"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
}
