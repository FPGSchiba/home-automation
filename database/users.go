package database

import (
	"context"
	"fpgschiba.com/automation/models"
	"fpgschiba.com/automation/util"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection *mongo.Collection

func GetUsersCollection(client *mongo.Client) *mongo.Collection {
	if usersCollection != nil {
		return usersCollection
	}
	usersCollection = client.Database(DatabaseName).Collection("users")
	return usersCollection
}

func CreateUser(user models.User) (error, string) {
	client = getClient()
	usersCollection = GetUsersCollection(client)

	result, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"user":      user,
			"component": "database",
			"func":      "CreateUser",
		}).Error()
		return err, ""
	}
	return nil, result.InsertedID.(primitive.ObjectID).Hex()
}

func DoesEmailExist(email string) (error, bool, string) {
	client = getClient()
	usersCollection = GetUsersCollection(client)

	result := usersCollection.FindOne(context.Background(), bson.M{"email": email})
	if result.Err() != nil {
		log.WithFields(log.Fields{
			"error":     result.Err(),
			"component": "database",
			"func":      "DoesEmailExist",
		}).Error()
		return result.Err(), false, ""
	}
	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return err, false, ""
	}
	return nil, true, user.ID.Hex()
}

func PasswordMatchesUser(userId primitive.ObjectID, password string) (error, bool, *models.User) {
	client = getClient()
	usersCollection = GetUsersCollection(client)

	result := usersCollection.FindOne(context.Background(), bson.M{"_id": userId})
	if result.Err() != nil {
		log.WithFields(log.Fields{
			"error":     result.Err(),
			"component": "database",
			"func":      "DoesEmailExist",
		}).Error()
		return result.Err(), false, nil
	}
	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return err, false, nil
	}
	if util.VerifyPassword(user.PasswordHash, password) {
		return nil, true, &user
	}
	return nil, false, nil
}

func GetUserIDByEmail(email string) (error, string) {
	client = getClient()
	usersCollection = GetUsersCollection(client)

	result := usersCollection.FindOne(context.Background(), bson.M{"email": email})
	if result.Err() != nil {
		log.WithFields(log.Fields{
			"error":     result.Err(),
			"component": "database",
			"func":      "GetUserIDByEmail",
			"email":     email,
		}).Error()
		return result.Err(), ""
	}
	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return err, ""
	}
	return nil, user.ID.Hex()
}

func GetAllUsers() (error, []models.User) {
	client = getClient()
	usersCollection = GetUsersCollection(client)

	cursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "GetAllUsers",
		}).Error()
		return err, nil
	}
	var users []models.User
	err = cursor.All(context.Background(), &users)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "database",
			"func":      "GetAllUsers",
		}).Error()
		return err, nil
	}
	return nil, users
}
