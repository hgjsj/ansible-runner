package services

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/models"

	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetJob retrieves a job by its ID.
func GetJob(id string) (*models.AnsibleJob, error) {
    j , err := db.Retrieve[models.AnsibleJob](db.NewMongoClient(), db.NewMongoClient().JobCollection, id)
	if err != nil{
		return nil, fmt.Errorf("job %s not found", id)
	} 
	return j, nil
	
}


// SaveJob saves a job.
func SaveJob(job *models.AnsibleJob) error {
	return db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID.Hex(), job)
}

// ListJobs retrieves all jobs.
func ListJobs(collection string) ([]models.AnsibleJob, error) {
	res, err := db.ListDocuments[models.AnsibleJob](db.NewMongoClient(),collection, bson.D{{Key: "_id", Value: 1}, {Key: "playbook", Value: 1}, {Key: "status", Value: 1}})
	if err != nil {
        if err == mongo.ErrEmptySlice{
			return make([]models.AnsibleJob, 0), nil
		}
		return nil, fmt.Errorf(err.Error())
	}
	return res, nil
}

// DeleteJob deletes a job by its ID.
func DeleteJob(id string) error {
	if id == "all" {
		//delete all jobs
		return db.DeleteDocument(db.NewMongoClient(), db.NewMongoClient().JobCollection, "")
	} 
	//delete job by id
	return db.DeleteDocument(db.NewMongoClient(), db.NewMongoClient().JobCollection, id)
	
}	

