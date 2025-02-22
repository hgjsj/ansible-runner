package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"fmt"
)

// GetJob retrieves a job by its ID.
func GetJob(id string) (*models.AnsibleJob, error) {
	if job, err := db.NewMongoClient().RetrieveJob(id); err != nil{
		return nil, fmt.Errorf("job %s not found", id)
	} else {
		return job, nil
	}
}

func GetJobInventory(id string) (*models.Inventory, error) {
	if i, err := db.NewMongoClient().Retrieve(db.NewMongoClient().InventoryCollection, id); err != nil{
		return nil, fmt.Errorf("inventory %s not found", id)
	} else {
		return i.(*models.Inventory), nil
	}
}
