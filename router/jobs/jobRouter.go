package jobs

import (
	"errors"
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/util"
	"fpgschiba.com/automation-meal/util/backup"
	"github.com/adhocore/gronx"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	uuid2 "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func validateJobType(jobType string, config map[string]interface{}) (bool, string, error) {
	// Validate the job type
	dbError, err := backup.ValidateConfiguration(jobType, config)
	if dbError && err != nil {
		return false, "Invalid job type identifier", err
	}
	if err != nil {
		return false, "Invalid configuration", err
	}
	return true, "", nil
}

func validateSchedule(schedule string) bool {
	gron := gronx.New()
	if !gron.IsValid(schedule) {
		return false
	}
	return true
}

/*
func saveJobToDatabase(job models.BackupJob, schedulerID uuid2.UUID) (bool, string, error) {
	jobID, err := database.InsertBackupJob(job)
	if err != nil {
		err = backup.RemoveJob(schedulerID)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "CreateJob",
				"method":    "saveJobToDatabase",
				"path":      "/jobs",
				"error":     err.Error(),
			}).Error("Failed to remove job from scheduler")
			return false, "Failed to create Job in Database and Failed remove job from scheduler", err
		}
		return false, "Failed to create SFTP job", err
	}
	log.WithFields(log.Fields{
		"component": "CreateJob",
		"method":    "saveJobToDatabase",
		"path":      "/jobs",
		"jobID":     jobID,
	}).Info("SFTP job created successfully")
	return true, jobID, nil
}
*/

func handleJobCreationOrUpdate(jobID string, jobType string, config map[string]interface{}, schedule string, oldSchedulerID *uuid2.UUID) (models.BackupJob, uuid2.UUID, error) {
	var schedulerID uuid2.UUID
	var err error
	var job models.BackupJob

	switch jobType {
	case "mongo":
		input := backup.MongoInput{
			Host:         config["Host"].(string),
			Username:     config["Username"].(string),
			Password:     config["Password"].(string),
			DatabaseName: config["DatabaseName"].(string),
		}
		schedulerID, err = backup.CreateMongoBackupJob(input, schedule, jobID)
	case "sftp":
		input := backup.SFTPInput{
			Host:     config["Host"].(string),
			Port:     int(config["Port"].(float64)),
			Username: config["Username"].(string),
			Password: config["Password"].(string),
			Path:     config["Path"].(string),
		}
		schedulerID, err = backup.CreateSFTPBackupJob(input, schedule, jobID)
	default:
		return job, schedulerID, errors.New("unsupported job type")
	}

	if err != nil {
		return job, schedulerID, err
	}

	// Remove old job from scheduler if updating
	if oldSchedulerID != nil {
		err = backup.RemoveJob(*oldSchedulerID)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "JobHandler",
				"error":     err.Error(),
			}).Error("Failed to remove old job from scheduler")
		}
	}

	job = models.BackupJob{
		Identifier:    jobType,
		Configuration: config,
		Schedule:      schedule,
		SchedulerID:   schedulerID.String(),
	}

	return job, schedulerID, nil
}

func handleJobResponse(c *gin.Context, success bool, update bool, message string, err error) {
	if !success {
		log.WithFields(log.Fields{
			"component": "JobHandler",
			"error":     err,
		}).Error(message)
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, message))
		return
	}
	if update {
		c.JSON(http.StatusOK, util.GetResponse("Successfully updated job", true))
		return
	}
	c.JSON(http.StatusCreated, createJobResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully created job",
		},
		JobID: message,
	})
}

func ListJobs(c *gin.Context) {
	jobs, err := database.ListBackupJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to fetch jobs"))
		return
	}
	c.JSON(http.StatusOK, listJobsResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully fetched jobs",
		},
		Jobs: jobs,
	})
}

