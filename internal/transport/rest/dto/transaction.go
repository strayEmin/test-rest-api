package dto

import "time"

type TransactionCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TransactionCreateResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
