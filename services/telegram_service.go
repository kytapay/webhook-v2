package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kytapay/webhook-v2/config"
)

type TelegramService struct {
	client *http.Client
	config *config.TelegramConfig
}

func NewTelegramService() *TelegramService {
	return &TelegramService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: config.GetTelegramConfig(),
	}
}

// SendMessage sends a message to Telegram
func (ts *TelegramService) SendMessage(message, parseMode string) error {
	if ts.config.Token == "" || ts.config.ChatID == "" {
		return nil // Skip if not configured
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", ts.config.Token)

	payload := map[string]interface{}{
		"chat_id":    ts.config.ChatID,
		"text":       message,
		"parse_mode": parseMode,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	return err
}

