package backups

import (
	"errors"
	"fpgschiba.com/automation/database"
	"fpgschiba.com/automation/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func ListBackups(c *gin.Context) {
	jobId := c.Query("jobId")
	backups, err := database.ListBackups(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to list backups"))
		return
	}
	c.JSON(http.StatusOK, listBackupsResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully fetched backups",
		},
		Backups: backups,
	})
}

func DeleteBackup(c *gin.Context) {
	backupId := c.Param("id")
	err := database.DeleteBackup(backupId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, util.GetErrorResponseAndMessage(err, "Backup not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to delete backup"))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse("Successfully deleted backup", true))
}

func DownloadBackup(c *gin.Context) {
	// Download a backup
	c.JSON(http.StatusNotImplemented, util.GetErrorResponseWithMessage("Not Implemented"))
}

func GetBackupLogs(c *gin.Context) {
	backupId := c.Param("id")
	severity := c.Query("severity")
	logs, err := database.GetBackupLogs(backupId, severity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, util.GetErrorResponseAndMessage(err, "Backup not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to fetch backup logs"))
		return
	}
	c.JSON(http.StatusOK, getBackupLogsResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully fetched backup logs",
		},
		Logs: logs,
	})
}

func GetBackupDetails(c *gin.Context) {
	backupId := c.Param("id")
	backup, err := database.GetBackupDetails(backupId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, util.GetErrorResponseAndMessage(err, "Backup not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseAndMessage(err, "Failed to fetch backup details"))
		return
	}
	c.JSON(http.StatusOK, getBackupResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Successfully fetched backup details",
		},
		Backup: backup,
	})
}
