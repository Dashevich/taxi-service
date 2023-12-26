package app

import (
	"context"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"net/http"
	"time"
	"trip_service/internal/config"
	"trip_service/internal/database"
	"trip_service/internal/model"
	"trip_service/internal/service"
)

type App struct {
	cfg *config.Config
	db  *sqlx.DB

	ClientWriter *kafka.Conn
	DriverWriter *kafka.Conn
	Reader       *kafka.Conn
}

func NewApp(cfg *config.Config, ctx context.Context) *App {
	db, err := database.InitDB(ctx, cfg)
	if err != nil {
		log.Fatal("no connection to db", err)
	}
	time.Sleep(10)
	conn1, err := service.ConnectKafka(ctx, cfg.Kafka.Address, "trip-to-client", 0)
	if err != nil {
		log.Fatal("Kafka connect error")
	}
	conn2, err := service.ConnectKafka(ctx, cfg.Kafka.Address, "trip-to-driver", 0)
	if err != nil {
		log.Fatal("Kafka connect error")
	}
	conn3, err := service.ConnectKafka(ctx, cfg.Kafka.Address, "trip-inbound", 0)
	if err != nil {
		log.Fatal("Kafka connect error")
	}
	return &App{
		cfg:          cfg,
		db:           db,
		ClientWriter: conn1,
		DriverWriter: conn2,
		Reader:       conn3,
	}
}

func (app *App) GetOffer(offerId string) (*model.Offer, error) {
	resp, err := http.Get("http://localhost:" + string(app.cfg.Http.Port) + "/offers/" + offerId)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var order model.Offer
	err = json.Unmarshal(bytes, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (app *App) Run(ctx context.Context) {
	for {
		bytes, err := service.ReadFromTopic(app.Reader)
		if err != nil {
			return
		}
		var request model.Command
		err = json.Unmarshal(bytes, &request)
		if err != nil {
			return
		}
		if request.DataContentType != "application/json" {
			return
		}
		response := model.Event{
			Id:              request.UserId,
			Source:          "/trip",
			Type:            "",
			DataContentType: "application/json",
			Time:            request.Time,
		}
		switch request.Type {
		case "trip.command.accept":
			response.Type = "trip.event.accepted"
			var commandData model.CommandDataAccept
			err := json.Unmarshal(bytes, &commandData)
			if err != nil {
				return
			}
			err = app.ChangeRow(commandData.TripeId, ctx, "DRIVER_FOUND")
			if err != nil {
				log.Fatal("db fall")
			}
			response.Data = model.EventDataAccept{TripId: commandData.TripeId}
			res, err := json.Marshal(response)
			if err != nil {
				log.Fatal("data error")
			}
			err = service.SendToTopic(app.ClientWriter, res)
		case "trip.command.cancel":
			response.Type = "trip.event.canceled"
			var commandData model.CommandDataCancel
			err = json.Unmarshal(bytes, &commandData)
			if err != nil {
				return
			}
			err = app.ChangeRow(commandData.TripeId, ctx, "CANCELED")
			if err != nil {
				log.Fatal("db fall")
			}
			response.Data = model.EventDataCancel{TripId: commandData.TripeId}
			res, err := json.Marshal(response)
			if err != nil {
				log.Fatal("AcceptTrip mistake")
			}
			err = service.SendToTopic(app.ClientWriter, res)
			if err != nil {
				log.Fatal("kafka send error")
			}
			err = service.SendToTopic(app.DriverWriter, res)
			if err != nil {
				log.Fatal("kafka send error")
			}

		case "trip.command.start":
			response.Type = "trip.event.started"
			var commandData model.CommandDataStartEnd
			err = json.Unmarshal(bytes, &commandData)
			if err != nil {
				return
			}
			err = app.ChangeRow(commandData.TripeId, ctx, "STARTED")
			if err != nil {
				log.Fatal("db fall")
			}
			response.Data = model.EventDataStartEnd{TripId: commandData.TripeId}
			res, err := json.Marshal(response)
			if err != nil {
				log.Fatal("AcceptTrip mistake")
			}
			err = service.SendToTopic(app.ClientWriter, res)
			if err != nil {
				log.Fatal("kafka send error")
			}
		case "trip.command.end":
			response.Type = "trip.event.ended"
			var commandData model.CommandDataStartEnd
			err = json.Unmarshal(bytes, &commandData)
			if err != nil {
				return
			}
			err = app.ChangeRow(commandData.TripeId, ctx, "ENDED")
			if err != nil {
				log.Fatal("db fall")
			}
			response.Data = model.EventDataStartEnd{TripId: commandData.TripeId}
			res, err := json.Marshal(response)
			if err != nil {
				log.Fatal("AcceptTrip mistake")
			}
			err = service.SendToTopic(app.ClientWriter, res)
			if err != nil {
				log.Fatal("kafka send error")
			}
		case "trip.command.create":
			response.Type = "trip.event.created"
			var commandData model.CommandDataCreate
			err = json.Unmarshal(bytes, &commandData)
			if err != nil {
				return
			}
			offer, err := app.GetOffer(commandData.OfferId)
			if err != nil {
				log.Fatal("wrong token for offering service")
			}
			event, err := app.CreateRow(commandData.OfferId, *offer, ctx)
			if err != nil {
				log.Fatal("db adding row error")
			}
			event.Time = response.Time
			res, err := json.Marshal(event)
			if err != nil {
				log.Fatal("AcceptTrip mistake")
			}
			err = service.SendToTopic(app.ClientWriter, res)
			if err != nil {
				log.Fatal("kafka send error")
			}
		}
		select {
		case <-ctx.Done():
			break
		default:
		}

	}

}

func (app *App) CreateRow(offer_id string, offer model.Offer, ctx context.Context) (model.Event, error) {
	from, _ := json.Marshal(offer.From)
	to, _ := json.Marshal(offer.To)
	price, _ := json.Marshal(offer.Price)
	trip_id := uuid.New().String()
	sql, args, err := squirrel.Insert("trips").
		Columns("trip_id", "offer_id", "from_offer", "to_offer", "price", "trip_status").
		Values(trip_id, offer_id, from, to, price, "CREATED").
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Event{}, err
	}
	var id int
	new_row := app.db.QueryRowContext(ctx, sql, args...)
	if err = new_row.Scan(&id); err != nil {
		return model.Event{}, err
	}
	event := model.Event{
		Id:     uuid.New().String(),
		Source: "trip",
		Type:   "trip.event.created",
		Data: model.EventDataCreate{
			TripId:  trip_id,
			OfferId: offer_id,
		},
		DataContentType: "application/json",
	}
	return event, nil
}

func (app *App) ChangeRow(tripId string, ctx context.Context, status string) error {
	sql, args, err := squirrel.Update("trips").
		Set("trip_status", status).
		Where(squirrel.Eq{"trip_id": tripId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	var id int
	new_row := app.db.QueryRowContext(ctx, sql, args...)
	if err = new_row.Scan(&id); err != nil {
		return err
	}
	return nil

}
