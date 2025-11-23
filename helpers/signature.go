package helpers

import (
	"crypto"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kytapay/webhook-v2/config"
)

// VerifyLinkQuSignature verifies LinkQu signature from callback
// LinkQu menggunakan client-id dan client-secret di header untuk validasi
func VerifyLinkQuSignature(clientID, clientSecret string) bool {
	linkQuConfig := config.GetLinkQuConfig()
	return clientID == linkQuConfig.ClientID && clientSecret == linkQuConfig.ClientSecret
}

// VerifyPakaiLinkSignature verifies PakaiLink signature for webhook callback
// Supports both symmetric (HMAC SHA-512) and asymmetric (RSA SHA-256) signatures
// Format: <HTTP METHOD> + ":" + <PATH URL CALLBACK> + ":" + LowerCase(HexEncode(SHA-256(Minify(<HTTP BODY>)))) + ":" + <X-TIMESTAMP>
func VerifyPakaiLinkSignature(method, path, body, timestamp, signature string) bool {
	pakaiLinkConfig := config.GetPakaiLinkConfig()

	// Minify request body
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		return false
	}

	minifiedBody, err := json.Marshal(bodyMap)
	if err != nil {
		return false
	}

	// SHA-256 hash of minified body
	hash := sha256.Sum256(minifiedBody)
	hashHex := strings.ToLower(hex.EncodeToString(hash[:]))

	// Compose string to sign: METHOD:PATH:HASH:TIMESTAMP
	stringToSign := fmt.Sprintf("%s:%s:%s:%s", method, path, hashHex, timestamp)

	// Try RSA verification first if public key is available
	if pakaiLinkConfig.RSAPublicKey != nil {
		return verifyRSASignature(stringToSign, signature, pakaiLinkConfig.RSAPublicKey)
	}

	// Fallback to symmetric signature (HMAC SHA-512)
	return verifySymmetricSignature(stringToSign, signature, pakaiLinkConfig.ClientSecret)
}

// verifySymmetricSignature verifies symmetric signature using HMAC SHA-512
func verifySymmetricSignature(stringToSign, signature, secret string) bool {
	// HMAC SHA-512
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	expectedHash := mac.Sum(nil)

	// Encode expected hash to Base64
	expectedSignature := base64.StdEncoding.EncodeToString(expectedHash)

	// Compare signatures (both are base64 encoded)
	return signature == expectedSignature
}

// verifyRSASignature verifies RSA signature using SHA-256
func verifyRSASignature(stringToSign string, signature string, publicKey *rsa.PublicKey) bool {
	// Decode base64 signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	// Hash the string to sign
	hashed := sha256.Sum256([]byte(stringToSign))

	// Verify signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
	return err == nil
}

// NormalizeStatus converts various status formats to standard format
func NormalizeStatus(status string) string {
	successStatuses := []string{"SUCCEEDED", "PAID", "CAPTURED", "SUCCESS", "COMPLETED", "00"}
	expiredStatuses := []string{"EXPIRED", "CANCELLED", "VOIDED", "FAILED"}
	pendingStatuses := []string{"PENDING", "IN_PROGRESS", "OPEN"}

	statusUpper := strings.ToUpper(status)

	for _, s := range successStatuses {
		if statusUpper == s {
			return "success"
		}
	}

	for _, e := range expiredStatuses {
		if statusUpper == e {
			return "expires"
		}
	}

	for _, p := range pendingStatuses {
		if statusUpper == p {
			return "pending"
		}
	}

	return "success" // Default
}

// MerchantNormalizeStatus converts status to merchant format
func MerchantNormalizeStatus(status string) string {
	successStatuses := []string{"SUCCEEDED", "PAID", "CAPTURED", "SUCCESS", "COMPLETED", "00"}
	expiredStatuses := []string{"EXPIRED", "CANCELLED", "VOIDED", "FAILED"}
	pendingStatuses := []string{"PENDING", "IN_PROGRESS", "OPEN"}

	statusUpper := strings.ToUpper(status)

	for _, s := range successStatuses {
		if statusUpper == s {
			return "Success"
		}
	}

	for _, e := range expiredStatuses {
		if statusUpper == e {
			return "Blocked"
		}
	}

	for _, p := range pendingStatuses {
		if statusUpper == p {
			return "Pending"
		}
	}

	return "Success" // Default
}

