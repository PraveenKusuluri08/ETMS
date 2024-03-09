package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Praveenkusuluri08/api/users"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type SignInDetails struct {
	Email string
	Uid   string
	Role  string
	jwt.StandardClaims
}

func GenerateToken(user users.User) (string, error) {
	claims := &SignInDetails{
		Email: user.Email,
		Role:  user.Role,
		Uid:   user.Uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panicln(err)
		return "", err
	}
	return token, err
}

func ValidateToken(token string) (claims *SignInDetails, msg string) {
	var message string
	tokenString, err := jwt.ParseWithClaims(
		token,
		&SignInDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	claims, ok := tokenString.Claims.(*SignInDetails)

	if !ok {
		message = fmt.Sprintf("token is expired")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		message = fmt.Sprintf("Token is expired please check")
		return
	}
	fmt.Println(message)
	return claims, message

}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 20)
	return string(hash), err
}

func DecryptPassword(password, hash string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", err
	}
	return password, nil
}
