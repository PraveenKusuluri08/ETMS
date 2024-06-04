package expenses

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Praveenkusuluri08/api/expenses_tracker"
	"github.com/Praveenkusuluri08/bootstrap"
	"github.com/Praveenkusuluri08/handlers"
	endpoints "github.com/Praveenkusuluri08/types"
	"github.com/gin-gonic/gin"
)

// TODO: First check the model is properly initialized like api got the response clearly
// TODO: Simple group checks like check if the group exists or not
// TODO: Check if the user is exists or not in the group like involved peers
// TODO: Check if the expenses with the same data is aldeary exists in the database
// TODO: Perform the math opearations like splitting the expenses based on the user requirements
// TODO: If user wants to add the notes while creating the expenses then use the create note endpoint
// TODO: All this must be into the expense tracker which is basically like the logger for the expenses whhich are creating
// TODO: Store the expense tracker or the logger for the expense based on the expenses ID
// TODO: Check if the expense is perosonal or not if it is perosonal then store the expenses based on the expense for the user it self like no group involed and it purely personal and make the calculations accordingly

var exepnsesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses")
var usersCollection = bootstrap.GetCollection(bootstrap.ClientDB, "users")
var expenses_trackerCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses_tracker")

// This endpoint is specifically for the non group expenses
func CreateExpense() gin.HandlerFunc {
	expensesService := &ExpensesService{}
	return expensesService.CreateExpense()
}
func (e *ExpensesService) CreateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		userId := c.GetString("uid")
		// var wg sync.WaitGroup

		defer cancel()
		var expense Expenses
		if err := c.BindJSON(&expense); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Invalid json data",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)

			return
		}

		isExpenseTitleExists := IsExpenseWithSameTitleExists(userId, expense.Title)
		if isExpenseTitleExists {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Expense with the same title already exists",
				},
				Status: "400",
				Error:  "expense_title_already_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		if expense.IsGroup && len(expense.Split.InvolvedPeers) == 2 {
			splitExpense(&expense, expense.Amount, userId)
		}
		if expense.IsGroup && len(expense.Split.InvolvedPeers) > 2 && expense.Split.SplitType == "GROUP_EXPENSE" {
			splitExpenseWithGroup(&expense, userId)
		}
		currentTime := time.Now()
		expense.CreatedBy = userId
		expense.SplitNeedToClearBy = currentTime.Format(time.ANSIC)

		_, err := exepnsesCollection.InsertOne(ctx, expense)
		if err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Error while creating the expense",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)

			return
		}
		//wait group
		// wg.Add(len(expense.Split.InvolvedPeers))

		for _, peer := range expense.Split.InvolvedPeers {
			fmt.Println("peer", peer.PeerID)
			expense_tracker_info := expenses_tracker.ExpenseTracker_Info{
				Expense_Created_By:  userId,
				Expense_Title:       expense.Title,
				Expense_Description: expense.Description,
				Expense_Amount:      peer.Amount,
				Expense_Activity:    fmt.Sprintf(`Expense is created by %s`, userId),
				Expense_Involved_By: peer.PeerID,
			}
			handlers.PushExpense_Tracker(&expense_tracker_info)
		}
		c.JSON(http.StatusOK, expense)
	}
}

func splitExpense(expense *Expenses, amount float64, userId string) error {
	switch expense.Split.SplitType {
	case "YOU_PAID_TOTAL_SPLIT_TO_PEERS":
		splitAmount := amount / float64(len(expense.Split.InvolvedPeers))
		for i, peer := range expense.Split.InvolvedPeers {
			peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
			expense.Split.InvolvedPeers[i] = peer
		}
	case "YOU_OWED_FULL_AMOUNT_TO_PEER":
		splitAmount := amount
		for i, peer := range expense.Split.InvolvedPeers {
			peer.PeerID = userId
			peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
			expense.Split.InvolvedPeers[i] = peer
		}
	case "PEER_OWED_FULL_AMOUNT_TO_YOU":
		splitAmount := amount
		for i, peer := range expense.Split.InvolvedPeers {
			peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
			expense.Split.InvolvedPeers[i] = peer
		}

	default:
		return errors.New("Invalid split type")
	}
	return nil
}

func splitExpenseWithGroup(expense *Expenses, userId string) {
	splitAmount := expense.Amount / float64(len(expense.Split.InvolvedPeers)+1)
	for i, peer := range expense.Split.InvolvedPeers {
		peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
		expense.Split.InvolvedPeers[i] = peer
	}
	currentUser := Peer{
		PeerID: userId,
		Amount: strconv.FormatFloat(splitAmount, 'f', -1, 64),
	}
	expense.Split.InvolvedPeers = append(expense.Split.InvolvedPeers, currentUser)
}
