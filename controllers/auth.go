package controllers

import (
	"net/http"
	"payment/core"
	"payment/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		Role:      core.USER,
	}
	err := c.System.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "uuid": user.UUID})
}

func (c *Controller) ChangeRole(ctx *gin.Context) {
	UUIDstr := ctx.Param("user_uuid")
	adminUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input ChangeRoleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !(input.Role == core.USER || input.Role == core.ADMIN) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": UnknownRoleError})
		return
	}
	userUUID, err := uuid.Parse(input.UserUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": BadRequestError})
		return
	}
	err = c.System.ChangeRole(adminUUID, userUUID, input.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "change role"})

}

func (c *Controller) BlockUser(ctx *gin.Context) {
	UUIDstr := ctx.Param("target_uuid")
	userUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.System.BlockUser(userUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user is blocked"})

}

func (c *Controller) UnblockUser(ctx *gin.Context) {
	UUIDstr := ctx.Param("target_uuid")
	userUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.System.UnblockUser(userUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user is active"})

}
