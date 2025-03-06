package db

import (
	"ansible-runner/utils"
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents a MongoDB configuration.
type MongoDB struct {
	mongoClient   *mongo.Client
	db            string
	JobCollection string
	InventoryCollection string
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
		Username: username,
		Password: password,
        AuthMechanism: "SCRAM-SHA-256",
        AuthSource: MongoClient.db,
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

func (m *MongoDB) Persist(collection string, keyFilter string, value any) error{
	objID, err := primitive.ObjectIDFromHex(keyFilter)

	if err != nil{
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
		utils.Log.Error(err.Error())
		return err
	}
	return nil
}

func  Retrieve[T any](m *MongoDB,collection string, keyFilter string) (*T, error) {
	var res T
	objID, err := primitive.ObjectIDFromHex(keyFilter)

	if err != nil{
		utils.Log.Error(err.Error())
		return &res, err
	}

	filter := bson.M{"_id": objID}
	err = m.mongoClient.Database(m.db).Collection(collection).FindOne(m.c, filter).Decode(&res)
	if err != nil {
        utils.Log.Error(err.Error())
		return &res, err
	}
	return &res, nil
}

func ListDocuments[T any](m *MongoDB, collection string, projection bson.D) ([]T, error) {
    var results []T

    cursor, err := m.mongoClient.Database(m.db).Collection(collection).Find(m.c, bson.M{}, options.Find().SetProjection(projection))
    if err != nil {
        utils.Log.Error("Failed to list documents", "error", err)
        return nil, err
    }
    defer cursor.Close(m.c)

    for cursor.Next(m.c) {
        var result T
        if err := cursor.Decode(&result); err != nil {
            utils.Log.Error("Failed to decode document", "error", err)
            return nil, err
        }
        results = append(results, result)
    }

    if err := cursor.Err(); err != nil {
        utils.Log.Error("Cursor error", "error", err)
        return nil, err
    }

    return results, nil
}


// Close closes the MongoDB connection.
func (m *MongoDB) Close() {
    if m != nil && m.mongoClient != nil {
	    m.mongoClient.Disconnect(m.c)
    }
}
