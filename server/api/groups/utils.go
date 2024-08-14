package groups

import (
	"context"
	"time"

	"github.com/Praveenkusuluri08/api/users"
	"github.com/Praveenkusuluri08/bootstrap"
	"go.mongodb.org/mongo-driver/bson"
)

var friendsCollection = bootstrap.GetCollection(bootstrap.ClientDB, "Friends")

func MakeUsersFriendsEachOther() {

}

func GetInviter(uid string) (users.UserStruct, error) {
	// TODO: Get the inviter from the token
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var u users.UserStruct
	filter := bson.M{"uid": uid}
	user := usersCollection.FindOne(ctx, filter)
	if err := user.Decode(&u); err != nil {
		return users.UserStruct{}, err
	}
	return u, nil
}
