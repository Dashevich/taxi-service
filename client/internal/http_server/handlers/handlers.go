package handlers

import (
	"client/client/internal/config"
	"client/client/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type Handlers struct {
	Server     *http.Server
	cfg        *config.Config
	controller *storage.ControllerDB
}

func NewTodoHandlers(config *config.Config, logger *zap.Logger, controller *storage.ControllerDB) *Handlers {
	return &Handlers{
		cfg:        config,
		controller: controller,
	}
}
