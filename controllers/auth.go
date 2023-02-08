package controllers

import (
	"net/http"
	"pay/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RegisterInput struct {
	FisrtName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type Controller struct {
	DB *gorm.DB
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
	_, err := user.CreateUser(c.DB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "uuid": user.UUID})
}
