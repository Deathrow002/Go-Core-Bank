package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountType string
type AccountStatus string
type Currency string

const (
    // Account Types
    AccountTypeSavings		AccountType = "savings"
	AccountTypeCurrent		AccountType = "current"
	AccountTypeFixed		AccountType = "fixed"
	AccountTypeNonResident	AccountType = "non_resident"

    // Account Status
    AccountStatusActive		AccountStatus = "active"
    AccountStatusInactive	AccountStatus = "inactive"
    AccountStatusSuspended	AccountStatus = "suspended"
    AccountStatusClosed		AccountStatus = "closed"

	// Currency
	DefaultCurrency = "USD"
	CurrencyTHB   Currency = "THB"
	CurrencyUSD   Currency = "USD"
	CurrencyEUR   Currency = "EUR"
	CurrencyJPY   Currency = "JPY"
	CurrencyGBP   Currency = "GBP"
	CurrencyAUD   Currency = "AUD"
	CurrencyCNY   Currency = "CNY"
	CurrencySGD   Currency = "SGD"
	CurrencyHKD   Currency = "HKD"
	CurrencyMYR   Currency = "MYR"
	CurrencyIDR   Currency = "IDR"
	CurrencyPHP   Currency = "PHP"
	CurrencyVND   Currency = "VND"
	CurrencyKRW   Currency = "KRW"
	CurrencyINR   Currency = "INR"
	CurrencyRUB   Currency = "RUB"
)

type Account struct {
    ID           uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CustomerID   uuid.UUID          `json:"customer_id" gorm:"type:uuid;not null;index"`
    AccountType  AccountType        `json:"account_type" gorm:"type:varchar(20);not null"`
    AccountNumber string            `json:"account_number" gorm:"type:varchar(20);unique;not null;index"`
    Balance      int64              `json:"balance" gorm:"type:bigint;default:0"` // Store in cents
    Currency     string             `json:"currency" gorm:"type:varchar(3);default:'USD'"`
    Status       AccountStatus      `json:"status" gorm:"type:varchar(20);default:'active'"`
    CreatedAt    time.Time          `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt    time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt    gorm.DeletedAt     `json:"-" gorm:"index"`
}

type AccountRequest struct {
    CustomerID  uuid.UUID       `json:"customer_id" binding:"required"`
    AccountType AccountType     `json:"account_type" binding:"required,oneof=savings checking credit loan"`
    Currency    string          `json:"currency" binding:"omitempty,len=3"`
    InitialBalance float64      `json:"initial_balance" binding:"omitempty,gt=0"` // In dollars
}

type AccountResponse struct {
    ID            uuid.UUID     `json:"id"`
    CustomerID    uuid.UUID     `json:"customer_id"`
    AccountType   AccountType   `json:"account_type"`
    AccountNumber string        `json:"account_number"`
    Balance       float64       `json:"balance"` // Convert from cents to dollars
    Currency      string        `json:"currency"`
    Status        AccountStatus `json:"status"`
    CreatedAt     time.Time     `json:"created_at"`
    UpdatedAt     time.Time     `json:"updated_at"`
}

type AccountListResponse struct {
    Accounts   []AccountResponse `json:"accounts"`
    Total      int64             `json:"total"`
    Page       int               `json:"page"`
    PageSize   int               `json:"page_size"`
    TotalPages int               `json:"total_pages"`
}

// Convert cents to dollars for display
func (a *Account) ToResponse() AccountResponse {
    return AccountResponse{
        ID:            a.ID,
        CustomerID:    a.CustomerID,
        AccountType:   a.AccountType,
        AccountNumber: a.AccountNumber,
        Balance:       float64(a.Balance) / 100.0, // Convert cents to dollars
        Currency:      a.Currency,
        Status:        a.Status,
        CreatedAt:     a.CreatedAt,
        UpdatedAt:     a.UpdatedAt,
    }
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
    if a.ID == uuid.Nil {
        a.ID = uuid.New()
    }
    return nil
}