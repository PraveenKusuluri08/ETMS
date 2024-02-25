package groups

import (
	"context"
	"fmt"
	"github.com/Praveenkusuluri08/bootstrap"
	"github.com/Praveenkusuluri08/endpoints"
	"github.com/Praveenkusuluri08/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type GroupController struct {
}

var groupCollection = bootstrap.GetCollection(bootstrap.ClientDB, "Groups")

func CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var group Group
		if err := c.BindJSON(&group); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Message: "Please provide fields properly",
				Status:  "400",
				Error:   err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		count, err := groupCollection.CountDocuments(ctx, bson.M{"group_name": group.GroupName})
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
				Message: "Group Name already exists. Please try again with different group name",
				Status:  "400",
				Error:   "group_name_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		group.ID = primitive.NewObjectID()

		inserted, err := groupCollection.InsertOne(ctx, group)
		if err != nil {
			statusInternalServerErrorResponse := endpoints.InternalServerResponse{
				Message: fmt.Sprintf("Failed to insert group"),
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, statusInternalServerErrorResponse)
			return
		}
		message := fmt.Sprintf("%s insertedDocumentId", inserted.InsertedID)
		statusCreatedResponse := endpoints.CreatedResponse{
			Message: message,
			Status:  "201",
		}
		c.JSON(http.StatusCreated, statusCreatedResponse)
	}
}

func InviteGroupMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var invitation Invitation
		if err := c.BindJSON(&invitation); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Message: "Please provide fields properly",
				Status:  "400",
				Error:   err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		fmt.Println(invitation.Users)

		// TODO:first check the user is already exists in the users array in db
		//TODO: if so then send the error message like user already exists
		// if not then perform another query to check the user is already exists in the
		// invites array if so then no need to insert the user to the array
		// perform invitation again.

		//matchStage := bson.D{
		//	{"$match", bson.D{{"$and", bson.A{bson.D{{"group_name", invitation.GroupName}},
		//		bson.D{{"invites", bson.D{{"$elemMatch", bson.D{{"$in", invitation.Users}}}}}},
		//		bson.D{{"users.email", bson.D{{"$in", invitation.Users}}}}}}}},
		//}
		//unwindStage := bson.D{{"$unwind", "$users"}}
		filter := bson.M{
			"group_name": invitation.GroupName,
			"users.email": bson.M{
				"$nin": invitation.Users, // Exclude users that are already present in the invitation.Users array
			},
		}
		update := bson.M{"$addToSet": bson.M{"invites": bson.M{"$each": invitation.Users}}}

		result, err := groupCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and modified %v documents\n", result.MatchedCount, result.ModifiedCount)

		email := &utils.SendEmailTypes{
			To:        invitation.Users,
			GroupName: invitation.GroupName}

		utils.SendEmail(email)

		c.JSON(http.StatusOK, "Invitation")
	}
}

func contains(slice []interface{}, value string) bool {
	for _, item := range slice {
		if item.(string) == value {
			return true
		}
	}
	return false
}
