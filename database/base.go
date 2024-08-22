package database

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var client *mongo.Client

const (
	DatabaseName = "HomeAutomation"
)

func getClient() *mongo.Client {
	uri := "mongodb://localhost:27017"
	//getting context
	if client != nil {
		return client
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//getting client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"uri":       uri,
			"component": "database",
			"func":      "getClient",
		}).Error()
	}
	return client
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if client == nil {
		return
	}
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
