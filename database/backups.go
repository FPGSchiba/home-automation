package database

import (
	"fmt"
	"fpgschiba.com/automation/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"os"
)

var backupsGridFS *gridfs.Bucket
var backupSettingsCollection *mongo.Collection

func GetBackupsBucket(client *mongo.Client) *gridfs.Bucket {
	if backupsGridFS != nil {
		return backupsGridFS
	}

	backupsGridFS, _ = gridfs.NewBucket(client.Database(DatabaseName))

	return backupsGridFS
}

func GetBackupSettingsCollection(client *mongo.Client) *mongo.Collection {
	if backupSettingsCollection != nil {
		return backupSettingsCollection
	}

	backupSettingsCollection = client.Database(DatabaseName).Collection("backupSettings")

	return backupSettingsCollection
}

func GetBackupSettings() (models.BackupSettings, error) {
	client := getClient()
	collection := GetBackupSettingsCollection(client)
	var settings models.BackupSettings
	err := collection.FindOne(nil, bson.M{}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.BackupSettings{}, nil
		}
		return models.BackupSettings{}, err
	}
	return settings, nil
}

func UpdateBackupSettings(settings models.BackupSettings) error {
	client := getClient()
	collection := GetBackupSettingsCollection(client)
	filter := bson.M{}
	update := bson.M{
		"$set": settings,
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(nil, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func UploadBackup(filename string, filePath string, meta models.BackupMetadata) error {
	client := getClient()
	bucket := GetBackupsBucket(client)
	uploadOpts := options.GridFSUpload().
		SetMetadata(meta)
	uploadStream, err := bucket.OpenUploadStream(filename, uploadOpts)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(uploadStream, file); err != nil {
		return fmt.Errorf("failed to upload to GridFS: %w", err)
	}
	return nil
}

func ListBackups(jobId string) ([]models.BackupMinimum, error) {
	client := getClient()
	bucket := GetBackupsBucket(client)
	var filter bson.M

	if jobId != "" {
		jobIdObj, err := primitive.ObjectIDFromHex(jobId)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"metadata.jobId": jobIdObj}
	} else {
		filter = bson.M{}
	}

	cursor, err := bucket.Find(filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var backups []models.Backup
	for cursor.Next(nil) {
		var backup models.Backup
		if err := cursor.Decode(&backup); err != nil {
			return nil, err
		}
		backups = append(backups, backup)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	var backupsMin []models.BackupMinimum
	for _, backup := range backups {
		backupsMin = append(backupsMin, models.BackupMinimum{
			ID:        backup.ID,
			Name:      backup.Metadata.JobName,
			StartedAt: backup.Metadata.StartedAt,
			Failed:    backup.Metadata.Failed,
		})
	}
	return backupsMin, nil
}

func UploadEmptyFile(filename string, meta models.BackupMetadata) error {
	client := getClient()
	bucket := GetBackupsBucket(client)
	uploadOpts := options.GridFSUpload().
		SetMetadata(meta)
	uploadStream, err := bucket.OpenUploadStream(filename, uploadOpts)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	tempFile, err := os.CreateTemp("", "empty-file-*.tmp")
	if err != nil {
		return err
	}

	defer os.Remove(tempFile.Name())
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(uploadStream, file); err != nil {
		return fmt.Errorf("failed to upload to GridFS: %w", err)
	}
	return nil
}

func DeleteBackup(id string) error {
	client := getClient()
	bucket := GetBackupsBucket(client)
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = bucket.Delete(idObj)
	if err != nil {
		return err
	}
	return nil
}

func GetBackupDetails(id string) (models.BackupDetails, error) {
	client := getClient()
	bucket := GetBackupsBucket(client)
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.BackupDetails{}, err
	}
	var backup models.Backup
	cursor, err := bucket.Find(bson.M{"_id": idObj})
	if err != nil {
		return models.BackupDetails{}, err
	}
	defer cursor.Close(nil)
	if cursor.Next(nil) {
		if err := cursor.Decode(&backup); err != nil {
			return models.BackupDetails{}, err
		}
	} else {
		return models.BackupDetails{}, fmt.Errorf("no backup found with id %s", id)
	}
	if err := cursor.Err(); err != nil {
		return models.BackupDetails{}, err
	}
	backupDetails := models.BackupDetails{
		ID:         backup.ID,
		Name:       backup.Metadata.JobName,
		StartedAt:  backup.Metadata.StartedAt,
		FinishedAt: backup.Metadata.FinishedAt,
		Size:       backup.Length,
		Failed:     backup.Metadata.Failed,
		Filename:   backup.Filename,
	}
	return backupDetails, nil
}

func GetBackupLogs(id, severity string) ([]models.BackupLog, error) {
	client := getClient()
	bucket := GetBackupsBucket(client)
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var backup models.Backup
	cursor, err := bucket.Find(bson.M{"_id": idObj})
	if err != nil {
		return []models.BackupLog{}, err
	}
	defer cursor.Close(nil)
	if cursor.Next(nil) {
		if err := cursor.Decode(&backup); err != nil {
			return []models.BackupLog{}, err
		}
	} else {
		return []models.BackupLog{}, mongo.ErrNoDocuments
	}
	if err := cursor.Err(); err != nil {
		return []models.BackupLog{}, err
	}

	severityHierarchy := map[string]int{
		"debug":   1,
		"info":    2,
		"warning": 3,
		"error":   4,
	}

	if severity != "" {
		var filteredLogs []models.BackupLog
		for _, log := range backup.Metadata.Logs {
			if severityHierarchy[log.Severity] >= severityHierarchy[severity] {
				filteredLogs = append(filteredLogs, log)
			}
		}
		return filteredLogs, nil
	}

	return backup.Metadata.Logs, nil
}
