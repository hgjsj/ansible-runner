package tasks

import (
	"ansible-runner/services"
	"ansible-runner/utils"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)


type AsynqServer struct {
	s *asynq.Server
	m *asynq.ServeMux
}

var asynqServer *AsynqServer

func NewAsyncServer() *AsynqServer {
	if asynqServer != nil {
		return asynqServer
	}
	redisAddr := viper.GetString("redis.addr")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	asynqServer = &AsynqServer{s: asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr, Password: password}, asynq.Config{Concurrency: 5}), m: asynq.NewServeMux()}
	return asynqServer
}

type JobTask struct {
	ID        string           `json:"id,omitempty"`
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
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func (s *AsynqServer) RegisterHandlers() {
	s.m.HandleFunc("ansiblejob", HandleAnsibleJobTask)
}

func (s *AsynqServer) LaunchServer() {
	if err := s.s.Run(s.m); err != nil {
		utils.Log.Error(err.Error())
	}
}

// Close closes the MongoDB connection.
func (s *AsynqServer) Close() {
	// close client connection
	
    if s != nil && s.s != nil {
	    s.s.Shutdown()
    }
}

