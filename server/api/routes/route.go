package routes

import "github.com/gin-gonic/gin"

func SetUp(gin *gin.Engine) {
	publicRouter := gin.Group("")
	GroupRouter(publicRouter)
}
