package middleware

import (
	"net/http"
	"pay/controllers"
	_ "pay/controllers"
	"pay/core"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var UnauthenticatedError = gin.H{"error": "unauthenticated"}

func Auth(c controllers.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UUIDstr := ctx.Param("user_uuid")
		UUID, err := uuid.FromBytes([]byte(UUIDstr))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		user, err := c.System.UserRepo.GetUserUUID(ctx, UUID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		userTok := ctx.GetHeader("Authorization")
		if userTok == "" {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		email, ok := core.GetEmail(userTok)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		if email != user.Email {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
