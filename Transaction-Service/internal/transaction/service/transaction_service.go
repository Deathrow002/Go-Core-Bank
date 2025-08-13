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

func (s *transactionService) CreateTransaction(ctx context.Context, req models.CreateTransactionRequest, authToken string) (*models.Transaction, error) {
    log.Printf("Creating transaction: Owner=%s, Target=%s, Amount=%.2f", req.AccNoOwner, req.AccNoTarget, req.Amount)
    
    // Validate request
    if req.AccNoOwner == uuid.Nil {
        log.Printf("Validation failed: account number owner is required")
        return nil, fmt.Errorf("account number owner is required")
    }
    if req.AccNoTarget == uuid.Nil {
        log.Printf("Validation failed: account number target is required")
        return nil, fmt.Errorf("account number target is required")
    }
    if req.Amount <= 0 {
        log.Printf("Validation failed: amount must be greater than zero (received: %.2f)", req.Amount)
        return nil, fmt.Errorf("amount must be greater than zero")
    }

    log.Printf("Fetching owner account balance: AccountID=%s", req.AccNoOwner)
    OwnerBalance, err := s.accountClient.GetAccountBalance(ctx, req.AccNoOwner.String(), authToken)
    if err != nil {
        log.Printf("Error getting owner account balance: %v", err)
        return nil, fmt.Errorf("failed to get owner account balance: %w", err)
    }
    log.Printf("Owner account balance retrieved: %.2f", OwnerBalance)
    
    if OwnerBalance < req.Amount {
        log.Printf("Insufficient balance: Available=%.2f, Required=%.2f", OwnerBalance, req.Amount)
        return nil, fmt.Errorf("insufficient balance in owner account")
    }

    OwnerBalance = OwnerBalance - req.Amount
    log.Printf("Owner balance after deduction: %.2f", OwnerBalance)

    log.Printf("Fetching target account balance: AccountID=%s", req.AccNoTarget)
    TargetBalance, err := s.accountClient.GetAccountBalance(ctx, req.AccNoTarget.String(), authToken)
    if err != nil {
        log.Printf("Error getting target account balance: %v", err)
        return nil, fmt.Errorf("failed to get target account balance: %w", err)
    }
    log.Printf("Target account balance retrieved: %.2f", TargetBalance)

    TargetBalance = TargetBalance + req.Amount
    log.Printf("Target balance after addition: %.2f", TargetBalance)
    
    transaction := &models.Transaction{
        ID:            uuid.New(),
        AccNoOwner:    req.AccNoOwner,
        AccNoTarget:   req.AccNoTarget,
        Amount:        req.Amount,
        Type:          models.TransactionTypeTransfer,
        Description:   req.Description,
        CreatedAt:     time.Now(),
    }
    log.Printf("Created transaction object: ID=%s, Type=%s", transaction.ID, transaction.Type)

    log.Printf("Sending balance update to Kafka for owner account: ID=%s, Balance=%.2f, Type=%s", 
        transaction.AccNoOwner, OwnerBalance, string(transaction.Type))
    if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoOwner, OwnerBalance, string(transaction.Type)); err != nil {
        log.Printf("Failed to send owner balance update to Kafka: %v", err)
    }

    log.Printf("Sending balance update to Kafka for target account: ID=%s, Balance=%.2f, Type=%s", 
        transaction.AccNoTarget, TargetBalance, string(transaction.Type))
    if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoTarget, TargetBalance, string(transaction.Type)); err != nil {
        log.Printf("Failed to send target balance update to Kafka: %v", err)
        return nil, fmt.Errorf("failed to send balance update: %w", err)
    }

    log.Printf("Saving transaction to database: ID=%s", transaction.ID)
    if err := s.repo.Create(ctx, transaction); err != nil {
        log.Printf("Failed to create transaction in database: %v", err)
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }
    log.Printf("Transaction created successfully: ID=%s", transaction.ID)

    return transaction, nil
}

