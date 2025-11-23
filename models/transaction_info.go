package models

import "time"

type TransactionInfo struct {
	ID              int        `json:"id" db:"id"`
	AppID           int        `json:"app_id" db:"app_id"`
	OrderID         string     `json:"order_id" db:"order_id"`
	PaymentMethod   string     `json:"payment_method" db:"payment_method"`
	Amount          int64      `json:"amount" db:"amount"`
	Currency        string     `json:"currency" db:"currency"`
	NotifyURL       string     `json:"notify_url" db:"notify_url"`
	SuccessURL      *string    `json:"success_url" db:"success_url"`
	CancelURL       *string    `json:"cancel_url" db:"cancel_url"`
	GrantID         string     `json:"grant_id" db:"grant_id"`
	Token           *string    `json:"token" db:"token"`
	BankNumber      *string    `json:"bank_number" db:"bank_number"`
	BankEwalletName *string    `json:"bank_ewallet_name" db:"bank_ewallet_name"`
	QrisString      *string    `json:"qris_string" db:"qris_string"`
	EwalletLink     *string    `json:"ewallet_link" db:"ewallet_link"`
	Status          string     `json:"status" db:"status"`
	Version         *int       `json:"version" db:"version"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
}

