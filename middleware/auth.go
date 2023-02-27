package middleware

import (
	"net/http"
	"payment/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var UnauthenticatedError = gin.H{"error": "unauthenticated"}
var UnkownUserError = gin.H{"error": "unknown user"}
var UserBlockedError = gin.H{"error": "user is blocked"}

func Auth(c controllers.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UUIDstr := ctx.Param("user_uuid")
		UUID, err := uuid.Parse(UUIDstr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		token := ctx.GetHeader("Authorization")
		err = c.System.CheckToken(UUID, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CheckAdmin(c controllers.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UUIDstr := ctx.Param("user_uuid")
		UUID, err := uuid.Parse(UUIDstr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}

		err = c.System.CheckAdmin(UUID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CheckBlockedUser(c controllers.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userUUIDstr := ctx.Param("user_uuid")
		userUUID, err := uuid.Parse(userUUIDstr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnkownUserError)
			ctx.Abort()
			return
		}
		ok, err := c.System.IsBlockedUser(userUUID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		if ok {
			ctx.JSON(http.StatusUnauthorized, UserBlockedError)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
