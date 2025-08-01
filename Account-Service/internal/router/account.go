package router

import (
	"account-service/internal/account/controllers"

	"github.com/gin-gonic/gin"
)

// setupAccountRoutes configures account-specific routes
func setupAccountRoutes(rg *gin.RouterGroup, accountController *controllers.AccountController) {
    accounts := rg.Group("/accounts")
    {
        // Basic CRUD operations
        accounts.POST("", accountController.CreateAccount)
        accounts.GET("/:id", accountController.GetAccount)
        accounts.PUT("/:id", accountController.UpdateAccount)
        accounts.DELETE("/:id", accountController.DeleteAccount)

        // Listing and search
        accounts.GET("", accountController.ListAccounts)
        accounts.GET("/search", accountController.SearchAccounts)

        // Special queries
        accounts.GET("/number/:account_number", accountController.GetAccountByNumber)
        accounts.GET("/customer/:customer_id", accountController.GetAccountsByCustomer)
    }
}