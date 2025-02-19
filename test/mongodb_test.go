package test

import (
	"ansible-runner/configs"
	"ansible-runner/db"
	"ansible-runner/models"
	"testing"
)

func init() {
	configs.LoadConfig("../")
}

func TestPersistJob(t *testing.T) {
	cli := db.NewMongoClient()

	job := new(models.AnsibleJob)
	job.ID = "unit-test"
	job.Status = "scheduled"
	job.Step = "pre job before"
	job.Playbook = "pre_job_run.yml"
	if err := cli.PersistJob(job); err != nil {
		t.Fatalf("insert doc to mongodb failed with err %s", err.Error())
	} else {
		t.Log("insert doc to mongodb successful")
	}
}

func TestRetrieveJob(t *testing.T) {
	cli := db.NewMongoClient()
	if job, err := cli.RetrieveJob("unit-test"); err != nil {
		t.Fatalf("get doc to mongodb failed with err %s", err.Error())
	} else {
		if job.ID == "unit-test" && job.Status == "scheduled" && job.Step == "pre job before" && job.Playbook == "pre_job_run.yml" {
			t.Log("get doc is the same as origin")
		} else {
			t.Fatal("get doc is not the same as origin")
		}
	}
}
