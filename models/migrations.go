package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Migration struct {
	ID   string             `bson:"_id,omitempty"`
	Name string             `bson:"name" json:"name"`
	Time primitive.DateTime `bson:"time" json:"time"`
}
