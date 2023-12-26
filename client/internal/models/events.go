package models

type EventAccept struct {
	TripId string `json:"trip_id"`
}

type EventCancel struct {
	TripId string `json:"trip_id"`
}

type EventEnd struct {
	TripId string `json:"trip_id"`
}

type EventStart struct {
	TripId string `json:"trip_id"`
}
