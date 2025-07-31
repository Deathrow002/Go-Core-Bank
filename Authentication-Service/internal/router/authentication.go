package router

import (
	"authentication-service/internal/authentication/controller"

	"github.com/gin-gonic/gin"
)

func setupAuthenticationRoutes(rg *gin.RouterGroup, authenticationController *controller.AuthenticationController){

	// Authentication routes
	authGroup := rg.Group("/auth")
	{
		{
			authGroup.POST("/authenticate", authenticationController.CreateAuthentication) // POST /api/v1/auth/authenticate
			authGroup.POST("/login", authenticationController.Login) // POST /api/v1/auth/login
			authGroup.GET("/email/:email", authenticationController.GetByEmail) // GET /api
			authGroup.GET("/customer/:customerID", authenticationController.GetByCustomerID) // GET
			authGroup.GET("/id/:id", authenticationController.GetByCustomerID) // GET /api/v1/auth/id/:id
			authGroup.PUT("/update/:id", authenticationController.UpdateAuthentication) // PUT /api/v
			authGroup.DELETE("/delete/:id", authenticationController.DeleteAuthentication) // DELETE /api/v1/auth/delete/:id
			authGroup.POST("/lock/:id", authenticationController.LockAccount) // POST /api/v
			authGroup.POST("/unlock/:id", authenticationController.UnlockAccount) // POST /api/v1/auth/unlock/:id
			authGroup.POST("/change-password", authenticationController.ChangePassword) // POST /api/v1
		}
	}
}