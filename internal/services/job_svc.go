package services

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/models"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetJob retrieves a job by its ID.
func GetJob(id string) (*models.AnsibleJob, error) {
    j , err := db.GetDBClient().GetAnsibleJob(id)
	if err != nil{
		return nil, fmt.Errorf("job %s not found", id)
	} 
	return j, nil
	
}

// SaveJob saves a job.
func SaveJob(job *models.AnsibleJob) error {
	return db.GetDBClient().CreateAnsibleJob(job)
}

// ListJobs retrieves all jobs.
func ListJobs() ([]models.AnsibleJob, error) {
	res, err := db.GetDBClient().ListAnsibleJob()
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
		id = ""
	} 
	//delete job by id
	return  db.GetDBClient().DeleteAnsibleJob(id)
	
}	

// UpdateJobStatus updates a job status by its ID.	
func UpdateJobStatus(job *models.AnsibleJob) error {
	j, err := GetJob(job.ID)
	if err != nil {
		return err
	}

	if len(job.Status) > 0{
		j.Status = job.Status
	}
	
	if len(job.Step) > 0{
		j.Step = job.Step
	}
	if len(job.Error) > 0{
		j.Error = job.Error
	}

	if len(job.Output) > 0{
		j.Output = job.Output
	}

	if time.Time.IsZero(job.StartTime) == false{
		j.StartTime = job.StartTime
	}
	
	if time.Time.IsZero(job.EndTime) == false{
		j.EndTime = job.EndTime
	}
	return SaveJob(j)
}
