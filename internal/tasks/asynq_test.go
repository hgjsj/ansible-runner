package tasks

import (
	"ansible-runner/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hibiken/asynq"
)

type TestJobTask struct {
	Job       string `json:"job,omitempty"`
	Inventory string `json:"inventory,omitempty"`
}

type TestUpdateJobTask struct {
	Job   string `json:"job,omitempty"`
	State string `json:"state,omitempty"`
}

func init() {
	utils.LoadConfig("../../configs/config")
	utils.InitLogger()
}

func handleTestjobTask(ctx context.Context, t *asynq.Task) error {
	// Simulate task processing
	var j TestJobTask
	utils.Log.Info("Unmarsh Test job task start")
	if err := json.Unmarshal(t.Payload(), &j); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	utils.Log.Info("Processing testjob task")
	time.Sleep(3 * time.Second) // Simulate some work
	TestUpdateJobTask := TestUpdateJobTask{Job: j.Job, State: "running"}
	b, err := json.Marshal(TestUpdateJobTask)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}

	go func() {
		newAsynqClient().Enqueue(asynq.NewTask("testupdatejob", b), asynq.MaxRetry(0), asynq.Timeout(3*time.Minute), asynq.Deadline(time.Now().Add(3*time.Minute)))

		time.Sleep(3 * time.Second) // Simulate some work
		TestUpdateJobTask.State = "completed"
		b, err = json.Marshal(TestUpdateJobTask)
		if err != nil {
			utils.Log.Error(err.Error())
		}
		newAsynqClient().Enqueue(asynq.NewTask("testupdatejob", b), asynq.MaxRetry(0), asynq.Timeout(3*time.Minute), asynq.Deadline(time.Now().Add(3*time.Minute)))
		utils.Log.Info("Enqueue testupdatejob task")
	}()
	return nil
}

func handleTestupdatejobTask(ctx context.Context, t *asynq.Task) error {
	var uj TestUpdateJobTask
	utils.Log.Info("Unmarsh Test update job task start")
	if err := json.Unmarshal(t.Payload(), &uj); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	go func() {
		time.Sleep(time.Second) // Simulate some work
		utils.Log.Info("Processing testupdatejob task")
		fmt.Printf("Job: %s, State: %s\n", uj.Job, uj.State)
	}()

	return nil
}

func TestAsynqWorker(t *testing.T) {
	s := newAsyncServer("default")
	mux.HandleFunc("testjob", handleTestjobTask)
	go s.Run(mux)

}

func TestAsynqServer(t *testing.T) {
	s := newAsyncServer("default")
	mux.HandleFunc("testupdatejob", handleTestupdatejobTask)
	go s.Run(mux)

	for i := range 5 {
		task := TestJobTask{Job: fmt.Sprintf("testjob-%d", i), Inventory: "testinventory"}
		b, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("failed to marshal task: %v", err)
		}
		t.Logf("Enqueue task with payload: %s", string(b))
		newAsynqClient().Enqueue(asynq.NewTask("testjob", b), asynq.MaxRetry(0), asynq.Timeout(3*time.Minute), asynq.Deadline(time.Now().Add(3*time.Minute)))
		time.Sleep(5 * time.Second) // Simulate some work
	}
	s.Shutdown()
}
