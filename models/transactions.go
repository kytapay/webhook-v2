package models

import "time"

type Transactions struct {
	ID                    int        `json:"id" db:"id"`
	UserID                *int       `json:"user_id" db:"user_id"`
	CurrencyID            *int       `json:"currency_id" db:"currency_id"`
	PaymentMethodID       *int       `json:"payment_method_id" db:"payment_method_id"`
	MerchantID            *int       `json:"merchant_id" db:"merchant_id"`
	UUID                  *string    `json:"uuid" db:"uuid"`
	GrantID               *string    `json:"grant_id" db:"grant_id"`
	TransactionReferenceID int       `json:"transaction_reference_id" db:"transaction_reference_id"`
	TransactionTypeID     *int       `json:"transaction_type_id" db:"transaction_type_id"`
	UserType              string     `json:"user_type" db:"user_type"`
	Subtotal              float64    `json:"subtotal" db:"subtotal"`
	Percentage            float64    `json:"percentage" db:"percentage"`
	ChargePercentage      float64    `json:"charge_percentage" db:"charge_percentage"`
	ChargeFixed           float64    `json:"charge_fixed" db:"charge_fixed"`
	Total                 float64    `json:"total" db:"total"`
	PaymentStatus         *string    `json:"payment_status" db:"payment_status"`
	Status                string     `json:"status" db:"status"`
	CreatedAt             *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at" db:"updated_at"`
}

type TransactionsData struct {
	UserID                *int
	CurrencyID            *int
	PaymentMethodID       *int
	MerchantID            *int
	UUID                  *string
	GrantID               *string
	TransactionReferenceID int
	TransactionTypeID     *int
	UserType              string
	Subtotal              float64
	Percentage            float64
	ChargePercentage      float64
	ChargeFixed           float64
	Total                 float64
	PaymentStatus         *string
	Status                string
}

