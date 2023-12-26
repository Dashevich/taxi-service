package server

import (
	"context"
	"fmt"
	"net/http"
	"offering_service/internal/config"
	"offering_service/internal/server/handlers"
	services "offering_service/internal/service"
)

type Server struct {
	cfg             *config.Config
	offeringService services.OfferingService
	jwtService      services.JWTService
	server          *http.Server
}

func NewServer(cfg *config.Config) *Server {
	addr := fmt.Sprintf(":%d", cfg.Http.Port)
	offering := services.NewSimpleOfferingService()
	jwt := services.NewJWTService(cfg)
	return &Server{
		cfg:             cfg,
		offeringService: &offering,
		jwtService:      jwt,
		server: &http.Server{
			Addr:    addr,
			Handler: initApi(cfg, &offering, jwt),
		},
	}
}

func (s *Server) Run() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	fmt.Println(s.server.Shutdown(ctx))
}

func initApi(cfg *config.Config, offer services.OfferingService, jwt services.JWTService) http.Handler {
	mux := http.NewServeMux()
	handler := handlers.NewHandler(cfg, &offer, &jwt)
	mux.HandleFunc("/offers/", handler.CreateOffer)
	mux.HandleFunc("/offers/{offer_id}", handler.ParseOffer)
	return mux
}
