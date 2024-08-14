package expenses_tracker

import "go.mongodb.org/mongo-driver/bson/primitive"

type Expense_tracker struct {
	Title                string  `json:"title,omitempty" bson:"title" validate:"required"`
	Description          string  `json:"description,omitempty" bson:"description" validate:"required"`
	Amount               float64 `json:"amount,omitempty" bson:"amount" validate:"required"`
	SettledPaymentMethod string  `json:"paymentMethod" bson:"paymentMethod" oneof:"CASH_PAYMENT,PayPal,ZELLE_PAYMENT,APPLE_PAY,Google_PAY"`
	ExpenseId            string  `json:"expenseId,omitempty" bson:"expense"`
	Expense_Activity     string  `json:"expense_info,omitempty" bson:"expense_activity"`
	IsExpenseSettled     bool    `json:"isExpenseSettled,omitempty" bson:"isExpenseSettled"`
}

const (
	CASH_PAYMENT = iota
	PayPal
	ZELLE_PAYMENT
	APPLE_PAY
	Google_PAY
)

type ExpenseTracker_Info struct {
	Expense_Created_By  string             `json:"userId,omitempty"`
	Expense_Title       string             `json:"expenseTitle,omitempty"`
	Expense_Description string             `json:"expenseDescription,omitempty"`
	Expense_Amount      string             `json:"expenseAmount,omitempty"`
	Expense_Activity    string             `json:"expenseActivity,omitempty" bson:"expense_activity"`
	Expense_Involved_By string             `json:"expenseInvolvedBy,omitempty" bson:"expense_involved_by"`
	Type                string             `json:"type,omitempty" bson:"type"`
	ExpenseId           primitive.ObjectID `json:"expenseId,omitempty" bson:"expenseId"`
	AmountPaidBy        string             `json:"amountPaidBy,omitempty" bson:"amountPaidBy"`
}
