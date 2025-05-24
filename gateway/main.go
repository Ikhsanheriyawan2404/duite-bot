package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	coreApiUrl := os.Getenv("CORE_API_URL")
	if coreApiUrl == "" {
		log.Fatal("CORE_API_URL is not set")
	}

	app := fiber.New()

	// Serve static HTML from ./public
	app.Static("/", "./public")

	// 1. GET or POST /api/users/:id
	app.All("/api/users/:id", func(c *fiber.Ctx) error {
		userId := c.Params("id")
		fullURL := coreApiUrl + "/users/" + userId
		method := c.Method()

		log.Printf("-> %s %s → %s", method, c.OriginalURL(), fullURL)

		reqBody := c.Body()
		req, err := http.NewRequest(method, fullURL, strings.NewReader(string(reqBody)))
		if err != nil {
			return c.Status(500).SendString("Failed to build request")
		}

		// Copy headers
		c.Request().Header.VisitAll(func(k, v []byte) {
			req.Header.Set(string(k), string(v))
		})
		req.URL.RawQuery = c.Context().QueryArgs().String()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error calling core-api:", err)
			return c.Status(502).JSON(fiber.Map{"error": "Failed to reach core-api"})
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		return c.Status(resp.StatusCode).Send(body)
	})

	// 2. GET or POST /api/users/:id/transactions
	app.All("/api/users/:id/transactions", func(c *fiber.Ctx) error {
		userId := c.Params("id")
		fullURL := coreApiUrl + "/users/" + userId + "/transactions"
		method := c.Method()

		log.Printf("-> %s %s → %s", method, c.OriginalURL(), fullURL)

		reqBody := c.Body()
		req, err := http.NewRequest(method, fullURL, strings.NewReader(string(reqBody)))
		if err != nil {
			return c.Status(500).SendString("Failed to build request")
		}

		c.Request().Header.VisitAll(func(k, v []byte) {
			req.Header.Set(string(k), string(v))
		})
		req.URL.RawQuery = c.Context().QueryArgs().String()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error calling core-api:", err)
			return c.Status(502).JSON(fiber.Map{"error": "Failed to reach core-api"})
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		return c.Status(resp.StatusCode).Send(body)
	})

	log.Println("Proxy running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
