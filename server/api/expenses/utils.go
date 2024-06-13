package expenses

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Praveenkusuluri08/api/expenses_tracker"
	"github.com/Praveenkusuluri08/bootstrap"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

var expensesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses")
var expensesTrackerCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses_tracker")

func IsExpenseWithSameTitleExists(userId string, expense_title string) bool {
	//TODO: Get all the expenses of the user
	//TODO: Filter the expenses based on the what user is create
	var ctx, cancel = context.WithTimeout(context.Background(), 10)
	fmt.Println(userId)
	defer cancel()
	today := time.Now().Format(time.ANSIC)
	filter := bson.M{"created_by": userId, "title": expense_title, "created_at": bson.M{
		"$regex": "^" + regexp.QuoteMeta(today[:11]) + ".*",
	}}
	var userExpense []bson.M
	if err := expensesCollection.FindOne(ctx, filter).Decode(&userExpense); err != nil {
		log.Println(err)
		return false
	}
	return len(userExpense) > 0
}

func GetExpensesCreatedByUser(userId string) *expenses_tracker.ExpenseTracker_Info {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	expenses := &expenses_tracker.ExpenseTracker_Info{}
	if err := expensesTrackerCollection.FindOne(ctx, bson.M{"expense_created_by": userId}).Decode(&expenses); err != nil {
		log.Println(err)
		return nil
	}
	return expenses
}

func InvolvedUserExpensesByUser(userId string) *expenses_tracker.ExpenseTracker_Info {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	expenses := &expenses_tracker.ExpenseTracker_Info{}
	if err := expensesTrackerCollection.FindOne(ctx, bson.M{"expense_involved_by": userId}).Decode(&expenses); err != nil {
		log.Println(err)
		return nil
	}
	return expenses
}
