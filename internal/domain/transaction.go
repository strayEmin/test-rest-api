package domain

import "time"

type TransactionStatus string

const (
	TransactionStatusCreated   TransactionStatus = "created"
	TransactionStatusInProcess TransactionStatus = "in_process"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

func (status TransactionStatus) String() string {
	return string(status)
}

var validStatuses = map[string]TransactionStatus{
	TransactionStatusCreated.String():   TransactionStatusCreated,
	TransactionStatusInProcess.String(): TransactionStatusInProcess,
	TransactionStatusCompleted.String(): TransactionStatusCompleted,
	TransactionStatusCancelled.String(): TransactionStatusCancelled,
}

func IsTransactionStatus(s string) bool {
	_, ok := validStatuses[s]
	return ok
}

func ToTransactionStatus(s string) TransactionStatus {
	if st, ok := validStatuses[s]; ok {
		return st
	}
	return ""
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
