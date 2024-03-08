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

	router.PUT("/update_group_name", groups.UpdateGroup())

	router.PUT("/remove_group_member", groups.RemoveGroupMember())
}
