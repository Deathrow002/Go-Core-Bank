package controller

import (
	"authentication-service/internal/authentication/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	service service.AuthenticationService
}

func NewAuthenticationController(service service.AuthenticationService) *AuthenticationController {
	return &AuthenticationController{
		service: service,
	}
}

// CreateAuthentication creates a new authentication record
func (c *AuthenticationController) CreateAuthentication(ctx *gin.Context) {
	var req service.CreateAuthenticationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Pass gin context and the correct request type
	authentication, err := c.service.CreateAuthentication(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, authentication)
}

// Login validates user credentials
func (c *AuthenticationController) Login(ctx *gin.Context) {
	var req service.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authentication, err := c.service.ValidateLogin(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authentication)
}

// GetByEmail retrieves authentication by email
func (c *AuthenticationController) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	authentication, err := c.service.GetByEmail(ctx.Request.Context(), email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authentication)
}

// GetByCustomerID retrieves authentication by customer ID
func (c *AuthenticationController) GetByCustomerID(ctx *gin.Context) {
	customerID := ctx.Param("customer_id")
	if customerID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	authentication, err := c.service.GetByCustomerID(ctx.Request.Context(), customerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authentication)
}

// UpdateAuthentication updates an existing authentication record
func (c *AuthenticationController) UpdateAuthentication(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req service.UpdateAuthenticationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authentication, err := c.service.UpdateAuthentication(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authentication)
}

// LockAccount locks a user account
func (c *AuthenticationController) LockAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := c.service.LockAccount(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account locked successfully"})
}

// UnlockAccount unlocks a user account
func (c *AuthenticationController) UnlockAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := c.service.UnlockAccount(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account unlocked successfully"})
}

// ChangePassword changes user password
func (c *AuthenticationController) ChangePassword(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	var req service.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.service.ChangePassword(ctx.Request.Context(), username, req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// DeleteAuthentication deletes an authentication record
func (c *AuthenticationController) DeleteAuthentication(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := c.service.DeleteAuthentication(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
