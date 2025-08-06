package router

import (
	"customer-service/internal/customer/controllers"

	"github.com/gin-gonic/gin"
)

func setupCustomerRoutes(rg *gin.RouterGroup, customerController *controllers.CustomerController) {
	customerGroup := rg.Group("/customers")
	{
		customerGroup.POST("/", customerController.CreateCustomer)
		customerGroup.GET("/:id", customerController.GetCustomer)
		customerGroup.PUT("/:id", customerController.UpdateCustomer)
		customerGroup.DELETE("/:id", customerController.DeleteCustomer)
		customerGroup.GET("/", customerController.ListCustomers)
		
		customerGroup.GET("/search", customerController.SearchCustomers)
		customerGroup.POST("/validate", customerController.ValidateCustomer)
	}
}