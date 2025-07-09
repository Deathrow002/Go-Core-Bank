package service

import (
	"customer-service/internal/customer/models"
	"customer-service/internal/customer/repository"
	"errors"
	"math"

	"github.com/google/uuid"
)

// CustomerService defines the interface for customer business logic
type CustomerService interface {
	CreateCustomer(req models.CustomerRequest) (*models.CustomerResponse, error)
	GetCustomer(id uuid.UUID) (*models.CustomerResponse, error)
	UpdateCustomer(id uuid.UUID, req models.CustomerRequest) (*models.CustomerResponse, error)
	DeleteCustomer(id uuid.UUID) error
	ListCustomers(page, pageSize int) (*models.CustomerListResponse, error)
	SearchCustomers(req models.CustomerSearchRequest) (*models.CustomerListResponse, error)
}

type customerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService creates a new customer service instance
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{
		repo: repo,
	}
}

// CreateCustomer creates a new customer
func (s *customerService) CreateCustomer(req models.CustomerRequest) (*models.CustomerResponse, error) {
	// Validate business rules
	if err := s.validateCustomerRequest(req); err != nil {
		return nil, err
	}

	// Check if customer with email already exists
	existingCustomer, _ := s.repo.GetByEmail(req.Email)
	if existingCustomer != nil {
		return nil, errors.New("customer with this email already exists")
	}

	// Create customer model
	customer := &models.Customer{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		DateOfBirth: req.DateOfBirth,
		Address:     req.Address,
		Status:      models.CustomerStatusActive,
	}

	// Save to database
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}

	// Convert to response
	response := customer.ToResponse()
	return &response, nil
}

// GetCustomer retrieves a customer by ID
func (s *customerService) GetCustomer(id uuid.UUID) (*models.CustomerResponse, error) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := customer.ToResponse()
	return &response, nil
}

// UpdateCustomer updates an existing customer
func (s *customerService) UpdateCustomer(id uuid.UUID, req models.CustomerRequest) (*models.CustomerResponse, error) {
	// Validate business rules
	if err := s.validateCustomerRequest(req); err != nil {
		return nil, err
	}

	// Get existing customer
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if new email already exists
	if customer.Email != req.Email {
		existingCustomer, _ := s.repo.GetByEmail(req.Email)
		if existingCustomer != nil {
			return nil, errors.New("customer with this email already exists")
		}
	}

	// Update customer fields
	customer.FirstName = req.FirstName
	customer.LastName = req.LastName
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.DateOfBirth = req.DateOfBirth
	customer.Address = req.Address

	// Save changes
	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}

	// Convert to response
	response := customer.ToResponse()
	return &response, nil
}

// DeleteCustomer deletes a customer
func (s *customerService) DeleteCustomer(id uuid.UUID) error {
	// Check if customer exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Perform soft delete
	return s.repo.Delete(id)
}

// ListCustomers lists customers with pagination
func (s *customerService) ListCustomers(page, pageSize int) (*models.CustomerListResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Limit maximum page size
	}

	customers, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	customerResponses := make([]models.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = customer.ToResponse()
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.CustomerListResponse{
		Customers:  customerResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchCustomers searches customers based on criteria
func (s *customerService) SearchCustomers(req models.CustomerSearchRequest) (*models.CustomerListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // Limit maximum page size
	}

	customers, total, err := s.repo.Search(req)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	customerResponses := make([]models.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = customer.ToResponse()
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &models.CustomerListResponse{
		Customers:  customerResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// validateCustomerRequest validates the customer request
func (s *customerService) validateCustomerRequest(req models.CustomerRequest) error {
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if req.LastName == "" {
		return errors.New("last name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Phone == "" {
		return errors.New("phone is required")
	}
	
	// Add more validation rules as needed
	return nil
}
