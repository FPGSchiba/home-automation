package backup

import (
	"fmt"
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/util"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var scheduler gocron.Scheduler

type extractedFields struct {
	Name  string
	Type  string
	Value interface{}
}

func getScheduler() gocron.Scheduler {
	if scheduler != nil {
		return scheduler
	}
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	return scheduler
}

func StartScheduler() {
	// Load all saved BackupJobs from the database & create them, then start the scheduler
	log.WithFields(log.Fields{
		"component": "backup",
		"func":      "StartScheduler",
	}).Info("Loading all backup jobs from the database")

	backupJobs, err := database.ListBackupJobs()
	if err != nil {
		panic(err)
	}

	for _, job := range backupJobs {
		log.WithFields(log.Fields{
			"component": "backup",
			"func":      "StartScheduler",
			"jobID":     job.ID,
			"jobType":   job.Identifier,
		}).Debug("loading backup job from database")
		switch job.Identifier {
		case "sftp":
			input := SFTPInput{
				Host:     job.Configuration["Host"].(string),
				Port:     int(job.Configuration["Port"].(float64)),
				Username: job.Configuration["Username"].(string),
				Password: job.Configuration["Password"].(string),
				Path:     job.Configuration["Path"].(string),
			}
			jobID, err := CreateSFTPBackupJob(input, job.Schedule, job.ID.Hex())
			if err != nil {
				panic(err)
			}
			// Update the job in the database with the new jobID
			err = database.UpdateBackupJobSchedulerID(job.ID.Hex(), jobID.String()) // This is still needed, as updates need to reference scheduler IDs
			if err != nil {
				panic(err)
			}
		case "mongo":
			// Create a mongo job

		}
	}
}

func StopScheduler() {
	// Stop the scheduler
	log.WithFields(log.Fields{
		"component": "backup",
		"func":      "StopScheduler",
	}).Info("Stopping scheduler")
	err := scheduler.Shutdown()
	if err != nil {
		panic(err)
	}
}

func getExtracts(config map[string]interface{}) []extractedFields {
	var extracts []extractedFields
	for key, value := range config {
		var valueType string
		switch value.(type) {
		case string:
			valueType = "string"
		case float64:
			valueType = "number"
		case int:
			valueType = "number"
		case bool:
			valueType = "boolean"
		case []interface{}:
			valueType = "array"
		case map[string]interface{}:
			valueType = "object"
		default:
			valueType = "unknown"
		}
		extracts = append(extracts, extractedFields{
			Name:  key,
			Type:  valueType,
			Value: value,
		})
	}
	return extracts
}

func RemoveJob(jobID uuid.UUID) error {
	scheduler := getScheduler()
	err := scheduler.RemoveJob(jobID)
	if err != nil {
		return err
	}
	return nil
}

func dummy() {
	// Dummy function to be able to add the JobID to a Job
}

func ValidateConfiguration(identifier string, config map[string]interface{}) (bool, error) {
	fields, err := database.GetConfigurationFieldsForJobType(identifier)
	if err != nil {
		return true, err
	}

	extracts := getExtracts(config)

	missing := fields
	var validationErrors util.JobTypeSchemaValidationErrors

	for _, extracted := range extracts {
		for _, field := range fields {
			if extracted.Name == field.Name {
				missing = util.Remove(missing, field)
				if extracted.Type != field.Type {
					validationErrors = append(validationErrors, util.JobTypeSchemaFieldError{
						FieldName:        field.Name,
						FieldType:        field.Type,
						ValidationResult: "invalid",
					})
				}
			}
		}
	}

	if len(missing) > 0 {
		for _, field := range missing {
			validationErrors = append(validationErrors, util.JobTypeSchemaFieldError{
				FieldName:        field.Name,
				FieldType:        field.Type,
				ValidationResult: "missing",
			})
		}
	}

	if len(validationErrors) > 0 {
		return false, &validationErrors
	}

	return true, nil
}

func handleRunFailure(jobID string, logs []models.BackupLog, startedAt time.Time, err error) {
	log.WithFields(log.Fields{
		"component": "backup",
		"func":      "handleRunFailure",
		"jobID":     jobID,
		"error":     err.Error(),
	}).Error("Backup job failed, beginning to handle failure")
	t := time.Now()
	jobName, err := database.GetJobNameFromID(jobID)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "backup",
			"func":      "handleRunFailure",
			"jobID":     jobID,
			"error":     err.Error(),
		}).Error("Failed to get job name from ID")
		return
	}

	jobIDObj, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "backup",
			"func":      "handleRunFailure",
			"jobID":     jobID,
			"error":     err.Error(),
		}).Error("Failed to convert jobID to ObjectID")
		return
	}

	meta := models.BackupMetadata{
		JobID:      jobIDObj,
		StartedAt:  primitive.NewDateTimeFromTime(startedAt),
		FinishedAt: primitive.NewDateTimeFromTime(t),
		Logs:       logs,
		JobName:    jobName,
		Failed:     true,
	}

	// Upload empty file to GridFS
	err = database.UploadEmptyFile(fmt.Sprintf("backup-%s-%s.zip", jobID, t.Format("2006-01-02-15-04")), meta)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "backup",
			"func":      "handleRunFailure",
			"jobID":     jobID,
			"error":     err.Error(),
		}).Error("Failed to upload failed Backup Results to GridFS")
		return
	}
}
