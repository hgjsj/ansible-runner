package services

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetInventory(id string) (*models.Inventory, error) {
	i, err := db.Retrieve[models.Inventory](db.NewMongoClient(), db.NewMongoClient().InventoryCollection, id)
	if err != nil {
		return nil, fmt.Errorf("inventory %s not found", id)
	} else {
		return i, nil
	}
}

func ListInventory(collection string) ([]models.Inventory, error) {
	res, err := db.ListDocuments[models.Inventory](db.NewMongoClient(), collection, bson.D{{"_id", 1}, {"name", 1}})
	if err != nil {
		if err == mongo.ErrEmptySlice{
			return make([]models.Inventory, 0), nil
		}
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
