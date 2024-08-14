package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var err error

func DBConnect() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// mongo_uri := os.Getenv("MONGODB_URI")

	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://praveen_admin:Praveen8919296298@cluster1.uyl3sto.mongodb.net/?retryWrites=true&w=majority&appName=Cluster1"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongodb connected")
	return client
}

var ClientDB *mongo.Client = DBConnect()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	appEnv := os.Getenv("APP_ENV")
	var collection *mongo.Collection
	if appEnv == "prod" {
		collection = client.Database("etms").Collection(collectionName)
	} else {
		collection = client.Database("etms_test").Collection(collectionName)
	}
	return collection
}
