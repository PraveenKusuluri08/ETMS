package groups

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Praveenkusuluri08/bootstrap"
	endpoints "github.com/Praveenkusuluri08/types"
	"github.com/Praveenkusuluri08/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupController struct {
}

var groupCollection = bootstrap.GetCollection(bootstrap.ClientDB, "Groups")
var usersCollection = bootstrap.GetCollection(bootstrap.ClientDB, "Users")

func CreateGroup() gin.HandlerFunc {
	groupService := &GroupService{}
	return groupService.CreateGroup()
}

func UpdateGroup() gin.HandlerFunc {
	groupService := &GroupService{}
	return groupService.UpdateGroup()
}

func InviteGroupMembers() gin.HandlerFunc {
	groupService := &GroupService{}
	return groupService.InviteGroupMembers()
}

func AcceptInvitation() gin.HandlerFunc {
	groupService := &GroupService{}
	return groupService.AcceptInvitation()
}

func DisplayUsers() gin.HandlerFunc {
	groupService := &GroupService{}
	return groupService.DisplayUsers()
}

func RemoveGroupMember() gin.HandlerFunc {
	groupService := &GroupService{}
	return groupService.RemoveGroupMember()
}

//	@Summary		Create a new group
//	@Description	Create a new group with the provided group name to manage the expenses between the tenets
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Param			group			body		Group	true	"Group"
//	@Success		200				{object}	endpoints.CreatedResponse
//	@Failure		400				{object}	endpoints.BadRequestResponse
//	@Failure		500				{object}	endpoints.InternalServerResponse
//	@Router			/api/v1/groups/creategroup [post]
//	@Tags			Groups
func (g *GroupService) CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var group Group
		fmt.Printf("This is the Body %s", c.Request.Body)
		if err := c.BindJSON(&group); err != nil {
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
		count, err := groupCollection.CountDocuments(ctx, bson.M{"group_name": group.GroupName})
		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Invalid json data",
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
					Name: "Group Name already exists. Please try again with different group name",
				},
				Status: "400",
				Error:  "group_name_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)

			return
		}
		group.ID = primitive.NewObjectID()
		group.CreatedBy = c.GetString("uid")

		inserted, err := groupCollection.InsertOne(ctx, group)
		if err != nil {
			statusInternalServerErrorResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: fmt.Sprintf("Failed to insert group"),
				},
				Status: "500",
				Error:  err.Error(),
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

//	@Summary		Update the group name
//	@Description	Update group name with the new group name
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string				true	"Bearer token"
//	@Param			updateGroup		body		UpdateGroupStruct	true	"Update Group"
//	@Success		200				{object}	endpoints.SuccessResponse
//	@Failure		400				{object}	endpoints.BadRequestResponse
//	@Failure		500				{object}	endpoints.InternalServerResponse
//	@Router			/api/v1/groups/update_group_name [put]
//	@Tags			Groups
func (group *GroupService) UpdateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var g UpdateGroupStruct
		if err := c.BindJSON(&g); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide fields properly",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)

			return
		}
		fmt.Println("groupName", g)
		filter := bson.M{"group_name": g.GroupName}
		count, err := groupCollection.CountDocuments(ctx, filter)
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
		if count == 0 {
			statusBadRequest := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Group is not exists. Please try again with different group name",
				},
				Status: "400",
				Error:  "group_name_not_exists",
			}
			c.JSON(http.StatusBadRequest, statusBadRequest)

			return
		}
		update := bson.M{"$set": bson.M{"group_name": g.NewGroupName}}

		result := groupCollection.FindOneAndUpdate(ctx, filter, update)

		if result.Err() != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to get count of the documents",
				},
				Status: "500",
				Error:  result.Err().Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)

			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Group name updated successfully"})
	}
}

