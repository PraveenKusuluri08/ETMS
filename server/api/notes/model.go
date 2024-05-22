package notes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notes struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Notes     string             `bson:"notes,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	CreatedAt string             `bson:"created_at,omitempty"`
	ExpenseId primitive.ObjectID `bson:"expense_id,omitempty"`
}

type NotesInterface interface {
	CreateNotes() gin.HandlerFunc
	UpdateNotes() gin.HandlerFunc
	DeleteNotes() gin.HandlerFunc
	GetAllNotesOfExpenses() gin.HandlerFunc
	GetNotesOfExpense() gin.HandlerFunc
	GetAllNotesOfUserForAllExpenses() gin.HandlerFunc
	GetNotesOfUserBasedOnExpense() gin.HandlerFunc
}

type NotesService struct{}
