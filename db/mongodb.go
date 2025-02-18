package db

import (
	"ansible-runner/models"
	"ansible-runner/utils"
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents a MongoDB configuration.
type MongoDB struct {
	mongoClient   *mongo.Client
	db            string
	jobCollection string
	c             context.Context
}

// NewMongoDB creates a new MongoDB configuration.
var MongoClient *MongoDB

func NewMongoClient() *MongoDB {
    if MongoClient != nil {
        return MongoClient
    }
	MongoClient = new(MongoDB)
	MongoClient.c = context.Background()
	var err error
	uri := viper.GetString("mongodb.uri")
    MongoClient.db = viper.GetString("mongodb.database")
	MongoClient.jobCollection = viper.GetString("mongodb.jobCollection")

	username := viper.GetString("mongodb.user")
	password := viper.GetString("mongodb.password")
	password, err = utils.Decrypt(password)
	if err != nil {
		panic(err)
	}

	credential := options.Credential{
		Username: username,
		Password: password,
        AuthMechanism: "SCRAM-SHA-256",
        AuthSource: MongoClient.db,
	}
    fmt.Println(credential.AuthSource)
    fmt.Println(credential.Username)
    fmt.Println(credential.Password)
    fmt.Println(uri)
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)

	MongoClient.mongoClient, err = mongo.Connect(MongoClient.c, clientOptions)
	if err != nil {
		panic(err)
	}


	databaseNames, err := MongoClient.mongoClient.ListDatabaseNames(MongoClient.c, bson.M{})
	if err != nil {
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
		err = MongoClient.mongoClient.Database(MongoClient.db).CreateCollection(MongoClient.c, MongoClient.jobCollection)
		if err != nil {
			panic(err)
		}
	}
	collectionNames, err := MongoClient.mongoClient.Database(MongoClient.db).ListCollectionNames(MongoClient.c, bson.M{})
	if err != nil {
		panic(err)
	}

	collectionExists := false
	for _, collName := range collectionNames {
		if collName == MongoClient.jobCollection {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		err = MongoClient.mongoClient.Database(MongoClient.db).CreateCollection(MongoClient.c, MongoClient.jobCollection)
		if err != nil {
			panic(err)
		}
	}
    return MongoClient
}

// PersistJob saves the AnsibleJob to the database.
func (m *MongoDB) PersistJob(job *models.AnsibleJob) error {
    update := bson.M{
        "$set": job,
    }
	updateOptions := options.Update().SetUpsert(true)
	filter := bson.M{"id": job.ID}
	_, err := m.mongoClient.Database(m.db).Collection(m.jobCollection).UpdateOne(m.c, filter, update, updateOptions)
	return err
}

// RetrieveJob retrieves an AnsibleJob from the database by its ID.
func (m *MongoDB) RetrieveJob(jobID string) (*models.AnsibleJob, error) {
	var job models.AnsibleJob
	filter := bson.M{"id": jobID}
	err := m.mongoClient.Database(m.db).Collection(m.jobCollection).FindOne(m.c, filter).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// Close closes the MongoDB connection.
func (m *MongoDB) Close() {
    if m.mongoClient != nil {
	    m.mongoClient.Disconnect(m.c)
    }
}
