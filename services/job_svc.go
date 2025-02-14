package services

import (
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
	if job, ok := jobList[id]; ok {
		return job, nil
	}
	return nil, fmt.Errorf("job %s not found", id)
}
