package controllers

import (
	"net/http"
	"pay/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountInput struct {
	IBAN string `json:"iban" binding:"required,len=29"`
}

func (c *Controller) NewAccount(ctx *gin.Context) {
	var input AccountInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	UUIDstr := ctx.Param("user_uuid")
	userUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account := models.Account{
		IBAN: input.IBAN,
	}
	c.System.NewAccount(userUUID, &account)
	ctx.JSON(http.StatusOK, gin.H{"message": "new account add", "uuid": account.UUID})

}
