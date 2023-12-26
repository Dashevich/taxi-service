package storage

import (
	"client/client/internal/config"
	"client/client/internal/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type ControllerDB struct {
	Cfg    *config.Config
	Client *mongo.Client
	Logger *zap.Logger
}

type Storage interface {
	ListTrips() ([]models.Trip, error)
	CreateTrip(trip *models.Trip) error
	GetTripByID(tripID string) (*models.Trip, error)
	CancelTrip(tripID string) error
	UpdateStatus(tripID string, status string) error
}

func (c *ControllerDB) ListTrips(user_id string) ([]models.Trip, error) {
	collection := c.Client.Database("my_mongodb").Collection("trips")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"client_id": user_id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trips []models.Trip
	for cursor.Next(ctx) {
		var trip models.Trip
		err := cursor.Decode(&trip)
		if err != nil {
			c.Logger.Warn("Decoding error")
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

func (c *ControllerDB) CreateTrip(trip *models.Trip) (*mongo.InsertOneResult, error) {
	collection := c.Client.Database("my_mongodb").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, trip)
	if err != nil {
		return &mongo.InsertOneResult{}, err
	}
	if err != nil {
		c.Logger.Warn("Insert error")
	}
	return res, nil
}

func (c *ControllerDB) GetTripByID(tripID string) (*models.Trip, error) {
	collection := c.Client.Database("my_mongodb").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, err
	}
	var trip models.Trip
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&trip)
	if err != nil {
		return nil, err
	}

	return &trip, nil
}

func (c *ControllerDB) CancelTrip(tripID string) error {
	collection := c.Client.Database("my_mongodb").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return err
	}
	res, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %v document(s)\n", res.DeletedCount)
	return nil
}

func (c *ControllerDB) UpdateStatus(tripID string, status string) error {
	collection := c.Client.Database("my_mongo").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": status}}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.Logger.Warn("Failed to update")
		return err
	}

	if res.ModifiedCount == 0 {
		c.Logger.Warn("Trip not found or not authorized")
		return err
	}
	return nil
}
