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

func DBConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	fmt.Println(os.Getenv("MONGODB_URI"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/local"))
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
	collection := client.Database("etms").Collection(collectionName)
	return collection
}
