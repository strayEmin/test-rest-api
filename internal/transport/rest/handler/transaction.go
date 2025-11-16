package handler

import (
	"context"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/transport/rest/dto"
	"test-task-rest-api/internal/utils/jwt"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

type TransactionService interface {
	Create(ctx context.Context, userId int64, name, description string) (domain.Transaction, error)
	TransactionsByUserId(ctx context.Context, userId int64) ([]domain.Transaction, error)
	//UpdateStatus(ctx context.Context, userId, transactionId int64, status domain.TransactionStatus) error
}

type TransactionHandler struct {
	transactionService TransactionService
}

func NewTransactionHandler(transactionService TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token metadata extraction failed",
		})
	}

	req := dto.TransactionCreateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	if req.Name == "" || req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Required fields are missing or empty",
		})
	}

	if utf8.RuneCountInString(req.Name) > 32 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name must be between 32 characters long",
		})
	}

	if utf8.RuneCountInString(req.Description) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description must be between 255 characters long",
		})
	}

	transaction, err := h.transactionService.Create(ctx, claims.UserID, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create transaction",
		})
	}

	resp := dto.TransactionCreateResponse{
		ID:          transaction.ID,
		Name:        transaction.Name,
		Description: transaction.Description,
		Status:      transaction.Status.String(),
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

//func (h *TransactionHandler) UpdateStatus(c *fiber.Ctx) error {
//	ctx := c.Context()
//
//	claims, err := jwt.ExtractTokenMetadata(c)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Token metadata extraction failed",
//		})
//	}
//
//	req :=
//}

func (h *TransactionHandler) TransactionsByUserId(c *fiber.Ctx) error {
	ctx := c.Context()

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token metadata extraction failed",
		})
	}

	transactions, err := h.transactionService.TransactionsByUserId(ctx, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot get transactions",
		})
	}

	resp := make([]dto.TransactionGetResponse, len(transactions))
	for i, transaction := range transactions {
		resp[i] = dto.TransactionGetResponse{
			ID:          transaction.ID,
			Name:        transaction.Name,
			Description: transaction.Description,
			Status:      transaction.Status.String(),
			CreatedAt:   transaction.CreatedAt,
			UpdatedAt:   transaction.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transactions": resp,
	})
}
