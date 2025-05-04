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

//var asynqClient *asynq.Client

func newAsynqClient() *asynq.Client {
	// if asynqClient != nil {
	// 	return asynqClient
	// }
	redisAddr := viper.GetString("redis.uri")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}
	
    return asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password, ReadTimeout: -1, WriteTimeout: -1, PoolSize: 100})
	//asynqClient = asynq.NewClientFromRedisClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password, ReadTimeout: -1, WriteTimeout: -1, PoolSize: 100}.MakeRedisClient().(redis.UniversalClient))

	// set logger to client
	//return asynqClient
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
	c := newAsynqClient()
	ti, err := c.Enqueue(task, asynq.MaxRetry(0), asynq.Timeout(3 * time.Minute), asynq.Deadline(time.Now().Add(3 * time.Minute)))
	if err != nil {
		c.Close()
		utils.Log.Error(err.Error())
		return "", err
	}
	c.Close()
	utils.Log.Info(fmt.Sprintf("Enqueue job task to queue with id %s and state %s", ti.ID, ti.State))
	return job.ID.Hex(), nil
	
}

// CloseClient closes the asynq client.
func CloseClient() {
	// if asynqClient != nil {
	// 	asynqClient.Close()
	// }
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
	c := newAsynqClient()
	ti, err := c.Enqueue(task, asynq.MaxRetry(0), asynq.Timeout(3*time.Minute), asynq.Deadline(time.Now().Add(3*time.Minute)))
	if err != nil {
		c.Close()
		utils.Log.Error(err.Error())
		return err
	}
	c.Close()
	utils.Log.Info("Enqueue job status task to queue with id %s and state %s", ti.ID, ti.State)
	return nil
	
}
