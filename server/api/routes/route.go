package routes

import "github.com/gin-gonic/gin"

func SetUp(gin *gin.Engine) {
	groupRouter := gin.Group("/group/v1")
	GroupRouter(groupRouter)

	userRouter := gin.Group("/user/v1")
	UserRoutes(userRouter)
}
