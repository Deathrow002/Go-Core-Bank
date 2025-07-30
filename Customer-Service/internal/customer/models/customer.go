package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomDate handles date-only values for PostgreSQL
type CustomDate struct {
	time.Time
}

// UnmarshalJSON handles multiple date formats and converts to date-only
func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	dateStr := strings.Trim(string(data), `"`)
	
	if dateStr == "null" || dateStr == "" {
		return nil
	}

	// List of supported date formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339 with timezone
		"2006-01-02T15:04:05Z",      // RFC3339 UTC
		"2006-01-02T15:04:05",       // RFC3339 without timezone
		"2006-01-02",                // Date only
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			// Convert to date-only (midnight UTC)
			cd.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			return nil
		}
	}

	return fmt.Errorf("unable to parse date: %s", dateStr)
}

// MarshalJSON returns the date in ISO format
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(cd.Time.Format("2006-01-02T15:04:05Z"))
}

// Value implements the driver.Valuer interface for database storage
func (cd CustomDate) Value() (driver.Value, error) {
	if cd.Time.IsZero() {
		return nil, nil
	}
	// Return only the date part as string for PostgreSQL date type
	return cd.Time.Format("2006-01-02"), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		// Convert to date-only (midnight UTC)
		cd.Time = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("failed to parse date string '%s': %w", v, err)
		}
		cd.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into CustomDate", value)
	}
}

// Customer represents a bank customer
type Customer struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName   string         `json:"first_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	LastName    string         `json:"last_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null;size:255" validate:"required,email"`
	Phone       string         `json:"phone" gorm:"size:20" validate:"required"`
	DateOfBirth *time.Time     `json:"date_of_birth" gorm:"type:date"`  // Use *time.Time with DATE type
	Address     Address        `json:"address" gorm:"embedded;embeddedPrefix:address_"`
	Status      CustomerStatus `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Address represents customer address information
type Address struct {
	Street     string `json:"street" gorm:"size:255"`
	City       string `json:"city" gorm:"size:100"`
	State      string `json:"state" gorm:"size:100"`
	PostalCode string `json:"postal_code" gorm:"size:20"`
	Country    string `json:"country" gorm:"size:100"`
}

// CustomerStatus represents the status of a customer
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusSuspended CustomerStatus = "suspended"
	CustomerStatusClosed    CustomerStatus = "closed"
)

// CustomerRequest represents the request payload for creating/updating a customer
type CustomerRequest struct {
	FirstName   string     `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string     `json:"last_name" validate:"required,min=2,max=100"`
	Email       string     `json:"email" validate:"required,email"`
	Phone       string     `json:"phone" validate:"required"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"required"`  // Use *time.Time
	Address     Address    `json:"address" validate:"required"`
}

// CustomerResponse represents the response payload for customer operations
type CustomerResponse struct {
	ID          uuid.UUID      `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	DateOfBirth *time.Time     `json:"date_of_birth"`
	Address     Address        `json:"address"`
	Status      CustomerStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// CustomerValidationRequest represents the validation request
type CustomerValidationRequest struct {
	// For JSON input (string UUID)
	CustomerID string `json:"customer_id,omitempty" validate:"omitempty,uuid"`

	// For service layer (parsed UUID)
	ID uuid.UUID `json:"-" validate:"required"`
}

// CustomerSearchRequest represents search criteria
type CustomerSearchRequest struct {
	Query  string `form:"query" json:"query"`
	Status string `form:"status" json:"status"`
}

// ToResponse converts Customer model to CustomerResponse
func (c *Customer) ToResponse() CustomerResponse {
	var dob *time.Time
	if c.DateOfBirth != nil && !c.DateOfBirth.IsZero() {
		dob = c.DateOfBirth
	}
	return CustomerResponse{
		ID:          c.ID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Email:       c.Email,
		Phone:       c.Phone,
		DateOfBirth: dob,
		Address:     c.Address,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// TableName returns the table name for Customer model
func (Customer) TableName() string {
	return "customers"
}

// Helper methods for CustomDate
func NewCustomDate(year int, month time.Month, day int) CustomDate {
	return CustomDate{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func ParseCustomDate(dateStr string) (CustomDate, error) {
	var cd CustomDate
	err := cd.UnmarshalJSON([]byte(`"` + dateStr + `"`))
	return cd, err
}
