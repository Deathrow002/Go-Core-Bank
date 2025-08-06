package service

import (
	"customer-service/internal/customer/models"

	"github.com/google/uuid"
)

type CustomerService interface {
	CreateCustomer(req models.CustomerRequest) (*models.Customer, error)
	GetCustomer(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(id uuid.UUID, req models.CustomerRequest) (*models.Customer, error)
	DeleteCustomer(id uuid.UUID) error
	ListCustomers() ([]models.Customer, error)
	SearchCustomers(req models.CustomerSearchRequest) ([]models.Customer, error)
	ValidateCustomer(req models.CustomerValidationRequest) (bool, error)
}