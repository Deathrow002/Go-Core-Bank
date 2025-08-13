package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer  TransactionType = "transfer"
)

type Transaction struct {
	ID            	uuid.UUID		`json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AccNoOwner  	uuid.UUID		`json:"acc_no_owner" gorm:"typeuuid;index"`
	AccNoTarget		uuid.UUID		`json:"acc_no_target" gorm:"typeuuid;index"`
	Amount			float64			`json:"amount" gorm:"type:bigint;not null"` // Store in cents
	Type			TransactionType	`json:"type" gorm:"type:varchar(20);not null"`
	Description		string			`json:"description" gorm:"type:varchar(255)"`
	CreatedAt		time.Time		`json:"created_at" gorm:"autoCreateTime"`
}

type CreateTransactionRequest struct {
	ID			uuid.UUID		`json:"id"`
	AccNoOwner  uuid.UUID		`json:"acc_no_owner"`
	AccNoTarget uuid.UUID		`json:"acc_no_target"`
	Amount      float64			`json:"amount"`
	Description string			`json:"description"`
}