package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bot-tele/config"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

var Client *APIClient

func NewAPIClient() *APIClient {
	return &APIClient{
		BaseURL: config.AppConfig.CoreApiUrl,
		Client:  &http.Client{},
	}
}

func InitClient() {
	Client = NewAPIClient()
}

func (c *APIClient) Request(method string, path string, reqBody any, result any) error {
	// Encode request body if provided
	var body io.Reader
	if reqBody != nil {
		jsonBytes, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	// Build full URL
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	fmt.Println("Request Bash:",  c.BaseURL)
	fmt.Println("Request URL:", fullURL)

	// Create HTTP request
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP error status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// Decode response into result
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
