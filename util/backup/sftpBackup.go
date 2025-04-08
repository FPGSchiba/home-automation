package backup

import (
	"archive/zip"
	"fmt"
	"fpgschiba.com/automation-meal/database"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func CreateSFTPBackupJob(input SFTPInput, schedule string) (uuid.UUID, error) {
	scheduler = getScheduler()
	job, err := scheduler.NewJob(gocron.CronJob(schedule, false), gocron.NewTask(dummy))
	if err != nil {
		return uuid.UUID{}, err
	}
	_, err = scheduler.Update(job.ID(), gocron.CronJob(schedule, false), gocron.NewTask(runSFTPBackup, job.ID(), input))
	if err != nil {
		return uuid.UUID{}, err
	}
	scheduler.Start() // Need to start the scheduler to run the job
	return job.ID(), nil
}

// Main function that orchestrates the backup process
func runSFTPBackup(input SFTPInput, jobID string, db *mongo.Database) error {
	// Connect to SFTP
	client, conn, err := connectToSFTP(input)
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP: %w", err)
	}
	defer conn.Close()
	defer client.Close()

	// Create temp directory
	tempDir, err := createTempDirectory()
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download files
	if err := downloadDirectory(client, input.Path, tempDir, jobID); err != nil {
		return fmt.Errorf("failed to download directory: %w", err)
	}

	// Create zip archive
	zipPath, err := createZipArchive(tempDir, jobID)
	if err != nil {
		return fmt.Errorf("failed to create zip archive: %w", err)
	}
	defer os.Remove(zipPath)

	// Upload to GridFS
	if err := uploadToGridFS(zipPath, jobID, input.Path); err != nil {
		return fmt.Errorf("failed to upload to GridFS: %w", err)
	}

	return nil
}

func connectToSFTP(input SFTPInput) (*sftp.Client, *ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: input.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(input.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%v", input.Host, input.Port), config)
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

func createTempDirectory() (string, error) {
	tempDir, err := os.MkdirTemp("", "sftp-backup-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
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

func downloadDirectory(client *sftp.Client, remotePath, tempDir, jobID string) error {
	w := client.Walk(remotePath)
	for w.Step() {
		if w.Err() != nil {
			log.WithFields(log.Fields{
				"component": "backup",
				"jobID":     jobID,
				"path":      w.Path(),
				"error":     w.Err(),
			}).Warn("Error accessing path")
			continue
		}

		currentPath := w.Path()
		localPath := filepath.Join(tempDir, strings.TrimPrefix(currentPath, remotePath))

		// Get file info
		fi, err := client.Stat(currentPath)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "backup",
				"jobID":     jobID,
				"path":      currentPath,
				"error":     err,
			}).Warn("Error getting file info")
			continue
		}

		if fi.IsDir() {
			if err := os.MkdirAll(localPath, 0755); err != nil {
				log.WithFields(log.Fields{
					"component": "backup",
					"jobID":     jobID,
					"path":      localPath,
					"error":     err,
				}).Warn("Error creating local directory")
			}
			continue
		}

		if err := downloadFile(client, currentPath, localPath); err != nil {
			log.WithFields(log.Fields{
				"component": "backup",
				"jobID":     jobID,
				"path":      currentPath,
				"error":     err,
			}).Warn("Error downloading file")
		}
	}

	return nil
}

func createZipArchive(tempDir, jobID string) (string, error) {
	zipPath := filepath.Join(os.TempDir(), fmt.Sprintf("backup-%s.zip", jobID))
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

	return zipPath, nil
}

// Function to upload to MongoDB GridFS
func uploadToGridFS(zipPath, jobID, sourcePath string) error {
	// Upload to GridFS
	t := time.Now()
	meta := bson.M{
		"jobID":      jobID,
		"sourcePath": sourcePath,
		"createdAt":  t,
	}
	err := database.UploadBackup(fmt.Sprintf("backup-%s-%s.zip", jobID, t.Format("2006-01-02-15-04")), zipPath, meta)
	if err != nil {
		return err
	}
	return nil
}
