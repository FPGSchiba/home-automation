package database

import (
	"context"
	"fmt"
	"fpgschiba.com/automation-meal/util"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var client *mongo.Client

var (
	DatabaseName = "HomeAutomation"
)

func getClient() *mongo.Client {
	config := util.Config{}
	config.GetConfig()
	DatabaseName = config.Database.Database
	uri := fmt.Sprintf("mongodb://%s:%d", config.Database.Host, config.Database.Port)
	//getting context
	if client != nil {
		return client
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//getting client
	opts := options.Client()
	if config.Database.User != "" && config.Database.Password != "" {
		credentials := options.Credential{
			Username:   config.Database.User,
			Password:   config.Database.Password,
			AuthSource: DatabaseName,
		}
		opts.SetAuth(credentials)
	}

	client, err := mongo.Connect(ctx, opts.ApplyURI(uri))
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
