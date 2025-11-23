package repositories

import (
	"database/sql"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// GetUserWallet gets user wallet by user_id
func (r *WalletRepository) GetUserWallet(userID int) (*models.Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = ? LIMIT 1`

	var wallet models.Wallet
	err := r.db.QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// UpdateWalletBalance updates wallet balance
func (r *WalletRepository) UpdateWalletBalance(userID int, balance float64) error {
	query := `UPDATE wallets SET balance = ?, updated_at = NOW() WHERE user_id = ?`
	_, err := r.db.Exec(query, balance, userID)
	return err
}

