package tasks

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/services"
	"ansible-runner/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

var asynqServer *asynq.Server
var mux *asynq.ServeMux 


func newAsyncServer() *asynq.Server {
	if asynqServer != nil {
		return asynqServer
	}
	redisAddr := viper.GetString("redis.uri")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		return nil
	}
	asynqServer = asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr, Password: password}, 
		asynq.Config{Concurrency: 5, ShutdownTimeout: 10 * time.Second})
	mux = asynq.NewServeMux()
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
	_, err := services.RunPlaybook(j.Job, j.Inventory)
	if err != nil {
		return err
	}
	return nil
}

// LaunchServer starts the asynq server and handles the job task.
func LaunchServer(srvType string) error {
	if newAsyncServer() == nil  {
		utils.Log.Error("can not create asynq server")
		return fmt.Errorf("can not create asynq server")
	}

	if srvType == "server" {
		utils.Log.Info("ansible-runner serer started with HandleAnsibleJobStatusTask")
		mux.HandleFunc("jobstatusupdate", HandleAnsibleJobStatusTask)
	} else if srvType == "worker" {
		utils.Log.Info("ansible-runner worker started with HandleAnsibleJobTask")
		mux.HandleFunc("ansiblejob", HandleAnsibleJobTask)
	} else {
		utils.Log.Error("Not supported server type, choose server or worker")
		return fmt.Errorf("Not supported server type, choose server or worker")
	}
	if err := newAsyncServer().Run(mux); err != nil {
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
	utils.Log.Info("Unmarsh job status task done")
	err := services.UpdateJobStatus(&js)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the MongoDB connection.
func Close() {
	// close Sever connection
    if  asynqServer != nil {
	    asynqServer.Shutdown()
    }
}

