package config

import "os"

type LinkQuConfig struct {
	ClientID     string
	ClientSecret string
}

func GetLinkQuConfig() *LinkQuConfig {
	return &LinkQuConfig{
		ClientID:     os.Getenv("LINKQU_CLIENT_ID"),
		ClientSecret: os.Getenv("LINKQU_CLIENT_SECRET"),
	}
}

