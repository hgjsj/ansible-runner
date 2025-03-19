package tasks

import (
	"ansible-runner/internal/services"
	"ansible-runner/internal/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

// type AsynqServer struct {
// 	s *asynq.Server
// 	m *asynq.ServeMux
// }

var asynqServer *asynq.Server
var mux *asynq.ServeMux 


func newAsyncServer() *asynq.Server {
	if asynqServer != nil {
		return asynqServer
	}
	redisAddr := viper.GetString("redis.addr")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		return nil
	}
	asynqServer = asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr, Password: password}, asynq.Config{Concurrency: 5})
	mux = asynq.NewServeMux()
	return asynqServer
}

type JobTask struct {
	ID          string           `json:"id,omitempty"`
	InventoryID string           `json:"inventoryid,omitempty"`
}



func HandleAnsibleJobTask(ctx context.Context, t *asynq.Task) error {
	var p JobTask
	utils.Log.Info("Unmarsh job task start")
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	utils.Log.Info("Unmarsh job task done")
	_, err := services.RunPlaybook(p.ID, p.InventoryID)
	if err != nil {
		return err
	}
	return nil
}


func LaunchServer() error {
	if newAsyncServer() == nil  {
		utils.Log.Error("can not create asynq server")
		return fmt.Errorf("can not create asynq server")
	}
	mux.HandleFunc("ansiblejob", HandleAnsibleJobTask)
	if err := newAsyncServer().Run(mux); err != nil {
		utils.Log.Error(err.Error())
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

