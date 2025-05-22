package handler

import (
	"finance-bot/model"
	"finance-bot/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req struct {
		ChatID          int64   `json:"chat_id"`
		OriginalText    string  `json:"original_text"`
		TransactionType string  `json:"transaction_type"` // Should be INCOME or EXPENSE
		Amount          float64 `json:"amount"`
		Category        string  `json:"category"`
		Description     string  `json:"description"`
		TransactionDate string  `json:"transaction_date"` // ISO 8601 e.g. "2025-07-30T10:00:00Z"
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Validate TransactionType
	if req.TransactionType != string(model.INCOME) && req.TransactionType != string(model.EXPENSE) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction type. Must be INCOME or EXPENSE",
		})
	}

	// Parse TransactionDate
	transactionDate, err := time.Parse(time.RFC3339, req.TransactionDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction_date format. Must be ISO8601 (e.g., 2025-07-30T10:00:00Z)",
		})
	}

	tx := &model.Transaction{
		ChatID:          req.ChatID,
		OriginalText:    req.OriginalText,
		TransactionType: model.TransactionType(req.TransactionType),
		Amount:          req.Amount,
		Category:        req.Category,
		Description:     req.Description,
		TransactionDate: transactionDate,
	}

	if err := h.transactionService.CreateTransaction(tx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tx)
}
