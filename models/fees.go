package models

type FeesLimit struct {
	ID                int      `json:"id" db:"id"`
	CurrencyID        *int     `json:"currency_id" db:"currency_id"`
	TransactionTypeID *int     `json:"transaction_type_id" db:"transaction_type_id"`
	PaymentMethodID   *int     `json:"payment_method_id" db:"payment_method_id"`
	ChargePercentage  float64  `json:"charge_percentage" db:"charge_percentage"`
	ChargeFixed       float64  `json:"charge_fixed" db:"charge_fixed"`
	MinLimit          float64  `json:"min_limit" db:"min_limit"`
	MaxLimit          *float64 `json:"max_limit" db:"max_limit"`
	ProcessingTime    string   `json:"processing_time" db:"processing_time"` // varchar(4)
	HasTransaction    string   `json:"has_transaction" db:"has_transaction"`   // varchar(3) - Yes or No
}

type FeesExpress struct {
	ID                int     `json:"id" db:"id"`
	TransactionTypeID *int    `json:"transaction_type_id" db:"transaction_type_id"`
	ChargePercentage  float64 `json:"charge_percentage" db:"charge_percentage"`
	ChargeFixed       float64 `json:"charge_fixed" db:"charge_fixed"`
}

