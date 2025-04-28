package tasks

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

var asynqClient *asynq.Client

func newAsynqClient() *asynq.Client {
	if asynqClient != nil {
		return asynqClient
	}
	redisAddr := viper.GetString("redis.uri")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	asynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password})
	// set logger to client
	return asynqClient
}

// EnqueueJobTask represents a job task.
func EnqueueJobTask(job *models.AnsibleJob, inventory *models.Inventory) (string, error) {
	JobTask := JobTask{
		Job:       job,
		Inventory: inventory,
	}

	jt, err := json.Marshal(JobTask)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}

	task := asynq.NewTask("ansiblejob", jt)
	if _, err := newAsynqClient().Enqueue(task, asynq.MaxRetry(0)); err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	return job.ID.Hex(), nil
}

// CloseClient closes the asynq client.
func CloseClient() {
	if asynqClient != nil {
		asynqClient.Close()
	}
}

// EnqueueJobStatusTask represents a job status task.
func EnqueueJobStatusTask(job *models.AnsibleJob) error {
	jst, err := json.Marshal(job)
	if err != nil {
		utils.Log.Error(err.Error())
		return  err
	}
	task := asynq.NewTask("ansiblejobstatus", jst)
	if _, err := newAsynqClient().Enqueue(task, asynq.MaxRetry(0)); err != nil {
		utils.Log.Error(err.Error())
		return  err
	}
	return nil
}
