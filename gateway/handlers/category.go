package handlers

import (
	"gateway/model"
	"gateway/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

// GetCategories returns all categories for the authenticated user.
func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userIDStr, _ := userIDRaw.(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	var categories []model.Category

	// Query categories that belong to the user
	if err := h.DB.Where("user_id = ? OR user_id IS NULL", userID).Order("name ASC").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return utils.JSONResponse(c, fiber.StatusOK, categories)
}
