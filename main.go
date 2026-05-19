package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	Err error
}

var mongoClient MongoClient

func connectToDatabase(dbName string) *MongoClient{
	database_uri := "mongodb://localhost:27017/"+dbName
	if dbName == "" {
		log.Fatal("No Database provided")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(database_uri))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to %s\n", dbName)

	return &MongoClient{Client: client, Err: err}

}

func main(){
	godotenv.Load()
	dabataseName := os.Getenv("DATABASE_NAME")
	mongoClient = *connectToDatabase(dabataseName)
	defer func(){
		if err := mongoClient.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}