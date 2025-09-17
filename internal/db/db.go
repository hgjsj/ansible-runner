package db

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"github.com/spf13/viper"
)

type DBService interface {
	GetInventory(id string) (*models.Inventory, error)
	GetAnsibleJob(id string) (*models.AnsibleJob, error)
	ListInventory() ([]models.Inventory, error)
	ListAnsibleJob() ([]models.AnsibleJob, error)
	CreateInventory(i *models.Inventory) error
	CreateAnsibleJob(j *models.AnsibleJob) error
	DeleteInventory(id string) error
	DeleteAnsibleJob(id string) error
	CloseService()
}


func GetDBClient() DBService {
	if viper.GetStringMap("postgresql") != nil {
		return newPostgreSQL()	
	}

	if viper.GetStringMap("mongodb") != nil {
		return newMongoClient()
	}

	utils.Log.Error("no database configuration found")
	panic("no database configuration found")
}

