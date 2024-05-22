package expenses

import (
	"net/http"

	"github.com/Praveenkusuluri08/bootstrap"
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
var groupsCollection = bootstrap.GetCollection(bootstrap.ClientDB, "group")
var notesCollection = bootstrap.GetCollection(bootstrap.ClientDB, "notes")
var expenses_trackerCollection = bootstrap.GetCollection(bootstrap.ClientDB, "expenses_tracker")

func CreateExpense() gin.HandlerFunc {
	expensesService := &ExpensesService{}
	return expensesService.CreateExpense()
}
func (e *ExpensesService) CreateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense Expenses
		if err := c.BindJSON(&expense); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Bad Request",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}

	}
}
