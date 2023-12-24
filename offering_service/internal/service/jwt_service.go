package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"offering_service/internal/config"
	"offering_service/internal/model"
)

type JWTService struct {
	SecretKey string
}

func NewJWTService(cfg *config.Config) JWTService {
	return JWTService{cfg.SigningKey}
}

func (service *JWTService) CreateToken(data model.Offer) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)

	claims["client_id"] = data.ClientId
	claims["from_lat"] = data.From.Lat
	claims["from_lng"] = data.From.Lng
	claims["to_lat"] = data.To.Lat
	claims["to_lng"] = data.To.Lng
	claims["price_amount"] = data.Price.Amount
	claims["price_currency"] = data.Price.Currency

	tokenString, err := token.SignedString(service.SecretKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
func (service *JWTService) ExtractClaims(offerId string) (model.Offer, error) {
	token, _ := jwt.Parse(offerId, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return service.SecretKey, nil
	})

	if token == nil {
		return model.Offer{}, nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.Offer{}, nil
	}
	var response = model.Offer{
		ClientId: claims["client_id"].(string),
		From:     model.LatLngLiteral{Lat: claims["from_lat"].(float64), Lng: claims["from_lng"].(float64)},
		To:       model.LatLngLiteral{Lat: claims["to_lat"].(float64), Lng: claims["to_lng"].(float64)},
		Price:    model.Price{Amount: claims["price_amount"].(float64), Currency: claims["price_currency"].(string)},
	}
	return response, nil
}
