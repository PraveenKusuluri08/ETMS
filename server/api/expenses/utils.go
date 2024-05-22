package expenses

import (
	"github.com/Praveenkusuluri08/bootstrap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

var expensesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses")

func IsExpenseWithSameTitleExists(userId primitive.ObjectID, expense_title string) bool {
	//TODO: Get all the expenses of the user
	//TODO: Filter the expenses based on the what user is create
	var ctx, cancel = context.WithTimeout(context.Background(), 10)
	defer cancel()
	filter := bson.M{"user_id": userId}
	cursor, err := expensesCollection.Find(ctx, filter)
	if err != nil {
		return true
	}
	defer cursor.Close(ctx)
	// cursor.All()

	return false
}
