package service

import (
	"math"
	"offering_service/internal/model"
)

type OfferingService interface {
	GetOffer(request model.Request) model.Offer
}

type SimpleOfferingService struct {
	OneDistPrice float64
	Currency     string
}

func NewSimpleOfferingService() SimpleOfferingService {
	return SimpleOfferingService{OneDistPrice: 10, Currency: "RUB"}
}

func (service *SimpleOfferingService) GetOffer(request model.Request) model.Offer {
	tmp := math.Pi / 180.0
	var radius float64 = 6400
	from_lat := request.From.Lat * tmp
	from_lng := request.From.Lng * tmp
	to_lat := request.To.Lat * tmp
	to_lng := request.To.Lng * tmp
	d := math.Acos(math.Sin(from_lat)*math.Sin(to_lat) + math.Cos(from_lat)*math.Cos(to_lat)*math.Cos(from_lng-to_lng))
	cost := service.OneDistPrice * (1 + d*radius)

	return model.Offer{
		ClientId: request.ClientId,
		Price:    model.Price{Amount: cost, Currency: service.Currency},
		From:     request.From,
		To:       request.To,
	}
}
