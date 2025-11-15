package domain

import "time"

type TransactionStatus string

const (
	TransactionStatusCreated    TransactionStatus = "created"
	TransactionStatusInProgress TransactionStatus = "in_progress"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
)

func (status TransactionStatus) String() string {
	return string(status)
}

type Transaction struct {
	ID          int64
	Name        string
	Description string
	UserID      int64
	Status      TransactionStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTransaction(name, description string, status TransactionStatus) *Transaction {
	return &Transaction{
		Name:        name,
		Description: description,
		Status:      status,
	}
}

type TransactionModel struct {
	ID          int64
	Name        string
	Description string
	UserID      int64
	Status      TransactionStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
