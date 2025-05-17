package handler

import (
	"finance-bot/service"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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
	
	user, err := h.userService.GetUser(userId)
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

