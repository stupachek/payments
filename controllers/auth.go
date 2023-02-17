package controllers

import (
	"net/http"
	"payment/core"
	"payment/models"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	FisrtName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type Controller struct {
	System core.PaymentSystem
}

func NewHttpController(system core.PaymentSystem) Controller {
	return Controller{
		System: system,
	}
}

func (c *Controller) Register(ctx *gin.Context) {
	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		FisrtName: input.FisrtName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
	}
	err := c.System.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "uuid": user.UUID})
}
