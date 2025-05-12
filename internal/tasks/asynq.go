package tasks

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/services"
	"ansible-runner/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"github.com/redis/go-redis/v9"
)

var asynqServer *asynq.Server
//var mux *asynq.ServeMux 


func newAsyncServer(queueName string) *asynq.Server {
	if asynqServer != nil {
		return asynqServer
	}
	redisAddr := viper.GetString("redis.uri")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		return nil
	}

	//asynqServer = asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr, Password: password, ReadTimeout: -1, WriteTimeout: -1, PoolSize: 100}, asynq.Config{Concurrency: 20, LogLevel: asynq.InfoLevel})
	asynqServer = asynq.NewServerFromRedisClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password, PoolSize: 100}.MakeRedisClient().(redis.UniversalClient), 
		asynq.Config{Concurrency: 100, LogLevel: asynq.InfoLevel, Queues: map[string]int{queueName: 10}})
	// ErrorHandler: func(ctx context.Context, t *asynq.Task, err error) {})
	//mux = asynq.NewServeMux()
	return asynqServer
}


// HandleAnsibleJobTask is a handler for the job task.
func HandleAnsibleJobTask(ctx context.Context, t *asynq.Task) error {
	var j JobTask
	utils.Log.Info("Unmarsh job task start")
	if err := json.Unmarshal(t.Payload(), &j); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	utils.Log.Info("Unmarsh job task done")
	_, err := services.RunPlaybook(j.Job, j.Inventory, services.UpdateJobStatus)
	if err != nil {
		return err
	}
	return nil
}

// LaunchServer starts the asynq server and handles the job task.
func LaunchServer(srvType string) error {

	var queueName string
	mux := asynq.NewServeMux()
	if srvType == "server" {
		utils.Log.Info("ansible-runner serer started with HandleAnsibleJobStatusTask for queue jobstatusupdate")
		mux.HandleFunc("jobstatusupdate", HandleAnsibleJobStatusTask)
		queueName = "jobstatusupdate"
	} else if srvType == "worker" {
		utils.Log.Info("ansible-runner worker started with HandleAnsibleJobTask for queue ansiblejob")
		mux.HandleFunc("ansiblejob", HandleAnsibleJobTask)
		queueName = "ansiblejob"
	} else {
		utils.Log.Error("Not supported server type, choose server or worker")
		return fmt.Errorf("Not supported server type, choose server or worker")
	}

	if err := newAsyncServer(queueName).Run(mux); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

// HandleAnsibleJobStatusTask is a handler for the job status task.
func HandleAnsibleJobStatusTask(ctx context.Context, t *asynq.Task) error {
	var js models.AnsibleJob
	utils.Log.Info("Unmarsh job status task start")
	if err := json.Unmarshal(t.Payload(), &js); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	go func ()  {
		utils.Log.Info("Unmarsh job status task done")
		err := services.UpdateJobStatus(&js)
		if err != nil {
			utils.Log.Error(err.Error())
		}
		utils.Log.Info(fmt.Sprintf("Update job %s status done", js.ID.Hex()))
	}()

	return nil
}

// Close closes the MongoDB connection.
func Close() {
	// close Sever connection
    if  asynqServer != nil {
	    asynqServer.Shutdown()
    }
}

