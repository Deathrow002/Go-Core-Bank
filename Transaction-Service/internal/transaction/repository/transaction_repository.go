package repository

import (
	"context"
	"errors"
	"transaction-service/internal/transaction/models"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	ListByAccNoOwner(ctx context.Context, transaction *models.Transaction) error
	ListAllTransactions(ctx context.Context) ([]models.Transaction, error)
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) ListByAccNoOwner(ctx context.Context, transaction *models.Transaction) error {
	if err := r.db.WithContext(ctx).Where("acc_no_owner = ?", transaction.AccNoOwner).Find(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) ListAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.WithContext(ctx).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}