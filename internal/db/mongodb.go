package db

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/utils"
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnsibleJobMongo struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Playbook     string             `bson:"playbook,omitempty"`
	Status       string             `bson:"status,omitempty"`
	Output       string             `bson:"output,omitempty"`
	Error        string             `bson:"error,omitempty"`
	Step         string             `bson:"step,omitempty"`
	StartTime    primitive.DateTime `bson:"start_time,omitempty"`
	EndTime      primitive.DateTime `bson:"end_time,omitempty"`
	InventoryID  primitive.ObjectID `bson:"inventory_id,omitempty"`
	OutputPlugin string             `bson:"output_plugin,omitempty"`
}

type inventoryMongo struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Localhost	*models.Host              `bson:"localhost,omitempty"`
	Groups     []models.Group            `bson:"groups,omitempty"`
}

// MongoDB represents a MongoDB configuration.
type MongoDB struct {
	mongoClient         *mongo.Client
	db                  string
	JobCollection       string
	InventoryCollection string
	c                   context.Context
	DBService
}

// NewMongoDB creates a new MongoDB configuration.
var MongoClient *MongoDB

func newMongoClient() *MongoDB {
	if MongoClient != nil {
		return MongoClient
	}
	MongoClient = new(MongoDB)
	MongoClient.c = context.Background()
	var err error
	uri := viper.GetString("mongodb.uri")
	MongoClient.db = viper.GetString("mongodb.database")
	MongoClient.JobCollection = viper.GetString("mongodb.jobCollection")
	MongoClient.InventoryCollection = viper.GetString("mongodb.inventoryCollection")

	username := viper.GetString("mongodb.user")
	password := viper.GetString("mongodb.password")
	password, err = utils.Decrypt(password)
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	credential := options.Credential{
		Username:      username,
		Password:      password,
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    MongoClient.db,
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)

	MongoClient.mongoClient, err = mongo.Connect(MongoClient.c, clientOptions)
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	databaseNames, err := MongoClient.mongoClient.ListDatabaseNames(MongoClient.c, bson.M{})
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	dbExists := false
	for _, dbName := range databaseNames {
		if dbName == MongoClient.db {
			dbExists = true
			break
		}
	}

	if !dbExists {
		err = MongoClient.mongoClient.Database(MongoClient.db).CreateCollection(MongoClient.c, MongoClient.JobCollection)
		if err != nil {
			utils.Log.Error(err.Error())
			panic(err)
		}
	}
	collectionNames, err := MongoClient.mongoClient.Database(MongoClient.db).ListCollectionNames(MongoClient.c, bson.M{})
	if err != nil {
		utils.Log.Error(err.Error())
		panic(err)
	}

	jobCollectionExists := false
	inventoryCollectionExists := false
	for _, collName := range collectionNames {
		if collName == MongoClient.JobCollection {
			jobCollectionExists = true
		}
		if collName == MongoClient.InventoryCollection {
			inventoryCollectionExists = true
		}
	}

	if !jobCollectionExists {
		err = MongoClient.mongoClient.Database(MongoClient.db).CreateCollection(MongoClient.c, MongoClient.JobCollection)
		if err != nil {
			utils.Log.Error(err.Error())
			panic(err)
		}
	}

	if !inventoryCollectionExists {
		err = MongoClient.mongoClient.Database(MongoClient.db).CreateCollection(MongoClient.c, MongoClient.InventoryCollection)
		if err != nil {
			utils.Log.Error(err.Error())
			panic(err)
		}
	}

	return MongoClient
}

func (m *MongoDB) persist(collection string, keyFilter string, value any) error {

	objID, err := primitive.ObjectIDFromHex(keyFilter)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}

	update := bson.M{
		"$set": value,
	}
	updateOptions := options.Update().SetUpsert(true)
	filter := bson.M{"_id": objID}
	_, err = m.mongoClient.Database(m.db).Collection(collection).UpdateOne(m.c, filter, update, updateOptions)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDocument deletes a document from a collection.