func (s *transactionService) CreateWithdrawalTransaction(ctx context.Context, req models.CreateTransactionRequest, authToken string) (*models.Transaction, error) {
    log.Printf("Creating withdrawal transaction: Owner=%s, Amount=%.2f", req.AccNoOwner, req.Amount)
    
    transaction := &models.Transaction{
        ID:            uuid.New(),
        AccNoOwner:    req.AccNoOwner,
        AccNoTarget:   uuid.Nil, // No target account for withdrawal
        Amount:        req.Amount,
        Type:          models.TransactionTypeWithdrawal,
        Description:   req.Description,
        CreatedAt:     time.Now(),
    }
    log.Printf("Created withdrawal transaction object: ID=%s", transaction.ID)

    log.Printf("Saving withdrawal transaction to database: ID=%s", transaction.ID)
    if err := s.repo.Create(ctx, transaction); err != nil {
        log.Printf("Failed to create withdrawal transaction in database: %v", err)
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }
    log.Printf("Withdrawal transaction saved successfully: ID=%s", transaction.ID)

    log.Printf("Sending balance update to Kafka for withdrawal: Account=%s, Amount=%.2f, Type=%s", 
        transaction.AccNoOwner, transaction.Amount, string(transaction.Type))
    if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoOwner, transaction.Amount, string(transaction.Type)); err != nil {
        log.Printf("Failed to send withdrawal balance update to Kafka: %v", err)
        return nil, fmt.Errorf("failed to send balance update: %w", err)
    }
    log.Printf("Withdrawal balance update sent successfully")

    return transaction, nil
}

func (s *transactionService) CreateDepositTransaction(ctx context.Context, req models.CreateTransactionRequest) (*models.Transaction, error) {
    log.Printf("Creating deposit transaction: Owner=%s, Amount=%.2f", req.AccNoOwner, req.Amount)
    
    transaction := &models.Transaction{
        ID:            uuid.New(),
        AccNoOwner:    req.AccNoOwner,
        AccNoTarget:   uuid.Nil, // No target account for deposit
        Amount:        req.Amount,
        Type:          models.TransactionTypeDeposit,
        Description:   req.Description,
        CreatedAt:     time.Now(),
    }
    log.Printf("Created deposit transaction object: ID=%s", transaction.ID)

    log.Printf("Saving deposit transaction to database: ID=%s", transaction.ID)
    if err := s.repo.Create(ctx, transaction); err != nil {
        log.Printf("Failed to create deposit transaction in database: %v", err)
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }
    log.Printf("Deposit transaction saved successfully: ID=%s", transaction.ID)

    log.Printf("Sending balance update to Kafka for deposit: Account=%s, Amount=%.2f, Type=%s", 
        transaction.AccNoOwner, transaction.Amount, string(transaction.Type))
    if err := producer.SendBalanceUpdate(s.kafkaWriter, transaction.AccNoOwner, transaction.Amount, string(transaction.Type)); err != nil {
        log.Printf("Failed to send deposit balance update to Kafka: %v", err)
        return nil, fmt.Errorf("failed to send balance update: %w", err)
    }
    log.Printf("Deposit balance update sent successfully")

    return transaction, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error) {
    log.Printf("Getting transaction by ID: %s", id)
    
    transaction, err := s.repo.GetByID(ctx, id)
    if err != nil {
        log.Printf("Failed to get transaction by ID %s: %v", id, err)
        return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
    }
    
    log.Printf("Transaction retrieved successfully: ID=%s, Owner=%s, Amount=%.2f, Type=%s", 
        transaction.ID, transaction.AccNoOwner, transaction.Amount, transaction.Type)
    return transaction, nil
}

func (s *transactionService) ListAllTransactions(ctx context.Context) ([]*models.Transaction, error) {
    log.Printf("Listing all transactions")
    
    transactions, err := s.repo.ListAllTransactions(ctx)
    if err != nil {
        log.Printf("Failed to list all transactions: %v", err)
        return nil, fmt.Errorf("failed to list all transactions: %w", err)
    }
    
    var result []*models.Transaction
    for _, transaction := range transactions {
        result = append(result, &transaction)
    }
    
    log.Printf("Retrieved %d transactions", len(result))
    return result, nil
}

// ListByAccNoOwner returns all transactions for a given owner account number.
func (s *transactionService) ListByAccNoOwner(ctx context.Context, accNoOwner uuid.UUID) ([]*models.Transaction, error) {
    log.Printf("Listing transactions by owner account: %s", accNoOwner)
    transactions, err := s.repo.ListByAccNoOwner(ctx, accNoOwner)
    if err != nil {
        log.Printf("Failed to list transactions by owner account: %v", err)
        return nil, fmt.Errorf("failed to list transactions by owner account: %w", err)
    }
    var result []*models.Transaction
    for _, transaction := range transactions {
        result = append(result, &transaction)
    }
    log.Printf("Retrieved %d transactions for owner %s", len(result), accNoOwner)
    return result, nil
}