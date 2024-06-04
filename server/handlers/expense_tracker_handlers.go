package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Praveenkusuluri08/api/expenses_tracker"
	"github.com/Praveenkusuluri08/bootstrap"
)

// This is used to track the expenses based on the created
// This must be like the different users different tracker information
// Iterate the involved peers and push tracker information to the tacker collection
// like created expenses, settled expenses etc...
var expenseTrackerCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses_tracker")

func PushExpense_Tracker(expense *expenses_tracker.ExpenseTracker_Info) {
	fmt.Println("From their")
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	result, err := expenseTrackerCollection.InsertOne(ctx, expense)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}
