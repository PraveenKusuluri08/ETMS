package users

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserName(uid string) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	var user UserStruct
	if err := usersCollection.FindOne(ctx, bson.M{"_id": uid}).Decode(&user); err != nil {
		log.Fatal(err.Error())
		return ""
	}
	return user.Username
}
