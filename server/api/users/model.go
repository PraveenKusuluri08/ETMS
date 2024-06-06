package users

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStruct struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Username  string              `bson:"username"`
	Email     string              `bson:"email" validate:"required"`
	Password  string              `bson:"password" validate:"required"`
	Groups    []map[string]string `bson:"groups"`
	Uid       string              `bson:"uid"`
	CreatedAt string              `bson:"created_at"`
}

type User struct {
	Username string              `bson:"username"`
	Email    string              `bson:"email" validate:"required"`
	Groups   []map[string]string `bson:"groups"`
	Password string              `bson:"password" validate:"required"`
}

type SignInDetails struct {
	Email string
	Uid   string
	jwt.StandardClaims
}

type UserSigninStruct struct {
	Email    string `json:"email" valudate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserSigninResponse struct {
	Token string `json:"token"`
}

type UsersInterface interface {
	CreateUser() gin.HandlerFunc
	SignInUser() gin.HandlerFunc
}

type UsersService struct{}
