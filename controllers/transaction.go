package controllers

import (
	"net/http"
	"payment/core"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CREATED = "created_at"
	UPDATED = "updated_at"
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
	query, err := query(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": UnknownQueryError})
		return
	}
	sort_by := ctx.DefaultQuery("sort_by", "uuid")
	sort_by = strings.ToLower(sort_by)
	order := ctx.DefaultQuery("order", "asc")
	order = strings.ToLower(order)

	if !(sort_by == UUID || sort_by == CREATED || sort_by == UPDATED) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": UnknownQueryError})
		return
	}
	if !(order == DESC || order == ASC) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": UnknownQueryError})
		return
	}
	query.Sort = sort_by + " " + order
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transactions, err := c.System.GetTransactions(accountUUID, query)
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
	transaction, err := c.System.SendTransaction(transactionUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sent transaction", "transaction": transaction})

}
