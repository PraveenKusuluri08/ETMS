package endpoints

import (
	"net/http"

	"github.com/Praveenkusuluri08/endpoints"
	"github.com/Praveenkusuluri08/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: "Please provide token",
				},
				Status: "400",
				Error:  "token_not_provided",
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		claims, err := utils.ValidateToken(token)
		if err != "" {
			badRequestResponse := endpoints.BadRequestResponse{
				Msg: endpoints.ErrorMessage{
					Name: err,
				},
				Status: "400",
				Error:  err,
			}
			c.JSON(http.StatusBadRequest, badRequestResponse)
			return
		}
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
