package routes

import (
	"github.com/Praveenkusuluri08/api/users"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup){
	
	router.POST("/signup",users.CreateUser())
}