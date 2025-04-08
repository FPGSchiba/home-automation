package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"os"
)

var backupsCollection *mongo.Collection
var backupsGridFS *gridfs.Bucket

func GetBackupsCollection(client *mongo.Client) *mongo.Collection {
	if backupsCollection != nil {
		return backupsCollection
	}
	backupsCollection = client.Database(DatabaseName).Collection("backups")
	return backupsCollection
}

func GetBackupsBucket(client *mongo.Client) *gridfs.Bucket {
	if backupsGridFS != nil {
		return backupsGridFS
	}

	backupsGridFS, _ = gridfs.NewBucket(client.Database(DatabaseName))

	return backupsGridFS
}

func UploadBackup(filename string, filePath string, meta bson.M) error {
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
