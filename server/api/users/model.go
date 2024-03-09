package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID  `bson:"_id"`
	Username string              `bson:"username"`
	Email    string              `bson:"email" validate:"required"`
	Password string              `bson:"password" validate:"required"`
	Groups   []map[string]string `bson:"groups"`
	Uid      string              `bson:"uid"`
	Role     string              `bson:"role" validate:"required"`
	Token    string              `bson:"token"`
}
