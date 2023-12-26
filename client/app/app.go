package app

import (
	"client/client/internal/config"
	"client/client/internal/http_server/handlers"
	"client/client/internal/kafka_client"
	"client/client/internal/models"
	"client/client/internal/storage"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"os"

	"context"
	_ "embed"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	//client       *mongo.Client
	store  *storage.ControllerDB
	server *http.Server
	Tracer trace.Tracer
}

func NewApp(cfg *config.Config) *App {
	address := fmt.Sprintf(":%d", cfg.Http.Port)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
		return nil
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		logger.Error("Failed to create:", zap.Error(err))
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Error("Failed to connect MongoDB:", zap.Error(err))
		log.Fatal(err)
	}

	controller := &storage.ControllerDB{
		Cfg:    cfg,
		Client: client,
		Logger: logger,
	}
	tracer := otel.Tracer("final")
	a := &App{
		cfg:    cfg,
		store:  controller,
		logger: logger,
		server: &http.Server{
			Addr:    address,
			Handler: Router(cfg, controller, logger),
		},
		Tracer: tracer,
	}
	return a
}

func Router(cfg *config.Config, controller *storage.ControllerDB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	todoHandlers := handlers.NewTodoHandlers(cfg, logger, controller)
	router.Post("/trips", todoHandlers.CreateTrip)
	router.Get("/trips", todoHandlers.ListTrips)
	router.Get("/trips/{trip_id}", todoHandlers.GetTripByID)
	router.Post("/trip/{trip_id}/cancel", todoHandlers.CancelTrip)

	return router
}

func (a *App) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		err := a.server.ListenAndServe()
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()
	<-done

	go func() {
		for {
			a.Listen()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.CloseDB()

	if a.store.Client != nil {
		if err := a.store.Client.Disconnect(context.Background()); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		}
	}
	a.Stop(ctx)
}

func (a *App) Listen() {
	ctx, span := a.Tracer.Start(context.Background(), "Iteration")
	defer span.End()
	connection, err := kafka_client.ConnectKafka(ctx, "kafka:9092", "trip-client-topic", 0)
	bytes, err := kafka_client.ReadFromTopic(connection)
	if err != nil {
		return
	}
	var req models.Request
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		return
	}

	switch req.Type {
	case "trip.event.accepted":
		var data models.EventAccept
		err := json.Unmarshal(req.Data, &data)
		if err != nil {
			return
		}
		err = a.store.UpdateStatus(data.TripId, "DRIVER_FOUND")
		if err != nil {
			return
		}
	case "trip.event.started":
		var data models.EventStart
		err := json.Unmarshal(req.Data, &data)
		if err != nil {
			return
		}
		err = a.store.UpdateStatus(data.TripId, "STARTED")
		if err != nil {
			return
		}
	case "trip.event.ended":
		var data models.EventEnd
		err := json.Unmarshal(req.Data, &data)
		if err != nil {
			return
		}
		err = a.store.UpdateStatus(data.TripId, "ENDED")
		if err != nil {
			return
		}
	}
}

func (a *App) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	fmt.Println(a.server.Shutdown(ctx))
}

func (a *App) CloseDB() {
	if a.store.Client != nil {
		if err := a.store.Client.Disconnect(context.Background()); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		}
	}
}
