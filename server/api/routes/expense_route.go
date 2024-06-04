package routes

import (
	"github.com/Praveenkusuluri08/api/expenses"
	"github.com/Praveenkusuluri08/middlewares"
	"github.com/gin-gonic/gin"
)

func ExpenseRouter(router *gin.RouterGroup) {
	router.Use(middlewares.AuthMiddleware())
	router.POST("/create", expenses.CreateExpense())
}
