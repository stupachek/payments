package middleware

import (
	"net/http"
	"pay/models"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := models.User{}
		uuid := c.Param("user_uuid")
		err := models.DB.Model(models.User{}).Where("UUID = ?", uuid).Take(&u).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			c.Abort()
			return
		}
		tok, ok := models.GetToken(u.Email)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user does not have token"})
			c.Abort()
			return
		}
		userTok := c.GetHeader("Authorization")
		if userTok == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}
		if tok != userTok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			c.Abort()
			return
		}
		c.Next()
	}
}
