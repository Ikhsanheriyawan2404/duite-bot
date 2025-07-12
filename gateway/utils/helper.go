package utils

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func JSONResponse(c *fiber.Ctx, statusCode int, data interface{}) error {
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    encoder.SetEscapeHTML(false)
    
    if err := encoder.Encode(data); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to encode response",
        })
    }
    
    // Remove trailing newline
    result := buf.Bytes()
    if len(result) > 0 && result[len(result)-1] == '\n' {
        result = result[:len(result)-1]
    }
    
    c.Set("Content-Type", "application/json")
    return c.Status(statusCode).Send(result)
}