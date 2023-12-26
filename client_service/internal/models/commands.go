package models

type CommandAccept struct {
	TripId   string `json:"trip_id"`
	DriverId string `json:"driver_id"`
}

type CommandStart struct {
	TripId string `json:"trip_id"`
}

type CommandCreate struct {
	OfferId string `json:"offer_id"`
}

type CommandCancel struct {
	TripId string `json:"trip_id"`
	Reason string `json:"reason"`
}

type CommandEnd struct {
	TripId string `json:"trip_id"`
}
