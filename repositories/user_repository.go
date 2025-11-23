package repositories

import (
	"database/sql"

	"github.com/kytapay/webhook-v2/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID gets user by ID
func (r *UserRepository) GetUserByID(userID int) (*models.User, error) {
	query := `SELECT id, email, role_id FROM users WHERE id = ? LIMIT 1`

	var user models.User
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.RoleID,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

