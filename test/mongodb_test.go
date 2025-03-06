package test

import (
	"ansible-runner/configs"
	"ansible-runner/db"
	"ansible-runner/models"
	"ansible-runner/utils"
	"fmt"
	"testing"
)

func init() {
	configs.LoadConfig("../")
	utils.InitLogger("../logs/test.log")
}

var job_id string
func TestPersistJob(t *testing.T) {
	cli := db.NewMongoClient()

	job := models.NewAnsibleJob(true)
	job.Status = "scheduled"
	job.Step = "pre job before"
	job.Playbook = "pre_job_run.yml"
	if err := cli.Persist(cli.JobCollection, job.ID.Hex(), job); err != nil {
		t.Fatalf("insert doc to mongodb failed with err %s", err.Error())
	} else {
		t.Log(fmt.Sprintf("insert doc id %s to mongodb successful", job.ID.Hex()))
		job_id = job.ID.Hex()
	}
}

func TestRetrieveJob(t *testing.T) {
	cli := db.NewMongoClient()
	job, err := db.Retrieve[models.AnsibleJob](cli, cli.JobCollection, job_id)
	if err != nil {
		t.Fatalf("get doc to mongodb failed with err %s", err.Error())
		return
	} 

	t.Log(fmt.Sprintf("get doc id %s from db", job.ID.Hex()))
	if  job.Status == "scheduled" && job.Step == "pre job before" && job.Playbook == "pre_job_run.yml" {
		t.Log("get doc is the same as origin")
	} else {
		t.Fatal("get doc is not the same as origin")
	}
	
}
