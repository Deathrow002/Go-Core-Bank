package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Customer represents a bank customer
type Customer struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName   string         `json:"first_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	LastName    string         `json:"last_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null;size:255" validate:"required,email"`
	Phone       string         `json:"phone" gorm:"size:20" validate:"required,min=10,max=20"`
	DateOfBirth *time.Time     `json:"date_of_birth" gorm:"type:date"`
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
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
	CustomerStatusSuspended CustomerStatus = "suspended"
	CustomerStatusClosed   CustomerStatus = "closed"
)

// CustomerRequest represents the request payload for creating/updating a customer
type CustomerRequest struct {
	FirstName   string     `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string     `json:"last_name" validate:"required,min=2,max=100"`
	Email       string     `json:"email" validate:"required,email"`
	Phone       string     `json:"phone" validate:"required,min=10,max=20"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Address     Address    `json:"address"`
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

// CustomerListResponse represents the response for listing customers
type CustomerListResponse struct {
	Customers  []CustomerResponse `json:"customers"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// CustomerSearchRequest represents search parameters
type CustomerSearchRequest struct {
	Query    string         `json:"query" form:"query"`
	Status   CustomerStatus `json:"status" form:"status"`
	Page     int            `json:"page" form:"page"`
	PageSize int            `json:"page_size" form:"page_size"`
}

// ToResponse converts Customer model to CustomerResponse
func (c *Customer) ToResponse() CustomerResponse {
	return CustomerResponse{
		ID:          c.ID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Email:       c.Email,
		Phone:       c.Phone,
		DateOfBirth: c.DateOfBirth,
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
