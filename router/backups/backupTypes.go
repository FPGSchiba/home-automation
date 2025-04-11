package backups

import (
	"fpgschiba.com/automation-meal/models"
	"fpgschiba.com/automation-meal/util"
)

type listBackupsResponse struct {
	util.Response
	Backups []models.BackupMinimum `json:"backups"`
}

type getBackupResponse struct {
	util.Response
	Backup models.BackupDetails `json:"backup"`
}

type getBackupLogsResponse struct {
	util.Response
	Logs []models.BackupLog `json:"logs"`
}
