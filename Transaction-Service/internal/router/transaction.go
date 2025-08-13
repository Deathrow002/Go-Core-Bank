package router

import (
	"transaction-service/internal/transaction/controller"
	"transaction-service/internal/transaction/controller/middleware"

	"github.com/gin-gonic/gin"
)

func setupTrasactionRoutes(rg *gin.RouterGroup, transactionController *controller.TransactionController) {
	// Transaction routes
	transactionGroup := rg.Group("/transactions")
	{
		transactionGroup.POST("/transaction", middleware.AuthorizeRole("admin", "support", "user"),transactionController.CreateTransaction) // POST /api/v1/transactions
		transactionGroup.POST("/withdraw", middleware.AuthorizeRole("admin", "support", "user"), transactionController.CreateWithdrawalTransaction) // POST /api/v1/transactions/withdraw
		transactionGroup.POST("/deposit", middleware.AuthorizeRole("admin", "support", "user"), transactionController.CreateDepositTransaction) // POST /api/v1/transactions/deposit
		transactionGroup.GET("/:id", middleware.AuthorizeRole("admin", "support", "user"), transactionController.GetTransactionByID) // GET /api/v1/transactions/:id
		transactionGroup.GET("/", middleware.AuthorizeRole("admin", "support", "user"), transactionController.GetAllTransactions) // GET /api/v1/transactions
	}
}