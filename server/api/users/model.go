package users

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStruct struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Username  string              `bson:"username"`
	Email     string              `bson:"email" validate:"required"`
	Password  string              `bson:"password" validate:"required"`
	Groups    []map[string]string `bson:"groups"`
	Uid       string              `bson:"uid"`
	Token     string              `bson:"token"`
	CreatedAt string              `bson:"created_at"`
}

type SignInDetails struct {
	Email string
	Uid   string
	jwt.StandardClaims
}
