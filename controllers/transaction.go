package controllers

import (
	"net/http"
	"pay/core"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionInput struct {
	DestinationUUID string `json:"destination_uuid" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
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
	var input TransactionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	destinationUUID, err := uuid.Parse(input.DestinationUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	amount, err := strconv.ParseUint(input.Amount, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tr := core.Transaction{
		UserUUID:        userUUID,
		SourceUUID:      accountUUID,
		DestinationUUID: destinationUUID,
		Amount:          uint(amount),
	}
	transaction, err := c.System.NewTransaction(tr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "create new transaction", "transaction": transaction})

}

func (c *Controller) GetTransactions(ctx *gin.Context) {
	accountUUIDstr := ctx.Param("account_uuid")
	accountUUID, err := uuid.Parse(accountUUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transactions, err := c.System.GetTransactions(accountUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"transactions": transactions})

}

func (c *Controller) SendTransaction(ctx *gin.Context) {
	transactionUUIDstr := ctx.Param("transaction_uuid")
	transactionUUID, err := uuid.Parse(transactionUUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.System.SendTransaction(transactionUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sent transaction"})

}
