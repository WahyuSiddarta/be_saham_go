package models

import "time"

// Stock represents a stock entity
type Stock struct {
	ID            int       `json:"id"`
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
	Volume        int64     `json:"volume"`
	MarketCap     float64   `json:"market_cap"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// StockPrice represents price history
type StockPrice struct {
	ID        int       `json:"id"`
	StockID   int       `json:"stock_id"`
	Price     float64   `json:"price"`
	Volume    int64     `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
