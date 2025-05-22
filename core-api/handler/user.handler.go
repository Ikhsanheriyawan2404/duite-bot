package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"finance-bot/config"
	"finance-bot/service"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model string `json:"model"`
	Messages []ChatMessage `json:"messages"`
	MaxTokens int `json:"max_tokens,omitempty"`
	ResponseFormat map[string]any `json:"response_format,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Result struct {
	TransactionType string `json:"type"`
	Amount float64 `json:"amount"`
	Category string `json:"category"`
	Date string `json:"date"`
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("userId")

	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	
	user, err := h.userService.GetByChatId(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetTransactions(c *fiber.Ctx) error {
	uuid := c.Params("userId")

	transactions, err := h.userService.GetTransactions(uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data transaksi",
		})
	}

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	income := make([]float64, 12)
	expense := make([]float64, 12)
	pieMap := make(map[string]float64)

	for _, t := range *transactions {
		month := int(t.TransactionDate.Month()) - 1
		if month < 0 || month > 11 {
			continue
		}

		if t.TransactionType == "INCOME" {
			income[month] += t.Amount
		} else if t.TransactionType == "EXPENSE" {
			expense[month] += t.Amount
			pieMap[t.Category] += t.Amount
		}
	}

	type pieEntry struct {
		Category string
		Amount   float64
	}
	
	var piceSlice []pieEntry
	for k, v := range pieMap {
		piceSlice = append(piceSlice, pieEntry{Category: k, Amount: v})
	}

	sort.Slice(piceSlice, func(i, j int) bool {
		return piceSlice[i].Amount > piceSlice[j].Amount
	})

	finalPieMap := make(map[string]float64)
	for i, entry := range piceSlice {
		if i < 5 {
			finalPieMap[entry.Category] = entry.Amount
		} else {
			finalPieMap["lainnya"] += entry.Amount
		}
	}

	var pieLabels []string
	var pieData []float64
	for k, v := range finalPieMap {
		pieLabels = append(pieLabels, k)
		pieData = append(pieData, v)
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"line": fiber.Map{
			"labels": months,
			"datasets": []fiber.Map{
				{ "label": "Pemasukan", "data": income },
				{ "label": "Pengeluaran", "data": expense },
			},
		},
		"pie": fiber.Map{
			"labels": pieLabels,
			"data":   pieData,
		},
	})
}

func (h *UserHandler) GetDailyReport(c *fiber.Ctx) error {
	idParam := c.Params("chatId")
	chatId, _ := strconv.ParseInt(idParam, 10, 64)

	transactions, err := h.userService.GetDailyReport(chatId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get daily report",
		})
	}

	return c.JSON(transactions)
}

func (h *UserHandler) GetMonthlyReport(c *fiber.Ctx) error {
	idParam := c.Params("chatId")
	chatId, _ := strconv.ParseInt(idParam, 10, 64)

	transactions, err := h.userService.GetMonthlyReport(chatId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get monthly report",
		})
	}

	return c.JSON(transactions)
}

func (h *UserHandler) DeleteTransactionByID(c *fiber.Ctx) error {
	idParam := c.Params("chatId")
	chatId, _ := strconv.ParseInt(idParam, 10, 64)

	transactionIdParam := c.Params("transactionId")
	transactionId, _ := strconv.Atoi(transactionIdParam)

	err := h.userService.DeleteTransactionByID(uint(transactionId), chatId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction deleted successfully",
	})
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	type RegisterRequest struct {
		ChatID int64  `json:"chat_id"`
		Name   string `json:"name"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ChatID == 0 || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chat ID and name are required",
		})
	}

	user, err := h.userService.RegisterUser(req.ChatID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *UserHandler) CheckUser(c *fiber.Ctx) error {
	idParam := c.Params("chatId")
	chatId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	exists := h.userService.CheckUser(chatId)
	return c.JSON(fiber.Map{
		"exists": exists,
	})
}

func (h *UserHandler) AIClassifyTransaction(c *fiber.Ctx) error {
	chatIdParam := c.Params("chatId")
	chatId, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	var input struct {
		Prompt string `json:"prompt"`
	}
	if err := c.BodyParser(&input); err != nil || input.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Prompt is required",
		})
	}

	// Gunakan fungsi classifyTransaction()
	result, llmResp, err := hitChatGpt(input.Prompt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengklasifikasi transaksi",
		})
	}

	// Validasi hasil LLM (opsional)
	if result.TransactionType == "" || result.Category == "" || result.Amount == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Hasil klasifikasi tidak valid",
			"llm_raw": llmResp.Choices[0].Message.Content,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"chat_id": chatId,
		"result":  result,
		"usage":   llmResp.Usage,
	})
}

func classifyTransaction(desc string) (*Result, *ChatResponse, error) {
	prompt := desc

	requestBody := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		ResponseFormat: map[string]any{ 
			"type": "json_object",
		},
		// MaxTokens: 50,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", config.AppConfig.LLMApiUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + config.AppConfig.LLMApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chatResp ChatResponse
	json.Unmarshal(body, &chatResp)
	jsonStr := chatResp.Choices[0].Message.Content

	var result Result
	json.Unmarshal([]byte(jsonStr), &result)

	return &result, &chatResp, err
}

func hitChatGpt(desc string) (*Result, *ChatResponse, error) {

	prompt := desc

	client := openai.NewClient(
		option.WithAPIKey(config.AppConfig.LLMApiKey), // Gantilah dengan API key Anda
	)

	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: "gpt-4.1-nano",
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{
					Type: "json_object",
				},
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}
	jsonBytes, err := json.MarshalIndent(chatCompletion, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	var chatResp ChatResponse
	json.Unmarshal(jsonBytes, &chatResp)
	jsonStr := chatResp.Choices[0].Message.Content

	var result Result
	json.Unmarshal([]byte(jsonStr), &result)

	return &result, &chatResp, nil
}

