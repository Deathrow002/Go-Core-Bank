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
	Update(customer *models.Customer) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]models.Customer, int64, error)
	Search(req models.CustomerSearchRequest) ([]models.Customer, int64, error)
}

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository instance
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

// List retrieves customers with pagination
func (r *customerRepository) List(page, pageSize int) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	// Count total records
	if err := r.db.Model(&models.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve customers with pagination
	if err := r.db.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}

	return customers, total, nil
}

// Search searches customers based on criteria
func (r *customerRepository) Search(req models.CustomerSearchRequest) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

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

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Apply pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	// Retrieve customers
	if err := query.Limit(req.PageSize).Offset(offset).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search customers: %w", err)
	}

	return customers, total, nil
}
