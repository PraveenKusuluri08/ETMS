package routes

import (
	"github.com/Praveenkusuluri08/api/groups"
	"github.com/gin-gonic/gin"
)

func GroupRouter(router *gin.RouterGroup) {

	router.POST("/creategroup", groups.CreateGroup())

	router.POST("/invite", groups.InviteGroupMembers())

	router.POST("/accept_invitation", groups.AcceptInvitation())

	router.POST("/get_users", groups.DisplaUsers())

	router.POST("/update_group_name", groups.UpdateGroup())
}
