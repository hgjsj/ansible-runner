package tasks

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"encoding/json"
	"fmt"

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
	utils.Log.Info("Enqueue job task to queue")
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
	if ti, err := newAsynqClient().Enqueue(task, asynq.MaxRetry(0)); err != nil {
		utils.Log.Error(err.Error())
		return "", err
	} else {
		utils.Log.Info(fmt.Sprintf("Enqueue job task to queue with id %s and state %s", ti.ID, ti.State))
		return job.ID.Hex(), nil
	}
	
}

// CloseClient closes the asynq client.
func CloseClient() {
	if asynqClient != nil {
		asynqClient.Close()
	}
}

