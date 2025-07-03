package handler

import (
	"finance-bot/model"
	"finance-bot/service"
	"finance-bot/utils"

	"log"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)


type Result struct {
	TransactionType string  `json:"type"`
	Amount          float64 `json:"amount"`
	CategoryID      *uint   `json:"category_id"`
	Date            string  `json:"date"`
}

type UserHandler struct {
	userService        service.UserService
	transactionService service.TransactionService
	categoryService    service.CategoryService
	aiService   	   service.LLMService
}

func NewUserHandler(
		userService service.UserService,
		transactionService service.TransactionService,
		categoryService service.CategoryService,
		aiService service.LLMService,
	) *UserHandler {
	
	return &UserHandler{
		userService:        userService,
		transactionService: transactionService,
		categoryService: 	categoryService,
		aiService:  	 	aiService,
	}
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
	pieMap := make(map[uint]float64)

	for _, t := range *transactions {
		month := int(t.TransactionDate.Month()) - 1
		if month < 0 || month > 11 {
			continue
		}

		if t.TransactionType == "INCOME" {
			income[month] += t.Amount
		} else if t.TransactionType == "EXPENSE" {
			expense[month] += t.Amount
			pieMap[*t.CategoryID] += t.Amount
		}
	}

	type pieEntry struct {
		Category uint
		Amount   float64
	}

	var piceSlice []pieEntry
	for k, v := range pieMap {
		piceSlice = append(piceSlice, pieEntry{Category: k, Amount: v})
	}

	sort.Slice(piceSlice, func(i, j int) bool {
		return piceSlice[i].Amount > piceSlice[j].Amount
	})

	finalPieMap := make(map[uint]float64)
	for i, entry := range piceSlice {
		if i < 5 {
			finalPieMap[entry.Category] = entry.Amount
		} else {
			finalPieMap[0] += entry.Amount
		}
	}

	var pieLabels []uint
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
				{"label": "Pemasukan", "data": income},
				{"label": "Pengeluaran", "data": expense},
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Transaksi tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaksi berhasil dihapus",
	})
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var request struct {
		ChatID int64  `json:"chat_id"`
		Name   string `json:"name"`
	}

	c.BodyParser(&request)

	user, err := h.userService.RegisterUser(request.ChatID, request.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) CheckUser(c *fiber.Ctx) error {
	idParam := c.Params("chatId")
	chatId, _ := strconv.ParseInt(idParam, 10, 64)

	exist := h.userService.CheckUser(chatId)
	if !exist {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"exist": exist,
		})
	}
	return c.JSON(fiber.Map{
		"exist": exist,
	})
}

func (h *UserHandler) GenerateMagicLink(c *fiber.Ctx) error {
	type request struct {
		ChatID int64 `json:"chat_id"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	token, _ := h.userService.GenerateMagicLoginToken(req.ChatID)

	return c.JSON(fiber.Map{
		"message": "Magic login link generated",
		"token":   token,
	})
}

func (h *UserHandler) ParseAndSaveTransaction(c *fiber.Ctx) error {
	chatIdParam := c.Params("chatId")
	chatId, _ := strconv.ParseInt(chatIdParam, 10, 64)

	var input struct {
		Prompt string `json:"prompt"`
	}

	if err := c.BodyParser(&input); err != nil || input.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Prompt is required",
		})
	}

	if !utils.ContainsNominal(input.Prompt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": `📌 Coba ketik gini:
➡️ Makan siang 20k  

Gampang banget kan? ✨
Kalau butuh bantuan lainnya, tinggal ketik bantuan yaa~
Biar aku bantuin kamu jadi lebih rapih ngatur duit!`,

			"help": "Contoh: 'Makan siang 25k' atau 'gaji masuk 1000'",
		})
	}

	// Step 0: Get Category Data
	typeRes, _, err := h.aiService.ClassifyTransactionType(input.Prompt)
    if err != nil {
		log.Printf("Klasifikasi Type: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "😓 Aduh, sistem lagi ngambek. Coba lagi lagi ya~"})
    }
	
	// Validasi hasil klasifikasi
	if typeRes.TransactionType != string(model.TransactionTypeINCOME) &&
	typeRes.TransactionType != string(model.TransactionTypeEXPENSE) {
		// Log actual value dari TransactionType supaya kita tahu apa yang dikembalikan LLM
		log.Printf("DEBUG: invalid transaction type from LLM: %q", typeRes.TransactionType)

		// Kalau kamu ingin tahu raw chat response juga, pastikan kamu simpan chatResp di outer scope:
		// log.Printf("DEBUG: raw LLM response: %s", chatResp.Choices[0].Message.Content)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "😓 Aduh, sistem lagi ngambek. Coba lagi ya~",
		})
	}

	filterCategories, err := h.categoryService.GetDefaultCategoriesByType(
        model.CategoryType(typeRes.TransactionType),
    )
	if err != nil {
		log.Printf("Load Categories: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "😓 Aduh, sistem lagi ngambek. Coba lagi lagi ya~"})
    }

	// Step 1: Klasifikasi menggunakan LLM
	result, llmResp, err := h.aiService.ClassifyTransactionFull(
        input.Prompt,
        typeRes.TransactionType,
        filterCategories,
    )
    if err != nil {
		log.Printf("Klasifikasi Full: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "😓 Aduh, sistem lagi ngambek. Coba lagi lagi ya~",
		})
    }

	var transactionDate time.Time
	if result.Date != "" {
		parsedDate, _ := time.Parse("2006-01-02", result.Date)
		transactionDate = parsedDate
	} else {
		transactionDate = time.Now() // jika kosong, fallback ke sekarang
	}

	// Step 3: Simpan transaksi
	tx := &model.Transaction{
		ChatID:          chatId,
		OriginalText:    input.Prompt,
		TransactionType: model.TransactionType(result.TransactionType),
		Amount:          result.Amount,
		CategoryID:      result.CategoryID,
		TransactionDate: transactionDate,
	}

	if err := h.transactionService.CreateTransaction(tx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan transaksi",
		})
	}

	txWithCategory, err := h.transactionService.GetTransactionWithCategory(tx.ID)
	if err != nil {
		// Log error tapi tetap return success karena create sudah berhasil
		log.Printf("Failed to load category: %v", err)
		txWithCategory = tx // fallback ke transaction tanpa category
	}

	// Step 4: Respon sukses
	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{
		"message": "Transaksi berhasil disimpan",
		"usage":   llmResp.Usage,
		"data":    txWithCategory,
	})
}