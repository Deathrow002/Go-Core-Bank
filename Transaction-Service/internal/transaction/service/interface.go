package service

import (
	"context"
	"transaction-service/internal/transaction/models"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req models.CreateTransactionRequest) (*models.Transaction, error)
	CreateWithdrawalTransaction(ctx context.Context, req models.CreateTransactionRequest) (*models.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error)
	CreateDepositTransaction(ctx context.Context, req models.CreateTransactionRequest) (*models.Transaction, error)
	ListByAccNoOwner(ctx context.Context, transaction *models.Transaction) error
	ListAllTransactions(ctx context.Context) ([]*models.Transaction, error)
}