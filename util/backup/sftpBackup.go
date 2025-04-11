package backup

import (
	"archive/zip"
	"fmt"
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SFTPInput struct {
	Host     string
	Port     int
	Username string
	Password string
	Path     string
}

func CreateSFTPBackupJob(input SFTPInput, schedule string, jobID string) (uuid.UUID, error) {
	scheduler = getScheduler()
	job, err := scheduler.NewJob(gocron.CronJob(schedule, false), gocron.NewTask(runSFTPBackup, input, jobID))
	if err != nil {
		return uuid.UUID{}, err
	}
	scheduler.Start() // Need to start the scheduler to run the job
	return job.ID(), nil
}

// Main function that orchestrates the backup process
func runSFTPBackup(input SFTPInput, jobID string) {
	// Connect to SFTP
	startedAt := time.Now()
	var logs []models.BackupLog
	client, conn, err := connectToSFTP(input, &logs)
	if err != nil {
		logs = append(logs, models.BackupLog{
			Message:   fmt.Sprintf("Failed to connect to SFTP server %s:%d with error '%s'", input.Host, input.Port, err.Error()),
			Severity:  "error",
			Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		})
		handleRunFailure(jobID, logs, startedAt, err)
		return
	}
	defer conn.Close()
	defer client.Close()

	// Create temp directory
	tempDir, err := createTempDirectory(&logs)
	if err != nil {
		logs = append(logs, models.BackupLog{
			Message:   fmt.Sprintf("Failed to create temp directory: %s", err.Error()),
			Severity:  "error",
			Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		})
		handleRunFailure(jobID, logs, startedAt, err)
		return
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "backup",
				"jobID":     jobID,
				"error":     err.Error(),
			}).Error("Failed to remove temp directory")
		}
	}(tempDir)

	// Download files
	if err := downloadDirectory(client, input.Path, tempDir, jobID, &logs); err != nil {
		logs = append(logs, models.BackupLog{
			Message:   fmt.Sprintf("Failed to download directory: %s", err.Error()),
			Severity:  "error",
			Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		})
		handleRunFailure(jobID, logs, startedAt, err)
		return
	}

	// Create zip archive
	zipPath, err := createZipArchive(tempDir, jobID, &logs)
	if err != nil {
		logs = append(logs, models.BackupLog{
			Message:   fmt.Sprintf("Failed to create zip archive: %s", err.Error()),
			Severity:  "error",
			Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		})
		handleRunFailure(jobID, logs, startedAt, err)
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "backup",
				"jobID":     jobID,
				"error":     err.Error(),
			}).Error("Failed to remove zip file")
		}
	}(zipPath)

	// Upload to GridFS
	if err := uploadToGridFS(zipPath, jobID, startedAt, logs); err != nil {
		logs = append(logs, models.BackupLog{
			Message:   fmt.Sprintf("Failed to upload to GridFS: %s", err.Error()),
			Severity:  "error",
			Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		})
		handleRunFailure(jobID, logs, startedAt, err)
		return
	}

	log.WithFields(log.Fields{
		"component": "backup",
		"jobID":     jobID,
		"duration":  time.Since(startedAt),
	}).Info("SFTP backup completed successfully")
}

