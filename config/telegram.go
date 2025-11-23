package config

import "os"

type TelegramConfig struct {
	Token  string
	ChatID string
}

func GetTelegramConfig() *TelegramConfig {
	return &TelegramConfig{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
	}
}

