package groups

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Praveenkusuluri08/bootstrap"
	"github.com/Praveenkusuluri08/endpoints"
	"github.com/Praveenkusuluri08/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func AcceptInvitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		email := c.Query("email")
		groupName := c.Query("groupName")

		fmt.Println(groupName, email)

		// check for the group is exists or not and and also check for the user is already in the users array in db
		filter := bson.M{"groupName": groupName}
		count, err := groupCollection.CountDocuments(ctx, filter)
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		if count == 0 {
			statusBadRequest := endpoints.BadRequestResponse{
				Message: "Group Name already exists. Please try again with different group name",
				Status:  "400",
				Error:   "group_name_exists",
			}
			c.JSON(http.StatusBadRequest, statusBadRequest)
			return
		}

		//check for the user is already in the users array in db
		var users []string
		users = append(users, email)
		userFilter := bson.M{
			"group_name": groupName,
			"users.email": bson.M{
				"$in": users,
			},
		}
		count, err = groupCollection.CountDocuments(ctx, userFilter)
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
				Message: "User already exists. Please try again with different user name",
				Status:  "400",
				Error:   "user_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}

		// if the user is not exists then need to insert the user to the users array in db
		update := bson.M{"$push": bson.M{"users": bson.M{"email": email}}}
		updateInfo := groupCollection.FindOneAndUpdate(ctx, filter, update)

		doc := bson.M{}
		err1 := updateInfo.Decode(&doc)
		if err1 != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   err1.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		c.JSON(http.StatusCreated, doc)
	}
}

func DisplaUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var group Group
		defer cancel()
		type D struct {
			groupname string
		}
		var d D
		//check if the group name is already exists or not and next find the users by using projection print only the users emails address
		err := c.BindJSON(&d)
		if err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Message: "Please provide fields properly",
				Status:  "400",
				Error:   err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		filter := bson.M{"group_name": d.groupname}
		count, err := groupCollection.CountDocuments(ctx, filter)
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		if count == 0 {
			statusBadRequest := endpoints.BadRequestResponse{
				Message: "Group is not exists. Please try again with different group name",
				Status:  "400",
				Error:   "group_name_not_exists",
			}
			c.JSON(http.StatusBadRequest, statusBadRequest)
			return
		}

		cursor, err := groupCollection.Find(ctx, filter)
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}

		if err := cursor.All(ctx, &group); err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get the documents",
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		c.JSON(http.StatusOK, group)
	}
}

func UpdateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		type group struct {
			group_name     string
			new_group_name string
		}
		var g group
		if err := c.BindJSON(&g); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Message: "Please provide fields properly",
				Status:  "400",
				Error:   err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		filter := bson.M{"group_name": g.group_name}
		count, err := groupCollection.CountDocuments(ctx, filter)
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		if count == 0 {
			statusBadRequest := endpoints.BadRequestResponse{
				Message: "Group is not exists. Please try again with different group name",
				Status:  "400",
				Error:   "group_name_not_exists",
			}
			c.JSON(http.StatusBadRequest, statusBadRequest)
			return
		}
		update := bson.M{"$set": bson.M{"group_name": g.new_group_name}}

		result := groupCollection.FindOneAndUpdate(ctx, filter, update)

		if result.Err() != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Message: "Failed to get count of the documents",
				Status:  "500",
				Error:   result.Err().Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Group name updated successfully"})
	}
}

func RemoveGroupMember() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),10*time.Second)

		defer cancel()

		type group struct{
			group_name string
			email string
		}

		var g group

		if err:=c.BindJSON(&g);err!=nil{
			badRequestResponse := endpoints.BadRequestResponse{
                Message: "Please provide fields properly",
                Status:  "400",
                Error:   err.Error(),
            }
            c.JSON(http.StatusBadRequest, badRequestResponse)
            return
		}

		// first check the group name is exists or not and if so then find the user and delete the user from the group
		// here user lies in the users array in db as users:[{email: "user@example.com}]
		filter:= bson.M{"group_name":g.group_name}

		count,err:=groupCollection.CountDocuments(ctx, filter)
		if err!=nil{
			internalServerResponse := endpoints.InternalServerResponse{
                Message: "Failed to get count of the documents",
                Status:  "500",
                Error:   err.Error(),
            }
            c.JSON(http.StatusInternalServerError, internalServerResponse)
            return
		}
		if count==0{
			badRequestResponse:= endpoints.BadRequestResponse{
				Message: "Group is not exists. Please try again with different group name",
                Status:  "400",
                Error:   "group_name_not_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		update := bson.M{"$pull":bson.M{"users":bson.M{"email":g.email}}}
		result:=groupCollection.FindOneAndUpdate(ctx, filter,update)
		if result.Err()!=nil{
			internalServerResponse := endpoints.InternalServerResponse{
                Message: "Failed to delete the user from the group",
                Status:  "500",
                Error:   result.Err().Error(),
            }
            c.JSON(http.StatusInternalServerError, internalServerResponse)
            return
		}
		var updatedGroup Group
		result.Decode(&updatedGroup)

		c.JSON(http.StatusOK,gin.H{"updatedGroup":updatedGroup})
	}
}