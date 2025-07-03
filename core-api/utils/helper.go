package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

func Slugify(s string) string {
    s = strings.ToLower(s)
    var b strings.Builder
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            b.WriteRune(r)
        } else if unicode.IsSpace(r) || r == '-' {
            b.WriteRune('-')
        }
        // abaikan karakter lain
    }
    return strings.Trim(b.String(), "-")
}

func ContainsNominal(s string) bool {
	re := regexp.MustCompile(`\d+`)
	return re.MatchString(s)
}

func GenerateRandomToken(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return "tok_" + hex.EncodeToString(bytes)
}

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