func (m *MongoDB) deleteDocument(collection string, keyFilter string) error {
	filter := bson.M{}
	if keyFilter != "" {

		objID, err := primitive.ObjectIDFromHex(keyFilter)

		if err != nil {
			utils.Log.Error(err.Error())
			return err
		}
		filter = bson.M{"_id": objID}
	} 

	_, err := m.mongoClient.Database(m.db).Collection(collection).DeleteMany(m.c, filter)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

// Close closes the MongoDB connection.
func (m *MongoDB) close() {
	if m != nil && m.mongoClient != nil {
		m.mongoClient.Disconnect(m.c)
	}
}

func (m *MongoDB) CloseService() {
	m.close()
}

func (m *MongoDB) CreateInventory(i *models.Inventory) error {
	var iID primitive.ObjectID
	var err error
	if i.ID == "" {
		iID = primitive.NewObjectID()
	} else {
		iID, err = primitive.ObjectIDFromHex(i.ID)
		if err != nil {
			utils.Log.Error(err.Error())
			return err
		}
	}

	iMongo := &inventoryMongo{
		ID:	  iID,
		Name:    i.Name,
		Groups:  i.Groups,
		Localhost: i.Localhost,
	}
	if  err := m.persist(m.InventoryCollection, iMongo.ID.Hex(), iMongo); err != nil {
		return err
	}
	i.ID = iMongo.ID.Hex()
	return nil
}


func (m *MongoDB) GetInventory(id string) (*models.Inventory, error) {
	var i models.Inventory
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	filter := bson.M{"_id": objID}
	iMongo := inventoryMongo{}
	err = m.mongoClient.Database(m.db).Collection(m.InventoryCollection).FindOne(m.c, filter).Decode(&iMongo)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	i.ID = iMongo.ID.Hex()
	i.Name = iMongo.Name
	i.Groups = iMongo.Groups
	i.Localhost = iMongo.Localhost	
	return &i, nil
}

func (m *MongoDB) DeleteInventory(id string) error {
	return m.deleteDocument(m.InventoryCollection, id)
}

func (m *MongoDB) ListInventory() ([]models.Inventory, error) {
	projection := bson.D{{Key: "_id", Value: 1}, {Key: "name", Value: 1}}
	cursor, err := m.mongoClient.Database(m.db).Collection(m.InventoryCollection).Find(m.c, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	defer cursor.Close(m.c)
	var inventories []models.Inventory
	for cursor.Next(m.c) {
		var i inventoryMongo
		if err := cursor.Decode(&i); err != nil {
			utils.Log.Error("Failed to decode document", "error", err)
			return nil, err
		}
		inventories = append(inventories, models.Inventory{
			ID:        i.ID.Hex(),
			Name:      i.Name,
			Groups:    i.Groups,
			Localhost: i.Localhost,
		})
	}	
	return inventories, nil
}

func (m *MongoDB) CreateAnsibleJob(j *models.AnsibleJob) error {
	var jID primitive.ObjectID
	var err error
	if j.ID == "" {
		jID = primitive.NewObjectID()
	} else {
		jID, err = primitive.ObjectIDFromHex(j.ID)
		if err != nil {
			utils.Log.Error(err.Error())
			return err
		}
	}

	invID, err := primitive.ObjectIDFromHex(j.InventoryID)
	if err != nil {
		utils.Log.Error(err.Error())
		return  err
	}
	jMongo := &AnsibleJobMongo{
		ID: 		  jID,
		Playbook:     j.Playbook,
		Status:       j.Status,
		Output:       j.Output,
		Error:        j.Error,
		Step:         j.Step,
		StartTime:    primitive.NewDateTimeFromTime(j.StartTime),
		EndTime:      primitive.NewDateTimeFromTime(j.EndTime),
		OutputPlugin: j.OutputPlugin,
		InventoryID: invID,
	}

	if  err := m.persist(m.JobCollection, jMongo.ID.Hex(), jMongo); err != nil {
		return err
	}
	j.ID = jMongo.ID.Hex()
	return nil
}

func (m *MongoDB) GetAnsibleJob(id string) (*models.AnsibleJob, error) {
	var j models.AnsibleJob
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	filter := bson.M{"_id": objID}
	jSQl := AnsibleJobMongo{}
	err = m.mongoClient.Database(m.db).Collection(m.JobCollection).FindOne(m.c, filter).Decode(&jSQl)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	j.ID = jSQl.ID.Hex()
	j.Playbook = jSQl.Playbook
	j.Status = jSQl.Status
	j.Output = jSQl.Output
	j.Error = jSQl.Error
	j.Step = jSQl.Step
	j.StartTime = jSQl.StartTime.Time()
	j.EndTime = jSQl.EndTime.Time()
	j.InventoryID = jSQl.InventoryID.Hex()
	j.OutputPlugin = jSQl.OutputPlugin
	return &j, nil
}

func (m *MongoDB) DeleteAnsibleJob(id string) error {
	return m.deleteDocument(m.JobCollection, id)
}

func (m *MongoDB) ListAnsibleJob() ([]models.AnsibleJob, error) {
	cursor, err := m.mongoClient.Database(m.db).Collection(m.JobCollection).Find(m.c, bson.M{})
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	defer cursor.Close(m.c)
	var jobs []models.AnsibleJob
	for cursor.Next(m.c) {
		var i AnsibleJobMongo
		if err := cursor.Decode(&i); err != nil {
			utils.Log.Error("Failed to decode document", "error", err)
			return nil, err
		}
		jobs = append(jobs, models.AnsibleJob{
			ID:        i.ID.Hex(),
			Playbook:  i.Playbook,
			Status:    i.Status,
			Output:    i.Output,
			Error:     i.Error,
			Step:      i.Step,
			StartTime: i.StartTime.Time(),
			EndTime:   i.EndTime.Time(),
			InventoryID: i.InventoryID.Hex(),
			OutputPlugin: i.OutputPlugin,

		})
	}	
	return jobs, nil
}
