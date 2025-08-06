package controllers

import (
	"account-service/internal/account/models"
	"account-service/internal/account/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountController struct {
	service service.AccountService
}

func NewAccountController(service service.AccountService) *AccountController {
	return &AccountController{
		service: service,
	}
}

func (c *AccountController) CreateAccount(ctx *gin.Context) {

	var req models.AccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	account, err := c.service.CreateAccount(ctx, req, ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create account",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (c *AccountController) GetAccount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	account, err := c.service.GetAccount(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Account not found",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) GetAccountsByCustomer(ctx *gin.Context) {
	customerIDStr := ctx.Param("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid customer ID",
			Message: err.Error(),
		})
		return
	}

	accounts, total, err := c.service.GetAccountsByCustomer(ctx, customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve accounts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  accounts,
		"total": total,
	})
}

func (c *AccountController) GetAccountByNumber(ctx *gin.Context) {
	accountNumber := ctx.Param("account_number")
	if accountNumber == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Account number is required",
			Message: "Please provide a valid account number",
		})
		return
	}

	account, err := c.service.GetAccountByNumber(ctx, accountNumber)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Account not found",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) SearchAccounts(ctx *gin.Context) {
	query := ctx.Query("query")

	// Parse account type filter
	var accountType *models.AccountType
	if at := ctx.Query("account_type"); at != "" {
		atype := models.AccountType(at)
		accountType = &atype
	}

	// Parse status filter
	var status *models.AccountStatus
	if s := ctx.Query("status"); s != "" {
		stat := models.AccountStatus(s)
		status = &stat
	}

	accounts, total, err := c.service.SearchAccounts(ctx, query, accountType, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to search accounts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accounts":    accounts,
		"total":       total,
	})
}

func (c *AccountController) UpdateAccount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	var req models.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	account, err := c.service.UpdateAccount(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update account",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) DeleteAccount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	if err := c.service.DeleteAccount(ctx, id); err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Failed to delete account",
			Message: err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *AccountController) ListAccounts(ctx *gin.Context) {
	accounts, total, err := c.service.ListAccounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve accounts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  accounts,
		"total": total,
	})
}

func (c *AccountController) ValidateAccount(ctx *gin.Context) {
	var req models.AccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	Validated , err := c.service.ValidateAccount(ctx, req.AccountNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid": Validated,
	})
}

func (c *AccountController) UpdateBalance(ctx *gin.Context) {
	var req models.AccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if err := c.service.UpdateBalance(ctx, req.AccountNumber, int64(req.InitialBalance* 100)); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update balance",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Balance updated successfully",
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}