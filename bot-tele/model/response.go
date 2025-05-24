package model

type LLMResult struct {
	Result struct {
		TransactionType string  `json:"type"`
		Amount          float64 `json:"amount"`
		Category        string  `json:"category"`
		Date            string  `json:"date"`
	} `json:"result"`
}