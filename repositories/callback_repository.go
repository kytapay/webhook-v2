package repositories

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/kytapay/webhook-v2/models"
)

type CallbackRepository struct {
	db *sql.DB
}

func NewCallbackRepository(db *sql.DB) *CallbackRepository {
	return &CallbackRepository{db: db}
}

// GetCallbackByTransactionInfoID gets callback by transaction_info_id
func (r *CallbackRepository) GetCallbackByTransactionInfoID(transactionInfoID int) (*models.CallbackStatus, error) {
	query := `SELECT id, transaction_info_id, merchant_id, notify_url, status, error_message, response_body, payload, retry_count, created_at, updated_at 
		FROM callback_status WHERE transaction_info_id = ? LIMIT 1`

	var callback models.CallbackStatus
	var errorMsg, responseBody, payload sql.NullString

	err := r.db.QueryRow(query, transactionInfoID).Scan(
		&callback.ID,
		&callback.TransactionInfoID,
		&callback.MerchantID,
		&callback.NotifyURL,
		&callback.Status,
		&errorMsg,
		&responseBody,
		&payload,
		&callback.RetryCount,
		&callback.CreatedAt,
		&callback.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if errorMsg.Valid {
		callback.ErrorMessage = &errorMsg.String
	}
	if responseBody.Valid {
		callback.ResponseBody = &responseBody.String
	}
	if payload.Valid {
		callback.Payload = &payload.String
	}

	return &callback, nil
}

// UpdateCallback updates callback status
func (r *CallbackRepository) UpdateCallback(transactionInfoID int, status, errorMessage, responseBody string, payloadData interface{}) error {
	payloadJSON, _ := json.Marshal(payloadData)
	payloadStr := string(payloadJSON)

	query := `UPDATE callback_status SET status = ?, error_message = ?, response_body = ?, payload = ?, updated_at = ? WHERE transaction_info_id = ?`
	now := time.Now()
	_, err := r.db.Exec(query, status, errorMessage, responseBody, payloadStr, now, transactionInfoID)
	return err
}

