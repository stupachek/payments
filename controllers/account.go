package controllers

import (
	"net/http"
	"payment/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddMoneyInput struct {
	Amount string `json:"amount" binding:"required"`
}

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

func query(ctx *gin.Context) (models.QueryParams, error) {
	limit, err := strconv.ParseUint(ctx.DefaultQuery("limit", "30"), 10, 32)
	if err != nil {
		return models.QueryParams{}, err
	}
	if limit > 30 {
		limit = 30
	}
	offset, err := strconv.ParseUint(ctx.DefaultQuery("offset", "0"), 10, 32)
	if err != nil {
		return models.QueryParams{}, err
	}
	return models.QueryParams{
		Limit:  uint(limit),
		Offset: uint(offset),
	}, nil
}

func (c *Controller) GetAccounts(ctx *gin.Context) {
	UUIDstr := ctx.Param("user_uuid")
	userUUID, err := uuid.Parse(UUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, err := query(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts, err := c.System.GetAccounts(userUUID, query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"accounts": accounts})

}

func (c *Controller) AddMoney(ctx *gin.Context) {
	accountUUIDstr := ctx.Param("account_uuid")
	accountUUID, err := uuid.Parse(accountUUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input AddMoneyInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := strconv.ParseUint(input.Amount, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := c.System.AddMoney(accountUUID, uint(amount))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "add money", "account": account})

}

func (c *Controller) GetAccount(ctx *gin.Context) {
	accountUUIDstr := ctx.Param("account_uuid")
	accountUUID, err := uuid.Parse(accountUUIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := c.System.GetAccount(accountUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"uuid": account.UUID, "iban": account.IBAN, "balance": account.Balance})

}
