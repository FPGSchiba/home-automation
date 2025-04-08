package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Components
const (
	ComponentAuth    = "auth"
	ComponentUsers   = "users"
	ComponentFinance = "finance"
	ComponentMeals   = "meals"
	ComponentBackups = "backups"
)

// Event types
const (
	EventLogin              = "login"
	EventPasswordReset      = "password-reset"
	EventBackupJobStarted   = "backup-job-started"
	EventBackupJobCompleted = "backup-job-completed"
	EventBackupJobFailed    = "backup-job-failed"
)

type Event struct {
	ID        string             `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	UserId    string             `bson:"userId"`
	TimeStamp primitive.DateTime `bson:"timeStamp"`
	Component string             `bson:"component"`
}
