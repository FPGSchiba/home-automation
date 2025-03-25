package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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
