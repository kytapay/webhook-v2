package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type PakaiLinkConfig struct {
	ClientSecret    string
	RSAPublicKey    *rsa.PublicKey
	RSAPublicKeyPath string
}

var pakaiLinkConfig *PakaiLinkConfig

func GetPakaiLinkConfig() *PakaiLinkConfig {
	if pakaiLinkConfig == nil {
		config := &PakaiLinkConfig{
			ClientSecret:     os.Getenv("PAKAILINK_CLIENT_SECRET"),
			RSAPublicKeyPath: os.Getenv("PAKAILINK_RSA_PUBLIC_KEY_PATH"),
		}

		// Load RSA public key if path is provided
		if config.RSAPublicKeyPath != "" {
			publicKey, err := loadRSAPublicKey(config.RSAPublicKeyPath)
			if err == nil {
				config.RSAPublicKey = publicKey
			}
		}

		pakaiLinkConfig = config
	}
	return pakaiLinkConfig
}

// loadRSAPublicKey loads RSA public key from PEM file
func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	return rsaPub, nil
}

