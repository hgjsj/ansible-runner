package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"fmt"
)

// Jobs is a map of job IDs to AnsibleJob objects.
var jobList map[string]*models.AnsibleJob

func init() {
	jobList = make(map[string]*models.AnsibleJob)
}

// GetJob retrieves a job by its ID.
func GetJob(id string) (*models.AnsibleJob, error) {
	if job, err := db.MongoClient.RetrieveJob(id); err != nil{
		return nil, fmt.Errorf("job %s not found", id)
	} else {
		return job, nil
	}
}
