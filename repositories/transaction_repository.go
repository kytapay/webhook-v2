package repositories

import (
	"database/sql"
	"time"

	"github.com/kytapay/webhook-v2/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// GetTransactionByGrantID gets transaction info by grant_id
func (r *TransactionRepository) GetTransactionByGrantID(grantID string) (*models.TransactionInfo, error) {
	query := `SELECT id, app_id, order_id, payment_method, amount, currency, notify_url, success_url, cancel_url, grant_id, token, bank_number, bank_ewallet_name, qris_string, ewallet_link, status, version, created_at, updated_at 
		FROM app_transactions_infos WHERE grant_id = ? LIMIT 1`

	var transaction models.TransactionInfo
	err := r.db.QueryRow(query, grantID).Scan(
		&transaction.ID,
		&transaction.AppID,
		&transaction.OrderID,
		&transaction.PaymentMethod,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.NotifyURL,
		&transaction.SuccessURL,
		&transaction.CancelURL,
		&transaction.GrantID,
		&transaction.Token,
		&transaction.BankNumber,
		&transaction.BankEwalletName,
		&transaction.QrisString,
		&transaction.EwalletLink,
		&transaction.Status,
		&transaction.Version,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransaction updates transaction status and amount
func (r *TransactionRepository) UpdateTransaction(grantID string, status string, amount int64) error {
	query := `UPDATE app_transactions_infos SET status = ?, amount = ?, updated_at = ? WHERE grant_id = ?`
	now := time.Now()
	_, err := r.db.Exec(query, status, amount, now, grantID)
	return err
}

// GetTransactionsByGrantID gets transactions table record by grant_id
func (r *TransactionRepository) GetTransactionsByGrantID(grantID string) (*models.Transactions, error) {
	query := `SELECT id, user_id, currency_id, payment_method_id, merchant_id, uuid, grant_id, transaction_reference_id, transaction_type_id, user_type, subtotal, percentage, charge_percentage, charge_fixed, total, payment_status, status, created_at, updated_at 
		FROM transactions WHERE grant_id = ? LIMIT 1`

	var transaction models.Transactions
	err := r.db.QueryRow(query, grantID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.CurrencyID,
		&transaction.PaymentMethodID,
		&transaction.MerchantID,
		&transaction.UUID,
		&transaction.GrantID,
		&transaction.TransactionReferenceID,
		&transaction.TransactionTypeID,
		&transaction.UserType,
		&transaction.Subtotal,
		&transaction.Percentage,
		&transaction.ChargePercentage,
		&transaction.ChargeFixed,
		&transaction.Total,
		&transaction.PaymentStatus,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// CreateTransaction creates a new transaction record
func (r *TransactionRepository) CreateTransaction(transaction models.TransactionsData) error {
	query := `INSERT INTO transactions 
		(user_id, currency_id, payment_method_id, merchant_id, uuid, grant_id, transaction_reference_id, transaction_type_id, user_type, subtotal, percentage, charge_percentage, charge_fixed, total, payment_status, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	_, err := r.db.Exec(query,
		transaction.UserID,
		transaction.CurrencyID,
		transaction.PaymentMethodID,
		transaction.MerchantID,
		transaction.UUID,
		transaction.GrantID,
		transaction.TransactionReferenceID,
		transaction.TransactionTypeID,
		transaction.UserType,
		transaction.Subtotal,
		transaction.Percentage,
		transaction.ChargePercentage,
		transaction.ChargeFixed,
		transaction.Total,
		transaction.PaymentStatus,
		transaction.Status,
		now,
		now,
	)

	return err
}

// UpdateTransactions updates transactions status
func (r *TransactionRepository) UpdateTransactions(grantID string, paymentStatus, status string) error {
	query := `UPDATE transactions SET payment_status = ?, status = ?, updated_at = ? WHERE grant_id = ?`
	now := time.Now()
	_, err := r.db.Exec(query, paymentStatus, status, now, grantID)
	return err
}

