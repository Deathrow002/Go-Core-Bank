package repository

import (
	"account-service/internal/account/models"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountRepository struct {
    db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
    return &accountRepository{
        db: db,
    }
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
    // Generate account number
    account.AccountNumber = r.generateAccountNumber(account.AccountType)
    
    return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
    var account models.Account
    err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &account, nil
}

func (r *accountRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID, page, pageSize int) ([]models.Account, int64, error) {
    var accounts []models.Account
    var total int64

    // Count total records
    r.db.WithContext(ctx).Model(&models.Account{}).Where("customer_id = ?", customerID).Count(&total)

    // Get paginated results
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Where("customer_id = ?", customerID).
        Order("created_at DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&accounts).Error

    return accounts, total, err
}

func (r *accountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
    var account models.Account
    err := r.db.WithContext(ctx).First(&account, "account_number = ?", accountNumber).Error
    if err != nil {
        return nil, err
    }
    return &account, nil
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) error {
    return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.Account{}, "id = ?", id).Error
}

func (r *accountRepository) List(ctx context.Context, page, pageSize int) ([]models.Account, int64, error) {
    var accounts []models.Account
    var total int64

    // Count total records
    r.db.WithContext(ctx).Model(&models.Account{}).Count(&total)

    // Get paginated results
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Order("created_at DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&accounts).Error

    return accounts, total, err
}

func (r *accountRepository) Search(ctx context.Context, query string, accountType *models.AccountType, status *models.AccountStatus, page, pageSize int) ([]models.Account, int64, error) {
    var accounts []models.Account
    var total int64

    // Build query
    db := r.db.WithContext(ctx).Model(&models.Account{})

    // Add search conditions
    if query != "" {
        searchQuery := "%" + strings.ToLower(query) + "%"
        db = db.Where("LOWER(account_number) LIKE ? OR LOWER(customer_id::text) LIKE ?", searchQuery, searchQuery)
    }

    if accountType != nil {
        db = db.Where("account_type = ?", *accountType)
    }

    if status != nil {
        db = db.Where("status = ?", *status)
    }

    // Count total records
    db.Count(&total)

    // Get paginated results
    offset := (page - 1) * pageSize
    err := db.Order("created_at DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&accounts).Error

    return accounts, total, err
}

// Additional methods for account operations
func (r *accountRepository) UpdateBalance(ctx context.Context, accountID uuid.UUID, newBalance int64) error {
    return r.db.WithContext(ctx).Model(&models.Account{}).
        Where("id = ?", accountID).
        Update("balance", newBalance).Error
}

func (r *accountRepository) GetAccountsByStatus(ctx context.Context, status models.AccountStatus, page, pageSize int) ([]models.Account, int64, error) {
    var accounts []models.Account
    var total int64

    // Count total records
    r.db.WithContext(ctx).Model(&models.Account{}).Where("status = ?", status).Count(&total)

    // Get paginated results
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Where("status = ?", status).
        Order("created_at DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&accounts).Error

    return accounts, total, err
}

func (r *accountRepository) UpdateStatus(ctx context.Context, accountID uuid.UUID, status models.AccountStatus) error {
    return r.db.WithContext(ctx).Model(&models.Account{}).
        Where("id = ?", accountID).
        Update("status", status).Error
}

func (r *accountRepository) CheckAccountExists(ctx context.Context, customerID uuid.UUID, accountType models.AccountType) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&models.Account{}).
        Where("customer_id = ? AND account_type = ? AND status != ?", customerID, accountType, models.AccountStatusClosed).
        Count(&count).Error
    
    return count > 0, err
}

func (r *accountRepository) GetAccountsWithLowBalance(ctx context.Context, threshold int64, page, pageSize int) ([]models.Account, int64, error) {
    var accounts []models.Account
    var total int64

    // Count total records
    r.db.WithContext(ctx).Model(&models.Account{}).
        Where("balance < ? AND status = ?", threshold, models.AccountStatusActive).
        Count(&total)

    // Get paginated results
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Where("balance < ? AND status = ?", threshold, models.AccountStatusActive).
        Order("balance ASC").
        Offset(offset).
        Limit(pageSize).
        Find(&accounts).Error

    return accounts, total, err
}

func (r *accountRepository) generateAccountNumber(accountType models.AccountType) string {
    var prefix string
    switch accountType {
    case models.AccountTypeSavings:
        prefix = "SAV"
    case models.AccountTypeCurrent:
        prefix = "CUR"
    case models.AccountTypeFixed:
        prefix = "FIX"
    case models.AccountTypeNonResident:
        prefix = "NR"
    default:
        prefix = "ACC"
    }

    // Generate random number (in production, use more sophisticated logic)
    id := uuid.New().ID()
    return fmt.Sprintf("%s%010d", prefix, id%1000000000) // Fixed to avoid overflow
}

// Helper method to generate unique account number with retry logic
func (r *accountRepository) generateUniqueAccountNumber(ctx context.Context, accountType models.AccountType) (string, error) {
    maxRetries := 10
    
    for i := 0; i < maxRetries; i++ {
        accountNumber := r.generateAccountNumber(accountType)
        
        // Check if account number already exists
        var count int64
        err := r.db.WithContext(ctx).Model(&models.Account{}).
            Where("account_number = ?", accountNumber).
            Count(&count).Error
        
        if err != nil {
            return "", err
        }
        
        if count == 0 {
            return accountNumber, nil
        }
    }
    
    return "", fmt.Errorf("failed to generate unique account number after %d attempts", maxRetries)
}