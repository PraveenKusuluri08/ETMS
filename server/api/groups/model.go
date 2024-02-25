package groups

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID        primitive.ObjectID  `bson:"_id"`
	GroupName string              `bson:"group_name" validate:"required"`
	Users     []map[string]string `bson:"users" validate:"required"`
	Type      string              `bson:"type" validate:"required"`
}

type Invitation struct {
	GroupName string   `json:"groupname"`
	Users     []string `json:"users"`
}