func connectToSFTP(input SFTPInput, logs *[]models.BackupLog) (*sftp.Client, *ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: input.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(input.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%v", input.Host, input.Port), config)
	*logs = append(*logs, models.BackupLog{
		Message:   fmt.Sprintf("Connected to SFTP server %s:%d", input.Host, input.Port),
		Severity:  "info",
		Timestamp: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("ssh connection failed: %w", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("sftp client creation failed: %w", err)
	}

	return client, conn, nil
}

func createTempDirectory(logs *[]models.BackupLog) (string, error) {
	tempDir := filepath.Join(os.TempDir(), "home-automation", fmt.Sprintf("sftp-backup-%s", time.Now().Format("2006-01-02-15-04")))
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	*logs = append(*logs, models.BackupLog{
		Message:   fmt.Sprintf("Created temp directory: %s", tempDir),
		Severity:  "info",
		Timestamp: primitive.NewDateTimeFromTime(time.Now()),
	})
	return tempDir, nil
}

func downloadFile(client *sftp.Client, remotePath, localPath string) error {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Open remote file
	remoteFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// Copy file contents
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func downloadDirectory(client *sftp.Client, remotePath, tempDir, jobID string, logs *[]models.BackupLog) error {
	w := client.Walk(remotePath)
	for w.Step() {
		if w.Err() != nil {
			*logs = append(*logs, models.BackupLog{
				Message:   fmt.Sprintf("Error accessing path %s: %s", w.Path(), w.Err().Error()),
				Severity:  "warning",
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
			})
			continue
		}

		currentPath := w.Path()
		localPath := filepath.Join(tempDir, strings.TrimPrefix(currentPath, remotePath))

		// Get file info
		fi, err := client.Stat(currentPath)
		if err != nil {
			*logs = append(*logs, models.BackupLog{
				Message:   fmt.Sprintf("Error getting file info for %s: %s", currentPath, err.Error()),
				Severity:  "warning",
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
			})
			continue
		}

		if fi.IsDir() {
			if err := os.MkdirAll(localPath, 0755); err != nil {
				*logs = append(*logs, models.BackupLog{
					Message:   fmt.Sprintf("Error creating local directory %s: %s", localPath, err.Error()),
					Severity:  "warning",
					Timestamp: primitive.NewDateTimeFromTime(time.Now()),
				})
			} else {
				*logs = append(*logs, models.BackupLog{
					Message:   fmt.Sprintf("Created local directory %s", localPath),
					Severity:  "debug",
					Timestamp: primitive.NewDateTimeFromTime(time.Now()),
				})
			}

			continue
		}

		if err := downloadFile(client, currentPath, localPath); err != nil {
			*logs = append(*logs, models.BackupLog{
				Message:   fmt.Sprintf("Error downloading file %s: %s", currentPath, err.Error()),
				Severity:  "error",
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
			})
		} else {
			*logs = append(*logs, models.BackupLog{
				Message:   fmt.Sprintf("Downloaded file %s to %s", currentPath, localPath),
				Severity:  "debug",
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
			})
		}
	}

	return nil
}

func createZipArchive(tempDir, jobID string, logs *[]models.BackupLog) (string, error) {
	localTempDir := filepath.Join(os.TempDir(), "home-automation")
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		*logs = append(*logs, models.BackupLog{
			Message:   fmt.Sprintf("Failed to create temp directory: %s", err.Error()),
			Severity:  "error",
			Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		})
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	zipPath := filepath.Join(localTempDir, fmt.Sprintf("backup-%s-%s.zip", jobID, time.Now().Format("2006-01-02-15-04")))
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create zip entry: %w", err)
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return fmt.Errorf("failed to copy to zip: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to create zip archive: %w", err)
	}

	*logs = append(*logs, models.BackupLog{
		Message:   fmt.Sprintf("Created zip archive: %s", zipPath),
		Severity:  "info",
		Timestamp: primitive.NewDateTimeFromTime(time.Now()),
	})

	return zipPath, nil
}

// Function to upload to MongoDB GridFS
func uploadToGridFS(zipPath string, jobID string, startedAt time.Time, logs []models.BackupLog) error {
	// Upload to GridFS
	t := time.Now()
	// Get the job from the database
	jobName, err := database.GetJobNameFromID(jobID)
	if err != nil {
		return err
	}
	// Convert jobID to primitive.ObjectID
	jobIDObj, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}
	meta := models.BackupMetadata{
		JobID:      jobIDObj,
		StartedAt:  primitive.NewDateTimeFromTime(startedAt),
		FinishedAt: primitive.NewDateTimeFromTime(t),
		Logs:       logs,
		JobName:    jobName,
		Failed:     false,
	}

	logs = append(logs, models.BackupLog{
		Message:   fmt.Sprintf("Uploading backup to GridFS...", meta),
		Severity:  "info",
		Timestamp: primitive.NewDateTimeFromTime(time.Now()),
	})

	err = database.UploadBackup(fmt.Sprintf("backup-%s-%s.zip", jobID, t.Format("2006-01-02-15-04")), zipPath, meta)
	if err != nil {
		return err
	}
	return nil
}
