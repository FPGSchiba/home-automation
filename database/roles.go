package database

import (
	"context"
	"fpgschiba.com/automation/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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

func GetRoleById(roleId primitive.ObjectID) (error, models.Role) {
	client = getClient()
	rolesCollection = GetRolesCollection(client)
	result := rolesCollection.FindOne(context.Background(), bson.M{"_id": roleId})
	if result.Err() != nil {
		log.WithFields(log.Fields{
			"error":     result.Err(),
			"component": "database",
			"func":      "GetRoleById",
			"roleId":    roleId,
		}).Error()
		return result.Err(), models.Role{}
	}
	var role models.Role
	err := result.Decode(&role)
	if err != nil {
		return err, models.Role{}
	}
	return nil, role
}

func GetAllRolesForUser(userId primitive.ObjectID) (error, []models.Role) {
	client = getClient()
	rolesCollection = GetRolesCollection(client)
	roleAssignmentsCollection = GetRoleAssignmentsCollection(client)

	err, roleAssignments := GetRoleAssignmentsByUserId(userId)
	if err != nil {
		return err, nil
	}

	var roles []models.Role
	for _, assignment := range roleAssignments {
		err, role := GetRoleById(assignment.RoleID)
		if err != nil {
			return err, nil
		}
		roles = append(roles, role)
	}
	return nil, roles
}
