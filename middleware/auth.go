package middleware

import (
	"net/http"
	"pay/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var UnauthenticatedError = gin.H{"error": "unauthenticated"}

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
		err = c.System.CheckToken(ctx, UUID, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