func CreateJob(c *gin.Context) {
	body := createJobRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}

	if !validateSchedule(body.Schedule) {
		c.JSON(http.StatusBadRequest, util.GetErrorResponseWithMessage("Invalid schedule"))
		return
	}

	valid, message, err := validateJobType(body.JobTypeIdentifier, body.Configuration)
	if !valid {
		c.JSON(http.StatusBadRequest, util.GetErrorResponseAndMessage(err, message))
		return
	}
	job := models.BackupJob{
		Identifier:    body.JobTypeIdentifier,
		Configuration: body.Configuration,
		Schedule:      body.Schedule,
		Name:          body.Name,
	}

	jobID, err := database.InsertBackupJob(job)
	if err != nil {
		return
	}

	_, schedulerID, err := handleJobCreationOrUpdate(jobID, body.JobTypeIdentifier, body.Configuration, body.Schedule, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to create job"))
		err := database.DeleteBackupJob(jobID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to create job and failed to delete job from database"))
			return
		}
		return
	}

	err = database.UpdateBackupJobSchedulerID(jobID, schedulerID.String())

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to update job in database"))
		err = backup.RemoveJob(schedulerID)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "CreateJob",
				"method":    c.Request.Method,
				"path":      c.Request.URL.Path,
				"error":     err.Error(),
			}).Error("Failed to remove job from scheduler")
			return
		}
		return
	}

	handleJobResponse(c, true, false, jobID, err)
}

func GetJob(c *gin.Context) {
	jobId := c.Param("id")
	job, err := database.GetBackupJobFromID(jobId)
	if err != nil {
		if errors.As(err, &mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, util.GetErrorResponseAndMessage(err, "Job not found"))
		} else {
			c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to fetch job"))
		}

		return
	}
	c.JSON(http.StatusOK, getJobResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully fetched job",
		},
		Job: job,
	})
}

func UpdateJob(c *gin.Context) {
	jobId := c.Param("id")
	body := updateJobRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}

	job, err := database.GetBackupJobFromID(jobId)
	if err != nil {
		if errors.As(err, &mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, util.GetErrorResponseAndMessage(err, "Job not found"))
		} else {
			c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to fetch job"))
		}
		return
	}

	if !validateSchedule(body.Schedule) {
		c.JSON(http.StatusBadRequest, util.GetErrorResponseWithMessage("Invalid schedule"))
		return
	}

	valid, message, err := validateJobType(body.JobTypeIdentifier, body.Configuration)
	if !valid {
		c.JSON(http.StatusBadRequest, util.GetErrorResponseAndMessage(err, message))
		return
	}

	oldSchedulerID := uuid2.UUID(uuid.FromStringOrNil(job.SchedulerID))
	updatedJob, _, err := handleJobCreationOrUpdate(jobId, body.JobTypeIdentifier, body.Configuration, body.Schedule, &oldSchedulerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to update job"))
		return
	}

	updatedJob.ID = job.ID // Preserve the original job ID
	err = database.UpdateBackupJob(jobId, updatedJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to update job in database"))
		return
	}

	c.JSON(http.StatusOK, util.GetResponse("Successfully updated job", true))
}

func DeleteJob(c *gin.Context) {
	jobId := c.Param("id")
	job, err := database.GetBackupJobFromID(jobId)
	if err != nil {
		if errors.As(err, &mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, util.GetErrorResponseAndMessage(err, "Job not found"))
		} else {
			c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to delete job"))
		}
		return
	}
	err = database.DeleteBackupJob(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to delete job"))
		return
	}
	err = backup.RemoveJob(uuid2.UUID(uuid.FromStringOrNil(job.SchedulerID)))
	if err != nil {
		log.WithFields(log.Fields{
			"component": "DeleteJob",
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"error":     err.Error(),
		}).Error("Failed to remove job from scheduler")

		// If the job could not be removed from the scheduler we will need to restore the Job in the database
		_, err = database.InsertBackupJob(job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to delete Backup Job in scheduler and failed to restore job in database"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to delete job"))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse("Job deleted successfully", true))
}

func GetJobTypes(c *gin.Context) {
	jobTypes, err := database.ListJobTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ListJobTypesResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully fetched job types",
		},
		JobTypes: jobTypes,
	})
}
