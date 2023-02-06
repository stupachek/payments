package middleware

import (
	"net/http"
	"pay/models"

	"github.com/gin-gonic/gin"
)

var UnauthenticatedError = gin.H{"error": "unauthenticated"}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := models.User{}
		uuid := c.Param("user_uuid")
		err := models.DB.Model(models.User{}).Where("UUID = ?", uuid).Take(&u).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, UnauthenticatedError)
			c.Abort()
			return
		}
		tok, ok := models.GetToken(u.Email)
		if !ok {
			c.JSON(http.StatusUnauthorized, UnauthenticatedError)
			c.Abort()
			return
		}
		userTok := c.GetHeader("Authorization")
		if userTok == "" {
			c.JSON(http.StatusUnauthorized, UnauthenticatedError)
			c.Abort()
			return
		}
		if tok != userTok {
			c.JSON(http.StatusUnauthorized, UnauthenticatedError)
			c.Abort()
			return
		}
		c.Next()
	}
}
