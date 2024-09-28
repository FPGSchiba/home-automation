package database

import (
	"context"
	"fpgschiba.com/automation-meal/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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

func GetRoleAssignmentsByUserId(userId primitive.ObjectID) (error, []models.RoleAssignment) {
	client = getClient()
	roleAssignmentsCollection = GetRoleAssignmentsCollection(client)

	cursor, err := roleAssignmentsCollection.Find(context.Background(), bson.M{"userId": userId})
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"func":      "GetRoleAssignmentsByUserId",
			"component": "database",
			"userId":    userId,
		}).Error()
		return err, nil
	}
	var assignments []models.RoleAssignment
	err = cursor.All(context.Background(), &assignments)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error()
		return err, nil
	}
	return nil, assignments
}
