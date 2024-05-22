package notes

import (
	"context"
	"net/http"
	"time"

	"github.com/Praveenkusuluri08/bootstrap"
	endpoints "github.com/Praveenkusuluri08/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var notesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "notes")

func CreateNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10)
		defer cancel()
		var note Notes
		if err := c.ShouldBindJSON(&note); err != nil {
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
		//check if the notes for the expense is exists
		count, err := notesCollection.CountDocuments(ctx, bson.M{"expense_id": note.ExpenseId})
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to get count of the documents",
				},
				Status: "500",
				Error:  err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		if count > 0 {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Notes for the expense already exists",
				},
				Status: "400",
				Error:  "notes_already_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		createdAt := time.Now()
		note.CreatedAt = createdAt.Format(time.ANSIC)
		note.ID = primitive.NewObjectID()
		result, err := notesCollection.InsertOne(ctx, note)
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to create the notes",
				},
				Status: "500",
				Error:  err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		c.JSON(http.StatusOK, gin.H{"notes": result})
	}
}
