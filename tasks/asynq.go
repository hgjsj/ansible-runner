package tasks

import (
	"ansible-runner/models"
	"ansible-runner/services"
	"ansible-runner/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

	asynqServer = &AsynqServer{s: asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr, Password: password}, asynq.Config{}), m: asynq.NewServeMux()}
	return asynqServer
}

type JobTask struct {
	ID        string           `json:"id,omitempty"`
	Playbook  string           `json:"playbook"`
	Inventory models.Inventory `json:"inventory"`
}

func (c *AsynqClient) EnqueueJobTask(p string, i *models.Inventory) (string, error) {
	jobTask := JobTask{Playbook: p, Inventory: *i, ID: uuid.New().String()}
	b, err := json.Marshal(jobTask)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	task := asynq.NewTask("ansiblejob", b)
	if _, err := c.Enqueue(task); err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}
	return jobTask.ID, nil
}

func HandleAnsibleJobTask(ctx context.Context, t *asynq.Task) error {
	var p JobTask

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	job, err := services.RunPlaybook(p.Playbook, p.Inventory)
	if err != nil {
		utils.Log.Error(err.Error())
	}
	utils.Log.Info(fmt.Sprintf("Job %s starts to run", job.ID))
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
