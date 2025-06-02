package handlers

import (
	"context"
	"time"

	"gateway/model"
	"gateway/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewAuthHandler(db *gorm.DB, redis *redis.Client) *AuthHandler {
	return &AuthHandler{DB: db, Redis: redis}
}

func (h *AuthHandler) MagicLogin(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	ctx := context.Background()

	// Ambil userID dari Redis
	userID, err := h.Redis.Get(ctx, "magic_login:"+req.Token).Result()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}
	h.Redis.Del(ctx, "magic_login:"+req.Token)

	// Buat access token
	accessToken, err := utils.GenerateJWT(userID, time.Minute*15)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	// Buat refresh token ID
	refreshTokenID := uuid.NewString()

	// Simpan refresh token ID di Redis
	err = h.Redis.Set(ctx, "refresh_token:"+refreshTokenID, userID, time.Hour*24*7).Err()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to store refresh token")
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshTokenID,
	})
}	

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	ctx := context.Background()

	// Verifikasi refresh token di Redis
	userID, err := h.Redis.Get(ctx, "refresh_token:"+req.RefreshToken).Result()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}

	// Hapus refresh token lama untuk mencegah reuse
	h.Redis.Del(ctx, "refresh_token:"+req.RefreshToken)

	// Buat token baru
	accessToken, err := utils.GenerateJWT(userID, time.Minute*15)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	// Buat refresh token baru (rotating token)
	newRefreshTokenID := uuid.NewString()
	err = h.Redis.Set(ctx, "refresh_token:"+newRefreshTokenID, userID, time.Hour*24*7).Err()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to store new refresh token")
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshTokenID,
	})
}


func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var user model.User
	if err := h.DB.First(&user, "chat_id = ?", userID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	return c.JSON(user)
}
