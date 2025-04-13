package database

import (
	"context"
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/router/roles"
	"fpgschiba.com/automation-meal/util"
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
		"first_user":       migrationFirstUser,
		"add_job_types_v1": migrationAddJobTypesV1,
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
	hashedPassword, err := util.HashPassword("Password1")
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "migrationFirstUser",
		}).Error()
		return err
	}
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

func migrationAddJobTypesV1(client *mongo.Client) error {
	jobTypesCollection := GetJobTypesCollection(client)
	_, err := jobTypesCollection.InsertOne(context.Background(), models.JobType{
		Identifier: "sftp",
		Name:       "SFTP",
		ConfigurationFields: []models.ConfigurationField{
			{
				Name:        "Host",
				Description: "The Host of the SFTP Server",
				Type:        "string",
			},
			{
				Name:        "Port",
				Description: "The Port of the SFTP Server",
				Type:        "number",
			},
			{
				Name:        "Username",
				Description: "The Username for the SFTP Server",
				Type:        "string",
			},
			{
				Name:        "Password",
				Description: "The Password for the SFTP Server",
				Type:        "string",
			},
			{
				Name:        "Path",
				Description: "The Path to the File on the SFTP Server",
				Type:        "string",
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = jobTypesCollection.InsertOne(context.Background(), models.JobType{
		Identifier: "mongo",
		Name:       "MongoDB",
		ConfigurationFields: []models.ConfigurationField{
			{
				Name:        "Host",
				Description: "The Host of the MongoDB Server",
				Type:        "string",
			},
			{
				Name:        "DatabaseName",
				Description: "The Name of the Database on the MongoDB Server to authenticate the user on, use '/' for no authentication",
				Type:        "string",
			},
			{
				Name:        "Username",
				Description: "The Username for the MongoDB Server",
				Type:        "string",
			},
			{
				Name:        "Password",
				Description: "The Password for the MongoDB Server",
				Type:        "string",
			},
		},
	})
	return err
}
