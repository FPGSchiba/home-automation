package jobs

import (
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListJobs(c *gin.Context) {
	// List all jobs
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func CreateJob(c *gin.Context) {
	// Create a job
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func GetJob(c *gin.Context) {
	// Get a job
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func UpdateJob(c *gin.Context) {
	// Update a job
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func DeleteJob(c *gin.Context) {
	// Delete a job
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func GetJobTypes(c *gin.Context) {
	err, jobTypes := database.ListJobTypes()
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
