package repository

import (
	"customer-service/internal/customer/models"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	Create(customer *models.Customer) error
	GetByID(id uuid.UUID) (*models.Customer, error)
	GetByEmail(email string) (*models.Customer, error)
	ValidateCustomer(req models.CustomerValidationRequest) (bool, error)
	Update(customer *models.Customer) error
	Delete(id uuid.UUID) error
	List() ([]models.Customer, error) // Remove pagination parameters
	Search(req models.CustomerSearchRequest) ([]models.Customer, error) // Remove pagination
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

// Create creates a new customer record
func (r *customerRepository) Create(customer *models.Customer) error {
	if err := r.db.Create(customer).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("customer with this email already exists")
		}
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

// GetByID retrieves a customer by ID
func (r *customerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.Where("id = ?", id).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

// GetByEmail retrieves a customer by email
func (r *customerRepository) GetByEmail(email string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.Where("email = ?", email).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

// ValidateCustomer checks if customer exists
func (r *customerRepository) ValidateCustomer(req models.CustomerValidationRequest) (bool, error) {
	var customer models.Customer
	if err := r.db.Where("ID = ?", req.ID).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("customer not found")
		}
		return false, fmt.Errorf("failed to get customer: %w", err)
	}
	return true, nil
}

// Update updates an existing customer record
func (r *customerRepository) Update(customer *models.Customer) error {
	if err := r.db.Save(customer).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("customer with this email already exists")
		}
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

// Delete soft deletes a customer record
func (r *customerRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Customer{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete customer: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}
	return nil
}

// List retrieves all customers (removed pagination)
func (r *customerRepository) List() ([]models.Customer, error) {
	var customers []models.Customer

	// Retrieve all customers ordered by created_at DESC
	if err := r.db.Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	return customers, nil
}

// Search searches customers based on criteria (removed pagination)
func (r *customerRepository) Search(req models.CustomerSearchRequest) ([]models.Customer, error) {
	var customers []models.Customer

	query := r.db.Model(&models.Customer{})

	// Apply search filters
	if req.Query != "" {
		searchTerm := "%" + req.Query + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Retrieve all matching customers
	if err := query.Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	return customers, nil
}
