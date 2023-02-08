package middleware

import (
	"net/http"
	_ "pay/controllers"
	"pay/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var UnauthenticatedError = gin.H{"error": "unauthenticated"}

func Auth(DB *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := models.User{}
		uuid := ctx.Param("user_uuid")
		err := DB.Model(models.User{}).Where("UUID = ?", uuid).Take(&u).Error
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
		email, ok := models.GetEmail(userTok)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		if email != u.Email {
			ctx.JSON(http.StatusUnauthorized, UnauthenticatedError)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
