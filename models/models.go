package models

import (
	"fpgschiba.com/automation-meal/router/roles"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                *primitive.ObjectID `bson:"_id,omitempty"`
	Email             string              `bson:"email" json:"email"`
	DisplayName       string              `bson:"displayName" json:"displayName"`
	ProfilePictureURL string              `bson:"profilePictureUrl" json:"profilePictureUrl"`
	PasswordHash      string              `bson:"passwordHash" json:"passwordHash"`
}

type Role struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty"`
	Name        string              `bson:"name" json:"name"`
	Permissions []roles.Permission  `bson:"permissions" json:"permissions"`
}

type RoleAssignment struct {
	ID     *primitive.ObjectID `bson:"_id,omitempty"`
	RoleID primitive.ObjectID  `bson:"roleId" json:"roleId"`
	UserID primitive.ObjectID  `bson:"userId" json:"userId"`
}
