package services

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetInventory retrieves an inventory by its ID.
func GetInventory(id string) (*models.Inventory, error) {
	i, err := db.Retrieve[models.Inventory](db.NewMongoClient(), db.NewMongoClient().InventoryCollection, id)
	if err != nil {
		return nil, fmt.Errorf("inventory %s not found", id)
	} 
	return i, nil
	
}

// ListInventory retrieves all inventories.
func ListInventory(collection string) ([]models.Inventory, error) {
	res, err := db.ListDocuments[models.Inventory](db.NewMongoClient(), collection, bson.D{{Key: "_id", Value: 1}, {Key: "name", Value: 1}})
	if err != nil {
		if err == mongo.ErrEmptySlice{
			return make([]models.Inventory, 0), nil
		}
		return nil, fmt.Errorf(err.Error())
	}
	return res, nil
}

// CreateInventory creates a new inventory.
func CreateInventory(i *models.Inventory) (string, error) {
	err := db.NewMongoClient().Persist(db.NewMongoClient().InventoryCollection, i.ID.Hex(), i)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}

	return i.ID.Hex(), nil
}

// DeleteInventory deletes an inventory by its ID.
func DeleteInventory(id string) error {
	if id == "all" {
		//delete all inventories
		return db.DeleteDocument(db.NewMongoClient(), db.NewMongoClient().InventoryCollection, "")
	} 
	//delete inventory by id
	return db.DeleteDocument(db.NewMongoClient(), db.NewMongoClient().InventoryCollection, id)
}	
