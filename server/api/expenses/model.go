package expenses

import (
	"github.com/Praveenkusuluri08/api/notes"
	"github.com/gin-gonic/gin"
)

type ExpensesControllers interface {
	CreateExpense() gin.HandlerFunc
}

type Expenses struct {
	ID                 string       `json:"expenses_id,omitempty" bson:"expenses_id" validate:"required"`
	Title              string       `json:"title,omitempty" bson:"title" validate:"required"`
	Description        string       `json:"description,omitempty" bson:"description" validate:"required"`
	Amount             string       `json:"amount,omitempty" bson:"amount" validate:"required"`
	Category           string       `json:"category,omitempty" bson:"category" validate:"required"`
	CreatedBy          string       `json:"created_by,omitempty" bson:"created_by" validate:"required"`
	IsGroup            bool         `json:"is_group_expense,omitempty" bson:"is_group_expense"`
	IsPersonal         bool         `json:"is_personal_expense,omitempty" bson:"is_personal_expense"`
	Group              Group        `json:"group_expense,omitempty" bson:"group_expense"`
	Split              Split        `json:"group_expense_split,omitempty" bson:"group_expense_split"`
	SplitNeedToClearBy string       `json:"split_need_to_clear_by,omitempty" bson:"split_need_to_clear_by"`
	CreatedAt          string       `json:"created_at,omitempty" bson:"created_at"`
	Notes              *notes.Notes `json:"notes,omitempty" bson:"notes"`
}

type Group struct {
	ID     string `json:"group_id,omitempty" bson:"group_id"`
	PaidBy string `json:"paid_by,omitempty" bson:"paid_by" enum:"YOU_PAID_TOTAL_SPLIT_TO_PEERS=>You, YOU_OWED_FULL_AMOUNT_TO_PEER=>You, PEER_OWED_FULL_AMOUNT_TO_YOU=>Peer"`
}

type Split struct {
	GroupID       string `json:"group_id,omitempty" bson:"group_id"`
	InvolvedPeers []Peer `json:"involved_peers,omitempty" bson:"involved_peers"`
}

type Peer struct {
	PeerID string `json:"peer_id,omitempty" bson:"peer_id"`
	Amount string `json:"amount,omitempty" bson:"amount"`
}

// this is for the single person split
const (
	YOU_PAID_TOTAL_SPLIT_TO_PEERS = iota
	YOU_OWED_FULL_AMOUNT_TO_PEER
	PEER_OWED_FULL_AMOUNT_TO_YOU
)

type ExpensesInterface interface {
	CreateExpense() gin.HandlerFunc
}

type ExpensesService struct{}
