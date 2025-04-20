package database

import (
	"context"
	"fpgschiba.com/automation/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var backupJobsCollection *mongo.Collection

func GetBackupJobsCollection(client *mongo.Client) *mongo.Collection {
	if backupJobsCollection != nil {
		return backupJobsCollection
	}
	backupJobsCollection = client.Database(DatabaseName).Collection("backupJobs")
	return backupJobsCollection
}

func InsertBackupJob(job models.BackupJob) (string, error) {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	result, err := collection.InsertOne(context.Background(), job)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func ListBackupJobs() ([]models.BackupJob, error) {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var jobs []models.BackupJob
	err = cursor.All(context.Background(), &jobs)
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return []models.BackupJob{}, nil
	}
	return jobs, nil
}

func UpdateBackupJobSchedulerID(jobID string, schedulerID string) error {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"schedulerId": schedulerID}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetBackupJobFromID(jobID string) (models.BackupJob, error) {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return models.BackupJob{}, err
	}
	filter := bson.M{"_id": id}
	var job models.BackupJob
	err = collection.FindOne(context.Background(), filter).Decode(&job)
	if err != nil {
		return models.BackupJob{}, err
	}
	return job, nil
}

func DeleteBackupJob(jobID string) error {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBackupJob(jobID string, job models.BackupJob) error {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": job}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetJobNameFromID(jobID string) (string, error) {
	client := getClient()
	collection := GetBackupJobsCollection(client)
	id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return "", err
	}
	filter := bson.M{"_id": id}
	var job models.BackupJob
	err = collection.FindOne(context.Background(), filter).Decode(&job)
	if err != nil {
		return "", err
	}
	return job.Name, nil
}
