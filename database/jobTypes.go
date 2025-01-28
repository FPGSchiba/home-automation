package database

import (
	"context"
	"fpgschiba.com/automation-meal/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jobTypesCollection *mongo.Collection

func GetJobTypesCollection(client *mongo.Client) *mongo.Collection {
	if jobTypesCollection != nil {
		return jobTypesCollection
	}
	jobTypesCollection = client.Database(DatabaseName).Collection("jobTypes")
	return jobTypesCollection
}

func ListJobTypes() (error, []models.JobType) {
	client := getClient()
	jobTypesCollection = GetJobTypesCollection(client)

	cursor, err := jobTypesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "ListJobTypes",
		}).Error()
		return err, nil
	}
	var jobTypes []models.JobType
	err = cursor.All(context.Background(), &jobTypes)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "ListJobTypes",
		}).Error()
		return err, nil
	}
	return nil, jobTypes
}
