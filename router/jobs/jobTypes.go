package jobs

import (
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/util"
)

type ListJobTypesResponse struct {
	util.Response
	JobTypes []models.JobType `json:"jobTypes"`
}

type createJobRequest struct {
	Name              string                 `json:"name" binding:"required"`
	JobTypeIdentifier string                 `json:"jobTypeIdentifier" binding:"required"`
	Configuration     map[string]interface{} `json:"configuration" binding:"required"`
	Schedule          string                 `json:"schedule" binding:"required"`
}

type createJobResponse struct {
	util.Response
	Job models.BackupJob `json:"job"`
}

type listJobsResponse struct {
	util.Response
	Jobs []models.BackupJob `json:"jobs"`
}

type getJobResponse struct {
	util.Response
	Job models.BackupJob `json:"job"`
}

type updateJobRequest struct {
	Name              string                 `json:"name" binding:"required"`
	JobTypeIdentifier string                 `json:"jobTypeIdentifier" binding:"required"`
	Configuration     map[string]interface{} `json:"configuration" binding:"required"`
	Schedule          string                 `json:"schedule" binding:"required"`
}
