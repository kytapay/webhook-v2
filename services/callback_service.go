package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kytapay/webhook-v2/models"
)

type CallbackService struct {
	client *http.Client
}

func NewCallbackService() *CallbackService {
	return &CallbackService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendCallback sends callback to merchant
func (cs *CallbackService) SendCallback(url string, payload interface{}, token *string) (int, string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if token != nil && *token != "" {
		req.Header.Set("X-CALLBACK-TOKEN", *token)
	}

	resp, err := cs.client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}

	return resp.StatusCode, string(body), nil
}

// BuildPayloadV2 builds payload for version 2
func BuildPayloadV2(transaction *models.TransactionInfo, paymentID, merchantName, status, date string) map[string]interface{} {
	payloads := map[string]interface{}{
		"QRIS": map[string]interface{}{
			"callback_code":    "2001200",
			"callback_message": "Callback Payment Status",
			"callback_data": map[string]interface{}{
				"id":           paymentID,
				"reference_id": transaction.OrderID,
				"amount":       int(transaction.Amount),
				"status":       status,
				"payment_type": "QR",
				"payment_data": map[string]interface{}{
					"qr_string": transaction.QrisString,
				},
				"merchant_url": map[string]interface{}{
					"notify_url":  transaction.NotifyURL,
					"success_url": transaction.SuccessURL,
					"failed_url":  transaction.CancelURL,
				},
				"callback_time": date,
			},
		},
		"VA": map[string]interface{}{
			"callback_code":    "2001200",
			"callback_message": "Callback Payment Status",
			"callback_data": map[string]interface{}{
				"id":           paymentID,
				"reference_id": transaction.OrderID,
				"amount":       int(transaction.Amount),
				"status":       status,
				"payment_type": "VIRTUAL_ACCOUNT",
				"payment_data": map[string]interface{}{
					"bank_code":      transaction.BankEwalletName,
					"account_number": transaction.BankNumber,
					"account_name":   fmt.Sprintf("%s", merchantName),
				},
				"merchant_url": map[string]interface{}{
					"notify_url":  transaction.NotifyURL,
					"success_url": transaction.SuccessURL,
					"failed_url":  transaction.CancelURL,
				},
				"callback_time": date,
			},
		},
		"EWALLET": map[string]interface{}{
			"callback_code":    "2001200",
			"callback_message": "Callback Payment Status",
			"callback_data": map[string]interface{}{
				"id":           paymentID,
				"reference_id": transaction.OrderID,
				"amount":       int(transaction.Amount),
				"status":       status,
				"payment_type": "E-WALLET",
				"payment_data": map[string]interface{}{
					"channel_code": transaction.BankEwalletName,
					"redirect_url": transaction.EwalletLink,
				},
				"merchant_url": map[string]interface{}{
					"notify_url":  transaction.NotifyURL,
					"success_url": transaction.SuccessURL,
					"failed_url":  transaction.CancelURL,
				},
				"callback_time": date,
			},
		},
	}

	return payloads
}

// BuildPayloadV2Payout builds payload for version 2 payout callback
func BuildPayloadV2Payout(transaction *models.TransactionInfo, paymentID, status, date string) map[string]interface{} {
	payloads := map[string]interface{}{
		"PAYOUTS": map[string]interface{}{
			"callback_code":    "2001400",
			"callback_message": "Callback Payout Status",
			"callback_data": map[string]interface{}{
				"id":           paymentID,
				"reference_id": transaction.OrderID,
				"amount":       int(transaction.Amount),
				"status":       status,
				"payout_data": map[string]interface{}{
					"code":          transaction.PaymentMethod,
					"account_number": transaction.BankNumber,
					"account_name":   transaction.BankEwalletName,
				},
				"merchant_url": map[string]interface{}{
					"notify_url": transaction.NotifyURL,
				},
				"callback_time": date,
			},
		},
	}

	return payloads
}

