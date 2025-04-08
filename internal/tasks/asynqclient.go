package tasks

import (
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
func EnqueueJobTask(jobid string, inventoryid string, outputplugin string) (string, error) {
	jobTask := JobTask{InventoryID: inventoryid, ID: jobid, OutputPlugin: outputplugin}
	b, err := json.Marshal(jobTask)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	task := asynq.NewTask("ansiblejob", b)
	if _, err := newAsynqClient().Enqueue(task, asynq.MaxRetry(0)); err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	return jobTask.ID, nil
}

// CloseClient closes the asynq client.
func CloseClient() {
	if asynqClient != nil {
		asynqClient.Close()
	}
}
