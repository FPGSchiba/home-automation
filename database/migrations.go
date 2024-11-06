package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/router/roles"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var migrationCollection *mongo.Collection

func StartMigrations() {
	client = getClient()
	//getting collection
	migrations := map[string]interface{}{
		"first_user": migrationFirstUser,
	}
	migrationCollection := GetMigrationCollection(client)
	doneMigrations, err := migrationCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "StartMigrations",
		}).Error()
		return
	}
	for doneMigrations.Next(context.Background()) {
		var migration models.Migration
		err := doneMigrations.Decode(&migration)
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err,
				"component": "database",
				"func":      "StartMigrations",
			}).Error()
			return
		}
		delete(migrations, migration.Name)
	}
	for name, migration := range migrations {
		log.WithFields(log.Fields{
			"migration": name,
			"component": "database",
			"func":      "StartMigrations",
		}).Info()
		err := migration.(func(*mongo.Client) error)(client)
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err,
				"component": "database",
				"func":      "StartMigrations",
			}).Error()
			return
		}
		_, err = migrationCollection.InsertOne(context.Background(), models.Migration{
			Name: name,
			Time: primitive.NewDateTimeFromTime(time.Now()),
		})
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err,
				"migration": name,
				"component": "database",
				"func":      "StartMigrations",
			}).Error()
			return
		}
	}
}

func GetMigrationCollection(client *mongo.Client) *mongo.Collection {
	if migrationCollection != nil {
		return migrationCollection
	}
	migrationCollection = client.Database(DatabaseName).Collection("migrations")
	return migrationCollection
}

func migrationFirstUser(client *mongo.Client) error {
	usersCollection := GetUsersCollection(client)
	rolesCollection := GetRolesCollection(client)
	roleAssignmentsCollection := GetRoleAssignmentsCollection(client)
	//creating roles
	roleResult, err := rolesCollection.InsertOne(context.Background(), models.Role{
		Name:        "SuperUser",
		Description: "SuperUser has all permissions",
		Permissions: []roles.Permission{
			{
				Routes: []roles.Route{
					{
						Path: "**",
						Methods: []string{
							"GET",
							"POST",
							"PUT",
							"DELETE",
							"PATCH",
						},
					},
				},
			},
		},
	})
	//creating User
	h := sha256.New()
	h.Write([]byte("Password1")) // TODO: Change this to a not clear text password
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	userResult, err := usersCollection.InsertOne(context.Background(), models.User{
		Email:             "jann.erhardt@icloud.com",
		DisplayName:       "Jann Erhardt",
		ProfilePictureURL: "",
		PasswordHash:      hashedPassword,
	})
	if err != nil {
		return err
	}
	_, err = roleAssignmentsCollection.InsertOne(context.Background(), models.RoleAssignment{
		RoleID: roleResult.InsertedID.(primitive.ObjectID),
		UserID: userResult.InsertedID.(primitive.ObjectID),
	})
	return err
}
