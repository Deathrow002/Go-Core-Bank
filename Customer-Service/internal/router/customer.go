package router

import (
	"customer-service/internal/customer/controllers"
	"customer-service/internal/customer/controllers/middleware"

	"github.com/gin-gonic/gin"
)

func setupCustomerRoutes(rg *gin.RouterGroup, customerController *controllers.CustomerController) {
	customerGroup := rg.Group("/customers")
	{
		customerGroup.POST("/", middleware.AuthorizeRole("admin", "support", "user"), customerController.CreateCustomer)
		customerGroup.GET("/:id", middleware.AuthorizeRole("admin", "support", "user"), customerController.GetCustomer)
		customerGroup.PUT("/:id", middleware.AuthorizeRole("admin", "support", "user"), customerController.UpdateCustomer)
		customerGroup.DELETE("/:id", middleware.AuthorizeRole("admin", "support", "user"), customerController.DeleteCustomer)
		customerGroup.GET("/", middleware.AuthorizeRole("admin", "support", "user"), customerController.ListCustomers)
		
		customerGroup.GET("/search", middleware.AuthorizeRole("admin", "support", "user"), customerController.SearchCustomers)
		customerGroup.POST("/validate", customerController.ValidateCustomer)
	}
}