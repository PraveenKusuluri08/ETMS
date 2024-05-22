package bootstrap

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Mongo *mongo.Client
}

func App(mongodbURI string, isTest bool) Application {
	app := &Application{}

	app.Mongo = DBConnect()

	fmt.Println("MONGODB_URI ", mongodbURI)
	return *app
}
