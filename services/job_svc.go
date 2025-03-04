package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"fmt"
)

// GetJob retrieves a job by its ID.
func GetJob(id string) (*models.AnsibleJob, error) {
    j := &models.AnsibleJob{}
	if err := db.NewMongoClient().Retrieve(db.MongoClient.JobCollection, id, j); err != nil{
		return nil, fmt.Errorf("job %s not found", id)
	} else {
		return j, nil
	}
}

func GetJobInventory(id string) (*models.Inventory, error) {
	i := &models.Inventory{}
	if err := db.NewMongoClient().Retrieve(db.NewMongoClient().InventoryCollection, id, i); err != nil{
		return nil, fmt.Errorf("inventory %s not found", id)
	} else {
		return i, nil
	}
}

func SaveJob(job *models.AnsibleJob) error {
	return db.NewMongoClient().Persist(db.NewMongoClient().JobCollection, job.ID, job)
}
