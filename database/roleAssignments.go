package database

import (
	"context"
	"fpgschiba.com/automation-meal/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var roleAssignmentsCollection *mongo.Collection

func GetRoleAssignmentsCollection(client *mongo.Client) *mongo.Collection {
	if roleAssignmentsCollection != nil {
		return roleAssignmentsCollection
	}
	roleAssignmentsCollection = client.Database(DatabaseName).Collection("roleAssignments")
	return roleAssignmentsCollection
}

func CreateRoleAssignment(assignment models.RoleAssignment) (error, string) {
	client = getClient()
	roleAssignmentsCollection = GetRoleAssignmentsCollection(client)

	result, err := roleAssignmentsCollection.InsertOne(context.Background(), assignment)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err,
			"roleAssignment": assignment,
			"component":      "database",
			"func":           "CreateRoleAssignment",
		}).Error()
		return err, ""
	}
	return nil, result.InsertedID.(primitive.ObjectID).Hex()
}
