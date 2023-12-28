package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

type MongoDbConnect struct{}

type MongoDbInterface interface {
	ConnectDB() (*mongo.Client, error)
	CloseDb() error
	GetCollection(string) *mongo.Collection
}

// var MongoDBConfig mongoDbInterface = MongoDbConnect{}

func NewDatabase() MongoDbInterface {
	return &MongoDbConnect{}
}

func (dbCon MongoDbConnect) ConnectDB() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(ENV.MONGODB_IRI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to check if the connection is successful
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db = client.Database("quickstart")

	fmt.Println("Connected to MongoDB!")
	return client, nil

}

func (dbCon MongoDbConnect) GetCollection(collectionName string) *mongo.Collection {

	return db.Collection(collectionName)
}

func (dbCon MongoDbConnect) CloseDb() error {

	return db.Client().Disconnect(context.Background())
}
