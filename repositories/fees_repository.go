package repositories

import (
	"database/sql"

	"github.com/kytapay/webhook-v2/models"
)

type FeesRepository struct {
	db *sql.DB
}

func NewFeesRepository(db *sql.DB) *FeesRepository {
	return &FeesRepository{db: db}
}

// GetFeesLimit gets fees limit by transaction_type_id and payment_method_id
func (r *FeesRepository) GetFeesLimit(transactionTypeID, paymentMethodID int) (*models.FeesLimit, error) {
	query := `SELECT id, currency_id, transaction_type_id, payment_method_id, charge_percentage, charge_fixed, min_limit, max_limit, processing_time, has_transaction 
		FROM fees_limits WHERE transaction_type_id = ? AND payment_method_id = ? LIMIT 1`

	var fee models.FeesLimit
	err := r.db.QueryRow(query, transactionTypeID, paymentMethodID).Scan(
		&fee.ID,
		&fee.CurrencyID,
		&fee.TransactionTypeID,
		&fee.PaymentMethodID,
		&fee.ChargePercentage,
		&fee.ChargeFixed,
		&fee.MinLimit,
		&fee.MaxLimit,
		&fee.ProcessingTime,
		&fee.HasTransaction,
	)

	if err != nil {
		return nil, err
	}

	return &fee, nil
}

// GetFeesExpress gets fees express by transaction_type_id
func (r *FeesRepository) GetFeesExpress(transactionTypeID int) (*models.FeesExpress, error) {
	query := `SELECT id, transaction_type_id, charge_percentage, charge_fixed 
		FROM fees_express WHERE transaction_type_id = ? LIMIT 1`

	var fee models.FeesExpress
	err := r.db.QueryRow(query, transactionTypeID).Scan(
		&fee.ID,
		&fee.TransactionTypeID,
		&fee.ChargePercentage,
		&fee.ChargeFixed,
	)

	if err != nil {
		return nil, err
	}

	return &fee, nil
}

