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

var exepnsesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses")

func CreateExpense() gin.HandlerFunc {
	expensesService := &ExpensesService{}
	return expensesService.CreateExpense()
}

// @Summary		Create new expense
// @Description	Create a new expense based on the user's amount and the preferences
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			Authorization	header		string		true	"Bearer token"
// @Param			expense			body		Expenses	true	"Expenses"
// @Success		200				{object}	endpoints.CreatedResponse
// @Failure		400				{object}	endpoints.BadRequestResponse
// @Failure		500				{object}	endpoints.InternalServerResponse
// @Router			/api/v1/expenses/create [post]
// @Tags			Expenses
func (e *ExpensesService) CreateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		userId := c.GetString("uid")

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
				Expense_Activity:    `Expense is modifed with the amount change by the previous non settled expenses`,
				Expense_Involved_By: peer.PeerID,
				Type:                "EXPENSE_AMOUNT_MODIFIED",
				ExpenseId:           expenseId,
			}
			handlers.PushExpense_Tracker(&expense_tracker_info)
		}
	} else {
		for _, peer := range expense.Split.InvolvedPeers {
			fmt.Println("peer", peer.PeerID)
			if peer.PeerID != userId {
				expense_tracker_info := expenses_tracker.ExpenseTracker_Info{
					Expense_Created_By:  userId,
					Expense_Title:       expense.Title,
					Expense_Description: expense.Description,
					Expense_Amount:      peer.Amount,
					Expense_Activity:    fmt.Sprintf(`Expense is created by %s`, userId),
					Expense_Involved_By: peer.PeerID,
					Type:                expense.Split.SplitType,
					ExpenseId:           expenseId,
					AmountPaidBy:        expense.PaidBy,
				}
				handlers.PushExpense_Tracker(&expense_tracker_info)
			}
		}
	}
}

// TODO:Get Current user expenses
// TODO: This function is used to get the expenses for the current user which is created by that
// user or involved by the current user with any other expenses
// settled expenses non settled expenses bascically all the expenses belongs to the current user
// TODO: Based on the user login needs to fetch the expenses for the user and make sure that need to fetch the expenses which is createdBy this user and fetch all the expenses
// like this user involved other users involved in this user

func GetCurrentUserExpenses() gin.HandlerFunc {
	expensesService := &ExpensesService{}
	return expensesService.GetCurrentUserExpenses()
}

func (e *ExpensesService) GetCurrentUserExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// TODO:Get Group Expenses based on the users
// TODO: Get the amount based on the person in the group and calculate the amount for all the expenses for all the users
