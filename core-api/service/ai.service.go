package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"finance-bot/model"
	"finance-bot/static"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// LLMProvider represents different LLM providers
type LLMProvider string

const (
	ProviderDeepSeek LLMProvider = "deepseek"
	ProviderChatGPT  LLMProvider = "chatgpt"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model          string         `json:"model"`
	Messages       []ChatMessage  `json:"messages"`
	MaxTokens      int            `json:"max_tokens,omitempty"`
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

// ClassificationResult represents the result of transaction classification
type ClassificationResult struct {
	TransactionType string  `json:"type"`
	Amount          float64 `json:"amount"`
	CategoryID      *uint   `json:"category_id"`
	Date            string  `json:"date"`
}

// TypeClassificationResult represents the result of type classification only
type TypeClassificationResult struct {
	TransactionType string `json:"type"`
}

// LLMService handles LLM operations
type LLMService struct {
	Provider   LLMProvider
	APIKey     string
	APIURL     string
	Model      string
	Categories []model.Category
}

// NewLLMService creates a new LLM service instance
func NewLLMService(provider LLMProvider, apiKey, apiURL, model string) *LLMService {
	return &LLMService{
		Provider:   provider,
		APIKey:     apiKey,
		APIURL:     apiURL,
		Model:      model,
	}
}

// PromptTemplate represents different prompt templates
type PromptTemplate struct {
	TypeClassification string
	FullClassification string
}

// GetPromptTemplates returns prompt templates
func (s *LLMService) GetPromptTemplates() PromptTemplate {
	return PromptTemplate{
		TypeClassification: static.PromptTypeClassification,
		FullClassification: static.PromptFullClassification,
	}
}

// ClassifyTransactionType classifies transaction type only (INCOME/EXPENSE)
func (s *LLMService) ClassifyTransactionType(description string) (*TypeClassificationResult, *ChatResponse, error) {
	templates := s.GetPromptTemplates()
	prompt := fmt.Sprintf(
		templates.TypeClassification,
		description,
	)
	
	chatResp, err := s.hitLLM(prompt)
	if err != nil {
		return nil, chatResp, err
	}
	
	var result TypeClassificationResult
	err = json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, chatResp, err
	}
	
	return &result, chatResp, nil
}

// ClassifyTransactionFull performs full transaction classification
func (s *LLMService) ClassifyTransactionFull(description string, transactionType string, categories []model.Category) (*ClassificationResult, *ChatResponse, error) {
	currentDate := time.Now().Format("2006-01-02")
	templates := s.GetPromptTemplates()

	var catsBuilder strings.Builder
	for _, cat := range categories {
		catsBuilder.WriteString(fmt.Sprintf("#%d %s\n", cat.ID, cat.Name))
	}
	categoriesStr := catsBuilder.String()

	// placeholder value:
	// 1) tanggal, 2) deskripsi, 3) kalimat tipe, 4) JSON tipe, 5) daftar kategori
	prompt := fmt.Sprintf(
		templates.FullClassification,
		currentDate,
		description,
		transactionType,
		transactionType,
		categoriesStr,
	)
	
	chatResp, err := s.hitLLM(prompt)
	if err != nil {
		return nil, chatResp, err
	}
	
	var result ClassificationResult
	err = json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, chatResp, err
	}
	
	return &result, chatResp, nil
}

// hitLLM is the unified method to hit different LLM providers
func (s *LLMService) hitLLM(prompt string) (*ChatResponse, error) {
	switch s.Provider {
	case ProviderDeepSeek:
		return s.hitDeepSeek(prompt)
	case ProviderChatGPT:
		return s.hitChatGPT(prompt)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", s.Provider)
	}
}

// hitDeepSeek handles DeepSeek API calls
func (s *LLMService) hitDeepSeek(prompt string) (*ChatResponse, error) {
	requestBody := ChatRequest{
		Model: s.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		ResponseFormat: map[string]any{
			"type": "json_object",
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.APIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chatResp ChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return nil, err
	}

	// Tampilkan prompt yang dikirim ke LLM (string)
	fmt.Printf("DEBUG: prompt = %q\n", prompt)

	// Tampilkan konten respon LLM (string)
	fmt.Printf("DEBUG: response content = %s\n", chatResp.Choices[0].Message.Content)

	// Tampilkan total token yang terpakai (integer)
	fmt.Printf("DEBUG: total tokens used = %d\n", chatResp.Usage.TotalTokens)

	return &chatResp, nil
}

// hitChatGPT handles ChatGPT API calls
func (s *LLMService) hitChatGPT(prompt string) (*ChatResponse, error) {
	client := openai.NewClient(
		option.WithAPIKey(s.APIKey),
	)

	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: openai.ChatModel(s.Model),
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{
					Type: "json_object",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.MarshalIndent(chatCompletion, "", "  ")
	if err != nil {
		return nil, err
	}

	var chatResp ChatResponse
	err = json.Unmarshal(jsonBytes, &chatResp)
	if err != nil {
		return nil, err
	}

		// Tampilkan prompt yang dikirim ke LLM (string)
	fmt.Printf("DEBUG: prompt = %q\n", prompt)

	// Tampilkan konten respon LLM (string)
	fmt.Printf("DEBUG: response content = %s\n", chatResp.Choices[0].Message.Content)

	// Tampilkan total token yang terpakai (integer)
	fmt.Printf("DEBUG: total tokens used = %d\n", chatResp.Usage.TotalTokens)

	return &chatResp, nil
}