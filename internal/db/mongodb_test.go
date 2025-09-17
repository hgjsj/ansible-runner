package db

import (
	"ansible-runner/internal/utils"
	"ansible-runner/internal/models"
	"fmt"
	"testing"
)

func init() {
	utils.LoadConfig("../../configs/config")
	utils.InitLogger()
}

var job_id string

func TestPersistJob(t *testing.T) {
	cli := newMongoClient()

	job := models.NewAnsibleJob()
	job.Status = "scheduled"
	job.Step = "pre job before"
	job.Playbook = "pre_job_run.yml"
	if err := cli.CreateAnsibleJob(job); err != nil {
		t.Fatalf("insert doc to mongodb failed with err %s", err.Error())
	} else {
		t.Log(fmt.Sprintf("insert doc id %s to mongodb successful", job.ID))
		job_id = job.ID
	}
}

func TestRetrieveJob(t *testing.T) {
	cli := newMongoClient()
	job, err := cli.GetAnsibleJob(job_id)
	if err != nil {
		t.Fatalf("get doc to mongodb failed with err %s", err.Error())
		return
	}

	t.Log(fmt.Sprintf("get doc id %s from db", job.ID))
	if job.Status == "scheduled" && job.Step == "pre job before" && job.Playbook == "pre_job_run.yml" {
		t.Log("get doc is the same as origin")
	} else {
		t.Fatal("get doc is not the same as origin")
	}

}
