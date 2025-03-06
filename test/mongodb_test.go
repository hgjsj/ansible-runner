package test

import (
	"ansible-runner/configs"
	"ansible-runner/db"
	"ansible-runner/models"
	"fmt"
	"testing"
)

func init() {
	configs.LoadConfig("../")
}

var job_id string
func TestPersistJob(t *testing.T) {
	cli := db.NewMongoClient()

	job := models.NewAnsibleJob()
	job.Status = "scheduled"
	job.Step = "pre job before"
	job.Playbook = "pre_job_run.yml"
	if id, err := cli.Persist(cli.JobCollection, "", job); err != nil {
		t.Fatalf("insert doc to mongodb failed with err %s", err.Error())
	} else {
		t.Log(fmt.Sprintf("insert doc id %s to mongodb successful", job.ID.Hex()))
		t.Log(fmt.Sprintf("insert doc id %s to mongodb successful", id))
		job_id = job.ID.Hex()
	}
}

func TestRetrieveJob(t *testing.T) {
	cli := db.NewMongoClient()
	job := new(models.AnsibleJob)
	if err := cli.Retrieve(cli.JobCollection, job_id, job); err != nil {
		t.Fatalf("get doc to mongodb failed with err %s", err.Error())
	} else {
		if  job.Status == "scheduled" && job.Step == "pre job before" && job.Playbook == "pre_job_run.yml" {
			t.Log("get doc is the same as origin")
		} else {
			t.Fatal("get doc is not the same as origin")
		}
	}
}
