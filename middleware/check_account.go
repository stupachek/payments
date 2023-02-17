package middleware

import (
	"net/http"
	"payment/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var AccountError = gin.H{"error": "wrong account"}

func CheckAccount(c controllers.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userUUIDstr := ctx.Param("user_uuid")
		userUUID, err := uuid.Parse(userUUIDstr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		accountUUIDstr := ctx.Param("account_uuid")
		accountUUID, err := uuid.Parse(accountUUIDstr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, AccountError)
			ctx.Abort()
			return
		}
		err = c.System.CheckAccountExists(userUUID, accountUUID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, AccountError)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}