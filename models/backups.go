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
	Name          string                 `json:"name" bson:"name"`
	Identifier    string                 `json:"identifier" bson:"identifier"`
	Configuration map[string]interface{} `json:"configuration" bson:"configuration"`
	Schedule      string                 `json:"schedule" bson:"schedule"`
	SchedulerID   string                 `json:"schedulerId" bson:"schedulerId"` // Still needed as we use it to reference the job in the scheduler
}

type BackupJobRunStatus struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	JobID      primitive.ObjectID `json:"jobId" bson:"jobId"`
	StartedAt  primitive.DateTime `json:"startedAt" bson:"startedAt"`
	FinishedAt primitive.DateTime `json:"finishedAt" bson:"finishedAt"`
	Status     string             `json:"status" bson:"status"`
	Logs       []BackupLog        `json:"logs" bson:"logs"`
}

type BackupMinimum struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	StartedAt primitive.DateTime `json:"startedAt" bson:"startedAt"`
	Failed    bool               `json:"failed" bson:"failed"`
}

type BackupDetails struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	StartedAt  primitive.DateTime `json:"startedAt" bson:"startedAt"`
	FinishedAt primitive.DateTime `json:"finishedAt" bson:"finishedAt"`
	Size       int64              `json:"size" bson:"size"`
	Failed     bool               `json:"failed" bson:"failed"`
	Filename   string             `json:"filename" bson:"filename"`
}

type Backup struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Length     int64              `json:"length" bson:"length"`
	ChunkSize  int64              `json:"chunkSize" bson:"chunkSize"`
	UploadDate primitive.DateTime `json:"uploadDate" bson:"uploadDate"`
	Filename   string             `json:"filename" bson:"filename"`
	Metadata   BackupMetadata     `json:"metadata" bson:"metadata"`
}

type BackupMetadata struct {
	JobID      primitive.ObjectID `json:"jobId" bson:"jobId"`
	JobName    string             `json:"jobName" bson:"jobName"`
	StartedAt  primitive.DateTime `json:"startedAt" bson:"startedAt"`
	FinishedAt primitive.DateTime `json:"finishedAt" bson:"finishedAt"`
	Logs       []BackupLog        `json:"logs" bson:"logs"`
	Failed     bool               `json:"failed" bson:"failed"`
}

type BackupLog struct {
	Timestamp primitive.DateTime `json:"timestamp" bson:"timestamp"`
	Message   string             `json:"message" bson:"message"`
	Severity  string             `json:"severity" bson:"severity"`
}

type BackupSettings struct {
	BackupRetentionDays int `json:"backupRetentionDays" bson:"backupRetentionDays"`
}
