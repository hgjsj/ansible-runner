package services

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"ansible-runner/utils"
	"fmt"
    "go.mongodb.org/mongo-driver/bson"
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
	res, err := db.ListDocuments[models.Inventory](db.NewMongoClient(), collection, bson.D{{"id", 1}, {"name", 1}})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return res, nil
}

func CreateInventory(i *models.Inventory) (string, error) {
	err := db.NewMongoClient().Persist(db.NewMongoClient().InventoryCollection, i.ID.Hex(), i)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}

	return i.ID.Hex(), nil
}
