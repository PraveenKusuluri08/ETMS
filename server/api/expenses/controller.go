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
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

//TODO: While creating an expenses i need to check the previous involved expenses between the two users or the group

//TODO: Like check previous expenses with the userId's in the expenses and based on the amount involved in the expense between the two users then make the new expense with the amount and update the expense tracker

var exepnsesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses")
var expensesTracker = bootstrap.GetCollection(bootstrap.ClientDB, "expenses_tracker")

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
		if expense.IsGroup && len(expense.Split.InvolvedPeers)+1 == 2 {
			splitExpense(&expense, expense.Amount, userId)
		}
		if expense.IsGroup && len(expense.Split.InvolvedPeers) > 2 && expense.Split.SplitType == "GROUP_EXPENSE" {
			splitExpenseWithGroup(&expense, userId, expense.PaidBy)
		}

		currentTime := time.Now()
		expense.CreatedBy = userId
		expense.SplitNeedToClearBy = currentTime.Format(time.ANSIC)
		expense.CreatedAt = currentTime.Format(time.ANSIC)
		expenseCreatedInfo, err := exepnsesCollection.InsertOne(ctx, expense)
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
		manage_previous_expenses_amount(&expense, userId, expenseCreatedInfo.InsertedID.(primitive.ObjectID))
		c.JSON(http.StatusOK, expense)
	}
}
func splitExpense(expense *Expenses, amount float64, userId string) error {
	switch expense.Split.SplitType {
	case "YOU_PAID_TOTAL_SPLIT_TO_PEERS":
		splitAmount := amount / float64(len(expense.Split.InvolvedPeers)+1)
		for i, peer := range expense.Split.InvolvedPeers {
			peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
			expense.Split.InvolvedPeers[i] = peer
			expense.Split.OwesTo = userId
			expense.OwesAmount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
		}
	case "YOU_OWED_FULL_AMOUNT_TO_PEER":
		splitAmount := amount
		for i, peer := range expense.Split.InvolvedPeers {
			expense.Split.OwesTo = peer.PeerID
			peer.PeerID = userId
			peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
			expense.Split.InvolvedPeers[i] = peer
			expense.OwesAmount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
		}
	case "PEER_OWED_FULL_AMOUNT_TO_YOU":
		splitAmount := amount
		for i, peer := range expense.Split.InvolvedPeers {
			expense.Split.OwesTo = userId
			peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
			expense.Split.InvolvedPeers[i] = peer
			expense.OwesAmount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
		}

	default:
		return errors.New("INVALID_SPLIT_TYPE")
	}
	return nil
}

func splitExpenseWithGroup(expense *Expenses, userId string, paidBy string) {
	splitAmount := expense.Amount / float64(len(expense.Split.InvolvedPeers)+1)
	var owesAmount float64

	paidUser := Peer{
		PeerID: paidBy,
		Amount: strconv.FormatFloat(splitAmount, 'f', -1, 64),
	}
	expense.Split.InvolvedPeers = append(expense.Split.InvolvedPeers, paidUser)

	for i, peer := range expense.Split.InvolvedPeers {
		if peer.PeerID != paidBy {
			owesAmount += splitAmount
		}
		peer.Amount = strconv.FormatFloat(splitAmount, 'f', -1, 64)
		expense.Split.InvolvedPeers[i] = peer
	}

	expense.Split.OwesTo = paidBy
	expense.OwesAmount = strconv.FormatFloat(owesAmount, 'f', -1, 64)

	for i, peer := range expense.Split.InvolvedPeers {
		if peer.PeerID == paidBy {
			expense.Split.InvolvedPeers = append(expense.Split.InvolvedPeers[:i], expense.Split.InvolvedPeers[i+1:]...)
			break
		}
	}
	currentUser := Peer{
		PeerID: userId,
		Amount: strconv.FormatFloat(splitAmount, 'f', -1, 64),
	}
	expense.Split.InvolvedPeers = append(expense.Split.InvolvedPeers, currentUser)
}

func manage_previous_expenses_amount(expense *Expenses, userId string, expenseId primitive.ObjectID) {
	currentUserExpenses := GetExpensesCreatedByUser(userId)
	isExpenseAmountModified := false
	if currentUserExpenses != nil {
		for _, peer := range expense.Split.InvolvedPeers {
			if peer.Amount != "" {
				amount, err := decimal.NewFromString(peer.Amount)
				if err != nil {
					return
				}
				if currentUserExpenses.Expense_Amount != "" {
					expenseAmount, err := decimal.NewFromString(currentUserExpenses.Expense_Amount)
					if err != nil {
						return
					}
					if currentUserExpenses.Expense_Involved_By == peer.PeerID {
						expenseAmount = expenseAmount.Add(amount)
						expense.OwesAmount = expenseAmount.String()
						isExpenseAmountModified = true
					} else {
						expenseAmount = expenseAmount.Sub(amount)
						expense.OwesAmount = expenseAmount.String()
						isExpenseAmountModified = true
					}
					currentUserExpenses.Expense_Amount = expenseAmount.String()
				} else {
					currentUserExpenses.Expense_Amount = peer.Amount
					isExpenseAmountModified = false
				}
			}
		}
	}
	if isExpenseAmountModified {
		for _, peer := range expense.Split.InvolvedPeers {
			fmt.Println("peer", peer.PeerID)
			expense_tracker_info := expenses_tracker.ExpenseTracker_Info{
				Expense_Created_By:  userId,
				Expense_Title:       expense.Title,
				Expense_Description: expense.Description,
				Expense_Amount:      peer.Amount,
				Expense_Activity:    `Expense is modifed with the amount change by the previous no settled expenses`,
				Expense_Involved_By: peer.PeerID,
				Type:                "EXPENSE_AMOUNT_MODIFIED",
				ExpenseId:           expenseId,
			}
			handlers.PushExpense_Tracker(&expense_tracker_info)
		}
	} else {
		for _, peer := range expense.Split.InvolvedPeers {
			fmt.Println("peer", peer.PeerID)
			expense_tracker_info := expenses_tracker.ExpenseTracker_Info{
				Expense_Created_By:  userId,
				Expense_Title:       expense.Title,
				Expense_Description: expense.Description,
				Expense_Amount:      peer.Amount,
				Expense_Activity:    fmt.Sprintf(`Expense is created by %s`, userId),
				Expense_Involved_By: peer.PeerID,
				Type:                "CREATED",
				ExpenseId:           expenseId,
			}
			handlers.PushExpense_Tracker(&expense_tracker_info)
		}
	}
}
