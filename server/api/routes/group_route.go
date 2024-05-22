package routes

import (
	"github.com/Praveenkusuluri08/api/groups"
	"github.com/gin-gonic/gin"
)

func GroupRouter(router *gin.RouterGroup) {
	var groupInterface groups.GroupInterface = &groups.GroupService{}
	router.POST("/creategroup", groupInterface.CreateGroup())

	router.POST("/invite", groups.InviteGroupMembers())

	router.POST("/accept_invitation", groupInterface.AcceptInvitation())

	router.POST("/get_users", groupInterface.DisplayUsers())

	router.PUT("/update_group_name", groupInterface.UpdateGroup())

	router.PUT("/remove_group_member", groupInterface.RemoveGroupMember())
}
