package repositories

import (
	"database/sql"

	"github.com/kytapay/webhook-v2/models"
)

type MerchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

// GetMerchantPaymentByGatewayRef gets merchant payment by gateway_reference
func (r *MerchantRepository) GetMerchantPaymentByGatewayRef(gatewayRef string) (*models.MerchantPayment, error) {
	query := `SELECT id, merchant_id, payment_method_id, gateway_reference, order_no, uuid, fee_bearer, percentage, charge_percentage, charge_fixed, amount, total, status, created_at, updated_at 
		FROM merchant_payments WHERE gateway_reference = ? LIMIT 1`

	var payment models.MerchantPayment
	err := r.db.QueryRow(query, gatewayRef).Scan(
		&payment.ID,
		&payment.MerchantID,
		&payment.PaymentMethodID,
		&payment.GatewayReference,
		&payment.OrderNo,
		&payment.UUID,
		&payment.FeeBearer,
		&payment.Percentage,
		&payment.ChargePercentage,
		&payment.ChargeFixed,
		&payment.Amount,
		&payment.Total,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// UpdateMerchantPayment updates merchant payment status and amount
func (r *MerchantRepository) UpdateMerchantPayment(gatewayRef string, status string, amount float64) error {
	query := `UPDATE merchant_payments SET status = ?, amount = ?, updated_at = NOW() WHERE gateway_reference = ?`
	_, err := r.db.Exec(query, status, amount, gatewayRef)
	return err
}

// GetMerchantByID gets merchant by ID
func (r *MerchantRepository) GetMerchantByID(merchantID int) (*models.Merchant, error) {
	query := `SELECT id, user_id, business_name, merchant_uuid, site_url, status FROM merchants WHERE id = ? LIMIT 1`

	var merchant models.Merchant
	err := r.db.QueryRow(query, merchantID).Scan(
		&merchant.ID,
		&merchant.UserID,
		&merchant.BusinessName,
		&merchant.MerchantUUID,
		&merchant.SiteURL,
		&merchant.Status,
	)

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

// GetMerchantPayoutByGatewayRef gets merchant payout by gateway_reference
func (r *MerchantRepository) GetMerchantPayoutByGatewayRef(gatewayRef string) (*models.MerchantPayout, error) {
	query := `SELECT id, merchant_id, currency_id, payment_method_id, user_id, gateway_reference, order_no, item_name, uuid, fee_bearer, percentage, charge_percentage, charge_fixed, amount, total, status, bank_name, account_name, account_number, created_at, updated_at 
		FROM merchant_payouts WHERE gateway_reference = ? LIMIT 1`

	var payout models.MerchantPayout
	var currencyID, userID sql.NullInt64
	var itemName, bankName, accountName, accountNumber sql.NullString

	err := r.db.QueryRow(query, gatewayRef).Scan(
		&payout.ID,
		&payout.MerchantID,
		&currencyID,
		&payout.PaymentMethodID,
		&userID,
		&payout.GatewayReference,
		&payout.OrderNo,
		&itemName,
		&payout.UUID,
		&payout.FeeBearer,
		&payout.Percentage,
		&payout.ChargePercentage,
		&payout.ChargeFixed,
		&payout.Amount,
		&payout.Total,
		&payout.Status,
		&bankName,
		&accountName,
		&accountNumber,
		&payout.CreatedAt,
		&payout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if currencyID.Valid {
		val := int(currencyID.Int64)
		payout.CurrencyID = &val
	}
	if userID.Valid {
		val := int(userID.Int64)
		payout.UserID = &val
	}
	if itemName.Valid {
		payout.ItemName = &itemName.String
	}
	if bankName.Valid {
		payout.BankName = &bankName.String
	}
	if accountName.Valid {
		payout.AccountName = &accountName.String
	}
	if accountNumber.Valid {
		payout.AccountNumber = &accountNumber.String
	}

	return &payout, nil
}

// UpdateMerchantPayout updates merchant payout status and amount
func (r *MerchantRepository) UpdateMerchantPayout(gatewayRef string, status string, amount float64) error {
	query := `UPDATE merchant_payouts SET status = ?, amount = ?, updated_at = NOW() WHERE gateway_reference = ?`
	_, err := r.db.Exec(query, status, amount, gatewayRef)
	return err
}

