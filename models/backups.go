package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobType struct {
	ID                  string               `json:"id" bson:"_id,omitempty"`
	Identifier          string               `json:"identifier" bson:"identifier"`
	Name                string               `json:"name" bson:"name"`
	ConfigurationFields []ConfigurationField `json:"configurationFields" bson:"configurationFields"`
}

type ConfigurationField struct {
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}

type BackupJob struct {
	ID            *primitive.ObjectID    `json:"id" bson:"_id,omitempty"` // this is how we address the Jobs
	Identifier    string                 `json:"identifier" bson:"identifier"`
	Configuration map[string]interface{} `json:"configuration" bson:"configuration"`
	Schedule      string                 `json:"schedule" bson:"schedule"`
	SchedulerID   string                 `json:"schedulerId" bson:"schedulerId"` // Still needed as we use it to reference the job in the scheduler
}
