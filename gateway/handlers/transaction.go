package handlers

import (
	"gateway/model"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	DB *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

// Create Transaction
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	chatID, _ := strconv.ParseInt(userID, 10, 64)

	var input struct {
		Amount         float64 `json:"amount"`
		TransactionType string  `json:"transaction_type"` // income / expense
		CategoryID       uint  `json:"category_id"`
		TransactionDate string  `json:"transaction_date"` // Format: YYYY-MM-DD
		OriginalText   string  `json:"original_text"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be greater than 0",
		})
	}

	transactionType := strings.ToUpper(strings.TrimSpace(input.TransactionType))
	if transactionType != "INCOME" && transactionType != "EXPENSE" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction type must be either INCOME or EXPENSE",
		})
	}
	
	parsedDate, err := time.Parse("2006-01-02", input.TransactionDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format, use YYYY-MM-DD",
		})
	}

	transaction := model.Transaction{
		ChatID:           chatID,
		Amount:           input.Amount,
		TransactionType:  model.TransactionType(transactionType),
		CategoryID:       input.CategoryID,
		TransactionDate:  parsedDate,
		OriginalText:     input.OriginalText,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.DB.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(transaction)
}

// Get All Transactions with optional filters
func (h *TransactionHandler) GetTransactions(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var transactions []model.Transaction

	query := h.DB.Model(&model.Transaction{}).Where("chat_id = ?", userID)

	// Optional filter: transaction_type
	if txType := c.Query("type"); txType != "" {
		query = query.Where("transaction_type = ?", txType)
	}

	// Optional filter: category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Optional filter: date range
	startDate := c.Query("start_date") // format: YYYY-MM-DD
	endDate := c.Query("end_date")     // format: YYYY-MM-DD

	if startDate != "" && endDate != "" {
		query = query.Where("transaction_date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		query = query.Where("transaction_date >= ?", startDate)
	} else if endDate != "" {
		query = query.Where("transaction_date <= ?", endDate)
	}

	// Execute query
	if err := query.Order("transaction_date DESC").Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	return c.JSON(transactions)
}

// Get Detail Transaction by ID
func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	transactionID := c.Params("id")

	var transaction model.Transaction
	if err := h.DB.Where("id = ? AND chat_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}

// Update Transaction
func (h *TransactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	transactionID := c.Params("id")

	var input struct {
		Amount         float64 `json:"amount"`
		TransactionType string  `json:"transaction_type"` // income / expense
		CategoryID      uint  `json:"category_id"`
		TransactionDate string  `json:"transaction_date"`
		OriginalText   string  `json:"original_text"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var transaction model.Transaction
	if err := h.DB.Where("id = ? AND chat_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if input.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be greater than 0",
		})
	}

	transactionType := strings.ToUpper(strings.TrimSpace(input.TransactionType))
	if transactionType != "INCOME" && transactionType != "EXPENSE" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction type must be either INCOME or EXPENSE",
		})
	}
	
	parsedDate, err := time.Parse("2006-01-02", input.TransactionDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format, use YYYY-MM-DD",
		})
	}

	transaction.Amount = input.Amount
	transaction.CategoryID = input.CategoryID
	transaction.OriginalText = input.OriginalText
	if !parsedDate.IsZero() {
		transaction.TransactionDate = parsedDate
	}
	transaction.UpdatedAt = time.Now()

	if err := h.DB.Save(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update transaction",
		})
	}

	return c.JSON(transaction)
}

// Delete Transaction
func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	transactionID := c.Params("id")

	var transaction model.Transaction
	if err := h.DB.Where("id = ? AND chat_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
	}

	if err := h.DB.Delete(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
