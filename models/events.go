package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	EventAuth    = "auth"
	EventUsers   = "users"
	EventFinance = "finance"
	EventMeals   = "meals"
)

type Event struct {
	ID        string             `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	UserId    string             `bson:"userId"`
	TimeStamp primitive.DateTime `bson:"timeStamp"`
	Component string             `bson:"component"`
}
