package controller

import (
	"net/http"
	"transaction-service/internal/transaction/models"
	"transaction-service/internal/transaction/service"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service service.TransactionService
}

func NewTransactionController(service service.TransactionService) *TransactionController {
	return &TransactionController{
		service: service,
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	// Get and sanitize the token
    authHeader := ctx.GetHeader("Authorization")
    token := authHeader
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        token = authHeader[7:]
    }

	transaction, err := c.service.CreateTransaction(ctx.Request.Context(), req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transaction)
}

func (c *TransactionController) CreateWithdrawalTransaction(ctx *gin.Context) {
	var req models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	// Get and sanitize the token
    authHeader := ctx.GetHeader("Authorization")
    token := authHeader
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        token = authHeader[7:]
    }

	transaction, err := c.service.CreateWithdrawalTransaction(ctx.Request.Context(), req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transaction)
}

func (c *TransactionController) CreateDepositTransaction(ctx *gin.Context) {
	var req models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "message": err.Error()})
		return
	}

	transaction, err := c.service.CreateDepositTransaction(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transaction)
}

func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	transaction, err := c.service.GetTransactionByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found", "message": "Transaction with the given ID does not exist"})
		return
	}
	ctx.JSON(http.StatusOK, transaction)
}

func (c *TransactionController) GetAllTransactions(ctx *gin.Context) {
	transactions, err := c.service.ListAllTransactions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list transactions", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transactions)
}