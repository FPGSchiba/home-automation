package database

import (
	"context"
	"fpgschiba.com/automation-meal/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var rolesCollection *mongo.Collection

func GetRolesCollection(client *mongo.Client) *mongo.Collection {
	if rolesCollection != nil {
		return rolesCollection
	}
	rolesCollection = client.Database(DatabaseName).Collection("roles")
	return rolesCollection
}

func CreateRole(role models.Role) (error, string) {
	client = getClient()
	rolesCollection = GetRolesCollection(client)
	result, err := rolesCollection.InsertOne(context.Background(), role)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"role":      role,
			"component": "database",
			"func":      "CreateRole",
		}).Error()
		return err, ""
	}
	return nil, result.InsertedID.(primitive.ObjectID).Hex()
}
