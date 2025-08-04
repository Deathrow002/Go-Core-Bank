package router

import (
	"authentication-service/internal/authentication/controller"
	"authentication-service/internal/authentication/controller/middleware"

	"github.com/gin-gonic/gin"
)

func setupAuthenticationRoutes(rg *gin.RouterGroup, authenticationController *controller.AuthenticationController){

	// Authentication routes
	authGroup := rg.Group("/auth")
	{
		{
			authGroup.POST("/authenticate", middleware.AuthorizeRole("admin", "support", "user"), authenticationController.CreateAuthentication) // POST /api/v1/auth/authenticate
			authGroup.POST("/login", middleware.AuthorizeRole("admin", "support", "user"), authenticationController.Login) // POST /api/v1/auth/login
			authGroup.POST("/change-password", middleware.AuthorizeRole("admin", "support", "user"), authenticationController.ChangePassword) // POST /api/v1
			
			authGroup.GET("/email/:email", middleware.AuthorizeRole("admin", "support"),authenticationController.GetByEmail) // GET /api
			authGroup.GET("/customer/:customerID", middleware.AuthorizeRole("admin", "support"), authenticationController.GetByCustomerID) // GET
			authGroup.GET("/id/:id", middleware.AuthorizeRole("admin", "support"), authenticationController.GetByCustomerID) // GET /api/v1/auth/id/:id
			authGroup.PUT("/update/:id", middleware.AuthorizeRole("admin", "support"), authenticationController.UpdateAuthentication) // PUT /api/v
			authGroup.DELETE("/delete/:id", middleware.AuthorizeRole("admin", "support"), authenticationController.DeleteAuthentication) // DELETE /api/v1/auth/delete/:id
			
			authGroup.POST("/lock/:id", middleware.AuthorizeRole("admin", "support"), authenticationController.LockAccount) // POST /api/v
			authGroup.POST("/unlock/:id", middleware.AuthorizeRole("admin", "support"), authenticationController.UnlockAccount) // POST /api/v1/auth/unlock/:id
		}
	}
}