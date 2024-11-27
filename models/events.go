package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Components
const (
	ComponentAuth    = "auth"
	ComponentUsers   = "users"
	ComponentFinance = "finance"
	ComponentMeals   = "meals"
)

// Event types
const (
	EventLogin         = "login"
	EventPasswordReset = "password-reset"
)

type Event struct {
	ID        string             `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	UserId    string             `bson:"userId"`
	TimeStamp primitive.DateTime `bson:"timeStamp"`
	Component string             `bson:"component"`
}
