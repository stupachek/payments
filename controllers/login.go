package controllers

import (
	"net/http"
	_ "payment/core"
	"payment/models"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *Controller) Login(ctx *gin.Context) {
	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := models.User{
		Email:    input.Email,
		Password: input.Password,
	}
	out, err := c.System.LoginCheck(u.Email, u.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uuid": out.UUID, "token": out.Token})

}
