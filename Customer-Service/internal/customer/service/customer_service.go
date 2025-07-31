package service

import (
	"customer-service/internal/customer/models"
	"customer-service/internal/customer/repository"
	"errors"

	"github.com/google/uuid"
)



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
func (s *customerService) CreateCustomer(req models.CustomerRequest) (*models.Customer, error) {
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

	// Return the created customer
	return customer, nil
}

// GetCustomer retrieves a customer by ID
func (s *customerService) GetCustomer(id uuid.UUID) (*models.Customer, error) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

// UpdateCustomer updates an existing customer
func (s *customerService) UpdateCustomer(id uuid.UUID, req models.CustomerRequest) (*models.Customer, error) {
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

	// Return the updated customer
	return customer, nil
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

// ListCustomers lists all customers (removed pagination)
func (s *customerService) ListCustomers() ([]models.Customer, error) {
	customers, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// ValidateCustomer checks if a customer exists and returns boolean
func (s *customerService) ValidateCustomer(req models.CustomerValidationRequest) (bool, error) {
	err := s.repo.ValidateCustomer(req)
	if err != nil {
		if err.Error() == "customer not found" {
			return false, nil // Customer doesn't exist, but no error occurred
		}
		return false, err // Actual error occurred
	}
	return true, nil // Customer exists
}

// SearchCustomers searches customers based on criteria (removed pagination)
func (s *customerService) SearchCustomers(req models.CustomerSearchRequest) ([]models.Customer, error) {
	customers, err := s.repo.Search(req)
	if err != nil {
		return nil, err
	}

	return customers, nil
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
