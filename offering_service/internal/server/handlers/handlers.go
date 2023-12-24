package handlers

import (
	"offering_service/internal/config"
	"offering_service/internal/service"
)

type Handler struct {
	cfg             *config.Config
	offeringService service.OfferingService
	jwtService      service.JWTService
}

func NewHandler(config *config.Config, offer *service.OfferingService, jwt *service.JWTService) *Handler {
	return &Handler{cfg: config,
		offeringService: *offer,
		jwtService:      *jwt}
}