//	@Summary		Update the group name
//	@Description	Update group name with the new group name
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string		true	"Bearer token"
//
//	@Param			invite			body		Invitation	true	"Invite User"
//	@Success		200				{object}	endpoints.InviteGroupMembersResponse
//	@Failure		400				{object}	endpoints.BadRequestResponse
//	@Failure		500				{object}	endpoints.InternalServerResponse
//	@Router			/api/v1/groups/invite [POST]
//
//	@Tags			Groups
func (g *GroupService) InviteGroupMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var invitation Invitation
		if err := c.BindJSON(&invitation); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide fields properly",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)

			return
		}
		fmt.Println(invitation)

		//check if the group is exists or not in the database
		filter_group := bson.M{"group_name": invitation.GroupName}
		count, group_count_error := groupCollection.CountDocuments(ctx, filter_group)
		if group_count_error != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to check group exists or not in the db",
				},
				Status: "500",
				Error:  group_count_error.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		if count == 0 {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Group Name does not exists",
				},
				Status: "400",
				Error:  "group_name_not_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)

			return
		}
		//TODO:First iterate the users array in the request to check the user is exists or not then perform the compound qurery
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
		var users_not_exists []string
		var users_exists []string
		for _, u := range invitation.Users {
			fmt.Println(u)
			userFilter := bson.M{"email": u}
			userCount, userErr := usersCollection.CountDocuments(ctx, userFilter)
			if userErr != nil {
				log.Fatal(userErr)
			}
			if userCount > 0 {
				users_exists = append(users_exists, u)
			} else {
				users_not_exists = append(users_not_exists, u)
			}
		}
		filter := bson.M{
			"group_name": invitation.GroupName,
			"users.email": bson.M{
				"$nin": users_exists, // Exclude users that are already present in the invitation.Users array
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

		c.JSON(http.StatusOK, gin.H{"message": "Invitations send successfully", "non_existing_users": users_not_exists, "total_no_existing_users": len(users_not_exists)})
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

//	@Summary		Accept invitation to group member
//	@Description	Accept invitation to group member with the provided group name and email
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Authorization		header		string					true	"Bearer token"
//
//	@Param			acceptInvitation	body		AcceptInvitationStruct	true	"Accept Inviation"
//	@Success		200					{object}	endpoints.AcceptInvitationResponse
//	@Failure		400					{object}	endpoints.BadRequestResponse
//	@Failure		500					{object}	endpoints.InternalServerResponse
//	@Router			/api/v1/groups/accept_invitation [POST]
//
//	@Tags			Groups
func (g *GroupService) AcceptInvitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var acceptInv AcceptInvitationStruct
		if err := c.BindJSON(&acceptInv); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide fields properly",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		filter := bson.M{"group_name": acceptInv.GroupName}
		count, err := groupCollection.CountDocuments(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Group does not exist",
				},
				Status: "400",
				Error:  "group_does_not_exist",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}

		//check if the user is exists in the invitation array
		filterUserInv := bson.M{"invites": acceptInv.Email}

		invcount, inverr := groupCollection.CountDocuments(ctx, filterUserInv)
		if inverr != nil {
			log.Fatal(err)
		}
		if invcount < 1 {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "User is not invited. Please invite the user first",
				},
				Status: "400",
				Error:  "user_not_invited",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}

		var users []map[string]string
		users = append(users, map[string]string{"email": acceptInv.Email})
		userFilter := bson.M{
			"group_name": acceptInv.GroupName,
			"users.email": bson.M{
				"$in": users,
			},
		}

		count, err = groupCollection.CountDocuments(ctx, userFilter)

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
					Name: "User already exists. Please try again with different user name",
				},
				Status: "400",
				Error:  "user_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}

		update := bson.M{"$push": bson.M{"users": bson.M{"email": acceptInv.Email}}}
		result, err := groupCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and modified %v documents\n", result.MatchedCount, result.ModifiedCount)

		//aftter the user data is updated then pop the user invitation from the invitation array

		updateInvitationArray := bson.M{"$pull": bson.M{"invites": acceptInv.Email}}
		updateRes := groupCollection.FindOneAndUpdate(ctx, bson.M{"group_name": acceptInv.GroupName}, updateInvitationArray)

		if updateRes.Err() != nil {
			log.Fatal(updateRes.Err())
		}
		fmt.Println(updateRes)

		//after this in the users collection need to add the group name into the groups array.

		c.JSON(http.StatusOK, gin.H{"message": "Invitation Accepted"})
	}
}

//	@Summary		Get Users Present in the group
//	@Description	get all the users present in the group by the group name
//	@Tags			Groups
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Param			group_name		query		string	true	"Group name"
//	@Success		200				{object}	endpoints.GetUsersResponse
//	@Failure		400				{object}	endpoints.BadRequestResponse
//	@Failure		500				{object}	endpoints.InternalServerResponse
//	@Router			/api/v1/groups/get_users [post]
func (g *GroupService) DisplayUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		groupName, isQueryParamExists := c.GetQuery("group_name")
		groupName = strings.TrimSpace(groupName)
		if !isQueryParamExists {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide group name",
				},
				Status: "400",
				Error:  "group_name_required",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		var result bson.M
		err := groupCollection.FindOne(ctx, bson.M{"group_name": groupName}).Decode(&result)

		if err != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to get the data from db",
				},
				Status: "500",
				Error:  err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": result["users"]})
	}
}

//	@Summary		Remove member from the group
//	@Description	Remove member from the group based on the email
//	@Tags			Groups
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer token"
//	@Param			removeMember	body		AcceptInvitationStruct				true	"Remove Group Member"
//	@Success		200				{object}	Group								"Group updated successfully"
//	@Failure		400				{object}	endpoints.BadRequestResponse		"Invalid request"
//	@Failure		500				{object}	endpoints.InternalServerResponse	"Internal server error"
//	@Router			/api/v1/groups/remove_group_member [put]
func (g *GroupService) RemoveGroupMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var g AcceptInvitationStruct
		if err := c.BindJSON(&g); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide fields properly",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		fmt.Println(g.GroupName)

		// first check the group name is exists or not and if so then find the user and delete the user from the group
		// here user lies in the users array in db as users:[{email: "user@example.com}]
		filter := bson.M{"group_name": g.GroupName}

		count, err := groupCollection.CountDocuments(ctx, filter)
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
		if count == 0 {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Group is not exists. Please try again with different group name",
				},
				Status: "400",
				Error:  "group_name_not_exists",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		update := bson.M{"$pull": bson.M{"users": bson.M{"email": g.Email}}}
		result := groupCollection.FindOneAndUpdate(ctx, filter, update)
		if result.Err() != nil {
			internalServerResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to delete the user from the group",
				},
				Status: "500",
				Error:  result.Err().Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerResponse)
			return
		}
		// after the user email is removed from the users array in the groups collection then need to remove
		// the group form the users collection users data also
		var updatedGroup Group
		result.Decode(&updatedGroup)

		c.JSON(http.StatusOK, gin.H{"updatedGroup": updatedGroup})
	}
}
