package database

import (
	"context"
	"fpgschiba.com/automation-meal/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var eventsCollection *mongo.Collection

func GetEventsCollection(client *mongo.Client) *mongo.Collection {
	if eventsCollection != nil {
		return eventsCollection
	}
	eventsCollection = client.Database(DatabaseName).Collection("events")
	return eventsCollection
}

func insertEvent(event models.Event) error {
	client = getClient()
	eventsCollection := GetEventsCollection(client)

	_, err := eventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"event":     event,
			"component": "database",
			"func":      "InsertEvent",
		}).Error()
		return err
	}
	return nil
}

func InsertAuthEvent(eventName string, userId string) error {
	event := models.Event{
		Name:      eventName,
		UserId:    userId,
		Component: models.ComponentAuth,
		TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
	}
	return insertEvent(event)
}
