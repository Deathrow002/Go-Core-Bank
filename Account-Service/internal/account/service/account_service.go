package service

import (
	"account-service/internal/account/models"
	"account-service/internal/account/repository"
	"account-service/internal/external"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountService interface {
    CreateAccount(ctx context.Context, req models.AccountRequest, authToken string) (*models.Account, error)
    GetAccount(ctx context.Context, id uuid.UUID) (*models.Account, error)
    GetAccountsByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Account, int64, error)
    GetAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error)
    UpdateAccount(ctx context.Context, id uuid.UUID, req models.AccountRequest) (*models.Account, error)
    DeleteAccount(ctx context.Context, id uuid.UUID) error
    ListAccounts(ctx context.Context) ([]models.Account, int64, error)
    SearchAccounts(ctx context.Context, query string, accountType *models.AccountType, status *models.AccountStatus) ([]models.Account, int64, error)
    UpdateAccountStatus(ctx context.Context, id uuid.UUID, status models.AccountStatus) error
    UpdateBalance(ctx context.Context, AccountNumber string, newBalance int64) error
    ValidateAccount(ctx context.Context, AccountNumber string) (bool, error)
}

type accountService struct {
    repo           repository.AccountRepository
    customerClient *external.CustomerClient // Add customer client
}

// Update constructor to accept customer client
func NewAccountService(repo repository.AccountRepository, customerClient *external.CustomerClient) AccountService {
    return &accountService{
        repo:           repo,
        customerClient: customerClient,
    }
}

func (s *accountService) CreateAccount(ctx context.Context, req models.AccountRequest, authToken string) (*models.Account, error) {    
    // Validate request
    if req.CustomerID == uuid.Nil {
        return nil, errors.New("customer ID is required")
    }
    if req.AccountType == "" {
        return nil, errors.New("account type is required")
    }

    // **Call Customer Service to validate customer exists**
    customerExists, err := s.customerClient.ValidateCustomer(ctx, req.CustomerID, authToken)
    if err != nil {
        return nil, fmt.Errorf("failed to validate customer: %w", err)
    }

    if !customerExists {
        return nil, errors.New("customer does not exist")
    }

    // Set default currency if not provided
    if req.Currency == "" {
        req.Currency = "USD"
    }

    // Create account model
    account := &models.Account{
        CustomerID:  req.CustomerID,
        AccountType: req.AccountType,
        Currency:    req.Currency,
        Balance:     0, // Start with zero balance
        Status:      models.AccountStatusActive,
    }

    // Save to database (account number will be generated in repository)
    if err := s.repo.Create(ctx, account); err != nil {
        return nil, fmt.Errorf("failed to create account: %w", err)
    }

    return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id uuid.UUID) (*models.Account, error) {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("account not found")
        }
        return nil, fmt.Errorf("failed to get account: %w", err)
    }
    return account, nil
}

func (s *accountService) GetAccountsByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Account, int64, error) {
    accounts, total, err := s.repo.GetByCustomerID(ctx, customerID, 0, 0)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get accounts by customer: %w", err)
    }
    return accounts, total, nil
}

func (s *accountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
    if accountNumber == "" {
        return nil, errors.New("account number is required")
    }

    account, err := s.repo.GetByAccountNumber(ctx, accountNumber)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("account not found")
        }
        return nil, fmt.Errorf("failed to get account by number: %w", err)
    }
    return account, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, id uuid.UUID, req models.AccountRequest) (*models.Account, error) {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("account not found")
        }
        return nil, fmt.Errorf("failed to get account: %w", err)
    }

    if req.Currency != "" {
        account.Currency = req.Currency
    }

    if err := s.repo.Update(ctx, account); err != nil {
        return nil, fmt.Errorf("failed to update account: %w", err)
    }

    return account, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("account not found")
        }
        return fmt.Errorf("failed to get account: %w", err)
    }

    if account.Balance != 0 {
        return errors.New("cannot delete account with non-zero balance")
    }

    account.Status = models.AccountStatusClosed
    if err := s.repo.Update(ctx, account); err != nil {
        return fmt.Errorf("failed to close account: %w", err)
    }

    return nil
}

func (s *accountService) ListAccounts(ctx context.Context) ([]models.Account, int64, error) {
    accounts, total, err := s.repo.List(ctx)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list accounts: %w", err)
    }

    return accounts, total, nil
}

func (s *accountService) SearchAccounts(ctx context.Context, query string, accountType *models.AccountType, status *models.AccountStatus) ([]models.Account, int64, error) {
    accounts, total, err := s.repo.Search(ctx, query, accountType, status, 0, 0)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to search accounts: %w", err)
    }

    return accounts, total, nil
}

// Additional utility methods for business logic
func (s *accountService) UpdateAccountStatus(ctx context.Context, id uuid.UUID, status models.AccountStatus) error {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("account not found")
        }
        return fmt.Errorf("failed to get account: %w", err)
    }

    account.Status = status
    if err := s.repo.Update(ctx, account); err != nil {
        return fmt.Errorf("failed to update account status: %w", err)
    }

    return nil
}

func (s *accountService) UpdateBalance(ctx context.Context, AccountNumber string, newBalance int64) error {
    return s.repo.UpdateBalance(ctx, AccountNumber, newBalance)
}

func (s *accountService) ValidateAccount(ctx context.Context, AccountNumber string) (bool, error) {
    if AccountNumber == "" {
        return false, errors.New("account number is required")
    }

    exists, err := s.repo.CheckAccountExists(ctx, AccountNumber)
    if err != nil {
        return exists, fmt.Errorf("failed to check account existence: %w", err)
    }

    if !exists {
        return exists, errors.New("account does not exist")
    }

    return exists, nil
}