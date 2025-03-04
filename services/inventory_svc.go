package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"ansible-runner/utils"
	"fmt"

	"github.com/google/uuid"
)

func GetInventory(id string) (*models.Inventory, error) {
	i := &models.Inventory{}
	if err := db.NewMongoClient().Retrieve(db.NewMongoClient().InventoryCollection, id, i); err != nil {
		return nil, fmt.Errorf("inventory %s not found", id)
	} else {
		return i, nil
	}
}

func ListInventory(collection string) ([]models.Inventory, error) {
	res, err := db.NewMongoClient().ListAllDocuments(collection)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return res, nil
}

func CreateInventory(i *models.Inventory) error {
	i.ID = uuid.New().String()

	err := db.NewMongoClient().Persist(db.NewMongoClient().InventoryCollection, i.ID, i)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}

	return nil
}
