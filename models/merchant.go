package models

type Merchant struct {
	ID           int     `json:"id" db:"id"`
	UserID       int     `json:"user_id" db:"user_id"`
	BusinessName string  `json:"business_name" db:"business_name"`
	MerchantUUID *string `json:"merchant_uuid" db:"merchant_uuid"`
	SiteURL      *string `json:"site_url" db:"site_url"`
	Status       string  `json:"status" db:"status"`
}

type MerchantPayment struct {
	ID              int        `json:"id" db:"id"`
	MerchantID      *int       `json:"merchant_id" db:"merchant_id"`
	PaymentMethodID *int       `json:"payment_method_id" db:"payment_method_id"`
	GatewayReference *string   `json:"gateway_reference" db:"gateway_reference"`
	OrderNo         *string   `json:"order_no" db:"order_no"`
	UUID            *string   `json:"uuid" db:"uuid"`
	FeeBearer       string    `json:"fee_bearer" db:"fee_bearer"`
	Percentage      float64   `json:"percentage" db:"percentage"`
	ChargePercentage float64  `json:"charge_percentage" db:"charge_percentage"`
	ChargeFixed     float64   `json:"charge_fixed" db:"charge_fixed"`
	Amount          float64   `json:"amount" db:"amount"`
	Total           float64   `json:"total" db:"total"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       *string   `json:"created_at" db:"created_at"`
	UpdatedAt       *string   `json:"updated_at" db:"updated_at"`
}

