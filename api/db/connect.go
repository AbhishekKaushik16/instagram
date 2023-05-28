package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	err    error
)

// TODO: Add logger for mongo query
func connectToDb() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://abhishek:abhishek@cluster0.4gab7.mongodb.net/")

	// Connect to MongoDB
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func GetDBClient() (*mongo.Client, error) {
	connectToDb()
	return client, err
}

func SetupDbClient() (*mongo.Client, *mongo.Database, error) {
	// Set up MongoDB connection
	client, err = GetDBClient()
	if err != nil {
		log.Fatal(err)
	}
	// Access the "users" collection
	return client, client.Database("users"), err
}

func DisconnectDB() {

	// Disconnect from MongoDB when finished
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
}
