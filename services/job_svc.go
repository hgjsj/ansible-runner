package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetJob retrieves a job by its ID.
func GetJob(id string) (*models.AnsibleJob, error) {
    j , err := db.Retrieve[models.AnsibleJob](db.NewMongoClient(), db.NewMongoClient().JobCollection, id)
	if err != nil{
		return nil, fmt.Errorf("job %s not found", id)
	} else {
		return j, nil
	}
}



func SaveJob(job *models.AnsibleJob) error {
	return db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
}

func ListJobs(collection string) ([]models.AnsibleJob, error) {
	res, err := db.ListDocuments[models.AnsibleJob](db.NewMongoClient(),collection, bson.D{{"_id", 1}, {"playbook", 1}, {"status", 1}})
	if err != nil {
        if err == mongo.ErrEmptySlice{
			return make([]models.AnsibleJob, 0), nil
		}
		return nil, fmt.Errorf(err.Error())
	}
	return res, nil
}
