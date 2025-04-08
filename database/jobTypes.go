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

func GetConfigurationFieldsForJobType(jobType string) ([]models.ConfigurationField, error) {
	client := getClient()
	jobTypesCollection = GetJobTypesCollection(client)
	var jobTypeDoc models.JobType
	err := jobTypesCollection.FindOne(context.Background(), bson.M{"identifier": jobType}).Decode(&jobTypeDoc)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "GetConfigurationFieldsForJobType",
		}).Error()
		return nil, err
	}
	return jobTypeDoc.ConfigurationFields, nil
}

func ListJobTypes() ([]models.JobType, error) {
	client := getClient()
	jobTypesCollection = GetJobTypesCollection(client)

	cursor, err := jobTypesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "ListJobTypes",
		}).Error()
		return nil, err
	}
	var jobTypes []models.JobType
	err = cursor.All(context.Background(), &jobTypes)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "ListJobTypes",
		}).Error()
		return nil, err
	}
	return jobTypes, nil
}
