package controllers

import (
	"customer-service/internal/customer/models"
	"customer-service/internal/customer/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerController struct {
	service service.CustomerService
}

func NewCustomerController(service service.CustomerService) *CustomerController {
	return &CustomerController{
		service: service,
	}
}

func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
	var req models.CustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	customer, err := c.service.CreateCustomer(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, customer)
}

func (c *CustomerController) GetCustomer(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	customer, err := c.service.GetCustomer(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var req models.CustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	customer, err := c.service.UpdateCustomer(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

func (c *CustomerController) DeleteCustomer(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	if err := c.service.DeleteCustomer(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *CustomerController) ListCustomers(ctx *gin.Context) {
	customers, err := c.service.ListCustomers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list customers"})
		return
	}

	ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerController) SearchCustomers(ctx *gin.Context) {
	var req models.CustomerSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameters"})
		return
	}

	customers, err := c.service.SearchCustomers(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search customers"})
		return
	}

	ctx.JSON(http.StatusOK, customers)
}

func (c *CustomerController) ValidateCustomer(ctx *gin.Context) {
	var req models.CustomerValidationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := c.service.ValidateCustomer(req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"valid": result})
}