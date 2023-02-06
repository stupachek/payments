package controllers

import (
	"net/http"
	"pay/models"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	FisrtName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		FisrtName: input.FisrtName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
	}
	_, err := user.CreateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration success", "uuid": user.UUID})
}
