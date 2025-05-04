package tasks

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"encoding/json"
	"fmt"
	"time"
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
	
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password, PoolSize: 20})
	//asynqClient = asynq.NewClientFromRedisClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password, PoolSize: 20}.MakeRedisClient().(redis.UniversalClient))

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
	if ti, err := newAsynqClient().Enqueue(task, asynq.MaxRetry(0), asynq.Timeout(time.Minute),asynq.Deadline(time.Now().Add(time.Minute))); err != nil {
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

// EnqueueJobStatusTask represents a job status task.
func EnqueueJobStatusTask(job *models.AnsibleJob) error {
	utils.Log.Info("Enqueue job status task to queue")
	jst, err := json.Marshal(job)
	if err != nil {
		utils.Log.Error(err.Error())
		return  err
	}
	task := asynq.NewTask("jobstatusupdate", jst)
	if ti, err := newAsynqClient().Enqueue(task, asynq.MaxRetry(0),asynq.Timeout(time.Minute),asynq.Deadline(time.Now().Add(time.Minute))); err != nil {
		utils.Log.Error(err.Error())
		return  err
	} else {
		utils.Log.Info("Enqueue job status task to queue with id %s and state %s", ti.ID, ti.State)
		return nil
	}
}
