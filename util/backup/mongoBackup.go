package backup

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MongoInput struct {
	Host         string `json:"host" bson:"host"`
	DatabaseName string `json:"databaseName" bson:"databaseName"`
	Username     string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
}

func CreateMongoBackupJob(input MongoInput, schedule string) (uuid.UUID, error) {
	scheduler = getScheduler()
	job, err := scheduler.NewJob(gocron.CronJob(schedule, false), gocron.NewTask(dummy))
	if err != nil {
		return uuid.UUID{}, err
	}
	_, err = scheduler.Update(job.ID(), gocron.CronJob(schedule, false), gocron.NewTask(runMongoBackup, job.ID(), input))
	if err != nil {
		return uuid.UUID{}, err
	}
	scheduler.Start() // Need to start the scheduler to run the job
	return job.ID(), nil
}

func runMongoBackup(jobID uuid.UUID, input MongoInput) {
	log.WithFields(log.Fields{
		"component": "backup",
		"func":      "runMongoBackup",
		"jobID":     jobID,
		"input":     input,
	}).Info("Running MongoDB backup job")
}
