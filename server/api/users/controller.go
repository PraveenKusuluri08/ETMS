package users

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Praveenkusuluri08/bootstrap"
	endpoints "github.com/Praveenkusuluri08/types"
	"github.com/Praveenkusuluri08/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection = bootstrap.GetCollection(bootstrap.ClientDB, "Users")

// @Summary		Create new user account
// @Description	Create a new user account with the provided email and password
// @Accept			json
// @Produce		json
// @Param			user	body		User	true	"User"
// @Success		201		{string}	string	"User created successfully"
// @Failure		400		{object}	endpoints.BadRequestResponse
// @Failure		500		{object}	endpoints.InternalServerResponse
// @Router			/api/v1/users/signup [post]
//
// @Tags			Users
func CreateUser() gin.HandlerFunc {
	usersService := &UsersService{}
	return usersService.CreateUser()
}

// @Summary		Sign in the user to the account
// @Description	Sign in the user to the account with the provided email and password
// @Accept			json
// @Produce		json
// @Param			user	body		UserSigninStruct	true	"User Signin"
// @Success		200		{object}	UserSigninResponse	"User Signin response"
// @Failure		400		{object}	endpoints.BadRequestResponse
// @Failure		500		{object}	endpoints.InternalServerResponse
// @Router			/api/v1/users/signin [post]
// @Tags			Users
func SignInUser() gin.HandlerFunc {
	usersService := &UsersService{}
	return usersService.SignInUser()
}

func (u *UsersService) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user UserStruct
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide data properly",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		userCounter := bson.M{"email": user.Email}
		count, err := usersCollection.CountDocuments(ctx, userCounter)
		fmt.Println(count)
		if err != nil {
			internalservererrorresponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to get count of the documents",
				},
				Status: "500",
				Error:  err.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalservererrorresponse)
			return
		}
		if count > 0 {
			badrequestresponse := &endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "User already exists",
				},
				Status: "400",
				Error:  "user_exists",
			}
			c.JSON(http.StatusBadRequest, badrequestresponse)
			return
		}

		currentTime := time.Now()
		hashedPassword := hash_password(user.Password)
		user.Password = hashedPassword
		user.CreatedAt = currentTime.Format(time.ANSIC)
		user.Uid = uuid.NewString()
		user.ID = primitive.NewObjectID()
		_, insErr := usersCollection.InsertOne(ctx, user)
		if insErr != nil {
			internalServerErrorResponse := endpoints.InternalServerResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Failed to create the user",
				},
				Status: "500",
				Error:  insErr.Error(),
			}
			c.JSON(http.StatusInternalServerError, internalServerErrorResponse)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func signin(user UserSigninStruct) (endpoints.ErrorResponse, UserSigninResponse) {
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if user.Password == "" || user.Email == "" {
		badrequestresponse := endpoints.BadRequestResponse{
			Msg: endpoints.ErrorMessage{
				Name: "Please provide data properly",
			},
			Status: "400",
			Error:  "data_not_provided",
		}

		return endpoints.ErrorResponse{
			BadRequestResponse: badrequestresponse,
		}, UserSigninResponse{}
	}
	filter_user_email := bson.M{"email": user.Email}
	count, countError := usersCollection.CountDocuments(ctx, filter_user_email)
	if countError != nil {
		internalServerErrorResponse := endpoints.InternalServerResponse{
			Msg: endpoints.ErrorMessage{
				Name: "Failed to get count of the documents",
			},
			Status: "500",
			Error:  countError.Error(),
		}
		return endpoints.ErrorResponse{
			InternalServerResponse: internalServerErrorResponse,
		}, UserSigninResponse{}
	}
	if count < 1 {
		badrequestresponse := &endpoints.BadRequestResponse{
			Msg: endpoints.ErrorMessage{
				Name: "User does not exist",
			},
			Status: "400",
			Error:  "user_does_not_exist",
		}
		return endpoints.ErrorResponse{
			BadRequestResponse: *badrequestresponse,
		}, UserSigninResponse{}
	}
	var userStruct UserStruct
	if err := usersCollection.FindOne(ctx, filter_user_email).Decode(&userStruct); err != nil {
		internalServerErrorResponse := endpoints.InternalServerResponse{
			Msg: endpoints.ErrorMessage{
				Name: "Failed to get the user details",
			},
			Status: "500",
			Error:  err.Error(),
		}
		return endpoints.ErrorResponse{
			InternalServerResponse: internalServerErrorResponse,
		}, UserSigninResponse{}
	}
	db_password := userStruct.Password
	// match the current given password with the db password
	is_password_match := utils.DecryptPassword(user.Password, db_password)
	fmt.Println(is_password_match)
	if is_password_match {
		badRequestResponse := endpoints.BadRequestResponse{
			Msg: endpoints.ErrorMessage{
				Name: "Password does not match",
			},
			Status: "400",
			Error:  "password_does_not_match",
		}

		return endpoints.ErrorResponse{
			BadRequestResponse: badRequestResponse,
		}, UserSigninResponse{}
	}
	token, tokenError := generateToken(userStruct)
	if tokenError != nil {
		internalServerResponse := endpoints.InternalServerResponse{
			Msg: endpoints.ErrorMessage{
				Name: "Failed to generate token",
			},
			Status: "500",
			Error:  tokenError.Error(),
		}
		return endpoints.ErrorResponse{
			InternalServerResponse: internalServerResponse,
		}, UserSigninResponse{}
	}
	return endpoints.ErrorResponse{}, UserSigninResponse{Token: token}
}

func (u *UsersService) SignInUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user UserSigninStruct
		if err := c.BindJSON(&user); err != nil {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide data properly",
				},
				Status: "400",
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		errorResponse, userSigninResponse := signin(user)
		fmt.Println(errorResponse)
		if errorResponse.InternalServerResponse.Error != "" {
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}
		if errorResponse.BadRequestResponse.Error != "" {
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": userSigninResponse.Token})
	}
}

func generateToken(user UserStruct) (string, error) {
	claims := &SignInDetails{
		Email: user.Email,
		Uid:   user.Uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Panicln(err)
		return "", err
	}
	return token, err
}

func hash_password(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash)
}
