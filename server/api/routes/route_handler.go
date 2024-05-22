package routes

import "github.com/gin-gonic/gin"

func SetUp(r *gin.Engine) {
	group_router := r.Group("/api/v1/groups")
	GroupRouter(group_router)
	user_router := r.Group("/api/v1/users")
	UserRoutes(user_router)
}
