package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) NewAccount(ctx *gin.Context) {
	UUIDstr := ctx.Param("user_uuid")
	userUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := c.System.NewAccount(userUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "new account add", "uuid": account.UUID})

}

func (c *Controller) GetAccounts(ctx *gin.Context) {
	UUIDstr := ctx.Param("user_uuid")
	userUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accounts, err := c.System.GetAccounts(userUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"accounts": accounts})

}
