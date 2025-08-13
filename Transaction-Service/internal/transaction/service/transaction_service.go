package service

import (
	"context"
	"fmt"
	"log"
	"time"
	external "transaction-service/internal/external"
	producer "transaction-service/internal/external/producer"
	"transaction-service/internal/transaction/models"
	"transaction-service/internal/transaction/repository"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type transactionService struct {
	repo			repository.TransactionRepository
	accountClient	*external.AccountClient
	kafkaWriter		*kafka.Writer
}

func NewTransactionService(repo repository.TransactionRepository, accountClient *external.AccountClient, kafkaWriter *kafka.Writer) TransactionService {
	return &transactionService{
		repo:			repo,
		accountClient:	accountClient,
		kafkaWriter:	kafkaWriter,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req models.CreateTransactionRequest) (*models.Transaction, error) {
	// Validate request
	if req.AccNoOwner == uuid.Nil {
		return nil, fmt.Errorf("account number owner is required")
	}
	if req.AccNoTarget == uuid.Nil {
		return nil, fmt.Errorf("account number target is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	OwnerBalance, err := s.accountClient.GetAccountBalance(ctx, req.AccNoOwner.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get owner account balance: %w", err)
	}
	if OwnerBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance in owner account")
	}

	OwnerBalance = OwnerBalance - req.Amount

	TargetBalance, err := s.accountClient.GetAccountBalance(ctx, req.AccNoTarget.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get target account balance: %w", err)
	}

	TargetBalance = TargetBalance + req.Amount
	
	transaction := &models.Transaction{
		ID:				uuid.New(),
		AccNoOwner:		req.AccNoOwner,
		AccNoTarget:	req.AccNoTarget,
		Amount:			req.Amount,
		Type:			models.TransactionTypeTransfer,
		Description:	req.Description,
		CreatedAt:		time.Now(),
	}

	log.Printf("Sending Kafka message with type='%s'", string(transaction.Type))
	if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoOwner, OwnerBalance, string(transaction.Type)); err != nil {
    log.Printf("Failed to send balance update to Kafka: %v", err)
	}

	log.Printf("Sending Kafka message with type='%s'", string(transaction.Type))
	if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoTarget, TargetBalance, string(transaction.Type)); err != nil {
		// Log the error
		log.Printf("Failed to send balance update to Kafka: %v", err)

		// Handle error
		return nil, fmt.Errorf("failed to send balance update: %w", err)
	}

	// Save to database
	if err := s.repo.Create(ctx, transaction); err != nil {
		// Log the error
		log.Printf("Failed to create transaction: %v", err)

		// Handle error
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (s *transactionService) CreateWithdrawalTransaction(ctx context.Context, req models.CreateTransactionRequest) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:				uuid.New(),
		AccNoOwner:		req.AccNoOwner,
		AccNoTarget:	uuid.Nil, // No target account for withdrawal
		Amount:			req.Amount,
		Type:			models.TransactionTypeWithdrawal,
		Description:	req.Description,
		CreatedAt:		time.Now(),
	}

	// Save to database
	if err := s.repo.Create(ctx, transaction); err != nil {
		// Log the error
		log.Printf("Failed to create transaction: %v", err)

		// Handle error
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	log.Printf("Sending Kafka message with type='%s'", string(transaction.Type))
	// After saving transaction
	if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoOwner, transaction.Amount, string(transaction.Type)); err != nil {
		// Log the error
		log.Printf("Failed to send balance update to Kafka: %v", err)
	
		// Handle error
		return nil, fmt.Errorf("failed to send balance update: %w", err)
	}

	return transaction, nil
}

func (s *transactionService) CreateDepositTransaction(ctx context.Context,req models.CreateTransactionRequest) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:				uuid.New(),
		AccNoOwner:		req.AccNoOwner,
		AccNoTarget:	uuid.Nil, // No target account for deposit
		Amount:			req.Amount,
		Type:			models.TransactionTypeDeposit,
		Description:	req.Description,
		CreatedAt:		time.Now(),
	}

	// Save to database
	if err := s.repo.Create(ctx, transaction); err != nil {
		// Log the error
		log.Printf("Failed to create transaction: %v", err)
		
		// Handle error
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	log.Printf("Sending Kafka message with type='%s'", string(transaction.Type))
	if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoOwner, transaction.Amount, string(transaction.Type)); err != nil {
		// Log the error
		log.Printf("Failed to send balance update to Kafka: %v", err)

		// Handle error
		return nil, fmt.Errorf("failed to send balance update: %w", err)
	}

	return transaction, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Log the error
		log.Printf("Failed to get transaction by ID: %v", err)

		// Handle error
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}
	return transaction, nil
}

func (s *transactionService) ListByAccNoOwner(ctx context.Context, transaction *models.Transaction) error {
	if err := s.repo.ListByAccNoOwner(ctx, transaction); err != nil {
		return fmt.Errorf("failed to list transactions by account number owner: %w", err)
	}
	return nil
}

func (s *transactionService) ListAllTransactions(ctx context.Context) ([]*models.Transaction, error) {
	var transactions []models.Transaction
	transactions, err := s.repo.ListAllTransactions(ctx)
	if err != nil {		return nil, fmt.Errorf("failed to list all transactions: %w", err)
	}
	var result []*models.Transaction
	for _, transaction := range transactions {
		result = append(result, &transaction)
	}
	return result, nil
}