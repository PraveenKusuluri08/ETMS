package bootstrap

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Mongo *mongo.Client
}

func App(mongodbURI string) Application {
	app := &Application{}
	fmt.Println("MONGODB_URI ", mongodbURI)
	app.Mongo = DBConnect()
	return *app
}
