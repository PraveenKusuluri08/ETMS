package users

import (
	"net/http"
	"time"

	"github.com/Praveenkusuluri08/bootstrap"
	"github.com/Praveenkusuluri08/endpoints"
	"github.com/Praveenkusuluri08/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

var usercollection = bootstrap.GetCollection(bootstrap.ClientDB, "Users")

func (u User) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.BindJSON(&u); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Message: "Please provide the fields properly",
				Status:  "400",
				Error:   err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		filter := bson.M{"email": u.Email}
		count, err := usercollection.CountDocuments(ctx, filter)
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		if count > 0 {
			badRequestResponse := endpoints.BadRequestResponse{
				Message: "Email already exists. Please try again with different email address",
				Status:  "400",
				Error:   "email_already_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		current_time := time.Now()
		hasPassword, _ := utils.HashPassword(u.Password)
		// createdAt := current_time.
	}
}
