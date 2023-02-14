package controllers

import (
	"net/http"
	"pay/core"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type transactionInput struct {
	destinationUUID uuid.UUID
	amount          uint
}

func (c *Controller) NewTransaction(ctx *gin.Context) {
	userUUIDstr := ctx.Param("user_uuid")
	userUUID, err := uuid.Parse(userUUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountUUIDstr := ctx.Param("account_uuid")
	accountUUID, err := uuid.Parse(accountUUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input transactionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tr := core.Transaction{
		UserUUID:        userUUID,
		SourseUUID:      accountUUID,
		DestinationUUID: input.destinationUUID,
		Amount:          input.amount,
	}
	transaction, err := c.System.NewTransaction(tr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "create new transaction", "transaction": transaction})

}
