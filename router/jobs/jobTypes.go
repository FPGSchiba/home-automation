package jobs

import (
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/util"
)

type ListJobTypesResponse struct {
	util.Response
	JobTypes []models.JobType `json:"jobTypes"`
}
