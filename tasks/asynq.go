package tasks

import (
	"ansible-runner/services"
	"ansible-runner/utils"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

type AsynqClient struct {
	*asynq.Client
}

type AsynqServer struct {
	s *asynq.Server
	m *asynq.ServeMux
}

var asynqClient *AsynqClient
var asynqServer *AsynqServer

func NewAsynqClient() *AsynqClient {
	if asynqClient != nil {
		return asynqClient
	}
	redisAddr := viper.GetString("redis.addr")
	password, err := utils.Decrypt(viper.GetString("redis.password"))
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	asynqClient = &AsynqClient{asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: password})}
	// set logger to client
	return asynqClient
}

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

func (c *AsynqClient) EnqueueJobTask(jobid string, inventoryid string) (string, error) {
	jobTask := JobTask{ InventoryID: inventoryid, ID: jobid}
	b, err := json.Marshal(jobTask)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	task := asynq.NewTask("ansiblejob", b)
	if _, err := c.Enqueue(task, asynq.MaxRetry(0)); err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	return jobTask.ID, nil
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
    if s != nil && s.s != nil {
	    s.s.Shutdown()
    }
}

func (c *AsynqClient) Close(){
	if c != nil {
	    c.Close()
    }
}