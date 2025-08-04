package repository

import (
	"account-service/internal/account/models"
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, page, pageSize int) ([]models.Account, int64, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]models.Account, int64, error)
	Search(ctx context.Context, query string, accountType *models.AccountType, status *models.AccountStatus, page, pageSize int) ([]models.Account, int64, error)

	// Additional methods
	UpdateBalance(ctx context.Context, AccountNumber string, newBalance int64) error
	GetAccountsByStatus(ctx context.Context, status models.AccountStatus, page, pageSize int) ([]models.Account, int64, error)
	UpdateStatus(ctx context.Context, accountID uuid.UUID, status models.AccountStatus) error
	CheckAccountExists(ctx context.Context, AccountNumber string) (bool, error)
	GetAccountsWithLowBalance(ctx context.Context, threshold int64, page, pageSize int) ([]models.Account, int64, error)
}
