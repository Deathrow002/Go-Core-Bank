package handlers

import (
	"customer-service/internal/customer/models"
	"customer-service/internal/customer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CustomerHandler handles HTTP requests for customer operations
type CustomerHandler struct {
	service service.CustomerService
}

// NewCustomerHandler creates a new customer handler instance
func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

// CreateCustomer creates a new customer
// @Summary Create a new customer
// @Description Create a new customer with the provided information
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body models.CustomerRequest true "Customer information"
// @Success 201 {object} models.CustomerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req models.CustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	customer, err := h.service.CreateCustomer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create customer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomer retrieves a customer by ID
// @Summary Get customer by ID
// @Description Get a customer's information by their ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID" Format(uuid)
// @Success 200 {object} models.CustomerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid customer ID",
			Message: "Customer ID must be a valid UUID",
		})
		return
	}

	customer, err := h.service.GetCustomer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Customer not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer updates an existing customer
// @Summary Update customer
// @Description Update a customer's information
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID" Format(uuid)
// @Param customer body models.CustomerRequest true "Updated customer information"
// @Success 200 {object} models.CustomerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid customer ID",
			Message: "Customer ID must be a valid UUID",
		})
		return
	}

	var req models.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	customer, err := h.service.UpdateCustomer(id, req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "customer not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to update customer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer deletes a customer
// @Summary Delete customer
// @Description Delete a customer by ID
// @Tags customers
// @Param id path string true "Customer ID" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid customer ID",
			Message: "Customer ID must be a valid UUID",
		})
		return
	}

	if err := h.service.DeleteCustomer(id); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "customer not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to delete customer",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCustomers lists customers with pagination
// @Summary List customers
// @Description Get a paginated list of customers
// @Tags customers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.CustomerListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/customers [get]
func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	customers, err := h.service.ListCustomers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list customers",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// SearchCustomers searches customers based on criteria
// @Summary Search customers
// @Description Search customers by query, status, and other criteria
// @Tags customers
// @Produce json
// @Param query query string false "Search query"
// @Param status query string false "Customer status" Enums(active, inactive, suspended, closed)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.CustomerListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/customers/search [get]
func (h *CustomerHandler) SearchCustomers(c *gin.Context) {
	var req models.CustomerSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid search parameters",
			Message: err.Error(),
		})
		return
	}

	customers, err := h.service.SearchCustomers(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to search customers",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
