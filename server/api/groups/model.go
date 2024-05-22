package groups

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupInterface interface {
	UpdateGroup() gin.HandlerFunc
	CreateGroup() gin.HandlerFunc
	RemoveGroupMember() gin.HandlerFunc
	InviteGroupMembers() gin.HandlerFunc
	AcceptInvitation() gin.HandlerFunc
	DisplayUsers() gin.HandlerFunc
}

type Group struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	GroupName string              `bson:"group_name,omitempty" validate:"required"`
	Users     []map[string]string `bson:"users,omitempty" validate:"required"`
	Type      string              `bson:"type,omitempty" validate:"required"`
}

type Invitation struct {
	GroupName string   `json:"groupname"`
	Users     []string `json:"users"`
}

type AcceptInvitationStruct struct {
	GroupName string `json:"groupname"`
	Email     string `json:"email"`
}

type UpdateGroupStruct struct {
	GroupName    string `json:"groupname"`
	NewGroupName string `json:"new_group_name"`
}

type GroupService struct{}
