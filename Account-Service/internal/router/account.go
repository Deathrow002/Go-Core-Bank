package router

import (
	"account-service/internal/account/controllers"
	"account-service/internal/account/controllers/middleware"

	"github.com/gin-gonic/gin"
)

// setupAccountRoutes configures account-specific routes
func setupAccountRoutes(rg *gin.RouterGroup, accountController *controllers.AccountController) {
    accounts := rg.Group("/accounts")
    {
        // Basic CRUD operations
        accounts.POST("", middleware.AuthorizeRole("admin", "support", "user"),accountController.CreateAccount)
        accounts.GET("/:id", middleware.AuthorizeRole("admin", "support"), accountController.GetAccount)
        accounts.PUT("/:id", middleware.AuthorizeRole("admin", "support", "user"), accountController.UpdateAccount)
        accounts.DELETE("/:id", middleware.AuthorizeRole("admin", "support", "user"), accountController.DeleteAccount)

        // Listing and search
        accounts.GET("", middleware.AuthorizeRole("admin", "support"), middleware.AuthorizeRole("admin", "support", "user"), accountController.ListAccounts)
        accounts.GET("/search", middleware.AuthorizeRole("admin", "support"), accountController.SearchAccounts)

        // Special queries
        accounts.GET("/number/:account_number", middleware.AuthorizeRole("admin", "support"), accountController.GetAccountByNumber)
        accounts.GET("/customer/:customer_id", middleware.AuthorizeRole("admin", "support"), accountController.GetAccountsByCustomer)
    }
}